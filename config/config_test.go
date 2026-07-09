package config

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDefault(t *testing.T) {
	Convey("Given no configuration", t, func() {
		Convey("When Default is called", func() {
			cfg := Default()

			Convey("Then all expected defaults are returned", func() {
				So(cfg.IndentSize, ShouldEqual, 2)
				So(cfg.IndentStyle, ShouldEqual, IndentStyleSpace)
				So(cfg.CheckTrailingWhitespace, ShouldBeTrue)
				So(cfg.CheckTrailingNewline, ShouldBeTrue)
				So(cfg.CheckTemplateSyntax, ShouldBeTrue)
				So(cfg.CheckObjectBrace, ShouldBeTrue)
				So(cfg.CheckTemplateSpacing, ShouldBeTrue)
				So(cfg.CheckKeyColon, ShouldBeTrue)
				So(cfg.CheckTemplateAction, ShouldBeTrue)
				So(cfg.CheckIndentStep, ShouldBeTrue)
				So(cfg.CheckIndentOpener, ShouldBeTrue)
			})
		})
	})
}

func TestLoadFromFile(t *testing.T) {
	Convey("Given a user attempts to load from a specific file", t, func() {
		Convey("When the file does not exist", func() {
			_, err := LoadFromFile("/nonexistent/.dis-json-template-linter.yaml")

			Convey("Then an error is returned", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When the file contains invalid YAML", func() {
			path := writeTempConfig(t, "indent_size: [unclosed\n")
			_, err := LoadFromFile(path)

			Convey("Then an error is returned", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When the file is empty", func() {
			path := writeTempConfig(t, "")
			cfg, err := LoadFromFile(path)

			Convey("Then all defaults are returned", func() {
				So(err, ShouldBeNil)
				So(cfg, ShouldResemble, Default())
			})
		})

		Convey("When the file overrides a single field", func() {
			path := writeTempConfig(t, "indent_size: 4\n")
			cfg, err := LoadFromFile(path)

			Convey("Then that field is updated and all others remain as defaults", func() {
				So(err, ShouldBeNil)
				So(cfg.IndentSize, ShouldEqual, 4)
				So(cfg.IndentStyle, ShouldEqual, IndentStyleSpace)
				So(cfg.CheckTrailingWhitespace, ShouldBeTrue)
			})
		})

		Convey("When the file overrides multiple fields", func() {
			path := writeTempConfig(t, "indent_size: 4\nindent_style: tab\ncheck_trailing_whitespace: false\n")
			cfg, err := LoadFromFile(path)

			Convey("Then all overridden fields are updated", func() {
				So(err, ShouldBeNil)
				So(cfg.IndentSize, ShouldEqual, 4)
				So(cfg.IndentStyle, ShouldEqual, IndentStyleTab)
				So(cfg.CheckTrailingWhitespace, ShouldBeFalse)
			})
		})

		Convey("When indent_size is set to 0", func() {
			path := writeTempConfig(t, "indent_size: 0\n")
			cfg, err := LoadFromFile(path)

			Convey("Then the zero value is preserved and not replaced by the default", func() {
				So(err, ShouldBeNil)
				So(cfg.IndentSize, ShouldEqual, 0)
			})
		})
	})
}

func TestLoad(t *testing.T) {
	Convey("Given a user attempts to load the config without specifying a file", t, func() {
		Convey("When no config file exists anywhere in the directory tree", func() {
			dir := t.TempDir()
			cfg, err := Load(dir)

			Convey("Then defaults are returned", func() {
				So(err, ShouldBeNil)
				So(cfg, ShouldResemble, Default())
			})
		})

		Convey("When a config file exists in the start directory", func() {
			dir := t.TempDir()
			writeConfigFile(t, dir, "indent_size: 4\n")
			cfg, err := Load(dir)

			Convey("Then it is loaded", func() {
				So(err, ShouldBeNil)
				So(cfg.IndentSize, ShouldEqual, 4)
			})
		})

		Convey("When a config file exists in a parent directory", func() {
			parent := t.TempDir()
			child := filepath.Join(parent, "sub", "dir")
			if err := os.MkdirAll(child, 0o755); err != nil {
				t.Fatal(err)
			}
			writeConfigFile(t, parent, "indent_size: 6\n")
			cfg, err := Load(child)

			Convey("Then it is found by walking up", func() {
				So(err, ShouldBeNil)
				So(cfg.IndentSize, ShouldEqual, 6)
			})
		})

		Convey("When a config file exists in both the start directory and a parent", func() {
			parent := t.TempDir()
			child := filepath.Join(parent, "subdir")
			if err := os.MkdirAll(child, 0o755); err != nil {
				t.Fatal(err)
			}
			writeConfigFile(t, parent, "indent_size: 6\n")
			writeConfigFile(t, child, "indent_size: 2\nindent_style: tab\n")
			cfg, err := Load(child)

			Convey("Then the closer file takes precedence", func() {
				So(err, ShouldBeNil)
				So(cfg.IndentStyle, ShouldEqual, IndentStyleTab)
			})
		})
	})
}

// writeTempConfig writes content to a temporary file and returns its path.
func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	if err = f.Close(); err != nil {
		t.Fatal(err)
	}
	return f.Name()
}

// writeConfigFile writes content to configFileName inside dir.
func writeConfigFile(t *testing.T, dir, content string) {
	t.Helper()
	path := filepath.Join(dir, configFileName)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}
