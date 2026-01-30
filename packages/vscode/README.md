# GitHub Browser for VS Code

一键打开 GitHub 仓库和 PR，支持完整的 LSP 功能。

## 功能特性

- ✅ 从 GitHub URL 快速打开仓库
- ✅ 支持 Pull Request
- ✅ 从剪贴板打开
- ✅ 智能缓存，重复打开速度快
- ✅ 完整的 LSP 支持（代码跳转、智能提示等）
- ✅ 状态栏显示服务状态

## 前置要求

需要先安装并运行 GitHub Browser 服务：

```bash
cd service
./install.sh
```

服务会在后台运行（http://localhost:9527）。

## 使用方法

### 方式 1：从命令面板

1. 按 `Ctrl+Shift+P`（Mac: `Cmd+Shift+P`）
2. 输入 "GitHub Browser"
3. 选择命令：
   - **Open Repository**: 输入 GitHub URL
   - **Open from Clipboard**: 从剪贴板读取 URL
   - **Open Pull Request**: 输入仓库和 PR 号

### 方式 2：使用快捷键

- `Ctrl+Shift+G Ctrl+Shift+O`（Mac: `Cmd+Shift+G Cmd+Shift+O`）
  - 从剪贴板打开 GitHub URL

### 方式 3：从状态栏

点击右下角的 "GitHub Browser" 图标。

## 支持的 URL 格式

- 仓库: `https://github.com/microsoft/vscode`
- 文件: `https://github.com/microsoft/vscode/blob/main/src/vs/code/electron-main/main.ts`
- 文件+行号: `https://github.com/microsoft/vscode/blob/main/src/vs/code/electron-main/main.ts#L42`
- Pull Request: `https://github.com/microsoft/vscode/pull/12345`
- 目录: `https://github.com/microsoft/vscode/tree/main/src/vs`

## 工作流程

1. 在 GitHub 上浏览代码
2. 复制 URL
3. 在 VS Code 中按快捷键 `Cmd+Shift+G Cmd+Shift+O`
4. 自动克隆并打开！

## 配置

### 服务 URL

默认：`http://localhost:9527`

如果服务运行在其他端口，可以修改：

```json
{
  "github-browser.serviceUrl": "http://localhost:9527"
}
```

### 自动检查服务状态

默认：`true`

```json
{
  "github-browser.autoCheckService": true
}
```

## 故障排除

### 服务未运行

如果状态栏显示警告图标：

1. 检查服务是否运行：
   ```bash
   curl http://localhost:9527/health
   ```

2. 启动服务：
   ```bash
   # Linux
   sudo systemctl start github-browser
   
   # macOS
   launchctl load ~/Library/LaunchAgents/com.github-browser.service.plist
   
   # 或手动运行
   github-browser-service
   ```

### 无法连接服务

1. 检查防火墙设置
2. 确认服务端口（默认 9527）
3. 查看服务日志：
   ```bash
   # Linux
   sudo journalctl -u github-browser -f
   
   # macOS
   tail -f ~/.github-browser/service.log
   ```

## 开发

### 构建

```bash
npm install
npm run compile
```

### 调试

1. 在 VS Code 中打开此项目
2. 按 F5 启动调试
3. 在新窗口中测试扩展

### 打包

```bash
npm run package
```

生成 `.vsix` 文件，可以手动安装：

```bash
code --install-extension github-browser-1.0.0.vsix
```

## 许可证

MIT
