#!/bin/bash
# 安装 Native Messaging Host（同时支持 Chrome 系浏览器和 Firefox）
set -euo pipefail

HOST_NAME="com.github.browser"
CHROME_EXTENSION_ID="leikfbanhflejfmlnejhigjbfpaknljd"
FIREFOX_EXTENSION_ID="github-browser@example.com"
INSTALL_ROOT="$HOME/.github-browser"
BIN_DIR="$INSTALL_ROOT/bin"
BINARY_PATH="$BIN_DIR/github-browser-service"

echo "🔧 Installing GitHub Browser native host..."

OS="$(uname -s)"
case "${OS}" in
    Linux*)   MACHINE=linux ;;
    Darwin*)  MACHINE=macos ;;
    *)        echo "❌ Unsupported OS: ${OS}"; exit 1 ;;
esac

# 编译二进制
mkdir -p "$BIN_DIR"
echo "📦 Building native host binary..."
go build -o "$BINARY_PATH"
chmod +x "$BINARY_PATH"

# 初始化配置
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

# ── Chrome 系浏览器 ──────────────────────────────────
declare -a CHROME_DIRS=()
if [ "$MACHINE" = "macos" ]; then
    CHROME_DIRS=(
        "$HOME/Library/Application Support/Google/Chrome/NativeMessagingHosts"
        "$HOME/Library/Application Support/Microsoft Edge/NativeMessagingHosts"
        "$HOME/Library/Application Support/Chromium/NativeMessagingHosts"
    )
else
    CHROME_DIRS=(
        "$HOME/.config/google-chrome/NativeMessagingHosts"
        "$HOME/.config/microsoft-edge/NativeMessagingHosts"
        "$HOME/.config/chromium/NativeMessagingHosts"
    )
fi

echo "📝 Registering for Chrome/Edge/Chromium..."
for dir in "${CHROME_DIRS[@]}"; do
    mkdir -p "$dir"
    cat > "$dir/$HOST_NAME.json" <<EOF
{
  "name": "$HOST_NAME",
  "description": "GitHub Browser Native Host",
  "path": "$BINARY_PATH",
  "type": "stdio",
  "allowed_origins": [
    "chrome-extension://$CHROME_EXTENSION_ID/"
  ]
}
EOF
done

# ── Firefox ──────────────────────────────────────────
declare -a FIREFOX_DIRS=()
if [ "$MACHINE" = "macos" ]; then
    FIREFOX_DIRS=(
        "$HOME/Library/Application Support/Mozilla/NativeMessagingHosts"
    )
else
    FIREFOX_DIRS=(
        "$HOME/.mozilla/native-messaging-hosts"
    )
fi

echo "📝 Registering for Firefox..."
for dir in "${FIREFOX_DIRS[@]}"; do
    mkdir -p "$dir"
    cat > "$dir/$HOST_NAME.json" <<EOF
{
  "name": "$HOST_NAME",
  "description": "GitHub Browser Native Host",
  "path": "$BINARY_PATH",
  "type": "stdio",
  "allowed_extensions": [
    "$FIREFOX_EXTENSION_ID"
  ]
}
EOF
done

echo ""
echo "✅ Native host installed!"
echo "   Binary:    $BINARY_PATH"
echo "   Host name: $HOST_NAME"
echo ""
echo "   Chrome extension ID:  $CHROME_EXTENSION_ID"
echo "   Firefox extension ID: $FIREFOX_EXTENSION_ID"
