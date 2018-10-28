package validator

import (
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
	Number float32 `rest:"-"`
}

func TestQueryValidator(t *testing.T) {
	input := url.Values{
		"a.a": []string{"hello"},
		"a.b": []string{"world"},
	}

	var param *TestB

	_, err := Validate(NewQueryReader(input), reflect.TypeOf(param))

	require.NoError(t, err)
}
