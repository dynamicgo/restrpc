package client

import (
	"reflect"
	"testing"
)

type testA interface {
	Hello(string) error
}

func TestImplInterface(t *testing.T) {
	var i testA

	it := reflect.TypeOf(&i).Elem()

	mt := it.Method(0).Type

	fn := reflect.MakeFunc(mt, func(args []reflect.Value) []reflect.Value {
		return nil
	})

	reflect.ValueOf(&i).Elem().Set(fn)
}
