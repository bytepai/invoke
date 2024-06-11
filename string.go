package invoke

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// String is a package-level variable representing a string handler.
var String stringHandler

// stringHandler is a struct for string manipulation.
type stringHandler struct{}

// FixUTF8 ensures the input string is valid UTF-8.
func (stringHandler) FixUTF8(s string) string {
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
func (stringHandler) ToLowerCamelCase(s string) string {
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
func (stringHandler) ToUpperCamelCase(s string) string {
	parts := strings.Fields(s)
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, "")
}

// ReverseString reverses the given string.
func (stringHandler) ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// IsPalindrome checks if a string is a palindrome.
func (stringHandler) IsPalindrome(s string) bool {
	s = strings.ToLower(s)
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		if runes[i] != runes[j] {
			return false
		}
	}
	return true
}

// RemoveWhitespace removes all whitespace from a string.
func (stringHandler) RemoveWhitespace(s string) string {
	return strings.Join(strings.Fields(s), "")
}

// isDigit checks if a given rune is a numeric digit (0-9).
func (stringHandler) isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// CountVowels counts the number of vowels in a string.
func (stringHandler) CountVowels(s string) int {
	count := 0
	vowels := "aeiouAEIOU"
	for _, char := range s {
		if strings.ContainsRune(vowels, char) {
			count++
		}
	}
	return count
}

// ReplaceVowels replaces all vowels in a string with a given character.
func (stringHandler) ReplaceVowels(s, replacement string) string {
	vowels := "aeiouAEIOU"
	for _, char := range vowels {
		s = strings.ReplaceAll(s, string(char), replacement)
	}
	return s
}

// Capitalize capitalizes the first letter of each word in a string.
func (stringHandler) Capitalize(s string) string {
	return strings.Title(s)
}

// CountWords counts the number of words in a string.
func (stringHandler) CountWords(s string) int {
	return len(strings.Fields(s))
}

// IsAnagram checks if two strings are anagrams.
func (stringHandler) IsAnagram(s1, s2 string) bool {
	runeCount1 := make(map[rune]int)
	runeCount2 := make(map[rune]int)

	for _, char := range s1 {
		runeCount1[char]++
	}
	for _, char := range s2 {
		runeCount2[char]++
	}

	return fmt.Sprintf("%v", runeCount1) == fmt.Sprintf("%v", runeCount2)
}

// GenerateSlug generates a URL-friendly slug from a string.
func (stringHandler) GenerateSlug(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.Map(func(r rune) rune {
		if r == '-' || r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return -1
	}, s)
	return s
}

// RemoveDuplicates removes duplicate characters from a string.
func (stringHandler) RemoveDuplicates(s string) string {
	seen := make(map[rune]bool)
	var result []rune
	for _, char := range s {
		if !seen[char] {
			seen[char] = true
			result = append(result, char)
		}
	}
	return string(result)
}

// Truncate truncates a string to a specified length, adding "..." if truncated.
func (stringHandler) Truncate(s string, maxLength int) string {
	if utf8.RuneCountInString(s) <= maxLength {
		return s
	}
	return string([]rune(s)[:maxLength]) + "..."
}

// CountSubstring counts the occurrences of a substring in a string.
func (stringHandler) CountSubstring(s, substr string) int {
	return strings.Count(s, substr)
}

// RemoveNonAlphanumeric removes all non-alphanumeric characters from a string.
func (stringHandler) RemoveNonAlphanumeric(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return -1
	}, s)
}

// Rot13 applies the ROT13 algorithm to a string.
func (stringHandler) Rot13(s string) string {
	return strings.Map(func(r rune) rune {
		if 'a' <= r && r <= 'z' {
			return 'a' + (r-'a'+13)%26
		}
		if 'A' <= r && r <= 'Z' {
			return 'A' + (r-'A'+13)%26
		}
		return r
	}, s)
}

// RightPad pads a string to the right with a specified character to a specified length.
func (stringHandler) RightPad(s string, padChar rune, length int) string {
	if len(s) >= length {
		return s
	}
	return s + strings.Repeat(string(padChar), length-len(s))
}

// LeftPad pads a string to the left with a specified character to a specified length.
func (stringHandler) LeftPad(s string, padChar rune, length int) string {
	if len(s) >= length {
		return s
	}
	return strings.Repeat(string(padChar), length-len(s)) + s
}

// CenterPad centers a string and pads it with a specified character to a specified length.
func (stringHandler) CenterPad(s string, padChar rune, length int) string {
	if len(s) >= length {
		return s
	}
	padding := (length - len(s)) / 2
	leftPadding := strings.Repeat(string(padChar), padding)
	rightPadding := strings.Repeat(string(padChar), length-len(s)-padding)
	return leftPadding + s + rightPadding
}

// isNL checks if a given rune is a newline character ('\n' or '\r').
func (stringHandler) isNL(r rune) bool {
	return r == '\n' || r == '\r'
}

// isWhitespace checks if a given rune is a whitespace character (tab or space).
func (stringHandler) isWhitespace(r rune) bool {
	return r == '\t' || r == ' '
}

// isHexadecimal checks if a given rune is a valid hexadecimal digit (0-9, a-f, A-F).
func (stringHandler) isHexadecimal(r rune) bool {
	return (r >= '0' && r <= '9') ||
		(r >= 'a' && r <= 'f') ||
		(r >= 'A' && r <= 'F')
}

// isBareKeyChar checks if a given rune is a valid character for a bare key (A-Z, a-z, 0-9, '_', or '-').
func (stringHandler) isBareKeyChar(r rune) bool {
	return (r >= 'A' && r <= 'Z') ||
		(r >= 'a' && r <= 'z') ||
		(r >= '0' && r <= '9') ||
		r == '_' ||
		r == '-'
}

// isValidHex checks if a string contains only valid hexadecimal characters.
func (sh stringHandler) isValidHex(s string) bool {
	for _, r := range s {
		if !sh.isHexadecimal(r) {
			return false
		}
	}
	return true
}

// MD5HashUpper Calculate the MD5 hash value of a string (in uppercase)
func (sh stringHandler) MD5HashUpper(data string) string {
	hasher := md5.New()
	hasher.Write([]byte(data))
	hashBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	// Convert the hash string to uppercase
	return strings.ToUpper(hashString)
}

// MD5HashLower Calculate the MD5 hash value of a string (in lowercase)
func (sh stringHandler) MD5HashLower(data string) string {
	hasher := md5.New()
	hasher.Write([]byte(data))
	hashBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	// Return the hash string in lowercase
	return hashString
}
