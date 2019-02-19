package core

import (
	"testing"

	"github.com/Gonzih/wasm-mk2/component"
	"github.com/Gonzih/wasm-mk2/dom"
	"github.com/Gonzih/wasm-mk2/registry"
	"github.com/stretchr/testify/assert"
)

type EmptyDiv struct {
	Data string `wasm:"prop"`
}

func (c *EmptyDiv) Init() error { return nil }

type MyDiv struct {
	Input   string `wasm:"prop"`
	Counter int    `wasm:"state"`
}

func (c *MyDiv) Init() error {
	c.Counter = 11
	c.Input = "MyDynamicInput"
	return nil
}

func TestBasic(t *testing.T) {
	wrapper, err := component.Wasmify(&MyDiv{})
	assert.Nil(t, err)
	registry.Register("mydiv", wrapper)

	input := `<mydiv></mydiv>`
	dom.RegisterMockTemplate("mydiv", input)
}
