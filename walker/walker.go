package walker

import (
	"log"
	"strings"

	"github.com/Gonzih/wasm-mk2/ast"
	"github.com/Gonzih/wasm-mk2/dom"
	"github.com/Gonzih/wasm-mk2/parser"
	"github.com/Gonzih/wasm-mk2/registry"
	"github.com/Gonzih/wasm-mk2/scope"
	"github.com/Gonzih/wasm-mk2/tree"
	"golang.org/x/net/html"
)

type Walker struct {
	parser *parser.Parser
	root   *ast.Root
	errors []string
}

func NewByID(templateID string) *Walker {
	input := dom.New().TemplateContent(templateID)
	r := strings.NewReader(input)
	z := html.NewTokenizer(r)
	p := parser.New(z)

	return New(p)
}

func New(p *parser.Parser) *Walker {
	w := &Walker{
		parser: p,
		root:   p.ParseTree(),
		errors: p.Errors(),
	}

	return w
}

func (w *Walker) Errors() []string {
	return w.errors
}

func (w *Walker) WalkAST(s *scope.Scope) []tree.Node {
	components := w.walkComponent(w.root.Children(), s)

	for _, cmp := range components {
		cmp.Notify()
	}

	return components
}
func (w *Walker) convertHandlers(attrs []ast.Attribute, scope *scope.Scope) []*tree.Handler {
	result := make([]*tree.Handler, 0)

	for _, attr := range attrs {
		k := attr.Name
		v := attr.Value

		if strings.HasPrefix(k, "@") {
			if scope != nil {
				k = strings.Replace(k, "@", "", 1)
				handler, ok := scope.Handler(v)
				if !ok {
					log.Fatalf("Could not find handler %s", k)
				}

				result = append(result, &tree.Handler{
					Key: k,
					F:   handler,
				})
			}
		}
	}

	return result
}

func (w *Walker) convertProperties(attrs []ast.Attribute, scope *scope.Scope) []tree.Attribute {
	result := make([]tree.Attribute, 0)

	for _, attr := range attrs {
		var prop tree.Attribute
		k := attr.Name
		v := attr.Value

		if strings.HasPrefix(k, ":") {
			if scope != nil {
				k = strings.Replace(k, ":", "", 1)
				prop = newDynamicAttribute(k, v, scope)
			} else {
				log.Print("Instance was nil")
			}
		} else {
			prop = newStaticAttribute(k, v)
		}

		result = append(result, prop)
	}

	return result
}

func (w *Walker) walkComponent(nodes []ast.Node, parentScope *scope.Scope) []tree.Node {
	cmps := make([]tree.Node, 0)

	for _, astNode := range nodes {
		var cmp tree.Node
		tag := astNode.Tag()
		instance, isComponent := registry.Instance(tag)
		currScope := parentScope

		if isComponent {
			currScope = scope.New(instance, parentScope)

			templateID, ok := registry.TemplateID(tag)
			if !ok {
				log.Fatalf("Could not find template for %s", tag)
			}

			innerWalker := NewByID(templateID)
			ast := innerWalker.WalkAST(currScope)
			w.errors = append(w.errors, innerWalker.Errors()...)

			cmp = &tree.ComponentNode{
				NodeTag:      tag,
				NodeChildren: ast,
				NodeBody:     w.walkComponent(astNode.Children(), currScope),
				NodeProps:    w.convertProperties(astNode.Attributes(), currScope),
				NodeHandlers: w.convertHandlers(astNode.Attributes(), currScope),
				Instance:     instance,
			}
		} else {
			cmp = &tree.HTMLNode{
				NodeTag:      tag,
				NodeChildren: w.walkComponent(astNode.Children(), currScope),
				NodeProps:    w.convertProperties(astNode.Attributes(), currScope),
				NodeHandlers: w.convertHandlers(astNode.Attributes(), currScope),
			}
		}

		cmps = append(cmps, cmp)
	}

	return cmps
}

func newDynamicAttribute(k, v string, scope *scope.Scope) tree.Attribute {
	var f func() string

	getter, ok := scope.Getter(v)
	if !ok {
		log.Fatalf("Could not find getter for %s", v)
		f = func() string {
			return v
		}
	} else {
		f = func() string {
			raw := getter()
			s, ok := raw.(string)
			if !ok {
				log.Printf("Could not convert %v in to string", raw)
			}

			return s
		}
	}

	propName, isAProp := scope.Wrapper.IsAProp(k)

	if isAProp {
		setter, ok := scope.Wrapper.Setter(propName)
		if !ok {
			log.Fatalf("Could not find setter for %s with name %s", k, propName)
		}
		return &tree.LinkedAttribute{
			K: k,
			F: f,
			Sync: func() {
				setter(f())
			},
		}
	}

	return &tree.DynamicAttribute{
		K: k,
		F: f,
	}
}

func newStaticAttribute(k, v string) tree.Attribute {
	return &tree.StaticAttribute{
		K: k,
		V: v,
	}
}
