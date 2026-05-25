package tree

import "github.com/marvinlanhenke/terraview/internal/ui/action"

type criteria struct {
	matcher matcher
	filters map[action.Action]bool
}

type row struct {
	node       *Node
	depth      int
	parent     int
	expandable bool
	expanded   bool
}

func (r row) open(searchActive bool) bool {
	return r.expanded || searchActive
}

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

func includeRootChild(n *Node, filters map[action.Action]bool, filtering bool) bool {
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

func appendRows(rows []row, n *Node, parentIndex, depth int, expanded map[string]bool, m matcher) []row {
	if n == nil {
		return rows
	}

	if m.Active() && !subtreeMatches(n, m) {
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

	if r.open(m.Active()) {
		for _, child := range n.Children {
			rows = appendRows(rows, child, rowIndex, depth+1, expanded, m)
		}
	}

	return rows
}

func subtreeMatches(n *Node, m matcher) bool {
	if n == nil {
		return false
	}

	if m.MatchNode(n) {
		return true
	}

	for _, child := range n.Children {
		if subtreeMatches(child, m) {
			return true
		}
	}

	return false
}

func hasActiveFilters(filters map[action.Action]bool) bool {
	for _, isActive := range filters {
		if isActive {
			return true
		}
	}

	return false
}
