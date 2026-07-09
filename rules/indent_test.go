package rules

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIndent_Style(t *testing.T) {
	r := Indent{Size: 2, Style: IndentStyleSpace}

	tests := []struct {
		name      string
		content   string
		wantLines []int
	}{
		{
			name:      "When indentation uses 2 spaces (matching the configured size), Then no violations are reported",
			content:   "{\n  \"key\": \"value\"\n}\n",
			wantLines: nil,
		},
		{
			name:      "When indentation uses tabs but the style is space, Then a violation is reported on line 2",
			content:   "{\n\t\"key\": \"value\"\n}\n",
			wantLines: []int{2},
		},
		{
			name:      "When a line is a template directive, Then it is exempt from indentation checks",
			content:   "{{range $i,$e := .Items}}\n{{if $i}},{{end}}\n{\n  \"field\": \"{{.}}\"\n}\n{{end}}\n",
			wantLines: nil,
		},
		{
			name:      "When indentation uses 3 spaces (not a multiple of 2), Then a violation is reported on line 2",
			content:   "{\n   \"key\": \"value\"\n}\n",
			wantLines: []int{2},
		},
		{
			name:      "When indentation uses 4 spaces (a multiple of 2), Then no violations are reported",
			content:   "{\n    \"key\": \"value\"\n}\n",
			wantLines: nil,
		},
		{
			name:      "When a line is blank, Then it is exempt from indentation checks",
			content:   "{\n\n  \"key\": \"value\"\n}\n",
			wantLines: nil,
		},
	}

	Convey("Given an Indent rule with 2-space style", t, func() {
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

func TestIndent_TabStyle(t *testing.T) {
	r := Indent{Size: 1, Style: IndentStyleTab}

	tests := []struct {
		name      string
		content   string
		wantLines []int
	}{
		{
			name:      "When indentation uses a single tab (matching the configured style), Then no violations are reported",
			content:   "{\n\t\"key\": \"value\"\n}\n",
			wantLines: nil,
		},
		{
			name:      "When indentation uses spaces but the style is tab, Then a violation is reported on line 2",
			content:   "{\n  \"key\": \"value\"\n}\n",
			wantLines: []int{2},
		},
	}

	Convey("Given an Indent rule with tab style", t, func() {
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

func TestIsTemplateDirectiveLine(t *testing.T) {
	// Test exemption via the Indent rule: directive lines must produce no
	// violations even when they appear at column 0.
	r := Indent{Size: 2, Style: IndentStyleSpace}

	directives := []string{
		"{{range $i,$e := .Topic}}\n",
		"{{if $i}},{{end}}\n",
		"{{end}}\n",
		"{{/* topic query */}}\n",
		"{{ template \"canonicalFilters.tmpl\" . }}\n",
		"  {{end}}\n", // leading spaces + directive - still exempt
	}

	jsonLines := []string{
		"   \"key\": \"value\"\n",
	}

	Convey("Given an Indent rule with 2-space style", t, func() {
		Convey("When a line is a template directive, Then it is exempt and no violations are reported", func() {
			for _, d := range directives {
				d := d
				Convey("directive: "+d, func() {
					lines := strings.Split(d, "\n")
					got := r.Check("test.tmpl", []byte(d), lines)
					So(got, ShouldBeEmpty)
				})
			}
		})

		Convey("When a line is a JSON line with 3-space indent, Then a violation is reported", func() {
			for _, j := range jsonLines {
				j := j
				Convey("json: "+j, func() {
					lines := strings.Split(j, "\n")
					got := r.Check("test.tmpl", []byte(j), lines)
					So(got, ShouldNotBeEmpty)
				})
			}
		})
	})
}
