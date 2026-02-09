package pointer

// Ref returns a pointer to the given value
func Ref[T any](v T) *T {
	return &v
}

// SliceRef returns a slice of pointers to the given slice
func SliceRef[T any](v []T) []*T {
	ptrs := make([]*T, len(v))
	for i := range v {
		ptrs[i] = &v[i]
	}
	return ptrs
}

// SliceCopyRef returns a slice of pointers to a copy of the given slice
func SliceCopyRef[T any](v []T) []*T {
	ptrs := make([]*T, len(v))
	for i := range v {
		vcpy := v[i]
		ptrs[i] = &vcpy
	}
	return ptrs
}
