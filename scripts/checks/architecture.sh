#!/usr/bin/env bash
set -euo pipefail

module=$(go list -m)
status=0

forbid() {
  local layer=$1 pattern=$2 reason=$3 allow=${4:-^$} violations
  [[ -d internal/$layer ]] || return 0
  violations=$(go list -f '{{.ImportPath}}: {{join .Imports " "}}' "./internal/$layer/..." |
    awk -v re="$pattern" -v allow="$allow" '{ for (i = 2; i <= NF; i++) if ($i ~ re && $i !~ allow) print $1 " imports " $i }')
  if [[ -n $violations ]]; then
    echo "architecture: $reason (ADR 0010):" >&2
    echo "$violations" >&2
    status=1
  fi
}

forbid domain "^$module/internal/(application|adapters)" "domain must not depend on outer layers"
forbid domain "^(os|os/.*|net|net/.*|io/ioutil)$" "domain must stay free of I/O"
forbid domain "^[^/]+[.][^/]+/" "domain may only use the standard library" "^$module/internal/domain"
forbid application "^$module/internal/adapters" "application must not depend on adapters"
forbid application "^(os|os/.*|net|net/.*|io/ioutil)$" "application reaches I/O only through ports"

exit $status
