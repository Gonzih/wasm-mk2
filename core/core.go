package core

import (
	"github.com/Gonzih/wasm-mk2/component"
	"github.com/Gonzih/wasm-mk2/registry"
)

func Component(strukt interface{}, name, templateID string) {
	wrapper := component.Wasmify(strukt)
	registry.Register(name, wrapper)
}
