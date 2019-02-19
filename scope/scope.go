package scope

import (
	"github.com/Gonzih/wasm-mk2/component"
	"github.com/Gonzih/wasm-mk2/event"
)

type Scope struct {
	Parent  *Scope
	Wrapper *component.Wrapper
}

func New(w *component.Wrapper, parent *Scope) *Scope {
	return &Scope{Parent: parent, Wrapper: w}
}

func Empty() *Scope {
	return &Scope{}
}

func (s *Scope) Getter(name string) (func() interface{}, bool) {
	if s.Wrapper == nil {
		return nil, false
	}

	getter, ok := s.Wrapper.Getter(name)

	if !ok {
		if s.Parent != nil {
			getter, ok = s.Parent.Getter(name)
			return getter, ok
		}
	}

	return getter, ok
}

func (s *Scope) Handler(name string) (func(*event.Event), bool) {
	if s.Wrapper == nil {
		return nil, false
	}

	handler, ok := s.Wrapper.Handler(name)

	if !ok {
		if s.Parent != nil {
			handler, ok = s.Parent.Handler(name)
			return handler, ok
		}
	}

	return handler, ok
}
