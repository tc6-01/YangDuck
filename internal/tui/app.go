package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/yangduck/yduck/internal/config"
	"github.com/yangduck/yduck/internal/installer"
	ylog "github.com/yangduck/yduck/internal/log"
	"github.com/yangduck/yduck/internal/recipe"
)

var (
	divider    = lipgloss.NewStyle().Foreground(ColorMuted).Render(strings.Repeat("─", 50))
	stepStyle  = lipgloss.NewStyle().Foreground(ColorSecondary).Bold(true)
	noteStyle  = lipgloss.NewStyle().Foreground(ColorSecondary).Italic(true).PaddingLeft(2)
	warnStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")).Bold(true)
	cmdStyle   = lipgloss.NewStyle().Foreground(ColorPrimary).PaddingLeft(4)
	checkStyle = lipgloss.NewStyle().Foreground(ColorSuccess)
	crossStyle = lipgloss.NewStyle().Foreground(ColorError)
)

type App struct {
	registry *recipe.Registry
	config   *config.Config
	brew     *installer.BrewInstaller
	mcp      *installer.MCPInstaller
	skill    *installer.SkillInstaller
	command  *installer.CommandInstaller
	rule     *installer.RuleInstaller
	bundle   *installer.BundleInstaller
}

func NewApp(reg *recipe.Registry, cfg *config.Config) *App {
	return &App{
		registry: reg,
		config:   cfg,
		brew:     installer.NewBrewInstaller(),
		mcp:      installer.NewMCPInstaller(),
		skill:    installer.NewSkillInstaller(),
		command:  installer.NewCommandInstaller(),
		rule:     installer.NewRuleInstaller(),
		bundle:   installer.NewBundleInstaller(reg),
	}
}

func (a *App) Run() error {
	if a.config.IsBeginner() {
		return a.runFirstTime()
	}
	return a.runMainMenu()
}

func (a *App) runFirstTime() error {
	fmt.Println()
	fmt.Println(BannerStyle.Render(DuckBanner))
	fmt.Println(TitleStyle.Render("  欢迎使用 YangDuck！"))
	fmt.Println()
	fmt.Println(DescStyle.Render("  我是你的开发环境配置助手。"))
	fmt.Println(DescStyle.Render("  我可以帮你："))
	fmt.Println(DescStyle.Render("  · 安装好用的命令行工具"))
	fmt.Println(DescStyle.Render("  · 配置 AI 编码助手（Cursor、Claude Code）的扩展能力"))
	fmt.Println()
	fmt.Println(divider)
	fmt.Println()

	var identity string
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
	)
	if err := form.Run(); err != nil {
		return err
	}

	if identity == "advanced" {
		_ = a.config.SetMode(config.ModeAdvanced)
		fmt.Println(SuccessStyle.Render("  ✓ 已切换到高级模式"))
		fmt.Println()
	}
	return a.runMainMenu()
}

func (a *App) runMainMenu() error {
	for {
		fmt.Println()
		var choice string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("请选择你想做什么").
					Options(
						huh.NewOption("🚀 快速开始 — 一键安装推荐工具套餐", "quick"),
						huh.NewOption("🔧 CLI 工具 — 浏览和安装命令行工具", "cli"),
						huh.NewOption("🤖 AI 工具 — 配置 MCP / Skills / Commands", "ai"),
						huh.NewOption("📦 套餐安装 — 按角色安装工具组合", "bundle"),
						huh.NewOption("⚙️  设置", "settings"),
						huh.NewOption("👋 退出", "exit"),
					).
					Value(&choice),
			),
		)
		if err := form.Run(); err != nil {
			return err
		}

		switch choice {
		case "quick":
			if err := a.runQuickStart(); err != nil {
				fmt.Println(ErrorStyle.Render("  错误: " + err.Error()))
			}
		case "cli":
			if err := a.runCLITools(); err != nil {
				fmt.Println(ErrorStyle.Render("  错误: " + err.Error()))
			}
		case "ai":
			if err := a.runAITools(); err != nil {
				fmt.Println(ErrorStyle.Render("  错误: " + err.Error()))
			}
		case "bundle":
			if err := a.runBundles(); err != nil {
				fmt.Println(ErrorStyle.Render("  错误: " + err.Error()))
			}
		case "settings":
			if err := a.runSettings(); err != nil {
				fmt.Println(ErrorStyle.Render("  错误: " + err.Error()))
			}
		case "exit":
			fmt.Println()
			fmt.Println("  👋 再见！下次见~")
			fmt.Println()
			return nil
		}
	}
}

