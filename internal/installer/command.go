package installer

import (
	"github.com/yangduck/yduck/internal/recipe"
)

type CommandInstaller struct {
	skill *SkillInstaller
}

func NewCommandInstaller() *CommandInstaller {
	return &CommandInstaller{skill: NewSkillInstaller()}
}

func (c *CommandInstaller) Install(rec *recipe.Recipe) error {
	return c.skill.Install(rec)
}

func (c *CommandInstaller) IsInstalled(rec *recipe.Recipe) bool {
	return c.skill.IsInstalled(rec)
}
