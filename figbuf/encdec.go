// Package figbuf impelments determinstic binary encoding/decoding
package figbuf

// EncoderDecoder pairs an encoder and decoder for convenience
type EncoderDecoder struct {
	Encoder
	Decoder
}
