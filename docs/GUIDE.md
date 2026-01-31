# GitHub Browser - 完整使用指南

## 项目概述

GitHub Browser 是一个完整的解决方案，让你可以从 GitHub 网页一键打开仓库和 PR 到本地 IDE，支持完整的 LSP 功能（代码跳转、智能提示等）。

### 架构

```
┌─────────────────────────────────────────┐
│  浏览器扩展 / IDE 插件 / 命令行          │
└──────────────────┬──────────────────────┘
                   │ HTTP
                   ▼
┌─────────────────────────────────────────┐
│   本地服务 (localhost:9527)              │
│   - 解析 GitHub URL                      │
│   - 克隆/更新仓库                         │
│   - 处理 PR 分支                         │
│   - 启动 IDE                             │
└─────────────────────────────────────────┘
```

### 组件

1. **本地服务**（Go）：核心服务，处理 Git 操作和 IDE 启动
2. **VS Code 插件**（TypeScript）：VS Code 集成
3. **Zed 插件**（Rust）：Zed 编辑器集成
4. **浏览器扩展**（JavaScript）：Chrome/Firefox 扩展

---

## 快速开始

### 第一步：安装本地服务

```bash
cd service
./install.sh
```

这将：
- 构建 Go 二进制文件
- 安装到 `/usr/local/bin`
- 创建配置文件 `~/.github-browser/config.json`
- 设置系统服务（自动启动）

### 第二步：验证服务

```bash
# 检查服务状态
curl http://localhost:9527/health

# 应该返回：
# {"status":"ok","version":"1.0.0","uptime":"..."}
```

### 第三步：选择客户端

根据你的工作流程选择一个或多个客户端：

#### 选项 A：VS Code 插件（推荐）

```bash
cd vscode-plugin
pnpm install
pnpm run compile

# 在 VS Code 中：
# 1. 按 F5 启动调试
# 2. 或打包安装：pnpm run package
```

#### 选项 B：浏览器扩展

```bash
# Chrome/Edge:
# 1. 打开 chrome://extensions/
# 2. 启用"开发者模式"
# 3. 点击"加载已解压的扩展程序"
# 4. 选择 browser-ext 目录
```

#### 选项 C：命令行（最简单）

创建 `/usr/local/bin/gho`：

```bash
#!/bin/bash
curl -s -X POST http://localhost:9527/open \
  -H "Content-Type: application/json" \
  -d "{\"url\": \"$1\", \"ide\": \"code\"}"
```

使用：
```bash
gho https://github.com/microsoft/vscode
```

---

## 详细安装指南

### 1. 本地服务安装

#### 系统要求

- Go 1.21+ （安装脚本会自动安装）
- Git
- 支持的操作系统：Linux, macOS

#### 手动安装

```bash
cd service

# 构建
go build -o github-browser-service

# 运行
./github-browser-service
```

#### 配置

编辑 `~/.github-browser/config.json`：

```json
{
  "port": 9527,
  "defaultIDE": "code",
  "githubToken": "",
  "cacheDir": "/home/user/.github-browser/repos",
  "pathMappings": [
    { "pattern": "microsoft", "localPath": "~/projects/microsoft" },
    { "pattern": "my-org/my-repo", "localPath": "~/work/my-repo" },
    { "pattern": "*", "localPath": "~/github" }
  ]
}
```

**配置项说明**：

- `port`: 服务端口
- `defaultIDE`: 默认 IDE（code, zed, idea, etc.）
- `githubToken`: GitHub Personal Access Token（可选）
  - 用于访问私有仓库
  - 提高 API 限制
  - 获取方式：https://github.com/settings/tokens
- `cacheDir`: 默认仓库缓存目录
- `pathMappings`: 路径映射规则（可选）
  - 将 GitHub owner/repo 映射到本地目录
  - 支持三种匹配模式（按优先级排序）：
    1. `owner/repo` - 精确匹配特定仓库
    2. `owner` - 匹配该用户/组织下的所有仓库
    3. `*` - 通配符，匹配所有其他仓库

#### 系统服务管理

**Linux (systemd)**:

```bash
# 启动
sudo systemctl start github-browser

# 停止
sudo systemctl stop github-browser

# 查看状态
sudo systemctl status github-browser

# 查看日志
sudo journalctl -u github-browser -f
```

**macOS (LaunchAgent)**:

```bash
# 启动
launchctl load ~/Library/LaunchAgents/com.github-browser.service.plist

# 停止
launchctl unload ~/Library/LaunchAgents/com.github-browser.service.plist

# 查看日志
tail -f ~/.github-browser/service.log
```

---

### 2. VS Code 插件安装

#### 开发模式

```bash
cd vscode-plugin
pnpm install
pnpm run compile

# 在 VS Code 中按 F5 启动调试
```

