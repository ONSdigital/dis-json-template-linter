package linter

import (
	"path/filepath"
	"runtime"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dis-json-template-linter/config"
)

func testdataPath(parts ...string) string {
	pcs := [1]uintptr{}
	runtime.Callers(1, pcs[:])
	frame, _ := runtime.CallersFrames(pcs[:]).Next()
	// frame.File is linter/linter_test.go; testdata is one level up at the repo root.
	root := filepath.Join(filepath.Dir(frame.File), "..")
	return filepath.Join(append([]string{root, "testdata"}, parts...)...)
}

func TestRunFile_Clean(t *testing.T) {
	Convey("Given a clean template file, When RunFile is called, Then there is no error and no violations", t, func() {
		violations, err := RunFile(testdataPath("good", "clean.tmpl"), config.Default())
		So(err, ShouldBeNil)
		So(violations, ShouldBeEmpty)
	})
}

func TestRunFile_TrailingWhitespace(t *testing.T) {
	Convey("Given a template file with trailing whitespace, When RunFile is called, Then trailing whitespace violations are reported", t, func() {
		violations, err := RunFile(testdataPath("bad", "trailing_whitespace.tmpl"), config.Default())
		So(err, ShouldBeNil)
		So(violations, ShouldNotBeEmpty)
	})
}

func TestRunFile_NoTrailingNewline(t *testing.T) {
	Convey("Given a template file with no trailing newline, When RunFile is called, Then a 'no newline at end of file' violation is reported", t, func() {
		violations, err := RunFile(testdataPath("bad", "no_trailing_newline.tmpl"), config.Default())
		So(err, ShouldBeNil)
		found := false
		for _, v := range violations {
			if v.Message == "no newline at end of file" {
				found = true
			}
		}
		So(found, ShouldBeTrue)
	})
}

func TestRunFile_IndentThreeSpaces(t *testing.T) {
	Convey("Given a template file indented with 3 spaces, When RunFile is called, Then indentation violations are reported", t, func() {
		violations, err := RunFile(testdataPath("bad", "indent_three_spaces.tmpl"), config.Default())
		So(err, ShouldBeNil)
		So(violations, ShouldNotBeEmpty)
	})
}

func TestRunFile_IndentTabs(t *testing.T) {
	Convey("Given a template file indented with tabs, When RunFile is called, Then indent style violations are reported", t, func() {
		violations, err := RunFile(testdataPath("bad", "indent_tabs.tmpl"), config.Default())
		So(err, ShouldBeNil)
		So(violations, ShouldNotBeEmpty)
	})
}

func TestRunFile_TemplateSyntax(t *testing.T) {
	Convey("Given a template file with invalid syntax, When RunFile is called, Then a template syntax violation is reported", t, func() {
		violations, err := RunFile(testdataPath("bad", "template_syntax.tmpl"), config.Default())
		So(err, ShouldBeNil)
		found := false
		for _, v := range violations {
			if len(v.Message) >= 8 && v.Message[:8] == "template" {
				found = true
			}
		}
		So(found, ShouldBeTrue)
	})
}
