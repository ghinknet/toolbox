package data

import (
	"maps"
	"testing"
)

func TestMergeMapsString(t *testing.T) {
	// Test cases
	var (
		emptyMap   = map[int]string{}
		singleMapA = map[int]string{1: "Hello"}
		singleMapB = map[int]string{2: "World"}
		overlapMap = map[int]string{1: "World", 2: "!"}
	)

	// Test case 1: Both maps empty
	result1 := MergeMapsString(emptyMap, emptyMap)
	if len(result1) != 0 {
		t.Error("empty + empty should result in empty map")
	}

	// Test case 2: First map empty, second map not empty
	result2 := MergeMapsString(emptyMap, singleMapB)
	if len(result2) != 1 || result2[2] != "World" {
		t.Error("empty + {2:World} should result in {2:World}")
	}

	// Test case 3: First map not empty, second map empty
	result3 := MergeMapsString(singleMapA, emptyMap)
	if len(result3) != 1 || result3[1] != "Hello" {
		t.Error("{1:Hello} + empty should result in {1:Hello}")
	}

	// Test case 4: No overlapping keys
	result4 := MergeMapsString(singleMapA, singleMapB)
	if len(result4) != 2 {
		t.Error("{1:Hello} + {2:World} should have 2 entries")
	}
	if result4[1] != "Hello" || result4[2] != "World" {
		t.Error("{1:Hello} + {2:World} should result in {1:Hello, 2:World}")
	}

	// Test case 5: Overlapping keys with string concatenation
	result5 := MergeMapsString(singleMapA, overlapMap)
	if len(result5) != 2 {
		t.Error("{1:Hello} + {1:World, 2:!} should have 2 entries")
	}
	if result5[1] != "HelloWorld" {
		t.Error("Overlapping key 1 should concatenate 'Hello' and 'World'")
	}
	if result5[2] != "!" {
		t.Error("Non-overlapping key 2 should be '!'")
	}

	// Test case 6: Multiple overlapping keys
	mapA := map[string]string{"a": "foo", "b": "bar"}
	mapB := map[string]string{"b": "baz", "c": "qux", "d": "quux"}
	result6 := MergeMapsString(mapA, mapB)
	expected6 := map[string]string{"a": "foo", "b": "barbaz", "c": "qux", "d": "quux"}
	if !maps.Equal(result6, expected6) {
		t.Error("Failed with multiple overlapping keys", result6)
	}

	// Test case 7: Test with different key types (float64)
	mapFloatA := map[float64]string{1.5: "one point five", 2.0: "two"}
	mapFloatB := map[float64]string{2.0: " point zero", 3.14: "pi"}
	result7 := MergeMapsString(mapFloatA, mapFloatB)
	if result7[1.5] != "one point five" {
		t.Error("Key 1.5 should be 'one point five'")
	}
	if result7[2.0] != "two point zero" {
		t.Error("Key 2.0 should concatenate to 'two point zero'")
	}
	if result7[3.14] != "pi" {
		t.Error("Key 3.14 should be 'pi'")
	}
}

func TestMergeMapsInt(t *testing.T) {
	// Test cases
	var (
		emptyMap   = map[string]int{}
		singleMapA = map[string]int{"a": 10}
		singleMapB = map[string]int{"b": 20}
		overlapMap = map[string]int{"a": 5, "c": 30}
	)

	// Test case 1: Both maps empty
	result1 := MergeMapsInt(emptyMap, emptyMap)
	if len(result1) != 0 {
		t.Error("empty + empty should result in empty map")
	}

	// Test case 2: First map empty, second map not empty
	result2 := MergeMapsInt(emptyMap, singleMapB)
	if len(result2) != 1 || result2["b"] != 20 {
		t.Error("empty + {b:20} should result in {b:20}")
	}

	// Test case 3: First map not empty, second map empty
	result3 := MergeMapsInt(singleMapA, emptyMap)
	if len(result3) != 1 || result3["a"] != 10 {
		t.Error("{a:10} + empty should result in {a:10}")
	}

	// Test case 4: No overlapping keys
	result4 := MergeMapsInt(singleMapA, singleMapB)
	if len(result4) != 2 {
		t.Error("{a:10} + {b:20} should have 2 entries")
	}
	if result4["a"] != 10 || result4["b"] != 20 {
		t.Error("{a:10} + {b:20} should result in {a:10, b:20}")
	}

	// Test case 5: Overlapping keys with integer addition
	result5 := MergeMapsInt(singleMapA, overlapMap)
	if len(result5) != 2 {
		t.Error("{a:10} + {a:5, c:30} should have 2 entries")
	}
	if result5["a"] != 15 {
		t.Error("Overlapping key 'a' should sum 10 and 5")
	}
	if result5["c"] != 30 {
		t.Error("Non-overlapping key 'c' should be 30")
	}

	// Test case 6: Test with negative numbers
	negativeMap := map[string]int{"a": -5, "b": -10}
	result6 := MergeMapsInt(singleMapA, negativeMap)
	if result6["a"] != 5 { // 10 + (-5) = 5
		t.Error("10 + (-5) should equal 5")
	}
	if result6["b"] != -10 {
		t.Error("Non-overlapping key 'b' should be -10")
	}

	// Test case 7: Multiple overlapping keys
	mapA := map[int]int{1: 100, 2: 200, 3: 300}
	mapB := map[int]int{2: 50, 3: -100, 4: 400}
	result7 := MergeMapsInt(mapA, mapB)
	expected7 := map[int]int{1: 100, 2: 250, 3: 200, 4: 400}
	if !maps.Equal(result7, expected7) {
		t.Error("Failed with multiple overlapping keys", result7)
	}

	// Test case 8: Test with different key types (bool)
	mapBoolA := map[bool]int{true: 1, false: 0}
	mapBoolB := map[bool]int{true: 2, false: 3}
	result8 := MergeMapsInt(mapBoolA, mapBoolB)
	if result8[true] != 3 { // 1 + 2 = 3
		t.Error("true key should sum to 3")
	}
	if result8[false] != 3 { // 0 + 3 = 3
		t.Error("false key should sum to 3")
	}
}
