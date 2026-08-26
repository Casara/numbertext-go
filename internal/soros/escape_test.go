package soros

import "testing"

// TestHideRevealEscapesRoundTrip exercises every special character listed
// in spec 2.6: hiding an escaped form and then revealing it must produce
// the literal character, for every one of them, not just the handful
// real-world ".sor" files happen to use.
func TestHideRevealEscapesRoundTrip(t *testing.T) {
	cases := map[string]string{
		`\\`: `\`,
		`\"`: `"`,
		`\$`: `$`,
		`\(`: `(`,
		`\)`: `)`,
		`\[`: `[`,
		`\]`: `]`,
		`\|`: `|`,
		`\#`: `#`,
		`\;`: `;`,
		`\n`: "\n",
	}
	for escaped, want := range cases {
		got := revealEscapes(hideEscapes("x" + escaped + "y"))
		if want2 := "x" + want + "y"; got != want2 {
			t.Errorf("round-trip(%q) = %q, want %q", escaped, got, want2)
		}
	}
}

// TestRevealEscapesForRegex checks the regex-context reveal maps each
// special character to a RE2-safe literal (metacharacters escaped, the
// rest passed through as-is), used when compiling a pattern.
func TestRevealEscapesForRegex(t *testing.T) {
	cases := map[string]string{
		`\\`: `\\`,
		`\"`: `"`,
		`\$`: `\$`,
		`\(`: `\(`,
		`\)`: `\)`,
		`\[`: `\[`,
		`\]`: `\]`,
		`\|`: `\|`,
		`\#`: `#`,
		`\;`: `;`,
	}
	for escaped, want := range cases {
		got := revealEscapesForRegex(hideEscapes(escaped))
		if got != want {
			t.Errorf("revealEscapesForRegex(hideEscapes(%q)) = %q, want %q", escaped, got, want)
		}
	}
}

// TestHideEscapesTrailingBackslash checks that a lone, unescaped
// backslash at the very end of the string (no following character to
// combine with) is left alone rather than panicking on an out-of-range
// read.
func TestHideEscapesTrailingBackslash(t *testing.T) {
	if got, want := hideEscapes(`x\`), `x\`; got != want {
		t.Errorf("hideEscapes(%q) = %q, want %q", `x\`, got, want)
	}
}

// TestHideEscapesLeavesBackreferencesAlone checks that "\0".."\9" (a
// backreference, not an escape sequence) is left as a literal backslash
// followed by the digit, at both ends of the digit range.
func TestHideEscapesLeavesBackreferencesAlone(t *testing.T) {
	for _, digit := range []string{"0", "9"} {
		in := `\` + digit
		if got := hideEscapes(in); got != in {
			t.Errorf("hideEscapes(%q) = %q, want %q (unchanged)", in, got, in)
		}
	}
}

// TestCompileWithEscapedSpecialChars proves the escape handling works
// end-to-end, not just at the hide/reveal-function level: a pattern that
// needs to literally match one of the special characters, and a
// template that needs to literally emit one, both via the compiled
// Program.
func TestCompileWithEscapedSpecialChars(t *testing.T) {
	p := mustCompile(t, `"\(x\)" literal parens: \(x\)`, "")
	if got, want := p.Run("(x)"), "literal parens: (x)"; got != want {
		t.Errorf("Run(%q) = %q, want %q", "(x)", got, want)
	}
}
