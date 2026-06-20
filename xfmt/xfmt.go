package xfmt

import (
	"fmt"
	"strconv"
	"sync"
)

// appenderFunc formats a single argument for a plain verb, appending the result
// to buf. ok reports whether the (verb, argument type) pair was handled on the
// fast path; when false the caller falls back to fmt for that directive.
type appenderFunc = func(buf []byte, verb byte, arg any) (out []byte, ok bool)

// typedAppender is the generic-path analogue: it receives T directly — zero
// type assertions in the hot loop.
type typedAppender[T any] func(buf []byte, verb byte, arg T) ([]byte, bool)

const (
	hexLower     = "0123456789abcdef"
	hexUpper     = "0123456789ABCDEF"
	maxPooledCap = 1 << 16
	stackBufSize = 8 // ≤8 args use a stack buffer, zero heap
)

var bufPool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, 64)
		return &b
	},
}

// Sprintf formats according to format using fast type-assertion paths for the
// common verbs and falling back to fmt.Sprintf for everything else. Its output
// is identical to fmt.Sprintf.
func Sprintf(format string, a ...any) string {
	return run(format, a, appendVerb)
}

// SprintfGeneric behaves like Sprintf but, when every argument shares concrete
// type T, the arguments arrive as ...T (not ...any) — zero boxing at the call
// site, zero type assertions in the fast path. For mixed-type format strings,
// use SprintfGeneric[any]; it delegates to the same engine as Sprintf.
func SprintfGeneric[T any](format string, a ...T) string {
	// Fast-path: when T is any, the []T slice is already []any.
	if _, ok := any(a).([]any); ok {
		return run(format, any(a).([]any), appendVerb)
	}

	// Narrow T to a concrete type and dispatch to a fully typed inner loop.
	switch any(a).(type) {
	case []int:
		return runGeneric[int](format, any(a).([]int), tAppSigned[int])
	case []int8:
		return runGeneric[int8](format, any(a).([]int8), tAppSigned[int8])
	case []int16:
		return runGeneric[int16](format, any(a).([]int16), tAppSigned[int16])
	case []int32:
		return runGeneric[int32](format, any(a).([]int32), tAppSigned[int32])
	case []int64:
		return runGeneric[int64](format, any(a).([]int64), tAppSigned[int64])
	case []uint:
		return runGeneric[uint](format, any(a).([]uint), tAppUnsigned[uint])
	case []uint8:
		return runGeneric[uint8](format, any(a).([]uint8), tAppUnsigned[uint8])
	case []uint16:
		return runGeneric[uint16](format, any(a).([]uint16), tAppUnsigned[uint16])
	case []uint32:
		return runGeneric[uint32](format, any(a).([]uint32), tAppUnsigned[uint32])
	case []uint64:
		return runGeneric[uint64](format, any(a).([]uint64), tAppUnsigned[uint64])
	case []uintptr:
		return runGeneric[uintptr](format, any(a).([]uintptr), tAppUnsigned[uintptr])
	case []string:
		return runGeneric[string](format, any(a).([]string), tAppString)
	case [][]byte:
		return runGeneric[[]byte](format, any(a).([][]byte), tAppBytes)
	case []float64:
		return runGeneric[float64](format, any(a).([]float64), tAppFloat64)
	case []float32:
		return runGeneric[float32](format, any(a).([]float32), tAppFloat32)
	}
	// Unrecognised concrete T: box each arg once into any and use the general engine.
	tmp := boxToAny(a)
	return run(format, tmp, appendVerb)
}

// boxToAny converts []T to []any using a stack buffer when n ≤ stackBufSize.
func boxToAny[T any](a []T) []any {
	n := len(a)
	var small [stackBufSize]any
	var out []any
	if n <= stackBufSize {
		out = small[:n]
	} else {
		out = make([]any, n)
	}
	for i := range a {
		out[i] = a[i]
	}
	return out
}

// anyFromStack copies up to stackBufSize elements from a []T into a stack-based
// []any, returning the slice and whether the whole slice fit.
func anyFromStack[T any](a []T, start, n int) ([]any, bool) {
	if start+n > len(a) {
		return nil, false
	}
	var small [stackBufSize]any
	if n > stackBufSize {
		return nil, false
	}
	out := small[:n]
	for i := range n {
		out[i] = a[start+i]
	}
	return out, true
}

