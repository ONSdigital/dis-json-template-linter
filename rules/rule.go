package rules

// Violation represents a single lint rule failure in a file.
type Violation struct {
	File    string
	Line    int // 1-based
	Col     int // 1-based
	Message string
}

// Rule is the interface implemented by every lint rule.
// content is the raw file bytes; lines is content split on "\n".
type Rule interface {
	Check(file string, content []byte, lines []string) []Violation
}
