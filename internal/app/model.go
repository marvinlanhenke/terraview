package app

import (
	"github.com/marvinlanhenke/terraview/internal/planview"
	"github.com/marvinlanhenke/terraview/internal/ui/details"
	"github.com/marvinlanhenke/terraview/internal/ui/filter"
	"github.com/marvinlanhenke/terraview/internal/ui/search"
	"github.com/marvinlanhenke/terraview/internal/ui/status"
	"github.com/marvinlanhenke/terraview/internal/ui/theme"
	"github.com/marvinlanhenke/terraview/internal/ui/tree"
)

const defaultMargin = 4

type Focus int

const (
	FocusTree Focus = iota
	FocusSearch
	FocusDetails
	FocusFilter
)

type TreeControls struct {
	query   string
	filters map[tree.Action]bool
}

func (t *TreeControls) filterView() map[filter.Action]bool {
	f := make(map[filter.Action]bool, len(t.filters))
	for k, v := range t.filters {
		f[filter.Action(k)] = v
	}

	return f
}

type Components struct {
	search  search.Search
	filter  filter.Modal
	status  status.Status
	tree    tree.Tree
	details details.Details
}

type size struct {
	width  int
	height int
}

type Model struct {
	theme      theme.Theme
	size       size
	focus      Focus
	controls   TreeControls
	components Components
}

func New(root *planview.Node) Model {
	t := theme.Default()

	c := Components{
		search:  search.New(t),
		filter:  filter.New(t),
		status:  status.New(t),
		tree:    tree.New(t),
		details: details.New(t),
	}

	controls := TreeControls{
		filters: make(map[tree.Action]bool),
	}

	m := Model{
		theme:      t,
		size:       size{},
		focus:      FocusTree,
		controls:   controls,
		components: c,
	}

	m.components.search.SetWidth(m.size.width - defaultMargin)

	m.components.status.SetWidth(m.size.width - defaultMargin)

	m.components.filter.SetOptions(buildFilterOptions(root.Children))

	treeWidth, treeHeight := treePaneSize(0, 0)
	m.components.tree.SetSize(treeWidth, treeHeight)
	m.components.tree.SetRoot(buildTreeNode(root))

	detailsWidth := detailsWidth(m.size.width, treeWidth)
	detailsHeight := treeHeight
	m.components.details.SetSize(detailsWidth, detailsHeight)

	m.refreshTreeFromControls()

	return m
}

func (m *Model) refreshTreeFromControls() {
	m.components.tree.SetCriteria(m.controls.query, m.controls.filters)
	m.syncTreeOutputs()
}

func (m *Model) syncTreeOutputs() {
	m.components.details.SetContent(buildDetailsContent(m.components.tree.Selected()))
	m.components.search.SetMatches(m.components.tree.VisibleResourceCount())
}

func buildTreeNode(n *planview.Node) *tree.Node {
	if n == nil {
		return nil
	}

	out := &tree.Node{
		Id:         n.Id,
		Label:      n.Label,
		LabelCount: n.LabelCount,
		Kind:       tree.NodeKind(n.Kind),
		Action:     tree.Action(n.Action),
		Payload:    n.Payload,
		Changes:    tree.ChangeSet{Before: n.ChangeSetBefore(), After: n.ChangeSetAfter()},
	}

	if len(n.Children) > 0 {
		out.Children = make([]*tree.Node, len(n.Children))
		for i, child := range n.Children {
			out.Children[i] = buildTreeNode(child)
		}
	}

	return out
}

func buildDetailsContent(n *tree.Node) details.Content {
	if n == nil {
		return details.Content{Kind: details.KindNone}
	}

	content := details.Content{
		Key:   n.Id,
		Label: n.Label,
	}

	switch n.Kind {
	case tree.NodeGroup:
		content.Kind = details.KindGroup
	case tree.NodeResource:
		content.Kind = details.KindResource

		content.Changes = details.ChangeSet{
			Before: n.Changes.Before,
			After:  n.Changes.After,
		}

		content.Payload = n.Payload
	default:
		content.Kind = details.KindNone
	}

	return content
}

func buildFilterOptions(nodes []*planview.Node) []filter.Option {
	seen := make(map[filter.Action]struct{})
	options := make([]filter.Option, 0, len(nodes))

	for _, n := range nodes {
		if n == nil {
			continue
		}

		action := filter.Action(n.Action)

		if _, exists := seen[action]; !exists {
			option := filter.Option{
				Action: action,
				Label:  n.Label,
				Count:  n.LabelCount,
			}

			options = append(options, option)
		}

		seen[action] = struct{}{}
	}

	return options
}
