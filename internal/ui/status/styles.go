package status

import (
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/ui"
)

type styles struct {
	base      lipgloss.Style
	borderBar lipgloss.Style
}

// newStyles builds the styles used by Status.
func newStyles(t ui.Theme) styles {
	p := t.Palette

	base := lipgloss.NewStyle().
		Padding(0, 1).
		Background(p.Surface)

	borderBar := base.
		Foreground(p.Text).
		Border(lipgloss.NormalBorder(), true, false, true, false).
		BorderForeground(p.Secondary)

	return styles{
		base:      base,
		borderBar: borderBar,
	}
}
