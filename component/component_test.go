package component

import "testing"

func TestBasic(t *testing.T) {
	Wasmify(&MyDiv{}, "my-div", "#my-div-template")
}

type MyDiv struct {
	Input   string `wasm:"prop"`
	Counter int    `wasm:"state"`
}

func (c *MyDiv) HandleClick(e Event) {
	c.Counter++
}
