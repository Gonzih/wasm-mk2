package walker

import (
	"github.com/Gonzih/wasm-mk2/ast"
	"github.com/Gonzih/wasm-mk2/parser"
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

func (w *Walker) WalkAST() []*tree.Component {
	return w.walkComponent(w.root.Children())
}

func (w *Walker) walkComponent(nodes []ast.Node) []*tree.Component {
	cmps := make([]*tree.Component, 0)

	for _, astNode := range nodes {
		cmp := &tree.Component{
			Tag:      astNode.Tag(),
			Children: w.walkComponent(astNode.Children()),
		}
		cmps = append(cmps, cmp)
	}

	return cmps
}
