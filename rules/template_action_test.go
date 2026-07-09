package rules

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTemplateAction(t *testing.T) {
	r := TemplateAction{}

	tests := []struct {
		name      string
		content   string
		wantLines []int
	}{
		// --- Correct: template ---
		{
			name:      "When a template action uses a dot argument with correct spacing, Then no violations are reported",
			content:   "{{ template \"header.tmpl\" . }}\n",
			wantLines: nil,
		},
		{
			name:      "When a template action uses a field argument with correct spacing, Then no violations are reported",
			content:   "{{ template \"header.tmpl\" .Data }}\n",
			wantLines: nil,
		},
		{
			name:      "When a template action uses a variable argument with correct spacing, Then no violations are reported",
			content:   "{{ template \"header.tmpl\" $ctx }}\n",
			wantLines: nil,
		},
		{
			name:      "When a template action has no argument, Then no violations are reported",
			content:   "{{ template \"header.tmpl\" }}\n",
			wantLines: nil,
		},
		{
			name:      "When a template action has a left-trim modifier with correct spacing, Then no violations are reported",
			content:   "{{- template \"header.tmpl\" . }}\n",
			wantLines: nil,
		},
		{
			name:      "When a template action has both trim modifiers with correct spacing, Then no violations are reported",
			content:   "{{- template \"header.tmpl\" . -}}\n",
			wantLines: nil,
		},

		// --- Correct: block ---
		{
			name:      "When a block action uses a dot argument with correct spacing, Then no violations are reported",
			content:   "{{ block \"content\" . }}\n{{ end }}\n",
			wantLines: nil,
		},
		{
			name:      "When a block action uses a field argument with correct spacing, Then no violations are reported",
			content:   "{{ block \"content\" .Page }}\n{{ end }}\n",
			wantLines: nil,
		},

		// --- Correct: other actions unaffected ---
		{
			name:      "When a range action is present, Then it is not checked by this rule",
			content:   "{{ range .Items }}\n{{ end }}\n",
			wantLines: nil,
		},
		{
			name:      "When a define action is present, Then it is not checked by this rule",
			content:   "{{ define \"footer\" }}\n{{ end }}\n",
			wantLines: nil,
		},
		{
			name:      "When an if action is present, Then it is not checked by this rule",
			content:   "{{ if .Show }}\n{{ end }}\n",
			wantLines: nil,
		},
		{
			name:      "When a template name contains an escaped quote, Then no violations are reported",
			content:   "{{ template \"it\\'s\" . }}\n",
			wantLines: nil,
		},

		// --- Violations: template ---
		{
			name:      "When a template dot argument has no preceding space, Then a violation is reported on line 1",
			content:   "{{ template \"header.tmpl\". }}\n",
			wantLines: []int{1},
		},
		{
			name:      "When a template field argument has no preceding space, Then a violation is reported on line 1",
			content:   "{{ template \"header.tmpl\".Data }}\n",
			wantLines: []int{1},
		},
		{
			name:      "When a template variable argument has no preceding space, Then a violation is reported on line 1",
			content:   "{{ template \"header.tmpl\"$ctx }}\n",
			wantLines: []int{1},
		},

		// --- Violations: block ---
		{
			name:      "When a block dot argument has no preceding space, Then a violation is reported on line 1",
			content:   "{{ block \"content\". }}\n{{ end }}\n",
			wantLines: []int{1},
		},

		// --- Multiple blocks on one line ---
		{
			name:      "When two template calls on one line both have missing spaces, Then two violations are reported on line 1",
			content:   "{{ template \"a\". }}{{ template \"b\". }}\n",
			wantLines: []int{1, 1},
		},
		{
			name:      "When two template calls are on one line and only the first has a missing space, Then one violation is reported on line 1",
			content:   "{{ template \"a\". }}{{ template \"b\" . }}\n",
			wantLines: []int{1},
		},
	}

	Convey("Given a TemplateAction rule", t, func() {
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
