#!/bin/bash
# 打包 Chrome 和 Firefox 浏览器扩展
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

VERSION=$(grep '"version"' manifest.json | head -1 | sed 's/.*: *"\(.*\)".*/\1/')
DIST_DIR="$SCRIPT_DIR/dist"
CHROME_DIR="$DIST_DIR/chrome"
FIREFOX_DIR="$DIST_DIR/firefox"

# 需要打包的文件列表
FILES=(
  manifest.json
  background.js
  content.js
  content.css
  popup.html
  popup.js
  options.html
  options.js
  icons/icon16.png
  icons/icon48.png
  icons/icon128.png
)

echo "📦 Packaging GitHub Browser Extension v${VERSION}..."

rm -rf "$DIST_DIR"
mkdir -p "$CHROME_DIR" "$FIREFOX_DIR"

# 复制源文件到两个目录
for f in "${FILES[@]}"; do
  mkdir -p "$CHROME_DIR/$(dirname "$f")" "$FIREFOX_DIR/$(dirname "$f")"
  cp "$f" "$CHROME_DIR/$f"
  cp "$f" "$FIREFOX_DIR/$f"
done

# Chrome: 移除 browser_specific_settings（Chrome 不支持）
python3 -c "
import json, sys
with open('$CHROME_DIR/manifest.json') as f:
    m = json.load(f)
m.pop('browser_specific_settings', None)
with open('$CHROME_DIR/manifest.json', 'w') as f:
    json.dump(m, f, indent=2)
"

# Firefox: 移除 key，将 service_worker 改为 scripts
python3 -c "
import json
with open('$FIREFOX_DIR/manifest.json') as f:
    m = json.load(f)
m.pop('key', None)
if 'background' in m and 'service_worker' in m['background']:
    sw = m['background'].pop('service_worker')
    m['background']['scripts'] = [sw]
with open('$FIREFOX_DIR/manifest.json', 'w') as f:
    json.dump(m, f, indent=2)
"

# 打包 zip
cd "$CHROME_DIR"
zip -r -q "$DIST_DIR/github-browser-chrome-v${VERSION}.zip" .

cd "$FIREFOX_DIR"
zip -r -q "$DIST_DIR/github-browser-firefox-v${VERSION}.zip" .

echo ""
echo "✅ Done!"
echo "   Chrome:  dist/github-browser-chrome-v${VERSION}.zip"
echo "   Firefox: dist/github-browser-firefox-v${VERSION}.zip"
echo ""
echo "Chrome:  加载 dist/chrome/ 目录或上传 zip 到 Chrome Web Store"
echo "Firefox: 上传 zip 到 addons.mozilla.org 或用 about:debugging 加载 dist/firefox/"
