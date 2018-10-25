package validator

import (
	"errors"
	"reflect"
	"strings"

	"github.com/dynamicgo/xerrors"
)

// Errors
var (
	ErrInvalidType = errors.New("invalid parameter type")
)

// MetadataTag .
const MetadataTag = "rest"

// Metadata .
type Metadata struct {
	Required bool   // required parameter flag
	Name     string // parameter name
}

// ParseMetadata .
func ParseMetadata(tag string) *Metadata {
	metadata := &Metadata{}

	if tag == "" {
		return metadata
	}

	tokens := strings.Split(tag, ",")

	for _, token := range tokens {
		switch token {
		case "required":
			metadata.Required = true
		default:
			metadata.Name = token
		}
	}

	return metadata
}

// Validator .
type Validator interface {
	Validate(reader Reader, paramT reflect.Type) ([]reflect.Value, error)
}

// Reader parameter reader
type Reader interface {
	Search(key string) ([]string, error)
	Get() ([]string, error)
	Path() string // reader path
}

// Validate validate parameter with reflect type and reader
func Validate(reader Reader, paramT reflect.Type) ([]reflect.Value, error) {

	if paramT.Kind() == reflect.Ptr {
		paramT = paramT.Elem()
	}

	var validator Validator

	switch paramT.Kind() {
	case reflect.Int, reflect.Int8,
		reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint,
		reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
	case reflect.Struct:
	case reflect.String:
	case reflect.Array:
	case reflect.Map:
	default:
		if paramT.Implements(reflect.TypeOf(&validator).Elem()) {
			value := reflect.New(paramT)
			validator = value.Interface().(Validator)
		}
	}

	if validator == nil {
		return nil, xerrors.Wrapf(ErrInvalidType, "not support parameter type %s with path %s", paramT, reader.Path())
	}

	return validator.Validate(reader, paramT)
}

type numberValidator struct {
}

type structValidator struct {
}

type stringValidator struct {
}

type arrayValidator struct {
}

type mapValidator struct {
}
