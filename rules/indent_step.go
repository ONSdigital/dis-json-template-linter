package rules

import "fmt"

// IndentStep checks that the indentation level changes by at most one step
// between consecutive non-blank, non-template-directive lines. With the
// default indent_size of 2 the only valid changes are -2, 0, or +2.
//
// Blank lines and template directive lines (e.g. {{range ...}}, {{end}}) are
// skipped when building the comparison sequence, because they sit at indent 0
// regardless of JSON nesting depth and would generate false positives.
//
// The rule is skipped entirely when Size <= 0.
type IndentStep struct {
	Size int
}

func (r IndentStep) Check(file string, _ []byte, lines []string) []Violation {
	if r.Size <= 0 {
		return nil
	}

	var violations []Violation
	prev := -1 // leading-char count of the last content line; -1 = none seen yet

	for i, line := range lines {
		if i == len(lines)-1 && line == "" {
			continue
		}
		if isBlank(line) || isTemplateDirectiveLine(line) {
			continue
		}

		cur := len(leadingWhitespace(line))

		if prev >= 0 {
			if ok, diff := checkIndentChange(prev, cur, r.Size); !ok {
				violations = append(violations, Violation{
					File:    file,
					Line:    i + 1,
					Col:     1,
					Message: generateMessage(diff, r.Size),
				})
			}
		}

		prev = cur
	}
	return violations
}

// checkIndentChange returns true if the change in indentation between two
// consecutive content lines is valid (0, +size, or -size). It also returns the
// actual difference
func checkIndentChange(prev, cur, size int) (ok bool, diff int) {
	diff = cur - prev
	ok = diff == 0 || diff == size || diff == -size
	return
}

// generateMessage returns a human-readable message describing the violation
func generateMessage(diff, size int) string {
	var message string
	if diff > 0 {
		message = fmt.Sprintf(
			"indent step: indented by %d (expected 0 or +%d from previous content line)",
			diff, size,
		)
	} else {
		message = fmt.Sprintf(
			"indent step: de-indented by %d (expected 0 or -%d from previous content line)",
			-diff, size,
		)
	}
	return message
}
