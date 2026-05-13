package details

import (
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/marvinlanhenke/terraview/internal/theme"
)

type Details struct {
	width    int
	height   int
	viewport viewport.Model
	styles   styles
}

func New(t theme.Theme) Details {
	s := newStyles(t)

	vp := viewport.New()
	vp.FillHeight = true
	vp.Style = s.background

	return Details{
		viewport: vp,
		styles:   s,
	}
}

func (d *Details) SetSize(width, height int) {
	d.width = max(0, width)
	d.height = max(0, height)

	d.viewport.SetWidth(d.width)
	d.viewport.SetHeight(d.height)
}

func (d Details) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (d Details) View() string {
	return d.viewport.View()
}
