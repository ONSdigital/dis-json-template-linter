package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ONSdigital/dis-json-template-linter/config"
	"github.com/ONSdigital/dis-json-template-linter/linter"
	"github.com/ONSdigital/dis-json-template-linter/rules"
)

// errViolations is a sentinel returned by RunE when lint violations are found.
// main() maps it to exit code 1 without printing cobra's "Error: ..." prefix.
var errViolations = errors.New("lint violations found")

var configPath string

var rootCmd = &cobra.Command{
	Use:   "dis-json-template-linter <file|glob> [<file|glob>...]",
	Short: "Lint JSON+Go-template (.tmpl) files for style violations",
	Long: `dis-json-template-linter checks .tmpl files that contain JSON written in
Go templating syntax, enforcing rules such as indentation, trailing
whitespace, missing trailing newlines, and Go template syntax validity.

Rules are configured via a .dis-json-template-linter.yaml file discovered by walking up the
directory tree from each linted file (or supplied explicitly with --config).`,
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	RunE:         run,
}

func main() {
	rootCmd.Flags().StringVarP(
		&configPath, "config", "c", "",
		"path to .dis-json-template-linter.yaml (default: search up from each file's directory)",
	)
	// SilenceErrors prevents cobra from printing "Error: lint violations found"
	// after we have already printed the individual violation lines.
	rootCmd.SilenceErrors = true
	if err := rootCmd.Execute(); err != nil {
		if errors.Is(err, errViolations) {
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}
}

func run(_ *cobra.Command, args []string) error {
	files, err := expandArgs(args)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no files matched the given patterns")
	}
	sort.Strings(files)

	var allViolations []rules.Violation
	for _, file := range files {
		cfg, cfgErr := resolveConfig(configPath, file)
		if cfgErr != nil {
			return fmt.Errorf("loading config for %s: %w", file, cfgErr)
		}

		violations, lintErr := linter.RunFile(file, cfg)
		if lintErr != nil {
			return fmt.Errorf("linting %s: %w", file, lintErr)
		}
		allViolations = append(allViolations, violations...)
	}

	for _, v := range allViolations {
		fmt.Printf("%s:%d:%d: %s\n", v.File, v.Line, v.Col, v.Message)
	}

	nFiles := len(files)
	nViolations := len(allViolations)

	if nViolations == 0 {
		fmt.Printf("%d %s checked, 0 violations found\n",
			nFiles, pluralise(nFiles, "file", "files"),
		)
	} else {
		nAffected := countFilesWithViolations(allViolations)
		fmt.Printf("\n%d %s checked, %d %s found across %d %s\n",
			nFiles, pluralise(nFiles, "file", "files"),
			nViolations, pluralise(nViolations, "violation", "violations"),
			nAffected, pluralise(nAffected, "file", "files"),
		)
		return errViolations
	}
	return nil
}

// countFilesWithViolations returns the number of distinct files that have at
// least one violation.
func countFilesWithViolations(violations []rules.Violation) int {
	seen := make(map[string]bool)
	for _, v := range violations {
		seen[v.File] = true
	}
	return len(seen)
}

func pluralise(n int, singular, plural string) string {
	if n == 1 {
		return singular
	}
	return plural
}

// expandArgs expands file paths and glob patterns into a deduplicated,
// sorted list of regular files (directories are silently skipped).
func expandArgs(args []string) ([]string, error) {
	seen := make(map[string]bool)
	var files []string
	for _, arg := range args {
		matches, err := globMatches(arg)
		if err != nil {
			return nil, fmt.Errorf("invalid glob %q: %w", arg, err)
		}
		if len(matches) == 0 {
			fmt.Fprintf(os.Stderr, "warning: no files matched pattern %q\n", arg)
		}
		for _, m := range matches {
			abs, err := filepath.Abs(m)
			if err != nil {
				return nil, err
			}
			info, err := os.Stat(abs)
			if err != nil {
				return nil, err
			}
			if info.IsDir() {
				continue
			}
			if !seen[abs] {
				seen[abs] = true
				files = append(files, m)
			}
		}
	}
	return files, nil
}

func globMatches(pattern string) ([]string, error) {
	if !strings.Contains(pattern, "**") {
		return filepath.Glob(pattern)
	}

	root := globWalkRoot(pattern)
	var matches []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		matched, matchErr := recursiveGlobMatch(pattern, path)
		if matchErr != nil {
			return matchErr
		}
		if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if os.IsNotExist(err) {
		return nil, nil
	}
	return matches, err
}

func globWalkRoot(pattern string) string {
	clean := filepath.Clean(pattern)
	parts := strings.Split(clean, string(filepath.Separator))
	var rootParts []string

	for index, part := range parts {
		if part == "" && index == 0 {
			rootParts = append(rootParts, part)
			continue
		}
		if hasGlobMeta(part) || part == "**" {
			break
		}
		rootParts = append(rootParts, part)
	}

	if len(rootParts) == 0 {
		return "."
	}

	if filepath.IsAbs(clean) {
		if len(rootParts) == 1 && rootParts[0] == "" {
			return string(filepath.Separator)
		}
		return filepath.Join(append([]string{string(filepath.Separator)}, rootParts[1:]...)...)
	}

	root := filepath.Join(rootParts...)
	if root == "" {
		return string(filepath.Separator)
	}
	return root
}

func recursiveGlobMatch(pattern, target string) (bool, error) {
	patternParts := splitPathParts(pattern)
	targetParts := splitPathParts(target)
	return matchPathParts(patternParts, targetParts)
}

func splitPathParts(path string) []string {
	clean := filepath.Clean(path)
	if clean == "." {
		return []string{"."}
	}
	return strings.Split(clean, string(filepath.Separator))
}

func matchPathParts(patternParts, targetParts []string) (bool, error) {
	if len(patternParts) == 0 {
		return len(targetParts) == 0, nil
	}

	if patternParts[0] == "**" {
		matched, err := matchPathParts(patternParts[1:], targetParts)
		if matched || err != nil {
			return matched, err
		}
		if len(targetParts) == 0 {
			return false, nil
		}
		return matchPathParts(patternParts, targetParts[1:])
	}

	if len(targetParts) == 0 {
		return false, nil
	}

	matched, err := filepath.Match(patternParts[0], targetParts[0])
	if err != nil {
		return false, err
	}
	if !matched {
		return false, nil
	}
	return matchPathParts(patternParts[1:], targetParts[1:])
}

func hasGlobMeta(part string) bool {
	return strings.ContainsAny(part, "*?[")
}

// resolveConfig loads config from an explicit path if provided, otherwise
// searches up the directory tree from the file being linted.
func resolveConfig(cfgPath, file string) (config.Config, error) {
	if cfgPath != "" {
		return config.LoadFromFile(cfgPath)
	}
	dir, err := filepath.Abs(filepath.Dir(file))
	if err != nil {
		return config.Default(), err
	}
	return config.Load(dir)
}
