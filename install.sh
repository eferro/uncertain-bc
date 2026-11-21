#!/bin/sh
# Installs ubc on Linux and macOS:
#   curl -fsSL https://raw.githubusercontent.com/eferro/uncertain-bc/main/install.sh | sh
# UBC_VERSION pins a release (default: latest); UBC_INSTALL_DIR picks the
# directory (default: /usr/local/bin if writable, else ~/.local/bin).
# POSIX sh, not bash, so it also runs where sh is dash (ADR 0026).
set -eu

repo=eferro/uncertain-bc
github=${UBC_GITHUB:-https://github.com}

fail() {
  echo "install.sh: $*" >&2
  exit 1
}

detect_os() {
  case $(uname -s) in
    Linux) echo linux ;;
    Darwin) echo darwin ;;
    *) fail "unsupported system $(uname -s); on Windows use install.ps1" ;;
  esac
}

detect_arch() {
  case $(uname -m) in
    x86_64 | amd64) echo amd64 ;;
    arm64 | aarch64) echo arm64 ;;
    *) fail "unsupported architecture $(uname -m)" ;;
  esac
}

latest_version() {
  url=$(curl -fsSLI -o /dev/null -w '%{url_effective}' "$github/$repo/releases/latest") ||
    fail "cannot reach $github to find the latest release"
  echo "${url##*/v}"
}

sha256_of() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | cut -d ' ' -f 1
  else
    shasum -a 256 "$1" | cut -d ' ' -f 1
  fi
}

default_dir() {
  if [ -w /usr/local/bin ]; then echo /usr/local/bin; else echo "$HOME/.local/bin"; fi
}

version=${UBC_VERSION:-$(latest_version)}
version=${version#v}
dir=${UBC_INSTALL_DIR:-$(default_dir)}
archive="ubc_${version}_$(detect_os)_$(detect_arch).tar.gz"
base="$github/$repo/releases/download/v$version"

work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT

curl -fsSL -o "$work/$archive" "$base/$archive" || fail "cannot download $base/$archive"
curl -fsSL -o "$work/checksums.txt" "$base/checksums.txt" || fail "cannot download the checksums"
expected=$(grep " $archive\$" "$work/checksums.txt" | cut -d ' ' -f 1)
[ -n "$expected" ] || fail "no checksum for $archive"
[ "$(sha256_of "$work/$archive")" = "$expected" ] || fail "checksum mismatch for $archive"

tar -xzf "$work/$archive" -C "$work" ubc
mkdir -p "$dir"
cp "$work/ubc" "$dir/ubc"
chmod 0755 "$dir/ubc"

echo "ubc $version installed to $dir/ubc"
case ":$PATH:" in
  *":$dir:"*) ;;
  *) echo "Add it to your PATH: export PATH=\"$dir:\$PATH\"" ;;
esac
