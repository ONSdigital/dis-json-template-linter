package rules

import (
	"fmt"
	"strings"
)

// ObjectBrace enforces that when an object-valued key has its opening "{"
// on the same line as the key (the SameLine check). Colon spacing is enforced
// separately by the KeyColon rule.
type ObjectBrace struct{}

func (ObjectBrace) Check(file string, _ []byte, lines []string) []Violation {
	var violations []Violation

	// prevContent tracks the last non-blank, non-directive line so the SameLine
	// check can look back one logical line.
	prevContent := ""
	prevLineNum := 0

	for i, line := range lines {
		if i == len(lines)-1 && line == "" {
			continue
		}
		if isBlank(line) || isTemplateDirectiveLine(line) {
			continue
		}

		lineNum := i + 1

		// SameLine: "{" on its own line when the preceding content line
		// ended with ":" (the brace belongs on that line).
		stripped := strings.TrimSpace(stripTemplateTokens(line))
		if stripped == "{" && prevContent != "" {
			prevStripped := strings.TrimRight(stripTemplateTokens(prevContent), " \t")
			if strings.HasSuffix(prevStripped, ":") {
				violations = append(violations, Violation{
					File:    file,
					Line:    lineNum,
					Col:     1,
					Message: fmt.Sprintf(`object brace: "{" should be on the same line as the key (line %d)`, prevLineNum),
				})
			}
		}

		prevContent = line
		prevLineNum = lineNum
	}
	return violations
}
