package walker

import (
	"github.com/Gonzih/wasm-mk2/ast"
	"github.com/Gonzih/wasm-mk2/component"
	"github.com/Gonzih/wasm-mk2/parser"
)

type Walker struct {
	parser *parser.Parser
	root   *ast.Root
}

func New(p *parser.Parser) *Walker {
	w := &Walker{parser: p, root: p.ParseTree()}

	return w
}

func (w *Walker) GenerateComponentTree() *component.Component {

}
