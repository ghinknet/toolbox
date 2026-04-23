package code

import (
	"crypto/rand"
	"errors"
	"math/big"
	"strings"
	"sync"
)

// Options holds the configuration for random code generation.
type Options struct {
	Digit         int    // Number of characters to generate
	UseNumbers    bool   // Include numeric characters
	UseLowercase  bool   // Include lowercase letters
	UseUppercase  bool   // Include uppercase letters
	UseSymbols    bool   // Include symbols
	CustomSymbols string // Custom symbol set; if not empty, replaces the default symbol set entirely
	ExcludeChars  string // Characters to exclude (e.g., "0OIl1" to avoid ambiguous characters)
	RequireEach   bool   // Guarantee at least one character from every enabled category
}

// DefaultSymbols is the default set of symbol characters.
const DefaultSymbols = "!@#$%^&*?"

var (
	// charSetCache caches the built character sets.
	// The key excludes Digit and RequireEach since they don't affect the set itself.
	charSetCache sync.Map
)

// charsetKey is the internal key used for caching; it only contains fields that affect the character set composition.
type charsetKey struct {
	UseNumbers    bool
	UseLowercase  bool
	UseUppercase  bool
	UseSymbols    bool
	CustomSymbols string
	ExcludeChars  string
}

// Code generates a random code. It returns the generated string and any error encountered.
// If crypto/rand fails repeatedly, an error is returned. Callers should handle it.
func Code(config Options) (string, error) {
	if config.Digit <= 0 {
		config.Digit = 6
	}

	charSet := buildCharSet(config)
	if len(charSet) == 0 {
		charSet = "0123456789" // fallback to digits only
	}

	// If RequireEach is set, ensure at least one character from each enabled category is present.
	if config.RequireEach {
		return codeWithRequire(config.Digit, charSet, config)
	}

	// Standard random generation
	b := make([]byte, config.Digit)
	for i := range b {
		n, err := randomInt(len(charSet))
		if err != nil {
			return "", err
		}
		b[i] = charSet[n]
	}
	return string(b), nil
}

// MustCode is like Code but panics on error. Useful for initialisation or scripting contexts.
func MustCode(config Options) string {
	code, err := Code(config)
	if err != nil {
		panic("RandomCode: failed to generate code: " + err.Error())
	}
	return code
}

// Number generates a random numeric code of the given length.
func Number(digit int) string {
	return MustCode(Options{
		Digit:      digit,
		UseNumbers: true,
	})
}

// Alpha generates a random alphabetic code of the given length, with control over case.
func Alpha(digit int, useLower, useUpper bool) string {
	return MustCode(Options{
		Digit:        digit,
		UseLowercase: useLower,
		UseUppercase: useUpper,
	})
}

// Mixed generates a random code using a mix of the specified character categories.
func Mixed(digit int, useNumbers, useLower, useUpper, useSymbols bool) string {
	return MustCode(Options{
		Digit:        digit,
		UseNumbers:   useNumbers,
		UseLowercase: useLower,
		UseUppercase: useUpper,
		UseSymbols:   useSymbols,
	})
}

// StrongPassword generates a high-strength password (default length 16, all categories, requires each).
func StrongPassword() string {
	return MustCode(Options{
		Digit:        16,
		UseNumbers:   true,
		UseLowercase: true,
		UseUppercase: true,
		UseSymbols:   true,
		RequireEach:  true,
		ExcludeChars: "0OIl1", // exclude commonly ambiguous characters
	})
}

