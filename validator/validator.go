package validator

import (
	"errors"
	"reflect"

	"github.com/dynamicgo/xerrors"
)

// Errors
var (
	ErrValidator = errors.New("invalid valididator,this inner error")
	ErrNotFound  = errors.New("resource not found")
)

// the validator include three part: reflect type validator, data reader and controller

// Validator type validator
type Validator interface {
	Validate(reader Reader, parmT reflect.Type) (reflect.Value, error)
}

// Validatable the type implement validatable
type Validatable interface {
	Validate() error
}

// Reader input data reader
type Reader interface {
	Get(path string) ([]string, error)
	GetReader(path string) Reader
}

// NewValidator create new validator for given parameter type
func NewValidator(parmT reflect.Type) (Validator, error) {

	if parmT.Kind() == reflect.Ptr {
		parmT = parmT.Elem()
	}

	switch parmT.Kind() {
	case reflect.Struct:
		return &structValidator{}, nil
	default:
		return nil, xerrors.Wrapf(ErrNotFound, "unsupport type %s validator", parmT)
	}
}
