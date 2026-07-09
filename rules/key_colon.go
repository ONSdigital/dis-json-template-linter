package rules

import "fmt"

// KeyColon enforces JSON key-value colon spacing:
//   - no whitespace between the closing quote of a key and ":"
//   - exactly one space between ":" and the value
//
// Examples:
//
//	"key": "value"    - correct
//	"key" : "value"   - violation: space before ':'
//	"key":"value"     - violation: 0 spaces after ':'
//	"key":  "value"   - violation: 2 spaces after ':'
//
// Template blocks ({{ ... }}) are skipped when they appear outside strings.
// A ':' at the end of a line (value on the next line) is not checked.
type KeyColon struct{}

func (KeyColon) Check(file string, _ []byte, lines []string) []Violation {
	var violations []Violation
	for i, line := range lines {
		if i == len(lines)-1 && line == "" {
			continue
		}
		if isBlank(line) || isTemplateDirectiveLine(line) {
			continue
		}
		violations = append(violations, checkKeyColon(file, i+1, line)...)
	}
	return violations
}

// checkKeyColon scans a single line of the original (non-stripped) content.
// It skips top-level {{ }} template blocks, finds quoted JSON keys, and
// validates the spacing around their ':' separators.
func checkKeyColon(file string, lineNum int, line string) []Violation {
	var violations []Violation
	i := 0

	for i < len(line) {
		// Skip top-level {{ ... }} template blocks so that ':' inside template
		// actions (e.g. {{ range $k, $v := .Map }}) is never misread as a
		// JSON key-value separator.
		if i+1 < len(line) && line[i] == '{' && line[i+1] == '{' {
			i = skipTemplateBlock(line, i)
			continue
		}

		if line[i] != '"' {
			i++
			continue
		}

		// Scan a quoted string (potential JSON key).
		i++ // skip opening "
		i = scanQuotedKey(line, i)
		if i >= len(line) {
			break // unclosed string - stop scanning this line
		}
		closeQuoteIdx := i
		i++ // skip closing "

		// Collect any whitespace between the closing '"' and the next ':'.
		j := i
		spacesBefore := 0
		for j < len(line) && (line[j] == ' ' || line[j] == '\t') {
			spacesBefore++
			j++
		}

		if j >= len(line) || line[j] != ':' {
			// This string is a value, not a key - resume scanning from j.
			i = j
			continue
		}

		// Found a key-value ':'.
		colonIdx := j

		// Check 1: no whitespace before ':'.
		if spacesBefore > 0 {
			violations = append(violations, Violation{
				File:    file,
				Line:    lineNum,
				Col:     closeQuoteIdx + 2, // 1-based col right after closing "
				Message: fmt.Sprintf("key colon: expected no space before \":\", got %d", spacesBefore),
			})
		}

		// Check 2: exactly one space after ':'.
		// Count spaces; a ':' at end-of-line (value on next line) is not checked.
		k := colonIdx + 1
		spacesAfter := 0
		for k < len(line) && line[k] == ' ' {
			spacesAfter++
			k++
		}
		if k < len(line) && spacesAfter != 1 {
			violations = append(violations, Violation{
				File:    file,
				Line:    lineNum,
				Col:     colonIdx + 1, // 1-based col of ':'
				Message: fmt.Sprintf("key colon: expected exactly 1 space after \":\", got %d", spacesAfter),
			})
		}

		i = colonIdx + 1
	}
	return violations
}

// skipTemplateBlock advances past the {{ ... }} block starting at line[i].
// It assumes line[i] == '{' && line[i+1] == '{' .
func skipTemplateBlock(line string, i int) int {
	depth := 1
	i += 2
	for i < len(line) && depth > 0 {
		if i+1 < len(line) && line[i] == '{' && line[i+1] == '{' {
			depth++
			i += 2
		} else if i+1 < len(line) && line[i] == '}' && line[i+1] == '}' {
			depth--
			i += 2
		} else {
			i++
		}
	}
	return i
}

// scanQuotedKey scans a JSON-style quoted string starting at line[i]
// (the opening '"' has already been consumed). Returns the 0-based index of
// the closing '"', or len(line) if the string is unclosed.
func scanQuotedKey(line string, i int) int {
	for i < len(line) {
		if line[i] == '\\' && i+1 < len(line) {
			i += 2 // skip escape sequence
			continue
		}
		if line[i] == '"' {
			return i
		}
		i++
	}
	return i
}
