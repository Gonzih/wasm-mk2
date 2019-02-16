package parser

import (
	"github.com/Gonzih/wasm-mk2/ast"
	"golang.org/x/net/html"
)

// Parser represents our parser state
type Parser struct {
	currToken html.Token
	peekToken html.Token
	tokenizer *html.Tokenizer
	errors    []string
}

// New creates new Parser out of html tokenizer
func New(z *html.Tokenizer) *Parser {
	p := &Parser{tokenizer: z}
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	_ = p.tokenizer.Next()
	p.currToken = p.peekToken
	p.peekToken = p.tokenizer.Token()
}

// Errors returns internal parser errors slice
func (p *Parser) Errors() []string {
	return p.errors
}

// ParseTree parses entire tree and returns ast Root
func (p *Parser) ParseTree() *ast.Root {
	root := &ast.Root{}

	for !p.curTokenIs(html.ErrorToken) {
		node := p.parseNode()
		if node != nil {
			root.Children = append(root.Children, node)
		}
		p.nextToken()
	}

	return root
}

func (p *Parser) curTokenIs(tt html.TokenType) bool {
	return p.currToken.Type == tt
}

func (p *Parser) peekTokenIs(tt html.TokenType) bool {
	return p.peekToken.Type == tt
}

func (p *Parser) parseNode() ast.Node {
	if p.curTokenIs(html.StartTagToken) && p.peekTokenIs(html.EndTagToken) {
		attrs := []ast.Attribute{}
		for _, attr := range p.currToken.Attr {
			at := ast.Attribute{
				Name:  attr.Key,
				Value: attr.Val,
			}

			attrs = append(attrs, at)
		}

		node := &ast.Element{
			HTMLTag:        p.currToken.Data,
			HTMLAttributes: attrs,
		}

		return node
	}

	return nil
}
