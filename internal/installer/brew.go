package installer

import (
	"fmt"
	"os/exec"
	"strings"

	ylog "github.com/yangduck/yduck/internal/log"
)

type BrewInstaller struct{}

func NewBrewInstaller() *BrewInstaller {
	return &BrewInstaller{}
}

func (b *BrewInstaller) IsAvailable() bool {
	_, err := exec.LookPath("brew")
	return err == nil
}

func (b *BrewInstaller) IsInstalled(pkg string) (bool, string) {
	out, err := exec.Command("brew", "list", "--versions", pkg).Output()
	if err != nil {
		return false, ""
	}
	version := strings.TrimSpace(string(out))
	if version == "" {
		return false, ""
	}
	parts := strings.Fields(version)
	if len(parts) >= 2 {
		return true, parts[len(parts)-1]
	}
	return true, version
}

func (b *BrewInstaller) Install(pkg string) error {
	ylog.S.Debugw("brew install", "package", pkg)
	cmd := exec.Command("brew", "install", pkg)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		ylog.S.Errorw("brew install failed", "package", pkg, "error", err)
		return err
	}
	ylog.S.Infow("brew install succeeded", "package", pkg)
	return nil
}

func (b *BrewInstaller) Upgrade(pkg string) error {
	cmd := exec.Command("brew", "upgrade", pkg)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func (b *BrewInstaller) RunPostInstall(commands []string) error {
	for _, c := range commands {
		ylog.S.Debugw("running post_install", "command", c)
		cmd := exec.Command("sh", "-c", c)
		if err := cmd.Run(); err != nil {
			ylog.S.Warnw("post_install failed", "command", c, "error", err)
			return fmt.Errorf("post_install %q: %w", c, err)
		}
	}
	return nil
}
