package ptp

type CodeError interface {
	Code() int
}

type PeerError struct {
	code int
	err  string
}

var _ = (error)((*PeerError)(nil))

// Error returns the string representation of the error.
func (e *PeerError) Error() string {
	return e.err
}

var _ = (CodeError)((*PeerError)(nil))

// Code returns the error code associated with the PeerError.
func (e *PeerError) Code() int {
	return e.code
}
