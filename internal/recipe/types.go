package recipe

type RecipeType string

const (
	TypeCLITool RecipeType = "cli-tool"
	TypeMCP     RecipeType = "mcp"
	TypeSkill   RecipeType = "skill"
	TypeCommand RecipeType = "command"
	TypeRule    RecipeType = "rule"
	TypeBundle  RecipeType = "bundle"
)

type Recipe struct {
	ID          string     `yaml:"id" json:"id"`
	Name        string     `yaml:"name" json:"name"`
	Type        RecipeType `yaml:"type" json:"type"`
	Description string     `yaml:"description" json:"description"`
	Tags        []string   `yaml:"tags,omitempty" json:"tags,omitempty"`
	Difficulty  string     `yaml:"difficulty,omitempty" json:"difficulty,omitempty"`
	Install     *Install   `yaml:"install,omitempty" json:"install,omitempty"`
	Targets     *Targets   `yaml:"targets,omitempty" json:"targets,omitempty"`
	Prompts     []Prompt   `yaml:"prompts,omitempty" json:"prompts,omitempty"`
	Files       []FileSpec `yaml:"files,omitempty" json:"files,omitempty"`
	Quickstart  []QSEntry  `yaml:"quickstart,omitempty" json:"quickstart,omitempty"`
	Includes    []string   `yaml:"includes,omitempty" json:"includes,omitempty"`
}

type Install struct {
	Method      string   `yaml:"method" json:"method"`
	Package     string   `yaml:"package" json:"package"`
	PostInstall []string `yaml:"post_install,omitempty" json:"post_install,omitempty"`
}

type Targets struct {
	Cursor    *TargetConfig `yaml:"cursor,omitempty" json:"cursor,omitempty"`
	ClaudeCode *TargetConfig `yaml:"claude-code,omitempty" json:"claude-code,omitempty"`
}

type TargetConfig struct {
	ConfigPath string                 `yaml:"config_path" json:"config_path"`
	Config     map[string]interface{} `yaml:"config" json:"config"`
}

type Prompt struct {
	Key     string `yaml:"key" json:"key"`
	Ask     string `yaml:"ask" json:"ask"`
	Default string `yaml:"default,omitempty" json:"default,omitempty"`
	Secret  bool   `yaml:"secret,omitempty" json:"secret,omitempty"`
}

type FileSpec struct {
	Source string `yaml:"source" json:"source"`
	Dest   string `yaml:"dest" json:"dest"`
}

type QSEntry struct {
	Title   string `yaml:"title" json:"title"`
	Command string `yaml:"command,omitempty" json:"command,omitempty"`
	Explain string `yaml:"explain" json:"explain"`
}

type RecipeIndex struct {
	Version string   `json:"version"`
	Updated string   `json:"updated"`
	Recipes []Recipe `json:"recipes"`
}