// buildCharSet builds the character set based on the configuration (with caching).
func buildCharSet(config Options) string {
	key := charsetKey{
		UseNumbers:    config.UseNumbers,
		UseLowercase:  config.UseLowercase,
		UseUppercase:  config.UseUppercase,
		UseSymbols:    config.UseSymbols,
		CustomSymbols: config.CustomSymbols,
		ExcludeChars:  config.ExcludeChars,
	}

	if cached, ok := charSetCache.Load(key); ok {
		return cached.(string)
	}

	var sb strings.Builder
	if config.UseNumbers {
		sb.WriteString("0123456789")
	}
	if config.UseLowercase {
		sb.WriteString("abcdefghijklmnopqrstuvwxyz")
	}
	if config.UseUppercase {
		sb.WriteString("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	}
	if config.UseSymbols {
		if config.CustomSymbols != "" {
			sb.WriteString(config.CustomSymbols)
		} else {
			sb.WriteString(DefaultSymbols)
		}
	}

	result := sb.String()
	if config.ExcludeChars != "" {
		result = removeChars(result, config.ExcludeChars)
	}

	charSetCache.Store(key, result)
	return result
}

// codeWithRequire generates a code where each enabled category contributes at least one character.
func codeWithRequire(total int, charSet string, config Options) (string, error) {
	// Group characters by category after applying exclusions.
	groups := make([][]byte, 0)
	if config.UseNumbers {
		group := filterCharset("0123456789", config.ExcludeChars)
		if len(group) > 0 {
			groups = append(groups, group)
		}
	}
	if config.UseLowercase {
		group := filterCharset("abcdefghijklmnopqrstuvwxyz", config.ExcludeChars)
		if len(group) > 0 {
			groups = append(groups, group)
		}
	}
	if config.UseUppercase {
		group := filterCharset("ABCDEFGHIJKLMNOPQRSTUVWXYZ", config.ExcludeChars)
		if len(group) > 0 {
			groups = append(groups, group)
		}
	}
	if config.UseSymbols {
		symbols := DefaultSymbols
		if config.CustomSymbols != "" {
			symbols = config.CustomSymbols
		}
		group := filterCharset(symbols, config.ExcludeChars)
		if len(group) > 0 {
			groups = append(groups, group)
		}
	}

	// If no groups remain, fall back to standard generation.
	if len(groups) == 0 {
		code, err := Code(Options{
			Digit:        total,
			UseNumbers:   true,
			RequireEach:  false,
			ExcludeChars: config.ExcludeChars,
		})
		return code, err
	}

	result := make([]byte, total)
	// Step 1: pick one random character from each group and place it at the start
	for i, group := range groups {
		if i >= total {
			break
		}
		n, err := randomInt(len(group))
		if err != nil {
			return "", err
		}
		result[i] = group[n]
	}

	// Step 2: fill remaining positions with random characters from the full set
	for i := len(groups); i < total; i++ {
		n, err := randomInt(len(charSet))
		if err != nil {
			return "", err
		}
		result[i] = charSet[n]
	}

	// Step 3: shuffle to avoid the required characters being clustered at the beginning
	if err := shuffle(result); err != nil {
		return "", err
	}
	return string(result), nil
}

// filterCharset returns the characters from source after removing those present in exclude.
func filterCharset(source, exclude string) []byte {
	var res []byte
	for _, r := range source {
		if !strings.ContainsRune(exclude, r) {
			res = append(res, byte(r))
		}
	}
	return res
}

// shuffle performs a Fisher‑Yates shuffle on the byte slice.
func shuffle(b []byte) error {
	for i := len(b) - 1; i > 0; i-- {
		j, err := randomInt(i + 1)
		if err != nil {
			return err
		}
		b[i], b[j] = b[j], b[i]
	}
	return nil
}

// randomInt returns a random integer in [0, max) with up to 3 retries on failure.
func randomInt(max int) (int, error) {
	if max <= 0 {
		return 0, errors.New("randomInt: max must be > 0")
	}
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
		if err == nil {
			return int(n.Int64()), nil
		}
		lastErr = err
	}
	return 0, lastErr
}

// removeChars returns a new string with all characters from charsToRemove removed.
func removeChars(source, charsToRemove string) string {
	var sb strings.Builder
	for _, char := range source {
		if !strings.ContainsRune(charsToRemove, char) {
			sb.WriteRune(char)
		}
	}
	return sb.String()
}
