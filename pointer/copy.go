package pointer

// Copy a ptr as new
func Copy[T any](p *T) *T {
	c := SafeDeref(p)
	return &c
}
