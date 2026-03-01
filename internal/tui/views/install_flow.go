package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/yangduck/yduck/internal/config"
	"github.com/yangduck/yduck/internal/installer"
	ylog "github.com/yangduck/yduck/internal/log"
	"github.com/yangduck/yduck/internal/recipe"
	s "github.com/yangduck/yduck/internal/tui/styles"
)

var divider = lipgloss.NewStyle().Foreground(s.ColorMuted).Render(strings.Repeat("─", 50))

type InstallFlow struct {
	Registry *recipe.Registry
	Config   *config.Config
	Brew     *installer.BrewInstaller
	MCP      *installer.MCPInstaller
	Skill    *installer.SkillInstaller
	Command  *installer.CommandInstaller
	Rule     *installer.RuleInstaller
	Bundle   *installer.BundleInstaller
}

func NewInstallFlow(reg *recipe.Registry, cfg *config.Config) *InstallFlow {
	return &InstallFlow{
		Registry: reg,
		Config:   cfg,
		Brew:     installer.NewBrewInstaller(),
		MCP:      installer.NewMCPInstaller(),
		Skill:    installer.NewSkillInstaller(),
		Command:  installer.NewCommandInstaller(),
		Rule:     installer.NewRuleInstaller(),
		Bundle:   installer.NewBundleInstaller(reg),
	}
}

func (f *InstallFlow) InstallRecipe(rec recipe.Recipe) error {
	switch rec.Type {
	case recipe.TypeCLITool:
		return f.installCLITool(rec)
	case recipe.TypeMCP:
		return f.installMCP(rec)
	case recipe.TypeSkill:
		return f.installFileRecipe(rec, f.Skill)
	case recipe.TypeCommand:
		return f.installFileRecipe(rec, f.Command)
	case recipe.TypeRule:
		return f.installFileRecipe(rec, f.Rule)
	case recipe.TypeBundle:
		return f.installBundle(rec)
	default:
		return fmt.Errorf("unsupported type: %s", rec.Type)
	}
}

func (f *InstallFlow) InstallRecipeSilent(rec recipe.Recipe) error {
	switch rec.Type {
	case recipe.TypeCLITool:
		if rec.Install == nil {
			return fmt.Errorf("no install config")
		}
		if err := f.Brew.Install(rec.Install.Package); err != nil {
			return err
		}
		return f.Brew.RunPostInstall(rec.Install.PostInstall)
	case recipe.TypeMCP:
		for _, t := range f.mcpTargets(rec) {
			if err := f.MCP.Install(&rec, t, nil); err != nil {
				return err
			}
		}
		return nil
	case recipe.TypeSkill:
		return f.Skill.Install(&rec)
	case recipe.TypeCommand:
		return f.Command.Install(&rec)
	case recipe.TypeRule:
		return f.Rule.Install(&rec)
	default:
		return fmt.Errorf("unsupported type: %s", rec.Type)
	}
}

func (f *InstallFlow) IsRecipeInstalled(rec recipe.Recipe) (bool, string) {
	switch rec.Type {
	case recipe.TypeCLITool:
		if rec.Install == nil {
			return false, ""
		}
		return f.Brew.IsInstalled(rec.Install.Package)
	case recipe.TypeMCP:
		for _, t := range f.mcpTargets(rec) {
			if f.MCP.IsConfigured(&rec, t) {
				return true, ""
			}
		}
		return false, ""
	case recipe.TypeSkill:
		return f.Skill.IsInstalled(&rec), ""
	case recipe.TypeCommand:
		return f.Command.IsInstalled(&rec), ""
	case recipe.TypeRule:
		return f.Rule.IsInstalled(&rec), ""
	default:
		return false, ""
	}
}

func (f *InstallFlow) installCLITool(rec recipe.Recipe) error {
	fmt.Println()
	fmt.Println(divider)
	fmt.Println()

	if f.Config.IsBeginner() {
		fmt.Println(s.TitleStyle.Render("  🔧 " + rec.Name))
		fmt.Println()
		fmt.Println(s.NoteStyle.Render("ℹ " + rec.Description))
		fmt.Println()
	}

	if installed, ver := f.Brew.IsInstalled(rec.Install.Package); installed {
		fmt.Printf("  %s %s 已安装 (v%s)\n", s.CheckStyle.Render("✓"), rec.Name, ver)
		fmt.Println()
		if f.Config.IsBeginner() {
			fmt.Println(s.NoteStyle.Render("💡 虽然已安装，还是给你看看怎么用:"))
		}
		f.showQuickstart(rec)
		return nil
	}

	fmt.Println(s.StepStyle.Render("  ▸ 正在安装 " + rec.Name + "..."))
	fmt.Println(s.DescStyle.Render("    通过 Homebrew 安装: brew install " + rec.Install.Package))
	fmt.Println()

	if err := f.Brew.Install(rec.Install.Package); err != nil {
		return fmt.Errorf("安装失败: %w", err)
	}
	if len(rec.Install.PostInstall) > 0 {
		fmt.Println(s.DescStyle.Render("    运行安装后配置..."))
		if err := f.Brew.RunPostInstall(rec.Install.PostInstall); err != nil {
			fmt.Println(s.WarnStyle.Render("    ⚠ 安装后配置有警告: " + err.Error()))
			fmt.Println(s.DescStyle.Render("    不影响使用，可以忽略"))
		}
	}

	fmt.Println()
	fmt.Printf("  %s %s 安装完成！\n", s.CheckStyle.Render("✓"), rec.Name)
	f.showQuickstart(rec)
	return nil
}

