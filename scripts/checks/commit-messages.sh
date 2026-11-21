#!/usr/bin/env bash
set -euo pipefail

range=${1:?usage: commit-messages.sh REVISION_RANGE}
status=0

while IFS= read -r subject; do
  "$(dirname "$0")/commit-subject.sh" "$subject" || status=1
done < <(git log --format=%s "$range")

exit $status
