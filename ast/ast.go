package ast

import (
	"fmt"
	"strings"
)

const indententionCharacter = "  "

// Node represents node interface
type Node interface {
	nodeType()
	Tag() string
	Attributes() []Attribute
	Children() []Node
	String() string
	indentedString(int) string
}

// Attribute represents singe html attribute
type Attribute struct {
	Name  string
	Value string
}

// Root represents AST tree root node
type Root struct {
	HTMLChildren []Node
}

func (rt *Root) nodeType()               {}
func (rt *Root) Tag() string             { return "root" }
func (rt *Root) Attributes() []Attribute { return []Attribute{} }
func (rt *Root) Children() []Node        { return rt.HTMLChildren }
func (rt *Root) String() string {
	var out strings.Builder

	out.WriteString("\n<root>\n")
	for _, ch := range rt.Children() {
		out.WriteString(ch.indentedString(1))
	}
	out.WriteString("</root>\n")

	return out.String()
}

// Element represents simple html element
type Element struct {
	HTMLTag        string
	HTMLAttributes []Attribute
	HTMLChildren   []Node
}

func (el *Element) nodeType()               {}
func (el *Element) Tag() string             { return el.HTMLTag }
func (el *Element) Attributes() []Attribute { return el.HTMLAttributes }
func (el *Element) Children() []Node        { return el.HTMLChildren }
func (el *Element) String() string {
	return el.indentedString(0)
}

func (el *Element) indentedString(level int) string {
	var out strings.Builder

	ident := strings.Repeat(indententionCharacter, level)

	out.WriteString(ident)
	out.WriteString("<")
	out.WriteString(el.Tag())
	for _, attr := range el.Attributes() {
		out.WriteString(fmt.Sprintf(` %s="%s"`, attr.Name, attr.Value))
	}
	out.WriteString(">\n")

	for _, ch := range el.Children() {
		out.WriteString(ch.indentedString(level + 1))
	}

	out.WriteString(ident)
	out.WriteString("</")
	out.WriteString(el.Tag())
	out.WriteString(">\n")

	return out.String()
}

// Component represents user defined component
// type Component struct {
// }

// func (co *Component) nodeType()               {}
// func (co *Component) Tag() string             { return "component" }
// func (co *Component) Attributes() []Attribute { return []Attribute{} }
