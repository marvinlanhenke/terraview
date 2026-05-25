package search

import (
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/ui"
)

type styles struct {
	palette         ui.Palette
	background      lipgloss.Style
	backgroundMuted lipgloss.Style
	nugget          lipgloss.Style
	status          lipgloss.Style
	input           lipgloss.Style
	inputAlt        lipgloss.Style
	banner          lipgloss.Style
}

func newStyles(t ui.Theme) styles {
	p := t.Palette

	background := lipgloss.NewStyle().Background(p.Surface)
	backgroundMuted := background.Faint(true)

	base := lipgloss.NewStyle().Padding(0, 1)

	input := base.
		Foreground(p.TextMuted).
		Background(p.Surface)

	inputAlt := base.
		Foreground(p.Text).
		Background(p.SurfaceAlt)

	status := base.
		Foreground(p.TextMuted).
		Background(p.Primary)

	nugget := status.Bold(true)

	banner := base.
		Background(p.Surface).
		Foreground(p.Info)

	return styles{
		palette:         p,
		background:      background,
		backgroundMuted: backgroundMuted,
		nugget:          nugget,
		status:          status,
		input:           input,
		inputAlt:        inputAlt,
		banner:          banner,
	}
}
