package param

import (
	"net/http"

	"github.com/dynamicgo/xerrors"
)

type queryReader struct {
	request *http.Request
}

func (reader *queryReader) Read(path string) ([]string, error) {

	keys, ok := reader.request.URL.Query()[path]

	if !ok {
		return nil, xerrors.Wrapf(ErrNotFound, "param with path %s not found", path)
	}

	return keys, nil
}
