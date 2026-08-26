package soros

import (
	"testing"
)

const wantRuleAfterBad = "one"

// TestCompileSkipsInvalidPattern checks that a rule whose pattern isn't
// a valid regex is recorded in Program.Skipped (see SkippedRule) instead
// of failing the whole Compile, and that every other rule still works.
func TestCompileSkipsInvalidPattern(t *testing.T) {
	p := mustCompile(t, "(unbalanced bad-rule\n1 one", "")
	if len(p.Skipped) != 1 {
		t.Fatalf("len(Skipped) = %d, want 1", len(p.Skipped))
	}
	if got, want := p.Run("1"), wantRuleAfterBad; got != want {
		t.Errorf("Run(1) = %q, want %q (the rule after the bad one must still work)", got, want)
	}
}

// TestCompileSkipsInvalidTemplate checks the same for a rule whose
// replacement template fails to parse (an unterminated "[").
func TestCompileSkipsInvalidTemplate(t *testing.T) {
	p := mustCompile(t, "2 two[unterminated\n1 one", "")
	if len(p.Skipped) != 1 {
		t.Fatalf("len(Skipped) = %d, want 1", len(p.Skipped))
	}
	if got, want := p.Run("1"), wantRuleAfterBad; got != want {
		t.Errorf("Run(1) = %q, want %q (the rule after the bad one must still work)", got, want)
	}
}

// TestApplyPrefixHoistsCaret exercises spec 2.8's own worked example
// directly against the unexported helper: a macro-prefixed rule whose
// own pattern starts with '^' must have the caret hoisted in front of
// the whole concatenated pattern, not left attached to the original
// fragment - no bundled ".sor" file happens to combine a macro with a
// '^'-anchored line, so this needs its own test.
func TestApplyPrefixHoistsCaret(t *testing.T) {
	const prefix = "ordinal"
	cases := []struct{ pattern, want string }{
		{"2", prefix + " 2"},
		{"^2", "^" + prefix + " 2"},
		{"", prefix},
	}
	for _, tc := range cases {
		if got := applyPrefix(prefix, tc.pattern); got != tc.want {
			t.Errorf("applyPrefix(%q, %q) = %q, want %q", prefix, tc.pattern, got, tc.want)
		}
	}
}
