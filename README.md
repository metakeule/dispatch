dispatch
========

a simple method dispatcher for go


example
-------

	package main

	import (
		"fmt"
		"github.com/metakeule/dispatch"
	)

	func p(i interface{}) {
		fmt.Println(i)
	}

	type I int

	func main() {

		// register once the types that we want to use in dispatchers
		// by passing an instance to dispatch.Register()
		dispatch.Register("")
		dispatch.Register(I(0))
		dispatch.Register(3)

		// a fallback function for errors and unhandled types
		fallback := func(o interface{}, err error) {
			if err != nil {
				p(err)
				return
			}
			fmt.Printf("unhandled content %#v\n", o)
		}

		// here the functions that are used for the different types
		// they have to cast to the type they use, but they don't
		// need to check for errors, since the casting will always work
		strFn := func(o interface{}) {
			fmt.Printf("%s is a string\n", o.(string))
		}

		iFn := func(o interface{}) {
			fmt.Printf("%d is a I\n", o.(I))
		}

		// setup the dispatcher
		d := dispatch.Dispatcher(fallback)

		// register the handler funcs
		err := d.Handle("string", strFn)
		if err != nil {
			fmt.Println(err)
		}

		err = d.Handle("I", iFn)
		if err != nil {
			fmt.Println(err)
		}

		// let the fun begin!
		d.Dispatch("my string")             // my string is a string
		d.Dispatch(I(3))                    // 3 is a I
		p(dispatch.IsRegistered("float64")) // false
		d.Dispatch(34.0)                    // Error: type float64 is not registered, use dispatch.Register(typeInstance)
		d.Dispatch(34)                      // unhandled content 34
	}
