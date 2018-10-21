package param

import (
	"reflect"

	"github.com/dynamicgo/xerrors"
)

// New create new param with paramT
func newParam(paramT reflect.Type) (reflect.Value, error) {

	if paramT.Kind() != reflect.Ptr {
		return reflect.Value{}, xerrors.Wrapf(ErrParamType, "param type must be ptr, got %s", paramT)
	}

	paramT = paramT.Elem()

	value := reflect.New(paramT)

	return value, nil
}

// Read read param
func Read(paramT reflect.Type, reader Reader) (reflect.Value, error) {

	if paramT.Kind() == reflect.Ptr {
		paramT = paramT.Elem()
	}

	if paramT.Kind() != reflect.Struct {
		return reflect.Value{}, xerrors.Wrapf(ErrParamType, "param type must be ptr ot struct, got %s", paramT)
	}

	return read("", paramT, reader)
}

func read(path string, paramT reflect.Type, reader Reader) (reflect.Value, error) {

	value, err := newParam(paramT)

	if err != nil {
		return reflect.Value{}, nil
	}

	if paramT.Kind() == reflect.Ptr {
		paramT = paramT.Elem()
	}

	for i := 0; i < paramT.NumField(); i++ {

	}

	return reflect.Value{}, nil
}
