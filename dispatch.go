package dispatch

import (
	"fmt"
	"reflect"
)

type TypeHandler func(in interface{}, out interface{}) error
type Fallback func(in interface{}, out interface{}) (handled bool, err error)

type NoFallback struct{ t string }
type NotHandled struct {
	v interface{}
	t string
}
type Dispatcher struct {
	handlers  map[string]TypeHandler
	fallbacks []Fallback
}

func (ø NoFallback) Error() string { return fmt.Sprintf("Error: no fallback function, type %s", ø.t) }
func (ø NotHandled) Error() string {
	return fmt.Sprintf("Error: not handled value %#v type %s", ø.v, ø.t)
}

func New() (ø *Dispatcher) { return &Dispatcher{map[string]TypeHandler{}, []Fallback{}} }

/*
Add a function as fallback.

A fallback function is called by Dispatch(), if no type handler could be found for a certain type of the given value.

A fallback functions is expected to return a boolean to indicate, if it did handle the given value and/or an error if some error occured.

A the fallback returns false, the next (LIFO) fallback function will be called from Dispatch() (see there).
*/
func (ø *Dispatcher) AddFallback(f func(interface{}, interface{}) (bool, error)) {
	ø.fallbacks = append(ø.fallbacks, Fallback(f))
}

/*
Set the handler for a specific type

Only one handler may be specified for a certain type.

You may check if a handler already exists with HasHandler()

You may also get the handler with GetHandler()

If the type is unknown for the registry, an error is returned
*/
func (ø *Dispatcher) SetHandler(i interface{}, f func(interface{}, interface{}) error) {
	ø.handlers[ø.Type(i)] = TypeHandler(f)
}

// removes all fallback functions
func (ø *Dispatcher) RemoveFallbacks()                               { ø.fallbacks = []Fallback{} }
func (ø *Dispatcher) HasHandler(i interface{}) bool                  { return ø.handlers[ø.Type(i)] != nil }
func (ø *Dispatcher) GetHandler(i interface{}) (handler TypeHandler) { return ø.handlers[ø.Type(i)] }
func (ø *Dispatcher) RemoveHandler(i interface{})                    { delete(ø.handlers, ø.Type(i)) }
func (ø *Dispatcher) Type(i interface{}) string                      { return reflect.TypeOf(i).String() }

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
func (ø *Dispatcher) Dispatch(in interface{}, out interface{}) error {
	tt := ø.Type(in)
	m := ø.handlers[tt]

	if m == nil {
		if len(ø.fallbacks) == 0 {
			return NoFallback{tt}
		}
		didHandle := false
		lenfb := len(ø.fallbacks)
		for i := lenfb - 1; i > -1; i-- {
			didHandle, e := ø.fallbacks[i](in, out)
			if e != nil {
				return e
			}
			if didHandle {
				return nil
			}
		}
		if !didHandle {
			return NotHandled{in, tt}
		}
	}
	return m(in, out)
}
