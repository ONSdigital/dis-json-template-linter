package rules

import (
	"fmt"
	"strings"
)

// TrailingWhitespace flags lines that end with spaces or tabs before the newline.
type TrailingWhitespace struct{}

func (TrailingWhitespace) Check(file string, _ []byte, lines []string) []Violation {
	var violations []Violation
	for i, line := range lines {
		// The last element produced by strings.Split on a "\n"-terminated file
		// is always an empty string - skip it; it has no trailing whitespace.
		if i == len(lines)-1 && line == "" {
			continue
		}
		trimmed := strings.TrimRight(line, " \t")
		if trimmed == line {
			continue
		}
		violations = append(violations, Violation{
			File:    file,
			Line:    i + 1,
			Col:     len(trimmed) + 1,
			Message: fmt.Sprintf("trailing whitespace (%d character(s))", len(line)-len(trimmed)),
		})
	}
	return violations
}
