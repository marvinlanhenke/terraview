package details

import (
	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/ui"
)

type Kind int

const (
	KindNone Kind = iota
	KindGroup
	KindResource
)

type Content struct {
	Key     string
	Kind    Kind
	Label   string
	Changes ui.ChangeSet
	Payload any
	IsError bool
}

type Details struct {
	content  Content
	changes  []change
	header   string
	focus    bool
	showPlan bool

	width    int
	height   int
	viewport viewport.Model
	styles   styles
}

func New(t ui.Theme) Details {
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

	d.setHeader()

	contentHeight := max(0, d.height-lipgloss.Height(d.header))
	d.viewport.SetHeight(contentHeight)
	d.viewport.SetWidth(d.width)
	d.syncViewport()
}

func (d *Details) SetContent(content Content) {
	changed := d.content.Key != content.Key

	d.content = content
	d.changes = flattenChanges(content.Changes)
	d.showPlan = content.IsError

	if changed {
		d.viewport.SetYOffset(0)
	}

	d.setHeader()
	d.syncViewport()
}

func (d *Details) Focus() {
	d.focus = true
	d.viewport.Style = d.styles.backgroundAlt
	d.syncViewport()
}

func (d *Details) Blur() {
	d.focus = false
	d.viewport.Style = d.styles.background
	d.syncViewport()
}

func (d *Details) setHeader() {
	label := "▤ Details · Diff"

	if d.showPlan && d.content.Payload != nil {
		label = "▤ Details · Plan"
	}

	d.header = d.styles.header.
		Width(d.width).
		Render(label)
}
