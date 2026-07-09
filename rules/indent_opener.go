package rules

import (
	"strings"
)

// IndentOpener checks that an increase in indentation only occurs when the
// previous non-blank, non-template-directive line ends with an opening brace
// '{' or opening bracket '['.
type IndentOpener struct{}

func (IndentOpener) Check(file string, _ []byte, lines []string) []Violation {
	var violations []Violation

	prevLine := ""
	prevIndent := -1

	for i, line := range lines {
		// Skip the synthetic empty element after a trailing newline.
		if i == len(lines)-1 && line == "" {
			continue
		}
		if isBlank(line) || isTemplateDirectiveLine(line) {
			continue
		}

		curIndent := len(leadingWhitespace(line))

		if prevIndent >= 0 && curIndent > prevIndent {
			trimmed := strings.TrimRight(prevLine, " \t")
			if trimmed != "" {
				last := trimmed[len(trimmed)-1]
				if last != '{' && last != '[' {
					violations = append(violations, Violation{
						File:    file,
						Line:    i + 1,
						Col:     1,
						Message: "indent opener: unexpected indent - previous content line does not end with '{' or '['",
					})
				}
			}
		}

		prevLine = line
		prevIndent = curIndent
	}
	return violations
}