// run drives the any-based formatting engine with a borrowed, recycled buffer.
func run(format string, a []any, app appenderFunc) string {
	bp := bufPool.Get().(*[]byte)
	buf, ok := appendf((*bp)[:0], format, a, app)

	var s string
	if ok {
		s = string(buf)
	} else {
		s = fmt.Sprintf(format, a...)
	}

	if cap(buf) <= maxPooledCap {
		*bp = buf[:0]
		bufPool.Put(bp)
	}
	return s
}

// runGeneric is the typed analogue of run.
func runGeneric[T any](format string, a []T, app typedAppender[T]) string {
	bp := bufPool.Get().(*[]byte)
	buf, ok := appendfGeneric[T]((*bp)[:0], format, a, app)

	var s string
	if ok {
		s = string(buf)
	} else {
		s = fmt.Sprintf(format, boxToAny(a)...)
	}

	if cap(buf) <= maxPooledCap {
		*bp = buf[:0]
		bufPool.Put(bp)
	}
	return s
}

// appendf walks format once, writing into buf. It returns ok=false to signal
// that the whole call should be re-done by fmt.
func appendf(buf []byte, format string, a []any, app appenderFunc) ([]byte, bool) {
	if n := len(format) + 16*len(a); cap(buf) < n {
		buf = make([]byte, 0, n)
	}

	argIdx := 0
	end := len(format)
	for i := 0; i < end; {
		if format[i] != '%' {
			j := i + 1
			for j < end && format[j] != '%' {
				j++
			}
			buf = append(buf, format[i:j]...)
			i = j
			continue
		}

		if i+1 >= end {
			return buf, false
		}

		if next := format[i+1]; isFastVerb(next) {
			if next == '%' {
				buf = append(buf, '%')
				i += 2
				continue
			}
			if argIdx >= len(a) {
				return buf, false
			}
			arg := a[argIdx]
			var ok bool
			if buf, ok = app(buf, next, arg); !ok {
				buf = fmt.Appendf(buf, format[i:i+2], arg)
			}
			argIdx++
			i += 2
			continue
		}

		verbPos, stars, ok := scanDirective(format, i)
		if !ok {
			return buf, false
		}
		need := stars
		if format[verbPos] != '%' {
			need++
		}
		if argIdx+need > len(a) {
			return buf, false
		}
		buf = fmt.Appendf(buf, format[i:verbPos+1], a[argIdx:argIdx+need]...)
		argIdx += need
		i = verbPos + 1
	}

	if argIdx != len(a) {
		return buf, false
	}
	return buf, true
}

// appendfGeneric is the typed analogue of appendf.
func appendfGeneric[T any](buf []byte, format string, a []T, app typedAppender[T]) ([]byte, bool) {
	if n := len(format) + 16*len(a); cap(buf) < n {
		buf = make([]byte, 0, n)
	}

	argIdx := 0
	end := len(format)
	for i := 0; i < end; {
		if format[i] != '%' {
			j := i + 1
			for j < end && format[j] != '%' {
				j++
			}
			buf = append(buf, format[i:j]...)
			i = j
			continue
		}

		if i+1 >= end {
			return buf, false
		}

		if next := format[i+1]; isFastVerb(next) {
			if next == '%' {
				buf = append(buf, '%')
				i += 2
				continue
			}
			if argIdx >= len(a) {
				return buf, false
			}
			arg := a[argIdx]
			var ok bool
			if buf, ok = app(buf, next, arg); !ok {
				buf = fmt.Appendf(buf, format[i:i+2], any(arg))
			}
			argIdx++
			i += 2
			continue
		}

		verbPos, stars, ok := scanDirective(format, i)
		if !ok {
			return buf, false
		}
		need := stars
		if format[verbPos] != '%' {
			need++
		}
		if argIdx+need > len(a) {
			return buf, false
		}
		anyBuf, fits := anyFromStack(a, argIdx, need)
		if !fits {
			return buf, false
		}
		buf = fmt.Appendf(buf, format[i:verbPos+1], anyBuf...)
		argIdx += need
		i = verbPos + 1
	}

	if argIdx != len(a) {
		return buf, false
	}
	return buf, true
}

