package tree

import "github.com/marvinlanhenke/terraview/internal/planview"

type criteria struct {
	matcher matcher
	filters map[planview.Action]bool
}

func buildVisible(root *planview.Node, c criteria) []*planview.Node {
	if root == nil {
		return nil
	}

	filtering := hasActiveFilters(c.filters)
	visible := make([]*planview.Node, 0, len(root.Children))

	for _, child := range root.Children {
		if !includeRootChild(child, c.filters, filtering) {
			continue
		}

		visible = appendVisible(visible, child, c.matcher)
	}

	return visible
}

func includeRootChild(n *planview.Node, filters map[planview.Action]bool, filtering bool) bool {
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

func appendVisible(dst []*planview.Node, n *planview.Node, m matcher) []*planview.Node {
	if n == nil {
		return dst
	}

	if m.Active() && !m.MatchNode(n) && !hasMatchingDescendant(n, m) {
		return dst
	}

	dst = append(dst, n)

	if n.Expanded || m.Active() {
		for _, child := range n.Children {
			dst = appendVisible(dst, child, m)
		}
	}

	return dst
}

func hasMatchingDescendant(n *planview.Node, m matcher) bool {
	for _, child := range n.Children {
		if child == nil {
			continue
		}

		if m.MatchNode(child) || hasMatchingDescendant(child, m) {
			return true
		}
	}

	return false
}

func hasActiveFilters(filters map[planview.Action]bool) bool {
	for _, isActive := range filters {
		if isActive {
			return true
		}
	}

	return false
}
