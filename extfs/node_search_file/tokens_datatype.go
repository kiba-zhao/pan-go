package nodesearchfile

import (
	"errors"
	"strings"
)

var ErrInvalidTokens = errors.New("invalid tokens")

func (t *Tokens) Scan(value interface{}) error {
	text, ok := value.(string)
	if !ok {
		return ErrInvalidTokens
	}
	*t = strings.Split(text, ",")
	return nil
}

func (t Tokens) Value() (interface{}, error) {
	if len(t) == 0 {
		return nil, nil
	}
	return strings.Join(t, ","), nil
}
