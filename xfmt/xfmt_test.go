package xfmt

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"testing"
)

type stringer struct{ s string }

func (s stringer) String() string { return "S(" + s.s + ")" }

type myError struct{ s string }

func (e myError) Error() string { return "err:" + e.s }

// parity asserts that Sprintf matches fmt.Sprintf.
func parity(t *testing.T, format string, args ...any) {
	t.Helper()
	want := fmt.Sprintf(format, args...)
	if got := Sprintf(format, args...); got != want {
		t.Errorf("Sprintf(%q, %#v)\n  got  %q\n  want %q", format, args, got, want)
	}
}

func TestSprintfParity(t *testing.T) {
	ptr := &struct{ x int }{}
	cases := []struct {
		format string
		args   []any
	}{
		// %s
		{"%s", []any{"hello"}},
		{"%s", []any{[]byte("bytes")}},
		{"[%s]", []any{""}},
		{"%s", []any{stringer{"x"}}},     // Stringer -> fallback
		{"%s", []any{myError{"boom"}}},   // error -> fallback
		{"%s", []any{42}},                // mismatch -> %!s(int=42)
		{"a %s b %s c", []any{"X", "Y"}}, // multiple

		// %d and integer widths
		{"%d", []any{0}},
		{"%d", []any{-12345}},
		{"%d", []any{int8(-128)}},
		{"%d", []any{int16(32767)}},
		{"%d", []any{int32(-7)}},
		{"%d", []any{int64(math.MaxInt64)}},
		{"%d", []any{uint(42)}},
		{"%d", []any{uint8(255)}},
		{"%d", []any{uint16(65535)}},
		{"%d", []any{uint32(4294967295)}},
		{"%d", []any{uint64(math.MaxUint64)}},
		{"%d", []any{uintptr(4096)}},
		{"%d", []any{"nope"}}, // mismatch

		// %x / %X
		{"%x", []any{255}},
		{"%X", []any{255}},
		{"%x", []any{-255}},
		{"%X", []any{-255}},
		{"%x", []any{uint64(0xdeadbeef)}},
		{"%X", []any{uint64(0xdeadbeef)}},
		{"%x", []any{"Go rocks"}},
		{"%X", []any{"Go rocks"}},
		{"%x", []any{[]byte{0x00, 0x0f, 0xa0, 0xff}}},
		{"%X", []any{[]byte{0x00, 0x0f, 0xa0, 0xff}}},
		{"%x", []any{3.14}}, // float hex -> fallback

		// %o / %b
		{"%o", []any{64}},
		{"%o", []any{-64}},
		{"%o", []any{uint(8)}},
		{"%b", []any{5}},
		{"%b", []any{-5}},
		{"%b", []any{uint8(0xff)}},
		{"%b", []any{3.14}}, // float -> fallback

		// %U
		{"%U", []any{'A'}},
		{"%U", []any{0}},
		{"%U", []any{0x1F600}},
		{"%U", []any{rune(0x10FFFF)}},
		{"%U", []any{uint16(0x20AC)}},
		{"%U", []any{-1}},     // negative -> fallback
		{"%U", []any{"nope"}}, // mismatch

		// %f
		{"%f", []any{3.14159}},
		{"%f", []any{-2.5}},
		{"%f", []any{0.0}},
		{"%f", []any{float32(0.1)}},
		{"%f", []any{math.NaN()}},
		{"%f", []any{math.Inf(1)}},
		{"%f", []any{math.Inf(-1)}},
		{"%f", []any{1e20}},
		{"%f", []any{42}}, // mismatch

		// %p (delegates)
		{"%p", []any{ptr}},
		{"ptr=%p done", []any{ptr}},

		// %% literal
		{"100%%", nil},
		{"%d%%", []any{50}},
		{"%%%%", nil},

		// special formats -> per-directive fallback
		{"%.2f", []any{3.14159}},
		{"%8.3f", []any{3.14159}},
		{"%-10d|", []any{42}},
		{"%+d", []any{42}},
		{"%05d", []any{42}},
		{"%#x", []any{255}},
		{"%#o", []any{64}},
		{"%6.2f%%", []any{99.5}},
		{"%q", []any{"quoted"}},
		{"%v", []any{[]int{1, 2, 3}}},
		{"%T", []any{3.14}},
		{"%c", []any{'A'}},
		{"%e", []any{1234.5678}},
		{"%g", []any{1234.5678}},
		{"%t", []any{true}},

		// star width / precision
		{"%*d", []any{5, 42}},
		{"%.*f", []any{2, 3.14159}},
		{"%*.*f", []any{8, 2, 3.14159}},

		// explicit index -> whole fallback
		{"%[1]d %[1]d", []any{7}},
		{"%[2]d %[1]d", []any{1, 2}},

		// arg-count mismatches -> whole fallback
		{"%d %d", []any{1}},          // missing
		{"%d", []any{1, 2}},          // extra
		{"no verbs", []any{1, 2, 3}}, // all extra
		{"plain text", nil},
		{"", nil},
		{"trailing %", nil}, // NOVERB

		// dense mixed line
		{"user=%s id=%d ok=%t hex=%x f=%.1f", []any{"bob", 7, true, 4096, 9.9}},
	}

	for _, c := range cases {
		parity(t, c.format, c.args...)
	}
}

