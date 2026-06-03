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
