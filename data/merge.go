package data

// MergeMapsString merges values of two maps of strings
func MergeMapsString[K comparable, V string](a, b map[K]V) map[K]V {
	result := make(map[K]V)

	// Copy map a
	for k, v := range a {
		result[k] = v
	}

	// Merge map b
	for k, v := range b {
		if val, exists := result[k]; exists {
			result[k] = val + v
		} else {
			result[k] = v
		}
	}

	return result
}

// MergeMapsInt merges values of two maps of ints
func MergeMapsInt[K comparable, V int](a, b map[K]V) map[K]V {
	result := make(map[K]V)

	// Copy map a
	for k, v := range a {
		result[k] = v
	}

	// Merge map b
	for k, v := range b {
		if val, exists := result[k]; exists {
			result[k] = val + v
		} else {
			result[k] = v
		}
	}

	return result
}
