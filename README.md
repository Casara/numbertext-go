# numbertext-go

[![CI](https://github.com/casara/numbertext-go/actions/workflows/ci.yml/badge.svg)](https://github.com/casara/numbertext-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/casara/numbertext-go.svg)](https://pkg.go.dev/github.com/casara/numbertext-go)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A Go library that spells out numbers as words — cardinals, ordinals,
gender/case variants, negatives, and currency amounts — in every
language shipped by [Numbertext.org](https://numbertext.github.io/).

Leia isto em [português (Brasil)](README.pt-BR.md).

## Why

[Numbertext.org](https://numbertext.github.io/) (upstream:
[libnumbertext](https://github.com/Numbertext/libnumbertext)) describes
number-to-text conversion for dozens of languages as small, declarative
rule files written in a tiny regex-rewriting DSL called **Soros**, rather
than as hand-written code per language. Adding a language — or a
regional/gender variant of one — means writing a new `.sor` file, not
writing new code.

This repository is a from-scratch Go implementation of the Soros
interpreter (see [`internal/soros`](internal/soros)) plus a friendly
public API (see [`numbertext.go`](numbertext.go)), reusing the exact
`.sor` rule files from upstream (see [`data/`](data), synced from
[`data/SOURCE.md`](data/SOURCE.md)) unmodified.

## Install

```sh
go get github.com/casara/numbertext-go
```

## Usage

```go
package main

import (
	"fmt"

	numbertext "github.com/casara/numbertext-go"
)

func main() {
	cardinal, _ := numbertext.Cardinal("en", "12345")
	fmt.Println(cardinal) // "twelve thousand three hundred forty-five"

	ordinal, _ := numbertext.Ordinal("en", "21")
	fmt.Println(ordinal) // "twenty-first"

	year, _ := numbertext.Year("en", "1999")
	fmt.Println(year) // "nineteen ninety-nine"

	amount, _ := numbertext.Currency("en", "USD", "2.50")
	fmt.Println(amount) // "two U.S. dollars and fifty cents"

	negative, _ := numbertext.Cardinal("pt", "-5")
	fmt.Println(negative) // "menos cinco"
}
```

### API

| Function | Purpose |
| --- | --- |
| `Cardinal(lang, n)` | "123" → "one hundred twenty-three" |
| `Ordinal(lang, n)` | "123" → "one hundred twenty-third" |
| `OrdinalNumber(lang, n)` | "123" → "123rd" |
| `Year(lang, n)` | "1999" → "nineteen ninety-nine" |
| `Currency(lang, isoCode, amount)` | NUMBERTEXT()-style: "USD", "2.5" → "two U.S. dollars and fifty cents" |
| `Money(lang, isoCode, amount)` | MONEYTEXT()-style: "USD", "2.5" → "two and 50/100 U.S. dollars" |
| `Help(lang)` | the locale's self-documenting usage section |
| `Convert(lang, prefix, arg)` | the primitive above everything else — see below |
| `RegisterLocale(code, sorSource)` | add a language at runtime from raw `.sor` text |
| `Languages()` | every currently known language code |

All numeric arguments are strings (`"123"`, `"-5"`, `"3.14"`), not `int`/
`float64`, because the underlying engine is text-based and this keeps
arbitrary-precision numbers (a `.sor` file can spell out decillions)
exact. `CardinalInt(lang, n int64)` is provided as a convenience for the
common case.

### Regional variants (e.g. "en-GB", "pt-BR")

A regional code doesn't need (and won't have) its own `.sor` file:
`data/en.sor` and `data/pt.sor` each carry the base language's rules plus
a handful of lines tagged `[:en-GB:]`, `[:pt-BR:]`, etc. for wording that
differs by region. Just pass the full code — `Languages()` still only
lists base codes (`"en"`, `"pt"`, ...), since the region isn't a separate
file, but any `"<base>-<REGION>"` string works and activates that
region's tagged lines automatically:

```go
gb, _ := numbertext.Cardinal("en-GB", "101")
fmt.Println(gb) // "one hundred and one"

us, _ := numbertext.Cardinal("en", "101")
fmt.Println(us) // "one hundred one" (no "and" — the plain/default wording)

br, _ := numbertext.Cardinal("pt-BR", "16")
fmt.Println(br) // "dezesseis"

pt, _ := numbertext.Cardinal("pt", "16")
fmt.Println(pt) // "dezasseis" (European Portuguese default)
```

If the base part before the first `-` isn't a known language, `Cardinal`
(and friends) return an error — there's no file to fall back to.

### Selecting which languages get embedded

By default, `go build` embeds all 52 bundled `.sor` files into your
binary via `go:embed` (see `locale_all.go`), regardless of which
languages your program actually calls. If you only need a few, you can
opt into a smaller binary with build tags: pass `numbertext_select` plus
one `numbertext_lang_<code>` tag per language you want (`<code>` is the
`.sor` file's base name, lowercased, e.g. `en`, `pt`, `hu_hung`):

```sh
go build -tags "numbertext_select numbertext_lang_en numbertext_lang_pt" .
```

This embeds only `en` and `pt` — `numbertext.Languages()` returns exactly
`["en", "pt"]`, and calling e.g. `Cardinal("de", ...)` returns an
"unknown language" error. Passing `numbertext_select` **without** any
`numbertext_lang_<code>` tag embeds nothing at all (`Languages()` returns
an empty slice); that's only useful if you plan to supply every locale
yourself via `RegisterLocale`.

Omitting `numbertext_select` entirely (the default) always embeds every
bundled language — this switch never changes behavior unless you opt in.
The per-language files (`locale_<code>.go`) are generated by
`scripts/gen-locale-embeds.py`; regenerate them after
`scripts/sync-data.sh` adds or removes a `data/*.sor` file.

### Gender, case, and other locale-specific variants

Numbertext's rule files aren't limited to a fixed cardinal/ordinal
shape — a language can define arbitrary named sections for whatever
grammatical variants it needs (`pt.sor` has `feminine`/`masculine`;
`ru.sor` has `cardinal-feminine`/`cardinal-neuter`; `it.sor` has
`ordinal-masculine`; and so on). These section names are **not**
standardized across languages, so this library does not pretend
otherwise with a one-size-fits-all `Gender(...)` function. Instead, call
`Convert` with the section name you need:

```go
fem, _ := numbertext.Convert("pt", "feminine", "2")
fmt.Println(fem) // "duas"
```

Run `Help(lang)` (or open the `.sor` file directly) to see which
sections a given language defines.

### Adding a new language

Because the whole library is data-driven, you don't need to touch any
Go code to add a language:

```go
sorSource := "1 uno\n2 dos\n3 tres\n" // a real file needs far more rules
err := numbertext.RegisterLocale("es-mini", sorSource)
```

or add a file under [`data/`](data) — see
[`data/SOURCE.md`](data/SOURCE.md#adding-a-new-language) for a short
guide and a link to the upstream Soros language specification.

## API stability

`numbertext-go` is released as **v0.x**, and in Semantic Versioning that
carries a specific meaning worth stating plainly rather than leaving you
to infer it: **a minor release may break the public API.** Pin a
version, read the [CHANGELOG](CHANGELOG.md) before upgrading, and expect
to make small mechanical changes when you do.

What counts as the public API: every exported symbol in the root
`numbertext` package and `cmd/numbertext`'s flags. Not public: anything
under `internal/` (including the Soros interpreter itself), and the
exact wording of error strings — `ErrUnknownLanguage`/`ErrEmptyLocaleCode`
(meant for `errors.Is`) are the stable part of error handling, matching
on message text is not supported.

Breaking changes will be listed under `### Changed` or `### Removed` in
the changelog, with the migration in the same entry. Where a rename can
keep the old spelling compiling, it will — as a deprecated alias.

## Skill for AI assistants

[skills/using-numbertext-go/SKILL.md](skills/using-numbertext-go/SKILL.md)
documents, in a format AI tools can load on demand (Claude Skill), how to
use `numbertext-go` idiomatically: the core functions, the `Convert`
primitive, regional variants, gender/case sections, and selective
language embedding. To use it in a project that depends on
`numbertext-go`, copy the directory into it:

```sh
cp -r skills/using-numbertext-go <your-project>/.claude/skills/using-numbertext-go
```

## Development

`make help` lists every command. Before opening a PR:

```sh
make check     # lint + arch-lint + verify-docs + tests with the race detector + test-select
make lint-md   # markdown lint (requires Node.js >= 20)
```

Style conventions are documented in
[docs/coding-style.md](docs/coding-style.md).

See [CONTRIBUTING.md](CONTRIBUTING.md) for the full contribution
workflow (commit message convention, branch/PR process) and
[CLAUDE.md](CLAUDE.md)/[AGENTS.md](AGENTS.md) for a deeper architecture
overview aimed at contributors and coding agents.

## License

The Go code in this repository is MIT-licensed, see [LICENSE](LICENSE).
The `.sor` rule files under [`data/`](data) are copied unmodified from
upstream libnumbertext and remain under its BSD-3-Clause license, see
[`data/LICENSE`](data/LICENSE) and [`data/SOURCE.md`](data/SOURCE.md).
