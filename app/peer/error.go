package peer

type PeerError struct {
	code int
	err  string
}

func (e *PeerError) Error() string {
	return e.err
}

func (e *PeerError) Code() int {
	return e.code
}
