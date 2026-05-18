package details

import (
	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/planview"
	"github.com/marvinlanhenke/terraview/internal/ui/theme"
)

type Details struct {
	node     *planview.Node
	changes  []change
	header   string
	showPlan bool

	width    int
	height   int
	viewport viewport.Model
	styles   styles
}

func New(t theme.Theme) Details {
	s := newStyles(t)

	vp := viewport.New()
	vp.FillHeight = true
	vp.SoftWrap = true
	vp.Style = s.background

	return Details{
		viewport: vp,
		styles:   s,
		showPlan: false,
	}
}

func (d *Details) SetSize(width, height int) {
	d.width = max(0, width)
	d.height = max(0, height)

	d.header = lipgloss.NewStyle().
		Width(d.width).
		Render("▤ Details")

	contentHeight := max(0, d.height-lipgloss.Height(d.header))
	d.viewport.SetHeight(contentHeight)
	d.viewport.SetWidth(d.width)
	d.syncViewport()
}

func (d *Details) SetNode(n *planview.Node) {
	hasChanged := d.node != n
	d.node = n

	if n == nil {
		d.changes = nil
		d.viewport.SetYOffset(0)
		d.syncViewport()
		return
	}

	d.changes = flattenChanges(n)

	if hasChanged {
		d.viewport.SetYOffset(0)
	}

	d.syncViewport()
}

func (d *Details) Focus() {
	d.viewport.Style = d.styles.backgroundAlt
}

func (d *Details) Blur() {
	d.viewport.Style = d.styles.background
}
