package tree

import (
	"strings"

	"charm.land/lipgloss/v2"
)

func (t *Tree) syncViewport() {
	if t.width <= 0 || t.height <= 0 {
		t.viewport.SetContentLines(nil)
		return
	}

	if len(t.rows) == 0 {
		t.viewport.SetContentLines(nil)
		t.viewport.SetYOffset(0)
		return
	}

	lines := make([]string, len(t.rows))
	for i, r := range t.rows {
		lines[i] = t.renderRow(r, i == t.cursor)
	}

	t.viewport.SetContentLines(lines)
	t.keepCursorVisible()
}

func (t *Tree) renderRow(r row, selected bool) string {
	n := r.node
	indent := strings.Repeat(" ", r.depth+1)

	icon := " "
	if r.expandable {
		if r.open(t.matcher.active()) {
			icon = "◉"
		} else {
			icon = "○"
		}
	}

	actionMarker := t.styles.actionMarker(n.Action)
	rawPrefix := indent + " " + icon + " "
	wrap := n.Kind == NodeResource

	return t.renderLine(rawPrefix, n.Label, n.LabelCount, actionMarker, selected, wrap)
}

func (t *Tree) renderLine(rawPrefix, rawLabel, rawLabelCount string, actionMarker actionStyle, selected, wrap bool) string {
	actionBackground := t.styles.palette.Surface
	labelStyle := t.styles.label
	labelCountStyle := labelStyle.Faint(true)
	style := t.styles.base

	if selected {
		actionBackground = t.styles.palette.SurfaceAlt
		labelStyle = t.styles.labelAlt
		labelCountStyle = labelStyle.Faint(true)
		style = t.styles.selected
	}

	prefix := labelStyle.Render(rawPrefix)
	action := actionMarker.style.Background(actionBackground).Render(actionMarker.marker + " ")
	label := labelStyle.Render(rawLabel)
	left := prefix + action + label

	width := t.width - style.GetHorizontalPadding()

	if wrap {
		return style.Width(t.width).MaxWidth(t.width).Render(lipgloss.Wrap(left, width, "."))
	}

	right := labelCountStyle.Render(rawLabelCount)

	gapSize := max(0, width-lipgloss.Width(left)-lipgloss.Width(right))
	gap := labelStyle.Render(strings.Repeat(" ", gapSize))

	line := left + gap + right

	return style.Width(t.width).MaxWidth(t.width).Render(line)
}

func (t *Tree) keepCursorVisible() {
	h := t.viewport.Height()

	if len(t.rows) == 0 || h <= 0 {
		return
	}

	top := t.viewport.YOffset()
	bottom := top + h - 1

	if t.cursor < top {
		t.viewport.SetYOffset(t.cursor)
		return
	}

	if t.cursor > bottom {
		t.viewport.SetYOffset(t.cursor - h + 1)
	}
}

func (t *Tree) clampCursor() {
	if len(t.rows) == 0 {
		t.cursor = 0
		return
	}

	if t.cursor < 0 {
		t.cursor = 0
	}

	if t.cursor >= len(t.rows) {
		t.cursor = len(t.rows) - 1
	}
}
