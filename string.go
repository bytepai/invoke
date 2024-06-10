package invoke

import (
	"strings"
	"unicode/utf8"
)

// String is a package-level variable representing a string handler.
var String stringHandler

// stringHandler is a struct for string manipulation.
type stringHandler struct{}


// FixUTF8 ensures the input string is valid UTF-8.
func (stringHandler)FixUTF8(s string) string {
	if utf8.ValidString(s) {
		return s
	}
	v := make([]rune, 0, len(s))
	for i, r := range s {
		if r == utf8.RuneError {
			_, size := utf8.DecodeRuneInString(s[i:])
			if size == 1 {
				continue // Skip invalid UTF-8 sequences
			}
		}
		v = append(v, r)
	}
	return string(v)
}

// ToLowerCamelCase converts a string to lower camel case.
func (stringHandler)ToLowerCamelCase(s string) string {
	parts := strings.Fields(s)
	for i := range parts {
		if i == 0 {
			parts[i] = strings.ToLower(parts[i])
		} else {
			parts[i] = strings.Title(parts[i])
		}
	}
	return strings.Join(parts, "")
}

// ToUpperCamelCase converts a string to upper camel case.
func (stringHandler)ToUpperCamelCase(s string) string {
	parts := strings.Fields(s)
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, "")
}

// ReverseString reverses a string.
func (stringHandler)ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// isDigit checks if a given rune is a numeric digit (0-9).
func (stringHandler)isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// isNL checks if a given rune is a newline character ('\n' or '\r').
func (stringHandler)isNL(r rune) bool {
	return r == '\n' || r == '\r'
}

// isWhitespace checks if a given rune is a whitespace character (tab or space).
func (stringHandler)isWhitespace(r rune) bool {
	return r == '\t' || r == ' '
}

// isHexadecimal checks if a given rune is a valid hexadecimal digit (0-9, a-f, A-F).
func (stringHandler)isHexadecimal(r rune) bool {
	return (r >= '0' && r <= '9') ||
		(r >= 'a' && r <= 'f') ||
		(r >= 'A' && r <= 'F')
}

// isBareKeyChar checks if a given rune is a valid character for a bare key (A-Z, a-z, 0-9, '_', or '-').
func (stringHandler)isBareKeyChar(r rune) bool {
	return (r >= 'A' && r <= 'Z') ||
		(r >= 'a' && r <= 'z') ||
		(r >= '0' && r <= '9') ||
		r == '_' ||
		r == '-'
}

// isValidHex checks if a string contains only valid hexadecimal characters.
func (sh stringHandler)isValidHex(s string) bool {
	for _, r := range s {
		if !sh.isHexadecimal(r) {
			return false
		}
	}
	return true
}
