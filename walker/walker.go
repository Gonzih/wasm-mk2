package walker

import (
	"log"
	"strings"

	"github.com/Gonzih/wasm-mk2/ast"
	"github.com/Gonzih/wasm-mk2/component"
	"github.com/Gonzih/wasm-mk2/parser"
	"github.com/Gonzih/wasm-mk2/registry"
	"github.com/Gonzih/wasm-mk2/tree"
)

type Walker struct {
	parser *parser.Parser
	root   *ast.Root
	errors []string
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

func (w *Walker) WalkAST() []tree.Node {
	return w.walkComponent(w.root.Children())
}

func (w *Walker) convertProperties(attrs []ast.Attribute, instance *component.Wrapper) []tree.Attribute {
	result := make([]tree.Attribute, 0)

	for _, attr := range attrs {
		var prop tree.Attribute
		k := attr.Name
		v := attr.Value

		if strings.HasPrefix(k, ":") {
			if instance != nil {
				k = strings.Replace(k, ":", "", 1)

				var f func() string

				getter, ok := instance.Getter(v)
				if !ok {
					log.Printf("Could not find getter for %s", v)
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

				prop = &tree.DynamicAttribute{
					K: k,
					F: f,
				}
			} else {
				log.Print("Instance was nil")
			}
		} else {
			prop = &tree.StaticAttribute{
				K: k,
				V: v,
			}
		}

		result = append(result, prop)
	}

	return result
}

func (w *Walker) walkComponent(nodes []ast.Node) []tree.Node {
	cmps := make([]tree.Node, 0)

	for _, astNode := range nodes {
		var cmp tree.Node
		tag := astNode.Tag()
		children := w.walkComponent(astNode.Children())
		instance, isComponent := registry.Instance(tag)
		props := w.convertProperties(astNode.Attributes(), instance)

		if isComponent {
			cmp = &tree.ComponentNode{
				NodeTag:      tag,
				NodeChildren: children,
				NodeProps:    props,
				Instance:     instance,
			}
		} else {
			cmp = &tree.HTMLNode{
				NodeTag:      tag,
				NodeChildren: children,
				NodeProps:    props,
			}
		}
		cmps = append(cmps, cmp)
	}

	return cmps
}
