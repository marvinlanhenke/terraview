package details

import (
	"encoding/json"
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

func (d *Details) syncViewport() {
	if d.width <= 0 || d.height <= 0 {
		d.viewport.SetContentLines(nil)
		return
	}

	if d.showEmptyState() {
		d.viewport.SetContentLines(nil)
		d.viewport.SetYOffset(0)
		return
	}

	lines := d.renderLines()
	d.viewport.SetContentLines(lines)
}

func (d *Details) renderLines() []string {
	lines := make([]string, 0)

	lines = append(lines, "")

	if d.showingPlan() {
		lines = append(lines, formatPayload(d.content.Payload))
		return lines
	}

	indent := strings.Repeat(" ", 2)
	beforeIcon := "(−) "
	afterIcon := "(+) "

	for _, cl := range d.changes {
		beforeLine := indent + beforeIcon + formatPayload(cl.before)
		afterLine := indent + afterIcon + formatPayload(cl.after)

		path := d.styles.subheader.
			Width(d.width).
			Render("attribute: " + cl.path)

		line := lipgloss.JoinVertical(lipgloss.Left, path, beforeLine, afterLine, "")

		lines = append(lines, line)
	}

	return lines
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

func formatPayload(v any) string {
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
		return string(b)
	}
}
