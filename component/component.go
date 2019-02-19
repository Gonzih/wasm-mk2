package component

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Gonzih/wasm-mk2/event"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

const tagKey = "wasm"

type ComponentInput interface {
	Init() error
}

type Wrapper struct {
	uuid     string
	input    interface{}
	instance interface{}
	setters  map[string]func(interface{}) error
	getters  map[string]func() interface{}
	handlers map[string]func(*event.Event)
	props    map[string]string
}

func Wasmify(comp interface{}) (*Wrapper, error) {
	wrapper := &Wrapper{input: comp}

	in := reflect.ValueOf(comp)
	if in.Kind() != reflect.Ptr {
		return nil, errors.New("Wasmify only accepts pointers")
	}

	val := in.Elem()

	if val.Kind() != reflect.Struct {
		return nil, errors.New("Wasmify only accepts pointers to structs")
	}

	_, ok := comp.(ComponentInput)
	if !ok {
		return nil, errors.New("Wasmify only accepts structs that implement ComponentInput interface")
	}

	return wrapper, nil
}

func (w *Wrapper) Instance() (*Wrapper, error) {
	wrapperCpy := *w
	result := &wrapperCpy

	result.getters = make(map[string]func() interface{}, 0)
	result.setters = make(map[string]func(interface{}) error, 0)
	result.props = make(map[string]string, 0)
	result.handlers = make(map[string]func(*event.Event), 0)

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
	result.findHandlers()

	result.uuid = uuid.NewV4().String()

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
				w.props[strings.ToLower(typeField.Name)] = typeField.Name
			}
		}
	}
}

func (w *Wrapper) findHandlers() {
	val := reflect.ValueOf(w.instance)

	for i := 0; i < val.NumMethod(); i++ {
		method := val.Type().Method(i)
		if strings.HasPrefix(method.Name, "Handle") {
			w.handlers[method.Name] = func(e *event.Event) {
				arg := reflect.ValueOf(e)
				method.Func.Call([]reflect.Value{val, arg})
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

func (w *Wrapper) IsAProp(name string) (string, bool) {
	field, ok := w.props[name]
	return field, ok
}

func (w *Wrapper) IsAHandler(name string) bool {
	_, ok := w.handlers[name]
	return ok
}

func (w *Wrapper) Handler(name string) (func(*event.Event), bool) {
	h, ok := w.handlers[name]
	return h, ok
}

func (w *Wrapper) UUID() string {
	return w.uuid
}
