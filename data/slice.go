package data

// MakeSliceNotNil makes sure slice not nil
func MakeSliceNotNil[T any, S ~[]T](slice S) S {
	if slice == nil {
		return make(S, 0)
	}
	return slice
}
