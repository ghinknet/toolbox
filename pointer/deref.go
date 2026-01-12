package pointer

// SafeDeref provides a safe method to deref an any type value
func SafeDeref[T any](ptr *T) T {
	if ptr == nil {
		var zero T
		return zero
	}
	return *ptr
}
