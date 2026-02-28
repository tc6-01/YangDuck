# 🐤 YangDuck (yduck)

快速配置你的 Mac 开发环境 + AI 编码工具。

一键安装 CLI 工具、配置 MCP 服务器、Cursor Skills/Commands/Rules，面向新手的傻瓜式体验。

## 安装

```bash
curl -fsSL https://yduck.dev/install | sh
```

或手动安装：

```bash
# macOS arm64 (Apple Silicon)
curl -L -o /usr/local/bin/yduck https://github.com/yangduck/yduck/releases/latest/download/yduck-darwin-arm64
chmod +x /usr/local/bin/yduck

# macOS amd64 (Intel)
curl -L -o /usr/local/bin/yduck https://github.com/yangduck/yduck/releases/latest/download/yduck-darwin-amd64
chmod +x /usr/local/bin/yduck
```

## 快速开始

```bash
# 进入交互式界面（推荐新手）
yduck

# 环境检查
yduck doctor

# 安装单个工具
yduck install fzf

# 安装套餐
yduck install --bundle cli-essentials

# 搜索配方
yduck search database

# 列出所有配方
yduck list

# 生成配方
yduck recipe generate fzf
yduck recipe generate --type mcp @some/mcp-package
```

## 功能

### CLI 工具安装

通过 Homebrew 安装常用命令行工具，安装后自动展示使用教程：

- fzf, ripgrep, bat, eza, jq, lazygit, tldr, httpie, fd ...

### AI 工具配置

一键配置 MCP 服务器到 Cursor 和 Claude Code：

- MySQL MCP — 让 AI 直接操作数据库
- Filesystem MCP — 让 AI 读写文件
- Fetch MCP — 让 AI 访问网页和 API

### 套餐安装

按场景组合的工具包，一键安装多个工具：

- **CLI 必备工具包** — fzf, ripgrep, bat, eza, jq, fd, lazygit, tldr
- **AI 入门套餐** — Fetch MCP, Filesystem MCP

### 配方生成器

自动采集工具信息 + AI 生成配方，批量创建无需手写：

```bash
yduck recipe generate fzf                          # 单个 CLI 工具
yduck recipe generate --type mcp @some/mcp-pkg     # 单个 MCP
yduck recipe generate --from-brewfile ~/Brewfile    # 从 Brewfile 批量
yduck recipe generate --from-mcp-config .cursor/mcp.json  # 从 MCP 配置批量
```

## 双模式

- **新手模式**（默认）：逐步引导、术语解释、安装后教程
- **高级模式**：精简输出、批量安装、跳过引导

```bash
yduck config mode advanced   # 切换到高级模式
yduck config mode beginner   # 切换回新手模式
```

## 配方贡献

配方是 YAML 文件，定义了工具的安装方式和使用引导。欢迎贡献！

1. Fork 仓库
2. 在 `recipes/` 下添加 YAML 配方（可用 `yduck recipe generate` 自动生成初稿）
3. 确保符合 `schemas/` 中的 JSON Schema
4. 提交 PR

配方变更合入 main 后会自动发布。

## 开发

```bash
go build ./cmd/yduck/           # 构建
go run ./cmd/yduck/             # 运行
go run scripts/build-index.go   # 构建配方索引
```

## License

MIT
