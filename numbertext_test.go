package numbertext_test

import (
	"testing"

	numbertext "github.com/casara/numbertext-go"
)

func TestCardinalEnglish(t *testing.T) {
	// Plain "en" (no country variant) intentionally omits "and" before
	// the tens/units part of a hundred (see en.sor: the en-AU/GB/IE/NZ
	// tagged rule adds "and", the untagged default rule doesn't), and
	// "0.x" is read as "point x" without a leading "zero" (en.sor has a
	// dedicated "0[.,] point" rule for exactly that case).
	cases := map[string]string{
		"0":             "zero",
		"1":             "one",
		"9":             "nine",
		"10":            "ten",
		"11":            "eleven",
		"19":            "nineteen",
		"20":            "twenty",
		"21":            "twenty-one",
		"99":            "ninety-nine",
		"100":           "one hundred",
		"101":           "one hundred one",
		"999":           "nine hundred ninety-nine",
		"1000":          "one thousand",
		"1001":          "one thousand and one",
		"12345":         "twelve thousand three hundred forty-five",
		"1000000":       "one million",
		"2000000":       "two million",
		"-5":            "negative five",
		"-123":          "negative one hundred twenty-three",
		"0.5":           "point five",
		"3.14":          "three point one four",
		"1000000000":    "one billion",
		"1000000000000": "one trillion",
	}
	for in, want := range cases {
		got, err := numbertext.Cardinal("en", in)
		if err != nil {
			t.Fatalf("Cardinal(en, %q): %v", in, err)
		}
		if got != want {
			t.Errorf("Cardinal(en, %q) = %q, want %q", in, got, want)
		}
	}
}

func TestCardinalIntEnglish(t *testing.T) {
	got, err := numbertext.CardinalInt("en", 42)
	if err != nil {
		t.Fatalf("CardinalInt: %v", err)
	}
	if want := "forty-two"; got != want {
		t.Errorf("CardinalInt(en, 42) = %q, want %q", got, want)
	}
}

func TestOrdinalEnglish(t *testing.T) {
	cases := map[string]string{
		"1":   "first",
		"2":   "second",
		"3":   "third",
		"5":   "fifth",
		"8":   "eighth",
		"9":   "ninth",
		"12":  "twelfth",
		"20":  "twentieth",
		"21":  "twenty-first",
		"100": "one hundredth",
	}
	for in, want := range cases {
		got, err := numbertext.Ordinal("en", in)
		if err != nil {
			t.Fatalf("Ordinal(en, %q): %v", in, err)
		}
		if got != want {
			t.Errorf("Ordinal(en, %q) = %q, want %q", in, got, want)
		}
	}
}

func TestOrdinalNumberEnglish(t *testing.T) {
	cases := map[string]string{
		"1":  "1st",
		"2":  "2nd",
		"3":  "3rd",
		"4":  "4th",
		"11": "11th",
		"12": "12th",
		"13": "13th",
		"21": "21st",
		"22": "22nd",
		"23": "23rd",
	}
	for in, want := range cases {
		got, err := numbertext.OrdinalNumber("en", in)
		if err != nil {
			t.Fatalf("OrdinalNumber(en, %q): %v", in, err)
		}
		if got != want {
			t.Errorf("OrdinalNumber(en, %q) = %q, want %q", in, got, want)
		}
	}
}

func TestYearEnglish(t *testing.T) {
	// en.sor reads a "19xx" year as two two-digit groups ("nineteen
	// ninety-nine"), and strips the cardinal "and" for years outside
	// that range via its "year-remove-and" section ("two thousand one").
	cases := map[string]string{
		"1999": "nineteen ninety-nine",
		"2000": "two thousand",
		"2001": "two thousand one",
	}
	for in, want := range cases {
		got, err := numbertext.Year("en", in)
		if err != nil {
			t.Fatalf("Year(en, %q): %v", in, err)
		}
		if got != want {
			t.Errorf("Year(en, %q) = %q, want %q", in, got, want)
		}
	}
}

func TestCurrencyEnglish(t *testing.T) {
	cases := []struct{ code, amount, want string }{
		{"USD", "1", "one U.S. dollar"},
		{"USD", "2.5", "two U.S. dollars and fifty cents"},
		{"EUR", "1", "one euro"},
		{"EUR", "2", "two euro"},
	}
	for _, tc := range cases {
		got, err := numbertext.Currency("en", tc.code, tc.amount)
		if err != nil {
			t.Fatalf("Currency(en, %q, %q): %v", tc.code, tc.amount, err)
		}
		if got != tc.want {
			t.Errorf("Currency(en, %q, %q) = %q, want %q", tc.code, tc.amount, got, tc.want)
		}
	}
}

func TestMoneyEnglish(t *testing.T) {
	// The "money" (MONEYTEXT-style) section mirrors classic check-writing
	// wording: amount and fraction first, unit name last.
	got, err := numbertext.Money("en", "USD", "2.5")
	if err != nil {
		t.Fatalf("Money: %v", err)
	}
	if want := "two and 50/100 U.S. dollars"; got != want {
		t.Errorf("Money(en, USD, 2.5) = %q, want %q", got, want)
	}
}

