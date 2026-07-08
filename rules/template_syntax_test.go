package rules

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTemplateSyntax(t *testing.T) {
	r := TemplateSyntax{}

	tests := []struct {
		name    string
		content string
		wantAny bool // true: expect at least one violation (line not asserted)
	}{
		{
			name:    "When the template has valid syntax, Then no violations are reported",
			content: "{{range $i,$e := .Items}}\n{{if $i}},{{end}}\n{\n  \"field\": \"{{.}}\"\n}\n{{end}}\n",
		},
		{
			// text/template reports "unexpected EOF" at the EOF line, not at
			// the opening {{if}}/{{range}}, so we only assert a violation exists.
			name:    "When an if block is missing its end, Then a violation is reported",
			content: "{{if .Enabled}}\n{\n  \"field\": \"value\"\n}\n",
			wantAny: true,
		},
		{
			name:    "When a range block is missing its end, Then a violation is reported",
			content: "{{range $i,$e := .Items}}\n{{if $i}},{{end}}\n{}\n",
			wantAny: true,
		},
		{
			name:    "When the template contains a comment block and valid syntax, Then no violations are reported",
			content: "{{/* topic query */}}\n{{ range .Items }}\n{}\n{{end}}\n",
		},
	}

	Convey("Given a TemplateSyntax rule", t, func() {
		for _, tc := range tests {
			tc := tc
			Convey(tc.name, func() {
				lines := strings.Split(tc.content, "\n")
				got := r.Check("test.tmpl", []byte(tc.content), lines)
				if tc.wantAny {
					So(got, ShouldNotBeEmpty)
					return
				}
				assertViolationLines(got, nil)
			})
		}
	})
}
