#!/usr/bin/env bash
set -euo pipefail

gremlins=${1:?usage: mutation.sh GREMLINS MIN_EFFICACY PACKAGE}
min=${2:?usage: mutation.sh GREMLINS MIN_EFFICACY PACKAGE}
package=${3:?usage: mutation.sh GREMLINS MIN_EFFICACY PACKAGE}

report=$("$gremlins" unleash --workers 2 --timeout-coefficient 20 "$package" 2>&1) || {
  echo "$report" >&2
  exit 1
}
echo "$report" | grep -E 'LIVED|NOT COVERED|Killed|Timed out|Test efficacy|Mutator coverage' || true

efficacy=$(echo "$report" | awk '/^Test efficacy:/ { sub("%", "", $3); print $3 }')
if [[ -z $efficacy ]]; then
  echo "mutation: no efficacy found in the gremlins report" >&2
  exit 1
fi
if awk -v e="$efficacy" -v m="$min" 'BEGIN { exit !(e < m) }'; then
  echo "mutation: test efficacy is ${efficacy}%, minimum is ${min}% (ADR 0012)" >&2
  exit 1
fi
