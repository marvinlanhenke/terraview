// Package details renders and updates the selected resource details pane.
package details

import (
	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/ui"
)

// Kind identifies the type of content shown in the details pane.
type Kind int

const (
	// KindNone represents an empty details selection.
	KindNone Kind = iota
	// KindGroup represents a selected resource group.
	KindGroup
	// KindResource represents a selected resource.
	KindResource
)

// Content contains the data rendered by the details pane.
type Content struct {
	Key     string
	Kind    Kind
	Label   string
	Changes ui.ChangeSet
	Payload any
	IsError bool
}

// Details renders and updates the selected resource details pane.
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

// New returns an initialized Details using t for styling.
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

// SetSize sets the rendered details pane dimensions.
func (d *Details) SetSize(width, height int) {
	d.width = max(0, width)
	d.height = max(0, height)

	d.setHeader()

	contentHeight := max(0, d.height-lipgloss.Height(d.header))
	d.viewport.SetHeight(contentHeight)
	d.viewport.SetWidth(d.width)
	d.syncViewport()
}

// SetContent replaces the content shown by Details.
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

// Focus moves keyboard focus to the details pane.
func (d *Details) Focus() {
	d.focus = true
	d.viewport.Style = d.styles.backgroundAlt
	d.syncViewport()
}

// Blur removes keyboard focus from the details pane.
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

func (d *Details) showEmptyState() bool {
	return !d.hasViewportContent()
}

func (d *Details) hasViewportContent() bool {
	return d.showingPlan() || d.hasDiffContent()
}

func (d *Details) showingPlan() bool {
	return d.showPlan && d.content.Payload != nil
}

func (d *Details) hasDiffContent() bool {
	return len(d.changes) > 0
}

func (d *Details) emptyStateMessage() string {
	switch d.content.Kind {
	case KindNone, KindGroup:
		return "Select a resource to inspect changes."
	case KindResource:
		return "No changed attributes for this resource."
	default:
		return "Nothing to show yet..."
	}
}
