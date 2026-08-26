#!/usr/bin/env bash
#
# Reports the documentation a change to the public API has left behind.
#
# `make check` proves the code is correct and `verify-docs` compiles the
# snippets that can be compiled, but neither can tell that a section of
# README.md still describes a signature that changed, or that AGENTS.md
# names an invariant that no longer holds. Those are fragments and
# prose; nothing compiles them. This is the residue, and this script
# names it.
#
# It is a reminder, not a gate: it always exits 0. A diff can legitimately
# touch an exported symbol without any of these needing an edit - a doc
# comment fix, for instance - and a check that cried wolf would be turned
# off within a week.
#
# Compares against HEAD by default, or against the ref given as $1.
set -euo pipefail

cd "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

base="${1:-HEAD}"

changed=$(git diff --name-only "$base" -- '*.go' ':!*_test.go')

if [ -z "$changed" ]; then
    exit 0
fi

# A line added or removed at the start of an exported declaration: a
# top-level func/type/const/var, or an exported struct field (one tab,
# capital letter). Deliberately a heuristic over the diff rather than a
# real API diff - the point is to prompt a human, and `go doc` on two
# worktrees would cost more than the reminder is worth.
surface_changed=$(
    git diff -U0 "$base" -- '*.go' ':!*_test.go' |
        grep -E '^[+-](func|type|const|var) [A-Z]|^[+-]\t[A-Z][A-Za-z0-9_]* ' |
        grep -vE '^[+-](func|type|const|var) [A-Z][A-Za-z0-9_]*Test' |
        head -1 ||
        true
)

if [ -z "$surface_changed" ]; then
    exit 0
fi

# Everything worth checking by hand when the public API moves, plus the
# reason each one can't be checked automatically. locale_<code>.go and
# locale_all.go are deliberately absent - they're generated (see
# scripts/gen-locale-embeds.py) and never describe the API themselves.
declare -a targets=(
    "README.md|the prose around Usage/Regional variants/Gender, which verify-docs does not read"
    "README.pt-BR.md|same, and it has to stay in sync with README.md"
    "AGENTS.md|if an invariant changed (CLAUDE.md just imports it)"
    "CHANGELOG.md|every public change belongs here"
    "skills/using-numbertext-go/SKILL.md|it documents the same public API for consumer projects"
)

stale=()

for entry in "${targets[@]}"; do
    path="${entry%%|*}"
    reason="${entry#*|}"

    if ! git diff --quiet "$base" -- "$path" 2>/dev/null; then
        continue
    fi

    stale+=("  $path
      $reason")
done

if [ ${#stale[@]} -eq 0 ]; then
    exit 0
fi

printf '\nThis diff changes the exported API. Untouched documentation:\n\n'
printf '%s\n' "${stale[@]}"
printf '\nIf none of the above needs an edit, nothing to do.\n\n'
