package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yangduck/yduck/internal/config"
	"github.com/yangduck/yduck/internal/recipe"
	s "github.com/yangduck/yduck/internal/tui/styles"
	"github.com/yangduck/yduck/internal/tui/views"
)

type installDoneMsg struct{ err error }

// installExecCmd implements tea.ExecCommand to run the install flow
// in the real terminal (alt-screen is temporarily released).
type installExecCmd struct {
	flow   *views.InstallFlow
	recipe recipe.Recipe
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func (c *installExecCmd) Run() error {
	err := c.flow.InstallRecipe(c.recipe)
	if err != nil {
		fmt.Fprintln(c.stdout)
		fmt.Fprintln(c.stdout, s.ErrorStyle.Render("  ✗ 安装失败: "+err.Error()))
	}
	fmt.Fprintln(c.stdout)
	fmt.Fprintln(c.stdout, s.DescStyle.Render("  按回车键返回..."))
	buf := make([]byte, 1)
	c.stdin.Read(buf)
	return err
}

func (c *installExecCmd) SetStdin(r io.Reader)  { c.stdin = r }
func (c *installExecCmd) SetStdout(w io.Writer) { c.stdout = w }
func (c *installExecCmd) SetStderr(w io.Writer) { c.stderr = w }

type App struct {
	registry    *recipe.Registry
	config      *config.Config
	installFlow *views.InstallFlow
	currentView views.View
	width       int
	height      int
}

func NewApp(reg *recipe.Registry, cfg *config.Config) *App {
	return &App{
		registry:    reg,
		config:      cfg,
		installFlow: views.NewInstallFlow(reg, cfg),
	}
}

func (a *App) Run() error {
	if a.config.IsFirstTime() {
		if err := a.runFirstTime(); err != nil {
			return err
		}
	}
	return a.runTUI("")
}

func (a *App) RunBrowse(category string) error {
	return a.runTUI(category)
}

func (a *App) RunSearch(keyword string) error {
	return a.runTUISearch(keyword)
}

func (a *App) runTUI(category string) error {
	switch {
	case category == "__installed__":
		a.currentView = views.NewInstalledView(a.registry, a.config, a.installFlow)
	case category != "":
		a.currentView = views.NewBrowseView(a.registry, a.config, a.installFlow, category)
	default:
		a.currentView = views.NewHomeView(a.registry, a.config, a.installFlow)
	}

	p := tea.NewProgram(a, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func (a *App) runTUISearch(keyword string) error {
	bv := views.NewBrowseView(a.registry, a.config, a.installFlow, "")
	bv.SetSearchMode(keyword)
	a.currentView = bv

	p := tea.NewProgram(a, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func (a *App) Init() tea.Cmd {
	return a.currentView.Init()
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return a, tea.Quit
		}

	case views.SwitchViewMsg:
		return a.switchView(msg)

	case views.InstallRecipeMsg:
		return a, a.doInstall(msg.Recipe)

	case views.ToggleModeMsg:
		if a.config.IsBeginner() {
			_ = a.config.SetMode(config.ModeAdvanced)
		} else {
			_ = a.config.SetMode(config.ModeBeginner)
		}
		a.currentView = views.NewHomeView(a.registry, a.config, a.installFlow)
		return a, a.currentView.Init()

	case installDoneMsg:
		if dv, ok := a.currentView.(*views.DetailView); ok {
			a.currentView = views.NewDetailView(a.registry, a.config, a.installFlow, dv.RecipeID())
		}
	}

	newView, cmd := a.currentView.Update(msg)
	a.currentView = newView
	return a, cmd
}

func (a *App) View() string {
	return a.currentView.View()
}

func (a *App) switchView(msg views.SwitchViewMsg) (tea.Model, tea.Cmd) {
	switch msg.Target {
	case views.TargetHome:
		a.currentView = views.NewHomeView(a.registry, a.config, a.installFlow)
	case views.TargetBrowse:
		bv := views.NewBrowseView(a.registry, a.config, a.installFlow, msg.Category)
		if msg.SearchTerm == "/" {
			bv.SetSearchMode("")
		}
		a.currentView = bv
	case views.TargetDetail:
		a.currentView = views.NewDetailView(a.registry, a.config, a.installFlow, msg.RecipeID)
	case views.TargetInstalled:
		a.currentView = views.NewInstalledView(a.registry, a.config, a.installFlow)
	}
	return a, a.currentView.Init()
}

func (a *App) doInstall(rec recipe.Recipe) tea.Cmd {
	cmd := &installExecCmd{flow: a.installFlow, recipe: rec}
	return tea.Exec(cmd, func(err error) tea.Msg {
		return installDoneMsg{err: err}
	})
}

func (a *App) runFirstTime() error {
	fmt.Println()
	fmt.Println(s.BannerStyle.Render(s.DuckBanner))
	fmt.Println(s.TitleStyle.Render("  欢迎使用 YangDuck！"))
	fmt.Println()
	fmt.Println(s.DescStyle.Render("  我是你的开发环境配置助手。"))
	fmt.Println(s.DescStyle.Render("  我可以帮你："))
	fmt.Println(s.DescStyle.Render("  · 安装好用的命令行工具"))
	fmt.Println(s.DescStyle.Render("  · 配置 AI 编码助手（Cursor、Claude Code）的扩展能力"))
	fmt.Println()
	fmt.Println(lipgloss.NewStyle().Foreground(s.ColorMuted).Render("──────────────────────────────────────────────────"))
	fmt.Println()

	var identity string
	var editor string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("你现在的身份是？").
				Options(
					huh.NewOption("🌱 我是编程新手，带我一步步来", "beginner"),
					huh.NewOption("🔧 我有经验，但想快速配好 AI 工具", "intermediate"),
					huh.NewOption("⚡ 我是老手，直接上", "advanced"),
				).
				Value(&identity),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("你使用的 AI 编码工具是？").
				Description("MCP、Skill、Rule 等配方会安装到对应工具的配置目录").
				Options(
					huh.NewOption("🖥️  Cursor", "cursor"),
					huh.NewOption("🤖 Claude Code", "claude-code"),
					huh.NewOption("📦 两个都用", "both"),
				).
				Value(&editor),
		),
	)
	if err := form.Run(); err != nil {
		return err
	}

	if identity == "advanced" {
		_ = a.config.SetMode(config.ModeAdvanced)
		fmt.Println(s.SuccessStyle.Render("  ✓ 已切换到高级模式"))
	}

	a.config.Editor = config.Editor(editor)
	editorNames := map[string]string{"cursor": "Cursor", "claude-code": "Claude Code", "both": "Cursor + Claude Code"}
	fmt.Println(s.SuccessStyle.Render("  ✓ AI 工具: " + editorNames[editor]))
	fmt.Println()

	return a.config.Save()
}
