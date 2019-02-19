package core

import (
	"testing"

	"github.com/Gonzih/wasm-mk2/dom"
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
	dom.RegisterMockTemplate("app-root", `<mydiv></mydiv>`)
	dom.RegisterMockTemplate("mydiv-template", `<div :class="Input"></div>`)
	Component(&MyDiv{}, "mydiv", "mydiv-template")

	app := New()
	err := app.Mount("app-root")

	assert.Nil(t, err)
}
