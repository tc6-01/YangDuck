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

type Config struct {
	Mode          Mode   `yaml:"mode"`
	CachePath     string `yaml:"cache_path"`
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

func (c *Config) SetMode(m Mode) error {
	c.Mode = m
	return c.Save()
}
