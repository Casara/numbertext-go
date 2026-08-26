package soros

import (
	"regexp"
	"strings"
)

// langTagRE matches a "[:lang-code:]" tag inside a comment (spec 2.7.1).
var langTagRE = regexp.MustCompile(`\[:([^:\]]+):\]`)

// macroDirectiveRE matches a "== prefix-literal ==" statement (spec 2.8).
var macroDirectiveRE = regexp.MustCompile(`^==\s*(.*?)\s*==$`)

// stripCommentsAndFilterLang removes '#' comments from every physical
// line of source (comments run to the end of the physical line, spec
// 2.7) and drops any line whose comment names one or more "[:lang:]"
// tags that do not include the requested lang (spec 2.7.1). Escaped
// characters must already be hidden (see hideEscapes) so that a real
// '#' is unambiguous.
func stripCommentsAndFilterLang(hidden, lang string) string {
	lang = strings.ToLower(strings.TrimSpace(lang))
	lines := strings.Split(hidden, "\n")
	kept := make([]string, 0, len(lines))
	for _, line := range lines {
		code, comment, _ := strings.Cut(line, "#")
		if tags := langTagRE.FindAllStringSubmatch(comment, -1); len(tags) > 0 {
			matched := false
			for _, t := range tags {
				if strings.ToLower(strings.TrimSpace(t[1])) == lang {
					matched = true

					break
				}
			}
			if !matched {
				kept = append(kept, "")

				continue
			}
		}
		kept = append(kept, code)
	}

	return strings.Join(kept, "\n")
}

// splitStatements splits comment-free source into trimmed statements,
// separated by ';' or newline (spec 2.1), dropping empty ones.
func splitStatements(codeOnly string) []string {
	raw := strings.FieldsFunc(codeOnly, func(r rune) bool {
		return r == ';' || r == '\n'
	})
	statements := make([]string, 0, len(raw))
	for _, s := range raw {
		s = strings.TrimSpace(s)
		if s != "" {
			statements = append(statements, s)
		}
	}

	return statements
}

// splitPatternValue splits one statement into its regex and value parts
// (spec 2.2), honoring optional ASCII-quote wrapping of either part
// (spec 2.3.3). The statement text is expected to still carry hidden
// escape placeholders (see hideEscapes), which makes any remaining
// unescaped '"' an unambiguous quote delimiter.
func splitPatternValue(stmt string) (string, string) {
	var pattern, rest string
	switch {
	case strings.HasPrefix(stmt, `"`):
		if end := strings.IndexByte(stmt[1:], '"'); end >= 0 {
			pattern, rest = stmt[1:1+end], stmt[1+end+1:]
		} else {
			pattern = stmt[1:]
		}
	default:
		if i := strings.IndexAny(stmt, " \t"); i >= 0 {
			pattern, rest = stmt[:i], stmt[i:]
		} else {
			pattern = stmt
		}
	}

	rest = strings.TrimLeft(rest, " \t")
	if len(rest) >= 2 && strings.HasPrefix(rest, `"`) && strings.HasSuffix(rest, `"`) {
		return pattern, rest[1 : len(rest)-1]
	}

	return pattern, rest
}
