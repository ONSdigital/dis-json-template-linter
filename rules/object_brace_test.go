package rules

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestObjectBrace_SameLine(t *testing.T) {
	r := ObjectBrace{}

	tests := []struct {
		name      string
		content   string
		wantLines []int
	}{
		{
			name:      "When the opening brace appears on the same line as its key, Then no violations are reported",
			content:   "{\n  \"match\": {\n    \"field\": \"value\"\n  }\n}\n",
			wantLines: nil,
		},
		{
			name:      "When the opening brace appears on the next line after its key, Then a violation is reported on line 3",
			content:   "{\n  \"match\":\n  {\n    \"field\": \"value\"\n  }\n}\n",
			wantLines: []int{3},
		},
		{
			name:      "When a lone brace is used as an array element, Then it is not flagged",
			content:   "[\n  {\n    \"field\": \"value\"\n  }\n]\n",
			wantLines: nil,
		},
		{
			name:      "When a lone brace has no prior content line, Then it is not flagged",
			content:   "{{range $i,$e := .Items}}\n{{if $i}},{{end}}\n{\n  \"field\": \"{{.}}\"\n}\n{{end}}\n",
			wantLines: nil,
		},
	}

	Convey("Given an ObjectBrace rule", t, func() {
		for _, tc := range tests {
			tc := tc
			Convey(tc.name, func() {
				lines := strings.Split(tc.content, "\n")
				got := r.Check("test.tmpl", []byte(tc.content), lines)
				assertViolationLines(got, tc.wantLines)
			})
		}
	})
}
