package tree

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/plan"
	"github.com/marvinlanhenke/terraview/internal/theme"
)

type Tree struct {
	root    *plan.Node
	visible []*plan.Node
	cursor  int
	query   string
	queryRE *regexp.Regexp
	filters map[plan.Action]bool
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

func (t Tree) GetVisible() int {
	return len(t.visible)
}

func (t *Tree) SetRoot(n *plan.Node) {
	t.root = n
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

func (t *Tree) Selected() *plan.Node {
	if len(t.visible) == 0 {
		return nil
	}

	return t.visible[t.cursor]
}

func (t *Tree) ApplyFilters(f map[plan.Action]bool) {
	t.filters = f
	t.rebuildVisible()
	t.clampCursor()
	t.syncViewport()
}

func (t *Tree) ApplyQuery(query string) {
	t.query = strings.TrimSpace(query)
	t.queryRE = nil

	if isRegexQuery(t.query) {
		re, err := regexp.Compile("(?i)" + unwrapRegex(t.query))
		if err == nil {
			t.queryRE = re
		}
	}

	t.rebuildVisible()
	t.clampCursor()
	t.syncViewport()
}

func (t *Tree) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		// Up
		case key.Matches(msg, keys.up):
			t.cursor--

		// Down
		case key.Matches(msg, keys.down):
			t.cursor++

		// Expand
		case key.Matches(msg, keys.expand):
			selected := t.Selected()
			if selected != nil && selected.HasChildren() {
				selected.Expanded = !selected.Expanded
				t.rebuildVisible()
			}

		// Collapse
		case key.Matches(msg, keys.collapse):
			selected := t.Selected()
			if selected == nil {
				break
			}
			if selected.HasChildren() && selected.Expanded {
				selected.Expanded = false
				t.rebuildVisible()
			} else if selected.Parent != nil {
				for i, n := range t.visible {
					if n == selected.Parent {
						t.cursor = i
						break
					}
				}
			}
		}
	}

	t.clampCursor()
	t.syncViewport()

	return nil
}

func (t Tree) View() string {
	if len(t.visible) == 0 {
		empty := t.styles.empty.
			Width(t.width).
			MaxWidth(t.width).
			Height(t.height - lipgloss.Height(t.header)).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render("Nothing to show...")

		return lipgloss.JoinVertical(lipgloss.Left, t.header, empty)
	}

	return lipgloss.JoinVertical(lipgloss.Left, t.header, t.viewport.View())
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

func (t Tree) renderNode(n *plan.Node, selected bool) string {
	indent := strings.Repeat(" ", max(1, n.Depth))

	icon := " "
	if n.HasChildren() {
		if n.Expanded || t.query != "" {
			icon = "◉"
		} else {
			icon = "○"
		}
	}

	actionMarker := t.styles.actionMarker(n.Action)
	rawPrefix := indent + " " + icon + " "
	wrap := n.Kind == plan.NodeResource

	return t.renderLine(rawPrefix, n.Label, n.LabelCount, actionMarker, selected, wrap)
}

func (t Tree) renderLine(rawPrefix, rawLabel, rawLabelCount string, actionMarker actionStyle, selected, wrap bool) string {
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
	t.visible = t.visible[:0]

	if t.root == nil {
		return
	}

	for _, child := range t.root.Children {
		// We filter out empty group nodes
		// We only show active filters
		if child.HasChildren() && (t.filters[child.Action] || !hasActiveFilters(t.filters)) {
			t.walk(child)
		}
	}
}

func (t *Tree) walk(n *plan.Node) {
	if t.query != "" && !t.matches(n) && !t.hasMatchingDescendant(n) {
		return
	}

	t.visible = append(t.visible, n)

	if n.Expanded || t.query != "" {
		for _, child := range n.Children {
			t.walk(child)
		}
	}
}

func (t *Tree) matches(n *plan.Node) bool {
	return t.matchField(n.Id) ||
		t.matchField(n.Label) ||
		t.matchField(string(n.Action)) ||
		t.matchField(convertPayload(n.Payload))
}

func (t *Tree) matchField(v string) bool {
	if t.queryRE != nil {
		return t.queryRE.MatchString(v)
	}
	return strings.Contains(strings.ToLower(v), strings.ToLower(t.query))
}

func (t *Tree) hasMatchingDescendant(n *plan.Node) bool {
	for _, child := range n.Children {
		if t.matches(child) || t.hasMatchingDescendant(child) {
			return true
		}
	}

	return false
}

func hasActiveFilters(f map[plan.Action]bool) bool {
	for _, isActive := range f {
		if isActive {
			return true
		}
	}

	return false
}

func isRegexQuery(query string) bool {
	if len(query) < 3 {
		// Require at lest /x/, so "/" and "//" stay plain text.
		return false
	}

	if query[0] != '/' || query[len(query)-1] != '/' {
		return false
	}

	return strings.TrimSpace(unwrapRegex(query)) != ""
}

func unwrapRegex(query string) string {
	if len(query) < 3 {
		return ""
	}
	return query[1 : len(query)-1]
}

func convertPayload(v any) string {
	if v == nil {
		return "null"
	}

	switch t := v.(type) {
	case string:
		return t
	default:
		b, err := json.Marshal(t)
		if err != nil {
			return fmt.Sprintf("%v", t)
		}
		return "\n" + string(b)
	}
}
