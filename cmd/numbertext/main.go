// Command numbertext spells out a number as words using the
// github.com/casara/numbertext-go library, for quick manual testing and
// as a runnable usage example.
package main

import (
	"flag"
	"fmt"
	"os"

	numbertext "github.com/casara/numbertext-go"
)

const (
	exitUsageError   = 2
	exitRuntimeError = 1
)

// options holds every flag main understands; see registerFlags.
type options struct {
	lang, cardinal, ordinal, ordinalNumber, year, currency, amount string
	money, help, list                                              bool
}

func registerFlags() *options {
	var o options
	flag.StringVar(&o.lang, "lang", "en", "language code (see -list)")
	flag.BoolVar(&o.list, "list", false, "list available language codes and exit")
	flag.StringVar(&o.cardinal, "cardinal", "", "spell out a cardinal number, e.g. -cardinal 123")
	flag.StringVar(&o.ordinal, "ordinal", "", "spell out an ordinal number, e.g. -ordinal 3")
	flag.StringVar(
		&o.ordinalNumber, "ordinal-number", "",
		"format an abbreviated ordinal, e.g. -ordinal-number 3",
	)
	flag.StringVar(&o.year, "year", "", "spell out a calendar year, e.g. -year 1999")
	flag.StringVar(&o.currency, "currency", "", "ISO 4217 currency code, used with -amount")
	flag.BoolVar(&o.money, "money", false, "use the MONEYTEXT-style section with -currency/-amount")
	flag.StringVar(&o.amount, "amount", "", "currency amount, used with -currency")
	flag.BoolVar(&o.help, "help-section", false, "print the locale's self-documenting help section")
	flag.Parse()

	return &o
}

func main() {
	o := registerFlags()

	if o.list {
		for _, l := range numbertext.Languages() {
			fmt.Println(l)
		}

		return
	}

	out, err := run(o)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(exitRuntimeError)
	}
	fmt.Println(out)
}

// run dispatches to the single numbertext function the given flags
// select. It exits the process directly (instead of returning an error)
// when no action flag was given, since that case prints usage rather
// than a numbertext error.
func run(o *options) (string, error) {
	switch {
	case o.cardinal != "":
		return numbertext.Cardinal(o.lang, o.cardinal)
	case o.ordinal != "":
		return numbertext.Ordinal(o.lang, o.ordinal)
	case o.ordinalNumber != "":
		return numbertext.OrdinalNumber(o.lang, o.ordinalNumber)
	case o.year != "":
		return numbertext.Year(o.lang, o.year)
	case o.currency != "" && o.money:
		return numbertext.Money(o.lang, o.currency, o.amount)
	case o.currency != "":
		return numbertext.Currency(o.lang, o.currency, o.amount)
	case o.help:
		return numbertext.Help(o.lang)
	default:
		printUsage()
		os.Exit(exitUsageError)

		return "", nil
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "usage: numbertext -lang <code> [action]")
	fmt.Fprintln(os.Stderr, "  actions: -cardinal|-ordinal|-ordinal-number|-year <n>")
	fmt.Fprintln(os.Stderr, "           -currency <ISO> -amount <n> [-money]")
	fmt.Fprintln(os.Stderr, "           -help-section | -list")
	flag.PrintDefaults()
}
