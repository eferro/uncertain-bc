#!/usr/bin/env bash
set -euo pipefail

if matches=$(grep -rnE '\bt\.(Error|Errorf|Fatal|Fatalf)\(' --include='*_test.go' .); then
  echo "test-style: use testify require/assert instead of t.Error*/t.Fatal* (ADR 0009):" >&2
  echo "$matches" >&2
  exit 1
fi
