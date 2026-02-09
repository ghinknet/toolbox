package fingerprint

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
	"os"
)

type hashType string

const (
	SHA1   hashType = "sha1"
	SHA256 hashType = "sha256"
	SHA512 hashType = "sha512"
	MD5    hashType = "md5"
)

// getExecutableHash returns the selected hash of the executable
func getExecutableHash(hashType hashType) (string, error) {
	// Get executable path
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	// Open file
	file, err := os.Open(exePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Calculate hash
	var hasher hash.Hash
	switch hashType {
	case SHA1:
		hasher = sha1.New()
	case SHA256:
		hasher = sha256.New()
	case SHA512:
		hasher = sha512.New()
	case MD5:
		hasher = md5.New()
	}
	if _, err = io.Copy(hasher, file); err != nil {
		return "", err
	}

	// Return hexadecimal string
	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}
