# Data source

The `.sor` files in this directory are copied verbatim from the `data/`
directory of the upstream [libnumbertext](https://github.com/Numbertext/libnumbertext)
project (the reference implementation behind [Numbertext.org](https://numbertext.github.io/)),
with one documented exception - see [Local patches](#local-patches-pending-upstream)
below.

- Upstream repository: <https://github.com/Numbertext/libnumbertext>
- Synced commit: `a4b0225813b015a0f796754bc6718be20dd9943c`
- Upstream `VERSION` file: `1.0.11`
- Sync date: 2026-08-21
- Verification: every local `*.sor` file's git blob SHA1 (`git hash-object`)
  matches the corresponding blob SHA reported by the GitHub Contents API
  for the commit above, **except `bg.sor`** (one line patched locally,
  see below).

## Local patches (pending upstream)

### `bg.sor` line 146: misplaced `]`/`)` in an ordinal-thousands rule

```diff
- ([0-9]{1,3})([0-9]{3}) $(f:|$1) хиляди[ $(and:$(ordinal \2)])
+ ([0-9]{1,3})([0-9]{3}) $(f:|$1) хиляди[ $(and:$(ordinal \2))]
```

Upstream's line ends `\2)])` (closing bracket before the second closing
paren); the correct, intended form - confirmed by comparing against the
singular counterpart one line above it, `1([0-9]{3}) хиляда[
$(and:$(ordinal \1))]`, which already has it right - closes both
parentheses before the bracket: `\2))]`. It reads as a plain
two-character transposition, not a deliberate construct: both rules
build the same "N thousand [and REMAINDER-th]" shape (`e3:`'s cardinal
counterpart a few lines up, `$(eN:\1)[ $(and:$2)]`, is the same pattern
again for every other magnitude), and the singular rule right next to
it already has the correct nesting.

This is why: a rule Compile can't parse (unclosed `[...]`, spec-wise -
see AGENTS.md's regex-engine note for what Compile *can* recover from
and why this isn't in that category) goes to `Program.Skipped` rather
than being guessed at - correct, since a parser has no way to tell "this
input is malformed" apart from "this is a legitimate construct this
port doesn't handle yet." A one-character-transposition typo, spotted by
comparing against a working sibling rule, is a different kind of claim:
verifiable by inspection, not a guess. See AGENTS.md's "Soros engine
and data files" section for the fuller reasoning trail (why patching
the *interpreter* to tolerate this was rejected in favor of patching
the *data*).

**Status**: reported upstream, not yet merged - see
[Numbertext/libnumbertext#139](https://github.com/Numbertext/libnumbertext/issues/139).
Once merged and picked up by a future `scripts/sync-data.sh` run, remove
this section and the matching re-patch step in that script.

## License

These files are distributed under the upstream project's BSD-3-Clause
license, reproduced in [LICENSE](LICENSE). This is a different license
from the one covering the Go source code in the rest of this repository
(see the top-level [LICENSE](../LICENSE)); only the contents of this
`data/` directory are BSD-3-Clause.

## Updating

Run `scripts/sync-data.sh` from the repository root to re-download every
`*.sor` file from upstream `master`, print the commit it synced to, and
reapply the `bg.sor` patch above (or report that upstream already has
the fix, in which case remove the patch step and the section above).
Review the diff, update the commit/date above, and re-run the test suite
(especially `TestSmokeAllLanguages`) before committing a refresh.

## Adding a new language

Numbertext.org's whole design point is that a new language only needs a
new `.sor` rule file, no code changes. To add one:

1. Read the [Soros language specification](https://github.com/Numbertext/libnumbertext/blob/master/doc/sorosspec.pdf)
   (or study an existing file such as `en.sor`, `pt.sor`, or `de.sor`).
2. Write `data/<code>.sor` (or use `numbertext.RegisterLocale` to add it
   at runtime without touching this repository at all).
3. Add a test asserting a handful of known conversions, similar to
   `numbertext_test.go`'s English/Portuguese tests.
4. Consider contributing the file back to upstream `libnumbertext` too,
   so other users of the Soros ecosystem benefit from it.
