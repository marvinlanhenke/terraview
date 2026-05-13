package tree

import (
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/theme"
)

type Tree struct {
	root    *Node
	visible []*Node
	cursor  int
	query   string

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

func (t *Tree) SetRoot(root *Node) {
	t.root = root
	t.rebuildVisible()
	t.clampCursor()
	t.syncViewport()
}

func (t *Tree) SetSize(width, height int) {
	t.width = max(0, width)
	t.height = max(0, height)

	t.viewport.SetWidth(t.width)
	t.viewport.SetHeight(t.height)
	t.syncViewport()
}

func (t *Tree) Selected() *Node {
	if len(t.visible) == 0 {
		return nil
	}

	return t.visible[t.cursor]
}

func (t *Tree) ApplyFilter(query string) {
	t.query = strings.TrimSpace(strings.ToLower(query))
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
			if selected != nil && selected.hasChildren() {
				selected.Expanded = !selected.Expanded
				t.rebuildVisible()
			}

		// Collapse
		case key.Matches(msg, keys.collapse):
			selected := t.Selected()
			if selected == nil {
				break
			}
			if selected.hasChildren() && selected.Expanded {
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
		return t.styles.empty.
			Width(t.width).
			MaxWidth(t.width).
			Height(t.height).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render("Nothing to show...")
	}

	return t.viewport.View()
}

func (t Tree) renderNode(n *Node, selected bool) string {
	indent := strings.Repeat(" ", n.Depth)

	icon := " "
	if n.hasChildren() {
		if n.Expanded || t.query != "" {
			icon = "▾"
		} else {
			icon = "▸"
		}
	}

	actionMarker := t.styles.actionMarker(n.Action)
	rawPrefix := indent + icon + " " + actionMarker.marker + " "

	if selected {
		prefixSelected := actionMarker.style.Background(t.styles.palette.SurfaceAlt).Render(rawPrefix)
		labelSelected := t.styles.labelAlt.Render(n.Label)
		return t.styles.selected.
			Width(t.width).
			MaxWidth(t.width).
			Render(prefixSelected + labelSelected)
	}

	prefix := actionMarker.style.Faint(true).Render(rawPrefix)
	label := t.styles.label.Render(n.Label)

	return t.styles.base.
		Width(t.width).
		MaxWidth(t.width).
		Render(prefix + label)
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
	if len(t.visible) == 0 || t.height <= 0 {
		return
	}

	top := t.viewport.YOffset()
	bottom := top + t.height - 1

	if t.cursor < top {
		t.viewport.SetYOffset(t.cursor)
		return
	}

	if t.cursor > bottom {
		t.viewport.SetYOffset(t.cursor - t.height + 1)
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
		t.walk(child)
	}
}

func (t *Tree) walk(n *Node) {
	if t.query != "" && !matches(n, t.query) && !hasMatchingDescendant(n, t.query) {
		return
	}

	t.visible = append(t.visible, n)

	if n.Expanded || t.query != "" {
		for _, child := range n.Children {
			t.walk(child)
		}
	}
}

func matches(n *Node, query string) bool {
	return strings.Contains(strings.ToLower(n.Label), query) ||
		strings.Contains(strings.ToLower(n.Id), query) ||
		strings.Contains(strings.ToLower(string(n.Action)), query)
}

func hasMatchingDescendant(n *Node, query string) bool {
	for _, child := range n.Children {
		if matches(child, query) || hasMatchingDescendant(child, query) {
			return true
		}
	}

	return false
}