func isFastVerb(c byte) bool {
	switch c {
	case 's', 'd', 'x', 'X', 'U', 'f', 'p', 'o', 'b', '%':
		return true
	}
	return false
}

func scanDirective(s string, start int) (verbPos, stars int, ok bool) {
	n := len(s)
	i := start + 1

	for i < n {
		switch s[i] {
		case '+', '-', '#', ' ', '0':
			i++
			continue
		}
		break
	}
	if i < n && s[i] == '[' {
		return 0, 0, false
	}

	if i < n && s[i] == '*' {
		stars++
		i++
	} else {
		for i < n && s[i] >= '0' && s[i] <= '9' {
			i++
		}
	}

	if i < n && s[i] == '.' {
		i++
		if i < n && s[i] == '*' {
			stars++
			i++
		} else {
			for i < n && s[i] >= '0' && s[i] <= '9' {
				i++
			}
		}
	}

	if i < n && s[i] == '[' {
		return 0, 0, false
	}
	if i >= n || s[i] >= 0x80 {
		return 0, 0, false
	}
	return i, stars, true
}

func appendVerb(buf []byte, verb byte, arg any) ([]byte, bool) {
	switch verb {
	case 's':
		switch v := arg.(type) {
		case string:
			return append(buf, v...), true
		case []byte:
			return append(buf, v...), true
		}
	case 'd':
		return appendInt(buf, arg, 10, false)
	case 'b':
		return appendInt(buf, arg, 2, false)
	case 'o':
		return appendInt(buf, arg, 8, false)
	case 'x':
		switch v := arg.(type) {
		case string:
			return appendHexStr(buf, v, false), true
		case []byte:
			return appendHexBytes(buf, v, false), true
		}
		return appendInt(buf, arg, 16, false)
	case 'X':
		switch v := arg.(type) {
		case string:
			return appendHexStr(buf, v, true), true
		case []byte:
			return appendHexBytes(buf, v, true), true
		}
		return appendInt(buf, arg, 16, true)
	case 'U':
		if v, ok := codePoint(arg); ok {
			return appendU(buf, v), true
		}
	case 'f':
		switch v := arg.(type) {
		case float64:
			return strconv.AppendFloat(buf, v, 'f', 6, 64), true
		case float32:
			return strconv.AppendFloat(buf, float64(v), 'f', 6, 32), true
		}
	}
	return buf, false
}

func appendInt(buf []byte, arg any, base int, upper bool) ([]byte, bool) {
	start := len(buf)
	switch v := arg.(type) {
	case int:
		buf = strconv.AppendInt(buf, int64(v), base)
	case int8:
		buf = strconv.AppendInt(buf, int64(v), base)
	case int16:
		buf = strconv.AppendInt(buf, int64(v), base)
	case int32:
		buf = strconv.AppendInt(buf, int64(v), base)
	case int64:
		buf = strconv.AppendInt(buf, v, base)
	case uint:
		buf = strconv.AppendUint(buf, uint64(v), base)
	case uint8:
		buf = strconv.AppendUint(buf, uint64(v), base)
	case uint16:
		buf = strconv.AppendUint(buf, uint64(v), base)
	case uint32:
		buf = strconv.AppendUint(buf, uint64(v), base)
	case uint64:
		buf = strconv.AppendUint(buf, v, base)
	case uintptr:
		buf = strconv.AppendUint(buf, uint64(v), base)
	default:
		return buf, false
	}
	if upper {
		upperHexTail(buf, start)
	}
	return buf, true
}

func appendHexStr(buf []byte, s string, upper bool) []byte {
	digits := hexLower
	if upper {
		digits = hexUpper
	}
	for i := 0; i < len(s); i++ {
		b := s[i]
		buf = append(buf, digits[b>>4], digits[b&0xf])
	}
	return buf
}

