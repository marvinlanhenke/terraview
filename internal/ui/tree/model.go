package tree

import (
	"maps"
	"strings"

	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/planview"
	"github.com/marvinlanhenke/terraview/internal/ui/theme"
)

type Tree struct {
	root     *planview.Node
	rows     []row
	expanded map[string]bool
	cursor   int

	filters map[planview.Action]bool
	matcher matcher

	width    int
	height   int
	header   string
	viewport viewport.Model
	styles   styles
}

func New(t theme.Theme) Tree {
	s := newStyles(t)
	expanded := make(map[string]bool)

	vp := viewport.New()
	vp.FillHeight = true
	vp.Style = s.background

	return Tree{
		expanded: expanded,
		viewport: vp,
		styles:   s,
	}
}

func (t *Tree) SetRoot(n *planview.Node) {
	t.root = n
	t.rebuildRows()
	t.clampCursor()
	t.syncViewport()
}

func (t *Tree) SetCriteria(query string, filters map[planview.Action]bool) {
	t.matcher = newMatcher(query)

	t.filters = maps.Clone(filters)
	if t.filters == nil {
		t.filters = make(map[planview.Action]bool)
	}

	t.rebuildRows()
	t.clampCursor()
	t.syncViewport()
}

func (t *Tree) SetSize(width, height int) {
	t.width = max(0, width)
	t.height = max(0, height)

	t.setHeader("⌘ Resources")

	t.viewport.SetWidth(t.width)
	t.viewport.SetHeight(t.contentHeight())
	t.syncViewport()
}

// TODO do we need to expose full node?
func (t *Tree) Selected() *planview.Node {
	if len(t.rows) == 0 {
		return nil
	}

	return t.rows[t.cursor].node
}

func (t *Tree) VisibleCount() int {
	return len(t.rows)
}

func (t *Tree) contentHeight() int {
	return max(0, t.height-lipgloss.Height(t.header))
}

func (t *Tree) selectedRow() (row, bool) {
	if len(t.rows) == 0 {
		return row{}, false
	}

	return t.rows[t.cursor], true
}

func (t *Tree) setHeader(text string) {
	t.header = lipgloss.
		NewStyle().
		Width(t.width).
		Render(text)
}

func (t *Tree) setExpanded(id string, expanded bool) {
	if expanded {
		t.expanded[id] = true
		return
	}

	delete(t.expanded, id)
}

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
		if r.open(t.matcher.Active()) {
			icon = "◉"
		} else {
			icon = "○"
		}
	}

	actionMarker := t.styles.actionMarker(n.Action)
	rawPrefix := indent + " " + icon + " "
	wrap := n.Kind == planview.NodeResource

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

func (t *Tree) rebuildRows() {
	t.rows = buildRows(t.root, t.expanded, criteria{
		matcher: t.matcher,
		filters: t.filters,
	})
}
