/*
Package lthash implements LtHash — a homomorphic (multi-)set hash — as used by ATProto spaces, for the Go programming language.

An LtHash is 1024 lanes of 16-bit unsigned integers (2048 bytes). Each element is expanded, with BLAKE3 in XOF mode,
into 2048 bytes — read as 1024 little-endian 16-bit lanes — and added to (or subtracted from) the state, lane by lane, modulo 2¹⁶.
Because addition commutes, the state depends only on which elements are in the set (and how many times), not on the order they were added in.
So a set's hash can be kept up to date one element at a time, without re-hashing the whole set.

The digest of an LtHash is the SHA-256 of its 2048-byte state.
The empty set's state is 2048 zero bytes.

In an ATProto space repo, each record contributes the element "{collection}/{rkey}/{cid}",
and the commit's "data" is the digest of the repo's LtHash.

# Example

	var hash lthash.LtHash // the empty set

	hash.Add("com.example.post/3aaa/bafyrei...")
	hash.Add("com.example.post/3bbb/bafyrei...")
	hash.Remove("com.example.post/3aaa/bafyrei...")

	digest := hash.Digest()

	// Persist the state; later, pick up where it left off.
	state := hash.State()
	hash2, err := lthash.FromState(state)
*/
package lthash
