package registry

import (
	"log"

	"github.com/Gonzih/wasm-mk2/component"
)

var registry = make(map[string]*component.Wrapper, 0)
var templateRegistry = make(map[string]string, 0)

func Register(name string, wrapper *component.Wrapper) {
	registry[name] = wrapper
}

func RegisterTemplate(name, templateID string) {
	templateRegistry[name] = templateID
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

func TemplateID(name string) (string, bool) {
	templateID, ok := templateRegistry[name]
	return templateID, ok
}
