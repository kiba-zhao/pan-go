package app

import "errors"

const (
	CodeOK            = 0
	CodeInternalError = 500
	CodeNotFound      = 404
)

var ErrNotFound = errors.New("not found")
