package pointer

// Copy a ptr as new
func Copy[T any](p *T) *T {
	if p == nil {
		return nil
	}
	c := *p
	return &c
}
