package data

import (
	"math"
	"testing"
)

func TestMapKeys(t *testing.T) {
	// Test empty map
	emptyMap := make(map[string]int)
	if keys := MapKeys(emptyMap); len(keys) != 0 {
		t.Error("empty map should return empty slice", keys)
	}

	// Test map with string keys
	stringMap := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	keys := MapKeys(stringMap)
	if len(keys) != 3 {
		t.Error("keys length not match", len(keys))
	}

	// Verify all expected keys exist
	expectedKeys := map[string]bool{"a": true, "b": true, "c": true}
	for _, k := range keys {
		if !expectedKeys[k] {
			t.Error("unexpected key", k)
		}
	}

	// Test map with integer keys
	intMap := map[int]string{
		1: "one",
		2: "two",
		3: "three",
	}
	intKeys := MapKeys(intMap)
	if len(intKeys) != 3 {
		t.Error("int keys length not match", len(intKeys))
	}
}

func TestMapValues(t *testing.T) {
	// Test empty map
	emptyMap := make(map[int]string)
	if values := MapValues(emptyMap); len(values) != 0 {
		t.Error("empty map should return empty slice", values)
	}

	// Test map with string values
	stringMap := map[int]string{
		1: "apple",
		2: "banana",
		3: "cherry",
	}
	values := MapValues(stringMap)
	if len(values) != 3 {
		t.Error("values length not match", len(values))
	}

	// Verify all expected values exist
	expectedValues := map[string]bool{"apple": true, "banana": true, "cherry": true}
	for _, v := range values {
		if !expectedValues[v] {
			t.Error("unexpected value", v)
		}
	}

	// Test map with duplicate values
	duplicateMap := map[string]int{
		"a": 10,
		"b": 10,
		"c": 20,
	}
	dupValues := MapValues(duplicateMap)
	if len(dupValues) != 3 {
		t.Error("duplicate values length not match", len(dupValues))
	}
}

func TestMapKeysValues(t *testing.T) {
	// Test empty map
	emptyMap := make(map[string]float64)
	keys, values := MapKeysValues(emptyMap)
	if len(keys) != 0 || len(values) != 0 {
		t.Error("empty map should return empty slices", keys, values)
	}

	// Test map with string keys and int values
	testMap := map[string]int{
		"x": 100,
		"y": 200,
		"z": 300,
	}
	keys, valuesInt := MapKeysValues(testMap)

	// Verify consistent lengths
	if len(keys) != len(valuesInt) {
		t.Error("keys and values length mismatch", len(keys), len(valuesInt))
	}

	if len(keys) != 3 {
		t.Error("keys length not match", len(keys))
	}

	// Verify key-value correspondence
	for i, key := range keys {
		expectedValue, exists := testMap[key]
		if !exists {
			t.Error("key not found in original map", key)
		}
		if valuesInt[i] != expectedValue {
			t.Error("value mismatch for key", key, valuesInt[i], expectedValue)
		}
	}

	// Test map with integer keys and struct values
	structMap := map[int]struct{ name string }{
		1: {name: "Alice"},
		2: {name: "Bob"},
	}
	intKeys, structValues := MapKeysValues(structMap)
	if len(intKeys) != 2 || len(structValues) != 2 {
		t.Error("struct map keys/values length not match", len(intKeys), len(structValues))
	}

	// Test map with float64 values
	floatMap := map[string]float64{
		"pi": 3.14159,
		"e":  2.71828,
	}
	strKeys, floatValues := MapKeysValues(floatMap)
	if len(strKeys) != 2 || len(floatValues) != 2 {
		t.Error("float map keys/values length not match", len(strKeys), len(floatValues))
	}

	// Verify float values with tolerance
	expectedFloatValues := map[string]float64{
		"pi": 3.14159,
		"e":  2.71828,
	}
	for i, key := range strKeys {
		expected, exists := expectedFloatValues[key]
		if !exists {
			t.Error("key not found in float map", key)
		}
		// Compare float values with tolerance
		if math.Abs(floatValues[i]-expected) > 0.00001 {
			t.Error("float value mismatch for key", key, floatValues[i], expected)
		}
	}
}