#### 打包安装

```bash
cd vscode-plugin
pnpm run package

# 生成 github-browser-1.0.0.vsix
# 安装：
code --install-extension github-browser-1.0.0.vsix
```

#### 使用方法

1. **从命令面板**：
   - 按 `Cmd+Shift+P` (Mac) 或 `Ctrl+Shift+P` (Windows/Linux)
   - 输入 "GitHub Browser"
   - 选择命令

2. **快捷键**：
   - `Cmd+Shift+G Cmd+Shift+O` (Mac)
   - `Ctrl+Shift+G Ctrl+Shift+O` (Windows/Linux)
   - 从剪贴板打开 GitHub URL

3. **状态栏**：
   - 点击右下角的 "GitHub Browser" 图标

---

### 3. Zed 插件安装

#### 注意

由于 Zed 扩展 API 的限制，推荐使用命令行工具配合 Zed。

#### 命令行工具

创建 `~/.local/bin/gho-zed`：

```bash
#!/bin/bash
URL="$1"
if [ -z "$URL" ]; then
  URL=$(pbpaste)  # macOS
  # URL=$(xclip -o)  # Linux
fi

curl -s -X POST http://localhost:9527/open \
  -H "Content-Type: application/json" \
  -d "{\"url\": \"$URL\", \"ide\": \"zed\"}"
```

使用：
```bash
chmod +x ~/.local/bin/gho-zed
gho-zed https://github.com/microsoft/vscode
```

---

### 4. 浏览器扩展安装

#### Chrome/Edge

1. 打开 `chrome://extensions/`
2. 启用"开发者模式"（右上角）
3. 点击"加载已解压的扩展程序"
4. 选择 `browser-ext` 目录

#### Firefox

1. 打开 `about:debugging#/runtime/this-firefox`
2. 点击"临时载入附加组件"
3. 选择 `browser-ext/manifest.json` 文件

#### 使用方法

1. **页面按钮**：
   - 在 GitHub 页面会自动添加 "Open in IDE" 按钮

2. **键盘快捷键**：
   - 在 GitHub 页面按 `Shift+O`

3. **扩展图标**：
   - 点击工具栏的扩展图标
   - 选择 "Open Current Page" 或 "Open from Clipboard"

4. **配置**：
   - 点击扩展图标 → Settings
   - 配置服务 URL 和默认 IDE

---

## 使用场景

### 场景 1：Review Pull Request

```bash
# 1. 在 GitHub 上打开 PR
https://github.com/microsoft/vscode/pull/12345

# 2. 点击 "Open in IDE" 按钮（或按 Shift+O）

# 3. 服务会：
#    - 克隆仓库（如果未缓存）
#    - 获取 PR 信息
#    - Checkout PR 分支
#    - 在 IDE 中打开

# 4. 在 IDE 中：
#    - 完整的 LSP 支持
#    - 代码跳转、查找引用
#    - 运行测试、调试
```

### 场景 2：快速查看代码

```bash
# 1. 在 GitHub 上浏览代码
https://github.com/golang/go/blob/master/src/runtime/proc.go#L123

# 2. 按 Shift+O

# 3. 服务会：
#    - 克隆仓库
#    - 在 IDE 中打开文件
#    - 跳转到第 123 行

# 4. 在 IDE 中：
#    - 查看函数定义
#    - 查找所有引用
#    - 理解代码结构
```

### 场景 3：从链接快速打开

```bash
# 1. 在 Slack/Email 中看到 GitHub 链接
# 2. 复制链接
# 3. 在 VS Code 中按 Cmd+Shift+G Cmd+Shift+O
# 4. 自动打开！
```

### 场景 4：批量打开多个仓库

```bash
# 使用命令行工具
for repo in microsoft/vscode golang/go rust-lang/rust; do
  gho "https://github.com/$repo"
done
```

---

## API 文档

### POST /open

打开 GitHub 仓库或 PR。

**请求**：

```json
{
  "url": "https://github.com/microsoft/vscode/pull/12345",
  "ide": "code",
  "filePath": "src/vs/code/electron-main/main.ts",
  "line": 42
}
```

**参数**：

- `url` (必需): GitHub URL
- `ide` (可选): IDE 名称
- `filePath` (可选): 文件路径
- `line` (可选): 行号

**响应**：

```json
{
  "status": "ok",
  "message": "Opened successfully",
  "path": "/home/user/.github-browser/repos/microsoft-vscode"
}
```

### GET /health

健康检查。

**响应**：

```json
{
  "status": "ok",
  "version": "1.0.0",
  "uptime": "1h30m"
}
```

### GET /cache

列出缓存的仓库。

**响应**：

```json
{
  "repos": [
    {
      "name": "microsoft-vscode",
      "path": "/home/user/.github-browser/repos/microsoft-vscode",
      "modified": "2024-01-28T10:30:00Z"
    }
  ],
  "count": 1
}
```

