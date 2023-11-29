package gen

import (
	"strings"
)

func SnakeCase(s string) string {
	return delimiterCase(s, '_')
}

func toLowerRune(ch rune) rune {
	if ch >= 'A' && ch <= 'Z' {
		return ch + 32
	}
	return ch
}

func isSpaceRune(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isDelimiter(ch rune) bool {
	return ch == '-' || ch == '_' || isSpaceRune(ch)
}

func ToLowerFirst(s string) string {
	if s == "" {
		return s
	}
	if s[0] >= 'A' && s[0] <= 'Z' {
		n := s[0] + 32
		b := strings.Builder{}
		b.WriteByte(n)
		b.WriteString(s[1:])
		return b.String()
	}
	return s
}

func isUpperRune(ch rune) bool {
	return ch >= 'A' && ch <= 'Z'
}

func isLowerRune(ch rune) bool {
	return ch >= 'a' && ch <= 'z'
}

// delimiterCase converts a string into snake_case or kebab-case depending on the delimiter passed
// as second argument. When upperCase is true the result will be UPPER_SNAKE_CASE or UPPER-KEBAB-CASE.
func delimiterCase(s string, delimiter rune) string {
	s = strings.TrimSpace(s)
	//buffer := make([]rune, 0, len(s)+3)
	buffer := &strings.Builder{}
	adjustCase := toLowerRune

	var prev rune
	var curr rune
	for _, next := range s {
		if isDelimiter(curr) {
			if !isDelimiter(prev) {
				buffer.WriteRune(delimiter)
			}
		} else if isUpperRune(curr) {
			if isLowerRune(prev) || (isUpperRune(prev) && isLowerRune(next)) {
				buffer.WriteRune(delimiter)
			}
			buffer.WriteRune(adjustCase(curr))
		} else if curr != 0 {
			buffer.WriteRune(adjustCase(curr))
		}
		prev = curr
		curr = next
	}

	if len(s) > 0 {
		if isUpperRune(curr) && isLowerRune(prev) && prev != 0 {
			buffer.WriteRune(delimiter)
		}
		buffer.WriteRune(adjustCase(curr))
	}

	return buffer.String()
}
