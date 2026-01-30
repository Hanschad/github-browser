# GitHub Browser for Zed

一键打开 GitHub 仓库和 PR，支持完整的 LSP 功能。

## 功能特性

- ✅ 从 GitHub URL 快速打开仓库
- ✅ 支持 Pull Request
- ✅ 智能缓存，重复打开速度快
- ✅ 完整的 LSP 支持（代码跳转、智能提示等）

## 前置要求

需要先安装并运行 GitHub Browser 服务：

```bash
cd service
./install.sh
```

服务会在后台运行（http://localhost:9527）。

## 安装

### 方式 1：从 Zed Extensions

1. 打开 Zed
2. 按 `Cmd+Shift+P`（Mac）或 `Ctrl+Shift+P`（Linux/Windows）
3. 输入 "Extensions"
4. 搜索 "GitHub Browser"
5. 点击 "Install"

### 方式 2：手动安装

```bash
cd zed-plugin
cargo build --release
cp target/release/libgithub_browser.dylib ~/.config/zed/extensions/github-browser/
```

## 使用方法

### 从命令面板

1. 复制 GitHub URL
2. 在 Zed 中按 `Cmd+Shift+P`
3. 输入 "GitHub Browser: Open from Clipboard"
4. 自动克隆并打开！

### 支持的 URL 格式

- 仓库: `https://github.com/microsoft/vscode`
- 文件: `https://github.com/microsoft/vscode/blob/main/src/vs/code/electron-main/main.ts`
- 文件+行号: `https://github.com/microsoft/vscode/blob/main/src/vs/code/electron-main/main.ts#L42`
- Pull Request: `https://github.com/microsoft/vscode/pull/12345`
- 目录: `https://github.com/microsoft/vscode/tree/main/src/vs`

## 配置

在 Zed 设置中配置服务 URL：

```json
{
  "github-browser": {
    "service_url": "http://localhost:9527"
  }
}
```

## 故障排除

### 服务未运行

检查服务是否运行：

```bash
curl http://localhost:9527/health
```

启动服务：

```bash
# Linux
sudo systemctl start github-browser

# macOS
launchctl load ~/Library/LaunchAgents/com.github-browser.service.plist

# 或手动运行
github-browser-service
```

## 开发

### 构建

```bash
cargo build --release
```

### 测试

```bash
cargo test
```

## 许可证

MIT

## 注意

由于 Zed 扩展 API 的限制，当前版本的功能可能不如 VS Code 版本完整。
我们建议主要使用 VS Code 插件或直接使用命令行工具配合浏览器扩展。

### 替代方案：使用命令行工具

创建一个简单的 shell 脚本 `~/.local/bin/gho`：

```bash
#!/bin/bash
# GitHub Open for Zed

URL="$1"
if [ -z "$URL" ]; then
  URL=$(pbpaste)  # macOS
  # URL=$(xclip -o)  # Linux
fi

curl -s -X POST http://localhost:9527/open \
  -H "Content-Type: application/json" \
  -d "{\"url\": \"$URL\", \"ide\": \"zed\"}"
```

使用方式：

```bash
# 从参数
gho https://github.com/microsoft/vscode

# 从剪贴板
gho
```

配合 Alfred/Raycast 等工具可以实现快速打开。
