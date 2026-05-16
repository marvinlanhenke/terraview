package filter

import (
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/theme"
)

type styles struct {
	palette theme.Palette
	header  lipgloss.Style
	modal   lipgloss.Style
	row     lipgloss.Style
	rowAlt  lipgloss.Style
}

func newStyles(t theme.Theme) styles {
	p := t.Palette

	header := lipgloss.NewStyle().
		Foreground(p.Text).
		Background(p.Surface)

	modal := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderBackground(p.Surface).
		Background(p.Surface).
		Padding(0, 1)

	row := lipgloss.NewStyle().
		Foreground(p.Text).
		Background(p.Surface)

	return styles{
		palette: p,
		header:  header,
		modal:   modal,
		row:     row,
		rowAlt:  row,
	}
}
