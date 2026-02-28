package installer

import (
	"github.com/yangduck/yduck/internal/recipe"
)

type RuleInstaller struct {
	skill *SkillInstaller
}

func NewRuleInstaller() *RuleInstaller {
	return &RuleInstaller{skill: NewSkillInstaller()}
}

func (r *RuleInstaller) Install(rec *recipe.Recipe) error {
	return r.skill.Install(rec)
}

func (r *RuleInstaller) IsInstalled(rec *recipe.Recipe) bool {
	return r.skill.IsInstalled(rec)
}
