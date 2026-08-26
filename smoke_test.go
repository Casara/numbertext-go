package numbertext_test

import (
	"testing"

	numbertext "github.com/casara/numbertext-go"
)

// TestSmokeAllLanguages compiles and exercises every bundled locale with
// a handful of representative inputs. It does not assert exact wording
// (that would require a native speaker per language); it exists to catch
// regressions that would break a locale outright, such as a parser
// change that makes a real-world ".sor" file fail to compile.
func TestSmokeAllLanguages(t *testing.T) {
	numbers := []string{"0", "1", "15", "42", "100", "1999", "-7", "3.14"}
	for _, lang := range numbertext.Languages() {
		for _, n := range numbers {
			_, err := numbertext.Cardinal(lang, n)
			if err != nil {
				t.Errorf("Cardinal(%s, %s): %v", lang, n, err)
			}
		}

		_, err := numbertext.Ordinal(lang, "3")
		if err != nil {
			t.Errorf("Ordinal(%s, 3): %v", lang, err)
		}

		_, err = numbertext.Year(lang, "2001")
		if err != nil {
			t.Errorf("Year(%s, 2001): %v", lang, err)
		}
	}
}
