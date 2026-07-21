package metadata

import "errors"

var (
	ERR_NOT_FOUND      = errors.New("not found")
	ERR_BAD_REQUEST    = errors.New("bad request")
	ERR_INVALID_ENTITY = errors.New("invalid entity type")
	ERR_INVALID_TYPE   = errors.New("invalid data type")
)
