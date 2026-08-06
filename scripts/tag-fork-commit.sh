#!/usr/bin/env bash
set -euo pipefail

usage() {
  echo "usage: $0 [--dry-run] [COMMIT]" >&2
  exit 2
}

dry_run=false
if [[ "${1:-}" == "--dry-run" ]]; then
  dry_run=true
  shift
fi
if [[ "$#" -gt 1 ]]; then
  usage
fi

commit="${1:-HEAD}"
git rev-parse --verify "$commit^{commit}" >/dev/null

existing="$(git tag --points-at "$commit" --list 'v[0-9]*.[0-9]*.[0-9]*' | sort -V | tail -1)"
if [[ -n "$existing" ]]; then
  echo "$existing"
  exit 0
fi

latest="$(git tag --list 'v[0-9]*.[0-9]*.[0-9]*' | sort -V | tail -1)"
if [[ -z "$latest" ]]; then
  next="v0.0.0"
else
  version="${latest#v}"
  IFS=. read -r major minor patch <<<"$version"
  next="v${major}.${minor}.$((patch + 1))"
fi

if git rev-parse --verify "refs/tags/$next" >/dev/null 2>&1; then
  echo "tag $next already exists but does not point at $commit" >&2
  exit 1
fi

if "$dry_run"; then
  echo "$next"
  exit 0
fi

git tag -a "$next" "$commit" -m "nix-openclaw fork release $next"
echo "$next"
