// Package figbuf impelments determinstic binary encoding/decoding
package figbuf

// EncoderDecoder pairs an encoder and decoder for convenience
type EncoderDecoder struct {
	Encoder
	// Decoder
}

// Decode is a placeholder until the real thing is built
func (ed *EncoderDecoder) Decode(dest interface{}, data []byte) error {
	return nil
}
