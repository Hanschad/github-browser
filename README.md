# GitHub Browser

一键打开 GitHub 仓库和 PR 到本地 IDE，支持完整的 LSP 功能。

## ✨ 功能特性

- ✅ **支持 GitHub 仓库和 Pull Request**
- ✅ **完整的 LSP 支持**（代码跳转、智能提示、查找引用）
- ✅ **智能缓存管理**（重复打开速度快）
- ✅ **支持多种 IDE**（VS Code, Zed, IntelliJ IDEA, etc.）
- ✅ **自动处理 PR 分支**（包括 fork 的 PR）
- ✅ **多种使用方式**（浏览器扩展、IDE 插件、命令行）

## 🚀 快速开始

### 1. 安装本地服务

```bash
cd packages/service
./install.sh
```

### 2. 验证服务

```bash
curl http://localhost:9527/health
```

### 3. 选择客户端

#### 选项 A：VS Code 插件（推荐）

```bash
cd packages/vscode
pnpm install
pnpm run compile
# 在 VS Code 中按 F5 启动调试
```

#### 选项 B：浏览器扩展

1. 打开 `chrome://extensions/`
2. 启用"开发者模式"
3. 点击"加载已解压的扩展程序"
4. 选择 `packages/browser-ext` 目录

#### 选项 C：命令行（最简单）

```bash
curl -s -X POST http://localhost:9527/open \
  -H "Content-Type: application/json" \
  -d '{"url": "https://github.com/microsoft/vscode", "ide": "code"}'
```

## 📖 文档

- **[完整使用指南](docs/GUIDE.md)** - 详细的安装、配置和使用说明
- **[服务文档](packages/service/README.md)** - 本地服务 API 文档
- **[VS Code 插件](packages/vscode/README.md)** - VS Code 插件使用说明
- **[Zed 插件](packages/zed/README.md)** - Zed 编辑器集成
- **[浏览器扩展](packages/browser-ext/README.md)** - Chrome/Firefox 扩展使用说明

## 🎯 使用场景

### 场景 1：Review Pull Request

```
1. 在 GitHub 上打开 PR
2. 点击 "Open in IDE" 按钮（或按 Shift+O）
3. 自动克隆并 checkout PR 分支
4. 在 IDE 中查看，完整的 LSP 支持！
```

### 场景 2：快速查看代码

```
1. 在 GitHub 上浏览代码
2. 看到感兴趣的文件
3. 按 Shift+O
4. 在 IDE 中打开，可以跳转定义、查找引用
```

### 场景 3：从链接快速打开

```
1. 在 Slack/Email 中看到 GitHub 链接
2. 复制链接
3. 在 VS Code 中按 Cmd+Shift+G Cmd+Shift+O
4. 自动打开！
```

## 🏗️ 架构

```
┌─────────────────────────────────────────┐
│  浏览器扩展 / IDE 插件 / 命令行          │
└──────────────────┬──────────────────────┘
                   │ HTTP
                   ▼
┌─────────────────────────────────────────┐
│   本地服务 (localhost:9527)              │
│   - 解析 GitHub URL (repo/PR)            │
│   - 克隆/更新仓库                         │
│   - 处理 PR (checkout branch)            │
│   - 启动 IDE                             │
└─────────────────────────────────────────┘
```

## 📦 项目结构

```
github-browser/
├── packages/
│   ├── service/          # 本地服务 (Go)
│   │   ├── main.go
│   │   ├── github.go     # GitHub API 处理
│   │   ├── git.go        # Git 操作
│   │   ├── ide.go        # IDE 启动
│   │   └── config.go     # 配置管理
│   ├── vscode/           # VS Code 插件 (TypeScript)
│   │   ├── src/
│   │   │   └── extension.ts
│   │   └── package.json
│   ├── zed/              # Zed 插件 (Rust)
│   │   ├── src/
│   │   │   └── lib.rs
│   │   └── Cargo.toml
│   └── browser-ext/      # 浏览器扩展 (JavaScript)
│       ├── content.js    # 内容脚本
│       ├── popup.js      # 弹出窗口
│       └── manifest.json
├── docs/
│   └── GUIDE.md          # 完整使用指南
└── scripts/
    ├── quick-start.sh
    └── examples.sh
```

## 🎨 支持的 URL 格式

| 类型 | 示例 |
|------|------|
| 仓库 | `https://github.com/microsoft/vscode` |
| 文件 | `https://github.com/microsoft/vscode/blob/main/src/vs/code/electron-main/main.ts` |
| 文件+行号 | `https://github.com/microsoft/vscode/blob/main/src/vs/code/electron-main/main.ts#L42` |
| Pull Request | `https://github.com/microsoft/vscode/pull/12345` |
| 目录 | `https://github.com/microsoft/vscode/tree/main/src/vs` |

## 💻 支持的 IDE

- VS Code / VS Code Insiders
- Zed
- Cursor
- IntelliJ IDEA
- PyCharm
- WebStorm
- GoLand
- Neovim
- Sublime Text

## 🔧 配置

配置文件位置：`~/.github-browser/config.json`

```json
{
  "port": 9527,
  "defaultIDE": "code",
  "githubToken": "",
  "cacheDir": "/home/user/.github-browser/repos"
}
```

### 获取 GitHub Token（可选）

用于访问私有仓库和提高 API 限制：

1. 访问 https://github.com/settings/tokens
2. 点击 "Generate new token (classic)"
3. 选择权限：`repo`
4. 复制 token 并填入配置文件

## 🐛 故障排除

### 服务未启动

```bash
# 检查服务状态
curl http://localhost:9527/health

# Linux
sudo systemctl status github-browser

# macOS
tail -f ~/.github-browser/service.log
```

### Git 克隆失败

- 检查 Git 是否安装
- 检查网络连接
- 对于私有仓库，配置 GitHub Token

### IDE 无法打开

- 检查 IDE 命令是否在 PATH 中
- 确认配置文件中的 IDE 名称正确

详细的故障排除指南请参考 [完整使用指南](docs/GUIDE.md)。

## 📊 API 文档

### POST /open

```bash
curl -X POST http://localhost:9527/open \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://github.com/microsoft/vscode/pull/12345",
    "ide": "code"
  }'
```

### GET /health

```bash
curl http://localhost:9527/health
```

### GET /cache

```bash
curl http://localhost:9527/cache
```

### DELETE /cache/:repo

```bash
curl -X DELETE http://localhost:9527/cache/microsoft-vscode
```

完整的 API 文档请参考 [服务文档](packages/service/README.md)。

## 🚦 开发状态

- [x] Phase 1: 本地服务（支持 PR）
- [x] Phase 2: VS Code 插件
- [x] Phase 3: Zed 插件
- [x] Phase 4: 浏览器扩展
- [x] Phase 5: 集成测试和文档

## 🤝 贡献

欢迎贡献！请提交 Issue 或 Pull Request。

## 📄 许可证

MIT

## 🙏 致谢

感谢所有开源项目的贡献者！

---

**开始使用**：阅读 [完整使用指南](docs/GUIDE.md) 了解详细信息。
