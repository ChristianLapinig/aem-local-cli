#!/bin/sh
set -e

REPO="ChristianLapinig/aem-local-cli"
BINARY="aemlocal"
INSTALL_DIR="/usr/local/bin"

# Determine OS
OS="$(uname -s)"
case "$OS" in
  Linux)  OS="linux" ;;
  Darwin) OS="darwin" ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

# Determine architecture
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Resolve version: use argument if provided, otherwise fetch latest
if [ -n "$1" ]; then
  VERSION="$1"
else
  VERSION="$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/')"
  if [ -z "$VERSION" ]; then
    echo "Failed to fetch latest version"
    exit 1
  fi
fi

ARCHIVE="${BINARY}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE}"

echo "Installing ${BINARY} ${VERSION} (${OS}/${ARCH})..."

TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

CHECKSUMS_URL="https://github.com/${REPO}/releases/download/${VERSION}/checksums.txt"

curl -fsSL "$URL" -o "${TMP}/${ARCHIVE}"
curl -fsSL "$CHECKSUMS_URL" -o "${TMP}/checksums.txt"

# Verify checksum — supports sha256sum (Linux) and shasum (macOS)
EXPECTED="$(grep "  ${ARCHIVE}$" "${TMP}/checksums.txt" | awk '{print $1}')"
if [ -z "$EXPECTED" ]; then
  echo "Checksum for ${ARCHIVE} not found in checksums.txt"
  exit 1
fi

if command -v sha256sum > /dev/null 2>&1; then
  ACTUAL="$(sha256sum "${TMP}/${ARCHIVE}" | awk '{print $1}')"
elif command -v shasum > /dev/null 2>&1; then
  ACTUAL="$(shasum -a 256 "${TMP}/${ARCHIVE}" | awk '{print $1}')"
else
  echo "No sha256 utility found (sha256sum or shasum required)"
  exit 1
fi

if [ "$ACTUAL" != "$EXPECTED" ]; then
  echo "Checksum mismatch for ${ARCHIVE}"
  echo "  expected: ${EXPECTED}"
  echo "  actual:   ${ACTUAL}"
  exit 1
fi

tar -xzf "${TMP}/${ARCHIVE}" -C "$TMP"

# Write to a temp location first, then move (avoids partial writes to PATH)
install -m 755 "${TMP}/${BINARY}" "${TMP}/${BINARY}.install"

if [ -w "$INSTALL_DIR" ]; then
  mv "${TMP}/${BINARY}.install" "${INSTALL_DIR}/${BINARY}"
else
  sudo mv "${TMP}/${BINARY}.install" "${INSTALL_DIR}/${BINARY}"
fi

echo "${BINARY} installed to ${INSTALL_DIR}/${BINARY}"
${BINARY} --version
