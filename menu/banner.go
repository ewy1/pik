package menu

import (
	"github.com/charmbracelet/lipgloss"
	"os/exec"
	"pik/flags"
	"pik/menu/style"
	"pik/model"
	"pik/paths"
	"strings"
)

var (
	BannerStyle = style.New(func() lipgloss.Style {
		return lipgloss.NewStyle()
	})
	BannerSourceLabelStyle = style.New(func() lipgloss.Style {
		return lipgloss.NewStyle().Faint(true).MarginRight(1)
	})
	BannerSubItemStyle = style.New(func() lipgloss.Style {
		return lipgloss.NewStyle().Faint(true).MarginRight(1)
	})
	BannerSubStyle = style.New(func() lipgloss.Style {
		return lipgloss.NewStyle()
	})
	BannerSelfStyle = style.New(func() lipgloss.Style {
		return lipgloss.NewStyle().MarginRight(1).Bold(true)
	})
	BannerPromptStyle = style.New(func() lipgloss.Style {
		return lipgloss.NewStyle()
	})
	BannerArgsStyle = style.New(func() lipgloss.Style {
		return lipgloss.NewStyle().MarginLeft(1)
	})
	BannerArgStyle = style.New(func() lipgloss.Style {
		return lipgloss.NewStyle()
	})
	BannerTerminatorColor = lipgloss.Color("1")
	BannerTerminatorStyle = style.New(func() lipgloss.Style {
		return lipgloss.NewStyle().Faint(true).Foreground(BannerTerminatorColor)
	})
	BannerDryColor = lipgloss.Color("1")
	BannerDryStyle = style.New(func() lipgloss.Style {
		return lipgloss.NewStyle().Foreground(BannerDryColor).Bold(true).MarginRight(1)
	})
)

func Banner(source *model.Source, target model.Target, args ...string) string {
	var parts, argParts []string
	if *flags.Dry {
		parts = append(parts, BannerDryStyle.Render("DRY"))
	}
	parts = append(parts, BannerPromptStyle.Render("> "))
	parts = append(parts, BannerSelfStyle.Render("pik"))
	parts = append(parts, BannerSourceLabelStyle.Render(source.Label()))
	if sub := target.Sub(); sub != nil {
		for i, s := range sub {
			sub[i] = BannerSubItemStyle.Render(s)
		}
		parts = append(parts, BannerSubStyle.Render(sub...))
	}
	parts = append(parts, target.ShortestId())
	if args != nil {
		needsTerminator := false
		for _, a := range args {
			if strings.HasPrefix(a, "-") {
				needsTerminator = true
			}
			argParts = append(argParts, BannerArgStyle.Render(a))
		}

		if needsTerminator {
			argParts = append([]string{BannerTerminatorStyle.Render("--")}, argParts...)
		}

		parts = append(parts, BannerArgsStyle.Render(argParts...))
	}
	result := BannerStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, parts...))
	return result
}

var (
	CmdStyle = style.New(func() lipgloss.Style {
		return lipgloss.NewStyle().Faint(true)
	})
	CmdDirStyle = style.New(func() lipgloss.Style {
		return lipgloss.NewStyle()
	})
	CmdArgStyle = style.New(func() lipgloss.Style {
		return lipgloss.NewStyle()
	})
)

func InlineCmd(cmd *exec.Cmd) string {
	var args []string
	for _, a := range cmd.Args {
		args = append(args, paths.ReplaceHome(a))
	}
	return CmdStyle.Render("  # "+CmdDirStyle.Render(paths.ReplaceHome(cmd.Dir)+":"), CmdArgStyle.Render(args...))
}

var (
	OverrideStyle = style.New(func() lipgloss.Style {
		return lipgloss.NewStyle()
	})
	OverrideCaretColor = lipgloss.Color("1")
	OverrideCaretStyle = style.New(func() lipgloss.Style {
		return lipgloss.NewStyle().Foreground(OverrideCaretColor).Bold(true)
	})
	OverrideTextStyle = style.New(func() lipgloss.Style {
		return lipgloss.NewStyle().Faint(true)
	})
)

func OverrideWarning(t model.Target) string {
	return OverrideStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left,
		OverrideCaretStyle.Render("! "),
		OverrideTextStyle.Render("overridden by "+t.Label()),
	))
}
