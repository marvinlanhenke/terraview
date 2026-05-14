package details

import (
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/components/tree"
	"github.com/marvinlanhenke/terraview/internal/theme"
)

type Details struct {
	node     *tree.Node
	summary  *summary
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
	d.syncViewport()
}

func (d *Details) SetNode(n *tree.Node) {
	d.node = n
	d.summary = newSummary(d.node)
	d.syncViewport()
}

func (d *Details) Focus() {
	d.viewport.Style = d.styles.backgroundAlt
}

func (d *Details) Blur() {
	d.viewport.Style = d.styles.background
}

func (d *Details) Update(msg tea.Msg) tea.Cmd {
	d.syncViewport()

	return nil
}

func (d Details) View() string {
	if d.node == nil {
		return d.styles.empty.
			Width(d.width).
			MaxWidth(d.width).
			Height(d.height).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render("Nothing to show...")
	}

	return d.viewport.View()
}

func (d *Details) syncViewport() {
	if d.width <= 0 || d.height <= 0 {
		d.viewport.SetContentLines(nil)
		return
	}

	if d.node == nil {
		d.viewport.SetContentLines(nil)
		d.viewport.SetYOffset(0)
		return
	}

	header := lipgloss.
		NewStyle().
		Width(d.width).
		Border(lipgloss.ASCIIBorder(), true, false).
		Render(d.summary.header)

	// TODO renderFunction
	lines := make([]string, 2)
	lines[0] = header
	lines[1] = d.node.Label + " " + string(d.node.Action)

	d.viewport.SetContentLines(lines)
}
