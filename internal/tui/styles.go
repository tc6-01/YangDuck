package tui

import "github.com/charmbracelet/lipgloss"

var (
	ColorPrimary   = lipgloss.Color("#FFD700")
	ColorSecondary = lipgloss.Color("#87CEEB")
	ColorSuccess   = lipgloss.Color("#32CD32")
	ColorError     = lipgloss.Color("#FF6347")
	ColorMuted     = lipgloss.Color("#808080")
	ColorWhite     = lipgloss.Color("#FFFFFF")

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			MarginBottom(1)

	DescStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 2)

	HintStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Italic(true)

	BannerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			MarginBottom(1)
)

const DuckBanner = `
  🐤 YangDuck
`
