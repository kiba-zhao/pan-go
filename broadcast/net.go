package broadcast

type Net interface {
	Read(int) ([]byte, []byte, error)
	Write([]byte) error
	Close() error
}

type NetWriteError struct {
	writenSize int
}

// Error ...
func (e *NetWriteError) Error() string {
	return "Truncate Write"
}
