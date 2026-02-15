package data

import (
	"testing"
)

func TestMakeSliceNotNil(t *testing.T) {
	// Test case 1: nil slice should become empty slice
	var nilSlice []int
	result1 := MakeSliceNotNil(nilSlice)
	if result1 == nil {
		t.Error("nil slice should not remain nil")
	}
	if len(result1) != 0 {
		t.Error("nil slice should become empty slice with length 0")
	}

	// Test case 2: already initialized empty slice should remain unchanged
	emptySlice := make([]string, 0)
	result2 := MakeSliceNotNil(emptySlice)
	if len(result2) != 0 {
		t.Error("empty slice should remain empty")
	}
	if cap(result2) != cap(emptySlice) {
		t.Error("empty slice capacity should remain unchanged")
	}

	// Test case 3: non-empty slice should remain unchanged
	nonEmptySlice := []int{1, 2, 3, 4, 5}
	result3 := MakeSliceNotNil(nonEmptySlice)
	if len(result3) != 5 {
		t.Error("non-empty slice length should remain unchanged")
	}
	for i, v := range result3 {
		if v != nonEmptySlice[i] {
			t.Errorf("slice element at index %d should remain %d, got %d", i, nonEmptySlice[i], v)
		}
	}

	// Test case 4: slice with capacity but zero length should remain unchanged
	capacitySlice := make([]float64, 0, 10)
	result4 := MakeSliceNotNil(capacitySlice)
	if len(result4) != 0 {
		t.Error("zero-length slice should remain zero-length")
	}
	if cap(result4) != 10 {
		t.Error("slice capacity should remain 10")
	}

	// Test case 5: test with different types
	// String slice
	stringSlice := []string{"hello", "world"}
	result5 := MakeSliceNotNil(stringSlice)
	if len(result5) != 2 {
		t.Error("string slice length should remain 2")
	}

	// Struct slice
	type Person struct {
		Name string
		Age  int
	}
	personSlice := []Person{{Name: "Alice", Age: 30}, {Name: "Bob", Age: 25}}
	result6 := MakeSliceNotNil(personSlice)
	if len(result6) != 2 {
		t.Error("struct slice length should remain 2")
	}

	// Pointer slice
	ptrSlice := []*int{new(int), new(int)}
	result7 := MakeSliceNotNil(ptrSlice)
	if len(result7) != 2 {
		t.Error("pointer slice length should remain 2")
	}

	// Test case 6: verify that returned slice can be appended to
	nilSlice2 := []int(nil)
	result8 := MakeSliceNotNil(nilSlice2)
	result8 = append(result8, 42)
	if len(result8) != 1 {
		t.Error("returned slice should be appendable")
	}
	if result8[0] != 42 {
		t.Error("appended value should be 42")
	}

	// Test case 7: test with custom slice type
	type CustomSlice []int
	var customNilSlice CustomSlice
	result9 := MakeSliceNotNil(customNilSlice)
	if result9 == nil {
		t.Error("custom nil slice should not remain nil")
	}
	if len(result9) != 0 {
		t.Error("custom nil slice should become empty")
	}
}
