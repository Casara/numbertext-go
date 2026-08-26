<!-- markdownlint-disable-next-line MD041 -- GitHub PR templates aren't rendered as a titled page, no H1 needed. -->
## What and why

<!-- What changes, and the motivation - link an issue if there is one. -->

## Checklist

- [ ] `make check` passes locally (lint + arch-lint + verify-docs + tests
      with the race detector + the selective-embedding build-tag sanity check).
- [ ] `make lint-md` passes, if any Markdown file changed.
- [ ] New/changed behavior has a test. For a new language, a handful of
      known-correct conversions (see `numbertext_test.go`'s
      English/Portuguese tests) - verified by actually running the code,
      not guessed (see AGENTS.md's "Verify locale output" note).
- [ ] Commit messages follow [CONTRIBUTING.md](../CONTRIBUTING.md)'s
      Conventional Commits convention.

### If an exported symbol changed

<!-- Delete this section if the change is internal only. -->

- [ ] `README.md` / `README.pt-BR.md` - prose and examples still match.
- [ ] `AGENTS.md` (`CLAUDE.md` just imports it), if an invariant changed.
- [ ] `CHANGELOG.md` - and if this is a breaking change, read
      [Changing the public API](../CONTRIBUTING.md#changing-the-public-api).
- [ ] `make doc-sync` - lists what a change to an exported symbol may
      have left stale; it never fails the build, so check it by hand.

### If `data/*.sor` changed

<!-- Delete this section otherwise. -->

- [ ] Came from `scripts/sync-data.sh` (upstream sync), not a hand edit -
      see AGENTS.md's "Don't hand-edit data/*.sor".
- [ ] `data/SOURCE.md`'s synced commit/date updated.
- [ ] `scripts/gen-locale-embeds.py` re-run if a locale was added/removed
      (`make gen-locales`).
