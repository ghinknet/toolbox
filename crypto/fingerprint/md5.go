package fingerprint

// GetExecutableMD5 returns the MD5 hash of the executable
func GetExecutableMD5() (string, error) {
	return getExecutableHash(MD5)
}
