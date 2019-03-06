package registry

import (
	"testing"

	"github.com/Gonzih/wasm-mk2/component"
	"github.com/stretchr/testify/require"
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
	require.Nil(t, err)

	Register("MyDiv", w)

	require.True(t, Exists("MyDiv"))
}

func TestBasicInstance(t *testing.T) {
	w, err := component.Wasmify(&MyDiv{})
	require.Nil(t, err)

	Register("MyDiv", w)

	wrapper, ok := Instance("MyDiv")
	require.True(t, ok)

	getter, ok := wrapper.Getter("Input")
	require.True(t, ok)
	require.Equal(t, "MyDiv", getter())
}

func TestBasicTemplateID(t *testing.T) {
	RegisterTemplate("mydiv", "mydiv-template")
	id, ok := TemplateID("mydiv")
	require.True(t, ok)
	require.Equal(t, "mydiv-template", id)
}