func (a *App) runQuickStart() error {
	bundles := a.registry.Bundles()
	if len(bundles) == 0 {
		fmt.Println(HintStyle.Render("  暂无可用套餐"))
		return nil
	}

	fmt.Println()
	fmt.Println(TitleStyle.Render("  🚀 快速开始"))
	fmt.Println()

	if a.config.IsBeginner() {
		fmt.Println(noteStyle.Render("💡 这会安装一组精选的开发工具，让你的终端更好用。"))
		fmt.Println(noteStyle.Render("   所有工具都可以单独卸载，不用担心。"))
		fmt.Println()
	}

	rec := bundles[0]
	fmt.Printf("  📦 套餐: %s\n", TitleStyle.Render(rec.Name))
	fmt.Println(DescStyle.Render("  " + rec.Description))
	fmt.Println()

	fmt.Println(DescStyle.Render("  包含以下工具:"))
	for i, id := range rec.Includes {
		sub, ok := a.registry.Get(id)
		if ok {
			installed, _ := a.isRecipeInstalled(sub)
			status := "  "
			if installed {
				status = checkStyle.Render("✓ ")
			}
			fmt.Printf("  %s%d. %s — %s\n", status, i+1, sub.Name, DescStyle.Render(sub.Description))
		} else {
			fmt.Printf("    %d. %s\n", i+1, id)
		}
	}
	fmt.Println()

	var confirm bool
	huh.NewConfirm().Title("开始安装？").Value(&confirm).Run()
	if !confirm {
		fmt.Println(DescStyle.Render("  好的，有需要随时回来~"))
		return nil
	}

	fmt.Println()
	fmt.Println(divider)
	total := len(rec.Includes)
	var installed, skipped, failed int
	for i, id := range rec.Includes {
		sub, ok := a.registry.Get(id)
		if !ok {
			fmt.Printf("  %s %s 配方未找到\n", crossStyle.Render("✗"), id)
			failed++
			continue
		}
		alreadyInstalled, _ := a.isRecipeInstalled(sub)
		if alreadyInstalled {
			fmt.Printf("  %s %s 已安装，跳过\n", checkStyle.Render("✓"), sub.Name)
			skipped++
			continue
		}
		fmt.Printf("\n  [%d/%d] 正在安装 %s...\n", i+1, total, stepStyle.Render(sub.Name))
		if a.config.IsBeginner() {
			fmt.Println(noteStyle.Render("ℹ " + sub.Description))
		}
		err := a.installRecipeSilent(sub)
		if err != nil {
			ylog.S.Errorw("quick start install failed", "recipe", sub.ID, "error", err)
			fmt.Printf("  %s 安装失败: %s\n", crossStyle.Render("✗"), err.Error())
			failed++
		} else {
			ylog.S.Infow("quick start installed", "recipe", sub.ID)
			fmt.Printf("  %s %s 安装完成\n", checkStyle.Render("✓"), sub.Name)
			installed++
		}
	}

	fmt.Println()
	fmt.Println(divider)
	fmt.Println()
	fmt.Println(TitleStyle.Render("  🎉 安装完成！"))
	fmt.Println()
	if installed > 0 {
		fmt.Printf("  %s 新安装 %d 个工具\n", checkStyle.Render("✓"), installed)
	}
	if skipped > 0 {
		fmt.Printf("  %s 跳过 %d 个（已安装）\n", checkStyle.Render("↷"), skipped)
	}
	if failed > 0 {
		fmt.Printf("  %s 失败 %d 个\n", crossStyle.Render("✗"), failed)
	}

	if a.config.IsBeginner() && installed > 0 {
		fmt.Println()
		fmt.Println(divider)
		fmt.Println()
		fmt.Println(SubtitleStyle.Render("  🐤 接下来试试这些:"))
		fmt.Println()
		a.showBundleQuickstart(rec)
	}

	return nil
}

