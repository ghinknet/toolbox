package data

import "strconv"

// Atoi converted to type int
func Atoi(s string) (int, error) {
	return strconv.Atoi(s)
}

// Atoi32 converted to type int32
func Atoi32(s string) (int32, error) {
	i, err := strconv.Atoi(s)
	return int32(i), err
}

// Atoi64 converted to type int64
func Atoi64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
