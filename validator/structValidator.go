package validator

import (
	"reflect"
	"strings"

	"github.com/dynamicgo/xerrors"
)

type structValidator struct {
}

func (validator *structValidator) Validate(reader Reader, parmT reflect.Type) (reflect.Value, error) {
	if parmT.Kind() == reflect.Ptr {
		parmT = parmT.Elem()
	}

	if parmT.Kind() != reflect.Struct {
		return reflect.Value{}, xerrors.Wrapf(ErrValidator, "expect struct param Type,got %s", parmT)
	}

	for i := 0; i < parmT.NumField(); i++ {
		field := parmT.Field(i)

		validator, err := NewValidator(field.Type)

		if err != nil {
			return reflect.Value{}, err
		}

		validator.Validate(reader.GetReader(strings.ToLower(field.Name())), field.Type)
	}

	return reflect.Value{}, nil
}
