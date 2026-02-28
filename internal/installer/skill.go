package installer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yangduck/yduck/internal/recipe"
)

type SkillInstaller struct{}

func NewSkillInstaller() *SkillInstaller {
	return &SkillInstaller{}
}

func (s *SkillInstaller) Install(rec *recipe.Recipe) error {
	for _, f := range rec.Files {
		dest := f.Dest
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return fmt.Errorf("create dir for %s: %w", dest, err)
		}
		content := []byte(f.Source)
		if _, err := os.Stat(f.Source); err == nil {
			var readErr error
			content, readErr = os.ReadFile(f.Source)
			if readErr != nil {
				return fmt.Errorf("read source %s: %w", f.Source, readErr)
			}
		}
		if err := os.WriteFile(dest, content, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", dest, err)
		}
	}
	return nil
}

func (s *SkillInstaller) IsInstalled(rec *recipe.Recipe) bool {
	for _, f := range rec.Files {
		if _, err := os.Stat(f.Dest); err != nil {
			return false
		}
	}
	return len(rec.Files) > 0
}
