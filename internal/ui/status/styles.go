package status

import (
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/ui/theme"
)

type styles struct {
	palette   theme.Palette
	base      lipgloss.Style
	borderBar lipgloss.Style
}

func newStyles(t theme.Theme) styles {
	p := t.Palette

	base := lipgloss.NewStyle().
		Padding(0, 1).
		Background(p.Surface)

	borderBar := base.
		Foreground(p.Text).
		Border(lipgloss.NormalBorder(), true, false, true, false).
		BorderForeground(p.Secondary)

	return styles{
		palette:   p,
		base:      base,
		borderBar: borderBar,
	}
}
