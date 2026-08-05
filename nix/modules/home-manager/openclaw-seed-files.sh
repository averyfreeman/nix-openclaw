#!/bin/sh
set -eu

if [ "$#" -ne 1 ]; then
  echo "usage: openclaw-seed-files <source-target-manifest>" >&2
  exit 1
fi

manifest="$1"

while IFS="$(printf '\t')" read -r source target; do
  [ -n "$source" ] || continue
  [ -n "$target" ] || continue

  if [ -e "$target" ] || [ -L "$target" ]; then
    continue
  fi

  mkdir -p "$(dirname "$target")"
  if [ -d "$source" ]; then
    cp -RL "$source" "$target"
  else
    cp -L "$source" "$target"
  fi
done < "$manifest"
