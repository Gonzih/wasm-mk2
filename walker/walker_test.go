package walker

import (
	"strings"
	"testing"

	"github.com/Gonzih/wasm-mk2/parser"
	"github.com/stretchr/testify/assert"

	"golang.org/x/net/html"
)

func testParse(t *testing.T, input string) *Walker {
	r := strings.NewReader(input)
	z := html.NewTokenizer(r)
	p := parser.New(z)
	w := New(p)

	assert.Len(t, w.parser.Errors(), 0)

	for _, e := range w.parser.Errors() {
		t.Error(e)
	}

	if len(w.parser.Errors()) > 0 {
		t.FailNow()
	}

	return w
}

func TestBasic(t *testing.T) {

}
