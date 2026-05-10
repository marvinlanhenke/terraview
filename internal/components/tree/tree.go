package tree

import (
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Tree struct {
	root    *Node
	visible []*Node
	cursor  int
	query   string
}

func New() Tree {
	return Tree{}
}

func (t *Tree) SetRoot(root *Node) {
	t.root = root
	t.rebuildVisible()
	t.clampCursor()
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

	return nil
}

func (t Tree) View(width, height int) string {
	if len(t.visible) == 0 {
		return treeEmpty.
			Width(width).
			MaxWidth(width).
			Height(height).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render("Nothing to show...")
	}

	var builder strings.Builder

	limit := min(len(t.visible), height)

	for i := range limit {
		n := t.visible[i]
		selected := i == t.cursor

		line := t.renderNode(n, width, selected)
		builder.WriteString(line)

		if i < limit-1 {
			builder.WriteRune('\n')
		}
	}

	lines := builder.String()

	return treeBackground.
		Height(height).
		Render(lines)
}

func (t Tree) renderNode(n *Node, width int, selected bool) string {
	indent := strings.Repeat(" ", n.Depth)

	icon := " "

	if n.HasChildren() {
		if n.Expanded || t.query != "" {
			icon = "▾"
		} else {
			icon = "▸"
		}
	}

	action, style := treeActionMarkerWithStyle(n.Action)
	line := indent + icon + " " + action + " " + n.Label

	if selected {
		style = treeSelected
	}

	return style.Width(width).MaxWidth(width).Render(line)
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
		strings.Contains(strings.ToLower(n.ID), query) ||
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
