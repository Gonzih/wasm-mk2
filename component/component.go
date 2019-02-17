package component

// type ComponentWrapper struct {
// 	component interface{}
// 	setters   map[string]func(interface{})
// 	getters   map[string]func() interface{}
// 	props     []string
// 	state     []string
// }

// func Wasmify(comp interface{}, name, id string) *ComponentWrapper {
// 	log.Println(comp)
// 	in := reflect.ValueOf(comp)
// 	if in.Kind() != reflect.Ptr {
// 		log.Fatalf("Wasmify only accepts pointers")
// 	}

// 	val := in.Elem()

// 	for i := 0; i < val.NumField(); i++ {
// 		// valueField := val.Field(i)
// 		typeField := val.Type().Field(i)
// 		log.Println(typeField.Tag.Lookup("wasm"))
// 	}

// 	wrapper := &componentWrapper{component: comp}
// 	log.Println(wrapper)

// 	return wrapper
// }

// type Event struct {
// }

type Component struct {
	Tag string
}
