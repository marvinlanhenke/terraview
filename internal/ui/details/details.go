package details

import (
	"encoding/json"
	"fmt"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/planview"
	"github.com/marvinlanhenke/terraview/internal/ui/theme"
)

type Details struct {
	node     *planview.Node
	changes  []changeLine
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

	// TODO show which mode (summary/plan)
	d.setHeader("▤ Details")

	d.viewport.SetWidth(d.width)
	d.viewport.SetHeight(d.contentHeight())
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

func (d *Details) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, keys.toggle):
			d.showPlan = !d.showPlan
			d.syncViewport()
		}
	}

	d.viewport, cmd = d.viewport.Update(msg)

	return cmd
}

func (d Details) View() string {
	if d.node == nil || d.node.Kind == planview.NodeGroup {
		empty := d.styles.empty.
			Width(d.width).
			MaxWidth(d.width).
			Height(d.height - lipgloss.Height(d.header)).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render("Nothing to show yet...")

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
	lines := make([]string, 0)

	if d.showPlan {
		plan, err := json.MarshalIndent(d.node.Payload, "", " ")
		if err != nil {
			line := "error while marshaling json payload."
			lines = append(lines, line)
			return lines
		}

		lines = append(lines, string(plan))
		return lines
	}

	header := lipgloss.
		NewStyle().
		Width(d.width).
		Render("Changed Attributes:")

	lines = append(lines, header)

	indent := " "
	beforeIcon := "−"
	afterIcon := "+"

	for _, cl := range d.changes {
		beforeLine := indent + beforeIcon + renderValue(cl.before) + "\n"
		afterLine := indent + afterIcon + renderValue(cl.after) + "\n"

		path := lipgloss.
			NewStyle().
			Border(lipgloss.ASCIIBorder(), false, false, true, false).
			Render(cl.path + ":")

		line := path + "\n" + beforeLine + afterLine
		lines = append(lines, line)
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
