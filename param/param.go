package param

import (
	"errors"
)

// Errors
var (
	ErrNotSupportRequest = errors.New("unsupport http request")
	ErrNotFound          = errors.New("resource not found")
	ErrMimeType          = errors.New("post only support mime-type application/json")
	ErrParamType         = errors.New("unsupport param type")
)

// Reader .
type Reader interface {
	Read(path string) ([]string, error)
}

// Writer .
type Writer interface {
	Write(path string, values ...string) error
}
