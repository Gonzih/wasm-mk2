package walker

import (
	"github.com/Gonzih/wasm-mk2/ast"
	"github.com/Gonzih/wasm-mk2/component"
	"github.com/Gonzih/wasm-mk2/parser"
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

func (w *Walker) WalkAST() []*component.Component {
	return w.walkComponent(w.root.Children())
}

func (w *Walker) walkComponent(nodes []ast.Node) []*component.Component {
	cmps := make([]*component.Component, 0)

	for _, astNode := range nodes {
		cmp := &component.Component{
			Tag: astNode.Tag(),
			Children: w.walkComponent(astNode.Children())
		}
		cmps = append(cmps, cmp)
	}

	return cmps
}