func appendHexBytes(buf, src []byte, upper bool) []byte {
	digits := hexLower
	if upper {
		digits = hexUpper
	}
	for _, b := range src {
		buf = append(buf, digits[b>>4], digits[b&0xf])
	}
	return buf
}

func appendU(buf []byte, v uint64) []byte {
	buf = append(buf, 'U', '+')
	var tmp [16]byte
	p := len(tmp)
	for v >= 16 {
		p--
		tmp[p] = hexUpper[v&0xf]
		v >>= 4
	}
	p--
	tmp[p] = hexUpper[v&0xf]
	for len(tmp)-p < 4 {
		p--
		tmp[p] = '0'
	}
	return append(buf, tmp[p:]...)
}

func upperHexTail(b []byte, start int) {
	for i := start; i < len(b); i++ {
		if c := b[i]; c >= 'a' && c <= 'f' {
			b[i] = c - 'a' + 'A'
		}
	}
}

func codePoint(arg any) (uint64, bool) {
	switch v := arg.(type) {
	case int:
		return nonNeg(int64(v))
	case int8:
		return nonNeg(int64(v))
	case int16:
		return nonNeg(int64(v))
	case int32:
		return nonNeg(int64(v))
	case int64:
		return nonNeg(v)
	case uint:
		return uint64(v), true
	case uint8:
		return uint64(v), true
	case uint16:
		return uint64(v), true
	case uint32:
		return uint64(v), true
	case uint64:
		return v, true
	case uintptr:
		return uint64(v), true
	}
	return 0, false
}

func nonNeg(v int64) (uint64, bool) {
	if v < 0 {
		return 0, false
	}
	return uint64(v), true
}

// --- typed appenders (receive T directly, zero type assertions) ---

type signedInt interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type unsignedInt interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

func tAppSigned[I signedInt](buf []byte, verb byte, v I) ([]byte, bool) {
	n := int64(v)
	switch verb {
	case 'd':
		return strconv.AppendInt(buf, n, 10), true
	case 'b':
		return strconv.AppendInt(buf, n, 2), true
	case 'o':
		return strconv.AppendInt(buf, n, 8), true
	case 'x':
		return strconv.AppendInt(buf, n, 16), true
	case 'X':
		start := len(buf)
		buf = strconv.AppendInt(buf, n, 16)
		upperHexTail(buf, start)
		return buf, true
	case 'U':
		if v < 0 {
			return buf, false
		}
		return appendU(buf, uint64(n)), true
	}
	return buf, false
}

func tAppUnsigned[U unsignedInt](buf []byte, verb byte, v U) ([]byte, bool) {
	n := uint64(v)
	switch verb {
	case 'd':
		return strconv.AppendUint(buf, n, 10), true
	case 'b':
		return strconv.AppendUint(buf, n, 2), true
	case 'o':
		return strconv.AppendUint(buf, n, 8), true
	case 'x':
		return strconv.AppendUint(buf, n, 16), true
	case 'X':
		start := len(buf)
		buf = strconv.AppendUint(buf, n, 16)
		upperHexTail(buf, start)
		return buf, true
	case 'U':
		return appendU(buf, n), true
	}
	return buf, false
}

func tAppString(buf []byte, verb byte, v string) ([]byte, bool) {
	switch verb {
	case 's':
		return append(buf, v...), true
	case 'x':
		return appendHexStr(buf, v, false), true
	case 'X':
		return appendHexStr(buf, v, true), true
	}
	return buf, false
}

func tAppBytes(buf []byte, verb byte, v []byte) ([]byte, bool) {
	switch verb {
	case 's':
		return append(buf, v...), true
	case 'x':
		return appendHexBytes(buf, v, false), true
	case 'X':
		return appendHexBytes(buf, v, true), true
	}
	return buf, false
}

func tAppFloat64(buf []byte, verb byte, v float64) ([]byte, bool) {
	if verb != 'f' {
		return buf, false
	}
	return strconv.AppendFloat(buf, v, 'f', 6, 64), true
}

func tAppFloat32(buf []byte, verb byte, v float32) ([]byte, bool) {
	if verb != 'f' {
		return buf, false
	}
	return strconv.AppendFloat(buf, float64(v), 'f', 6, 32), true
}
