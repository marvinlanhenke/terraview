package details

import (
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/theme"
)

type styles struct {
	palette       theme.Palette
	base          lipgloss.Style
	empty         lipgloss.Style
	background    lipgloss.Style
	backgroundAlt lipgloss.Style
	header        lipgloss.Style
}

func newStyles(t theme.Theme) styles {
	p := t.Palette
	s := t.Styles

	base := lipgloss.NewStyle().Padding(0, 1).Background(p.Surface)
	header := lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.ASCIIBorder(), true, false)

	return styles{
		palette:       p,
		base:          base,
		empty:         base.Faint(true),
		background:    s.Background,
		backgroundAlt: s.BackgroundAlt,
		header:        header,
	}
}
