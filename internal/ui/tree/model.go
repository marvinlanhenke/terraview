package tree

import (
	"maps"

	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/ui"
)

// NodeKind distinguishes between the two structural roles a node can play in
// the tree: a group (container) or a leaf resource.
type NodeKind int

const (
	// NodeGroup represents a container node that groups one or more resources.
	NodeGroup NodeKind = iota
	// NodeResource represents a leaf node corresponding to a single Terraform resource.
	NodeResource
)

// Node is a single element in the resource tree. It carries the display
// metadata, its planned action, an optional arbitrary payload used for
// full-text search, and any child nodes that make up its subtree.
type Node struct {
	Id         string
	Label      string
	LabelCount string
	Kind       NodeKind
	Action     ui.Action
	Children   []*Node
	Payload    any
	Changes    ui.ChangeSet

	// searchPayload is a pre-computed, searchable string representation of
	// Payload, populated by prepareSearchPayloads.
	searchPayload string
}

// HasChildren reports whether the node has at least one child.
func (n *Node) HasChildren() bool {
	return len(n.Children) > 0
}

// IsInspectable reports whether the node represents a resource whose details
// can be inspected (i.e. it is a resource node with a non-no-op action).
func (n *Node) IsInspectable() bool {
	return n != nil && n.Kind == NodeResource && n.Action != ui.ActionNoOp
}

// IsError reports whether the node represents a resource that encountered an
// error during planning.
func (n *Node) IsError() bool {
	return n != nil && n.Kind == NodeResource && n.Action == ui.ActionError
}

// Tree is the Bubble Tea component that renders a scrollable, filterable
// resource tree. It manages cursor position, expand/collapse state, active
// search/filter criteria, and delegates rendering to an embedded viewport.
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

// New returns an initialised Tree styled with the provided theme.
func New(t ui.Theme) Tree {
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

// SetRoot replaces the root node of the tree, carries over the previous
// expand state for any node ids that still exist, and refreshes the viewport.
func (t *Tree) SetRoot(n *Node) {
	t.root = n
	prepareSearchPayloads(t.root)

	t.expanded = rebaseExpanded(n, t.expanded)

	t.rebuildRows()
	t.clampCursor()
	t.syncViewport()
}

// SetCriteria applies a search query and a set of action filters to the tree.
// Rows are rebuilt and the cursor is re-clamped so the selection stays valid.
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

// SetSize updates the available rendering area and adjusts the inner viewport
// height to account for the fixed header.
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

// Selected returns the node currently under the cursor, or nil when the tree
// is empty.
func (t *Tree) Selected() *Node {
	if len(t.rows) == 0 {
		return nil
	}

	return t.rows[t.cursor].node
}

// VisibleResourceCount returns the number of resource (leaf) nodes currently
// present in the rendered row list, respecting active filters and search.
func (t *Tree) VisibleResourceCount() int {
	var count int

	for _, row := range t.rows {
		if row.node.Kind == NodeResource {
			count++
		}
	}

	return count
}

// selectedRow returns the row at the current cursor position together with a
// boolean indicating whether a valid row was found.
func (t *Tree) selectedRow() (row, bool) {
	if len(t.rows) == 0 {
		return row{}, false
	}

	return t.rows[t.cursor], true
}

// setExpanded marks the node with the given id as expanded or collapsed in the
// internal expand map.
func (t *Tree) setExpanded(id string, expanded bool) {
	if expanded {
		t.expanded[id] = true
		return
	}

	delete(t.expanded, id)
}

// rebuildRows regenerates the flat row slice from the current root node
// applying the active search matcher and action filters.
func (t *Tree) rebuildRows() {
	t.rows = buildRows(t.root, t.expanded, criteria{
		matcher: t.matcher,
		filters: t.filters,
	})
}

// rebaseExpanded computes a new expand map that contains only the node ids
// present in the updated tree, preserving the expanded state carried over from
// the previous map. This prevents stale ids from accumulating over time.
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
