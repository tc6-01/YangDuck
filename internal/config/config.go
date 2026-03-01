package config

import (
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

type Mode string

const (
	ModeBeginner Mode = "beginner"
	ModeAdvanced Mode = "advanced"
)

type Editor string

const (
	EditorCursor    Editor = "cursor"
	EditorClaudeCode Editor = "claude-code"
	EditorBoth      Editor = "both"
)

type Config struct {
	Mode           Mode   `yaml:"mode"`
	Editor         Editor `yaml:"editor,omitempty"`
	CachePath      string `yaml:"cache_path"`
	LastUpdateCheck string `yaml:"last_update_check,omitempty"`
}

var (
	cfg     *Config
	cfgOnce sync.Once
	cfgDir  string
)

func Dir() string {
	if cfgDir != "" {
		return cfgDir
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".yduck")
}

func CacheDir() string {
	return filepath.Join(Dir(), "cache")
}

func filePath() string {
	return filepath.Join(Dir(), "config.yaml")
}

func Load() *Config {
	cfgOnce.Do(func() {
		cfg = &Config{
			Mode:      ModeBeginner,
			CachePath: CacheDir(),
		}
		data, err := os.ReadFile(filePath())
		if err != nil {
			return
		}
		_ = yaml.Unmarshal(data, cfg)
	})
	return cfg
}

func (c *Config) Save() error {
	if err := os.MkdirAll(Dir(), 0o755); err != nil {
		return err
	}
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath(), data, 0o644)
}

func (c *Config) IsBeginner() bool {
	return c.Mode != ModeAdvanced
}

func (c *Config) IsFirstTime() bool {
	_, err := os.Stat(filePath())
	return os.IsNotExist(err)
}

func (c *Config) SetMode(m Mode) error {
	c.Mode = m
	return c.Save()
}

func (c *Config) SetEditor(e Editor) error {
	c.Editor = e
	return c.Save()
}

func (c *Config) ShouldInstallFor(target string) bool {
	switch c.Editor {
	case EditorCursor:
		return target == "cursor"
	case EditorClaudeCode:
		return target == "claude-code"
	case EditorBoth, "":
		return true
	}
	return true
}
