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

	// setup the dispatcher
	d := dispatch.New()

	// we need to register once the types that we want to use in dispatchers
	// by passing an instance to AddType()
	d.AddType("")
	d.AddType(I(0))
	d.AddType(3)

	// we may have fallback functions for unhandled types
	fallback := func(o interface{}) (didHandle bool, err error) {
		didHandle = true
		fmt.Printf("fallback for %#v\n", o)
		return
	}

	d.AddFallback(fallback)

	// here the functions that are used for the different types
	// they have to cast to the type they serve, but they don't
	// need to check for type casting of the interface.
	// they should however return an error for other situations
	strHandler := func(i interface{}) (err error) {
		fmt.Printf("%s is a string\n", i.(string))
		return
	}

	// returns an error if the string type is not registered with AddType()
	d.SetHandler("string", strHandler)

	iHandler := func(i interface{}) (err error) {
		fmt.Printf("%d is a I\n", i.(I))
		return
	}

	d.SetHandler("I", iHandler)

	// let the fun begin!
	d.Dispatch("my string")       // my string is a string
	d.Dispatch(I(3))              // 3 is a I
	p(d.HasType("float64"))       // false
	fmt.Println(d.Dispatch(34.0)) // Error: type float64 is not registered, use AddType()
	d.Dispatch(34)                // fallback for 34
}
