package tree

import (
	"github.com/Gonzih/wasm-mk2/component"
	"github.com/Gonzih/wasm-mk2/event"
)

type Node interface {
	Tag() string
	Children() []Node
	Body() []Node
	Props() []Attribute
	Refresh()
	Notify()
	Handle(string, *event.Event) bool
}

type Attribute interface {
	Key() string
	Value() string
	Refresh()
}

type StaticAttribute struct {
	K string
	V string
}

func (p *StaticAttribute) Key() string   { return p.K }
func (p *StaticAttribute) Value() string { return p.V }
func (p *StaticAttribute) Refresh()      {}

type DynamicAttribute struct {
	K string
	F func() string
}

func (p *DynamicAttribute) Key() string   { return p.K }
func (p *DynamicAttribute) Value() string { return p.F() }
func (p *DynamicAttribute) Refresh()      {}

type LinkedAttribute struct {
	K    string
	F    func() string
	Sync func()
}

func (p *LinkedAttribute) Key() string   { return p.K }
func (p *LinkedAttribute) Value() string { return p.F() }
func (p *LinkedAttribute) Refresh()      { p.Sync() }

type HTMLNode struct {
	NodeHandlers []*Handler
	NodeTag      string
	NodeChildren []Node
	NodeProps    []Attribute
}

func (n *HTMLNode) Tag() string        { return n.NodeTag }
func (n *HTMLNode) Children() []Node   { return n.NodeChildren }
func (n *HTMLNode) Body() []Node       { return []Node{} }
func (n *HTMLNode) Props() []Attribute { return n.NodeProps }
func (n *HTMLNode) Refresh()           {}
func (n *HTMLNode) Notify()            {}

func (n *HTMLNode) Handle(name string, e *event.Event) bool {
	for _, handle := range n.NodeHandlers {
		if handle.Key == name {
			handle.F(e)
			return true
		}
	}

	return false
}

type ComponentNode struct {
	NodeHandlers []*Handler
	NodeTag      string
	NodeChildren []Node
	NodeBody     []Node
	NodeProps    []Attribute
	Instance     *component.Wrapper
}

func (n *ComponentNode) Notify() {
	for _, sub := range n.NodeChildren {
		sub.Refresh()
	}
	for _, sub := range n.NodeBody {
		sub.Refresh()
	}
}

func (n *ComponentNode) Refresh() {
	for _, prop := range n.NodeProps {
		prop.Refresh()
	}
	n.Notify()
}

func (n *ComponentNode) Handle(name string, e *event.Event) bool {
	for _, handle := range n.NodeHandlers {
		if handle.Key == name {
			handle.F(e)
			return true
		}
	}

	return false
}

func (n *ComponentNode) Tag() string        { return n.NodeTag }
func (n *ComponentNode) Children() []Node   { return n.NodeChildren }
func (n *ComponentNode) Body() []Node       { return n.NodeBody }
func (n *ComponentNode) Props() []Attribute { return n.NodeProps }

type Handler struct {
	Key string
	F   func(*event.Event)
}
