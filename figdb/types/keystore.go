package types

// TODO: add merkle root as type and export in `figdb` for use in `Account`, `Block`, etc

// KeyStoreUpdate update represents a single update in a batch of updates
type KeyStoreUpdate struct {
	Key   []byte
	Value []byte
}

// KeyStoreUpdateBatch is an ordered list of KeyStoreUpdate
type KeyStoreUpdateBatch []KeyStoreUpdate

// KeyStore is an interface for a backing keystore
type KeyStore interface {
	Get(key Key) (value []byte, err error)
	Set(key Key, value []byte) error
	Delete(key Key) error
	Batch()
	Write() error
	BatchUpdate(updates KeyStoreUpdateBatch) error
}

// Key is a []byte that can be conveniently converted to a string
type Key []byte

func (k Key) String() string {
	return string(k)
}

// Cache is a generic cache interface used by this package
type Cache interface {
	Add(k Key, v []byte)
	Get(k Key) (v []byte, ok bool)
}
