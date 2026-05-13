package details

import (
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/theme"
)

type styles struct {
	palette    theme.Palette
	base       lipgloss.Style
	background lipgloss.Style
}

func newStyles(t theme.Theme) styles {
	p := t.Palette
	s := t.Styles

	base := lipgloss.NewStyle().Padding(0, 1).Background(p.Surface)

	return styles{
		palette:    p,
		base:       base,
		background: s.Background,
	}
}
