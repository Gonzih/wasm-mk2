package component

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MyDiv struct {
	Input   string `wasm:"prop"`
	Counter int    `wasm:"state"`
	Label   int    `wasm:"prop,state"`
}

func (c *MyDiv) Init() error {
	c.Counter = 10
	return nil
}

func (c *MyDiv) HandleClick(e Event) {
	c.Counter++
}

func TestBasic(t *testing.T) {
	_, err := Wasmify(&MyDiv{})
	assert.Nil(t, err)
}

func TestInstance(t *testing.T) {
	initial := &MyDiv{Counter: 5}
	wrapper, err := Wasmify(initial)
	assert.Nil(t, err)

	newWrapper, err := wrapper.Instance()
	assert.Nil(t, err)

	md, ok := newWrapper.instance.(*MyDiv)

	assert.True(t, ok)

	if ok {
		assert.Equal(t, 10, md.Counter)
		assert.Equal(t, 5, initial.Counter)
	}
}

func TestGetters(t *testing.T) {
	w, err := Wasmify(&MyDiv{})
	assert.Nil(t, err)

	wrapper, err := w.Instance()
	assert.Nil(t, err)

	getter, ok := wrapper.getters["Counter"]
	assert.True(t, ok)
	assert.Equal(t, 10, getter())
}

func TestGettersAndSetters(t *testing.T) {
	w, err := Wasmify(&MyDiv{})
	assert.Nil(t, err)

	wrapper, err := w.Instance()
	assert.Nil(t, err)

	setter, ok := wrapper.setters["Label"]
	assert.True(t, ok)

	getter, ok := wrapper.getters["Label"]
	assert.True(t, ok)

	err = setter(23)
	assert.Nil(t, err)
	assert.Equal(t, 23, getter())
}

func TestProps(t *testing.T) {
	w, err := Wasmify(&MyDiv{})
	assert.Nil(t, err)

	wrapper, err := w.Instance()
	assert.Nil(t, err)

	assert.True(t, wrapper.IsAProp("Label"))
	assert.False(t, wrapper.IsAProp("Counter"))
}
