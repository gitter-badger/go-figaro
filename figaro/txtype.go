// Package figaro is the main package for go-figaro
package figaro

import "errors"

var (
	// ErrInvalidTxTypeData is a self-explantory error.
	ErrInvalidTxTypeData = errors.New("figaro tx: invalid TxType data")
)

// TxType is a supported transaction type. It is used to save
// space by only including a single `Value` field and a single
// byte to represent the type of value being transferred.
type TxType byte

const (
	// BalanceTx transactions transfer Fia Balance from one account to another.
	BalanceTx TxType = iota
	// StakeTx transactions transfer FIG Stake from one account to another.
	StakeTx
)

// ValidTxType is returns whether a TxType is a valid TxType
func ValidTxType(t TxType) bool {
	switch t {
	case BalanceTx:
		return true
	case StakeTx:
		return true
	default:
		return false
	}
}

// MarshalBinary implements encoding.BinaryMarshaler
func (tx TxType) MarshalBinary() ([]byte, error) {
	return []byte{byte(tx)}, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (tx *TxType) UnmarshalBinary(b []byte) error {
	if len(b) > 1 {
		return ErrInvalidTxTypeData
	}
	t := TxType(b[0])
	if !ValidTxType(t) {
		return ErrInvalidTxTypeData
	}
	*tx = t
	return nil
}
