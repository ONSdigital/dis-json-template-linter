package rules

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTemplateSpacing(t *testing.T) {
	r := TemplateSpacing{}

	tests := []struct {
		name      string
		content   string
		wantLines []int
	}{
		// --- Correct forms (no violations) ---
		{
			name:      "When a template action has correct spaces on both sides, Then no violations are reported",
			content:   "{{ range .Items }}\n{{ end }}\n",
			wantLines: nil,
		},
		{
			name:      "When a left-trim modifier is followed by a space, Then no violations are reported",
			content:   "{{- range .Items }}\n{{ end }}\n",
			wantLines: nil,
		},
		{
			name:      "When a right-trim modifier is preceded by a space, Then no violations are reported",
			content:   "{{ range .Items -}}\n{{ end }}\n",
			wantLines: nil,
		},
		{
			name:      "When both trim modifiers are used with spaces, Then no violations are reported",
			content:   "{{- range .Items -}}\n{{- end -}}\n",
			wantLines: nil,
		},
		{
			name:      "When the action is a comment block, Then it is exempt from spacing checks",
			content:   "{{/* topic query */}}\n",
			wantLines: nil,
		},
		{
			name:      "When a template action includes a dot argument, Then no violations are reported",
			content:   "{{ template \"foo.tmpl\" . }}\n",
			wantLines: nil,
		},
		{
			name:      "When a string value contains }}, Then it is not confused with a closing delimiter",
			content:   "{{ printf \"val}}\" }}\n",
			wantLines: nil,
		},

		// --- Violations ---
		{
			name:      "When there is no space after {{ on both blocks, Then four violations are reported",
			content:   "{{range .Items}}\n{{end}}\n",
			wantLines: []int{1, 1, 2, 2}, // opening + closing for each block
		},
		{
			name:      "When a single block has no space after {{ and before }}, Then two violations are reported on line 1",
			content:   "{{.Value}}\n",
			wantLines: []int{1, 1},
		},
		{
			name:      "When a left-trim modifier is not followed by a space, Then a violation is reported on line 1",
			content:   "{{-range .Items }}\n{{ end }}\n",
			wantLines: []int{1},
		},
		{
			name:      "When there is no space before }}, Then a violation is reported on line 1",
			content:   "{{ range .Items}}\n{{ end }}\n",
			wantLines: []int{1},
		},
		{
			name:      "When there is no space before a right-trim modifier, Then a violation is reported on line 1",
			content:   "{{ range .Items-}}\n{{ end }}\n",
			wantLines: []int{1},
		},
		{
			name:      "When multiple blocks on one line all have missing spaces, Then four violations are reported on line 1",
			content:   "{{if $i}},{{end}}\n",
			wantLines: []int{1, 1, 1, 1}, // {{if open+close, {{end open+close
		},
	}

	Convey("Given a TemplateSpacing rule", t, func() {
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
