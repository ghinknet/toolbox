package pointer

import (
	"maps"
	"slices"
	"testing"
)

func TestSafeDeref(t *testing.T) {
	// Test cases
	var (
		// String
		exampleStringContent            = "homo"
		exampleStringContentPtr         = &exampleStringContent
		exampleStringPtr        *string = nil
		// Int
		exampleIntContent         = 114514
		exampleIntContentPtr      = &exampleIntContent
		exampleIntPtr        *int = nil
		// Int 64
		exampleInt64Content    int64  = 114514
		exampleInt64ContentPtr        = &exampleInt64Content
		exampleInt64Ptr        *int64 = nil
		// Bool
		exampleBoolContent          = true
		exampleBoolContentPtr       = &exampleBoolContent
		exampleBoolPtr        *bool = nil
		// Array
		exampleArrayContent               = [5]string{"life", "blooms", "like", "a", "flower"}
		exampleArrayContentPtr            = &exampleArrayContent
		exampleArrayPtr        *[5]string = nil
		// Slice
		exampleSliceContent              = []string{"tech", "otakus", "save", "the", "world"}
		exampleSliceContentPtr           = &exampleSliceContent
		exampleSlicePtr        *[]string = nil
		// Map
		exampleMapContent                       = map[string]string{"ghink": "Geek the Think"}
		exampleMapContentPtr                    = &exampleMapContent
		exampleMapPtr        *map[string]string = nil
	)

	// String
	if SafeDeref(exampleStringContentPtr) != exampleStringContent {
		t.Error("string ptr failed", SafeDeref(exampleStringContentPtr))
	}
	if SafeDeref(exampleStringPtr) != "" {
		t.Error("string ptr failed", SafeDeref(exampleStringContentPtr))
	}

	// Int
	if SafeDeref(exampleIntContentPtr) != exampleIntContent {
		t.Error("int ptr failed", SafeDeref(exampleIntContentPtr))
	}
	if SafeDeref(exampleIntPtr) != 0 {
		t.Error("int ptr failed", SafeDeref(exampleIntPtr))
	}

	// Int64
	if SafeDeref(exampleInt64ContentPtr) != exampleInt64Content {
		t.Error("int64 ptr failed", SafeDeref(exampleInt64ContentPtr))
	}
	if SafeDeref(exampleInt64Ptr) != 0 {
		t.Error("int64 ptr failed", SafeDeref(exampleInt64Ptr))
	}

	// Bool
	if SafeDeref(exampleBoolContentPtr) != exampleBoolContent {
		t.Error("bool ptr failed", SafeDeref(exampleBoolContentPtr))
	}
	if SafeDeref(exampleBoolPtr) != false {
		t.Error("bool ptr failed", SafeDeref(exampleBoolPtr))
	}

	// Array
	if SafeDeref(exampleArrayContentPtr) != exampleArrayContent {
		t.Error("array ptr failed", SafeDeref(exampleArrayContentPtr))
	}
	if SafeDeref(exampleArrayPtr) != [5]string{} {
		t.Error("array ptr failed", SafeDeref(exampleArrayPtr))
	}

	// Slice
	if !slices.Equal(SafeDeref(exampleSliceContentPtr), exampleSliceContent) {
		t.Error("slice ptr failed", SafeDeref(exampleSliceContentPtr))
	}
	if !slices.Equal(SafeDeref(exampleSlicePtr), make([]string, 0)) {
		t.Error("slice ptr failed", SafeDeref(exampleSlicePtr))
	}

	// Map
	if !maps.Equal(SafeDeref(exampleMapContentPtr), exampleMapContent) {
		t.Error("map ptr failed", SafeDeref(exampleMapContentPtr))
	}
	if !maps.Equal(SafeDeref(exampleMapPtr), make(map[string]string)) {
		t.Error("map ptr failed", SafeDeref(exampleMapPtr))
	}
}
