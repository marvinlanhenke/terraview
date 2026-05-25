package tree

import "github.com/marvinlanhenke/terraview/internal/ui"

// criteria bundles the active search matcher and action-filter map that are
// used together when building the visible row list.
type criteria struct {
	matcher matcher
	filters map[ui.Action]bool
}

// row is a single rendered entry in the flat list that backs the viewport. It
// stores a pointer to its source node along with the display metadata needed
// to render indentation, expand/collapse icons, and cursor highlighting.
type row struct {
	node       *Node
	depth      int
	parent     int
	expandable bool
	expanded   bool
}

// open reports whether this row's children should be shown. Children are
// always shown when a search is active so that matching descendants are
// visible regardless of the manual expand state.
func (r row) open(searchActive bool) bool {
	return r.expanded || searchActive
}

// buildRows converts the tree rooted at root into a flat, ordered slice of
// rows that can be rendered line-by-line in the viewport. Only root-level
// children that pass the active action filters are included; deeper nodes are
// filtered by the search matcher inside appendRows.
func buildRows(root *Node, expanded map[string]bool, c criteria) []row {
	if root == nil {
		return nil
	}

	filtering := hasActiveFilters(c.filters)
	rows := make([]row, 0, len(root.Children))

	for _, child := range root.Children {
		if !includeRootChild(child, c.filters, filtering) {
			continue
		}

		rows = appendRows(rows, child, -1, 0, expanded, c.matcher)
	}

	return rows
}

// includeRootChild decides whether a top-level group node should appear in the
// row list. A group is excluded when it has no children, or when action
// filtering is active and the group's action is not in the allowed set.
func includeRootChild(n *Node, filters map[ui.Action]bool, filtering bool) bool {
	if n == nil {
		return false
	}

	if !n.HasChildren() {
		return false
	}

	if !filtering {
		return true
	}

	return filters[n.Action]
}

// appendRows recursively appends n and its visible descendants to rows,
// respecting the active search matcher and expand state. parentIndex is the
// slice index of the parent row (-1 for top-level nodes) and depth controls
// indentation.
func appendRows(rows []row, n *Node, parentIndex, depth int, expanded map[string]bool, m matcher) []row {
	if n == nil {
		return rows
	}

	if m.active() && !subtreeMatches(n, m) {
		return rows
	}

	isExpanded := expanded[n.Id]
	rowIndex := len(rows)

	r := row{
		node:       n,
		depth:      depth,
		parent:     parentIndex,
		expandable: n.HasChildren(),
		expanded:   isExpanded,
	}

	rows = append(rows, r)

	if r.open(m.active()) {
		for _, child := range n.Children {
			rows = appendRows(rows, child, rowIndex, depth+1, expanded, m)
		}
	}

	return rows
}

// subtreeMatches reports whether n itself or any node in its subtree satisfies
// the matcher. This is used to keep ancestor nodes visible when a descendant
// matches the current search query.
func subtreeMatches(n *Node, m matcher) bool {
	if n == nil {
		return false
	}

	if m.matchNode(n) {
		return true
	}

	for _, child := range n.Children {
		if subtreeMatches(child, m) {
			return true
		}
	}

	return false
}

// hasActiveFilters reports whether any action filter is currently enabled.
func hasActiveFilters(filters map[ui.Action]bool) bool {
	for _, isActive := range filters {
		if isActive {
			return true
		}
	}

	return false
}
