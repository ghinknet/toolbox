package data

import "strconv"

// Itoa convert int to type string
func Itoa(i int) string {
	return strconv.Itoa(i)
}

// I32toa convert int32 to type string
func I32toa(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

// I64toa convert int64 to type string
func I64toa(i int64) string {
	return strconv.FormatInt(i, 10)
}
