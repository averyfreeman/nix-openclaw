#!/bin/sh
set -eu

script="${OPENCLAW_SEED_FILES:?OPENCLAW_SEED_FILES is required}"
work="$(mktemp -d)"
mkdir -p "$work/source" "$work/target"
printf '%s\n' seeded > "$work/source/config.json"
printf '%s\t%s\n' "$work/source/config.json" "$work/target/config.json" > "$work/manifest"

"$script" "$work/manifest"
grep -q '^seeded$' "$work/target/config.json"

printf '%s\n' user-owned > "$work/target/config.json"
"$script" "$work/manifest"
grep -q '^user-owned$' "$work/target/config.json"

printf '%s\n' new > "$work/source/new.txt"
printf '%s\t%s\n' "$work/source/new.txt" "$work/target/new.txt" >> "$work/manifest"
"$script" "$work/manifest"
grep -q '^new$' "$work/target/new.txt"
