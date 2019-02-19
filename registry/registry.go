package registry

import (
	"log"

	"github.com/Gonzih/wasm-mk2/component"
)

var registry = make(map[string]*component.Wrapper, 0)

func Register(name string, wrapper *component.Wrapper) {
	registry[name] = wrapper
}

func Exists(name string) bool {
	_, ok := registry[name]
	return ok
}

func Instance(name string) (*component.Wrapper, bool) {
	w, ok := registry[name]
	if !ok {
		return nil, ok
	}

	instance, err := w.Instance()
	if err != nil {
		log.Printf("Error creating instance: %s", err)
		return nil, false
	}

	return instance, true
}
