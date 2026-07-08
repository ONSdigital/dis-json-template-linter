package linter

import (
	"os"
	"sort"
	"strings"

	"github.com/ONSdigital/dis-json-template-linter/config"
	"github.com/ONSdigital/dis-json-template-linter/rules"
)

// RunFile lints a single file using the provided config and returns all
// violations sorted by line then column.
func RunFile(path string, cfg config.Config) ([]rules.Violation, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")

	var all []rules.Violation
	for _, rule := range buildRules(cfg) {
		all = append(all, rule.Check(path, content, lines)...)
	}

	sort.Slice(all, func(i, j int) bool {
		if all[i].Line != all[j].Line {
			return all[i].Line < all[j].Line
		}
		return all[i].Col < all[j].Col
	})
	return all, nil
}

// buildRules returns the ordered list of rules to run based on cfg.
func buildRules(cfg config.Config) []rules.Rule {
	var rs []rules.Rule
	if cfg.CheckTrailingWhitespace {
		rs = append(rs, rules.TrailingWhitespace{})
	}
	if cfg.CheckTrailingNewline {
		rs = append(rs, rules.TrailingNewline{})
	}
	// Indentation is always checked; Size<=0 disables the size sub-check.
	rs = append(rs, rules.Indent{Size: cfg.IndentSize, Style: rules.IndentStyle(cfg.IndentStyle)})
	if cfg.CheckIndentStep {
		rs = append(rs, rules.IndentStep{Size: cfg.IndentSize})
	}
	if cfg.CheckIndentOpener {
		rs = append(rs, rules.IndentOpener{})
	}
	if cfg.CheckTemplateSyntax {
		rs = append(rs, rules.TemplateSyntax{})
	}
	if cfg.CheckObjectBrace {
		rs = append(rs, rules.ObjectBrace{})
	}
	if cfg.CheckTemplateSpacing {
		rs = append(rs, rules.TemplateSpacing{})
	}
	if cfg.CheckKeyColon {
		rs = append(rs, rules.KeyColon{})
	}
	if cfg.CheckTemplateAction {
		rs = append(rs, rules.TemplateAction{})
	}
	return rs
}
