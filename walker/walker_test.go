package walker

import (
	"strings"
	"testing"

	"github.com/Gonzih/wasm-mk2/component"
	"github.com/Gonzih/wasm-mk2/parser"
	"github.com/Gonzih/wasm-mk2/registry"
	"github.com/Gonzih/wasm-mk2/tree"
	"github.com/stretchr/testify/assert"

	"golang.org/x/net/html"
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

func walkString(t *testing.T, input string) *Walker {
	r := strings.NewReader(input)
	z := html.NewTokenizer(r)
	p := parser.New(z)
	w := New(p)

	checkWalkErrors(t, w)

	return w
}

func checkWalkErrors(t *testing.T, w *Walker) {
	assert.Len(t, w.Errors(), 0)

	for _, e := range w.Errors() {
		t.Error(e)
	}

	if len(w.Errors()) > 0 {
		t.FailNow()
	}
}

func TestBasic(t *testing.T) {
	input := `<div></div>`
	w := walkString(t, input)
	cmp := w.WalkAST()
	checkWalkErrors(t, w)

	assert.Len(t, cmp, 1)
	assert.Equal(t, "div", cmp[0].Tag())
}

func TestNested(t *testing.T) {
	input := `<div><p></p><a></a></div>`
	w := walkString(t, input)
	cmp := w.WalkAST()
	checkWalkErrors(t, w)

	assert.Len(t, cmp, 1)
	assert.Equal(t, "div", cmp[0].Tag())
	assert.Equal(t, "p", cmp[0].Children()[0].Tag())
	assert.Equal(t, "a", cmp[0].Children()[1].Tag())
}

func TestSimpleComponent(t *testing.T) {
	wrapper, err := component.Wasmify(&MyDiv{})
	assert.Nil(t, err)

	registry.Register("mydiv", wrapper)

	input := `<mydiv></mydiv>`
	w := walkString(t, input)
	cmp := w.WalkAST()
	checkWalkErrors(t, w)

	assert.Len(t, cmp, 1)
	assert.IsType(t, &tree.ComponentNode{}, cmp[0])
	assert.Equal(t, "mydiv", cmp[0].Tag())
}

func TestSimpleComponentWithStaticProp(t *testing.T) {
	wrapper, err := component.Wasmify(&MyDiv{})
	assert.Nil(t, err)

	registry.Register("mydiv", wrapper)

	input := `<mydiv class="myclass"></mydiv>`
	w := walkString(t, input)
	cmp := w.WalkAST()
	checkWalkErrors(t, w)

	assert.Len(t, cmp, 1)
	assert.Equal(t, "class", cmp[0].Props()[0].Key())
	assert.Equal(t, "myclass", cmp[0].Props()[0].Value())
}

func TestSimpleComponentWithDynamicProp(t *testing.T) {
	wrapper, err := component.Wasmify(&MyDiv{})
	assert.Nil(t, err)

	registry.Register("mydiv", wrapper)

	input := `<mydiv :id="Input"></mydiv>`
	w := walkString(t, input)
	cmp := w.WalkAST()
	checkWalkErrors(t, w)

	assert.Len(t, cmp, 1)
	assert.Equal(t, "id", cmp[0].Props()[0].Key())
	assert.Equal(t, "MyDynamicInput", cmp[0].Props()[0].Value())
}

func TestSimpleComponentWithDynamicPropAndNestedScopes(t *testing.T) {
	wrapper, err := component.Wasmify(&MyDiv{})
	assert.Nil(t, err)
	registry.Register("mydiv", wrapper)

	wrapper, err = component.Wasmify(&EmptyDiv{})
	assert.Nil(t, err)
	registry.Register("empty-div", wrapper)

	input := `<mydiv><empty-div :class="Input"></empty-div></mydiv>`
	w := walkString(t, input)
	cmp := w.WalkAST()
	checkWalkErrors(t, w)

	assert.Len(t, cmp, 1)
	assert.Equal(t, "class", cmp[0].Children()[0].Props()[0].Key())
	assert.Equal(t, "MyDynamicInput", cmp[0].Children()[0].Props()[0].Value())
}

func TestSimpleComponentWithDynamicPropPassing(t *testing.T) {
	wrapper, err := component.Wasmify(&MyDiv{})
	assert.Nil(t, err)
	registry.Register("mydiv", wrapper)

	wrapper, err = component.Wasmify(&EmptyDiv{})
	assert.Nil(t, err)
	registry.Register("empty-div", wrapper)

	input := `<mydiv><empty-div :data="Input"></empty-div></mydiv>`
	w := walkString(t, input)
	cmp := w.WalkAST()
	checkWalkErrors(t, w)

	node := cmp[0].Children()[0]
	cmpn, ok := node.(*tree.ComponentNode)
	assert.True(t, ok)

	getter, ok := cmpn.Instance.Getter("data")
	assert.True(t, ok)
	assert.Equal(t, "MyDynamicInput", getter())
}
