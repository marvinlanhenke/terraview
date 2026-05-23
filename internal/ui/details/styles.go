package details

import (
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/ui/theme"
)

type styles struct {
	palette       theme.Palette
	base          lipgloss.Style
	empty         lipgloss.Style
	header        lipgloss.Style
	subheader     lipgloss.Style
	background    lipgloss.Style
	backgroundAlt lipgloss.Style
}

func newStyles(t theme.Theme) styles {
	p := t.Palette
	s := t.Styles

	base := lipgloss.NewStyle().Padding(0, 1).Background(p.Surface)

	header := lipgloss.NewStyle()

	subheader := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true, false, true, false).
		BorderForeground(p.TextMuted).
		Padding(0, 1)

	return styles{
		palette:       p,
		base:          base,
		empty:         base.Faint(true),
		header:        header,
		subheader:     subheader,
		background:    s.Background,
		backgroundAlt: s.BackgroundAlt,
	}
}
