package client

import (
	"bytes"
	"fmt"
	"io"
	"reflect"

	"github.com/dynamicgo/xerrors"

	"github.com/dynamicgo/restrpc"
)

func writeJSON(args interface{}, writer io.Writer) error {

	paramT := reflect.TypeOf(args)

	if paramT.Kind() == reflect.Ptr {
		paramT = paramT.Elem()
	}

	switch paramT.Kind() {
	case reflect.Int, reflect.Int8,
		reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint,
		reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return writeNumber(args, writer)
	case reflect.Struct:
		return writeStruct(args, writer)
	case reflect.String:
		return writeString(args, writer)
	case reflect.Slice, reflect.Array:
		return writeArray(args, writer)
	case reflect.Map:
		return writeMap(args, writer)
	case reflect.Bool:
		return writeBool(args, writer)
	default:

		var target restrpc.Writer

		if paramT.Implements(reflect.TypeOf(&target).Elem()) {
			target = args.(restrpc.Writer)
			return target.Write(writer)
		}

		return xerrors.Wrapf(restrpc.ErrInvalidType, "invalid param type %s", paramT)
	}
}

func writeBool(args interface{}, writer io.Writer) error {

	b := args.(bool)

	if b {
		_, err := writer.Write([]byte("true"))

		return err
	} else {
		_, err := writer.Write([]byte("false"))

		return err
	}
}

func writeNumber(args interface{}, writer io.Writer) error {

	v := fmt.Sprintf("%v", args)

	_, err := writer.Write([]byte(v))

	return err
}

func writeStruct(args interface{}, writer io.Writer) error {
	return nil
}

func writeString(args interface{}, writer io.Writer) error {

	_, err := writer.Write([]byte(fmt.Sprintf(`"%s"`, args.(string))))

	return err
}

func writeArray(args interface{}, writer io.Writer) error {

	var buff bytes.Buffer

	buff.WriteString("[")

	value := reflect.ValueOf(args)

	for i := 0; i < value.Len(); i++ {
		if err := writeJSON(value.Index(i).Interface(), &buff); err != nil {
			return err
		}

		if i+1 < value.Len() {
			buff.WriteString(",")
		}
	}

	buff.WriteString("]")

	_, err := writer.Write(buff.Bytes())

	return err
}

func writeMap(args interface{}, writer io.Writer) error {

	var buff bytes.Buffer

	buff.WriteString("{")

	value := reflect.ValueOf(args)

	keys := value.MapKeys()

	for i, key := range keys {

		writeJSON(key.Interface(), &buff)

		buff.WriteString(":")

		keyValue := value.MapIndex(key)

		if keyValue.Type().Kind() != reflect.String {
			return xerrors.Wrapf(restrpc.ErrMapKey, "only support string key,got %s", keyValue.Type())
		}

		if err := writeJSON(keyValue.Interface(), &buff); err != nil {
			return err
		}

		if i+1 < value.Len() {
			buff.WriteString(",")
		}
	}

	buff.WriteString("}")

	_, err := writer.Write(buff.Bytes())

	return err
}
