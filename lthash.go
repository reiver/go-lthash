package lthash

import (
	"crypto/sha256"
	"encoding/binary"

	"codeberg.org/reiver/go-erorr"
	"codeberg.org/reiver/go-field"
	"lukechampine.com/blake3"
)

const (
	// Lanes is the number of 16-bit lanes in an LtHash.
	Lanes = 1024

	// StateSize is the size, in bytes, of an LtHash's state.
	StateSize = Lanes * 2 // 2048

	// DigestSize is the size, in bytes, of an LtHash's digest.
	DigestSize = sha256.Size // 32
)

// LtHash is a homomorphic (multi-)set hash.
//
// The zero value is the empty set.
//
// An LtHash is a value (an array, not a slice): copying one copies its state.
type LtHash struct {
	state [StateSize]byte
}

// FromState returns the LtHash whose state is 'state' (as returned by [LtHash.State]).
//
// It returns [ErrBadStateSize] if 'state' is not 2048 bytes. (A nil or empty 'state' is an error too — for the empty set, use the zero value.)
func FromState(state []byte) (LtHash, error) {
	var hash LtHash

	if StateSize != len(state) {
		return hash, erorr.Wrap(ErrBadStateSize, "bad state size", field.Int("size", len(state)))
	}

	copy(hash.state[:], state)
	return hash, nil
}

// Add adds an element to the set.
//
// Adding an element twice counts it twice (it is a multi-set).
func (receiver *LtHash) Add(element string) {
	receiver.combine(element, +1)
}

// Remove removes an element from the set.
//
// (Removing an element that was never added does not fail — but the state no longer stands for any set; adding it back undoes it.)
func (receiver *LtHash) Remove(element string) {
	receiver.combine(element, -1)
}

func (receiver *LtHash) combine(element string, sign int) {
	if nil == receiver {
		return
	}

	var expanded [StateSize]byte = expand(element)

	for i := 0; i < StateSize; i += 2 {
		var lane uint16 = binary.LittleEndian.Uint16(receiver.state[i:])
		var delta uint16 = binary.LittleEndian.Uint16(expanded[i:])

		if 0 < sign {
			lane += delta
		} else {
			lane -= delta
		}

		binary.LittleEndian.PutUint16(receiver.state[i:], lane)
	}
}

// expand expands an element into 2048 bytes, with BLAKE3 in XOF mode.
func expand(element string) [StateSize]byte {
	var expanded [StateSize]byte

	hasher := blake3.New(DigestSize, nil)
	hasher.Write([]byte(element)) // never returns an error
	hasher.XOF().Read(expanded[:]) // never returns an error (for a read this small)

	return expanded
}

// State returns (a copy of) the 2048-byte state — ex: to persist it. See [FromState].
func (receiver LtHash) State() []byte {
	var state []byte = make([]byte, StateSize)
	copy(state, receiver.state[:])
	return state
}

// Digest returns the SHA-256 of the state.
func (receiver LtHash) Digest() [DigestSize]byte {
	return sha256.Sum256(receiver.state[:])
}

// IsEmpty returns whether the state is all zeros — as it is for the empty set.
func (receiver LtHash) IsEmpty() bool {
	return [StateSize]byte{} == receiver.state
}

// Equal returns whether two LtHashes have the same state.
func (receiver LtHash) Equal(other LtHash) bool {
	return receiver.state == other.state
}
