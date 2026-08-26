package soros

import "testing"

// TestCompilePatternRE2FastPath checks that an ordinary, backreference-
// free pattern compiles through the RE2 path (returning a
// *regexp.Regexp directly, not a fallback wrapper) - the path every
// rule but a handful across all of data/*.sor takes.
func TestCompilePatternRE2FastPath(t *testing.T) {
	m, err := compilePattern(`^(?:(\d+))$`)
	if err != nil {
		t.Fatalf("compilePattern: %v", err)
	}
	if _, ok := m.(regexp2Matcher); ok {
		t.Fatal("a backreference-free pattern took the regexp2 fallback path")
	}
	if got := m.FindStringSubmatch("42"); got == nil || got[1] != "42" {
		t.Errorf("FindStringSubmatch(42) = %v, want group 1 = 42", got)
	}
}

// TestCompilePatternBackreferenceFallback checks that a pattern with a
// backreference - which RE2 refuses to compile - falls back to regexp2
// and matches correctly, including the classic case a naive "probe and
// verify against one RE2 candidate split" approach gets wrong: repeated
// text where the greedy candidate split doesn't happen to be the one
// that satisfies the backreference (see AGENTS.md's regex-engine note).
func TestCompilePatternBackreferenceFallback(t *testing.T) {
	m, err := compilePattern(`^(?:([a-z]+)\1)$`)
	if err != nil {
		t.Fatalf("compilePattern: %v", err)
	}
	got := m.FindStringSubmatch("abab")
	if got == nil {
		t.Fatal("FindStringSubmatch(abab) = nil, want a match (group1=\"ab\", repeated)")
	}
	if got[1] != "ab" {
		t.Errorf("group 1 = %q, want %q", got[1], "ab")
	}
	if got := m.FindStringSubmatch("abc"); got != nil {
		t.Errorf("FindStringSubmatch(abc) = %v, want no match", got)
	}
}

// TestCompilePatternBothEnginesFail checks that a genuinely invalid
// pattern (not just RE2-incompatible) still reports an error, so
// Program.Skipped keeps working for real syntax mistakes.
func TestCompilePatternBothEnginesFail(t *testing.T) {
	_, err := compilePattern(`^(?:(unbalanced)$`)
	if err == nil {
		t.Error("compilePattern with unbalanced parens should return an error")
	}
}