func TestSprintfGenericParity(t *testing.T) {
	// Uniform string args.
	checkGEq(t, "%s/%s/%s", []any{"a", "b", "c"},
		SprintfGeneric[string]("%s/%s/%s", "a", "b", "c"))
	checkGEq(t, "%s and %x", []any{"x", "y"},
		SprintfGeneric[string]("%s and %x", "x", "y"))

	// Uniform ints across verbs.
	checkGEq(t, "%d %x %X %o %b %U", []any{255, 255, 255, 255, 255, 255},
		SprintfGeneric[int]("%d %x %X %o %b %U", 255, 255, 255, 255, 255, 255))
	checkGEq(t, "%d", []any{uint8(200)}, SprintfGeneric[uint8]("%d", uint8(200)))
	checkGEq(t, "%x", []any{int64(-1)}, SprintfGeneric[int64]("%x", int64(-1)))

	// Uniform floats.
	checkGEq(t, "%f|%f", []any{1.5, -2.25}, SprintfGeneric[float64]("%f|%f", 1.5, -2.25))
	checkGEq(t, "%f", []any{float32(0.1)}, SprintfGeneric[float32]("%f", float32(0.1)))

	// Uniform []byte.
	checkGEq(t, "%s-%x", []any{[]byte("hi"), []byte("hi")},
		SprintfGeneric[[]byte]("%s-%x", []byte("hi"), []byte("hi")))

	// Type with a verb it can't accelerate -> per-directive fmt.
	checkGEq(t, "%.2f", []any{3.14159}, SprintfGeneric[float64]("%.2f", 3.14159))

	// Unspecialised T falls back to the general engine.
	checkGEq(t, "%v", []any{true}, SprintfGeneric[bool]("%v", true))

	// Mixed args via SprintfGeneric[any].
	checkGEq(t, "%s id=%d", []any{"k", 7}, SprintfGeneric[any]("%s id=%d", "k", 7))
}

func checkGEq(t *testing.T, format string, args []any, got string) {
	t.Helper()
	if want := fmt.Sprintf(format, args...); got != want {
		t.Errorf("SprintfGeneric(%q, %#v)\n  got  %q\n  want %q", format, args, got, want)
	}
}

// TestRandomParity throws random verb/type combinations at both engines and
// requires them to agree with fmt.
func TestRandomParity(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	verbs := []byte{'s', 'd', 'x', 'X', 'U', 'f', 'o', 'b', 'v', 'q', 'p'}
	mkArg := func() any {
		switch rng.Intn(11) {
		case 0:
			return rng.Intn(1<<31) - (1 << 30)
		case 1:
			return int64(rng.Uint64())
		case 2:
			return uint(rng.Intn(1 << 20))
		case 3:
			return uint8(rng.Intn(256))
		case 4:
			return rng.NormFloat64() * 1000
		case 5:
			return float32(rng.NormFloat64())
		case 6:
			return randString(rng)
		case 7:
			return []byte(randString(rng))
		case 8:
			return rng.Intn(0x110000)
		case 9:
			return stringer{randString(rng)}
		default:
			return rng.Intn(2) == 0
		}
	}

	for iter := 0; iter < 5000; iter++ {
		n := rng.Intn(4)
		var sb strings.Builder
		args := make([]any, n)
		for i := 0; i < n; i++ {
			sb.WriteString("p")
			sb.WriteByte('%')
			sb.WriteByte(verbs[rng.Intn(len(verbs))])
			args[i] = mkArg()
		}
		format := sb.String()

		want := fmt.Sprintf(format, args...)
		if got := Sprintf(format, args...); got != want {
			t.Fatalf("Sprintf mismatch\n  format %q\n  args   %#v\n  got    %q\n  want   %q",
				format, args, got, want)
		}
		if got := SprintfGeneric[any](format, args...); got != want {
			t.Fatalf("SprintfGeneric mismatch\n  format %q\n  args   %#v\n  got    %q\n  want   %q",
				format, args, got, want)
		}
	}
}

func randString(rng *rand.Rand) string {
	n := rng.Intn(6)
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('!' + rng.Intn(90))
	}
	return string(b)
}
