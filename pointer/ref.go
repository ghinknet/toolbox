package pointer

// Ref returns a pointer to the given value
func Ref[T any](v T) *T {
	return &v
}
