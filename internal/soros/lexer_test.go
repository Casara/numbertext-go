package soros

import "testing"

// TestSplitPatternValueQuotedValue exercises the branch of the pattern/
// value split that isn't reachable through any bundled ".sor" file: a
// quoted *value* part, used to preserve leading/trailing whitespace in
// the template text (spec 2.3.3).
func TestSplitPatternValueQuotedValue(t *testing.T) {
	pattern, value := splitPatternValue(`1 "one "`)
	if pattern != "1" {
		t.Errorf("pattern = %q, want %q", pattern, "1")
	}
	if value != "one " {
		t.Errorf("value = %q, want %q (trailing space preserved)", value, "one ")
	}
}

// TestSplitPatternValueQuotedPattern checks the quoted-pattern branch
// directly (a pattern containing a literal space, spec 2.3.3), rather
// than only through a full Compile of a real rule.
func TestSplitPatternValueQuotedPattern(t *testing.T) {
	pattern, value := splitPatternValue(`"a b" c`)
	if pattern != "a b" {
		t.Errorf("pattern = %q, want %q", pattern, "a b")
	}
	if value != "c" {
		t.Errorf("value = %q, want %q", value, "c")
	}
}

// TestSplitPatternValueUnquotedNoValue checks a bare pattern with no
// value part at all (used by macro-directive-adjacent bookkeeping, e.g.
// the "__numbertext__" directive statement).
func TestSplitPatternValueUnquotedNoValue(t *testing.T) {
	pattern, value := splitPatternValue("justapattern")
	if pattern != "justapattern" || value != "" {
		t.Errorf("splitPatternValue(%q) = (%q, %q), want (%q, %q)",
			"justapattern", pattern, value, "justapattern", "")
	}
}

// TestSplitPatternValueEmptyQuotedValue checks the shortest possible
// quoted value, `""` (an explicitly empty template, as opposed to no
// value part at all): it must be recognized as "quoted, strip the
// quotes, empty result", not fall through to being returned verbatim as
// the two-character string `""`.
func TestSplitPatternValueEmptyQuotedValue(t *testing.T) {
	pattern, value := splitPatternValue(`1 ""`)
	if pattern != "1" {
		t.Errorf("pattern = %q, want %q", pattern, "1")
	}
	if value != "" {
		t.Errorf(`value = %q, want "" (the quotes stripped, not kept literally)`, value)
	}
}

// TestSplitPatternValueUnterminatedQuotedPattern checks a quoted pattern
// missing its closing '"': everything after the opening quote becomes
// the pattern, with no value part.
func TestSplitPatternValueUnterminatedQuotedPattern(t *testing.T) {
	pattern, value := splitPatternValue(`"unterminated`)
	if pattern != "unterminated" || value != "" {
		t.Errorf(`splitPatternValue(%q) = (%q, %q), want (%q, %q)`,
			`"unterminated`, pattern, value, "unterminated", "")
	}
}

// TestStripCommentsAndFilterLangUntaggedLineSurvives checks that a line
// with no "[:lang:]" tag at all is always kept (spec 2.7.1's untagged
// default), for a requested language that doesn't match any tag on a
// neighboring line - the case that would break first if the "no tags on
// this line" condition were ever inverted.
func TestStripCommentsAndFilterLangUntaggedLineSurvives(t *testing.T) {
	source := "1 one # [:en-GB:]\n2 two"
	got := stripCommentsAndFilterLang(hideEscapes(source), "de")
	want := "\n2 two"
	if got != want {
		t.Errorf("stripCommentsAndFilterLang(...) = %q, want %q", got, want)
	}
}
