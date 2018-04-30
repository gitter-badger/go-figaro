package mock

type hashFn func(b ...[]byte) []byte

// Hasher is a mock implementation of figaro.Hasher
type Hasher struct {
	HashFn           hashFn
	HashInvoked      bool
	HashInvokedCount int
	HashInvokedWith  []interface{}
}

// Hash invokes the mock implementation and marks it as invoked
func (h *Hasher) Hash(b ...[]byte) []byte {
	h.HashInvoked = true
	h.HashInvokedCount++
	h.HashInvokedWith = append(h.HashInvokedWith, b)
	return h.HashFn(b...)
}
