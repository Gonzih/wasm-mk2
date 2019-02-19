package scope

import (
	"testing"

	"github.com/Gonzih/wasm-mk2/component"
	"github.com/Gonzih/wasm-mk2/event"
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
	assert.Nil(t, err)

	wrapper, err := w.Instance()
	assert.Nil(t, err)

	s := New(wrapper, nil)
	getter, ok := s.Getter("Counter")
	assert.True(t, ok)
	assert.Equal(t, 11, getter())

	getter, ok = s.Getter("Input")
	assert.True(t, ok)
	assert.Equal(t, "MyDiv", getter())
}

func TestRecursiveLookup(t *testing.T) {
	w, err := component.Wasmify(&MyDivTwo{})
	assert.Nil(t, err)
	wrapper, err := w.Instance()
	assert.Nil(t, err)
	sParent := New(wrapper, nil)

	w, err = component.Wasmify(&MyDiv{})
	assert.Nil(t, err)
	wrapper, err = w.Instance()
	assert.Nil(t, err)
	s := New(wrapper, sParent)

	handler, ok := s.Handler("HandleClick")
	assert.True(t, ok)
	handler(&event.Event{})

	getter, ok := s.Getter("Num")
	assert.True(t, ok)
	assert.Equal(t, 1999, getter())
}
