package generator

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yangduck/yduck/internal/recipe"
	"gopkg.in/yaml.v3"
)

type Generator struct {
	ai        AIBackend
	outputDir string
}

func New(outputDir string) (*Generator, error) {
	ai, err := DetectBackend()
	if err != nil {
		return nil, err
	}
	return &Generator{ai: ai, outputDir: outputDir}, nil
}

func (g *Generator) GenerateCLITool(pkg string) (*recipe.Recipe, error) {
	info, err := CollectBrewInfo(pkg)
	if err != nil {
		return nil, fmt.Errorf("采集信息失败: %w", err)
	}

	prompt := BuildPrompt("cli-tool", info)
	yamlStr, err := g.ai.Generate(prompt)
	if err != nil {
		return nil, fmt.Errorf("AI 生成失败: %w", err)
	}

	rec, err := parseYAMLRecipe(yamlStr)
	if err != nil {
		retryPrompt := prompt + "\n\n上次生成的 YAML 有格式错误: " + err.Error() + "\n请修正后重新生成。"
		yamlStr, err = g.ai.Generate(retryPrompt)
		if err != nil {
			return nil, fmt.Errorf("AI 重试失败: %w", err)
		}
		rec, err = parseYAMLRecipe(yamlStr)
		if err != nil {
			return nil, fmt.Errorf("YAML 解析仍然失败: %w", err)
		}
	}

	outPath := filepath.Join(g.outputDir, "cli-tools", rec.ID+".yaml")
	if err := writeRecipe(outPath, yamlStr); err != nil {
		return nil, err
	}
	return rec, nil
}

func (g *Generator) GenerateMCP(pkg string) (*recipe.Recipe, error) {
	info, err := CollectNPMInfo(pkg)
	if err != nil {
		return nil, fmt.Errorf("采集信息失败: %w", err)
	}

	prompt := BuildPrompt("mcp", info)
	yamlStr, err := g.ai.Generate(prompt)
	if err != nil {
		return nil, fmt.Errorf("AI 生成失败: %w", err)
	}

	rec, err := parseYAMLRecipe(yamlStr)
	if err != nil {
		retryPrompt := prompt + "\n\n上次生成的 YAML 有格式错误: " + err.Error() + "\n请修正后重新生成。"
		yamlStr, err = g.ai.Generate(retryPrompt)
		if err != nil {
			return nil, fmt.Errorf("AI 重试失败: %w", err)
		}
		rec, err = parseYAMLRecipe(yamlStr)
		if err != nil {
			return nil, fmt.Errorf("YAML 解析仍然失败: %w", err)
		}
	}

	outPath := filepath.Join(g.outputDir, "mcps", rec.ID+".yaml")
	if err := writeRecipe(outPath, yamlStr); err != nil {
		return nil, err
	}
	return rec, nil
}

func (g *Generator) GenerateFromBrewfile(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var generated []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "brew ") {
			pkg := strings.Trim(strings.TrimPrefix(line, "brew "), "\"' ")
			pkg = strings.Split(pkg, ",")[0]
			pkg = strings.TrimSpace(pkg)
			if pkg == "" {
				continue
			}
			fmt.Printf("正在生成 %s...\n", pkg)
			rec, err := g.GenerateCLITool(pkg)
			if err != nil {
				fmt.Printf("⚠ %s 跳过: %v\n", pkg, err)
				continue
			}
			generated = append(generated, rec.ID)
			fmt.Printf("✓ %s.yaml\n", rec.ID)
		}
	}
	return generated, nil
}

func (g *Generator) GenerateFromMCPConfig(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config struct {
		MCPServers map[string]struct {
			Command string   `json:"command"`
			Args    []string `json:"args"`
		} `json:"mcpServers"`
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	var generated []string
	for name, srv := range config.MCPServers {
		pkg := ""
		for _, arg := range srv.Args {
			if arg != "-y" && !strings.HasPrefix(arg, "-") {
				pkg = arg
				break
			}
		}
		if pkg == "" {
			pkg = name
		}
		fmt.Printf("正在生成 %s...\n", name)
		rec, err := g.GenerateMCP(pkg)
		if err != nil {
			fmt.Printf("⚠ %s 跳过: %v\n", name, err)
			continue
		}
		generated = append(generated, rec.ID)
		fmt.Printf("✓ %s.yaml\n", rec.ID)
	}
	return generated, nil
}

func parseYAMLRecipe(yamlStr string) (*recipe.Recipe, error) {
	yamlStr = strings.TrimPrefix(yamlStr, "```yaml\n")
	yamlStr = strings.TrimPrefix(yamlStr, "```\n")
	yamlStr = strings.TrimSuffix(yamlStr, "\n```")
	yamlStr = strings.TrimSuffix(yamlStr, "```")
	yamlStr = strings.TrimSpace(yamlStr)

	var rec recipe.Recipe
	if err := yaml.Unmarshal([]byte(yamlStr), &rec); err != nil {
		return nil, err
	}
	return &rec, nil
}

func writeRecipe(path, content string) error {
	content = strings.TrimPrefix(content, "```yaml\n")
	content = strings.TrimPrefix(content, "```\n")
	content = strings.TrimSuffix(content, "\n```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content+"\n"), 0o644)
}
