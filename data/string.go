package data

import (
	"strconv"
	"strings"
)

// Atoi convert string to type int
func Atoi(s string) (int, error) {
	return strconv.Atoi(s)
}

// Atoi32 convert string to type int32
func Atoi32(s string) (int32, error) {
	i, err := strconv.Atoi(s)
	return int32(i), err
}

// Atoi64 convert string to type int64
func Atoi64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// StrConcat provides an easy method to concat strings directly
func StrConcat(strs ...string) string {
	return strings.Join(strs, "")
}
