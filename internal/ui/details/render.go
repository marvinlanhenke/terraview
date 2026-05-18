package details

import (
	"encoding/json"
	"fmt"

	"charm.land/lipgloss/v2"
)

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
