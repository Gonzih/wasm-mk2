package walker

import (
	"strings"
	"testing"

	"github.com/Gonzih/wasm-mk2/parser"
	"github.com/stretchr/testify/assert"

	"golang.org/x/net/html"
)

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
	assert.Equal(t, "div", cmp.Tag)
}

func TestNested(t *testing.T) {
	input := `<div><p></p><a></a></div>`
	w := walkString(t, input)
	cmp := w.WalkAST()
	checkWalkErrors(t, w)

	assert.Len(t, cmp, 1)
	assert.Equal(t, "div", cmp.Tag)
	assert.Equal(t, "p", cmp.Children[0].Tag)
	assert.Equal(t, "a", cmp.Children[1].Tag)
}
