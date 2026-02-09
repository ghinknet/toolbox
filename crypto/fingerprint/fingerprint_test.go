package fingerprint

import "testing"

func TestGetExecutableHash(t *testing.T) {
	for _, hashType := range []hashType{MD5, SHA1, SHA256, SHA512} {
		hash, err := getExecutableHash(hashType)
		if err != nil {
			t.Errorf("getExecutableHash(%s) failed: %v", hashType, err)
		}
		if hash == "" {
			t.Errorf("getExecutableHash(%s) returned empty hash", hashType)
		}
	}
}
