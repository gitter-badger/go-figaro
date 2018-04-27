package memory

import "fmt"

// Store sets up an in-memory key/value store
type Store struct {
	db map[string][]byte
}

// Get returns a trie value given a trie key
func (s *Store) Get(key []byte) []byte {
	v := s.db[string(key)]
	if v == nil {
		return nil
	}
	c := make([]byte, len(v))
	copy(c, v)
	return c
}

// Set updates a trie key with a trie value
func (s *Store) Set(key []byte, value []byte) {
	if value != nil {
		c := make([]byte, len(value))
		copy(c, value)
		s.db[string(key)] = c
	}
}

// Delete removes a trie key/value
func (s *Store) Delete(key []byte) {
	delete(s.db, string(key))
}

func (s *Store) String() string {
	var str string
	for k, v := range s.db {
		str = str + fmt.Sprintf("% x:% x\n", k, v)
	}
	return str
}

// Len returns the number of entries in the store
func (s *Store) Len() int {
	return len(s.db)
}

// NewStore returns an initialized *Store
func NewStore() *Store {
	db := make(map[string][]byte)
	return &Store{db: db}
}
