package lthash

import (
	"encoding/hex"
	"testing"
)

// TestExpand_blake3Vectors checks the expansion against official BLAKE3 test vectors (the unkeyed "hash" cases, 131 bytes of XOF output),
// so that a change in the BLAKE3 package (ex: an upgrade) can't change LtHash results unnoticed.
//
// The input of each case is 'InputLen' bytes, repeating 0, 1, 2, …, 250.
//
// Source: https://github.com/BLAKE3-team/BLAKE3/blob/master/test_vectors/test_vectors.json
func TestExpand_blake3Vectors(t *testing.T) {
	tests := []struct {
		InputLen int
		Expected string
	}{
		{
			InputLen: 0,
			Expected: "af1349b9f5f9a1a6a0404dea36dcc9499bcb25c9adc112b7cc9a93cae41f3262e00f03e7b69af26b7faaf09fcd333050338ddfe085b8cc869ca98b206c08243a26f5487789e8f660afe6c99ef9e0c52b92e7393024a80459cf91f476f9ffdbda7001c22e159b402631f277ca96f2defdf1078282314e763699a31c5363165421cce14d",
		},
		{
			InputLen: 1,
			Expected: "2d3adedff11b61f14c886e35afa036736dcd87a74d27b5c1510225d0f592e213c3a6cb8bf623e20cdb535f8d1a5ffb86342d9c0b64aca3bce1d31f60adfa137b358ad4d79f97b47c3d5e79f179df87a3b9776ef8325f8329886ba42f07fb138bb502f4081cbcec3195c5871e6c23e2cc97d3c69a613eba131e5f1351f3f1da786545e5",
		},
		{
			InputLen: 1023,
			Expected: "10108970eeda3eb932baac1428c7a2163b0e924c9a9e25b35bba72b28f70bd11a182d27a591b05592b15607500e1e8dd56bc6c7fc063715b7a1d737df5bad3339c56778957d870eb9717b57ea3d9fb68d1b55127bba6a906a4a24bbd5acb2d123a37b28f9e9a81bbaae360d58f85e5fc9d75f7c370a0cc09b6522d9c8d822f2f28f485",
		},
		{
			InputLen: 1024,
			Expected: "42214739f095a406f3fc83deb889744ac00df831c10daa55189b5d121c855af71cf8107265ecdaf8505b95d8fcec83a98a6a96ea5109d2c179c47a387ffbb404756f6eeae7883b446b70ebb144527c2075ab8ab204c0086bb22b7c93d465efc57f8d917f0b385c6df265e77003b85102967486ed57db5c5ca170ba441427ed9afa684e",
		},
		{
			InputLen: 1025,
			Expected: "d00278ae47eb27b34faecf67b4fe263f82d5412916c1ffd97c8cb7fb814b8444f4c4a22b4b399155358a994e52bf255de60035742ec71bd08ac275a1b51cc6bfe332b0ef84b409108cda080e6269ed4b3e2c3f7d722aa4cdc98d16deb554e5627be8f955c98e1d5f9565a9194cad0c4285f93700062d9595adb992ae68ff12800ab67a",
		},
		{
			InputLen: 3073,
			Expected: "7124b49501012f81cc7f11ca069ec9226cecb8a2c850cfe644e327d22d3e1cd39a27ae3b79d68d89da9bf25bc27139ae65a324918a5f9b7828181e52cf373c84f35b639b7fccbb985b6f2fa56aea0c18f531203497b8bbd3a07ceb5926f1cab74d14bd66486d9a91eba99059a98bd1cd25876b2af5a76c3e9eed554ed72ea952b603bf",
		},
	}

	for testNumber, test := range tests {
		var input []byte = make([]byte, test.InputLen)
		for i := range input {
			input[i] = byte(i % 251)
		}

		expanded := expand(string(input))

		expected, err := hex.DecodeString(test.Expected)
		if nil != err {
			t.Fatalf("For test #%d, did not expect an error but actually got one: (%T) %s", testNumber, err, err)
		}

		if actual := expanded[:len(expected)]; string(expected) != string(actual) {
			t.Errorf("For test #%d (input length %d), the actual expansion is not what was expected.", testNumber, test.InputLen)
			t.Logf("EXPECTED: %x", expected)
			t.Logf("ACTUAL:   %x", actual)
		}
	}
}
