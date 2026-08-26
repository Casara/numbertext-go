#!/usr/bin/env bash
# Re-download every *.sor file from the upstream libnumbertext repository
# (https://github.com/Numbertext/libnumbertext), overwriting data/*.sor.
#
# Usage: scripts/sync-data.sh [git-ref]
#   git-ref defaults to "master".
#
# After running this, review `git diff data/`, update the synced commit
# and date recorded in data/SOURCE.md, and re-run the test suite (in
# particular TestSmokeAllLanguages) before committing the refresh.
set -euo pipefail

ref="${1:-master}"
repo="Numbertext/libnumbertext"
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
data_dir="$script_dir/../data"

echo "Fetching data/*.sor from ${repo}@${ref} ..." >&2

commit="$(curl -fsSL "https://api.github.com/repos/${repo}/commits/${ref}" | python3 -c 'import json,sys; print(json.load(sys.stdin)["sha"])')"

curl -fsSL "https://api.github.com/repos/${repo}/contents/data?ref=${ref}" |
	python3 -c '
import json, sys
for entry in json.load(sys.stdin):
    if entry["name"].endswith(".sor"):
        print(entry["name"])
' |
	while IFS= read -r name; do
		echo "  $name" >&2
		curl -fsSL "https://raw.githubusercontent.com/${repo}/${ref}/data/${name}" -o "$data_dir/$name"
	done

# Local patch, pending upstream (see data/SOURCE.md#local-patches-pending-upstream):
# bg.sor line ~146 has "]" and ")" transposed in an ordinal-thousands
# rule, breaking on any parser that validates bracket/paren nesting
# (this port's included) even though it's a plain two-character typo,
# not a deliberate construct - confirmed by comparing against the
# (correctly nested) singular rule right above it. Reapply on every
# sync until upstream has the fix.
python3 - "$data_dir/bg.sor" <<'EOF'
import sys

path = sys.argv[1]
broken = r"\2)])"
fixed = r"\2))]"

with open(path, encoding="utf-8") as f:
    content = f.read()

if broken not in content:
    print("  bg.sor no longer needs the local patch - upstream likely merged the fix.", file=sys.stderr)
    print("  Remove the patch step above and the matching section in data/SOURCE.md.", file=sys.stderr)
else:
    with open(path, "w", encoding="utf-8") as f:
        f.write(content.replace(broken, fixed, 1))
    print("  reapplied the bg.sor local patch (upstream doesn't have the fix yet)", file=sys.stderr)
EOF

echo "Done. Synced commit: ${commit}" >&2
echo "Update data/SOURCE.md with this commit hash and today's date." >&2
