package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/yangduck/yduck/internal/recipe"
)

func main() {
	recipes, err := recipe.LoadFromDir("recipes")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading recipes: %v\n", err)
		os.Exit(1)
	}

	for i := range recipes {
		errs, err := recipe.Validate(&recipes[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Validation error for %s: %v\n", recipes[i].ID, err)
			os.Exit(1)
		}
		if len(errs) > 0 {
			fmt.Fprintf(os.Stderr, "Schema errors for %s:\n", recipes[i].ID)
			for _, e := range errs {
				fmt.Fprintf(os.Stderr, "  - %s\n", e)
			}
			os.Exit(1)
		}
	}

	index := recipe.RecipeIndex{
		Version: time.Now().Format("2006-01-02"),
		Updated: time.Now().UTC().Format(time.RFC3339),
		Recipes: recipes,
	}

	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling index: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile("recipes-index.json", data, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing index: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Generated recipes-index.json with %d recipes\n", len(recipes))
}
