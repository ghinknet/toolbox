package data

// hashBytes covers common cryptographic hash output types.
// ~ allows named types (e.g. type Digest [32]byte).
type hashBytes interface {
	~[16]byte | ~[20]byte | ~[28]byte | ~[32]byte | ~[48]byte | ~[64]byte | ~[]byte
}

// HexString converts a hash byte array or slice to a lower-case hex string.
// It accepts every common hash size (MD5, SHA-1, SHA-256/224, SHA-384, SHA-512),
// named types, and plain []byte.
func HexString[T hashBytes](h T) string {
	n := len(h)
	buf := make([]byte, n*2)
	for i := 0; i < n; i++ {
		b := h[i]
		buf[i*2] = encHex[b>>4]
		buf[i*2+1] = encHex[b&0xf]
	}
	return string(buf)
}

var encHex = [16]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}
