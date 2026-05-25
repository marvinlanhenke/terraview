package tree

import (
	"maps"

	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/ui"
	"github.com/marvinlanhenke/terraview/internal/ui/theme"
)

type NodeKind int

const (
	NodeGroup NodeKind = iota
	NodeResource
)

type Node struct {
	Id         string
	Label      string
	LabelCount string
	Kind       NodeKind
	Action     ui.Action
	Children   []*Node
	Payload    any
	Changes    ui.ChangeSet
}

func (n *Node) HasChildren() bool {
	return len(n.Children) > 0
}

func (n *Node) IsResource() bool {
	return n != nil && n.Kind == NodeResource && n.Action != ui.ActionNoOp
}

func (n *Node) IsError() bool {
	return n != nil && n.Kind == NodeResource && n.Action == ui.ActionError
}

type Tree struct {
	root     *Node
	rows     []row
	expanded map[string]bool
	cursor   int

	filters map[ui.Action]bool
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

func (t *Tree) SetRoot(n *Node) {
	t.root = n
	t.expanded = rebaseExpanded(n, t.expanded)

	t.rebuildRows()
	t.clampCursor()
	t.syncViewport()
}

func (t *Tree) SetCriteria(query string, filters map[ui.Action]bool) {
	t.matcher = newMatcher(query)

	t.filters = maps.Clone(filters)
	if t.filters == nil {
		t.filters = make(map[ui.Action]bool)
	}

	t.rebuildRows()
	t.clampCursor()
	t.syncViewport()
}

func (t *Tree) SetSize(width, height int) {
	t.width = max(0, width)
	t.height = max(0, height)

	t.header = lipgloss.NewStyle().
		Width(t.width).
		Render("⌘ Resources")

	contentHeight := max(0, t.height-lipgloss.Height(t.header))
	t.viewport.SetHeight(contentHeight)
	t.viewport.SetWidth(t.width)
	t.syncViewport()
}

func (t *Tree) Selected() *Node {
	if len(t.rows) == 0 {
		return nil
	}

	return t.rows[t.cursor].node
}

func (t *Tree) VisibleResourceCount() int {
	var count int

	for _, row := range t.rows {
		if row.node.Kind == NodeResource {
			count++
		}
	}

	return count
}

func (t *Tree) selectedRow() (row, bool) {
	if len(t.rows) == 0 {
		return row{}, false
	}

	return t.rows[t.cursor], true
}

func (t *Tree) setExpanded(id string, expanded bool) {
	if expanded {
		t.expanded[id] = true
		return
	}

	delete(t.expanded, id)
}

func (t *Tree) rebuildRows() {
	t.rows = buildRows(t.root, t.expanded, criteria{
		matcher: t.matcher,
		filters: t.filters,
	})
}

func rebaseExpanded(root *Node, previous map[string]bool) map[string]bool {
	next := make(map[string]bool)

	if root == nil || len(previous) == 0 {
		return next
	}

	var visit func(*Node)

	visit = func(n *Node) {
		if n == nil {
			return
		}

		if previous[n.Id] {
			next[n.Id] = true
		}

		for _, child := range n.Children {
			visit(child)
		}
	}

	visit(root)

	return next
}
