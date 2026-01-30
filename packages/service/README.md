# GitHub Browser Service

本地 HTTP 服务，用于克隆和打开 GitHub 仓库和 PR。

## 功能特性

- ✅ 支持 GitHub 仓库 URL
- ✅ 支持 GitHub Pull Request URL
- ✅ 智能缓存管理
- ✅ 支持多种 IDE（VS Code, Zed, IntelliJ IDEA, etc.）
- ✅ 自动处理 PR 分支（包括 fork 的 PR）
- ✅ 支持打开特定文件和行号

## 安装

### 方式 1：使用安装脚本（推荐）

```bash
cd service
./install.sh
```

这将：
1. 构建二进制文件
2. 安装到 `/usr/local/bin`
3. 创建配置文件
4. 设置系统服务（自动启动）

### 方式 2：手动安装

```bash
# 构建
go build -o github-browser-service

# 运行
./github-browser-service
```

## 配置

配置文件位置：`~/.github-browser/config.json`

```json
{
  "port": 9527,
  "defaultIDE": "code",
  "githubToken": "",
  "cacheDir": "/home/user/.github-browser/repos"
}
```

### 配置项说明

- `port`: 服务端口（默认 9527）
- `defaultIDE`: 默认 IDE（code, zed, idea, etc.）
- `githubToken`: GitHub Personal Access Token（可选，用于访问私有仓库和提高 API 限制）
- `cacheDir`: 仓库缓存目录

### 获取 GitHub Token（可选）

如果需要访问私有仓库或提高 API 限制：

1. 访问 https://github.com/settings/tokens
2. 点击 "Generate new token (classic)"
3. 选择权限：`repo`（访问私有仓库）
4. 复制 token 并填入配置文件

## API 接口

### POST /open

打开 GitHub 仓库或 PR。

**请求**：

```json
{
  "url": "https://github.com/microsoft/vscode",
  "ide": "code",
  "filePath": "src/vs/code/electron-main/main.ts",
  "line": 42
}
```

**参数说明**：

- `url` (必需): GitHub URL
  - 仓库: `https://github.com/owner/repo`
  - 文件: `https://github.com/owner/repo/blob/main/file.go`
  - 文件+行号: `https://github.com/owner/repo/blob/main/file.go#L42`
  - PR: `https://github.com/owner/repo/pull/123`
- `ide` (可选): IDE 名称，默认使用配置中的 `defaultIDE`
- `filePath` (可选): 文件路径（相对于仓库根目录）
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

列出所有缓存的仓库。

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

删除指定的缓存仓库。

**示例**：

```bash
curl -X DELETE http://localhost:9527/cache/microsoft-vscode
```

### GET /config

获取当前配置。

### PUT /config

更新配置。

**请求**：

```json
{
  "port": 9527,
  "defaultIDE": "zed",
  "githubToken": "ghp_xxx",
  "cacheDir": "/custom/path"
}
```

## 支持的 IDE

| IDE | 命令 | 行号支持 |
|-----|------|---------|
| VS Code | `code` | ✅ |
| VS Code Insiders | `code-insiders` | ✅ |
| Zed | `zed` | ✅ |
| Cursor | `cursor` | ✅ |
| IntelliJ IDEA | `idea` | ✅ |
| PyCharm | `pycharm` | ✅ |
| WebStorm | `webstorm` | ✅ |
| GoLand | `goland` | ✅ |
| Neovim | `nvim` | ✅ |
| Sublime Text | `subl` | ✅ |

## Pull Request 处理

服务会自动处理 PR：

1. **同仓库的 PR**：直接 checkout PR 分支
2. **Fork 的 PR**：
   - 添加 fork 仓库为 remote
   - Fetch PR 分支
   - 创建本地分支 `pr-{number}`

**示例**：

```bash
curl -X POST http://localhost:9527/open \
  -H "Content-Type: application/json" \
  -d '{"url": "https://github.com/microsoft/vscode/pull/12345"}'
```

服务会：
1. 克隆 `microsoft/vscode`（如果未缓存）
2. 获取 PR 信息
3. Checkout PR 分支
4. 在 IDE 中打开

## 使用示例

### 命令行测试

```bash
# 打开仓库
curl -X POST http://localhost:9527/open \
  -H "Content-Type: application/json" \
  -d '{"url": "https://github.com/microsoft/vscode"}'

# 打开特定文件
curl -X POST http://localhost:9527/open \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://github.com/microsoft/vscode/blob/main/src/vs/code/electron-main/main.ts",
    "ide": "code"
  }'

# 打开 PR
curl -X POST http://localhost:9527/open \
  -H "Content-Type: application/json" \
  -d '{"url": "https://github.com/microsoft/vscode/pull/12345"}'

# 查看缓存
curl http://localhost:9527/cache

# 健康检查
curl http://localhost:9527/health
```

### 从 IDE 插件调用

IDE 插件会调用此服务：

```typescript
// VS Code 插件示例
async function openGitHubRepo(url: string) {
  const response = await fetch('http://localhost:9527/open', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ url })
  });
  
  const result = await response.json();
  console.log(result);
}
```

## 故障排除

### 服务未启动

```bash
# 检查服务状态（Linux）
sudo systemctl status github-browser

# 查看日志（Linux）
sudo journalctl -u github-browser -f

# 查看日志（macOS）
tail -f ~/.github-browser/service.log
```

### Git 克隆失败

确保：
1. Git 已安装
2. 有网络连接
3. 对于私有仓库，需要配置 GitHub Token

### IDE 无法打开

确保：
1. IDE 已安装
2. IDE 命令在 PATH 中
3. 配置文件中的 IDE 名称正确

## 开发

### 运行测试

```bash
go test ./...
```

### 本地开发

```bash
# 运行服务
go run .

# 或使用 air 进行热重载
go install github.com/cosmtrek/air@latest
air
```

## 许可证

MIT
