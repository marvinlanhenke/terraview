package filter

import (
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/ui"
)

type styles struct {
	header lipgloss.Style
	modal  lipgloss.Style
	row    lipgloss.Style
	rowAlt lipgloss.Style
}

func newStyles(t ui.Theme) styles {
	p := t.Palette

	header := lipgloss.NewStyle().
		Foreground(p.Text).
		Background(p.SurfaceMuted)

	modal := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderBackground(p.SurfaceMuted).
		Background(p.SurfaceMuted).
		Padding(0, 1)

	row := lipgloss.NewStyle().
		Foreground(p.Text).
		Background(p.SurfaceMuted)

	rowAlt := lipgloss.NewStyle().
		Foreground(p.Text).
		Background(p.SurfaceAlt)

	return styles{
		header: header,
		modal:  modal,
		row:    row,
		rowAlt: rowAlt,
	}
}
