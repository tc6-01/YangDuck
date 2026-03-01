package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type AIBackend interface {
	Generate(prompt string) (string, error)
}

func DetectBackend() (AIBackend, error) {
	if _, err := exec.LookPath("claude"); err == nil {
		return &ClaudeCLI{}, nil
	}
	if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		return &AnthropicAPI{Key: key}, nil
	}
	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		return &OpenAIAPI{Key: key}, nil
	}
	return nil, fmt.Errorf("未找到可用的 AI 后端。请安装 Claude CLI (npm install -g @anthropic-ai/claude-code) 或设置 ANTHROPIC_API_KEY / OPENAI_API_KEY 环境变量")
}

type ClaudeCLI struct{}

func (c *ClaudeCLI) Generate(prompt string) (string, error) {
	cmd := exec.Command("claude", "--print", prompt)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("claude cli: %w", err)
	}
	return string(out), nil
}

type AnthropicAPI struct{ Key string }

func (a *AnthropicAPI) Generate(prompt string) (string, error) {
	body := map[string]interface{}{
		"model":      "claude-sonnet-4-20250514",
		"max_tokens": 2000,
		"messages":   []map[string]string{{"role": "user", "content": prompt}},
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", a.Key)
	req.Header.Set("anthropic-version", "2023-06-01")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result struct {
		Content []struct{ Text string } `json:"content"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	if len(result.Content) == 0 {
		return "", fmt.Errorf("empty response")
	}
	return result.Content[0].Text, nil
}

type OpenAIAPI struct{ Key string }

func (o *OpenAIAPI) Generate(prompt string) (string, error) {
	body := map[string]interface{}{
		"model":      "gpt-4o",
		"max_tokens": 2000,
		"messages":   []map[string]string{{"role": "user", "content": prompt}},
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.Key)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result struct {
		Choices []struct {
			Message struct{ Content string } `json:"message"`
		} `json:"choices"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("empty response")
	}
	return result.Choices[0].Message.Content, nil
}

func BuildPrompt(recipeType string, info *ToolInfo) string {
	var sb strings.Builder
	sb.WriteString("请根据以下工具信息生成一个 YAML 配方文件。\n\n")
	sb.WriteString("要求：\n")
	sb.WriteString("- 使用中文描述，面向编程新手\n")
	sb.WriteString("- description 简洁明了，一句话说明这个工具能干什么\n")
	sb.WriteString("- quickstart 提供 2-4 个最实用的使用示例\n")
	sb.WriteString("- 只输出 YAML 内容，不要包含 ``` 标记\n\n")

	sb.WriteString(fmt.Sprintf("工具名称: %s\n", info.Name))
	sb.WriteString(fmt.Sprintf("官方描述: %s\n", info.Description))
	if info.Homepage != "" {
		sb.WriteString(fmt.Sprintf("官网: %s\n", info.Homepage))
	}

	switch recipeType {
	case "cli-tool":
		sb.WriteString(fmt.Sprintf("\n配方类型: cli-tool\n"))
		sb.WriteString("安装方式: brew\n")
		sb.WriteString(fmt.Sprintf("包名: %s\n", info.Name))
		if len(info.Examples) > 0 {
			sb.WriteString("\ntldr 示例:\n")
			for _, e := range info.Examples {
				sb.WriteString("- " + e + "\n")
			}
		}
		sb.WriteString("\nYAML schema:\n")
		sb.WriteString(`id: <kebab-case-name>
name: <name>
type: cli-tool
description: <中文一句话描述>
tags: [<相关标签>]
difficulty: beginner
install:
  method: brew
  package: <包名>
quickstart:
  - title: <使用场景标题>
    command: "<示例命令>"
    explain: <用中文解释这个命令做什么>
`)
	case "mcp":
		sb.WriteString(fmt.Sprintf("\n配方类型: mcp\n"))
		if len(info.EnvVars) > 0 {
			sb.WriteString("发现的环境变量: " + strings.Join(info.EnvVars, ", ") + "\n")
		}
		if len(info.ReadmeParts) > 0 {
			sb.WriteString("\nREADME 摘要:\n" + info.ReadmeParts[0] + "\n")
		}
		sb.WriteString("\nYAML schema:\n")
		sb.WriteString(`id: <kebab-case-name>
name: <name>
type: mcp
description: <中文一句话描述，说明安装后AI能干什么>
tags: [<相关标签>]
difficulty: beginner
install:
  method: npm
  package: "<npm包名>"
targets:
  cursor:
    config_path: ".cursor/mcp.json"
    config:
      mcpServers:
        <server-name>:
          command: "npx"
          args: ["-y", "<npm包名>"]
          env: <如果需要环境变量>
  claude-code:
    config_path: "~/.claude.json"
    config:
      mcpServers:
        <server-name>:
          command: "npx"
          args: ["-y", "<npm包名>"]
          env: <如果需要环境变量，和 cursor 保持一致>
prompts: <如果有需要用户填写的配置>
  - key: <环境变量名>
    ask: "<用中文问用户>"
    default: "<默认值>"
quickstart:
  - title: 试试看
    explain: <用中文告诉用户安装后怎么用>
`)
		sb.WriteString("\n注意：targets 必须同时包含 cursor 和 claude-code 两个目标。config_path 固定为上述值，config 下必须有 mcpServers 字段。\n")
	}
	return sb.String()
}
