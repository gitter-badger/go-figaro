// Package figbuf implements figaro domain specific wrappers for figbuf
package figbuf

import "errors"

// ErrInvalidData is a self-explantory error
var ErrInvalidData = errors.New("figbuf: invalid data for type")

// EncoderDecoder maps the figaro domain to figbuf
type EncoderDecoder struct{}
