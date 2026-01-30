# GitHub Browser - Browser Extension

一键在 GitHub 页面打开仓库和 PR 到本地 IDE。

## 功能特性

- ✅ 在 GitHub 页面添加 "Open in IDE" 按钮
- ✅ 支持仓库、文件、PR 等所有页面
- ✅ 键盘快捷键：`Shift+O`
- ✅ 从剪贴板打开
- ✅ 自动检测服务状态
- ✅ 支持多种 IDE

## 安装

### Chrome/Edge

1. 下载扩展文件
2. 打开 `chrome://extensions/`
3. 启用"开发者模式"
4. 点击"加载已解压的扩展程序"
5. 选择 `browser-ext` 目录

### Firefox

1. 下载扩展文件
2. 打开 `about:debugging#/runtime/this-firefox`
3. 点击"临时载入附加组件"
4. 选择 `manifest.json` 文件

## 前置要求

需要先安装并运行 GitHub Browser 服务：

```bash
cd service
./install.sh
```

## 使用方法

### 方式 1：点击按钮

在 GitHub 页面上，你会看到 "Open in IDE" 按钮：

- **仓库主页**：按钮在侧边栏
- **文件页面**：按钮在文件操作栏
- **PR 页面**：按钮在标题旁边

### 方式 2：键盘快捷键

在任何 GitHub 页面按 `Shift+O` 即可打开当前页面。

### 方式 3：扩展图标

点击浏览器工具栏的扩展图标：

- **Open Current Page**：打开当前 GitHub 页面
- **Open from Clipboard**：从剪贴板读取 URL 并打开
- **Settings**：配置服务 URL 和默认 IDE

## 支持的页面类型

- ✅ 仓库主页：`https://github.com/microsoft/vscode`
- ✅ 文件：`https://github.com/microsoft/vscode/blob/main/src/vs/code/electron-main/main.ts`
- ✅ 文件+行号：`https://github.com/microsoft/vscode/blob/main/src/vs/code/electron-main/main.ts#L42`
- ✅ Pull Request：`https://github.com/microsoft/vscode/pull/12345`
- ✅ 目录：`https://github.com/microsoft/vscode/tree/main/src/vs`

## 配置

点击扩展图标 → Settings，配置：

### 服务 URL

默认：`http://localhost:9527`

如果服务运行在其他端口，修改此设置。

### 默认 IDE

选择你想使用的 IDE：

- VS Code
- VS Code Insiders
- Zed
- Cursor
- IntelliJ IDEA
- PyCharm
- WebStorm
- GoLand
- Neovim
- Sublime Text

## 工作流程示例

### 场景 1：Review PR

1. 在 GitHub 上打开 PR
2. 点击 "Open in IDE" 按钮（或按 `Shift+O`）
3. 自动克隆仓库并 checkout PR 分支
4. 在 IDE 中查看代码，完整的 LSP 支持！

### 场景 2：快速查看代码

1. 在 GitHub 上浏览代码
2. 看到感兴趣的文件
3. 按 `Shift+O`
4. 在 IDE 中打开，可以跳转定义、查找引用等

### 场景 3：从其他地方打开

1. 在 Slack/Email 等地方看到 GitHub 链接
2. 复制链接
3. 点击扩展图标 → "Open from Clipboard"
4. 自动打开！

## 故障排除

### 服务未运行

如果扩展显示 "Service not running"：

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

### 按钮未显示

1. 刷新 GitHub 页面
2. 检查是否在支持的页面类型上
3. 查看浏览器控制台是否有错误

### CORS 错误

如果看到 CORS 错误，确保服务配置正确：

```bash
curl http://localhost:9527/health
```

应该返回 JSON 响应。

## 开发

### 修改代码

修改代码后，在扩展管理页面点击"重新加载"。

### 调试

1. 打开 GitHub 页面
2. 按 F12 打开开发者工具
3. 查看 Console 标签页的日志

### 打包发布

```bash
# 创建 zip 文件
cd browser-ext
zip -r github-browser-extension.zip * -x "*.git*" "*.DS_Store"
```

## 隐私说明

此扩展：

- ✅ 只在 GitHub 页面运行
- ✅ 不收集任何数据
- ✅ 不发送数据到外部服务器
- ✅ 只与本地服务通信

## 许可证

MIT
