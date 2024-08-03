package app

import "errors"

const (
	CodeOK            = 0
	CodeInternalError = 500
	CodeNotFound      = 404
)

var ErrNotFound = errors.New("not found")
var ErrUnavailable = errors.New("unavailable")
var ErrInvalidNode = errors.New("invalid node")
var ErrNodeClosed = errors.New("node closed")
var ErrTimeout = errors.New("time out")
