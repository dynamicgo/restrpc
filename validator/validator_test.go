package validator

import (
	"encoding/json"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type testA struct {
	A string
	B string
}

type TestB struct {
	A      *testA  `rest:"required"`
	Number float32 `rest:""`
	C      int     `rest:"required"`
}

var input = url.Values{
	"a.a": []string{"hello"},
	"a.b": []string{"world"},
	"c":   []string{"1.899999"},
}

func TestQueryValidator(t *testing.T) {

	var param *TestB

	values, err := Validate(NewQueryReader(input), reflect.TypeOf(param))

	require.NoError(t, err)

	println(printResult(values[0].Interface()))
}

func BenchmarkQueryValidator(t *testing.B) {
	var param *TestB

	for i := 0; i < t.N; i++ {
		Validate(NewQueryReader(input), reflect.TypeOf(param))
	}
}

func printResult(v interface{}) string {
	val, _ := json.MarshalIndent(v, "", "\t")

	return string(val)
}
