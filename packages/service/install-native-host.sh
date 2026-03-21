#!/bin/bash

set -euo pipefail

HOST_NAME="com.github.browser"
EXTENSION_ID="leikfbanhflejfmlnejhigjbfpaknljd"
INSTALL_ROOT="$HOME/.github-browser"
BIN_DIR="$INSTALL_ROOT/bin"
BINARY_PATH="$BIN_DIR/github-browser-service"
MANIFEST_FILE="$HOST_NAME.json"

echo "Installing GitHub Browser native host..."

OS="$(uname -s)"
case "${OS}" in
    Linux*)   MACHINE=linux ;;
    Darwin*)  MACHINE=macos ;;
    *)        echo "Unsupported OS: ${OS}"; exit 1 ;;
esac

mkdir -p "$BIN_DIR"

echo "Building native host binary..."
go build -o "$BINARY_PATH"
chmod +x "$BINARY_PATH"

mkdir -p "$INSTALL_ROOT"
if [ ! -f "$INSTALL_ROOT/config.json" ]; then
    cat > "$INSTALL_ROOT/config.json" <<EOF
{
  "port": 9527,
  "defaultIDE": "code",
  "githubToken": "",
  "cacheDir": "$HOME/.github-browser/repos"
}
EOF
fi

declare -a HOST_DIRS=()
if [ "$MACHINE" = "macos" ]; then
    HOST_DIRS+=(
        "$HOME/Library/Application Support/Google/Chrome/NativeMessagingHosts"
        "$HOME/Library/Application Support/Microsoft Edge/NativeMessagingHosts"
        "$HOME/Library/Application Support/Chromium/NativeMessagingHosts"
    )
else
    HOST_DIRS+=(
        "$HOME/.config/google-chrome/NativeMessagingHosts"
        "$HOME/.config/microsoft-edge/NativeMessagingHosts"
        "$HOME/.config/chromium/NativeMessagingHosts"
    )
fi

echo "Writing native host manifests..."
for host_dir in "${HOST_DIRS[@]}"; do
    mkdir -p "$host_dir"
    cat > "$host_dir/$MANIFEST_FILE" <<EOF
{
  "name": "$HOST_NAME",
  "description": "GitHub Browser Native Host",
  "path": "$BINARY_PATH",
  "type": "stdio",
  "allowed_origins": [
    "chrome-extension://$EXTENSION_ID/"
  ]
}
EOF
done

echo
echo "Native host installed."
echo "Binary: $BINARY_PATH"
echo "Expected extension ID: $EXTENSION_ID"
echo "Native host name: $HOST_NAME"
echo
echo "Next steps:"
echo "1. Load packages/browser-ext as an unpacked extension in Chrome or Edge."
echo "2. Confirm the extension ID is $EXTENSION_ID."
echo "3. Use the extension normally; it will talk to the native host first and fall back to HTTP if needed."
