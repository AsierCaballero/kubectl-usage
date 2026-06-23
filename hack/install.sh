#!/bin/bash
set -euo pipefail

VERSION="${1:-latest}"
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported arch: $ARCH"; exit 1 ;;
esac

if [ "$VERSION" = "latest" ]; then
  VERSION=$(curl -s https://api.github.com/repos/AsierCaballero/kubectl-usage/releases/latest | grep tag_name | cut -d'"' -f4)
fi

URL="https://github.com/AsierCaballero/kubectl-usage/releases/download/${VERSION}/kubectl-usage_${VERSION}_${OS}_${ARCH}.tar.gz"

echo "Downloading kubectl-usage ${VERSION} (${OS}/${ARCH})..."
curl -sSL "$URL" | tar xz -C /usr/local/bin kubectl-usage
chmod +x /usr/local/bin/kubectl-usage

echo "Installed kubectl-usage ${VERSION}"
echo "Run: kubectl usage -n default"
