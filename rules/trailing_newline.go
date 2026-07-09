package rules

import "bytes"

// TrailingNewline flags files that do not end with exactly one newline, or
// that end with a blank trailing line.
type TrailingNewline struct{}

func (TrailingNewline) Check(file string, content []byte, lines []string) []Violation {
	if len(content) == 0 {
		return nil
	}
	if !bytes.HasSuffix(content, []byte("\n")) {
		lastLine := len(lines)
		return []Violation{{
			File:    file,
			Line:    lastLine,
			Col:     len(lines[lastLine-1]) + 1,
			Message: "no newline at end of file",
		}}
	}
	if bytes.HasSuffix(content, []byte("\n\n")) {
		// Point to the blank line itself (second-to-last element after split).
		blankLine := len(lines) - 1
		return []Violation{{
			File:    file,
			Line:    blankLine,
			Col:     1,
			Message: "trailing blank line at end of file",
		}}
	}
	return nil
}
