// Define peer application error

package peer

type PeerError struct {
	code int
	err  string
}

// Error returns the string representation of the error.
func (e *PeerError) Error() string {
	return e.err
}

// Code returns the error code associated with the PeerError.
func (e *PeerError) Code() int {
	return e.code
}
