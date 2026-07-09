package rules

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

// TemplateSyntax validates Go template syntax by parsing the file with
// text/template. This catches unclosed {{if}}/{{range}} blocks, unknown
// actions, and other parse-time errors.
type TemplateSyntax struct{}

func (TemplateSyntax) Check(file string, content []byte, _ []string) []Violation {
	name := filepath.Base(file)
	_, err := template.New(name).Parse(string(content))
	if err == nil {
		return nil
	}
	return []Violation{{
		File:    file,
		Line:    lineFromTemplateError(name, err.Error()),
		Col:     1,
		Message: fmt.Sprintf("template syntax: %s", err.Error()),
	}}
}

// lineFromTemplateError extracts the line number from a text/template parse
// error. Parse errors have the format:
//
//	template: <name>:<line>: <description>
//	template: <name>:<line>:<col>: <description>
func lineFromTemplateError(name, msg string) int {
	prefix := "template: " + name + ":"
	idx := strings.Index(msg, prefix)
	if idx < 0 {
		return 1
	}
	rest := msg[idx+len(prefix):]
	end := strings.IndexAny(rest, ": ")
	if end < 0 {
		end = len(rest)
	}
	n, err := strconv.Atoi(rest[:end])
	if err != nil || n < 1 {
		return 1
	}
	return n
}
