package restrpc

import (
	"github.com/dynamicgo/xerrors/apierr"
)

// Errors
var (
	ErrInternal = apierr.New(-1, "INNER_ERROR")
)
