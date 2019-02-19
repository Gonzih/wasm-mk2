package core

import (
	"log"
	"strings"

	"github.com/Gonzih/wasm-mk2/component"
	"github.com/Gonzih/wasm-mk2/dom"
	"github.com/Gonzih/wasm-mk2/parser"
	"github.com/Gonzih/wasm-mk2/registry"
	"github.com/Gonzih/wasm-mk2/walker"
	"golang.org/x/net/html"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Component(strukt component.ComponentInput, name, templateID string) {
	wrapper, err := component.Wasmify(strukt)
	must(err)
	registry.Register(name, wrapper)
	markup := dom.New().TemplateContent(templateID)
	r := strings.NewReader(markup)
	z := html.NewTokenizer(r)
	p := parser.New(z)
	w := walker.New(p)
	log.Println(w)
}

type App struct {
}

func New(targetID string) *App {
	return &App{}
}

func (a *App) Mount(targetID string) error {
	return nil
}
