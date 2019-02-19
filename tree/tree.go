package tree

import "github.com/Gonzih/wasm-mk2/component"

type Node interface {
	Tag() string
	Children() []Node
	Props() []Attribute
}

type Attribute interface {
	Key() string
	Value() string
}

type StaticAttribute struct {
	K string
	V string
}

func (p *StaticAttribute) Key() string   { return p.K }
func (p *StaticAttribute) Value() string { return p.V }

type DynamicAttribute struct {
	K string
	F func() string
}

func (p *DynamicAttribute) Key() string   { return p.K }
func (p *DynamicAttribute) Value() string { return p.F() }

type HTMLNode struct {
	NodeTag      string
	NodeChildren []Node
	NodeProps    []Attribute
}

func (n *HTMLNode) Tag() string        { return n.NodeTag }
func (n *HTMLNode) Children() []Node   { return n.NodeChildren }
func (n *HTMLNode) Props() []Attribute { return n.NodeProps }

type ComponentNode struct {
	NodeTag      string
	NodeChildren []Node
	NodeProps    []Attribute
	Instance     *component.Wrapper
}

func (n *ComponentNode) Tag() string        { return n.NodeTag }
func (n *ComponentNode) Children() []Node   { return n.NodeChildren }
func (n *ComponentNode) Props() []Attribute { return n.NodeProps }
