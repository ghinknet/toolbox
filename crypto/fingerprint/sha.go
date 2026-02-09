package fingerprint

// GetExecutableSHA1 returns the SHA1 hash of the executable
func GetExecutableSHA1() (string, error) {
	return getExecutableHash(SHA1)
}

// GetExecutableSHA256 returns the SHA256 hash of the executable
func GetExecutableSHA256() (string, error) {
	return getExecutableHash(SHA256)
}

// GetExecutableSHA512 returns the SHA512 hash of the executable
func GetExecutableSHA512() (string, error) {
	return getExecutableHash(SHA512)
}
