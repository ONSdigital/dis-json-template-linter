package rules

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIndentOpener(t *testing.T) {
	r := IndentOpener{}

	tests := []struct {
		name      string
		content   string
		wantLines []int
	}{
		// --- No violations ---
		{
			name:      "When an indented line follows a '{' opener, Then no violations are reported",
			content:   "{\n  \"key\": \"value\"\n}\n",
			wantLines: nil,
		},
		{
			name:      "When an indented line follows a '[' opener, Then no violations are reported",
			content:   "[\n  \"item\"\n]\n",
			wantLines: nil,
		},
		{
			name:      "When a nested object follows a '{' opener, Then no violations are reported",
			content:   "{\n  \"a\": {\n    \"b\": \"c\"\n  }\n}\n",
			wantLines: nil,
		},
		{
			name:      "When indent level stays the same between lines, Then no violations are reported",
			content:   "{\n  \"a\": \"1\",\n  \"b\": \"2\"\n}\n",
			wantLines: nil,
		},
		{
			name:      "When indent level decreases, Then no violations are reported",
			content:   "{\n  \"a\": {\n    \"b\": \"c\"\n  }\n}\n",
			wantLines: nil,
		},
		{
			name:      "When a blank line appears before an indented line, Then blank lines are skipped and no violations are reported",
			content:   "{\n\n  \"key\": \"value\"\n}\n",
			wantLines: nil,
		},
		{
			name:      "When a template directive appears between an opener and the indented line, Then the directive is skipped and no violations are reported",
			content:   "{\n  \"items\": [\n    {{ range .Items }}\n    \"{{.}}\"\n    {{ end }}\n  ]\n}\n",
			wantLines: nil,
		},
		{
			name:      "When a template directive line precedes the first content line, Then it is skipped and no violations are reported",
			content:   "{{ range .Items }}\n{\n  \"key\": \"{{.}}\"\n}\n{{ end }}\n",
			wantLines: nil,
		},

		// --- Violations ---
		{
			name:      "When an indented line follows a line ending with ':', Then a violation is reported on the indented line",
			content:   "{\n  \"key\":\n    \"value\"\n}\n",
			wantLines: []int{3},
		},
		{
			name:      "When an indented line follows a line ending with a value, Then a violation is reported on the indented line",
			content:   "\"outer\"\n  \"inner\"\n",
			wantLines: []int{2},
		},
	}

	Convey("Given an IndentOpener rule", t, func() {
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
