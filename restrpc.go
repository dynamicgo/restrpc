package restrpc

import (
	"errors"
	"io"
	"reflect"

	"github.com/dynamicgo/xerrors/apierr"
)

// Errors
var (
	ErrInternal    = apierr.New(-1, "INNER_ERROR")
	ErrInvalidType = errors.New("invalid param type")
	ErrMapKey      = errors.New("map key must be string")
)

// Reader parameter reader
type Reader interface {
	Search(key string) ([]string, error)
	Get() ([]string, error)
	Range(func(key string, reader Reader) error) error
	Path() string // reader path
	Reader(key string) Reader
}

// Validator .
type Validator interface {
	Validate(reader Reader, paramT reflect.Type) ([]reflect.Value, error)
}

//  Writer .
type Writer interface {
	Write(writer io.Writer) error
}
