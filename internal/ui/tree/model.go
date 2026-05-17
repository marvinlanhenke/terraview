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
	root    *planview.Node
	visible []*planview.Node
	cursor  int
	filters map[planview.Action]bool
	matcher matcher
	header  string

	width    int
	height   int
	viewport viewport.Model
	styles   styles
}

func New(t theme.Theme) Tree {
	s := newStyles(t)

	vp := viewport.New()
	vp.FillHeight = true
	vp.Style = s.background

	return Tree{
		viewport: vp,
		styles:   s,
	}
}

func (t *Tree) SetRoot(n *planview.Node) {
	t.root = n
	t.rebuildVisible()
	t.clampCursor()
	t.syncViewport()
}

func (t *Tree) SetCriteria(query string, filters map[planview.Action]bool) {
	t.matcher = newMatcher(query)

	t.filters = maps.Clone(filters)
	if t.filters == nil {
		t.filters = make(map[planview.Action]bool)
	}

	t.rebuildVisible()
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

func (t *Tree) Selected() *planview.Node {
	if len(t.visible) == 0 {
		return nil
	}

	return t.visible[t.cursor]
}

func (t *Tree) VisibleCount() int {
	return len(t.visible)
}

func (t Tree) contentHeight() int {
	return max(0, t.height-lipgloss.Height(t.header))
}

func (t *Tree) setHeader(text string) {
	t.header = lipgloss.
		NewStyle().
		Width(t.width).
		Render(text)
}

func (t *Tree) renderNode(n *planview.Node, selected bool) string {
	indent := strings.Repeat(" ", max(1, n.Depth))

	icon := " "
	if n.HasChildren() {
		if n.Expanded || t.matcher.Active() {
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

func (t *Tree) syncViewport() {
	if t.width <= 0 || t.height <= 0 {
		t.viewport.SetContentLines(nil)
		return
	}

	if len(t.visible) == 0 {
		t.viewport.SetContentLines(nil)
		t.viewport.SetYOffset(0)
		return
	}

	lines := make([]string, len(t.visible))
	for i, n := range t.visible {
		selected := i == t.cursor
		lines[i] = t.renderNode(n, selected)
	}

	t.viewport.SetContentLines(lines)
	t.keepCursorVisible()
}

func (t *Tree) keepCursorVisible() {
	h := t.viewport.Height()

	if len(t.visible) == 0 || h <= 0 {
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
	if len(t.visible) == 0 {
		t.cursor = 0
		return
	}

	if t.cursor < 0 {
		t.cursor = 0
	}

	if t.cursor >= len(t.visible) {
		t.cursor = len(t.visible) - 1
	}
}

func (t *Tree) rebuildVisible() {
	t.visible = buildVisible(t.root, criteria{
		matcher: t.matcher,
		filters: t.filters,
	})
}
