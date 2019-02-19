package registry

import (
	"testing"

	"github.com/Gonzih/wasm-mk2/component"
	"github.com/stretchr/testify/assert"
)

type MyDiv struct {
	Input   string `wasm:"prop"`
	Counter int    `wasm:"state"`
}

func (c *MyDiv) Init() error {
	c.Counter = 11
	c.Input = "MyDiv"
	return nil
}

func TestBasicExists(t *testing.T) {
	w, err := component.Wasmify(&MyDiv{})
	assert.Nil(t, err)

	Register("MyDiv", w)

	assert.True(t, Exists("MyDiv"))
}

func TestBasicInstance(t *testing.T) {
	w, err := component.Wasmify(&MyDiv{})
	assert.Nil(t, err)

	Register("MyDiv", w)

	wrapper, ok := Instance("MyDiv")
	assert.True(t, ok)

	getter, ok := wrapper.Getter("Input")
	assert.True(t, ok)
	assert.Equal(t, "MyDiv", getter())
}

func TestBasicTemplateID(t *testing.T) {
	RegisterTemplate("mydiv", "mydiv-template")
	id, ok := TemplateID("mydiv")
	assert.True(t, ok)
	assert.Equal(t, "mydiv-template", id)
}
