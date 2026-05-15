package details

import (
	"encoding/json"
	"fmt"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/components/tree"
	"github.com/marvinlanhenke/terraview/internal/theme"
)

type Details struct {
	node    *tree.Node
	changes []changeLine
	header  string

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

	d.setHeader("▤ Details")

	d.viewport.SetWidth(d.width)
	d.viewport.SetHeight(d.contentHeight())
	d.syncViewport()
}

func (d *Details) SetNode(n *tree.Node) {
	d.node = n
	d.changes = flattenChanges(d.node)
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
		empty := d.styles.empty.
			Width(d.width).
			MaxWidth(d.width).
			Height(d.height).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render("Nothing to show...")

		return lipgloss.JoinVertical(lipgloss.Left, d.header, empty)
	}

	return lipgloss.JoinVertical(lipgloss.Left, d.header, d.viewport.View())
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

	lines := d.renderLines()

	d.viewport.SetContentLines(lines)
}

func (d Details) contentHeight() int {
	return max(0, d.height-lipgloss.Height(d.header))
}

func (d *Details) setHeader(text string) {
	d.header = lipgloss.
		NewStyle().
		Width(d.width).
		Render(text)
}

func (d *Details) renderLines() []string {
	header := lipgloss.
		NewStyle().
		Width(d.width).
		Border(lipgloss.ASCIIBorder(), true, false).
		Render("Changed Attributes:")

	indent := " "
	beforeIcon := "−"
	afterIcon := "+"

	lines := make([]string, len(d.changes)+1)

	lines[0] = header
	for i, cl := range d.changes {
		beforeLine := indent + beforeIcon + renderValue(cl.before) + "\n"
		afterLine := indent + afterIcon + renderValue(cl.after) + "\n"

		path := lipgloss.
			NewStyle().
			Border(lipgloss.ASCIIBorder(), false, false, true, false).
			Render(cl.path + ":")

		lines[i+1] = path + "\n" + beforeLine + afterLine
	}

	return lines
}

func renderValue(v any) string {
	if v == nil {
		return "null"
	}

	switch t := v.(type) {
	case string:
		return t
	default:
		b, err := json.MarshalIndent(t, "", " ")
		if err != nil {
			return fmt.Sprintf("%v", t)
		}
		return "\n" + string(b)
	}
}
