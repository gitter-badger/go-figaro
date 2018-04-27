// Package figaro is the main package for go-figaro
package figaro

// EncoderDecoder provides determinstic binary encoding/decoding
type EncoderDecoder interface {
	Encoder
	Decoder
}

// Encoder provides deterministic binary encoding
type Encoder interface {
	// Encode binary encodes a given source
	Encode(src interface{}) ([]byte, error)
}

// Decoder provides deterministic binary decoding
type Decoder interface {
	// Decode binary decodes data into a given destination
	Decode(dest interface{}, b []byte) error
}
