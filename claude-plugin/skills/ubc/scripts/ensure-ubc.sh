#!/bin/sh
# Checks that ubc 1.1.1 or later is available for the ubc skill (ADR 0068):
#   ensure-ubc.sh            exit 0 and print "ubc X ready: PATH", or exit 1
#   ensure-ubc.sh --install  install the latest release first when ubc is
#                            missing or older
# It installs with the project's install.sh, which verifies the checksum, into
# UBC_INSTALL_DIR (default ~/.local/bin), so it never needs sudo.
# POSIX sh, like install.sh, so it also runs where sh is dash.
set -eu

minimum=1.1.1
installer=${UBC_INSTALL_SCRIPT:-https://raw.githubusercontent.com/eferro/uncertain-bc/main/install.sh}
dir=${UBC_INSTALL_DIR:-$HOME/.local/bin}

usage() {
  echo "usage: ensure-ubc.sh [--install]" >&2
  exit 2
}

# A release prints "ubc 1.1.1"; a build from a clone prints
# "ubc v1.1.1-3-gabc1234", counted as the release it starts from.
version_of() {
  "$1" --version 2>/dev/null |
    sed -n 's/^ubc v\{0,1\}\([0-9][0-9]*\.[0-9][0-9]*\.[0-9][0-9]*\).*/\1/p'
}

at_least_minimum() {
  lowest=$(printf '%s\n%s\n' "$minimum" "$1" | sort -t . -k 1,1n -k 2,2n -k 3,3n | head -n 1)
  [ "$lowest" = "$minimum" ]
}

# check BINARY: report it when it is recent enough, explain why not otherwise.
check() {
  if [ -z "$1" ]; then
    echo "ubc is not installed; run: ensure-ubc.sh --install" >&2
    return 1
  fi
  version=$(version_of "$1")
  if [ -z "$version" ] || ! at_least_minimum "$version"; then
    echo "ubc at $1 is version ${version:-unknown}, older than $minimum; run: ensure-ubc.sh --install" >&2
    return 1
  fi
  echo "ubc $version ready: $1"
}

on_path() {
  command -v ubc 2>/dev/null || true
}

in_dir() {
  if [ -x "$dir/ubc" ]; then echo "$dir/ubc"; fi
}

# An earlier --install may have left ubc in a directory that is not on PATH.
installed() {
  check "$(on_path)" 2>/dev/null || check "$(in_dir)" 2>/dev/null || {
    found=$(on_path)
    check "${found:-$(in_dir)}"
  }
}

install_latest() {
  echo "Installing the latest ubc release into $dir ..."
  script=$(curl -fsSL "$installer") || {
    echo "ensure-ubc.sh: cannot download $installer" >&2
    exit 1
  }
  printf '%s\n' "$script" | UBC_INSTALL_DIR="$dir" sh
  check "$dir/ubc"
  case $(on_path) in
    "$dir/ubc") ;;
    "") echo "$dir is not on PATH: call ubc as $dir/ubc" ;;
    *) echo "another ubc comes first on PATH ($(on_path)): call ubc as $dir/ubc" ;;
  esac
}

case $# in
  0) installed ;;
  1)
    [ "$1" = "--install" ] || usage
    installed 2>/dev/null || install_latest
    ;;
  *) usage ;;
esac
