#!/usr/bin/env bash
set -euo pipefail

subject=${1:?usage: commit-subject.sh SUBJECT}
pattern='^(build|chore|ci|docs|feat|fix|perf|refactor|revert|style|test)(\([a-z0-9 /-]+\))?!?: .+'

if [[ ! $subject =~ $pattern ]]; then
  echo "commit-subject: must follow Conventional Commits (ADR 0011), e.g. 'feat(cli): evaluate -e expressions'" >&2
  echo "got: $subject" >&2
  exit 1
fi
