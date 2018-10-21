package param

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/Jeffail/gabs"
	"github.com/dynamicgo/xerrors"
)

type bodyReader struct {
	request    *http.Request
	jsonparsed *gabs.Container
}

func (reader *bodyReader) Read(path string) ([]string, error) {

	var err error

	if reader.jsonparsed == nil {
		reader.jsonparsed, err = gabs.ParseJSONBuffer(reader.request.Body)
	}

	if err != nil {
		return nil, err
	}

	if !reader.jsonparsed.ExistsP(path) {
		return nil, xerrors.Wrapf(ErrNotFound, "param with path %s not found", path)
	}

	switch v := reader.jsonparsed.Path(path).Data().(type) {
	case []interface{}:
		var array []string
		for _, val := range v {
			array = append(array, fmt.Sprintf("%v", val))
			return array, nil
		}
	case interface{}:
		return []string{fmt.Sprintf("%v", v)}, nil
	default:
		return nil, xerrors.Wrapf(ErrParamType, "unspport parameter %s with type %s", path, reflect.TypeOf(v))
	}

	return nil, xerrors.Wrapf(ErrParamType, "unspport parameter %s with type %s", path, reader.jsonparsed.Path(path).Data())
}
