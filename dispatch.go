package dispatch

import (
	"fmt"
	"reflect"
)

type NotInRegistry struct{ t string }
type dispatcherType struct{ reflect.Type }
type dispatcher struct {
	d        map[*dispatcherType]func(o interface{})
	fallback func(o interface{}, e error)
}

var registry = map[string]*dispatcherType{}

func (ø NotInRegistry) Error() string {
	return fmt.Sprintf("Error: type %s is not registered, use dispatch.Register(typeInstance)", ø.t)
}
func (ø *dispatcherType) String() string { return ø.Name() }

func Dispatcher(f func(o interface{}, e error)) *dispatcher {
	return &dispatcher{d: map[*dispatcherType]func(o interface{}){}, fallback: f}
}

func (ø *dispatcher) Handle(ty string, f func(o interface{})) (err error) {
	real, err := GetType(ty)
	if err != nil {
		return
	}
	ø.d[real] = f
	return
}

func (ø *dispatcher) Dispatch(o interface{}) {
	tt := reflect.TypeOf(o).Name()
	real, err := GetType(tt)
	if err != nil {
		ø.fallback(o, err)
		return
	}
	m := ø.d[real]
	if m == nil {
		ø.fallback(o, nil)
		return
	}
	m(o)
	return
}

func Register(i interface{}) {
	t := reflect.TypeOf(i)
	registry[t.Name()] = &dispatcherType{t}
}

func GetType(name string) (t *dispatcherType, err error) {
	t = registry[name]
	if t == nil {
		err = NotInRegistry{name}
	}
	return
}

func IsRegistered(name string) bool {
	if registry[name] == nil {
		return false
	}
	return true
}
