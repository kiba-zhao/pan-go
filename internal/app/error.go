package app

import "pan/internal/servlet"

const (
	CodeOK            = servlet.CodeOK
	CodeInternalError = -500
	CodeNotFound      = -404
	CodeBadRequest    = -400
	CodeForbidden     = -403
)
