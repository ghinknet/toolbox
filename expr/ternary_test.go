package expr

import "testing"

func TestTernary(t *testing.T) {
	if Ternary(true, "Hello", "World") == "World" {
		t.Error("string not match", Ternary(true, "Hello", "World"))
	}
}
