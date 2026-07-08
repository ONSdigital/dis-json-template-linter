package rules

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTrailingNewline(t *testing.T) {
	r := TrailingNewline{}

	tests := []struct {
		name      string
		content   string
		wantLines []int
		wantMsg   string
	}{
		{
			name:      "When the file ends with a single newline, Then no violations are reported",
			content:   "{\n  \"key\": \"value\"\n}\n",
			wantLines: nil,
		},
		{
			name:      "When the file has no trailing newline, Then a violation is reported on the last line",
			content:   "{\n  \"key\": \"value\"\n}",
			wantLines: []int{3},
			wantMsg:   "no newline at end of file",
		},
		{
			name:      "When the file ends with a trailing blank line, Then a violation is reported",
			content:   "{\n  \"key\": \"value\"\n}\n\n",
			wantLines: []int{4},
			wantMsg:   "trailing blank line at end of file",
		},
		{
			name:      "When the file is empty, Then no violations are reported",
			content:   "",
			wantLines: nil,
		},
	}

	Convey("Given a TrailingNewline rule", t, func() {
		for _, tc := range tests {
			tc := tc
			Convey(tc.name, func() {
				lines := strings.Split(tc.content, "\n")
				got := r.Check("test.tmpl", []byte(tc.content), lines)
				assertViolationLines(got, tc.wantLines)
				if tc.wantMsg != "" && len(got) > 0 {
					So(got[0].Message, ShouldEqual, tc.wantMsg)
				}
			})
		}
	})
}
