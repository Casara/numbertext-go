package soros

import (
	"fmt"
	"strings"
)

// numbertextDirective is the bare statement that marks where the default
// NUMBERTEXT left-zero-deletion rule is inserted (spec 2.5). When absent,
// the rule is inserted at the very start of the program (default mode).
const numbertextDirective = "__numbertext__"

// defaultLeftZeroRule is the hidden rule every Soros program gets in
// NUMBERTEXT mode, stripping extra leading zeros (also after an optional
// macro-prefix word, so it keeps working under "== prefix ==" sections).
// The pattern is quoted because it contains a literal space (spec 2.3.3).
const defaultLeftZeroRulePattern = `"([a-z][-a-z]* )?0+(0|[1-9]\d*)"`
const defaultLeftZeroRuleValue = `$(\1\2)`

// rule is one compiled Soros program line: a full-match regular
// expression, its parsed replacement template, and whether the original
// pattern carried an explicit leading '^' / trailing '$' boundary anchor.
type rule struct {
	re        matcher
	tmpl      []node
	begin     bool
	end       bool
	statement string // original statement text, for error messages
}

// Program is a compiled Soros program: an ordered list of rules plus the
// single recursive entry point they implicitly define (the "$" function
// of the spec).
type Program struct {
	rules   []rule
	Skipped []SkippedRule
}

// SkippedRule records a program line that Compile could not use: its
// pattern is not valid regex syntax under either engine compilePattern
// tries (see regex.go), or its replacement template doesn't parse.
// Rather than fail the whole locale over one bad rule, Compile skips it:
// it simply never matches, so every other rule keeps working normally.
type SkippedRule struct {
	Statement string
	Err       error
}

// Compile parses Soros source code into a runnable Program for the given
// language code. lang selects which "[:lang-code:]" conditional lines
// apply (spec 2.7.1); pass "" if the source has none.
//
// Compile only fails outright if the source cannot be tokenized at all;
// individual rules it cannot compile are recorded in Program.Skipped
// instead (see SkippedRule) and otherwise ignored.
func Compile(source, lang string) (*Program, error) {
	hidden := hideEscapes(source)
	codeOnly := stripCommentsAndFilterLang(hidden, lang)
	statements := splitStatements(codeOnly)
	statements = injectDefaultLeftZeroRule(statements)

	p := &Program{}
	var prefix string
	prefixActive := false

	for _, stmt := range statements {
		if m := macroDirectiveRE.FindStringSubmatch(stmt); m != nil {
			prefix = m[1]
			prefixActive = true

			continue
		}

		pattern, value := splitPatternValue(stmt)
		if prefixActive && prefix != "" {
			pattern = applyPrefix(prefix, pattern)
		}

		begin := strings.HasPrefix(pattern, "^")
		if begin {
			pattern = pattern[1:]
		}
		end := strings.HasSuffix(pattern, "$")
		if end {
			pattern = pattern[:len(pattern)-1]
		}

		body := revealEscapesForRegex(pattern)
		re, err := compilePattern("^(?:" + body + ")$")
		if err != nil {
			p.Skipped = append(p.Skipped, SkippedRule{Statement: stmt, Err: fmt.Errorf("invalid pattern: %w", err)})

			continue
		}

		tmpl, err := parseTemplate(value)
		if err != nil {
			p.Skipped = append(p.Skipped, SkippedRule{Statement: stmt, Err: fmt.Errorf("invalid replacement: %w", err)})

			continue
		}

		p.rules = append(p.rules, rule{
			re:        re,
			tmpl:      tmpl,
			begin:     begin,
			end:       end,
			statement: stmt,
		})
	}

	return p, nil
}

// applyPrefix concatenates a macro prefix in front of a rule's pattern
// (spec 2.8), hoisting a leading '^' boundary anchor in front of the
// prefix when present, and using the prefix alone when the rule's own
// pattern is empty.
func applyPrefix(prefix, pattern string) string {
	if pattern == "" {
		return prefix
	}
	if strings.HasPrefix(pattern, "^") {
		return "^" + prefix + " " + pattern[1:]
	}

	return prefix + " " + pattern
}

// injectDefaultLeftZeroRule inserts the hidden left-zero-deletion rule
// (spec 2.5) at the position of an explicit "__numbertext__" directive,
// or at the very start of the program when no such directive is present.
func injectDefaultLeftZeroRule(statements []string) []string {
	hiddenRule := defaultLeftZeroRulePattern + " " + defaultLeftZeroRuleValue

	for i, s := range statements {
		if s == numbertextDirective {
			out := make([]string, 0, len(statements)+1)
			out = append(out, statements[:i]...)
			out = append(out, hiddenRule)
			out = append(out, statements[i+1:]...)

			return out
		}
	}

	out := make([]string, 0, len(statements)+1)
	out = append(out, hiddenRule)
	out = append(out, statements...)

	return out
}
