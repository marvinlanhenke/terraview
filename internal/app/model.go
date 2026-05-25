// Package app wires the top-level Bubble Tea model for Terraview.
package app

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"github.com/marvinlanhenke/terraview/internal/planview"
	"github.com/marvinlanhenke/terraview/internal/ui/details"
	"github.com/marvinlanhenke/terraview/internal/ui/filter"
	"github.com/marvinlanhenke/terraview/internal/ui/search"
	"github.com/marvinlanhenke/terraview/internal/ui/status"
	"github.com/marvinlanhenke/terraview/internal/ui/theme"
	"github.com/marvinlanhenke/terraview/internal/ui/tree"
)

const defaultMargin = 4

// focus identifies the active app pane.
type focus int

const (
	focusTree focus = iota
	focusSearch
	focusDetails
	focusFilter
)

type components struct {
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

// treeControls stores the query and filters applied to the tree.
type treeControls struct {
	query   string
	filters map[tree.Action]bool
}

// filterView converts tree filters into the filter modal action type.
func (t *treeControls) filterView() map[filter.Action]bool {
	f := make(map[filter.Action]bool, len(t.filters))
	for k, v := range t.filters {
		f[filter.Action(k)] = v
	}

	return f
}

// Model is the top-level Bubble Tea model for the Terraview app.
type Model struct {
	theme      theme.Theme
	size       size
	focus      focus
	controls   treeControls
	components components
	help       help.Model
}

// New creates a Model from the planview root node.
func New(root *planview.Node) Model {
	t := theme.Default()

	c := components{
		search:  search.New(t),
		filter:  filter.New(t),
		status:  status.New(t),
		tree:    tree.New(t),
		details: details.New(t),
	}

	controls := treeControls{
		filters: make(map[tree.Action]bool),
	}

	m := Model{
		theme:      t,
		size:       size{},
		focus:      focusTree,
		controls:   controls,
		components: c,
		help:       help.New(),
	}

	children := make([]*planview.Node, 0)
	if root != nil {
		children = root.Children
	}

	m.help.SetWidth(m.size.width - defaultMargin)

	m.components.search.SetWidth(m.size.width - defaultMargin)

	m.components.status.SetWidth(m.size.width - defaultMargin)
	m.components.status.SetStats(buildStats(root))

	m.components.filter.SetOptions(buildFilterOptions(children))

	treeWidth, treeHeight := treePaneSize(0, 0)
	m.components.tree.SetSize(treeWidth, treeHeight)
	m.components.tree.SetRoot(buildTreeNode(root))

	detailsWidth := detailsWidth(m.size.width, treeWidth)
	detailsHeight := treeHeight
	m.components.details.SetSize(detailsWidth, detailsHeight)

	m.refreshTreeFromControls()

	return m
}

// refreshTreeFromControls reapplies the current query and filters to the tree.
func (m *Model) refreshTreeFromControls() {
	m.components.tree.SetCriteria(m.controls.query, m.controls.filters)
	m.syncTreeOutputs()
}

// syncTreeOutputs updates selection-dependent UI state from the tree.
func (m *Model) syncTreeOutputs() {
	m.components.details.SetContent(buildDetailsContent(m.components.tree.Selected()))
	m.components.search.SetMatches(m.components.tree.VisibleResourceCount())
}

func (m *Model) generalFooterBindings() []key.Binding {
	bindings := []key.Binding{
		keys.Quit,
		keys.Escape,
		keys.Search,
		keys.Filter,
	}

	switch m.focus {
	case focusTree:
		bindings = append(bindings, keys.LeftPane)
		bindings = append(bindings, keys.RightPane)
	case focusDetails:
		bindings = append(bindings, keys.LeftPane)
		bindings = append(bindings, keys.RightPane)
	}

	return bindings
}

func (m *Model) specificFooterBindings() []key.Binding {
	bindings := make([]key.Binding, 0)

	switch m.focus {
	case focusTree:
		bindings = append(bindings, tree.KeyMap().ShortHelp()...)
	case focusDetails:
		bindings = append(bindings, details.KeyMap().ShortHelp()...)
	case focusFilter:
		bindings = append(bindings, filter.KeyMap().ShortHelp()...)
	}

	return bindings
}

// buildTreeNode converts a planview node into the tree component model.
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

// buildStats counts resource nodes by action for the status component.
func buildStats(n *planview.Node) *status.Stats {
	if n == nil {
		return &status.Stats{}
	}

	stats := &status.Stats{}
	collectStats(n, stats)

	return stats
}

// collectStats recursively accumulates action counts from the plan tree.
func collectStats(n *planview.Node, stats *status.Stats) {
	if n == nil {
		return
	}

	if n.Kind == planview.NodeResource {
		switch n.Action {
		case planview.ActionCreate:
			stats.Create++
		case planview.ActionUpdate:
			stats.Update++
		case planview.ActionDelete:
			stats.Delete++
		case planview.ActionReplace:
			stats.Replace++
		case planview.ActionNoOp:
			stats.NoOp++
		case planview.ActionError:
			stats.Errors++
		}
	}

	for _, child := range n.Children {
		collectStats(child, stats)
	}
}

// buildDetailsContent derives the details pane content from the selected tree node.
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

		content.IsError = n.IsError()
	default:
		content.Kind = details.KindNone
	}

	return content
}

// buildFilterOptions builds filter modal options from the action group nodes.
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
