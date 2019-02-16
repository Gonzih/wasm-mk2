package ast

// Node represents node interface
type Node interface {
	nodeType()
	Tag() string
	Attributes() []Attribute
}

// Attribute represents singe html attribute
type Attribute struct {
	Name  string
	Value string
}

// Root represents AST tree root node
type Root struct {
	Children []Node
}

func (rt *Root) nodeType()               {}
func (rt *Root) Tag() string             { return "root" }
func (rt *Root) Attributes() []Attribute { return []Attribute{} }

// Element represents simple html element
type Element struct {
	HTMLTag        string
	HTMLAttributes []Attribute
}

func (el *Element) nodeType()               {}
func (el *Element) Tag() string             { return el.HTMLTag }
func (el *Element) Attributes() []Attribute { return el.HTMLAttributes }

// Component represents user defined component
type Component struct {
}

func (co *Component) nodeType()               {}
func (co *Component) Tag() string             { return "component" }
func (co *Component) Attributes() []Attribute { return []Attribute{} }
