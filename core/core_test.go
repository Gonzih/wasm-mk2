package core

import (
	"testing"

	"github.com/Gonzih/wasm-mk2/dom"
	"github.com/Gonzih/wasm-mk2/event"
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

func (c *MyDiv) HandleClick(e *event.Event) {
	c.Counter = 200
}

func TestBasic(t *testing.T) {
	dom.RegisterMockTemplate("app-root", `<mydiv></mydiv>`)
	dom.RegisterMockTemplate("mydiv-template", `<div :class="Input" :data-id="Counter" @click="HandleClick"></div>`)
	Component(&MyDiv{}, "mydiv", "mydiv-template")

	app := New()
	err := app.Mount("app-root")

	assert.Len(t, app.Components, 1)

	child := app.Components[0].Children()[0]
	assert.True(t, child.Handle("click", &event.Event{}))

	for _, prop := range child.Props() {
		if prop.Key() == "data-id" {
			assert.Equal(t, "200", prop.Value())
		}
	}

	assert.Nil(t, err)
}
