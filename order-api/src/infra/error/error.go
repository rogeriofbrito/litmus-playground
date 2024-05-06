package infra_error

import "errors"

var (
	ErrQueryNotReturnValues = errors.New("query doesn't return values")
)
