package data

// MapKeys returns a slice of all keys in the map
func MapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// MapValues returns a slice of all values in the map
func MapValues[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// MapKeysValues returns slices of all keys and values in the map
func MapKeysValues[K comparable, V any](m map[K]V) ([]K, []V) {
	keys := make([]K, 0, len(m))
	values := make([]V, 0, len(m))

	for k, v := range m {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}
