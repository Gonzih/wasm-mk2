package component

import (
	"testing"

	"github.com/Gonzih/wasm-mk2/event"
	"github.com/stretchr/testify/require"
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
	require.Nil(t, err)
}

func TestInstance(t *testing.T) {
	initial := &MyDiv{Counter: 5}
	wrapper, err := Wasmify(initial)
	require.Nil(t, err)

	newWrapper, err := wrapper.Instance()
	require.Nil(t, err)

	md, ok := newWrapper.instance.(*MyDiv)

	require.True(t, ok)

	if ok {
		require.Equal(t, 10, md.Counter)
		require.Equal(t, 5, initial.Counter)
	}
}

func TestGetters(t *testing.T) {
	w, err := Wasmify(&MyDiv{})
	require.Nil(t, err)

	wrapper, err := w.Instance()
	require.Nil(t, err)

	getter, ok := wrapper.Getter("Counter")
	require.True(t, ok)
	require.Equal(t, 10, getter())
}

func TestGettersAndSetters(t *testing.T) {
	w, err := Wasmify(&MyDiv{})
	require.Nil(t, err)

	wrapper, err := w.Instance()
	require.Nil(t, err)

	setter, ok := wrapper.Setter("Label")
	require.True(t, ok)

	getter, ok := wrapper.Getter("Label")
	require.True(t, ok)

	err = setter(23)
	require.Nil(t, err)
	require.Equal(t, 23, getter())
}

func TestProps(t *testing.T) {
	w, err := Wasmify(&MyDiv{})
	require.Nil(t, err)

	wrapper, err := w.Instance()
	require.Nil(t, err)

	_, ok := wrapper.IsAProp("label")
	require.True(t, ok)
	_, ok = wrapper.IsAProp("counter")
	require.False(t, ok)
	require.NotEqual(t, "", wrapper.UUID())
}

func TestHandlerLookup(t *testing.T) {
	in := &MyDiv{}
	w, err := Wasmify(in)
	require.Nil(t, err)

	wrapper, err := w.Instance()
	require.Nil(t, err)

	require.True(t, wrapper.IsAHandler("HandleClick"))
	require.False(t, wrapper.IsAHandler("Init"))
	require.False(t, wrapper.IsAHandler("OtherMethod"))
}

func TestHandler(t *testing.T) {
	in := &MyDiv{}
	w, err := Wasmify(in)
	require.Nil(t, err)

	wrapper, err := w.Instance()
	require.Nil(t, err)

	getter, ok := wrapper.Getter("Counter")
	require.True(t, ok)

	require.Equal(t, 10, getter())

	handler, ok := wrapper.Handler("HandleClick")
	require.True(t, ok)
	handler(&event.Event{})

	require.Equal(t, 11, getter())
}
