package param

import (
	"net/http"

	"github.com/dynamicgo/xerrors"
)

// ReadFromHTTPRequest .
func ReadFromHTTPRequest(r *http.Request) (Reader, error) {
	switch r.Method {
	case http.MethodGet:
		return &queryReader{
			request: r,
		}, nil
	case http.MethodPost:

		contentType := r.Header.Get("Content-type")

		if contentType != "application/json" {
			return nil, xerrors.Wrapf(ErrMimeType, "unsupport mine-type %s", contentType)
		}

		return &bodyReader{
			request: r,
		}, nil
	default:
		return nil, xerrors.Wrapf(ErrNotSupportRequest, "unsupport http request %s", r.Method)
	}

}
