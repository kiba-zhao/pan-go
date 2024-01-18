package peer

const (
	BadRequestErrorCode   = 400
	UnauthorizedErrorCode = 401
	ForbiddenErrorCode    = 403
	NotFoundErrorCode     = 404
	InternalErrorCode     = 500
)

type ResponseError struct {
	code    int
	message string
}

// Code ...
func (re *ResponseError) Code() int {
	return re.code
}

// Error ...
func (re *ResponseError) Error() string {
	return re.message
}

// NewReponseError ...
func NewReponseError(code int, message string) *ResponseError {
	err := new(ResponseError)
	err.code = code
	err.message = message
	return err
}
