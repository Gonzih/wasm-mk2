package walker

import (
	"testing"

	"github.com/Gonzih/wasm-mk2/component"
	"github.com/Gonzih/wasm-mk2/dom"
	"github.com/Gonzih/wasm-mk2/event"
	"github.com/Gonzih/wasm-mk2/registry"
	"github.com/Gonzih/wasm-mk2/scope"
	"github.com/Gonzih/wasm-mk2/tree"
	"github.com/stretchr/testify/require"
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
	c.Counter += 6
}

func walkString(t *testing.T, input string) *Walker {
	dom.RegisterMockTemplate("app-root", input)
	w := NewByID("app-root")

	checkWalkErrors(t, w)

	return w
}

func checkWalkErrors(t *testing.T, w *Walker) {
	require.Len(t, w.Errors(), 0)

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
	cmp := w.WalkAST(scope.Empty())
	checkWalkErrors(t, w)

	require.Len(t, cmp, 1)
	require.Equal(t, "div", cmp[0].Tag())
}

func TestNested(t *testing.T) {
	input := `<div><p></p><a></a></div>`
	w := walkString(t, input)
	cmp := w.WalkAST(scope.Empty())
	checkWalkErrors(t, w)

	require.Len(t, cmp, 1)
	require.Equal(t, "div", cmp[0].Tag())
	require.Equal(t, "p", cmp[0].Children()[0].Tag())
	require.Equal(t, "a", cmp[0].Children()[1].Tag())
}

func TestSimpleComponent(t *testing.T) {
	wrapper, err := component.Wasmify(&MyDiv{})
	require.Nil(t, err)

	registry.Register("mydiv", wrapper)
	registry.RegisterTemplate("mydiv", "mydiv-template")
	dom.RegisterMockTemplate("mydiv-template", `<div></div>`)

	input := `<mydiv></mydiv>`
	w := walkString(t, input)
	cmp := w.WalkAST(scope.Empty())
	checkWalkErrors(t, w)

	require.Len(t, cmp, 1)
	require.IsType(t, &tree.ComponentNode{}, cmp[0])
	require.Equal(t, "mydiv", cmp[0].Tag())
}

func TestSimpleComponentWithStaticProp(t *testing.T) {
	wrapper, err := component.Wasmify(&MyDiv{})
	require.Nil(t, err)

	registry.Register("mydiv", wrapper)
	registry.RegisterTemplate("mydiv", "mydiv-template")
	dom.RegisterMockTemplate("mydiv-template", `<div></div>`)

	input := `<mydiv class="myclass"></mydiv>`
	w := walkString(t, input)
	cmp := w.WalkAST(scope.Empty())
	checkWalkErrors(t, w)

	require.Len(t, cmp, 1)
	require.Equal(t, "class", cmp[0].Props()[0].Key())
	require.Equal(t, "myclass", cmp[0].Props()[0].Value())
}

func TestSimpleComponentWithDynamicProp(t *testing.T) {
	wrapper, err := component.Wasmify(&MyDiv{})
	require.Nil(t, err)

	registry.Register("mydiv", wrapper)
	registry.RegisterTemplate("mydiv", "mydiv-template")
	dom.RegisterMockTemplate("mydiv-template", `<div></div>`)

	input := `<mydiv :id="Input"></mydiv>`
	w := walkString(t, input)
	cmp := w.WalkAST(scope.Empty())
	checkWalkErrors(t, w)

	require.Len(t, cmp, 1)
	require.Equal(t, "id", cmp[0].Props()[0].Key())
	require.Equal(t, "MyDynamicInput", cmp[0].Props()[0].Value())
}

func TestSimpleComponentWithDynamicPropAndNestedScopes(t *testing.T) {
	wrapper, err := component.Wasmify(&MyDiv{})
	require.Nil(t, err)
	registry.Register("mydiv", wrapper)
	registry.RegisterTemplate("mydiv", "mydiv-template")
	dom.RegisterMockTemplate("mydiv-template", `<div></div>`)

	wrapper, err = component.Wasmify(&EmptyDiv{})
	require.Nil(t, err)
	registry.Register("empty-div", wrapper)
	registry.RegisterTemplate("empty-div", "empty-div-template")
	dom.RegisterMockTemplate("empty-div-template", `<div></div>`)

	input := `<mydiv><empty-div :class="Input"></empty-div></mydiv>`
	w := walkString(t, input)
	cmp := w.WalkAST(scope.Empty())
	checkWalkErrors(t, w)

	require.Len(t, cmp, 1)
	require.Equal(t, "class", cmp[0].Body()[0].Props()[0].Key())
	require.Equal(t, "MyDynamicInput", cmp[0].Body()[0].Props()[0].Value())
}

func TestSimpleComponentWithDynamicPropPassing(t *testing.T) {
	wrapper, err := component.Wasmify(&MyDiv{})
	require.Nil(t, err)
	registry.Register("mydiv", wrapper)
	registry.RegisterTemplate("mydiv", "mydiv-template")
	dom.RegisterMockTemplate("mydiv-template", `<div></div>`)

	wrapper, err = component.Wasmify(&EmptyDiv{})
	require.Nil(t, err)
	registry.Register("empty-div", wrapper)
	registry.RegisterTemplate("empty-div", "empty-div-template")
	dom.RegisterMockTemplate("empty-div-template", `<div></div>`)

	input := `<mydiv><empty-div :data="Input"></empty-div></mydiv>`
	dom.RegisterMockTemplate("app-root", input)
	w := NewByID("app-root")
	cmp := w.WalkAST(scope.Empty())
	checkWalkErrors(t, w)

	node := cmp[0].Body()[0]
	cmpn, ok := node.(*tree.ComponentNode)
	require.True(t, ok)

	getter, ok := cmpn.Instance.Getter("Data")
	require.True(t, ok)
	require.Equal(t, "MyDynamicInput", getter())
}

func TestSimpleComponentWithChildrenProp(t *testing.T) {
	wrapper, err := component.Wasmify(&MyDiv{})
	require.Nil(t, err)
	registry.Register("mydiv", wrapper)
	registry.RegisterTemplate("mydiv", "mydiv-template")
	dom.RegisterMockTemplate("mydiv-template", `<div :class="Input"></div>`)

	input := `<mydiv></mydiv>`
	dom.RegisterMockTemplate("app-root", input)
	w := NewByID("app-root")
	cmp := w.WalkAST(scope.Empty())
	checkWalkErrors(t, w)

	require.Len(t, cmp, 1)
	require.Len(t, cmp[0].Children(), 1)
	require.Equal(t, "div", cmp[0].Children()[0].Tag())
}

func TestHandlersBasic(t *testing.T) {
	wrapper, err := component.Wasmify(&MyDiv{})
	require.Nil(t, err)
	registry.Register("mydiv", wrapper)
	registry.RegisterTemplate("mydiv", "mydiv-template")
	dom.RegisterMockTemplate("mydiv-template", `<div @click="HandleClick"></div>`)

	input := `<mydiv></mydiv>`
	dom.RegisterMockTemplate("app-root", input)
	w := NewByID("app-root")
	cmp := w.WalkAST(scope.Empty())
	checkWalkErrors(t, w)

	require.Len(t, cmp, 1)
	node := cmp[0]
	child := cmp[0].Children()[0]
	require.True(t, child.Handle("click", &event.Event{}))

	cmpn, ok := node.(*tree.ComponentNode)
	require.True(t, ok)

	getter, ok := cmpn.Instance.Getter("Counter")
	require.True(t, ok)
	require.Equal(t, 17, getter())
}
