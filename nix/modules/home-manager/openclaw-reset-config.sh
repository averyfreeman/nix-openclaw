#!/bin/sh
set -eu

if [ "$#" -ne 2 ]; then
  echo "usage: openclaw-reset-config <seed-file> <config-file>" >&2
  exit 1
fi

seed="$1"
target="$2"
backup_dir="$(dirname "$target")/.backups"
timestamp="$(date +%Y%m%d-%H%M%S)"

if [ -e "$target" ] || [ -L "$target" ]; then
  mkdir -p "$backup_dir"
  cp -RL "$target" "$backup_dir/openclaw.json.$timestamp"
fi

mkdir -p "$(dirname "$target")"
cp -L "$seed" "$target"
chmod u+rw "$target"
echo "Reset $target from the Nix seed; backups are in $backup_dir" >&2
