package rules

import (
	"fmt"
	"strings"
)

// TemplateAction enforces that Go template actions which take a quoted name
// followed by a pipeline argument - specifically "template" and "block" - have
// exactly one space between the closing quote of the name and the argument.
//
// Correct:
//
//	{{ template "header.tmpl" . }}
//	{{ block "content" .Data }}
//	{{ template "header.tmpl" }}   - no argument at all (space before }} is
//	                                 enforced separately by TemplateSpacing)
//
// Incorrect:
//
//	{{ template "header.tmpl". }}  - no space before argument
//	{{ block "content"$var }}      - no space before argument
type TemplateAction struct{}

func (TemplateAction) Check(file string, _ []byte, lines []string) []Violation {
	var violations []Violation
	for i, line := range lines {
		if i == len(lines)-1 && line == "" {
			continue
		}
		if isBlank(line) {
			continue
		}
		violations = append(violations, checkTemplateAction(file, i+1, line)...)
	}
	return violations
}

// checkTemplateAction scans every {{ }} block on the line. When the action
// keyword is "template" or "block", it locates the quoted name and checks that
// a space follows the closing quote.
func checkTemplateAction(file string, lineNum int, line string) []Violation {
	var violations []Violation
	i := 0
	for i < len(line) {
		// Look for {{ opening.
		if i+1 >= len(line) || line[i] != '{' || line[i+1] != '{' {
			i++
			continue
		}

		j := i + 2 // first character after {{

		// Find matching }}.
		closeIdx := findTemplateClose(line, j)
		if closeIdx < 0 {
			i = j
			continue
		}

		// Locate the keyword and quoted name inside this block.
		keyword, closeQuoteIdx, found := findActionName(line, j, closeIdx)
		if !found {
			i = closeIdx + 2
			continue
		}

		// Check: the character immediately after the closing quote must be a
		// space or tab. If the action has no argument (e.g. {{ template "x" }})
		// the TemplateSpacing rule already enforces a space before }}, so we
		// only flag when the non-space character is not } or -.
		k := closeQuoteIdx + 1
		if k < closeIdx && line[k] != ' ' && line[k] != '\t' {
			violations = append(violations, Violation{
				File: file,
				Line: lineNum,
				Col:  closeQuoteIdx + 2, // 1-based column of the offending character
				Message: fmt.Sprintf(
					"template action: expected a space after the %s name %q",
					keyword,
					strings.TrimSpace(line[j:closeQuoteIdx+1]), // include " for context
				),
			})
		}

		i = closeIdx + 2
	}
	return violations
}

// matchWord reports whether line[pos:] starts with word followed by whitespace
// or the end of the string - i.e. that word is not a prefix of a longer token.
func matchWord(line string, pos int, word string) bool {
	end := pos + len(word)
	if end > len(line) {
		return false
	}
	if line[pos:end] != word {
		return false
	}
	return end == len(line) || line[end] == ' ' || line[end] == '\t'
}

// findActionName extracts the keyword ("template" or "block") and the quoted
// name from a {{ }} action block (line[j:closeIdx]). It returns the keyword,
// the 0-based index of the closing quote of the name, and true on success.
// Returns ("", -1, false) when the block does not begin with a recognised
// keyword or when the quoted name is missing or unclosed.
func findActionName(line string, j, closeIdx int) (keyword string, closeQuoteIdx int, found bool) {
	k := j
	// Skip optional left-trim modifier.
	if k < closeIdx && line[k] == '-' {
		k++
	}
	// Skip leading whitespace.
	for k < closeIdx && (line[k] == ' ' || line[k] == '\t') {
		k++
	}
	// Identify keyword.
	switch {
	case matchWord(line, k, "template"):
		keyword = "template"
		k += len("template")
	case matchWord(line, k, "block"):
		keyword = "block"
		k += len("block")
	default:
		return "", -1, false
	}
	// Keyword must be followed by whitespace.
	if k >= closeIdx || (line[k] != ' ' && line[k] != '\t') {
		return keyword, -1, false
	}
	for k < closeIdx && (line[k] == ' ' || line[k] == '\t') {
		k++
	}
	// Expect opening quote.
	if k >= closeIdx || line[k] != '"' {
		return keyword, -1, false
	}
	k++ // skip opening "
	// Scan the quoted name.
	closeQuoteIdx, ok := scanActionName(line, k, closeIdx)
	if !ok {
		return keyword, -1, false
	}
	return keyword, closeQuoteIdx, true
}

// scanActionName scans a quoted template/block name starting at line[k] (the opening '"'
// must already have been consumed). Returns the index of the closing '"' and true,
// or -1 and false if the string is not closed before closeIdx.
func scanActionName(line string, k, closeIdx int) (int, bool) {
	for k < closeIdx {
		if line[k] == '\\' && k+1 < closeIdx {
			k += 2
			continue
		}
		if line[k] == '"' {
			return k, true
		}
		k++
	}
	return -1, false
}
