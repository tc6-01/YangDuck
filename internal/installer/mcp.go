package installer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yangduck/yduck/internal/recipe"
)

type MCPInstaller struct{}

func NewMCPInstaller() *MCPInstaller {
	return &MCPInstaller{}
}

func (m *MCPInstaller) Install(rec *recipe.Recipe, target string, promptValues map[string]string) error {
	var tc *recipe.TargetConfig
	switch target {
	case "cursor":
		if rec.Targets == nil || rec.Targets.Cursor == nil {
			return fmt.Errorf("recipe %s has no cursor target", rec.ID)
		}
		tc = rec.Targets.Cursor
	case "claude-code":
		if rec.Targets == nil || rec.Targets.ClaudeCode == nil {
			return fmt.Errorf("recipe %s has no claude-code target", rec.ID)
		}
		tc = rec.Targets.ClaudeCode
	default:
		return fmt.Errorf("unknown target: %s", target)
	}

	configPath := expandPath(tc.ConfigPath)
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	existing := make(map[string]interface{})
	if data, err := os.ReadFile(configPath); err == nil {
		_ = json.Unmarshal(data, &existing)
		backup := configPath + ".backup"
		_ = os.WriteFile(backup, data, 0o644)
	}

	newConfig := substitutePrompts(tc.Config, promptValues)
	merged := mergeConfig(existing, newConfig)

	data, err := json.MarshalIndent(merged, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	return os.WriteFile(configPath, data, 0o644)
}

func (m *MCPInstaller) IsConfigured(rec *recipe.Recipe, target string) bool {
	var tc *recipe.TargetConfig
	switch target {
	case "cursor":
		if rec.Targets == nil || rec.Targets.Cursor == nil {
			return false
		}
		tc = rec.Targets.Cursor
	case "claude-code":
		if rec.Targets == nil || rec.Targets.ClaudeCode == nil {
			return false
		}
		tc = rec.Targets.ClaudeCode
	default:
		return false
	}
	configPath := expandPath(tc.ConfigPath)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return false
	}
	var existing map[string]interface{}
	if err := json.Unmarshal(data, &existing); err != nil {
		return false
	}
	servers, ok := existing["mcpServers"]
	if !ok {
		return false
	}
	serversMap, ok := servers.(map[string]interface{})
	if !ok {
		return false
	}
	mcpServers, ok := tc.Config["mcpServers"]
	if !ok {
		return false
	}
	mcpServersMap, ok := mcpServers.(map[string]interface{})
	if !ok {
		return false
	}
	for name := range mcpServersMap {
		if _, exists := serversMap[name]; exists {
			return true
		}
	}
	return false
}

func expandPath(p string) string {
	if strings.HasPrefix(p, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, p[2:])
	}
	return p
}

func substitutePrompts(config map[string]interface{}, values map[string]string) map[string]interface{} {
	data, _ := json.Marshal(config)
	s := string(data)
	for k, v := range values {
		s = strings.ReplaceAll(s, k, v)
	}
	var result map[string]interface{}
	_ = json.Unmarshal([]byte(s), &result)
	return result
}

func mergeConfig(base, overlay map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range base {
		result[k] = v
	}
	for k, v := range overlay {
		if baseVal, ok := result[k]; ok {
			baseMap, baseOk := baseVal.(map[string]interface{})
			overlayMap, overlayOk := v.(map[string]interface{})
			if baseOk && overlayOk {
				result[k] = mergeConfig(baseMap, overlayMap)
				continue
			}
		}
		result[k] = v
	}
	return result
}