func (f *InstallFlow) installMCP(rec recipe.Recipe) error {
	fmt.Println()
	fmt.Println(divider)
	fmt.Println()

	fmt.Println(s.TitleStyle.Render("  🔌 " + rec.Name))
	fmt.Println()

	if f.Config.IsBeginner() {
		fmt.Println(s.BoxStyle.Render(
			"💡 这是什么？\n\n" +
				rec.Description + "\n\n" +
				"MCP（Model Context Protocol）是让 AI 助手连接外部工具的协议。\n" +
				"装了之后，你在 Cursor 里就可以让 AI 直接使用这个工具了。\n\n" +
				"🔧 会做什么？\n" +
				"  1. 把连接信息写入 Cursor/Claude Code 的配置文件\n" +
				"  2. 之后你在 AI 对话中就能直接用了",
		))
		fmt.Println()
	}

	if f.MCP.IsConfigured(&rec, "cursor") {
		fmt.Println(s.WarnStyle.Render("  ⚠ 这个 MCP 在 Cursor 中已经配置过了"))
		var overwrite bool
		huh.NewConfirm().Title("要覆盖现有配置吗？").Value(&overwrite).Run()
		if !overwrite {
			fmt.Println(s.DescStyle.Render("  好的，保留现有配置"))
			return nil
		}
	}

	if len(rec.Prompts) > 0 {
		fmt.Println(s.StepStyle.Render("  ▸ 需要你填几个信息"))
		if f.Config.IsBeginner() {
			fmt.Println(s.NoteStyle.Render("不知道填什么？直接按回车用默认值就行"))
		}
		fmt.Println()
	}

	promptValues := make(map[string]string)
	for _, p := range rec.Prompts {
		var val string
		title := p.Ask
		if p.Default != "" {
			title = fmt.Sprintf("%s (默认: %s)", p.Ask, p.Default)
		}
		input := huh.NewInput().Title(title)
		if p.Default != "" {
			input.Placeholder(p.Default)
		}
		input.Value(&val)
		huh.NewForm(huh.NewGroup(input)).Run()
		if val == "" {
			val = p.Default
		}
		promptValues[p.Key] = val
	}

	targets := f.mcpTargets(rec)

	fmt.Println()
	for _, t := range targets {
		targetName := map[string]string{"cursor": "Cursor", "claude-code": "Claude Code"}[t]
		fmt.Printf("  %s 正在写入 %s 配置...\n", s.StepStyle.Render("▸"), targetName)
		if err := f.MCP.Install(&rec, t, promptValues); err != nil {
			ylog.S.Errorw("mcp config write failed", "target", t, "recipe", rec.ID, "error", err)
			fmt.Printf("  %s %s 配置失败: %s\n", s.CrossStyle.Render("✗"), targetName, err.Error())
			continue
		}
		tc := rec.Targets.Cursor
		if t == "claude-code" {
			tc = rec.Targets.ClaudeCode
		}
		fmt.Printf("  %s 已写入 %s\n", s.CheckStyle.Render("✓"), s.DescStyle.Render(tc.ConfigPath))
		fmt.Printf("  %s 原配置已备份到 %s\n", s.CheckStyle.Render("✓"), s.DescStyle.Render(tc.ConfigPath+".backup"))
	}

	fmt.Println()
	fmt.Println(s.SuccessStyle.Render("  🎉 配置完成！"))
	f.showQuickstart(rec)

	if f.Config.IsBeginner() {
		fmt.Println(s.WarnStyle.Render("  ⚠ 记得重启 Cursor 让配置生效"))
		fmt.Println()
	}
	return nil
}

type fileInstaller interface {
	Install(rec *recipe.Recipe) error
}

