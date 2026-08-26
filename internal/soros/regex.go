package soros

import (
	"fmt"
	"regexp"
	"time"

	"github.com/dlclark/regexp2/v2"
)

// fallbackMatchTimeout bounds a single match attempt by the regexp2
// fallback engine (see compilePattern). It only ever applies to the rare
// patterns RE2 cannot express at all (see below), never to the ordinary
// RE2 path every other rule uses, so it is a safety net against a
// pathological pattern's catastrophic backtracking, not a normal-path
// budget.
const fallbackMatchTimeout = 100 * time.Millisecond

// matcher is the only capability a compiled rule pattern needs: try to
// match the (always full-string) input, and report capture groups if it
// did. *regexp.Regexp already satisfies this directly.
type matcher interface {
	FindStringSubmatch(s string) []string
}

// compilePattern compiles a full-match pattern (already wrapped in
// "^(?:...)$ " by the caller), preferring Go's standard regexp (RE2):
// linear-time in the input length, guaranteed. RE2 deliberately excludes
// backreferences and lookaround - a small number of real-world locale
// files (e.g. hu.sor's grammatical-suffix lookup) genuinely need them,
// the same way the reference C++/Java/Python implementations' native
// (backtracking) regex engines do. For exactly those patterns - the ones
// RE2 refuses to compile at all - this falls back to regexp2 (a ported
// .NET/PCRE-style backtracking engine), bounded by
// fallbackMatchTimeout so a pathological pattern degrades to "this rule
// doesn't match" instead of hanging.
//
// If neither engine can compile the pattern, the RE2 error is returned:
// that's the more informative one for an author fixing a genuine syntax
// mistake, since it's what the pattern needed to satisfy in the common
// case.
func compilePattern(wrapped string) (matcher, error) {
	re, err := regexp.Compile(wrapped)
	if err == nil {
		return re, nil
	}

	re2, err2 := regexp2.Compile(wrapped, regexp2.RE2)
	if err2 != nil {
		return nil, fmt.Errorf("regexp: %w", err)
	}
	re2.MatchTimeout = fallbackMatchTimeout

	return regexp2Matcher{re2}, nil
}

// regexp2Matcher adapts *regexp2.Regexp to the matcher interface,
// treating a match timeout or any other match-time error as simply "no
// match" rather than propagating it - the same degrade a rule with no
// matching input at all already has, and the correct outcome for a
// backtracking budget being exhausted (see fallbackMatchTimeout).
type regexp2Matcher struct {
	re *regexp2.Regexp
}

func (m regexp2Matcher) FindStringSubmatch(s string) []string {
	match, err := m.re.FindStringMatch(s)
	if err != nil || match == nil {
		return nil
	}

	groups := match.Groups()
	out := make([]string, len(groups))
	for i, g := range groups {
		out[i] = g.String()
	}

	return out
}
