#!/usr/bin/env bash
set -euo pipefail

profile=${1:?usage: coverage.sh PROFILE MIN_PERCENT}
min=${2:?usage: coverage.sh PROFILE MIN_PERCENT}

total=$(go tool cover -func="$profile" | awk '/^total:/ { sub("%", "", $3); print $3 }')

if awk -v t="$total" -v m="$min" 'BEGIN { exit !(t < m) }'; then
  echo "coverage: internal/ is at ${total}%, minimum is ${min}% (ADR 0012)" >&2
  exit 1
fi
echo "coverage: ${total}% (minimum ${min}%)"
