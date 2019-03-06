package scope

import (
	"testing"

	"github.com/Gonzih/wasm-mk2/component"
	"github.com/Gonzih/wasm-mk2/event"
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

type MyDivTwo struct {
	Input string `wasm:"prop"`
	Num   int    `wasm:"state"`
}

func (c *MyDivTwo) Init() error {
	c.Num = 99
	c.Input = "MyDivTwo"
	return nil
}

func (c *MyDivTwo) HandleClick(e *event.Event) {
	c.Num = 1999
}

func TestBasicLookup(t *testing.T) {
	w, err := component.Wasmify(&MyDiv{})
	require.Nil(t, err)

	wrapper, err := w.Instance()
	require.Nil(t, err)

	s := New(wrapper, nil)
	getter, ok := s.Getter("Counter")
	require.True(t, ok)
	require.Equal(t, 11, getter())

	getter, ok = s.Getter("Input")
	require.True(t, ok)
	require.Equal(t, "MyDiv", getter())
}

func TestRecursiveLookup(t *testing.T) {
	w, err := component.Wasmify(&MyDivTwo{})
	require.Nil(t, err)
	wrapper, err := w.Instance()
	require.Nil(t, err)
	sParent := New(wrapper, nil)

	w, err = component.Wasmify(&MyDiv{})
	require.Nil(t, err)
	wrapper, err = w.Instance()
	require.Nil(t, err)
	s := New(wrapper, sParent)

	handler, ok := s.Handler("HandleClick")
	require.True(t, ok)
	handler(&event.Event{})

	getter, ok := s.Getter("Num")
	require.True(t, ok)
	require.Equal(t, 1999, getter())
}
