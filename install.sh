#!/usr/bin/env sh
set -eu

REPO="kaihendry/ai-check-guardrails"
BINARY="ai-check-guardrails"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
  linux|darwin) ;;
  *) echo "Unsupported OS: $OS" >&2; exit 1 ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
  x86_64)        ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

# Fetch latest release tag
echo "Fetching latest release..."
TAG=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
  | grep '"tag_name"' \
  | sed 's/.*"tag_name": *"\(.*\)".*/\1/')

if [ -z "$TAG" ]; then
  echo "Failed to fetch latest release tag" >&2
  exit 1
fi

NAME="${BINARY}-${OS}-${ARCH}"
BASE_URL="https://github.com/${REPO}/releases/download/${TAG}"

echo "Downloading ${BINARY} ${TAG} (${OS}/${ARCH})..."
curl -fsSL "${BASE_URL}/${NAME}" -o "/tmp/${BINARY}"
curl -fsSL "${BASE_URL}/${NAME}.sha256" -o "/tmp/${BINARY}.sha256"

# Verify checksum
echo "Verifying checksum..."
EXPECTED=$(cat "/tmp/${BINARY}.sha256" | awk '{print $1}')
case "$OS" in
  linux)
    ACTUAL=$(sha256sum "/tmp/${BINARY}" | awk '{print $1}')
    ;;
  darwin)
    ACTUAL=$(shasum -a 256 "/tmp/${BINARY}" | awk '{print $1}')
    ;;
esac

if [ "$EXPECTED" != "$ACTUAL" ]; then
  echo "Checksum mismatch — download may be corrupted" >&2
  rm -f "/tmp/${BINARY}" "/tmp/${BINARY}.sha256"
  exit 1
fi

mkdir -p "${INSTALL_DIR}"
install -m 755 "/tmp/${BINARY}" "${INSTALL_DIR}/${BINARY}"
rm -f "/tmp/${BINARY}" "/tmp/${BINARY}.sha256"

echo "Installed ${BINARY} ${TAG} to ${INSTALL_DIR}/${BINARY}"

if ! echo ":$PATH:" | grep -q ":${INSTALL_DIR}:"; then
  echo "Note: ${INSTALL_DIR} is not in your PATH. Add it with:"
  echo "  export PATH=\"\$PATH:${INSTALL_DIR}\""
fi
