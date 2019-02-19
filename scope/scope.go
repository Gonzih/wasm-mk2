package scope

import "github.com/Gonzih/wasm-mk2/component"

type Scope struct {
	Parent           *Scope
	componentWrapper *component.Wrapper
}

func New(w *component.Wrapper, parent *Scope) *Scope {
	return &Scope{Parent: parent, componentWrapper: w}
}

func (s *Scope) Getter(name string) (func() interface{}, bool) {
	getter, ok := s.componentWrapper.Getter(name)

	if !ok {
		if s.Parent != nil {
			getter, ok = s.Parent.Getter(name)
			return getter, ok
		}
	}

	return getter, ok
}
