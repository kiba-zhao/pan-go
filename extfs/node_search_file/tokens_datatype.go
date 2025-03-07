// Custom data type for tokens field for gorm
package nodesearchfile

import (
	"errors"
	"strings"
)

var ErrInvalidTokens = errors.New("invalid tokens")

// Scan implements the Scanner interface.
//
// The value must be a string. The method splits the string by comma and assigns
// the result to the Tokens object. If the value is not a string, it returns
// ErrInvalidTokens.
func (t *Tokens) Scan(value interface{}) error {
	text, ok := value.(string)
	if !ok {
		return ErrInvalidTokens
	}
	*t = strings.Split(text, ",")
	return nil
}

// Value implements the Valuer interface.
//
// If the length of the Tokens slice is 0, it returns nil. Otherwise, it joins the
// slice with comma and returns the result as a string.
func (t Tokens) Value() (interface{}, error) {
	if len(t) == 0 {
		return nil, nil
	}
	return strings.Join(t, ","), nil
}
