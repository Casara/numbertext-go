// Package numbertext spells out numbers as words: cardinals, ordinals,
// years, and currency amounts, in every language shipped by
// Numbertext.org (https://numbertext.github.io/).
//
// Each supported language is described by a small rule file (a ".sor"
// program, see the "Soros" language implemented by the internal/soros
// package) rather than by generated Go code, so adding a new language –
// or a regional/gender variant of an existing one – never requires a
// code change: see RegisterLocale.
package numbertext

import "strconv"

// Convert is the primitive every other function in this package is
// built on. It runs lang's program against the global input built from
// prefix and arg:
//
//   - prefix == "": the input is just arg (used by Cardinal).
//   - arg == "": the input is just prefix (used by Help, or to call a
//     zero-argument section).
//   - otherwise: the input is prefix + " " + arg (used by Ordinal,
//     Year, Currency, Money, and any locale-specific section such as
//     pt.sor's "feminine"/"masculine" or ru.sor's "cardinal-neuter" –
//     see that language's Help output for the section names it defines).
//
// This mirrors how a ".sor" file's own "help" section documents and
// calls its sections, e.g. "$(ordinal 1)" or "$(year 1999)".
func Convert(lang, prefix, arg string) (string, error) {
	p, err := reg.program(lang)
	if err != nil {
		return "", err
	}
	var input string
	switch {
	case prefix == "":
		input = arg
	case arg == "":
		input = prefix
	default:
		input = prefix + " " + arg
	}

	return p.Run(input), nil
}

// Cardinal spells out n (e.g. "123", "-5", "3.14") as a cardinal number,
// such as "one hundred twenty-three".
func Cardinal(lang, n string) (string, error) {
	return Convert(lang, "", n)
}

// CardinalInt spells out n as a cardinal number.
func CardinalInt(lang string, n int64) (string, error) {
	return Cardinal(lang, strconv.FormatInt(n, 10))
}

// Ordinal spells out n as an ordinal number, such as "one hundred
// twenty-third".
func Ordinal(lang, n string) (string, error) {
	return Convert(lang, "ordinal", n)
}

// OrdinalNumber formats n as an abbreviated ordinal, such as "123rd".
func OrdinalNumber(lang, n string) (string, error) {
	return Convert(lang, "ordinal-number", n)
}

// Year spells out n as a calendar year, such as "nineteen ninety-nine"
// for 1999.
func Year(lang, n string) (string, error) {
	return Convert(lang, "year", n)
}

// Currency spells out amount (e.g. "2.5", "-10") as an amount of
// isoCode (an ISO 4217 currency code, e.g. "USD", "EUR"), such as "two
// dollars and fifty cents". This is the NUMBERTEXT()-compatible form;
// see Money for the MONEYTEXT()-compatible variant.
func Currency(lang, isoCode, amount string) (string, error) {
	return Convert(lang, "", isoCode+" "+amount)
}

// Money spells out amount of isoCode using the locale's "money" section
// (the MONEYTEXT()-compatible form, which some locales format slightly
// differently from Currency, e.g. using a fraction for sub-unit
// remainders).
func Money(lang, isoCode, amount string) (string, error) {
	return Convert(lang, "money", isoCode+" "+amount)
}

// Help returns the locale's self-documenting "help" section, if it
// defines one: a short usage summary along with the section names
// available for that language (cardinal/ordinal endings, gender
// variants, etc.), meant for interactive exploration rather than
// programmatic parsing.
func Help(lang string) (string, error) {
	return Convert(lang, "help", "")
}
