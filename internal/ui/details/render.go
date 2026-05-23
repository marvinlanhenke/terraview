package details

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/alecthomas/chroma/v2/quick"
)

var ansiBackground = regexp.MustCompile(`\x1b\[0m`)

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
		jsonStr := getPayloadStr(d.content.Payload)

		lines = append(lines, d.highlightJson(jsonStr))

		return lines
	}

	indent := strings.Repeat(" ", 2)
	beforeIcon := "(−) "
	afterIcon := "(+) "

	for _, cl := range d.changes {
		beforeLine := indent + beforeIcon + getPayloadStr(cl.before)
		afterLine := indent + afterIcon + getPayloadStr(cl.after)

		path := d.styles.subheader.
			Width(d.width).
			Render("attribute: " + cl.path)

		line := lipgloss.JoinVertical(lipgloss.Left, path, beforeLine, afterLine, "")

		lines = append(lines, line)
	}

	return lines
}

func (d *Details) highlightJson(jsonStr string) string {
	var buf bytes.Buffer

	err := quick.Highlight(&buf, jsonStr, "json", "terminal256", "monokai")
	if err != nil {
		return fmt.Sprintf("%v", jsonStr)
	}

	style := lipgloss.NewStyle().Background(d.styles.palette.Surface)
	if d.focus {
		style = lipgloss.NewStyle().Background(d.styles.palette.SurfaceAlt)
	}

	highlighted := style.
		Width(d.width).
		Render(ansiBackground.ReplaceAllString(buf.String(), ""))

	return highlighted
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

func getPayloadStr(v any) string {
	if v == nil {
		return "null"
	}

	b, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return fmt.Sprintf("%v", v)
	}

	return string(b)
}
