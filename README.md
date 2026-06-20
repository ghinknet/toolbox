# toolbox

[![Go Version](https://img.shields.io/badge/go-1.24%2B-00ADD8.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

Some easy tools make go easier.

`toolbox` is a small, dependency-free collection of generic utilities for everyday Go
development — cryptography helpers, secure random code/password generation, pointer
helpers, map/slice manipulation, and a few convenience types. Everything is built on
the standard library only.

## Features

- **Zero external dependencies** — pure Go standard library.
- **Generics-first** — leverages Go 1.18+ generics for type-safe helpers.
- **Crypto helpers** — RSA-OAEP encrypt/decrypt and executable fingerprinting.
- **Secure randomness** — configurable random codes and strong passwords built on `crypto/rand`.
- **Fast formatting** — `xfmt.Sprintf`, a (relatively) faster drop-in for `fmt.Sprintf` on the common verbs, with byte-identical output.
- **Well tested** — comprehensive unit tests and benchmarks across packages.

## Requirements

- Go **1.24.0** or newer.

## Installation

```bash
go get go.gh.ink/toolbox
```

## Packages

| Package | Import path | Description |
| --- | --- | --- |
| `crypto/fingerprint` | `go.gh.ink/toolbox/crypto/fingerprint` | Hash the running executable (MD5/SHA1/SHA256/SHA512). |
| `crypto/rsa` | `go.gh.ink/toolbox/crypto/rsa` | RSA-OAEP encryption/decryption and PEM key loading. |
| `data` | `go.gh.ink/toolbox/data` | Slice, map and string manipulation helpers. |
| `expr` | `go.gh.ink/toolbox/expr` | Expression helpers such as a generic ternary. |
| `pointer` | `go.gh.ink/toolbox/pointer` | Type-safe pointer creation, dereference and copy. |
| `random/code` | `go.gh.ink/toolbox/random/code` | Cryptographically secure random code / password generation. |
| `xfmt` | `go.gh.ink/toolbox/xfmt` | Faster drop-in `Sprintf` for the common format verbs, falling back to `fmt`. |
| `xtype` | `go.gh.ink/toolbox/xtype` | Convenience type aliases. |

## Usage

### crypto/rsa

RSA-OAEP encryption and decryption. Public keys must be PKIX-encoded, private keys
PKCS#8-encoded. `ReadKey` reads and decodes a PEM file into the raw DER block.

```go
import (
    "crypto/sha256"

    rsacrypt "go.gh.ink/toolbox/crypto/rsa"
)

pub, _ := rsacrypt.ReadKey("public.pem")
priv, _ := rsacrypt.ReadKey("private.pem")

// Encrypt to base64
cipher, err := rsacrypt.EncryptOAEPBase64([]byte("hello"), pub, sha256.New(), nil)
if err != nil {
    // handle error
}

// Decrypt from base64
plain, err := rsacrypt.DecryptOAEPBase64(cipher, priv, sha256.New(), nil)
```

API:

- `EncryptOAEP(plainText, block []byte, hash hash.Hash, label []byte) ([]byte, error)`
- `EncryptOAEPBase64(plainText, block []byte, hash hash.Hash, label []byte) (string, error)`
- `DecryptOAEP(cipherText, block []byte, hash hash.Hash, label []byte) ([]byte, error)`
- `DecryptOAEPBase64(cipherText string, block []byte, hash hash.Hash, label []byte) ([]byte, error)`
- `ReadKey(path string) ([]byte, error)`

### crypto/fingerprint

Compute a hash of the currently running executable, useful for integrity checks.

```go
import "go.gh.ink/toolbox/crypto/fingerprint"

sum, err := fingerprint.GetExecutableSHA256()
```

API:

- `GetExecutableMD5() (string, error)`
- `GetExecutableSHA1() (string, error)`
- `GetExecutableSHA256() (string, error)`
- `GetExecutableSHA512() (string, error)`

### random/code

Cryptographically secure random code and password generation, backed by `crypto/rand`.

```go
import "go.gh.ink/toolbox/random/code"

// 6-digit numeric code
otp := code.Number(6)

// 16-char strong password (all categories, ambiguous chars removed)
pw := code.StrongPassword()

// Fully configurable
s, err := code.Code(code.Options{
    Digit:        12,
    UseNumbers:   true,
    UseLowercase: true,
    UseUppercase: true,
    UseSymbols:   true,
    ExcludeChars: "0OIl1", // avoid ambiguous characters
    RequireEach:  true,    // guarantee at least one of each enabled category
})
```

`Options`:

| Field | Type | Description |
| --- | --- | --- |
| `Digit` | `int` | Length of the generated code (defaults to 6 if `<= 0`). |
| `UseNumbers` | `bool` | Include digits `0-9`. |
| `UseLowercase` | `bool` | Include lowercase `a-z`. |
| `UseUppercase` | `bool` | Include uppercase `A-Z`. |
| `UseSymbols` | `bool` | Include symbols (`DefaultSymbols` = `!@#$%^&*?`). |
| `CustomSymbols` | `string` | If set, replaces the default symbol set entirely. |
| `ExcludeChars` | `string` | Characters to exclude (e.g. `"0OIl1"`). |
| `RequireEach` | `bool` | Guarantee at least one character from every enabled category. |

API:

- `Code(config Options) (string, error)`
- `MustCode(config Options) string` — panics on error.
- `Number(digit int) string`
- `Alpha(digit int, useLower, useUpper bool) string`
- `Mixed(digit int, useNumbers, useLower, useUpper, useSymbols bool) string`
- `StrongPassword() string`

### pointer

Type-safe helpers for working with pointers and generics.

```go
import "go.gh.ink/toolbox/pointer"

p := pointer.Ref(42)          // *int
v := pointer.SafeDeref(p)     // 42; returns zero value when p is nil
cp := pointer.Copy(p)         // new *int pointing to a copy
ptrs := pointer.SliceRef([]int{1, 2, 3}) // []*int
```

API:

- `Ref[T any](v T) *T`
- `SliceRef[T any](v []T) []*T` — pointers share the slice's backing array.
- `SliceCopyRef[T any](v []T) []*T` — each element is copied first.
- `SafeDeref[T any](ptr *T) T` — returns the zero value if `ptr` is nil.
- `Copy[T any](p *T) *T`

### data

Generic slice, map and string utilities.

```go
import "go.gh.ink/toolbox/data"

n, _ := data.Atoi("123")

s := data.MakeSliceNotNil[int](nil) // []int{} instead of nil

m := map[string]int{"a": 1, "b": 2}
keys := data.MapKeys(m)
vals := data.MapValues(m)

merged := data.MergeMapsInt(
    map[string]int{"a": 1},
    map[string]int{"a": 2, "b": 3},
) // {"a": 3, "b": 3}
```

API:

- String: `Atoi(s string) (int, error)`, `Atoi32(s string) (int32, error)`, `Atoi64(s string) (int64, error)`
- Slice: `MakeSliceNotNil[T any, S ~[]T](slice S) S`
- Map: `MapKeys`, `MapValues`, `MapKeysValues`
- Merge: `MergeMapsString` (concatenates), `MergeMapsStringDropMismatch` (keeps shared keys only), `MergeMapsInt` (adds)
- Trim: `TrimPrefixMapsString`, `TrimSuffixMapsString`

### expr

```go
import "go.gh.ink/toolbox/expr"

max := expr.Ternary(a > b, a, b)
```

API:

- `Ternary[T any](condition bool, trueVal, falseVal T) T`

### xfmt

A faster, drop-in `Sprintf`. It uses type assertions instead of reflection for
the most common verbs — `%s %d %x %X %U %f %p %% %o %b` — and defers everything
else (decorated formats such as `%.2f` or `%-10d`, and any other verb) to the
standard library on a per-directive basis. Output is **byte-for-byte identical to
`fmt.Sprintf`**.

```go
import "go.gh.ink/toolbox/xfmt"

s := xfmt.Sprintf("user=%s id=%d hex=%x", "bob", 7, 255)
// user=bob id=7 hex=ff

// When every argument shares one concrete type, SprintfGeneric receives them as
// ...T (not ...any), so the compiler catches type mismatches at the call site.
s = xfmt.SprintfGeneric[int]("%d-%d-%d", 1, 22, 333)

// Heterogeneous args work too — use SprintfGeneric[any].
s = xfmt.SprintfGeneric[any]("%s id=%d", "k", 7)
```

API:

- `Sprintf(format string, a ...any) string`
- `SprintfGeneric[T any](format string, a ...T) string`

How it works:

- **Fast path** — the *plain* form of `%s %d %x %X %U %f %o %b` (and `%%`) is
  formatted directly with type assertions and `strconv`, into a pooled, size-hinted
  buffer.
- **Per-directive fallback** — any directive carrying a flag, width, precision or `*`
  (e.g. `%.2f`, `%-10d`, `%08x`), and any unhandled verb (`%v`, `%q`, `%p`, …), is
  handed to `fmt` for just that one directive.
- **Whole fallback** — explicit argument indices (`%[1]d`) and argument-count
  mismatches defer to `fmt.Sprintf` entirely, preserving its exact diagnostics
  (`%!d(MISSING)`, `%!(EXTRA ...)`).

#### Performance

Benchmarked against `fmt.Sprintf` on identical format + arguments, Go 1.24,
`go test -bench=Compare -benchmem`. Speedup = `fmt` ÷ implementation (>1 means
faster).

##### Apple M1 Pro (ARM64)

| Scenario | Example | `fmt` | `Sprintf` | Speedup | `SprintfGeneric` | Speedup | Allocs fmt/Sprintf/Generic |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: |
| int | `%d` | 49.2 | 38.0 | 1.30× | 37.2 | 1.32× | 1/1/1 |
| ints ×4 | `%d-%d-%d-%d` | 106.6 | 86.4 | 1.23× | 90.5 | 1.18× | 1/1/1 |
| uint hex | `%x` | 46.7 | 34.2 | 1.37× | 35.4 | 1.32× | 1/1/1 |
| string | `%s` | 35.9 | 29.4 | 1.22× | 30.7 | 1.17× | 1/1/1 |
| strings ×3 | `%s/%s/%s` | 68.0 | 56.3 | 1.21× | 53.8 | 1.26× | 1/1/1 |
| float | `%f` | 69.8 | 59.1 | 1.18× | 59.9 | 1.17× | 1/1/1 |
| bytes hex | `%x` (`[]byte`) | 81.2 | 63.8 | 1.27× | 49.1 | 1.65× | 3/3/2 |
| mixed (hot verbs) | `%s id=%d hex=%x f=%f` | 130.8 | 108.8 | 1.20× | — | — | 1/1/— |
| log line | `[%s] svc=%s code=%d` | 82.1 | 67.3 | 1.22× | — | — | 1/1/— |
| with fallback | `%s took %.2fms` | 94.1 | 106.5 | 0.88× | — | — | 1/1/— |
| all fallback | `%v %q %t` | 209.1 | 271.7 | 0.77× | — | — | 5/5/— |

##### Intel i7-13700K (x86‑64)

| Scenario | Example | `fmt` | `Sprintf` | Speedup | `SprintfGeneric` | Speedup | Allocs fmt/Sprintf/Generic |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: |
| int | `%d` | 40.2 | 36.1 | 1.11× | 38.6 | 1.04× | 1/1/1 |
| ints ×4 | `%d-%d-%d-%d` | 99.2 | 76.2 | 1.30× | 75.6 | 1.31× | 1/1/1 |
| uint hex | `%x` | 40.1 | 37.3 | 1.08× | 34.3 | 1.17× | 1/1/1 |
| string | `%s` | 46.8 | 43.1 | 1.08× | 47.3 | 0.99× | 1/1/1 |
| strings ×3 | `%s/%s/%s` | 68.1 | 47.5 | 1.43× | 56.5 | 1.21× | 1/1/1 |
| float | `%f` | 55.4 | 43.8 | 1.26× | 43.7 | 1.27× | 1/1/1 |
| bytes hex | `%x` (`[]byte`) | 108.8 | 113.9 | 0.96× | 63.2 | 1.72× | 3/3/2 |
| mixed (hot verbs) | `%s id=%d hex=%x f=%f` | 135.1 | 121.7 | 1.11× | — | — | 1/1/— |
| log line | `[%s] svc=%s code=%d` | 88.1 | 70.8 | 1.24× | — | — | 1/1/— |
| with fallback | `%s took %.2fms` | 103.8 | 112.2 | 0.92× | — | — | 1/1/— |
| all fallback | `%v %q %t` | 247.1 | 280.9 | 0.88× | — | — | 5/5/— |

#### Overhead vs bare `strconv`

Single‑value conversion benchmarks comparing `xfmt` against a raw `strconv` call
(the theoretical floor) and `fmt.Sprintf`. The "tax" is `xfmt`/`fmt` minus `raw`.

```
go test ./xfmt/ -bench=Overhead -benchmem
```

###### Apple M1 Pro (ARM64)

| Scenario | `raw` | `fmt` | `xfmt` | `fmt` tax | `xfmt` tax |
| --- | ---: | ---: | ---: | ---: | ---: |
| int → `%d` | 18.7 | 49.4 | 37.9 | +30.7 | +19.2 |
| uint → `%x` | 17.5 | 47.8 | 35.5 | +30.3 | +18.0 |
| int → `%o` | 17.0 | 41.6 | 36.9 | +24.7 | +19.9 |
| int → `%b` | 17.1 | 41.2 | 36.3 | +24.1 | +19.2 |
| float → `%f` | 43.4 | 69.3 | 59.2 | +25.9 | +15.8 |
| `%s` (no‑op on string) | 2.1 | 35.8 | 29.3 | +33.8 | +27.3 |
| `[]byte` → `%x` | 19.2 | 73.8 | 64.8 | +54.6 | +45.6 |
| `[]byte` → `%x` (`Generic`) | 19.2 | 73.8 | 49.1 | +54.6 | +29.9 |
| pre+int (`n=%d`) | 40.9 | 53.2 | 44.7 | +12.3 | +3.8 |

###### Intel i7-13700K (x86‑64)

| Scenario | `raw` | `fmt` | `xfmt` | `fmt` tax | `xfmt` tax |
| --- | ---: | ---: | ---: | ---: | ---: |
| int → `%d` | 23.4 | 43.4 | 35.5 | +20.0 | +12.1 |
| uint → `%x` | 25.5 | 44.8 | 34.9 | +19.2 | +9.4 |
| int → `%o` | 13.0 | 35.2 | 30.7 | +22.2 | +17.7 |
| int → `%b` | 14.0 | 35.5 | 30.3 | +21.5 | +16.3 |
| float → `%f` | 35.3 | 55.3 | 45.4 | +20.0 | +10.1 |
| `%s` (no‑op on string) | 1.0 | 48.7 | 41.4 | +47.7 | +40.3 |
| `[]byte` → `%x` | 32.7 | 126.5 | 102.8 | +93.8 | +70.0 |
| `[]byte` → `%x` (`Generic`) | 32.7 | 126.5 | 75.9 | +93.8 | +43.2 |
| pre+int (`n=%d`) | 52.9 | 61.0 | 57.1 | +8.1 | +4.2 |

Reproduce locally:

```bash
# head-to-head against fmt
go test ./xfmt -bench=Compare -benchmem

# overhead breakdown vs raw strconv
go test ./xfmt -bench=Overhead -benchmem

# at-a-glance table with speedups
XFMT_REPORT=1 go test ./xfmt -run TestComparisonReport -v
```

<details>
<summary>Raw <code>go test -bench=Compare -benchmem</code> output (both platforms)</summary>

**Apple M1 Pro (ARM64)**
```
BenchmarkCompare/int/fmt-8                       24138008        49.22 ns/op        8 B/op       1 allocs/op
BenchmarkCompare/int/xfmt-8                      32191862        38.01 ns/op        8 B/op       1 allocs/op
BenchmarkCompare/int/xfmt_generic-8              32659503        37.16 ns/op        8 B/op       1 allocs/op
BenchmarkCompare/ints_x4/fmt-8                   11154662       106.6 ns/op        16 B/op       1 allocs/op
BenchmarkCompare/ints_x4/xfmt-8                  13803064        86.43 ns/op       16 B/op       1 allocs/op
BenchmarkCompare/ints_x4/xfmt_generic-8          14480263        90.48 ns/op       16 B/op       1 allocs/op
BenchmarkCompare/uint_hex/fmt-8                  25228148        46.74 ns/op        8 B/op       1 allocs/op
BenchmarkCompare/uint_hex/xfmt-8                 35498538        34.22 ns/op        8 B/op       1 allocs/op
BenchmarkCompare/uint_hex/xfmt_generic-8         33413152        35.37 ns/op        8 B/op       1 allocs/op
BenchmarkCompare/string/fmt-8                    34052816        35.86 ns/op       16 B/op       1 allocs/op
BenchmarkCompare/string/xfmt-8                   40588132        29.39 ns/op       16 B/op       1 allocs/op
BenchmarkCompare/string/xfmt_generic-8           40808576        30.74 ns/op       16 B/op       1 allocs/op
BenchmarkCompare/strings_x3/fmt-8                17595027        68.00 ns/op       16 B/op       1 allocs/op
BenchmarkCompare/strings_x3/xfmt-8               22713846        56.33 ns/op       16 B/op       1 allocs/op
BenchmarkCompare/strings_x3/xfmt_generic-8       22474346        53.80 ns/op       16 B/op       1 allocs/op
BenchmarkCompare/float/fmt-8                     17246674        69.80 ns/op        8 B/op       1 allocs/op
BenchmarkCompare/float/xfmt-8                    20544720        59.10 ns/op        8 B/op       1 allocs/op
BenchmarkCompare/float/xfmt_generic-8            20012799        59.91 ns/op        8 B/op       1 allocs/op
BenchmarkCompare/bytes_hex/fmt-8                 15965230        81.23 ns/op       48 B/op       3 allocs/op
BenchmarkCompare/bytes_hex/xfmt-8                19082476        63.83 ns/op       48 B/op       3 allocs/op
BenchmarkCompare/bytes_hex/xfmt_generic-8        24523389        49.11 ns/op       24 B/op       2 allocs/op
BenchmarkCompare/mixed_hot/fmt-8                  9214952       130.8 ns/op        24 B/op       1 allocs/op
BenchmarkCompare/mixed_hot/xfmt-8                10965279       108.8 ns/op        24 B/op       1 allocs/op
BenchmarkCompare/log_line/fmt-8                  14383666        82.10 ns/op       24 B/op       1 allocs/op
BenchmarkCompare/log_line/xfmt-8                 17820006        67.32 ns/op       24 B/op       1 allocs/op
BenchmarkCompare/with_fallback/fmt-8             12770179        94.10 ns/op       24 B/op       1 allocs/op
BenchmarkCompare/with_fallback/xfmt-8            11301354       106.5 ns/op        24 B/op       1 allocs/op
BenchmarkCompare/all_fallback/fmt-8               5722921       209.1 ns/op        72 B/op       5 allocs/op
BenchmarkCompare/all_fallback/xfmt-8              4493806       271.7 ns/op        72 B/op       5 allocs/op
```

**Intel i7-13700K (x86‑64)**
```
BenchmarkCompare/int/fmt-24                     25824943        40.17 ns/op        8 B/op       1 allocs/op
BenchmarkCompare/int/xfmt-24                    34964701        36.12 ns/op        8 B/op       1 allocs/op
BenchmarkCompare/int/xfmt_generic-24            32031006        38.64 ns/op        8 B/op       1 allocs/op
BenchmarkCompare/ints_x4/fmt-24                 13032422        99.23 ns/op       16 B/op       1 allocs/op
BenchmarkCompare/ints_x4/xfmt-24                15398772        76.22 ns/op       16 B/op       1 allocs/op
BenchmarkCompare/ints_x4/xfmt_generic-24        14505160        75.62 ns/op       16 B/op       1 allocs/op
BenchmarkCompare/uint_hex/fmt-24                25859976        40.11 ns/op        8 B/op       1 allocs/op
BenchmarkCompare/uint_hex/xfmt-24               35570955        37.25 ns/op        8 B/op       1 allocs/op
BenchmarkCompare/uint_hex/xfmt_generic-24       35969205        34.25 ns/op        8 B/op       1 allocs/op
BenchmarkCompare/string/fmt-24                  38183049        46.78 ns/op       16 B/op       1 allocs/op
BenchmarkCompare/string/xfmt-24                 45570930        43.12 ns/op       16 B/op       1 allocs/op
BenchmarkCompare/string/xfmt_generic-24         35418469        47.25 ns/op       16 B/op       1 allocs/op
BenchmarkCompare/strings_x3/fmt-24              18044306        68.11 ns/op       16 B/op       1 allocs/op
BenchmarkCompare/strings_x3/xfmt-24             24852638        47.46 ns/op       16 B/op       1 allocs/op
BenchmarkCompare/strings_x3/xfmt_generic-24     28393248        56.48 ns/op       16 B/op       1 allocs/op
BenchmarkCompare/float/fmt-24                   21492949        55.42 ns/op        8 B/op       1 allocs/op
BenchmarkCompare/float/xfmt-24                  27909249        43.82 ns/op        8 B/op       1 allocs/op
BenchmarkCompare/float/xfmt_generic-24          27788601        43.66 ns/op        8 B/op       1 allocs/op
BenchmarkCompare/bytes_hex/fmt-24               10662582       108.8 ns/op        48 B/op       3 allocs/op
BenchmarkCompare/bytes_hex/xfmt-24              10773590       113.9 ns/op        48 B/op       3 allocs/op
BenchmarkCompare/bytes_hex/xfmt_generic-24      19580682        63.19 ns/op       24 B/op       2 allocs/op
BenchmarkCompare/mixed_hot/fmt-24                9462279       135.1 ns/op        24 B/op       1 allocs/op
BenchmarkCompare/mixed_hot/xfmt-24               9786537       121.7 ns/op        24 B/op       1 allocs/op
BenchmarkCompare/log_line/fmt-24                13548580        88.11 ns/op       24 B/op       1 allocs/op
BenchmarkCompare/log_line/xfmt-24               17041275        70.81 ns/op       24 B/op       1 allocs/op
BenchmarkCompare/with_fallback/fmt-24           12643534       103.8 ns/op        24 B/op       1 allocs/op
BenchmarkCompare/with_fallback/xfmt-24          11922300       112.2 ns/op        24 B/op       1 allocs/op
BenchmarkCompare/all_fallback/fmt-24             4864030       247.1 ns/op        72 B/op       5 allocs/op
BenchmarkCompare/all_fallback/xfmt-24            3856126       280.9 ns/op        72 B/op       5 allocs/op
```

</details>

### xtype

Convenience type aliases for common generic maps.

```go
import "go.gh.ink/toolbox/xtype"

var payload xtype.H        // map[string]any
var counts  xtype.MS[int]  // map[string]int
```

API:

- `type H = map[string]any`
- `type MS[V any] = map[string]V`

## Testing

```bash
go test ./...
```

Benchmarks (e.g. in `random/code`) can be run with:

```bash
go test ./random/code -bench .
```

## License

Licensed under the [Apache License 2.0](LICENSE).
