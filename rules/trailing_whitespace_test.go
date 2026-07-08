package rules

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTrailingWhitespace(t *testing.T) {
	r := TrailingWhitespace{}

	tests := []struct {
		name      string
		content   string
		wantLines []int
	}{
		{
			name:      "When the file has no trailing whitespace, Then no violations are reported",
			content:   "{\n  \"key\": \"value\"\n}\n",
			wantLines: nil,
		},
		{
			name:      "When line 2 has a trailing space, Then a violation is reported on line 2",
			content:   "{\n  \"key\": \"value\" \n}\n",
			wantLines: []int{2},
		},
		{
			name:      "When a line has a trailing tab, Then a violation is reported",
			content:   "line one\t\nline two\n",
			wantLines: []int{1},
		},
		{
			name:      "When multiple lines have trailing whitespace, Then a violation is reported for each",
			content:   "line one \nline two  \nline three\n",
			wantLines: []int{1, 2},
		},
		{
			name:      "When the file ends with a newline, Then the final empty line is not flagged",
			content:   "line\n",
			wantLines: nil,
		},
	}

	Convey("Given a TrailingWhitespace rule", t, func() {
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

// assertViolationLines checks that the violation line numbers exactly match wantLines.
func assertViolationLines(got []Violation, wantLines []int) {
	So(len(got), ShouldEqual, len(wantLines))
	for i, v := range got {
		So(v.Line, ShouldEqual, wantLines[i])
	}
}
