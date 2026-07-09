package rules

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestKeyColon(t *testing.T) {
	r := KeyColon{}

	tests := []struct {
		name      string
		content   string
		wantLines []int
	}{
		// --- Correct ---
		{
			name:      "When key-colon spacing is correct, Then no violations are reported",
			content:   "{\n  \"key\": \"value\"\n}\n",
			wantLines: nil,
		},
		{
			name:      "When an object value has correct spacing, Then no violations are reported",
			content:   "{\n  \"match\": {\n    \"field\": \"value\"\n  }\n}\n",
			wantLines: nil,
		},
		{
			name:      "When an array value has correct spacing, Then no violations are reported",
			content:   "{\n  \"items\": [\n    \"a\"\n  ]\n}\n",
			wantLines: nil,
		},
		{
			name:      "When a colon appears inside a URL string value, Then it is not flagged",
			content:   "{\n  \"url\": \"http://example.com\"\n}\n",
			wantLines: nil,
		},
		{
			name:      "When a template expression is used as a value, Then no violations are reported",
			content:   "{\n  \"field\": \"{{ .Value }}\"\n}\n",
			wantLines: nil,
		},
		{
			name:      "When a colon appears inside a template action, Then it is not flagged",
			content:   "{{ range $k, $v := .Map }}\n{{ end }}\n",
			wantLines: nil,
		},
		{
			name:      "When a line is a template directive, Then it is exempt from the rule",
			content:   "{{range $i,$e := .Items}}\n{{end}}\n",
			wantLines: nil,
		},

		// --- Violations: space before ':' ---
		{
			name:      "When there is one space before the colon, Then a violation is reported on line 2",
			content:   "{\n  \"key\" : \"value\"\n}\n",
			wantLines: []int{2},
		},
		{
			name:      "When there are two spaces before the colon, Then a violation is reported on line 2",
			content:   "{\n  \"key\"  : \"value\"\n}\n",
			wantLines: []int{2},
		},

		// --- Violations: wrong space after ':' ---
		{
			name:      "When there is no space after the colon, Then a violation is reported on line 2",
			content:   "{\n  \"key\":\"value\"\n}\n",
			wantLines: []int{2},
		},
		{
			name:      "When there is no space after the colon before an object brace, Then a violation is reported on line 2",
			content:   "{\n  \"match\":{\n    \"field\": \"value\"\n  }\n}\n",
			wantLines: []int{2},
		},
		{
			name:      "When there are two spaces after the colon, Then a violation is reported on line 2",
			content:   "{\n  \"key\":  \"value\"\n}\n",
			wantLines: []int{2},
		},

		// --- Both violations on same line ---
		{
			name:      "When there is a space before and no space after the colon, Then two violations are reported on line 2",
			content:   "{\n  \"key\" :\"value\"\n}\n",
			wantLines: []int{2, 2},
		},

		// --- Multiple keys on same line ---
		{
			name:      "When two keys on one line both have wrong spacing, Then two violations are reported on line 1",
			content:   "{ \"a\":\"1\", \"b\":\"2\" }\n",
			wantLines: []int{1, 1},
		},
		{
			name:      "When two keys on one line both have correct spacing, Then no violations are reported",
			content:   "{ \"a\": \"1\", \"b\": \"2\" }\n",
			wantLines: nil,
		},
	}

	Convey("Given a KeyColon rule", t, func() {
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
