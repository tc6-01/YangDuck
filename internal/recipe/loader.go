package recipe

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func LoadFromFS(fsys fs.FS, root string) ([]Recipe, error) {
	var recipes []Recipe
	err := fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".yaml") {
			return err
		}
		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		var r Recipe
		if err := yaml.Unmarshal(data, &r); err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		if r.ID == "" {
			base := filepath.Base(path)
			r.ID = strings.TrimSuffix(base, ".yaml")
		}
		recipes = append(recipes, r)
		return nil
	})
	return recipes, err
}

func LoadFromDir(dir string) ([]Recipe, error) {
	return LoadFromFS(os.DirFS(dir), ".")
}

func LoadFromFile(path string) (*Recipe, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var r Recipe
	if err := yaml.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	if r.ID == "" {
		base := filepath.Base(path)
		r.ID = strings.TrimSuffix(base, ".yaml")
	}
	return &r, nil
}
