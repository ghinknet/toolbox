package data

import "strings"

// TrimPrefixMapsString trims the prefix from all values in the map that have a matching key in the prefix map
func TrimPrefixMapsString[K comparable](m map[K]string, prefix map[K]string) map[K]string {
	result := make(map[K]string)

	for k, v := range m {
		if val, exists := prefix[k]; exists {
			result[k] = strings.TrimPrefix(v, val)
		} else {
			result[k] = v
		}
	}

	return result
}

// TrimSuffixMapsString trims the suffix from all values in the map that have a matching key in the suffix map
func TrimSuffixMapsString[K comparable](m map[K]string, suffix map[K]string) map[K]string {
	result := make(map[K]string)

	for k, v := range m {
		if val, exists := suffix[k]; exists {
			result[k] = strings.TrimSuffix(v, val)
		} else {
			result[k] = v
		}
	}

	return result
}
