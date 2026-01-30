#!/bin/bash

set -e

echo "🚀 Installing GitHub Browser Service..."

# 检测操作系统
OS="$(uname -s)"
case "${OS}" in
    Linux*)     MACHINE=linux;;
    Darwin*)    MACHINE=macos;;
    *)          echo "Unsupported OS: ${OS}"; exit 1;;
esac

# 构建二进制文件
echo "📦 Building service..."
go build -o github-browser-service

# 安装到 /usr/local/bin
echo "📥 Installing to /usr/local/bin..."
sudo cp github-browser-service /usr/local/bin/

# 创建配置目录
echo "📁 Creating config directory..."
mkdir -p ~/.github-browser

# 创建默认配置
if [ ! -f ~/.github-browser/config.json ]; then
    echo "⚙️  Creating default config..."
    cat > ~/.github-browser/config.json <<EOF
{
  "port": 9527,
  "defaultIDE": "code",
  "githubToken": "",
  "cacheDir": "$HOME/.github-browser/repos"
}
EOF
fi

# 创建 systemd 服务（Linux）
if [ "$MACHINE" = "linux" ]; then
    echo "🔧 Creating systemd service..."
    sudo tee /etc/systemd/system/github-browser.service > /dev/null <<EOF
[Unit]
Description=GitHub Browser Service
After=network.target

[Service]
Type=simple
User=$USER
ExecStart=/usr/local/bin/github-browser-service
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

    echo "🔄 Enabling and starting service..."
    sudo systemctl daemon-reload
    sudo systemctl enable github-browser
    sudo systemctl start github-browser
    
    echo "✅ Service installed and started!"
    echo "📊 Check status: sudo systemctl status github-browser"
fi

# 创建 LaunchAgent（macOS）
if [ "$MACHINE" = "macos" ]; then
    echo "🔧 Creating LaunchAgent..."
    mkdir -p ~/Library/LaunchAgents
    cat > ~/Library/LaunchAgents/com.github-browser.service.plist <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.github-browser.service</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/github-browser-service</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>$HOME/.github-browser/service.log</string>
    <key>StandardErrorPath</key>
    <string>$HOME/.github-browser/service.error.log</string>
</dict>
</plist>
EOF

    echo "🔄 Loading LaunchAgent..."
    launchctl load ~/Library/LaunchAgents/com.github-browser.service.plist
    
    echo "✅ Service installed and started!"
    echo "📊 Check logs: tail -f ~/.github-browser/service.log"
fi

echo ""
echo "🎉 Installation complete!"
echo ""
echo "📝 Configuration file: ~/.github-browser/config.json"
echo "🌐 Service URL: http://localhost:9527"
echo "💡 Test: curl http://localhost:9527/health"
echo ""
echo "Next steps:"
echo "1. Install IDE plugins (VS Code, Zed)"
echo "2. Install browser extension"
echo "3. Start browsing GitHub repos!"
