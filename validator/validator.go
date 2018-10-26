package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/dynamicgo/xerrors"
)

// Errors
var (
	ErrInvalidType = errors.New("invalid parameter type")
	ErrNumber      = errors.New("parse number error")
	ErrKey         = errors.New("restrpc map only support string key")
	ErrInner       = errors.New("inner error")
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
	Range(func(key string, reader Reader) error) error
	Path() string // reader path
	Reader(key string) Reader
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
		validator = &numberValidator{}
	case reflect.Struct:
		validator = &structValidator{}
	case reflect.String:
		validator = &stringValidator{}
	case reflect.Array:
		validator = &arrayValidator{}
	case reflect.Map:
		validator = &mapValidator{}
	case reflect.Bool:
		validator = &boolValidator{}
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

type boolValidator struct {
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

func (validator *boolValidator) Validate(reader Reader, paramT reflect.Type) ([]reflect.Value, error) {

	params, err := reader.Get()

	if err != nil {
		return nil, err
	}

	var values []reflect.Value

	for _, param := range params {

		boolean, err := strconv.ParseBool(param)

		if err != nil {
			return nil, xerrors.Wrapf(ErrNumber, "parse boolean parameter %s with value %s error", reader.Path(), param)
		}

		values = append(values, reflect.ValueOf(boolean))
	}

	return values, nil
}

func (validator *numberValidator) Validate(reader Reader, paramT reflect.Type) ([]reflect.Value, error) {

	params, err := reader.Get()

	if err != nil {
		return nil, err
	}

	var values []reflect.Value

	for _, param := range params {

		number, err := strconv.ParseFloat(param, 64)

		if err != nil {
			return nil, xerrors.Wrapf(ErrNumber, "parse number parameter %s with value %s error", reader.Path(), param)
		}

		values = append(values, reflect.ValueOf(number))
	}

	return values, nil
}

func (validator *structValidator) Validate(reader Reader, paramT reflect.Type) ([]reflect.Value, error) {

	mapValue := reflect.New(paramT)

	for i := 0; i < paramT.NumField(); i++ {
		field := paramT.Field(i)

		metadata := ParseMetadata(field.Tag.Get(MetadataTag))

		name := strings.ToLower(field.Name)

		if metadata.Name != "" {
			name = metadata.Name
		}

		fieldReader := reader.Reader(name)

		values, err := Validate(fieldReader, field.Type)

		if err != nil {
			return nil, err
		}

		if len(values) == 0 {
			if metadata.Required {
				return nil, fmt.Errorf("expect param %s", fieldReader.Path())
			}

			continue
		}

		mapValue.Field(i).Set(values[0])
	}

	return []reflect.Value{mapValue}, nil
}

func (validator *stringValidator) Validate(reader Reader, paramT reflect.Type) ([]reflect.Value, error) {
	params, err := reader.Get()

	if err != nil {
		return nil, err
	}

	var values []reflect.Value

	for _, param := range params {

		values = append(values, reflect.ValueOf(param))
	}

	return values, nil
}

func (validator *arrayValidator) Validate(reader Reader, paramT reflect.Type) ([]reflect.Value, error) {

	paramT = paramT.Elem()

	return Validate(reader, paramT)
}

func (validator *mapValidator) Validate(reader Reader, paramT reflect.Type) ([]reflect.Value, error) {
	keyT := paramT.Key()

	if keyT.Kind() != reflect.String {
		return nil, xerrors.Wrapf(ErrKey, "map %s key must be string", reader.Path())
	}

	valueT := paramT.Elem()

	mapValue := reflect.New(paramT)

	path := reader.Path()

	err := reader.Range(func(key string, reader Reader) error {

		values, err := Validate(reader, valueT)

		if err != nil {
			return err
		}

		if len(values) == 0 {
			return xerrors.Wrapf(ErrInner, "expect map %s key %s value", path, key)
		}

		mapValue.SetMapIndex(reflect.ValueOf(key), values[0])

		return nil
	})

	return []reflect.Value{mapValue}, err
}
