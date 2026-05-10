package summary

import (
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/theme"
)

var (
	summaryBar = lipgloss.NewStyle().
		Padding(0, 1).
		Foreground(theme.TextFocus).
		Background(theme.BackgroundBlur).
		Border(lipgloss.NormalBorder(), true, false, true, false).
		BorderForeground(theme.AccentSecondary)
)
