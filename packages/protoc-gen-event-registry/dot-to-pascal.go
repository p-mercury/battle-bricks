package main

import (
	"strings"
	"unicode"
)

func DotToPascal(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	parts := strings.Split(s, ".")
	var b strings.Builder
	b.Grow(len(s))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		start := -1
		flush := func(end int) {
			if start < 0 || end <= start {
				return
			}

			first := true
			for _, r := range part[start:end] {
				if first {
					b.WriteRune(unicode.ToUpper(r))
					first = false
				} else {
					b.WriteRune(unicode.ToLower(r))
				}
			}
			start = -1
		}

		for i, r := range part {
			if unicode.IsLetter(r) || unicode.IsDigit(r) {
				if start < 0 {
					start = i
				}
			} else {
				flush(i)
			}
		}
		flush(len(part))
	}

	return b.String()
}
