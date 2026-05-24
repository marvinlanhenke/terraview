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

	bg := lipgloss.NewStyle().Background(d.styles.palette.Surface)

	if d.showingPlan() {
		bg = lipgloss.NewStyle().Background(d.styles.palette.Surface)

		if d.focus {
			bg = lipgloss.NewStyle().Background(d.styles.palette.SurfaceAlt)
		}

		jsonStr := getJsonStr(d.content.Payload, strings.Repeat(" ", 2))

		lines = append(lines, d.highlightJson(jsonStr, bg))

		return lines
	}

	indent := strings.Repeat(" ", 2)
	beforeIcon := indent + "(−) before:"
	afterIcon := indent + "(+) after:"

	if d.focus {
		bg = lipgloss.NewStyle().Background(d.styles.palette.SurfaceEmbedded)
	}

	for _, cl := range d.changes {
		prefixIndent := indent + indent + indent

		beforeStr := getJsonStr(cl.before, prefixIndent)
		beforeLine := "\n" + beforeIcon + "\n" + d.highlightJson(beforeStr, bg) + "\n"
		beforeLine = bg.Render(beforeLine)

		afterStr := getJsonStr(cl.after, prefixIndent)
		afterLine := afterIcon + "\n" + d.highlightJson(afterStr, bg) + "\n"
		afterLine = bg.Render(afterLine)

		path := d.styles.subheader.
			Width(d.width).
			Render("attribute: " + cl.path)

		line := lipgloss.JoinVertical(lipgloss.Left, path, beforeLine, afterLine)

		lines = append(lines, line)
	}

	return lines
}

func (d *Details) highlightJson(jsonStr string, bg lipgloss.Style) string {
	var buf bytes.Buffer

	err := quick.Highlight(&buf, jsonStr, "json", "terminal256", "catppuccin-macchiato")
	if err != nil {
		return fmt.Sprintf("%v", jsonStr)
	}

	highlighted := bg.
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

func getJsonStr(v any, prefix string) string {
	if v == nil {
		return prefix + "null"
	}

	b, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return fmt.Sprintf("%v", v)
	}

	return prefix + strings.ReplaceAll(string(b), "\n", "\n"+prefix)
}
