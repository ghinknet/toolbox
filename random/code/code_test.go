package code

import (
	"strings"
	"testing"
)

// TestCodeDefault verifies default behavior when Options has zero values.
func TestCodeDefault(t *testing.T) {
	code, err := Code(Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(code) != 6 {
		t.Errorf("default length: expected 6, got %d (%q)", len(code), code)
	}
	// Should contain only digits because no category was enabled and fallback activated.
	for _, c := range code {
		if c < '0' || c > '9' {
			t.Errorf("expected only digits by default, got %q", code)
			break
		}
	}
}

// TestCodeLength ensures Digit is respected.
func TestCodeLength(t *testing.T) {
	for _, length := range []int{0, 1, 4, 10, 20} {
		code, err := Code(Options{Digit: length, UseNumbers: true})
		if err != nil {
			t.Errorf("length %d: unexpected error: %v", length, err)
			continue
		}
		expected := length
		if length <= 0 {
			expected = 6 // default
		}
		if len(code) != expected {
			t.Errorf("length %d: expected %d, got %d", length, expected, len(code))
		}
	}
}

// TestCodeNumbersOnly verifies Number convenience function and UseNumbers.
func TestCodeNumbersOnly(t *testing.T) {
	code := Number(8)
	if len(code) != 8 {
		t.Errorf("Number(8) length: expected 8, got %d", len(code))
	}
	for _, c := range code {
		if c < '0' || c > '9' {
			t.Fatalf("Number(8) contains non-digit: %q", code)
		}
	}

	// Also directly via Code
	code, err := Code(Options{Digit: 5, UseNumbers: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(code) != 5 {
		t.Errorf("expected length 5, got %d", len(code))
	}
}

// TestCodeAlpha verifies Alpha convenience and case controls.
func TestCodeAlpha(t *testing.T) {
	// Lower only
	lower := Alpha(10, true, false)
	if len(lower) != 10 {
		t.Errorf("Alpha length expected 10, got %d", len(lower))
	}
	for _, c := range lower {
		if c < 'a' || c > 'z' {
			t.Fatalf("expected only lowercase, got %q", lower)
		}
	}

	// Upper only
	upper := Alpha(12, false, true)
	for _, c := range upper {
		if c < 'A' || c > 'Z' {
			t.Fatalf("expected only uppercase, got %q", upper)
		}
	}

	// Both cases
	both := Alpha(15, true, true)
	for _, c := range both {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
			t.Fatalf("expected only letters, got %q", both)
		}
	}
}

// TestCodeMixed explores Mixed shortcut.
func TestCodeMixed(t *testing.T) {
	code := Mixed(12, true, true, true, true)
	if len(code) != 12 {
		t.Errorf("Mixed length expected 12, got %d", len(code))
	}
	// With all categories enabled, it should not panic.
}

// TestCodeSymbols checks symbol inclusion and custom symbols.
func TestCodeSymbols(t *testing.T) {
	code, err := Code(Options{
		Digit:      20,
		UseSymbols: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.ContainsAny(code, DefaultSymbols) {
		t.Errorf("expected at least one default symbol in %q", code)
	}

	// Custom symbols replace default set
	custom := "@#"
	customCode, err := Code(Options{
		Digit:         30,
		UseSymbols:    true,
		CustomSymbols: custom,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.ContainsAny(customCode, custom) {
		t.Errorf("expected custom symbols in %q", customCode)
	}
	// Should not contain any default symbol that is not in custom set
	if strings.ContainsAny(customCode, "$%^") {
		t.Errorf("custom symbols should replace default set, got %q", customCode)
	}
}

// TestCodeExclude checks that excluded characters never appear.
func TestCodeExclude(t *testing.T) {
	exclude := "abc123"
	code, err := Code(Options{
		Digit:        50,
		UseNumbers:   true,
		UseLowercase: true,
		ExcludeChars: exclude,
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.ContainsAny(code, exclude) {
		t.Errorf("excluded characters found in %q", code)
	}
}

// TestRequireEach verifies that every enabled category appears at least once.
func TestRequireEach(t *testing.T) {
	config := Options{
		Digit:        20,
		UseNumbers:   true,
		UseLowercase: true,
		UseUppercase: true,
		UseSymbols:   true,
		RequireEach:  true,
	}
	code, err := Code(config)
	if err != nil {
		t.Fatal(err)
	}

	hasDigit := strings.ContainsAny(code, "0123456789")
	hasLower := strings.ContainsAny(code, "abcdefghijklmnopqrstuvwxyz")
	hasUpper := strings.ContainsAny(code, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	hasSymbol := strings.ContainsAny(code, DefaultSymbols)

	if !hasDigit || !hasLower || !hasUpper || !hasSymbol {
		t.Errorf("RequireEach failed: code=%q, digit=%v lower=%v upper=%v symbol=%v",
			code, hasDigit, hasLower, hasUpper, hasSymbol)
	}
}

// TestRequireEachWithExclude checks that RequireEach still works when exclusions shrink some categories.
func TestRequireEachWithExclude(t *testing.T) {
	code, err := Code(Options{
		Digit:        20,
		UseNumbers:   true,
		UseLowercase: true,
		UseUppercase: true,
		UseSymbols:   true,
		RequireEach:  true,
		ExcludeChars: "abcdef", // remove some lowercase and possibly symbols
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.ContainsAny(code, "0123456789") {
		t.Error("expected at least one digit")
	}
	// Lowercases might not include a..f, but we expect at least one lowercase from g..z
	hasLower := false
	for _, c := range code {
		if c >= 'g' && c <= 'z' {
			hasLower = true
			break
		}
	}
	if !hasLower {
		t.Errorf("expected at least one remaining lowercase in %q", code)
	}
}

// TestRequireEachEdgeCase tests short length with many required categories.
func TestRequireEachShortLength(t *testing.T) {
	// Request length shorter than number of required categories
	code, err := Code(Options{
		Digit:        2,
		UseNumbers:   true,
		UseLowercase: true,
		UseUppercase: true,
		UseSymbols:   true,
		RequireEach:  true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(code) != 2 {
		t.Errorf("expected length 2, got %d", len(code))
	}
	// Should not panic; only as many categories as length can be filled.
}

// TestMustCodePanic verifies that MustCode panics on impossible configuration.
// We can trigger an error by breaking randomInt indirectly, but that's hard.
// Instead, we rely on MustCode's happy path regularly.
func TestMustCode(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("MustCode panicked unexpectedly: %v", r)
		}
	}()
	code := MustCode(Options{Digit: 10, UseNumbers: true})
	if len(code) != 10 {
		t.Errorf("MustCode: unexpected length %d", len(code))
	}
}

// TestStrongPassword verifies the StrongPassword convenience.
func TestStrongPassword(t *testing.T) {
	pwd := StrongPassword()
	if len(pwd) != 16 {
		t.Errorf("StrongPassword length: expected 16, got %d", len(pwd))
	}
	// Check presence of all categories
	if !strings.ContainsAny(pwd, "0123456789") {
		t.Error("StrongPassword missing digits")
	}
	if !strings.ContainsAny(pwd, "abcdefghijklmnopqrstuvwxyz") {
		t.Error("StrongPassword missing lowercase")
	}
	if !strings.ContainsAny(pwd, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		t.Error("StrongPassword missing uppercase")
	}
	if !strings.ContainsAny(pwd, DefaultSymbols) {
		t.Error("StrongPassword missing symbols")
	}
	// Excluded ambiguous characters
	ambiguous := "0OIl1"
	if strings.ContainsAny(pwd, ambiguous) {
		t.Errorf("StrongPassword contains ambiguous characters: %q", pwd)
	}
}

// TestCodeErrorHandling simulates error from randomInt by using zero-length charset?
// But we can't easily force crypto/rand to fail. We trust the error return paths.
// We can test that empty charSet fallback still works without error.
func TestCodeEmptyCharset(t *testing.T) {
	// If no categories are enabled and fallback triggers, code should be digits.
	code, err := Code(Options{Digit: 4, UseNumbers: false})
	if err != nil {
		t.Fatal(err)
	}
	if len(code) != 4 {
		t.Errorf("expected 4, got %d", len(code))
	}
	for _, c := range code {
		if c < '0' || c > '9' {
			t.Errorf("fallback should be digits, got %q", code)
			break
		}
	}
}

// TestCaching ensures that buildCharSet returns the same result for identical configs.
func TestCaching(t *testing.T) {
	opts := Options{
		UseNumbers:   true,
		UseLowercase: true,
		ExcludeChars: "xyz",
	}
	set1 := buildCharSet(opts)
	set2 := buildCharSet(opts)
	if set1 != set2 {
		t.Error("caching failed: different sets for same config")
	}
}

// TestRemoveChars verifies the exclusion logic.
func TestRemoveChars(t *testing.T) {
	result := removeChars("abcdef123", "abc")
	if strings.Contains(result, "a") || strings.Contains(result, "b") || strings.Contains(result, "c") {
		t.Errorf("removeChars failed, got %q", result)
	}
	if !strings.Contains(result, "def123") {
		t.Errorf("remaining part missing, got %q", result)
	}
}

// TestFilterCharset verifies filtering logic.
func TestFilterCharset(t *testing.T) {
	res := filterCharset("abc123", "a2")
	if len(res) == 0 {
		t.Fatal("filterCharset returned empty")
	}
	s := string(res)
	if strings.Contains(s, "a") || strings.Contains(s, "2") {
		t.Errorf("filter failed, got %q", s)
	}
	if !strings.Contains(s, "b") || !strings.Contains(s, "c") || !strings.Contains(s, "1") || !strings.Contains(s, "3") {
		t.Errorf("missing characters in %q", s)
	}
}

// TestShuffle verifies that shuffle changes order (non-deterministic but likely).
func TestShuffle(t *testing.T) {
	original := []byte("abcdefgh")
	dup := make([]byte, len(original))
	copy(dup, original)
	err := shuffle(dup)
	if err != nil {
		t.Fatal(err)
	}
	// Very unlikely to be exactly same order for 8 elements
	same := true
	for i := range original {
		if original[i] != dup[i] {
			same = false
			break
		}
	}
	if same {
		t.Log("shuffle did not change order (unlikely but possible)")
	}
}

// TestRandomIntInvalidMax ensures error for non-positive max.
func TestRandomIntInvalidMax(t *testing.T) {
	_, err := randomInt(0)
	if err == nil {
		t.Error("expected error for max=0")
	}
	_, err = randomInt(-5)
	if err == nil {
		t.Error("expected error for negative max")
	}
}

// BenchmarkCode provides a performance baseline.
func BenchmarkCode(b *testing.B) {
	config := Options{
		Digit:        16,
		UseNumbers:   true,
		UseLowercase: true,
		UseUppercase: true,
		UseSymbols:   true,
		RequireEach:  true,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Code(config)
		if err != nil {
			b.Fatal(err)
		}
	}
}
