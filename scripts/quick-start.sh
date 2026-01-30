#!/bin/bash

set -e

echo "🚀 GitHub Browser - Quick Start"
echo "================================"
echo ""

# 检测操作系统
OS="$(uname -s)"
case "${OS}" in
    Linux*)     MACHINE=linux;;
    Darwin*)    MACHINE=macos;;
    *)          echo "❌ Unsupported OS: ${OS}"; exit 1;;
esac

echo "📋 Detected OS: $MACHINE"
echo ""

# 步骤 1: 安装服务
echo "Step 1/3: Installing service..."
cd service
./install.sh
cd ..

echo ""
echo "✅ Service installed!"
echo ""

# 步骤 2: 测试服务
echo "Step 2/3: Testing service..."
sleep 2

if curl -s http://localhost:9527/health > /dev/null; then
    echo "✅ Service is running!"
else
    echo "❌ Service is not running. Please check the logs."
    exit 1
fi

echo ""

# 步骤 3: 安装 VS Code 插件（可选）
echo "Step 3/3: VS Code plugin (optional)"
read -p "Do you want to install VS Code plugin? (y/n) " -n 1 -r
echo ""

if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "Installing VS Code plugin..."
    cd vscode-plugin
    
    if command -v pnpm &> /dev/null; then
        pnpm install
        pnpm run compile
        pnpm run package
        
        if [ -f "github-browser-1.0.0.vsix" ]; then
            echo ""
            echo "✅ VS Code plugin built!"
            echo ""
            echo "To install:"
            echo "  code --install-extension github-browser-1.0.0.vsix"
        fi
    else
        echo "❌ pnpm not found. Please install pnpm first:"
        echo "  npm install -g pnpm"
    fi
    
    cd ..
fi

echo ""
echo "🎉 Setup complete!"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📝 Next steps:"
echo ""
echo "1. Test the service:"
echo "   curl -X POST http://localhost:9527/open \\"
echo "     -H 'Content-Type: application/json' \\"
echo "     -d '{\"url\": \"https://github.com/golang/go\", \"ide\": \"code\"}'"
echo ""
echo "2. Install browser extension:"
echo "   - Chrome: chrome://extensions/ → Load unpacked → browser-ext/"
echo "   - Firefox: about:debugging → Load Temporary Add-on → browser-ext/manifest.json"
echo ""
echo "3. Read the full guide:"
echo "   cat GUIDE.md"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "🌐 Service URL: http://localhost:9527"
echo "📁 Cache directory: ~/.github-browser/repos"
echo "⚙️  Config file: ~/.github-browser/config.json"
echo ""
echo "Enjoy! 🎊"
