package mock

type encodeFn func(src interface{}) ([]byte, error)
type decodeFn func(dest interface{}, b []byte) error

// EncoderDecoder is a mock implementation of figaro.EncoderDecoder
type EncoderDecoder struct {
	Encoder
	Decoder
}

// Encoder is a mock implementation of figaro.Encoder
type Encoder struct {
	EncodeFn           encodeFn
	EncodeInvoked      bool
	EncodeInvokedCount int
	EncodeInvokedWith  []interface{}
}

// Decoder is a mock implementation of figaro.Decoder
type Decoder struct {
	DecodeFn           decodeFn
	DecodeInvoked      bool
	DecodeInvokedCount int
	DecodeInvokedWith  []interface{}
}

// Encode invokes the mock implementation and marks it as invoked
func (e *Encoder) Encode(src interface{}) ([]byte, error) {
	e.EncodeInvoked = true
	e.EncodeInvokedCount++
	e.EncodeInvokedWith = append(e.EncodeInvokedWith, src)
	return e.EncodeFn(src)
}

// Decode invokes the mock implementation and marks it as invoked
func (d *Decoder) Decode(dest interface{}, b []byte) error {
	d.DecodeInvoked = true
	d.DecodeInvokedCount++
	d.DecodeInvokedWith = append(d.DecodeInvokedWith, []interface{}{dest, b})
	return d.DecodeFn(dest, b)
}
