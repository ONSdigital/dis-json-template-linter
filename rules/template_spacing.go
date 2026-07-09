package rules

import "fmt"

// TemplateSpacing enforces that every {{ }} template action has a space
// immediately inside the opening and closing delimiters.
//
// Allowed forms:
//   - {{ action }}         - standard with spaces
//   - {{- action -}}       - whitespace-trim modifiers followed/preceded by spaces
//   - {{- action }}        - left-trim only
//   - {{ action -}}        - right-trim only
//   - {{/* comment */}}    - comments (no inner space required)
//
// Flagged forms:
//   - {{action}}           - no space after {{ or before }}
//   - {{-action}}          - no space after the trim modifier
//   - {{ action-}}         - no space before the trim modifier
type TemplateSpacing struct{}

func (TemplateSpacing) Check(file string, _ []byte, lines []string) []Violation {
	var violations []Violation
	for i, line := range lines {
		if i == len(lines)-1 && line == "" {
			continue
		}
		violations = append(violations, checkTemplateSpacing(file, i+1, line)...)
	}
	return violations
}

// checkTemplateSpacing scans a single line for every {{ ... }} block and
// verifies internal spacing on both the opening and closing delimiter.
func checkTemplateSpacing(file string, lineNum int, line string) []Violation {
	var violations []Violation
	i := 0
	for i < len(line) {
		// Find the next {{ opening.
		if i+1 >= len(line) || line[i] != '{' || line[i+1] != '{' {
			i++
			continue
		}

		openCol := i + 1 // 1-based column of {{
		j := i + 2       // first character after {{

		// Detect optional left-trim modifier (-).
		hasLeftTrim := j < len(line) && line[j] == '-'
		if hasLeftTrim {
			j++
		}

		// Detect comment opening /*, which is exempt from the space requirement.
		isComment := j+1 < len(line) && line[j] == '/' && line[j+1] == '*'

		// --- Opening check ---
		if v, ok := openDelimViolation(file, lineNum, openCol, j, line, hasLeftTrim, isComment); ok {
			violations = append(violations, v)
		}

		// Scan forward to find the matching }}, tracking string literals so
		// that }} inside a quoted string (e.g. {{ printf "val}}" }}) is skipped.
		closePos := findTemplateClose(line, j)
		if closePos < 0 {
			// No closing }} on this line; advance past the {{ and continue.
			i = j
			continue
		}

		// --- Closing check ---
		m := closePos - 1 // character immediately before }}

		// Detect optional right-trim modifier (-).
		hasRightTrim := m >= 0 && line[m] == '-'
		if hasRightTrim {
			m-- // character before the -
		}

		// Detect comment closing */, which is exempt from the space requirement.
		isCommentEnd := m >= 1 && line[m-1] == '*' && line[m] == '/'

		if v, ok := closeDelimViolation(file, lineNum, closePos, m, line, hasRightTrim, isCommentEnd); ok {
			violations = append(violations, v)
		}

		// Advance past the closing }}.
		i = closePos + 2
	}
	return violations
}

// findTemplateClose scans line from start looking for }}, skipping }} that
// appear inside double-quoted strings. Returns the 0-based index of the first
// } in the closing }}, or -1 if not found.
func findTemplateClose(line string, start int) int {
	inString := false
	for k := start; k < len(line); k++ {
		if inString {
			if line[k] == '\\' && k+1 < len(line) {
				k++ // skip escaped character
				continue
			}
			if line[k] == '"' {
				inString = false
			}
			continue
		}
		if line[k] == '"' {
			inString = true
			continue
		}
		if k+1 < len(line) && line[k] == '}' && line[k+1] == '}' {
			return k
		}
	}
	return -1
}

// openDelimViolation returns a spacing Violation for a {{ opening delimiter when
// the required inner space is missing. Returns (zero, false) when no violation.
func openDelimViolation(file string, lineNum, openCol, j int, line string, hasLeftTrim, isComment bool) (Violation, bool) {
	if isComment || (j < len(line) && line[j] == ' ') {
		return Violation{}, false
	}
	open := "{{"
	if hasLeftTrim {
		open = "{{-"
	}
	return Violation{
		File:    file,
		Line:    lineNum,
		Col:     openCol,
		Message: fmt.Sprintf("template spacing: expected a space after %q", open),
	}, true
}

// closeDelimViolation returns a spacing Violation for a }} closing delimiter when
// the required inner space is missing. Returns (zero, false) when no violation.
func closeDelimViolation(file string, lineNum, closePos, m int, line string, hasRightTrim, isCommentEnd bool) (Violation, bool) {
	if isCommentEnd || (m >= 0 && line[m] == ' ') {
		return Violation{}, false
	}
	closingDelim := "}}"
	if hasRightTrim {
		closingDelim = "-}}"
	}
	return Violation{
		File:    file,
		Line:    lineNum,
		Col:     closePos + 1,
		Message: fmt.Sprintf("template spacing: expected a space before %q", closingDelim),
	}, true
}
