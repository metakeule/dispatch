package dispatch

import (
	"testing"
)

func err(t *testing.T, msg string, is interface{}, shouldbe interface{}) {
	t.Errorf(msg+": is %s, should be %s\n", is, shouldbe)
}

func init() {
	d.AddType("")
}

type known int
type notRegistered int

var lastUnhandledContent interface{} = nil
var lastString = ""
var d = New()

func fallback(o interface{}) (didHandle bool, err error) {
	didHandle = true
	lastUnhandledContent = o
	return
}

func handleString(ø interface{}) error {
	lastString = ø.(string)
	return nil
}

func TestRegister(t *testing.T) {
	if d.HasType("xyz") {
		err(t, "HasType with unknown type", true, false)
	}

	if !d.HasType("string") {
		err(t, "HasType with known type", false, true)
	}

	d.AddType(known(0))
	d.RemoveType("known")

	if d.HasType("known") {
		err(t, "RemoveType with known type", false, true)
	}
}

func TestGetType(t *testing.T) {
	if _, e := d.GetType("xyz"); e == nil {
		err(t, "GetType with error", true, false)
	}

	if _, e := d.GetType("string"); e != nil {
		err(t, "GetType without error", true, false)
	}
}

func TestDispatch(t *testing.T) {
	d.AddType(known(0))
	d.AddFallback(fallback)

	d.AddHandler("string", handleString)

	d.Dispatch("hello world")

	if lastString != "hello world" {
		err(t, "Dispatch string", lastString, "hello world")
	}

	lastError := d.Dispatch(notRegistered(0))

	if lastError == nil {
		err(t, "Dispatch error with notRegistered", false, true)
	}

	lastError = d.Dispatch(known(1))

	if lastError != nil {
		err(t, "Dispatch error with known", true, false)
	}

	if lastUnhandledContent != known(1) {
		err(t, "Dispatch error with unhandledContent", lastUnhandledContent, known(1))
	}
}

func TestFallbacks(t *testing.T) {
	res := []string{}
	fb1 := func(val interface{}) (didHandle bool, err error) {
		res = append(res, "fb1"+val.(string))
		return
	}
	fb2 := func(val interface{}) (didHandle bool, err error) {
		res = append(res, "fb2"+val.(string))
		didHandle = true
		return
	}
	fb3 := func(val interface{}) (didHandle bool, err error) {
		res = append(res, "fb3"+val.(string))
		didHandle = true
		return
	}

	disp := New()
	disp.AddType("")
	disp.AddFallback(fb3)
	disp.AddFallback(fb2)
	disp.AddFallback(fb1)

	disp.Dispatch("hiho")

	if res[0] != "fb1hiho" {
		err(t, "Fallback order 0", res[0], "fb1hiho")
	}

	if res[1] != "fb2hiho" {
		err(t, "Fallback order 1", res[1], "fb2hiho")
	}

	if len(res) > 2 {
		err(t, "Fallback order 2", res[2], "")
	}
}
