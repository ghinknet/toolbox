package expr

// Ternary provides an approximate ternary expression alternative function
func Ternary[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}