### DELETE /cache/:repo

删除缓存的仓库。

```bash
curl -X DELETE http://localhost:9527/cache/microsoft-vscode
```

---

## 支持的 URL 格式

| 类型 | 格式 | 示例 |
|------|------|------|
| 仓库 | `github.com/owner/repo` | `https://github.com/microsoft/vscode` |
| 文件 | `github.com/owner/repo/blob/branch/path` | `https://github.com/microsoft/vscode/blob/main/src/vs/code/electron-main/main.ts` |
| 文件+行号 | `github.com/owner/repo/blob/branch/path#L123` | `https://github.com/microsoft/vscode/blob/main/src/vs/code/electron-main/main.ts#L123` |
| PR | `github.com/owner/repo/pull/123` | `https://github.com/microsoft/vscode/pull/12345` |
| 目录 | `github.com/owner/repo/tree/branch/path` | `https://github.com/microsoft/vscode/tree/main/src/vs` |

---

## 支持的 IDE

| IDE | 命令 | 行号支持 | 备注 |
|-----|------|---------|------|
| VS Code | `code` | ✅ | 推荐 |
| VS Code Insiders | `code-insiders` | ✅ | |
| Zed | `zed` | ✅ | |
| Cursor | `cursor` | ✅ | |
| IntelliJ IDEA | `idea` | ✅ | |
| PyCharm | `pycharm` | ✅ | |
| WebStorm | `webstorm` | ✅ | |
| GoLand | `goland` | ✅ | |
| Neovim | `nvim` | ✅ | macOS 会在新终端打开 |
| Sublime Text | `subl` | ✅ | |

---

## Pull Request 处理详解

### 同仓库的 PR

```bash
# URL: https://github.com/microsoft/vscode/pull/12345
# PR 分支在同一个仓库

服务会：
1. 克隆 microsoft/vscode
2. 获取 PR 信息（分支名：feature-branch）
3. git checkout feature-branch
4. 在 IDE 中打开
```

### Fork 的 PR

```bash
# URL: https://github.com/microsoft/vscode/pull/12345
# PR 来自 fork: user/vscode

服务会：
1. 克隆 microsoft/vscode
2. 获取 PR 信息
   - Head: user/vscode:feature-branch
   - Base: microsoft/vscode:main
3. 添加 remote: git remote add pr-12345 https://github.com/user/vscode.git
4. Fetch 分支: git fetch pr-12345 feature-branch
5. Checkout: git checkout -b pr-12345 pr-12345/feature-branch
6. 在 IDE 中打开
```

---

## 故障排除

### 服务无法启动

**症状**：`curl http://localhost:9527/health` 失败

**解决**：

1. 检查端口是否被占用：
   ```bash
   lsof -i :9527
   ```

2. 查看日志：
   ```bash
   # Linux
   sudo journalctl -u github-browser -f
   
   # macOS
   tail -f ~/.github-browser/service.log
   ```

3. 手动运行查看错误：
   ```bash
   cd service
   ./github-browser-service
   ```

### Git 克隆失败

**症状**：服务返回 "git clone failed"

**解决**：

1. 检查 Git 是否安装：
   ```bash
   git --version
   ```

2. 检查网络连接：
   ```bash
   ping github.com
   ```

3. 对于私有仓库，配置 GitHub Token：
   ```json
   {
     "githubToken": "YOUR_GITHUB_TOKEN_HERE"
   }
   ```

### IDE 无法打开

**症状**：服务返回成功，但 IDE 没有打开

**解决**：

1. 检查 IDE 命令是否在 PATH 中：
   ```bash
   which code  # VS Code
   which zed   # Zed
   ```

2. 手动测试 IDE 命令：
   ```bash
   code /path/to/repo
   ```

3. 检查配置文件中的 IDE 名称是否正确

### VS Code 插件无法连接服务

**症状**：插件显示 "Service not running"

**解决**：

1. 检查服务是否运行：
   ```bash
   curl http://localhost:9527/health
   ```

2. 检查插件配置：
   - VS Code 设置 → 搜索 "github-browser"
   - 确认 Service URL 正确

3. 查看 VS Code 输出：
   - View → Output → 选择 "GitHub Browser"

### 浏览器扩展无法连接服务

**症状**：点击按钮后显示 "Cannot connect to service"

**解决**：

1. 检查服务是否运行

2. 检查 CORS 设置：
   - 服务已经配置了 CORS，应该可以正常工作

3. 查看浏览器控制台：
   - F12 → Console
   - 查看错误信息

### PR 分支 checkout 失败

**症状**：服务返回 "failed to checkout PR branch"

**解决**：

