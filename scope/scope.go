package scope

import "github.com/Gonzih/wasm-mk2/component"

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
