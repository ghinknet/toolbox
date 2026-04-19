package pointer

import "testing"

func TestCopy(t *testing.T) {
	value := "something"
	valuePtrA := &value
	valuePtrB := Copy(valuePtrA)
	if valuePtrB == nil {
		t.Fatal("value pointer cannot be nil")
	}
	if valuePtrA == valuePtrB {
		t.Fatal("value pointers cannot be equal")
	}
	if *valuePtrA != *valuePtrB {
		t.Fatal("value content mismatch")
	}
}