func (a *App) runCLITools() error {
	recipes := a.registry.ListByType(recipe.TypeCLITool)
	if len(recipes) == 0 {
		fmt.Println(HintStyle.Render("  暂无可用的 CLI 工具配方"))
		return nil
	}
	options := make([]huh.Option[string], 0, len(recipes)+1)
	for _, r := range recipes {
		installed, _ := a.brew.IsInstalled(r.Install.Package)
		mark := "  "
		if installed {
			mark = "✓ "
		}
		label := fmt.Sprintf("%s%s — %s", mark, r.Name, r.Description)
		options = append(options, huh.NewOption(label, r.ID))
	}
	options = append(options, huh.NewOption("◀ 返回", "back"))

	var selected string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().Title("CLI 工具（✓ 表示已安装）").Options(options...).Value(&selected),
		),
	)
	if err := form.Run(); err != nil {
		return err
	}
	if selected == "back" {
		return nil
	}
	return a.installCLITool(selected)
}

func (a *App) installCLITool(id string) error {
	rec, ok := a.registry.Get(id)
	if !ok || rec.Install == nil {
		return fmt.Errorf("配方未找到: %s", id)
	}

	fmt.Println()
	fmt.Println(divider)
	fmt.Println()

	if a.config.IsBeginner() {
		fmt.Println(TitleStyle.Render("  🔧 " + rec.Name))
		fmt.Println()
		fmt.Println(noteStyle.Render("ℹ " + rec.Description))
		fmt.Println()
	}

	if installed, ver := a.brew.IsInstalled(rec.Install.Package); installed {
		fmt.Printf("  %s %s 已安装 (v%s)\n", checkStyle.Render("✓"), rec.Name, ver)
		fmt.Println()
		if a.config.IsBeginner() {
			fmt.Println(noteStyle.Render("💡 虽然已安装，还是给你看看怎么用:"))
		}
		a.showQuickstart(rec)
		return nil
	}

	fmt.Println(stepStyle.Render("  ▸ 正在安装 " + rec.Name + "..."))
	fmt.Println(DescStyle.Render("    通过 Homebrew 安装: brew install " + rec.Install.Package))
	fmt.Println()

	if err := a.brew.Install(rec.Install.Package); err != nil {
		return fmt.Errorf("安装失败: %w", err)
	}
	if len(rec.Install.PostInstall) > 0 {
		fmt.Println(DescStyle.Render("    运行安装后配置..."))
		if err := a.brew.RunPostInstall(rec.Install.PostInstall); err != nil {
			fmt.Println(warnStyle.Render("    ⚠ 安装后配置有警告: " + err.Error()))
			fmt.Println(DescStyle.Render("    不影响使用，可以忽略"))
		}
	}

	fmt.Println()
	fmt.Printf("  %s %s 安装完成！\n", checkStyle.Render("✓"), rec.Name)
	a.showQuickstart(rec)
	return nil
}

