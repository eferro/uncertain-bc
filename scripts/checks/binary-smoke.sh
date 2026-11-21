#!/usr/bin/env bash
set -euo pipefail

bin=${1:?usage: binary-smoke.sh BINARY EXPECTED_VERSION}
version=${2:?usage: binary-smoke.sh BINARY EXPECTED_VERSION}

expect() {
  local want=$1 got=$2 what=$3
  if [[ $got != "$want" ]]; then
    echo "binary-smoke: $what printed '$got', want '$want'" >&2
    exit 1
  fi
}

expect "ubc $version" "$("$bin" --version)" "--version (is main.version injected by -ldflags?)"
expect "3" "$("$bin" -e '1 + 2')" "-e '1 + 2'"
