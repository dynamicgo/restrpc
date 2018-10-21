package param

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type Param1 struct {
	val *Param1
}

type Result struct {
	A *int
	P *Param1
}

func TestCreate(t *testing.T) {

	var r *Result

	value, err := newParam(reflect.TypeOf(r))

	require.NoError(t, err)

	r = value.Interface().(*Result)

	println(printResult(r))
}

func printResult(v interface{}) string {
	val, _ := json.MarshalIndent(v, "", "\t")

	return string(val)
}
