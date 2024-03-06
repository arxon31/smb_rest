package v1

import "errors"

var (
	errFileSizeIsTooBig = errors.New("file size is too big")
	errInternalError    = errors.New("internal error")
	errBadRequest       = errors.New("bad request")
)
