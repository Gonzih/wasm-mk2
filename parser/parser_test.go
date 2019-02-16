package parser

import (
	"strings"
	"testing"

	"github.com/Gonzih/wasm-mk2/ast"
	"github.com/stretchr/testify/assert"

	"golang.org/x/net/html"
)

func newTestParser(input string) *Parser {
	r := strings.NewReader(input)
	z := html.NewTokenizer(r)
	return New(z)
}

func checkParserErrors(t *testing.T, p *Parser) {
	assert.Len(t, p.Errors(), 0)

	for _, e := range p.Errors() {
		t.Error(e)
	}

	if len(p.Errors()) > 0 {
		t.FailNow()
	}
}

func TestNew(t *testing.T) {
	s := `<div></div>`
	p := newTestParser(s)

	checkParserErrors(t, p)

	assert.Equal(t, html.StartTagToken, p.currToken.Type)
	assert.Equal(t, html.EndTagToken, p.peekToken.Type)
}

func TestParseSingleDiv(t *testing.T) {
	s := `<div></div>`
	p := newTestParser(s)
	root := p.ParseTree()

	checkParserErrors(t, p)

	assert.Len(t, root.Children(), 1)
	assert.IsType(t, &ast.Element{}, root.Children()[0])
	assert.Equal(t, "div", root.Children()[0].Tag())
}

func TestParseSingleDivWithAttributes(t *testing.T) {
	s := `<div class="mydiv"></div>`
	p := newTestParser(s)
	root := p.ParseTree()

	checkParserErrors(t, p)

	assert.Len(t, root.Children(), 1)
	assert.Equal(t, "div", root.Children()[0].Tag())
	assert.Len(t, root.Children()[0].Attributes(), 1)

	at := root.Children()[0].Attributes()[0]
	assert.Equal(t, "class", at.Name)
	assert.Equal(t, "mydiv", at.Value)
}

func TestParseNestedElements(t *testing.T) {
	s := `<div><p></p></div>`
	p := newTestParser(s)
	root := p.ParseTree()

	checkParserErrors(t, p)

	assert.Len(t, root.Children(), 1)
	assert.Equal(t, "div", root.Children()[0].Tag())
	assert.Len(t, root.Children()[0].Children(), 1)

	ch := root.Children()[0].Children()[0]
	assert.Equal(t, "p", ch.Tag())
}

func TestParseMultipleNestedElements(t *testing.T) {
	s := `<div><p></p><a></a></div>`
	p := newTestParser(s)
	root := p.ParseTree()

	checkParserErrors(t, p)

	assert.Len(t, root.Children(), 1)
	assert.Equal(t, "div", root.Children()[0].Tag())
	assert.Len(t, root.Children()[0].Children(), 2)

	ch1 := root.Children()[0].Children()[0]
	assert.Equal(t, "p", ch1.Tag())

	ch2 := root.Children()[0].Children()[1]
	assert.Equal(t, "a", ch2.Tag())
}

func TestParseMultipleNestedElementsInRoot(t *testing.T) {
	s := `<div><p></p><a></a></div><span><div></div></span>`
	p := newTestParser(s)
	root := p.ParseTree()

	checkParserErrors(t, p)

	assert.Len(t, root.Children(), 2)

	assert.Equal(t, "div", root.Children()[0].Tag())
	assert.Len(t, root.Children()[0].Children(), 2)

	assert.Equal(t, "span", root.Children()[1].Tag())
	assert.Len(t, root.Children()[1].Children(), 1)

	ch := root.Children()[0].Children()[0]
	assert.Equal(t, "p", ch.Tag())
	ch = root.Children()[0].Children()[1]
	assert.Equal(t, "a", ch.Tag())

	ch = root.Children()[1].Children()[0]
	assert.Equal(t, "div", ch.Tag())
}

func TestParseSelfClosingElement(t *testing.T) {
	s := `<img src="http://google.com/"/>`
	p := newTestParser(s)
	root := p.ParseTree()

	checkParserErrors(t, p)

	assert.Len(t, root.Children(), 1)

	assert.Equal(t, "img", root.Children()[0].Tag())
	assert.Len(t, root.Children()[0].Children(), 0)
}

func TestParseNestedWithSelfClosingElement(t *testing.T) {
	s := `<div><img src="http://google.com/"/><hr/></div>`
	p := newTestParser(s)
	root := p.ParseTree()

	checkParserErrors(t, p)

	assert.Len(t, root.Children(), 1)
	assert.Len(t, root.Children()[0].Children(), 2)

	assert.Equal(t, "div", root.Children()[0].Tag())
	assert.Equal(t, "img", root.Children()[0].Children()[0].Tag())
	assert.Equal(t, "hr", root.Children()[0].Children()[1].Tag())
}
