package rules

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIndentStep(t *testing.T) {
	r := IndentStep{Size: 2}

	tests := []struct {
		name      string
		content   string
		wantLines []int
	}{
		// --- Correct ---
		{
			name:      "When indentation increases by one level, then no violations are reported",
			content:   "{\n  \"key\": \"value\"\n}\n",
			wantLines: nil,
		},
		{
			name:      "When content is nested two levels deep, then no violations are reported",
			content:   "{\n  \"a\": {\n    \"b\": \"c\"\n  }\n}\n",
			wantLines: nil,
		},
		{
			name:      "When indentation does not change between lines, then no violations are reported",
			content:   "{\n  \"a\": \"1\",\n  \"b\": \"2\"\n}\n",
			wantLines: nil,
		},
		{
			name:      "When blank lines appear between content lines, then blank lines are skipped and no violations are reported",
			content:   "{\n\n  \"key\": \"value\"\n}\n",
			wantLines: nil,
		},
		{
			name:      "When template directive lines appear, then they are skipped and no violations are reported",
			content:   "{{ range .Items }}\n{\n  \"key\": \"{{ . }}\"\n}\n{{ end }}\n",
			wantLines: nil,
		},
		{
			name:      "When a template directive appears between content lines, then the directive is skipped and no violations are reported",
			content:   "  \"a\": \"1\",\n{{ if .Show }}\n  \"b\": \"2\"\n{{ end }}\n",
			wantLines: nil,
		},
		{
			name:      "When the rule size is zero, then the rule is disabled and no violations are reported",
			content:   "{\n        \"key\": \"value\"\n}\n",
			wantLines: nil,
		},

		// --- Violations: over-indented ---
		{
			name: "When indentation jumps by +4 (two levels at once), then a violation is reported on line 3",
			// line 3 jumps from indent 2 to indent 6 - step of +4
			content:   "{\n  \"a\": {\n      \"b\": \"c\"\n    }\n  }\n}\n",
			wantLines: []int{3},
		},
		{
			name: "When indentation jumps by +6 (three levels at once), then a violation is reported on line 3",
			// line 3 jumps from indent 2 to indent 8 - step of +6
			content:   "{\n  \"a\": {\n        \"b\": \"c\"\n      }\n    }\n  }\n}\n",
			wantLines: []int{3},
		},

		// --- Violations: under-indented (de-indented by wrong amount) ---
		{
			name:      "When indentation de-dents by 4 (two levels at once), then a violation is reported on line 5",
			content:   "{\n  \"a\": {\n    \"b\": {\n      \"c\": \"d\"\n  }\n}\n",
			wantLines: []int{5},
		},
		{
			name:      "When indentation jumps from deeply nested to zero in one step, then a violation is reported on line 5",
			content:   "{\n  \"a\": {\n    \"b\": {\n      \"c\": \"d\"\n}\n",
			wantLines: []int{5},
		},

		// --- Multiple violations ---
		{
			name: "When there are two bad indent steps, then violations are reported on lines 3 and 5",
			// line 3: +4 (from 2 to 6), line 5: +4 (from 6 to 10)
			content:   "{\n  \"a\": {\n      \"b\": \"c\",\n      \"d\": {\n          \"e\": \"f\"\n        }\n      }\n    }\n  }\n}\n",
			wantLines: []int{3, 5},
		},
	}

	Convey("Given an IndentStep rule with size 2", t, func() {
		for _, tc := range tests {
			tc := tc
			Convey(tc.name, func() {
				// Override size for the "size zero disables rule" test.
				rule := r
				if tc.name == "When the rule size is zero, then the rule is disabled and no violations are reported" {
					rule = IndentStep{Size: 0}
				}
				lines := strings.Split(tc.content, "\n")
				got := rule.Check("test.tmpl", []byte(tc.content), lines)
				assertViolationLines(got, tc.wantLines)
			})
		}
	})
}
