package mock

type getFn func(key []byte) []byte
type setFn func(key []byte, value []byte)
type deleteFn func(key []byte)

// Store is a mock implementation of figaro.Store
type Store struct {
	GetFn           getFn
	GetInvoked      bool
	GetInvokedCount int
	GetInvokedWith  []interface{}

	SetFn           setFn
	SetInvoked      bool
	SetInvokedCount int
	SetInvokedWith  []interface{}

	DeleteFn           deleteFn
	DeleteInvoked      bool
	DeleteInvokedCount int
	DeleteInvokedWith  []interface{}
}

// Get invokes the mock implementation and marks it as invoked
func (s *Store) Get(key []byte) []byte {
	s.GetInvoked = true
	s.GetInvokedCount++
	s.GetInvokedWith = append(s.GetInvokedWith, key)
	return s.GetFn(key)
}

// Set invokes the mock implementation and marks it as invoked
func (s *Store) Set(key []byte, value []byte) {
	s.SetInvoked = true
	s.SetInvokedCount++
	s.SetInvokedWith = append(s.SetInvokedWith, []interface{}{key, value})
	s.SetFn(key, value)
}

// Delete invokes the mock implementation and marks it as invoked
func (s *Store) Delete(key []byte) {
	s.DeleteInvoked = true
	s.DeleteInvokedCount++
	s.DeleteInvokedWith = append(s.DeleteInvokedWith, key)
	s.DeleteFn(key)
}
