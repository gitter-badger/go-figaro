package rlp

// RLP Well-Known Types
//
// uint, uint8, uint16, uint32, uint64
// int, int8, int16, int32, int64
// []byte, string
// []interface{}
//
// RLP will use reflection for:
// slice, map, struct
//
// struct types will respect the `rpl:"-"` tag to skip encoding,
// otherwise it will serialize based on field order of public fields
//
// it will often be faster to serialize by converting to an []interface{}
// containing only well-known types

// Serializer prepares itself for RLP encoding by handling
// initial serialization into a well-known types
//
// This may be useful for complex types where reflection must
// be used by RLP otherwise
type Serializer interface {
	// SerializeRLP should return either a []byte or
	// an arbitrarily nested []interface{} of []byte for
	// final RLP encoding
	RLPSerialize() (interface{}, error)
}

// Deserializer hydrates itself after initial RLP decoding
//
// This should be capable of reversing the process of SerializeRLP
type Deserializer interface {
	// DeserializeRLP should expect either a []byte or
	// an arbitrarily nested []interface{} of []byte for
	// after initial RLP decoding
	RLPDeserialize(interface{}) error
}
