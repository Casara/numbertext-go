# Contributing

*[Leia em português](CONTRIBUTING.pt-BR.md)*

This project follows the [Code of Conduct](CODE_OF_CONDUCT.md). By
participating, you're expected to follow it.

## Before you start

Read, in this order:

1. [AGENTS.md](AGENTS.md) - architecture, the Soros language condensed,
   known/deliberate limitations, non-obvious design decisions.
2. [docs/coding-style.md](docs/coding-style.md) - code conventions
   (formatting, error handling/wrapcheck, explicit dependencies).
3. [README.md](README.md) - the user-facing API and usage examples.
4. `data/SOURCE.md` - how the bundled `.sor` locale files are sourced
   and how to add a new one.

## Environment

`make help` lists every command. Before opening a PR, run:

```sh
make check     # lint + arch-lint + verify-docs + tests with the race detector + test-select
make lint-md   # markdown lint (requires Node.js >= 20)
```

That's the minimum CI runs on every PR
([.github/workflows/ci.yml](.github/workflows/ci.yml)) - split across
separate jobs there (`test` per Go version, `lint` once, plus markdown
lint) rather than one `make check` call; see AGENTS.md's "Commands"
section for why. If the change touches behavior -
not just refactoring - include a test covering the new case, and, for a
new/changed language, a handful of known-correct conversions **verified
by actually running the code** (`go run ./cmd/numbertext -lang <x>
-cardinal <n>` or a throwaway test), not guessed - see AGENTS.md's
"Verify locale output before writing it down."

## Branch workflow

Rebase, not merge commits: update your branch with `git rebase`
against the base before opening/updating a PR, instead of merging the
base into your branch. Linear history, no merge commits.

The commit convention below is mandatory on `main`. A work branch that
will be squashed on merge doesn't need to follow it strictly - the
intermediate history isn't what stays.

## Commit messages

[Conventional Commits](https://www.conventionalcommits.org/), in
English, imperative mood:

```text
<type>(<scope>): <description>
```

* `<description>` in the imperative, describing the action (`add`,
  `fix`, `remove`, `harden`, `resolve`, `guarantee`), not the
  resulting state (`added`) or history (`fixed`, "adds/added").
* `<scope>` is the affected package or area. It's not a closed list,
  but the most common ones are: `soros` (the interpreter),
  `numbertext` (the public API/registry), `cmd` (the CLI), `data`
  (the bundled `.sor` locale files) - plus the cross-cutting
  `release`, `docs`, `ci`, `deps`.

### Types

| Type | When to use |
| --- | --- |
| `feat` | New feature (semver MINOR). |
| `fix` | Bug fix (semver PATCH). |
| `docs` | Documentation only (README, comments) - no code change. |
| `test` | Test only (added, changed, or removed) - no production code change. |
| `refactor` | Changes how code is written/organized without changing observable behavior. |
| `perf` | A change whose purpose is performance. |
| `style` | Formatting, semicolons, whitespace, lint - no code change. |
| `build` | Build and dependencies (`go.mod`, `Makefile`, tooling). |
| `ci` | Continuous integration (`.github/workflows/`). |
| `chore` | Maintenance tasks that don't fit the types above (config, `.gitignore`, ...). |
| `cleanup` | Removes commented-out, dead, or unnecessary code - no behavior change. |
| `remove` | Removes an obsolete/unused file, directory, or feature. |
| `raw` | Change to a configuration/data/parameter file that doesn't fit the types above (e.g. a `data/*.sor` sync). |

### Examples

```text
feat(numbertext): add RegisterLocale for adding a language at runtime
fix(soros): propagate the leading boundary through a piped call argument
docs(readme): explain regional variant fallback for en-GB/pt-BR
raw(data): sync data/*.sor to libnumbertext@a4b0225
```

## Changing the public API

The project is on **v0.x**, so a breaking change is allowed - see
"API stability" in the README for what that means for users. It still
has to be deliberate and visible:

* Say so in the PR description, and add a `### Changed` or `### Removed`
  entry to `CHANGELOG.md` with the migration in the same entry.
* Where a rename can keep the old spelling compiling, keep it as a
  deprecated alias rather than deleting it outright.
* Prefer an additive change when one exists: a new optional parameter
  via a variadic option, a new function alongside the old one.

Error strings are explicitly *not* part of the public API; matching on
their text is not supported. Sentinel errors (`ErrUnknownLanguage`,
`ErrEmptyLocaleCode`, meant for `errors.Is`) are.

## Contributing with an AI assistant

This project is built with one, and there is nothing to hide or disclose about
that: the code is judged the same either way, and a PR is not marked.

What is asked is the same thing asked of anyone:

* **Understand what you are submitting.** If you cannot explain why a change is
  correct, it is not ready - regardless of what wrote it.
* **Verify locale output, don't guess it.** A `.sor` file encodes real,
  sometimes surprising per-language stylistic choices (see AGENTS.md's
  own example: `en` vs `en-GB`'s "and", `Cardinal` vs `Year`'s "and"
  removal). Run the code before writing an expected value down.
* **Do not let a tool re-litigate a settled decision.** `AGENTS.md`
  documents this port's known, deliberate limitations (the two classes
  of upstream rule it cannot express) because that boundary was reached
  once, on purpose - a PR "fixing" it by guessing at the missing
  behavior gets closed with a pointer, not a debate.
* **Run the checks locally.** `make check` and `make lint-md`, before
  opening the PR rather than after CI says so.

`AGENTS.md` is read by Codex and Cursor directly; `CLAUDE.md` imports it for
Claude Code. Anything you add for one of them belongs in `AGENTS.md`, so the
others get it too.

## Pull requests

* `make check` passing is mandatory, not optional. If the change
  touches any Markdown file, `make lint-md` too.
* A small PR focused on one change is preferable to a large PR
  covering several unrelated things - but that's judgment, not a
  strict rule.
* If the change alters documented behavior, update the docs in the
  same PR.