// TestRegionalVariant checks that a regional code with no ".sor" file of
// its own (e.g. "en-GB", "pt-BR") falls back to its base language's file
// (data/en.sor, data/pt.sor) while still activating that region's
// "[:lang-code:]"-tagged lines (spec 2.7.1), which is how en.sor and
// pt.sor encode regional wording differences within a single file.
func TestRegionalVariant(t *testing.T) {
	if got, want := mustCardinal(t, "en-GB", "101"), "one hundred and one"; got != want {
		t.Errorf("Cardinal(en-GB, 101) = %q, want %q", got, want)
	}
	if got, want := mustCardinal(t, "en", "101"), "one hundred one"; got != want {
		t.Errorf("Cardinal(en, 101) = %q, want %q", got, want)
	}
	if got, want := mustCardinal(t, "pt-BR", "16"), "dezesseis"; got != want {
		t.Errorf("Cardinal(pt-BR, 16) = %q, want %q", got, want)
	}
	if got, want := mustCardinal(t, "pt", "16"), "dezasseis"; got != want {
		t.Errorf("Cardinal(pt, 16) = %q, want %q", got, want)
	}

	_, err := numbertext.Cardinal("xx-YY", "1")
	if err == nil {
		t.Error("Cardinal with an unknown base language should return an error")
	}
}

func mustCardinal(t *testing.T, lang, n string) string {
	t.Helper()
	got, err := numbertext.Cardinal(lang, n)
	if err != nil {
		t.Fatalf("Cardinal(%s, %s): %v", lang, n, err)
	}

	return got
}

func TestHelpEnglish(t *testing.T) {
	got, err := numbertext.Help("en")
	if err != nil {
		t.Fatalf("Help: %v", err)
	}
	if got == "" {
		t.Error("Help(en) = \"\", want non-empty usage text")
	}
}

// TestGenderPortuguese is a smoke test for locale-specific sections
// reached through Convert, e.g. pt.sor's "feminine"/"masculine" gender
// variants (not a universal API since section names are not
// standardized across locale files, see numbertext.go's Convert doc).
func TestGenderPortuguese(t *testing.T) {
	cardinal, err := numbertext.Cardinal("pt", "2")
	if err != nil {
		t.Fatalf("Cardinal(pt, 2): %v", err)
	}
	if want := "dois"; cardinal != want {
		t.Errorf("Cardinal(pt, 2) = %q, want %q", cardinal, want)
	}

	feminine, err := numbertext.Convert("pt", "feminine", "2")
	if err != nil {
		t.Fatalf("Convert(pt, feminine, 2): %v", err)
	}
	if want := "duas"; feminine != want {
		t.Errorf("Convert(pt, feminine, 2) = %q, want %q", feminine, want)
	}
}

// TestHungarianSuffixDatabaseLookup exercises hu.sor's grammatical-case
// suffix lookup, the one real-world use of a backreference inside a
// *matching* pattern across all of data/*.sor (see internal/soros's
// regex.go): a "database" of case->suffix pairs encoded as one string,
// searched via ".*:\1:(-[^-:]*).*" where \1 correlates the requested
// case name (captured earlier, from a macro prefix) with its entry.
// RE2 can't compile that pattern at all; this only works because
// Compile falls back to regexp2 for it. Expected values come straight
// from hu.sor's own "help" section example ("-szor 15 -> 15-ször").
func TestHungarianSuffixDatabaseLookup(t *testing.T) {
	got, err := numbertext.Convert("hu", "-szor", "15")
	if err != nil {
		t.Fatalf("Convert(hu, -szor, 15): %v", err)
	}
	if want := "15-ször"; got != want {
		t.Errorf("Convert(hu, -szor, 15) = %q, want %q", got, want)
	}
}

// TestBulgarianOrdinalThousands exercises the bg.sor rule that used to
// be the one entry in Program.Skipped for a template-grammar reason
// rather than a regex-engine one (see AGENTS.md's "Resolved: bg.sor's
// one non-nesting rule was a data typo, not a parser gap"): upstream's
// "N thousandth [and REMAINDERth]" rule had a bracket/paren transposed
// (`\2)])` instead of `\2))]`), patched locally in data/bg.sor pending
// an upstream fix - see data/SOURCE.md's "Local patches" section.
func TestBulgarianOrdinalThousands(t *testing.T) {
	cases := map[string]string{
		"1023": "хиляда двадесет и трети",
		"2023": "две хиляди двадесет и трети",
		"5023": "пет хиляди двадесет и трети",
	}
	for in, want := range cases {
		got, err := numbertext.Ordinal("bg", in)
		if err != nil {
			t.Fatalf("Ordinal(bg, %q): %v", in, err)
		}
		if got != want {
			t.Errorf("Ordinal(bg, %q) = %q, want %q", in, got, want)
		}
	}
}

func TestLanguagesIncludesBundledLocales(t *testing.T) {
	langs := numbertext.Languages()
	want := map[string]bool{"en": false, "pt": false, "de": false}
	for _, l := range langs {
		if _, ok := want[l]; ok {
			want[l] = true
		}
	}
	for l, found := range want {
		if !found {
			t.Errorf("Languages() missing %q", l)
		}
	}
}

func TestUnknownLanguage(t *testing.T) {
	_, err := numbertext.Cardinal("xx-not-a-lang", "1")
	if err == nil {
		t.Error("Cardinal with an unknown language code should return an error")
	}
}

func TestRegisterLocale(t *testing.T) {
	source := "1 uno\n2 dos\n3 tres\n"
	err := numbertext.RegisterLocale("test-es-mini", source)
	if err != nil {
		t.Fatalf("RegisterLocale: %v", err)
	}
	got, err := numbertext.Cardinal("test-es-mini", "2")
	if err != nil {
		t.Fatalf("Cardinal after RegisterLocale: %v", err)
	}
	if want := "dos"; got != want {
		t.Errorf("Cardinal(test-es-mini, 2) = %q, want %q", got, want)
	}

	found := false
	for _, l := range numbertext.Languages() {
		if l == "test-es-mini" {
			found = true
		}
	}
	if !found {
		t.Error("Languages() did not include the just-registered locale")
	}
}
