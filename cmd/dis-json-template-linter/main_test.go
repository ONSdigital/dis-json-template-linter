package main

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestExpandArgs(t *testing.T) {
	Convey("Given a directory tree with nested template files", t, func() {
		root := t.TempDir()
		mustWriteFile(t, filepath.Join(root, "top.tmpl"))
		mustWriteFile(t, filepath.Join(root, "nested", "child.tmpl"))
		mustWriteFile(t, filepath.Join(root, "nested", "ignore.txt"))

		Convey("When expandArgs is given a recursive glob pattern", func() {
			pattern := filepath.Join(root, "**", "*.tmpl")
			files, err := expandArgs([]string{pattern})
			sort.Strings(files)

			Convey("Then it returns matching template files from all nested directories", func() {
				So(err, ShouldBeNil)
				So(files, ShouldResemble, []string{
					filepath.Join(root, "nested", "child.tmpl"),
					filepath.Join(root, "top.tmpl"),
				})
			})
		})

		Convey("When expandArgs is given an invalid recursive glob pattern", func() {
			pattern := filepath.Join(root, "**", "[")
			_, err := expandArgs([]string{pattern})

			Convey("Then it returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func mustWriteFile(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("test"), 0o600); err != nil {
		t.Fatal(err)
	}
}