func (a *App) runAITools() error {
	fmt.Println()
	if a.config.IsBeginner() {
		fmt.Println(noteStyle.Render("💡 AI 工具可以增强 Cursor 和 Claude Code 的能力。"))
		fmt.Println(noteStyle.Render("   安装后 AI 助手就能帮你操作数据库、读写文件等。"))
		fmt.Println()
	}

	var aiType string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().Title("选择 AI 工具类别").Options(
				huh.NewOption("🔌 MCP 服务器 — 让 AI 连接外部工具（数据库、文件等）", "mcp"),
				huh.NewOption("📝 Skills — 给 AI 增加新技能", "skill"),
				huh.NewOption("⌨️  Commands — 添加斜杠命令", "command"),
				huh.NewOption("📏 Rules — 设置编码规范", "rule"),
				huh.NewOption("◀ 返回", "back"),
			).Value(&aiType),
		),
	)
	if err := form.Run(); err != nil {
		return err
	}
	if aiType == "back" {
		return nil
	}

	typeMap := map[string]recipe.RecipeType{
		"mcp": recipe.TypeMCP, "skill": recipe.TypeSkill,
		"command": recipe.TypeCommand, "rule": recipe.TypeRule,
	}
	typeNames := map[string]string{
		"mcp": "MCP 服务器", "skill": "Skill", "command": "Command", "rule": "Rule",
	}
	rt := typeMap[aiType]
	recipes := a.registry.ListByType(rt)
	if len(recipes) == 0 {
		fmt.Println(HintStyle.Render("  暂无可用的" + typeNames[aiType] + "配方"))
		return nil
	}
	options := make([]huh.Option[string], 0, len(recipes)+1)
	for _, r := range recipes {
		options = append(options, huh.NewOption(fmt.Sprintf("%s — %s", r.Name, r.Description), r.ID))
	}
	options = append(options, huh.NewOption("◀ 返回", "back"))

	var selected string
	huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().Title("选择要安装的 " + typeNames[aiType]).Options(options...).Value(&selected),
	)).Run()
	if selected == "back" {
		return nil
	}

	rec, _ := a.registry.Get(selected)
	switch rec.Type {
	case recipe.TypeMCP:
		return a.installMCP(rec)
	case recipe.TypeSkill:
		return a.installFileRecipe(rec, a.skill)
	case recipe.TypeCommand:
		return a.installFileRecipe(rec, a.command)
	case recipe.TypeRule:
		return a.installFileRecipe(rec, a.rule)
	}
	return nil
}

