package installer

import (
	"fmt"

	ylog "github.com/yangduck/yduck/internal/log"
	"github.com/yangduck/yduck/internal/recipe"
)

type BundleInstaller struct {
	registry *recipe.Registry
	brew     *BrewInstaller
	mcp      *MCPInstaller
	skill    *SkillInstaller
	command  *CommandInstaller
	rule     *RuleInstaller
}

func NewBundleInstaller(reg *recipe.Registry) *BundleInstaller {
	return &BundleInstaller{
		registry: reg,
		brew:     NewBrewInstaller(),
		mcp:      NewMCPInstaller(),
		skill:    NewSkillInstaller(),
		command:  NewCommandInstaller(),
		rule:     NewRuleInstaller(),
	}
}

type BundleResult struct {
	Installed []string
	Skipped   []string
	Failed    map[string]error
}

func (b *BundleInstaller) Install(bundle *recipe.Recipe, promptValues map[string]string, targets []string) (*BundleResult, error) {
	if bundle.Type != recipe.TypeBundle {
		return nil, fmt.Errorf("%s is not a bundle", bundle.ID)
	}
	ylog.S.Infow("installing bundle", "id", bundle.ID, "includes", len(bundle.Includes))
	result := &BundleResult{Failed: make(map[string]error)}
	for _, id := range bundle.Includes {
		rec, ok := b.registry.Get(id)
		if !ok {
			ylog.S.Warnw("recipe not found in bundle", "bundle", bundle.ID, "recipe", id)
			result.Failed[id] = fmt.Errorf("recipe not found: %s", id)
			continue
		}
		err := b.installOne(&rec, promptValues, targets)
		if err != nil {
			ylog.S.Errorw("bundle item install failed", "recipe", id, "error", err)
			result.Failed[id] = err
		} else {
			result.Installed = append(result.Installed, id)
		}
	}
	ylog.S.Infow("bundle install complete", "installed", len(result.Installed), "failed", len(result.Failed))
	return result, nil
}

func (b *BundleInstaller) installOne(rec *recipe.Recipe, promptValues map[string]string, targets []string) error {
	switch rec.Type {
	case recipe.TypeCLITool:
		if rec.Install == nil {
			return fmt.Errorf("no install config for %s", rec.ID)
		}
		if installed, _ := b.brew.IsInstalled(rec.Install.Package); installed {
			return nil
		}
		if err := b.brew.Install(rec.Install.Package); err != nil {
			return err
		}
		return b.brew.RunPostInstall(rec.Install.PostInstall)
	case recipe.TypeMCP:
		for _, t := range targets {
			if err := b.mcp.Install(rec, t, promptValues); err != nil {
				return err
			}
		}
		return nil
	case recipe.TypeSkill:
		return b.skill.Install(rec)
	case recipe.TypeCommand:
		return b.command.Install(rec)
	case recipe.TypeRule:
		return b.rule.Install(rec)
	default:
		return fmt.Errorf("unsupported recipe type: %s", rec.Type)
	}
}