func (f *InstallFlow) installFileRecipe(rec recipe.Recipe, inst fileInstaller) error {
	fmt.Println()
	fmt.Println(divider)
	fmt.Println()

	if f.Config.IsBeginner() {
		fmt.Println(s.TitleStyle.Render("  📝 " + rec.Name))
		fmt.Println()
		fmt.Println(s.NoteStyle.Render("ℹ " + rec.Description))
		fmt.Println()
	}

	fmt.Println(s.StepStyle.Render("  ▸ 正在安装 " + rec.Name + "..."))
	if err := inst.Install(&rec); err != nil {
		return err
	}

	fmt.Println()
	fmt.Printf("  %s %s 安装完成！\n", s.CheckStyle.Render("✓"), rec.Name)

	if len(rec.Files) > 0 && f.Config.IsBeginner() {
		fmt.Println()
		fmt.Println(s.DescStyle.Render("  文件已安装到:"))
		for _, file := range rec.Files {
			fmt.Println(s.CmdStyle.Render("→ " + file.Dest))
		}
	}
	f.showQuickstart(rec)

	if f.Config.IsBeginner() {
		fmt.Println(s.WarnStyle.Render("  ⚠ 记得重启 Cursor 让新功能生效"))
		fmt.Println()
	}
	return nil
}

func (f *InstallFlow) installBundle(rec recipe.Recipe) error {
	fmt.Println()
	fmt.Println(s.TitleStyle.Render("  📦 " + rec.Name))
	fmt.Println(s.DescStyle.Render("  " + rec.Description))
	fmt.Println()
	fmt.Println(s.DescStyle.Render("  包含:"))
	for _, id := range rec.Includes {
		sub, ok := f.Registry.Get(id)
		if ok {
			fmt.Printf("    · %s — %s\n", sub.Name, s.DescStyle.Render(sub.Description))
		}
	}
	fmt.Println()

	var confirm bool
	huh.NewConfirm().Title("开始安装？").Value(&confirm).Run()
	if !confirm {
		return nil
	}

	fmt.Println()
	fmt.Println(divider)
	total := len(rec.Includes)
	var installed, skippedN, failedN int
	for i, id := range rec.Includes {
		sub, ok := f.Registry.Get(id)
		if !ok {
			fmt.Printf("  %s %s 配方未找到\n", s.CrossStyle.Render("✗"), id)
			failedN++
			continue
		}
		already, _ := f.IsRecipeInstalled(sub)
		if already {
			fmt.Printf("  %s %s 已安装，跳过\n", s.CheckStyle.Render("✓"), sub.Name)
			skippedN++
			continue
		}
		fmt.Printf("\n  [%d/%d] 正在安装 %s...\n", i+1, total, s.StepStyle.Render(sub.Name))
		err := f.InstallRecipeSilent(sub)
		if err != nil {
			ylog.S.Errorw("bundle item install failed", "recipe", sub.ID, "error", err)
			fmt.Printf("  %s 失败: %s\n", s.CrossStyle.Render("✗"), err.Error())
			failedN++
		} else {
			ylog.S.Infow("bundle item installed", "recipe", sub.ID)
			fmt.Printf("  %s %s 完成\n", s.CheckStyle.Render("✓"), sub.Name)
			installed++
		}
	}

	fmt.Println()
	fmt.Println(divider)
	fmt.Println()
	fmt.Println(s.SuccessStyle.Render("  🎉 套餐安装完成！"))
	fmt.Println()
	if installed > 0 {
		fmt.Printf("  %s 新安装 %d 个\n", s.CheckStyle.Render("✓"), installed)
	}
	if skippedN > 0 {
		fmt.Printf("  %s 跳过 %d 个（已安装）\n", s.CheckStyle.Render("↷"), skippedN)
	}
	if failedN > 0 {
		fmt.Printf("  %s 失败 %d 个\n", s.CrossStyle.Render("✗"), failedN)
	}
	fmt.Println()
	return nil
}

func (f *InstallFlow) showQuickstart(rec recipe.Recipe) {
	if !f.Config.IsBeginner() || len(rec.Quickstart) == 0 {
		return
	}
	fmt.Println()
	fmt.Println(s.SubtitleStyle.Render("  ━━━━ 马上试试 ━━━━"))
	for i, qs := range rec.Quickstart {
		fmt.Println()
		fmt.Printf("  %d. %s\n", i+1, s.StepStyle.Render(qs.Title))
		if qs.Command != "" {
			fmt.Println(s.CmdStyle.Render("$ " + qs.Command))
		}
		fmt.Println(s.NoteStyle.Render(qs.Explain))
	}
	fmt.Println()
}

func (f *InstallFlow) mcpTargets(rec recipe.Recipe) []string {
	var targets []string
	if rec.Targets != nil && rec.Targets.Cursor != nil && f.Config.ShouldInstallFor("cursor") {
		targets = append(targets, "cursor")
	}
	if rec.Targets != nil && rec.Targets.ClaudeCode != nil && f.Config.ShouldInstallFor("claude-code") {
		targets = append(targets, "claude-code")
	}
	return targets
}
