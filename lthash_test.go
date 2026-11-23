package lthash_test

import (
	"encoding/hex"
	"testing"

	"codeberg.org/reiver/go-erorr"

	"github.com/reiver/go-lthash"
)

func TestLtHash_empty(t *testing.T) {
	var hash lthash.LtHash

	if !hash.IsEmpty() {
		t.Errorf("Expected the zero value to be empty.")
	}

	digest := hash.Digest()

	// The SHA-256 of 2048 zero bytes.
	const expected string = "e5a00aa9991ac8a5ee3109844d84a55583bd20572ad3ffcd42792f3c36b183ad"
	if actual := hex.EncodeToString(digest[:]); expected != actual {
		t.Errorf("The actual digest of the empty set is not what was expected.")
		t.Logf("EXPECTED: %s", expected)
		t.Logf("ACTUAL:   %s", actual)
	}
}

// The same known value as the ATProto reference implementation's test
// (packages/space/src/lthash.test.ts, at bluesky-social/atproto 9ded207).
func TestLtHash_knownDigest(t *testing.T) {
	var hash lthash.LtHash
	hash.Add("one")
	hash.Add("two")

	digest := hash.Digest()

	const expected string = "ae05cb6d224379d9710c290c8529945c5b0e0fde9ead30b9699057ce701c63e7"
	if actual := hex.EncodeToString(digest[:]); expected != actual {
		t.Errorf("The actual digest is not what was expected.")
		t.Logf("EXPECTED: %s", expected)
		t.Logf("ACTUAL:   %s", actual)
	}
}

func TestLtHash_addThenRemove(t *testing.T) {
	var hash lthash.LtHash

	hash.Add("a")
	if hash.IsEmpty() {
		t.Fatalf("Did not expect the set to be empty after an add.")
	}

	hash.Remove("a")
	if !hash.IsEmpty() {
		t.Errorf("Expected the set to be empty after an add then a remove.")
	}
}

// Removing first (the state wraps below zero), then adding, gets back to empty too.
func TestLtHash_removeThenAdd(t *testing.T) {
	var hash lthash.LtHash

	hash.Remove("a")
	hash.Add("a")
	if !hash.IsEmpty() {
		t.Errorf("Expected the set to be empty after a remove then an add.")
	}
}

func TestLtHash_orderIndependent(t *testing.T) {
	var elements = []string{
		"com.example.post/3aaa/bafyreiaaa",
		"com.example.post/3bbb/bafyreibbb",
		"com.example.like/3ccc/bafyreiccc",
		"",
		"ünïcödé/🪞/x",
	}

	var forward lthash.LtHash
	for _, element := range elements {
		forward.Add(element)
	}

	var backward lthash.LtHash
	for i := len(elements) - 1; 0 <= i; i-- {
		backward.Add(elements[i])
	}

	if !forward.Equal(backward) {
		t.Errorf("Expected the same set, added in a different order, to have the same state.")
	}
	if forward.Digest() != backward.Digest() {
		t.Errorf("Expected the same set, added in a different order, to have the same digest.")
	}
}

func TestLtHash_distinguishes(t *testing.T) {
	var a, b lthash.LtHash
	a.Add("a")
	b.Add("b")

	if a.Equal(b) {
		t.Errorf("Expected different elements to give different states.")
	}
}

func TestLtHash_multiset(t *testing.T) {
	var twice lthash.LtHash
	twice.Add("a")
	twice.Add("a")

	if twice.IsEmpty() {
		t.Fatalf("Did not expect adding an element twice to cancel out.")
	}

	var once lthash.LtHash
	once.Add("a")

	twice.Remove("a")
	if !twice.Equal(once) {
		t.Errorf("Expected adding twice then removing once to equal adding once.")
	}
}

// An update — the record at a path changes from one CID to another — is a remove plus an add.
func TestLtHash_update(t *testing.T) {
	var incremental lthash.LtHash
	incremental.Add("com.example.post/3aaa/bafyreiold")
	incremental.Add("com.example.post/3bbb/bafyreibbb")
	incremental.Remove("com.example.post/3aaa/bafyreiold")
	incremental.Add("com.example.post/3aaa/bafyreinew")

	var fromScratch lthash.LtHash
	fromScratch.Add("com.example.post/3bbb/bafyreibbb")
	fromScratch.Add("com.example.post/3aaa/bafyreinew")

	if !incremental.Equal(fromScratch) {
		t.Errorf("Expected the incrementally updated hash to equal the hash computed from scratch.")
	}
}

func TestFromState(t *testing.T) {
	var hash lthash.LtHash
	hash.Add("a")
	hash.Add("b")

	resumed, err := lthash.FromState(hash.State())
	if nil != err {
		t.Fatalf("Did not expect an error but actually got one: (%T) %s", err, err)
	}
	if !resumed.Equal(hash) {
		t.Errorf("Expected the resumed hash to equal the original.")
	}

	resumed.Remove("b")
	var justA lthash.LtHash
	justA.Add("a")
	if !resumed.Equal(justA) {
		t.Errorf("Expected the resumed hash to keep working.")
	}
}

func TestFromState_badSize(t *testing.T) {
	for _, size := range []int{0, 1, 32, 2047, 2049, 4096} {
		var state []byte
		if 0 < size {
			state = make([]byte, size)
		}

		_, err := lthash.FromState(state)
		if !erorr.Is(err, lthash.ErrBadStateSize) {
			t.Errorf("For size %d, expected ErrBadStateSize but actually got: %v", size, err)
		}
	}
}

func TestLtHash_noAliasing(t *testing.T) {
	var state []byte = make([]byte, lthash.StateSize)
	state[0] = 0xff

	hash, err := lthash.FromState(state)
	if nil != err {
		t.Fatalf("Did not expect an error but actually got one: (%T) %s", err, err)
	}

	state[0] = 0x00
	if 0xff != hash.State()[0] {
		t.Errorf("Expected the hash not to alias the state it was made from.")
	}

	handedOut := hash.State()
	handedOut[0] ^= 0xff
	if 0xff != hash.State()[0] {
		t.Errorf("Expected the hash not to alias the state it handed out.")
	}

	// Copying an LtHash copies its state.
	var copied lthash.LtHash = hash
	copied.Add("a")
	if copied.Equal(hash) {
		t.Errorf("Expected a copy to be independent of the original.")
	}
}

func TestLtHash_nilReceiver(t *testing.T) {
	defer func() {
		if r := recover(); nil != r {
			t.Errorf("Did not expect a panic but actually got one: %v", r)
		}
	}()

	var hash *lthash.LtHash
	hash.Add("a")
	hash.Remove("a")
}