1. 手动测试 Git 操作：
   ```bash
   cd ~/.github-browser/repos/owner-repo
   git fetch origin pull/123/head:pr-123
   git checkout pr-123
   ```

2. 对于私有仓库，确保配置了 GitHub Token

---

## 高级配置

### 路径映射（Path Mappings）

通过 `pathMappings` 配置项，可以将不同的 GitHub 用户/组织/仓库映射到不同的本地目录。

**配置示例**：

```json
{
  "pathMappings": [
    { "pattern": "my-company", "localPath": "~/work" },
    { "pattern": "my-company/important-repo", "localPath": "~/work/important" },
    { "pattern": "microsoft", "localPath": "~/opensource/microsoft" },
    { "pattern": "*", "localPath": "~/github" }
  ]
}
```

**匹配规则**（按优先级从高到低）：

| 模式 | 说明 | 示例 |
|------|------|------|
| `owner/repo` | 精确匹配特定仓库 | `microsoft/vscode` → `~/work/vscode` |
| `owner` | 匹配该用户/组织下的所有仓库 | `microsoft` → `~/opensource/microsoft/vscode` |
| `*` | 通配符，匹配所有其他仓库 | 任意仓库 → `~/github/owner-repo` |

**实际效果示例**：

使用上述配置，各 URL 会被克隆到以下位置：

| GitHub URL | 本地路径 |
|------------|---------|
| `github.com/my-company/important-repo` | `~/work/important` |
| `github.com/my-company/other-repo` | `~/work/other-repo` |
| `github.com/microsoft/vscode` | `~/opensource/microsoft/vscode` |
| `github.com/torvalds/linux` | `~/github/torvalds-linux` |

**通过浏览器扩展配置**：

1. 点击浏览器扩展图标 → Settings
2. 在 "Path Mappings" 区域添加映射规则
3. 点击 "Save Settings"

> **注意**：修改配置后需要重启服务才能生效。

### 自定义缓存目录

```json
{
  "cacheDir": "/mnt/ssd/github-cache"
}
```

### 自定义端口

```json
{
  "port": 8080
}
```

然后更新客户端配置：

- VS Code: 设置 → `github-browser.serviceUrl` → `http://localhost:8080`
- 浏览器扩展: 扩展设置 → Service URL → `http://localhost:8080`

### 使用 GitHub Token

获取 Token：

1. 访问 https://github.com/settings/tokens
2. 点击 "Generate new token (classic)"
3. 选择权限：
   - `repo`（访问私有仓库）
   - `read:org`（访问组织私有仓库）
4. 复制 token

配置：

```json
{
  "githubToken": "YOUR_GITHUB_TOKEN_HERE"
}
```

重启服务：

```bash
sudo systemctl restart github-browser  # Linux
launchctl unload ~/Library/LaunchAgents/com.github-browser.service.plist && \
launchctl load ~/Library/LaunchAgents/com.github-browser.service.plist  # macOS
```

---

## 性能优化

### Shallow Clone

服务默认使用 shallow clone：

```bash
git clone --depth=1 --filter=blob:none --single-branch
```

这大大减少了克隆时间和磁盘占用。

### 缓存管理

定期清理旧缓存：

```bash
# 查看缓存
curl http://localhost:9527/cache

# 删除特定仓库
curl -X DELETE http://localhost:9527/cache/microsoft-vscode

# 删除所有超过 30 天未访问的缓存
find ~/.github-browser/repos -type d -mtime +30 -exec rm -rf {} \;
```

---

## 开发指南

### 项目结构

```
github-browser/
├── service/          # 本地服务（Go）
│   ├── main.go
│   ├── github.go     # GitHub API
│   ├── git.go        # Git 操作
│   ├── ide.go        # IDE 启动
│   └── config.go     # 配置管理
├── vscode-plugin/    # VS Code 插件（TypeScript）
│   ├── src/
│   │   └── extension.ts
│   └── package.json
├── zed-plugin/       # Zed 插件（Rust）
│   ├── src/
│   │   └── lib.rs
│   └── Cargo.toml
└── browser-ext/      # 浏览器扩展（JavaScript）
    ├── content.js
    ├── popup.js
    └── manifest.json
```

### 开发服务

```bash
cd service

# 运行
go run .

# 或使用 air 热重载
go install github.com/cosmtrek/air@latest
air
```

### 开发 VS Code 插件

```bash
cd vscode-plugin

# 安装依赖
pnpm install

# 编译
pnpm run compile

# 监听文件变化
pnpm run watch

# 在 VS Code 中按 F5 启动调试
```

### 开发浏览器扩展

```bash
cd browser-ext

# 修改代码后，在扩展管理页面点击"重新加载"
```

---

## 许可证

MIT

---

## 贡献

欢迎贡献！请提交 Issue 或 Pull Request。

---

## 致谢

感谢所有开源项目的贡献者！
