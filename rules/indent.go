package rules

import (
	"fmt"
	"strings"
)

// IndentStyle specifies the whitespace character used for indentation.
type IndentStyle string

const (
	IndentStyleSpace IndentStyle = "space"
	IndentStyleTab   IndentStyle = "tab"
)

// Indent checks that every non-template-directive, non-blank line uses
// consistent indentation.
//
// Two sub-checks are performed:
//  1. Style: leading whitespace must use the configured character (space or tab).
//  2. Size:  the number of leading characters must be a multiple of IndentSize.
//
// "Template directive lines" - lines whose non-template content is empty or
// only commas - are exempt. Examples: {{range ...}}, {{end}}, {{if $i}},{{end}}.
type Indent struct {
	Size  int         // number of indent characters per level (e.g. 2)
	Style IndentStyle // IndentStyleSpace or IndentStyleTab
}

func (r Indent) Check(file string, _ []byte, lines []string) []Violation {
	var violations []Violation
	for i, line := range lines {
		// Skip the synthetic empty element after a trailing newline.
		if i == len(lines)-1 && line == "" {
			continue
		}
		if isBlank(line) || isTemplateDirectiveLine(line) {
			continue
		}
		leading := leadingWhitespace(line)
		if leading == "" {
			continue
		}
		if v := r.checkStyle(file, i+1, leading); v != nil {
			violations = append(violations, *v)
			// Skip the size check when the style is wrong: counting chars of the
			// wrong type (e.g. tabs when spaces are expected) is meaningless.
			continue
		}
		if v := r.checkSize(file, i+1, leading); v != nil {
			violations = append(violations, *v)
		}
	}
	return violations
}

func (r Indent) checkStyle(file string, lineNum int, leading string) *Violation {
	wantByte := byte(' ')
	wantName := "spaces"
	badName := "tabs"
	if r.Style == IndentStyleTab {
		wantByte = '\t'
		wantName = "tabs"
		badName = "spaces"
	}
	for i := range len(leading) {
		if leading[i] != wantByte {
			return &Violation{
				File:    file,
				Line:    lineNum,
				Col:     i + 1,
				Message: fmt.Sprintf("indent style: expected %s but found %s", wantName, badName),
			}
		}
	}
	return nil
}

func (r Indent) checkSize(file string, lineNum int, leading string) *Violation {
	if r.Size <= 0 {
		return nil
	}
	n := len(leading)
	if n%r.Size != 0 {
		return &Violation{
			File:    file,
			Line:    lineNum,
			Col:     1,
			Message: fmt.Sprintf("indent size: %d leading character(s) is not a multiple of %d", n, r.Size),
		}
	}
	return nil
}

// isBlank reports whether a line contains only whitespace.
func isBlank(line string) bool {
	return strings.TrimSpace(line) == ""
}

// isTemplateDirectiveLine reports whether a line is a Go template control
// directive with no meaningful JSON content. A line qualifies when:
//   - After trimming leading whitespace, it starts with "{{".
//   - After stripping all {{ ... }} tokens, the remaining trimmed text is empty
//     or consists only of commas (separator tokens emitted by templates).
func isTemplateDirectiveLine(line string) bool {
	trimmed := strings.TrimLeft(line, " \t")
	if !strings.HasPrefix(trimmed, "{{") {
		return false
	}
	remainder := stripTemplateTokens(trimmed)
	remainder = strings.Trim(remainder, " \t,")
	return remainder == ""
}

// stripTemplateTokens removes all {{ ... }} blocks from s, returning the
// interleaved non-template content.
func stripTemplateTokens(s string) string {
	var b strings.Builder
	depth := 0
	i := 0
	for i < len(s) {
		if i+1 < len(s) && s[i] == '{' && s[i+1] == '{' {
			depth++
			i += 2
			continue
		}
		if i+1 < len(s) && s[i] == '}' && s[i+1] == '}' && depth > 0 {
			depth--
			i += 2
			continue
		}
		if depth == 0 {
			b.WriteByte(s[i])
		}
		i++
	}
	return b.String()
}

// leadingWhitespace returns the prefix of line consisting solely of spaces and tabs.
func leadingWhitespace(line string) string {
	for i, ch := range line {
		if ch != ' ' && ch != '\t' {
			return line[:i]
		}
	}
	return line // entire line is whitespace (handled by isBlank)
}
