package soros

import "testing"

// mustCompile is a small test helper.
func mustCompile(t *testing.T, source, lang string) *Program {
	t.Helper()
	p, err := Compile(source, lang)
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}

	return p
}

// TestReverseString is the spec's Example 1 (section 3): a program with
// no NUMBERTEXT-mode concerns, since it doesn't use digits at all — but
// the hidden left-zero rule is always present and must simply never
// match non-numeric input.
func TestReverseString(t *testing.T) {
	p := mustCompile(t, `(.*)(.) \2$1`, "")
	got := p.Run("hello")
	if want := "olleh"; got != want {
		t.Errorf("Run(%q) = %q, want %q", "hello", got, want)
	}
}

// TestThousandSeparator is the spec's Example 2.
func TestThousandSeparator(t *testing.T) {
	source := `
(\d+)(\d{3}) $1,\2
(\d+) \1
`
	p := mustCompile(t, source, "")
	cases := map[string]string{
		"1":         "1",
		"999":       "999",
		"1000":      "1,000",
		"1234567":   "1,234,567",
		"123456789": "123,456,789",
	}
	for in, want := range cases {
		if got := p.Run(in); got != want {
			t.Errorf("Run(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestDevanagariDigits is the spec's Example 3.
func TestDevanagariDigits(t *testing.T) {
	source := `
(\d*)0 $1०
(\d*)1 $1१
(\d*)2 $1२
(\d*)3 $1३
(\d*)4 $1४
(\d*)5 $1५
(\d*)6 $1६
(\d*)7 $1७
(\d*)8 $1८
(\d*)9 $1९
`
	p := mustCompile(t, source, "")
	if got, want := p.Run("0"), "०"; got != want {
		t.Errorf("Run(0) = %q, want %q", got, want)
	}
	if got, want := p.Run("29"), "२९"; got != want {
		t.Errorf("Run(29) = %q, want %q", got, want)
	}
}

// TestEnglishNumberNameExample is the spec's Example 4, the small English
// number-name program shown inline in the spec text (distinct from the
// full data/en.sor shipped with this repo).
func TestEnglishNumberNameExample(t *testing.T) {
	source := `
__numbertext__
^0 zero; 1 one; 2 two; 3 three; 4 four; 5 five;
6 six; 7 seven; 8 eight; 9 nine; 10 ten; 11 eleven; 12 twelve;
13 thirteen; 15 fifteen; 18 eighteen; 1(\d) $1teen
2(\d) twenty[-$1]; 3(\d) thirty[-$1]; 4(\d) forty[-$1]
5(\d) fifty[-$1]; 8(\d) eighty[-$1]
(\d)(\d) $1ty[-$2]
(\d)(\d\d) $1 hundred[ and $2]
`
	p := mustCompile(t, source, "")
	cases := map[string]string{
		"0":   "zero",
		"1":   "one",
		"9":   "nine",
		"10":  "ten",
		"13":  "thirteen",
		"20":  "twenty",
		"21":  "twenty-one",
		"42":  "forty-two",
		"99":  "ninety-nine",
		"100": "one hundred",
		"101": "one hundred and one",
		"007": "seven",
		"00":  "zero",
	}
	for in, want := range cases {
		if got := p.Run(in); got != want {
			t.Errorf("Run(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestConditionalBracketEmpty verifies that a "[...]" block disappears
// entirely (bracket text and all) when the call inside it is empty, and
// is included (with the call substituted) when it is not (spec 2.4.5).
func TestConditionalBracketEmpty(t *testing.T) {
	source := `
1 one
2 two
(\d)(\d) $1[-$2]
0
`
	p := mustCompile(t, source, "")
	if got, want := p.Run("10"), "one"; got != want {
		t.Errorf("Run(10) = %q, want %q", got, want)
	}
	if got, want := p.Run("12"), "one-two"; got != want {
		t.Errorf("Run(12) = %q, want %q", got, want)
	}
}

// TestPipeBoundaryModifier checks that a '|' immediately before/after a
// call forces its leading/trailing boundary eligibility regardless of
// its position in the template (spec 2.4.4), the mechanism the real
// locale files rely on for negative numbers, e.g. "negative |$1".
func TestPipeBoundaryModifier(t *testing.T) {
	source := "^0 zero\n[-−](\\d+) negative |$1"
	p := mustCompile(t, source, "")
	if got, want := p.Run("-0"), "negative zero"; got != want {
		t.Errorf("Run(-0) = %q, want %q", got, want)
	}
}

// TestTrailingDollarAnchor checks that a pattern ending in a literal '$'
// (spec 2.3.2's trailing-boundary anchor) only matches when the current
// call is itself at the trailing boundary of the template it's embedded
// in - the same digit processed at a non-trailing vs. a trailing
// position of "$1 $2" must take two different rules.
func TestTrailingDollarAnchor(t *testing.T) {
	source := `
1 one
2$ two-at-end
2 two
(\d)(\d) $1 $2
`
	p := mustCompile(t, source, "")
	if got, want := p.Run("22"), "two two-at-end"; got != want {
		t.Errorf("Run(22) = %q, want %q", got, want)
	}
}

// TestGroupOutOfRange checks that referencing a capture group beyond
// what the matched pattern actually captured (\N for an N the rule's own
// pattern has no group for) returns "" rather than panicking.
func TestGroupOutOfRange(t *testing.T) {
	p := mustCompile(t, `(\d) \2`, "")
	if got, want := p.Run("5"), ""; got != want {
		t.Errorf("Run(5) = %q, want %q", got, want)
	}
}

// TestOptionalBracketBoundaryPropagation checks that a call inside a
// "[...]" block only inherits the enclosing begin/end boundary when
// it's actually the first/last child of the bracket *and* the bracket
// itself sits at that boundary of its own template (spec 2.4.4's
// "brackets are inline, not a new template" rule for renderOptional) -
// found via mutation testing: no other test in the module distinguishes
// this from "any child inherits the boundary regardless of position",
// which would wrongly let a mid-bracket call trigger a '^'/'$'-anchored
// rule it isn't entitled to.
func TestOptionalBracketBoundaryPropagation(t *testing.T) {
	source := `
a alpha
^z here-begin
z$ here-end
z neither
(a)(z) $1[x$2]
`
	p := mustCompile(t, source, "")
	// The outer call is "az", which doesn't itself match "^z"/"z$" (a
	// full-match requirement), so begin/end boundary status for "z" can
	// only come from how (a)(z)'s template propagates it through the
	// bracket - not from "z" accidentally being the whole input too.
	//
	// $2 ("z") is the bracket's *second* child (i=1, not i=0), so it
	// must not inherit begin=true even though the bracket is the
	// template's own last (and only) node - "^z here-begin" must not
	// fire. It *is* the bracket's last child, so it must inherit
	// end=true - "z$ here-end" must fire, ahead of unanchored "z neither".
	if got, want := p.Run("az"), "alphaxhere-end"; got != want {
		t.Errorf("Run(az) = %q, want %q", got, want)
	}
}

// TestGroupZero checks that "\0" (a raw backreference to the whole
// match, used by e.g. pt.sor's and hu.sor's "help" sections as
// "$(\0 feminine)"/"$(\0 ordinal)") is treated as a valid group index,
// not rejected the same way a negative index would be.
func TestGroupZero(t *testing.T) {
	p := mustCompile(t, `(\d)(\d) \0`, "")
	if got, want := p.Run("12"), "12"; got != want {
		t.Errorf("Run(12) = %q, want %q (\\0 = the whole match)", got, want)
	}
}

// TestMacroPrefix exercises "== name ==" sections (spec 2.8): calling a
// prefixed rule set from the outside by constructing "<prefix> <arg>",
// and confirms cardinal-to-ordinal ending rewriting works through it.
func TestMacroPrefix(t *testing.T) {
	source := `
1 one
2 two
3 three

== ordinal ==
(\d+) $(ordinal |$1)
(.*)one \1first
(.*)two \1second
(.*)three \1third
`
	p := mustCompile(t, source, "")
	if got, want := p.Run("ordinal 1"), "first"; got != want {
		t.Errorf("Run(%q) = %q, want %q", "ordinal 1", got, want)
	}
	if got, want := p.Run("ordinal 2"), "second"; got != want {
		t.Errorf("Run(%q) = %q, want %q", "ordinal 2", got, want)
	}
	if got, want := p.Run("ordinal 3"), "third"; got != want {
		t.Errorf("Run(%q) = %q, want %q", "ordinal 3", got, want)
	}
}

// TestEmptyMacroDisablesPrefix checks "== ==" turns macro expansion off
// for subsequent lines (spec 2.8).
func TestEmptyMacroDisablesPrefix(t *testing.T) {
	source := `
== ordinal ==
1 first
== ==
2 two
`
	p := mustCompile(t, source, "")
	if got, want := p.Run("ordinal 1"), "first"; got != want {
		t.Errorf("Run(%q) = %q, want %q", "ordinal 1", got, want)
	}
	if got, want := p.Run("2"), "two"; got != want {
		t.Errorf("Run(%q) = %q, want %q", "2", got, want)
	}
}

// TestLanguageConditionalLine checks "[:lang-code:]" comment tags (spec
// 2.7.1): a line only applies for the requested language, otherwise the
// next (untagged, default) rule takes over.
func TestLanguageConditionalLine(t *testing.T) {
	source := `
(\d)(\d\d) \1 hundred and \2 # [:en-GB:]
(\d)(\d\d) \1 hundred \2
`
	if got, want := mustCompile(t, source, "en-GB").Run("123"), "1 hundred and 23"; got != want {
		t.Errorf("en-GB: Run(123) = %q, want %q", got, want)
	}
	if got, want := mustCompile(t, source, "en-US").Run("123"), "1 hundred 23"; got != want {
		t.Errorf("en-US: Run(123) = %q, want %q", got, want)
	}
}
