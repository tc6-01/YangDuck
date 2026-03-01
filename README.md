# 🐤 YangDuck (yduck)

快速配置你的 Mac 开发环境 + AI 编码工具。

一键安装 CLI 工具、配置 MCP 服务器、Cursor Skills/Commands/Rules，面向新手的傻瓜式体验。

## 安装

```bash
curl -fsSL https://raw.githubusercontent.com/tc6-01/YangDuck/master/install.sh | sh
```

或手动安装：

```bash
# macOS arm64 (Apple Silicon)
curl -L -o /usr/local/bin/yduck https://github.com/tc6-01/YangDuck/releases/latest/download/yduck-darwin-arm64
chmod +x /usr/local/bin/yduck

# macOS amd64 (Intel)
curl -L -o /usr/local/bin/yduck https://github.com/tc6-01/YangDuck/releases/latest/download/yduck-darwin-amd64
chmod +x /usr/local/bin/yduck
```

## 快速开始

```bash
# 进入交互式界面（推荐新手）
yduck

# 环境检查
yduck doctor
```

## 命令参考

### `yduck`

无参数运行，进入交互式 TUI 界面。新手首次使用会进入引导流程。

### `yduck install`

安装配方（CLI 工具、MCP、Skill、Command、Rule）。

```bash
yduck install fzf                       # 安装单个工具
yduck install fzf ripgrep bat           # 安装多个工具
yduck install --bundle cli-essentials   # 安装套餐
```

### `yduck search`

按关键字搜索配方。

```bash
yduck search database
yduck search mcp
```

### `yduck list`

列出可用配方。

```bash
yduck list                # 列出所有配方
yduck list --installed    # 仅显示已安装的
yduck list --bundles      # 仅显示套餐
```

### `yduck doctor`

检查开发环境状态（Homebrew、Node.js、Cursor、Claude Code）。

### `yduck config`

管理配置。

```bash
yduck config mode advanced    # 切换到高级模式
yduck config mode beginner    # 切换回新手模式
```

### `yduck recipe generate`

自动采集工具信息并生成配方 YAML。

```bash
yduck recipe generate fzf                                    # 生成 CLI 工具配方
yduck recipe generate --type mcp @some/mcp-pkg               # 生成 MCP 配方
yduck recipe generate --from-brewfile ~/Brewfile              # 从 Brewfile 批量生成
yduck recipe generate --from-mcp-config .cursor/mcp.json     # 从 MCP 配置批量生成
```

### `yduck update`

更新配方索引（远程仓库配置后启用）。

## 功能

### CLI 工具安装

通过 Homebrew 安装常用命令行工具，安装后自动展示使用教程（新手模式）：

- fzf, ripgrep, bat, eza, jq, lazygit, tldr, httpie, fd ...

### MCP 服务器配置

一键配置 MCP 服务器到 Cursor 和 Claude Code：

- MySQL MCP — 让 AI 直接操作数据库
- Filesystem MCP — 让 AI 读写文件
- Fetch MCP — 让 AI 访问网页和 API

### Cursor 扩展

支持安装 Cursor 的 Skill、Command 和 Rule 配方。

### 套餐安装

按场景组合的工具包，一键安装多个工具：

- **CLI 必备工具包** — fzf, ripgrep, bat, eza, jq, fd, lazygit, tldr
- **AI 入门套餐** — Fetch MCP, Filesystem MCP

### 配方生成器

自动采集工具信息 + AI 生成配方，批量创建无需手写。支持从 Brewfile 和 MCP 配置文件批量导入。

## 双模式

- **新手模式**（默认）：逐步引导、术语解释、安装后教程
- **高级模式**：精简输出、批量安装、跳过引导

## 配方贡献

配方是 YAML 文件，定义了工具的安装方式和使用引导。欢迎贡献！

1. Fork 仓库
2. 在 `recipes/` 下添加 YAML 配方（可用 `yduck recipe generate` 自动生成初稿）
3. 确保符合 `schemas/` 中的 JSON Schema
4. 提交 PR

配方变更合入 main 后会自动发布。

## 项目结构

```
yduck-cli/
├── cmd/yduck/          # 入口
├── internal/
│   ├── config/         # 配置管理
│   ├── generator/      # 配方生成器
│   ├── installer/      # 安装器（Brew、MCP、Skill、Bundle）
│   ├── quickstart/     # 新手引导
│   ├── recipe/         # 配方加载与校验
│   └── tui/            # 交互式界面
├── recipes/            # 配方源文件
├── schemas/            # JSON Schema
├── scripts/            # 构建脚本
└── install.sh          # 一键安装脚本
```

## 开发

```bash
go build ./cmd/yduck/           # 构建
go run ./cmd/yduck/             # 运行
go run scripts/build-index.go   # 构建配方索引
```

## 致谢

感谢以下开源项目的支持：

- [Cobra](https://github.com/spf13/cobra) — CLI 框架
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI 框架
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — 终端样式
- [Huh](https://github.com/charmbracelet/huh) — 终端表单
- [gojsonschema](https://github.com/xeipuuv/gojsonschema) — JSON Schema 校验

## License

MIT
