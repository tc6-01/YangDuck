package recipe

import (
	"embed"
	"encoding/json"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

//go:embed all:schemas
var embeddedSchemas embed.FS

var schemaMap = map[RecipeType]string{
	TypeCLITool: "schemas/cli-tool.schema.json",
	TypeMCP:     "schemas/mcp.schema.json",
	TypeSkill:   "schemas/skill.schema.json",
	TypeCommand: "schemas/skill.schema.json",
	TypeRule:    "schemas/skill.schema.json",
	TypeBundle:  "schemas/bundle.schema.json",
}

func Validate(r *Recipe) ([]string, error) {
	schemaFile, ok := schemaMap[r.Type]
	if !ok {
		return nil, fmt.Errorf("unknown recipe type: %s", r.Type)
	}
	schemaData, err := embeddedSchemas.ReadFile(schemaFile)
	if err != nil {
		return nil, fmt.Errorf("read schema %s: %w", schemaFile, err)
	}
	recipeJSON, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshal recipe: %w", err)
	}
	schemaLoader := gojsonschema.NewBytesLoader(schemaData)
	docLoader := gojsonschema.NewBytesLoader(recipeJSON)
	result, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}
	if result.Valid() {
		return nil, nil
	}
	var errs []string
	for _, e := range result.Errors() {
		errs = append(errs, e.String())
	}
	return errs, nil
}
