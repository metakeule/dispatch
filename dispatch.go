package dispatch

import (
	"fmt"
	"reflect"
)

type TypeHandler func(o interface{}) error
type Fallback func(o interface{}) (handled bool, err error)

type NotInRegistry struct{ t string }
type NoFallback struct{ t string }
type NotHandled struct {
	v interface{}
	t string
}
type DispatcherType struct{ reflect.Type }
type Dispatcher struct {
	handlers  map[*DispatcherType]TypeHandler
	fallbacks []Fallback
	registry  map[string]*DispatcherType
}

func (ø *DispatcherType) String() string { return ø.Name() }

func (ø NotInRegistry) Error() string {
	return fmt.Sprintf("Error: type %s is not registered, use AddType()", ø.t)
}

func (ø NoFallback) Error() string {
	return fmt.Sprintf("Error: type %s has no handler and we have no fallback function, use AddFallback()", ø.t)
}

func (ø NotHandled) Error() string {
	return fmt.Sprintf("Error: value %#v type %s has no handler and no fallback function did handle it, use AddFallback()", ø.v, ø.t)
}

func New() (ø *Dispatcher) {
	ø = &Dispatcher{
		handlers:  map[*DispatcherType]TypeHandler{},
		fallbacks: []Fallback{},
		registry:  map[string]*DispatcherType{}}

	return
}

// import reflect types for reuse
func (ø *Dispatcher) ImportReflectType(types ...reflect.Type) {
	for _, t := range types {
		ø.registry[t.Name()] = &DispatcherType{t}
	}
}

// import types for reuse
func (ø *Dispatcher) ImportType(types ...DispatcherType) {
	for _, t := range types {
		ø.AddType(t)
	}
}

/*
Add a function as fallback.

A fallback function is called by Dispatch(), if no type handler could be found for a certain type of the given value.

A fallback functions is expected to return a boolean to indicate, if it did handle the given value and/or an error if some error occured.

A the fallback returns false, the next (LIFO) fallback function will be called from Dispatch() (see there).
*/
func (ø *Dispatcher) AddFallback(f func(interface{}) (bool, error)) {
	ø.fallbacks = append(ø.fallbacks, Fallback(f))
}

// removes all fallback functions
func (ø *Dispatcher) RemoveFallbacks() { ø.fallbacks = []Fallback{} }

/*
Set the handler for a specific type

Only one handler may be specified for a certain type.

You may check if a handler already exists with HasHandler()

You may also get the handler with GetHandler()

If the type is unknown for the registry, an error is returned
*/
func (ø *Dispatcher) SetHandler(ty string, f func(interface{}) error) (err error) {
	real, err := ø.GetType(ty)
	if err != nil {
		return
	}
	ø.handlers[real] = TypeHandler(f)
	return
}

func (ø *Dispatcher) HasHandler(ty string) (has bool, err error) {
	real, err := ø.GetType(ty)
	if err != nil {
		return
	}
	has = ø.handlers[real] != nil
	return
}

func (ø *Dispatcher) GetHandler(ty string) (handler TypeHandler, err error) {
	real, err := ø.GetType(ty)
	if err != nil {
		return
	}
	handler = ø.handlers[real]
	return
}

func (ø *Dispatcher) RemoveHandler(ty string) (err error) {
	real, err := ø.GetType(ty)
	if err != nil {
		return
	}
	delete(ø.handlers, real)
	return
}

/*
Dispatch() takes any value, looks if it can find a handler function to handle its type and calls that handler with the value.

If there is no specific handler for the type of value, the fallback functions are called in the reverse order of their registration (LIFO) until one of them
either returns an error or a boolean with true that means, that the function did handle the value.

Dispatch() returns an error if one of the following conditions are met:

	- the type of the value is not in the registry. fix it with AddType()
	- the handler function returns an error. the error is passed through
	- there is no handler function for the type of the value and there also is
	  no fallback function that could/did handle the value. fix it with SetHandler() or AddFallback()
	- a fallback function returned an error. the error is passed through
*/
func (ø *Dispatcher) Dispatch(o interface{}) error {
	tt := reflect.TypeOf(o).Name()
	m, err := ø.GetHandler(tt)
	if err != nil {
		return err
	}

	if m == nil {
		if len(ø.fallbacks) == 0 {
			return NoFallback{tt}
		}
		didHandle := false
		lenfb := len(ø.fallbacks)
		for i := lenfb - 1; i > -1; i-- {
			didHandle, e := ø.fallbacks[i](o)
			if e != nil {
				return e
			}
			if didHandle {
				return nil
			}
		}
		if !didHandle {
			return NotHandled{o, tt}
		}
	}
	return m(o)
}

func (ø *Dispatcher) AddType(i interface{}) {
	t := reflect.TypeOf(i)
	ø.registry[t.Name()] = &DispatcherType{t}
}

func (ø *Dispatcher) RemoveType(name string) { delete(ø.registry, name) }

func (ø *Dispatcher) GetType(name string) (t *DispatcherType, err error) {
	t = ø.registry[name]
	if t == nil {
		err = NotInRegistry{name}
	}
	return
}

func (ø *Dispatcher) HasType(name string) bool {
	if ø.registry[name] == nil {
		return false
	}
	return true
}