func (a *App) installMCP(rec recipe.Recipe) error {
	fmt.Println()
	fmt.Println(divider)
	fmt.Println()

	fmt.Println(TitleStyle.Render("  🔌 " + rec.Name))
	fmt.Println()

	if a.config.IsBeginner() {
		fmt.Println(BoxStyle.Render(
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

	if a.mcp.IsConfigured(&rec, "cursor") {
		fmt.Println(warnStyle.Render("  ⚠ 这个 MCP 在 Cursor 中已经配置过了"))
		var overwrite bool
		huh.NewConfirm().Title("要覆盖现有配置吗？").Value(&overwrite).Run()
		if !overwrite {
			fmt.Println(DescStyle.Render("  好的，保留现有配置"))
			return nil
		}
	}

	if len(rec.Prompts) > 0 {
		fmt.Println(stepStyle.Render("  ▸ 需要你填几个信息"))
		if a.config.IsBeginner() {
			fmt.Println(noteStyle.Render("不知道填什么？直接按回车用默认值就行"))
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

	targets := []string{}
	if rec.Targets != nil && rec.Targets.Cursor != nil {
		targets = append(targets, "cursor")
	}
	if rec.Targets != nil && rec.Targets.ClaudeCode != nil {
		targets = append(targets, "claude-code")
	}

	fmt.Println()
	for _, t := range targets {
		targetName := map[string]string{"cursor": "Cursor", "claude-code": "Claude Code"}[t]
		fmt.Printf("  %s 正在写入 %s 配置...\n", stepStyle.Render("▸"), targetName)
		if err := a.mcp.Install(&rec, t, promptValues); err != nil {
			ylog.S.Errorw("mcp config write failed", "target", t, "recipe", rec.ID, "error", err)
			fmt.Printf("  %s %s 配置失败: %s\n", crossStyle.Render("✗"), targetName, err.Error())
			continue
		}
		tc := rec.Targets.Cursor
		if t == "claude-code" {
			tc = rec.Targets.ClaudeCode
		}
		fmt.Printf("  %s 已写入 %s\n", checkStyle.Render("✓"), DescStyle.Render(tc.ConfigPath))
		fmt.Printf("  %s 原配置已备份到 %s\n", checkStyle.Render("✓"), DescStyle.Render(tc.ConfigPath+".backup"))
	}

	fmt.Println()
	fmt.Println(SuccessStyle.Render("  🎉 配置完成！"))
	a.showQuickstart(rec)

	if a.config.IsBeginner() {
		fmt.Println(warnStyle.Render("  ⚠ 记得重启 Cursor 让配置生效"))
		fmt.Println()
	}
	return nil
}

type fileInstaller interface {
	Install(rec *recipe.Recipe) error
}

func (a *App) installFileRecipe(rec recipe.Recipe, inst fileInstaller) error {
	fmt.Println()
	fmt.Println(divider)
	fmt.Println()

	if a.config.IsBeginner() {
		fmt.Println(TitleStyle.Render("  📝 " + rec.Name))
		fmt.Println()
		fmt.Println(noteStyle.Render("ℹ " + rec.Description))
		fmt.Println()
	}

	fmt.Println(stepStyle.Render("  ▸ 正在安装 " + rec.Name + "..."))
	if err := inst.Install(&rec); err != nil {
		return err
	}

	fmt.Println()
	fmt.Printf("  %s %s 安装完成！\n", checkStyle.Render("✓"), rec.Name)

	if len(rec.Files) > 0 && a.config.IsBeginner() {
		fmt.Println()
		fmt.Println(DescStyle.Render("  文件已安装到:"))
		for _, f := range rec.Files {
			fmt.Println(cmdStyle.Render("→ " + f.Dest))
		}
	}
	a.showQuickstart(rec)

	if a.config.IsBeginner() {
		fmt.Println(warnStyle.Render("  ⚠ 记得重启 Cursor 让新功能生效"))
		fmt.Println()
	}
	return nil
}

func (a *App) runBundles() error {
	bundles := a.registry.Bundles()
	if len(bundles) == 0 {
		fmt.Println(HintStyle.Render("  暂无可用套餐"))
		return nil
	}
	options := make([]huh.Option[string], 0, len(bundles)+1)
	for _, b := range bundles {
		label := fmt.Sprintf("%s — %s (%d 个工具)", b.Name, b.Description, len(b.Includes))
		options = append(options, huh.NewOption(label, b.ID))
	}
	options = append(options, huh.NewOption("◀ 返回", "back"))

	var selected string
	huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().Title("选择套餐").Options(options...).Value(&selected),
	)).Run()
	if selected == "back" {
		return nil
	}

	rec, _ := a.registry.Get(selected)

	fmt.Println()
	fmt.Println(TitleStyle.Render("  📦 " + rec.Name))
	fmt.Println(DescStyle.Render("  " + rec.Description))
	fmt.Println()
	fmt.Println(DescStyle.Render("  包含:"))
	for _, id := range rec.Includes {
		sub, ok := a.registry.Get(id)
		if ok {
			fmt.Printf("    · %s — %s\n", sub.Name, DescStyle.Render(sub.Description))
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
		sub, ok := a.registry.Get(id)
		if !ok {
			fmt.Printf("  %s %s 配方未找到\n", crossStyle.Render("✗"), id)
			failedN++
			continue
		}
		already, _ := a.isRecipeInstalled(sub)
		if already {
			fmt.Printf("  %s %s 已安装，跳过\n", checkStyle.Render("✓"), sub.Name)
			skippedN++
			continue
		}
		fmt.Printf("\n  [%d/%d] 正在安装 %s...\n", i+1, total, stepStyle.Render(sub.Name))
		err := a.installRecipeSilent(sub)
		if err != nil {
			ylog.S.Errorw("bundle item install failed", "recipe", sub.ID, "error", err)
			fmt.Printf("  %s 失败: %s\n", crossStyle.Render("✗"), err.Error())
			failedN++
		} else {
			ylog.S.Infow("bundle item installed", "recipe", sub.ID)
			fmt.Printf("  %s %s 完成\n", checkStyle.Render("✓"), sub.Name)
			installed++
		}
	}

	fmt.Println()
	fmt.Println(divider)
	fmt.Println()
	fmt.Println(SuccessStyle.Render("  🎉 套餐安装完成！"))
	fmt.Println()
	if installed > 0 {
		fmt.Printf("  %s 新安装 %d 个\n", checkStyle.Render("✓"), installed)
	}
	if skippedN > 0 {
		fmt.Printf("  %s 跳过 %d 个（已安装）\n", checkStyle.Render("↷"), skippedN)
	}
	if failedN > 0 {
		fmt.Printf("  %s 失败 %d 个\n", crossStyle.Render("✗"), failedN)
	}
	fmt.Println()
	return nil
}

func (a *App) runSettings() error {
	var choice string
	currentMode := "🌱 新手模式"
	if !a.config.IsBeginner() {
		currentMode = "⚡ 高级模式"
	}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(fmt.Sprintf("设置（当前: %s）", currentMode)).
				Options(
					huh.NewOption("🌱 切换到新手模式 — 逐步引导、详细解释", "beginner"),
					huh.NewOption("⚡ 切换到高级模式 — 精简输出、批量操作", "advanced"),
					huh.NewOption("◀ 返回", "back"),
				).
				Value(&choice),
		),
	)
	if err := form.Run(); err != nil {
		return err
	}
	switch choice {
	case "beginner":
		_ = a.config.SetMode(config.ModeBeginner)
		fmt.Println(SuccessStyle.Render("  ✓ 已切换到新手模式"))
	case "advanced":
		_ = a.config.SetMode(config.ModeAdvanced)
		fmt.Println(SuccessStyle.Render("  ✓ 已切换到高级模式"))
	}
	return nil
}

func (a *App) showQuickstart(rec recipe.Recipe) {
	if !a.config.IsBeginner() || len(rec.Quickstart) == 0 {
		return
	}
	fmt.Println()
	fmt.Println(SubtitleStyle.Render("  ━━━━ 马上试试 ━━━━"))
	for i, qs := range rec.Quickstart {
		fmt.Println()
		fmt.Printf("  %d. %s\n", i+1, stepStyle.Render(qs.Title))
		if qs.Command != "" {
			fmt.Println(cmdStyle.Render("$ " + qs.Command))
		}
		fmt.Println(noteStyle.Render(qs.Explain))
	}
	fmt.Println()
}

func (a *App) showBundleQuickstart(rec recipe.Recipe) {
	shown := 0
	for _, id := range rec.Includes {
		sub, ok := a.registry.Get(id)
		if !ok || len(sub.Quickstart) == 0 {
			continue
		}
		qs := sub.Quickstart[0]
		fmt.Printf("  · %s: %s\n", stepStyle.Render(sub.Name), qs.Title)
		if qs.Command != "" {
			fmt.Println(cmdStyle.Render("$ " + qs.Command))
		}
		fmt.Println(noteStyle.Render(qs.Explain))
		fmt.Println()
		shown++
		if shown >= 4 {
			fmt.Println(DescStyle.Render("  输入 yduck list --installed 查看全部已安装工具的用法"))
			break
		}
	}
}

func (a *App) isRecipeInstalled(rec recipe.Recipe) (bool, string) {
	switch rec.Type {
	case recipe.TypeCLITool:
		if rec.Install == nil {
			return false, ""
		}
		return a.brew.IsInstalled(rec.Install.Package)
	case recipe.TypeMCP:
		return a.mcp.IsConfigured(&rec, "cursor"), ""
	case recipe.TypeSkill:
		return a.skill.IsInstalled(&rec), ""
	case recipe.TypeCommand:
		return a.command.IsInstalled(&rec), ""
	case recipe.TypeRule:
		return a.rule.IsInstalled(&rec), ""
	default:
		return false, ""
	}
}

func (a *App) installRecipeSilent(rec recipe.Recipe) error {
	switch rec.Type {
	case recipe.TypeCLITool:
		if rec.Install == nil {
			return fmt.Errorf("no install config")
		}
		if err := a.brew.Install(rec.Install.Package); err != nil {
			return err
		}
		return a.brew.RunPostInstall(rec.Install.PostInstall)
	case recipe.TypeMCP:
		targets := []string{}
		if rec.Targets != nil && rec.Targets.Cursor != nil {
			targets = append(targets, "cursor")
		}
		if rec.Targets != nil && rec.Targets.ClaudeCode != nil {
			targets = append(targets, "claude-code")
		}
		for _, t := range targets {
			if err := a.mcp.Install(&rec, t, nil); err != nil {
				return err
			}
		}
		return nil
	case recipe.TypeSkill:
		return a.skill.Install(&rec)
	case recipe.TypeCommand:
		return a.command.Install(&rec)
	case recipe.TypeRule:
		return a.rule.Install(&rec)
	default:
		return fmt.Errorf("unsupported type: %s", rec.Type)
	}
}
