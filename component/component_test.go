package component

import (
	"testing"

	"github.com/Gonzih/wasm-mk2/event"
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

func (c *MyDiv) HandleClick(e *event.Event) {
	c.Counter++
}

func (c *MyDiv) OtherMethod() {
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

	getter, ok := wrapper.Getter("Counter")
	assert.True(t, ok)
	assert.Equal(t, 10, getter())
}

func TestGettersAndSetters(t *testing.T) {
	w, err := Wasmify(&MyDiv{})
	assert.Nil(t, err)

	wrapper, err := w.Instance()
	assert.Nil(t, err)

	setter, ok := wrapper.Setter("Label")
	assert.True(t, ok)

	getter, ok := wrapper.Getter("Label")
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

	_, ok := wrapper.IsAProp("label")
	assert.True(t, ok)
	_, ok = wrapper.IsAProp("counter")
	assert.False(t, ok)
	assert.NotEqual(t, "", wrapper.UUID())
}

func TestHandlerLookup(t *testing.T) {
	in := &MyDiv{}
	w, err := Wasmify(in)
	assert.Nil(t, err)

	wrapper, err := w.Instance()
	assert.Nil(t, err)

	assert.True(t, wrapper.IsAHandler("HandleClick"))
	assert.False(t, wrapper.IsAHandler("Init"))
	assert.False(t, wrapper.IsAHandler("OtherMethod"))
}

func TestHandler(t *testing.T) {
	in := &MyDiv{}
	w, err := Wasmify(in)
	assert.Nil(t, err)

	wrapper, err := w.Instance()
	assert.Nil(t, err)

	getter, ok := wrapper.Getter("Counter")
	assert.True(t, ok)

	assert.Equal(t, 10, getter())

	handler, ok := wrapper.Handler("HandleClick")
	assert.True(t, ok)
	handler(&event.Event{})

	assert.Equal(t, 11, getter())
}
