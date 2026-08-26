---
name: using-numbertext-go
description: >-
  Guidance for writing Go code that uses numbertext-go
  (github.com/casara/numbertext-go) - spelling out numbers, ordinals,
  years, currency amounts, and gender/case variants in any of its
  bundled languages, plus registering a new one at runtime. Use when
  the project imports github.com/casara/numbertext-go or when asked to
  add/change a number-to-text conversion, add a language, or reduce a
  binary's size by selecting which locales get embedded.
---

# Using numbertext-go

`numbertext-go` spells out numbers as words - cardinals, ordinals, years,
currency amounts, and gender/case variants - in every language shipped by
[Numbertext.org](https://numbertext.github.io/). This skill is for code that
*consumes* numbertext-go as a dependency, not for numbertext-go's own
source - if you're working inside the `numbertext-go` repo itself, follow
its `AGENTS.md` instead.

This file is meant to be copied into any project that depends on
numbertext-go, so it never lives next to numbertext-go's own source tree.
Every path mentioned below (`data/*.sor`, `README.md`, `scripts/*`) is
therefore a full link into the numbertext-go repository itself
(<https://github.com/casara/numbertext-go>), not a path relative to this
file.

## The core functions

All numeric arguments are strings (`"123"`, `"-5"`, `"3.14"`), not `int` or
`float64` - the underlying engine (the "Soros" rule interpreter) is
text-based, which keeps arbitrary-precision numbers exact (a `.sor` file
can spell out decillions). `CardinalInt(lang string, n int64)` is a
convenience wrapper for the common integer case; everything else takes a
string.

```go
c, err := numbertext.Cardinal("en", "123")   // "one hundred twenty-three"
o, err := numbertext.Ordinal("en", "123")    // "one hundred twenty-third"
n, err := numbertext.OrdinalNumber("en", "123") // "123rd"
y, err := numbertext.Year("en", "1999")      // "nineteen ninety-nine"
cur, err := numbertext.Currency("en", "USD", "2.5") // "two U.S. dollars and fifty cents"
m, err := numbertext.Money("en", "USD", "2.5")      // "two and 50/100 U.S. dollars"
h, err := numbertext.Help("en")              // the locale's own usage summary
```

Every one of these returns `(string, error)` - check the error. The most
common cause is an unknown language code (see "Errors" below), but a
custom locale registered via `RegisterLocale` can also fail to compile if
its own program is malformed.

`Currency` and `Money` are two different upstream-defined renderings of
the same amount, not interchangeable spellings - `Currency` is the
NUMBERTEXT()-compatible form, `Money` is the MONEYTEXT()-compatible one
some locales format differently (e.g. a fraction for sub-unit remainders
instead of a second currency name). Pick the one your output format
actually expects; don't assume they always agree.

## `Convert`: the primitive everything else is built on

`Convert(lang, prefix, arg string) (string, error)` builds one input string
and runs it through the language's compiled program:

- `prefix == ""`: the input is just `arg` (this is what `Cardinal` does).
- `arg == ""`: the input is just `prefix` (what `Help` does).
- otherwise: the input is `prefix + " " + arg` (what `Ordinal`, `Year`,
  `Currency`, `Money` do).

Reach for `Convert` directly only when you need a section name none of the
named wrappers cover - see "Gender, case, and other variants" below.

## Regional variants (e.g. "en-GB", "pt-BR")

A regional code doesn't need its own data file: the base language's file
carries a handful of lines tagged for that region, and passing the full
code (`"en-GB"`, `"pt-BR"`) activates them automatically while falling
back to the base file for everything else. `Languages()` only ever lists
base codes - the absence of `"en-GB"` from that list does **not** mean
`Cardinal("en-GB", ...)` will fail. It fails only if the part before the
first `-` isn't a known base language at all.

## Gender, case, and other locale-specific variants

A language's rule file isn't limited to a fixed cardinal/ordinal shape -
it can define arbitrary named sections for whatever grammatical variant it
needs (`pt.sor` has `feminine`/`masculine`; `ru.sor` has
`cardinal-feminine`/`cardinal-neuter`; `it.sor` has `ordinal-masculine`;
and so on). **These section names are not standardized across
languages** - there is no `Gender(...)` function, and guessing a name that
"should" exist for a language you haven't checked is a common mistake.
Call `Convert` with the section name instead, and confirm it exists first
via `Help(lang)` or by reading that language's `.sor` file directly:

```go
fem, err := numbertext.Convert("pt", "feminine", "2") // "duas"
```

## Selecting which languages get embedded

By default, every bundled `.sor` file (52 languages) is embedded via
`go:embed`, regardless of which ones your program actually calls. If your
program only ever needs a handful of languages and you want a smaller
binary, opt in with build tags: `numbertext_select` plus one
`numbertext_lang_<code>` tag per language you want
(`<code>` is the `.sor` file's base name, lowercased, e.g. `en`, `pt`,
`hu_hung`):

```sh
go build -tags "numbertext_select numbertext_lang_en numbertext_lang_pt" .
```

Get this wrong and the failure mode is silent at compile time but loud at
runtime: `Languages()` returns only the codes you tagged, and calling
`Cardinal` (or anything else) with an untagged language code returns
`ErrUnknownLanguage` - it will not fail to build. If you add a call to a
new language anywhere in the program, add its `numbertext_lang_<code>` tag
to every build/test/lint invocation, not just `go build`.

Passing `numbertext_select` **without** any `numbertext_lang_<code>` tag
embeds nothing at all - only useful if every locale your program needs
comes from `RegisterLocale` instead.

## Adding a language at runtime

Because the whole library is data-driven (a `.sor` program, not generated
Go code), a consumer can add a language numbertext-go doesn't bundle
without forking it or waiting on a release:

```go
err := numbertext.RegisterLocale("es-mini", sorSource) // sorSource is a ".sor" file's contents
```

`RegisterLocale` compiles `sorSource` immediately and returns an error if
it doesn't parse - it never registers a broken program silently. Registering
a code that already exists (bundled or previously registered) replaces it.
Do this once, at startup, before any call that might use the new code -
`RegisterLocale` is safe to call concurrently with lookups, but a call
racing its own registration is still a logic bug in the caller, not
something to rely on the library to order for you.

Treat `sorSource` as configuration you control, not as end-user input:
compiling a `.sor` program builds regular expressions from it (see
numbertext-go's own `SECURITY.md`/`AGENTS.md` for the RE2/regexp2 hybrid
engine and its bounded fallback for the rare backreference pattern) and,
like any recursive interpreter, a hostile custom program could in
principle be crafted to recurse deeply on a crafted input.

## Errors

Two sentinel errors are returned across the whole package, meant to be
matched with `errors.Is`, not by inspecting the message text:

- `numbertext.ErrUnknownLanguage` - the language code (and, for a regional
  variant, its base language too) has no bundled or registered source.
  This is the one every call site touching a caller-supplied language
  code should check for.
- `numbertext.ErrEmptyLocaleCode` - returned by `RegisterLocale` when
  `code == ""`.

```go
_, err := numbertext.Cardinal(userLang, n)
if errors.Is(err, numbertext.ErrUnknownLanguage) {
    // fall back to a default language, or report it to the user
}
```

## Reference

Full rationale for the above (the Soros rule-file format, the regex
engine's RE2/regexp2 hybrid design, and why gender/case sections aren't
standardized) lives in the numbertext-go repo itself: `README.md` (usage
and API stability), `AGENTS.md` (architecture and the Soros language), and
`data/SOURCE.md` (adding a language as a data file instead of via
`RegisterLocale`).
