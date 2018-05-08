package types

// KeyStoreUpdate update represents a single update in a batch of updates
type KeyStoreUpdate struct {
	Key   []byte
	Value []byte
}

// KeyStoreUpdateBatch is an ordered list of KeyStoreUpdate
type KeyStoreUpdateBatch []KeyStoreUpdate

// KeyStore is an interface for a backing keystore
type KeyStore interface {
	Get(key []byte) (value []byte, err error)
	Set(key []byte, value []byte) error
	Delete(key []byte) error
	Batch()
	Write() error
	BatchUpdate(updates KeyStoreUpdateBatch) error
}
