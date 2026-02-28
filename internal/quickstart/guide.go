package quickstart

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/yangduck/yduck/internal/recipe"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
	descStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
)

func Show(rec recipe.Recipe) {
	if len(rec.Quickstart) == 0 {
		return
	}
	fmt.Println()
	fmt.Println(titleStyle.Render("━━━━ " + rec.Name + " 快速开始 ━━━━"))
	for i, qs := range rec.Quickstart {
		fmt.Printf("\n  %d. %s\n", i+1, qs.Title)
		if qs.Command != "" {
			fmt.Printf("     $ %s\n", qs.Command)
		}
		fmt.Printf("     %s\n", descStyle.Render(qs.Explain))
	}
	fmt.Println()
}
