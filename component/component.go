package component

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

const tagKey = "wasm"

type ComponentInput interface {
	Init() error
}

type Wrapper struct {
	input    interface{}
	instance interface{}
	setters  map[string]func(interface{}) error
	getters  map[string]func() interface{}
	props    map[string]bool
}

func Wasmify(comp ComponentInput) (*Wrapper, error) {
	wrapper := &Wrapper{input: comp}

	in := reflect.ValueOf(comp)
	if in.Kind() != reflect.Ptr {
		return nil, errors.New("Wasmify only accepts pointers")
	}

	val := in.Elem()

	if val.Kind() != reflect.Struct {
		return nil, errors.New("Wasmify only accepts pointers to structs")
	}

	return wrapper, nil
}

func (w *Wrapper) Instance() (*Wrapper, error) {
	wrapperCpy := *w
	result := &wrapperCpy

	result.getters = make(map[string]func() interface{}, 0)
	result.setters = make(map[string]func(interface{}) error, 0)
	result.props = make(map[string]bool, 0)

	result.instance = reflect.New(reflect.ValueOf(result.input).Elem().Type()).Interface()
	in, ok := result.instance.(ComponentInput)
	if ok {
		in.Init()
	} else {
		return nil, fmt.Errorf("Could not cast type %s to ComponentInput interface", reflect.ValueOf(result.instance).Elem().Type())
	}

	err := result.constructGetters()
	if err != nil {
		return nil, errors.Wrap(err, "Could not create Wrapper instance")
	}

	err = result.constructSetters()
	if err != nil {
		return nil, errors.Wrap(err, "Could not create Wrapper instance")
	}

	result.findProps()

	return result, nil
}

func (w *Wrapper) constructGetters() error {
	val := reflect.ValueOf(w.instance).Elem()

	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		name := typeField.Name

		getter := func() interface{} {
			return reflect.ValueOf(w.instance).Elem().FieldByName(name).Interface()
		}

		w.getters[name] = getter
	}

	return nil
}

func (w *Wrapper) constructSetters() error {
	val := reflect.ValueOf(w.instance).Elem()

	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		name := typeField.Name

		setter := func(in interface{}) error {
			targetField := reflect.ValueOf(w.instance).Elem().FieldByName(name)
			input := reflect.ValueOf(in)

			if targetField.Type() != input.Type() {
				return errors.New(fmt.Sprintf("Mismatched target and input types %s != %s", input.Type(), targetField.Type()))
			}

			targetField.Set(input)

			return nil
		}

		w.setters[name] = setter
	}

	return nil
}

func (w *Wrapper) findProps() {
	val := reflect.ValueOf(w.instance).Elem()

	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		ts, _ := typeField.Tag.Lookup(tagKey)
		tags := strings.Split(ts, ",")
		for _, tag := range tags {
			if tag == "prop" {
				w.props[typeField.Name] = true
			}
		}
	}
}

func (w *Wrapper) Getter(name string) (func() interface{}, bool) {
	f, ok := w.getters[name]

	return f, ok
}

func (w *Wrapper) Setter(name string) (func(interface{}) error, bool) {
	f, ok := w.setters[name]

	return f, ok
}

func (w *Wrapper) IsAProp(name string) bool {
	_, ok := w.props[name]
	return ok
}
