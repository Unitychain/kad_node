package identity

import (
	"crypto/sha256"
	"encoding/hex"
)

// Identity ...
type Identity struct {
	commitment string
}

// NewIdentity ...
func NewIdentity(commitment string) *Identity {
	return &Identity{commitment: commitment}
}

// Hash ...
type Hash []byte

// Byte ...
func (h Hash) Byte() []byte { return []byte(h) }

// Hex ...
func (h Hash) Hex() HashHex {
	return HashHex(hex.EncodeToString(h.Byte()))
}

// Index ...
type Index map[string]HashSet

// NewIndex ...
func NewIndex() Index {
	// TODO : workaround for cycle import
	// return Index(make(map[subject.HashHex]HashSet))
	return Index(make(map[string]HashSet))
}

// HashSet ...
type HashSet map[HashHex]string

// NewHashSet ...
func NewHashSet() HashSet {
	result := HashSet(make(map[HashHex]string))
	return result
}

// HashHex ...
type HashHex string

// String ...
func (h HashHex) String() string { return string(h) }

// Hash ...
func (i *Identity) Hash() *Hash {
	h := sha256.Sum256([]byte(i.commitment))
	result := Hash(h[:])
	return &result
}
