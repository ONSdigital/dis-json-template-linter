package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const configFileName = ".dis-json-template-linter.yaml"

type IndentStyle string

const (
	IndentStyleSpace IndentStyle = "space"
	IndentStyleTab   IndentStyle = "tab"
)

// Config holds the linter configuration.
type Config struct {
	IndentSize              int         `yaml:"indent_size"`
	IndentStyle             IndentStyle `yaml:"indent_style"`
	CheckTrailingWhitespace bool        `yaml:"check_trailing_whitespace"`
	CheckTrailingNewline    bool        `yaml:"check_trailing_newline"`
	CheckTemplateSyntax     bool        `yaml:"check_template_syntax"`
	CheckObjectBrace        bool        `yaml:"check_object_brace"`
	CheckTemplateSpacing    bool        `yaml:"check_template_spacing"`
	CheckKeyColon           bool        `yaml:"check_key_colon"`
	CheckTemplateAction     bool        `yaml:"check_template_action"`
	CheckIndentStep         bool        `yaml:"check_indent_step"`
	CheckIndentOpener       bool        `yaml:"check_indent_opener"`
}

// Default returns the default configuration.
func Default() Config {
	return Config{
		IndentSize:              2,
		IndentStyle:             IndentStyleSpace,
		CheckTrailingWhitespace: true,
		CheckTrailingNewline:    true,
		CheckTemplateSyntax:     true,
		CheckObjectBrace:        true,
		CheckTemplateSpacing:    true,
		CheckKeyColon:           true,
		CheckTemplateAction:     true,
		CheckIndentStep:         true,
		CheckIndentOpener:       true,
	}
}

// Load searches for a .dis-json-template-linter.yaml file starting at startDir and walking up
// the directory tree. Returns defaults if no config file is found.
func Load(startDir string) (Config, error) {
	cfg := Default()
	path, found := findConfigFile(startDir, configFileName)
	if !found {
		return cfg, nil
	}
	return loadFromPath(path, cfg)
}

// LoadFromFile loads config from an explicit file path, falling back to
// defaults for any unset fields.
func LoadFromFile(path string) (Config, error) {
	return loadFromPath(path, Default())
}

func loadFromPath(path string, base Config) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return base, err
	}
	if err := yaml.Unmarshal(data, &base); err != nil {
		return base, err
	}
	return base, nil
}

// findConfigFile walks up from startDir looking for filename, stopping at the
// filesystem root. Returns the absolute path and true when found.
func findConfigFile(startDir, filename string) (string, bool) {
	dir, err := filepath.Abs(startDir)
	if err != nil {
		return "", false
	}
	for {
		candidate := filepath.Join(dir, filename)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}
