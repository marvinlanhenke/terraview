// Package app wires the top-level Bubble Tea model for Terraview.
package app

import (
	"log/slog"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"github.com/marvinlanhenke/terraview/internal/planview"
	"github.com/marvinlanhenke/terraview/internal/ui"
	"github.com/marvinlanhenke/terraview/internal/ui/details"
	"github.com/marvinlanhenke/terraview/internal/ui/filter"
	"github.com/marvinlanhenke/terraview/internal/ui/search"
	"github.com/marvinlanhenke/terraview/internal/ui/status"
	"github.com/marvinlanhenke/terraview/internal/ui/tree"
)

// defaultMargin reserves horizontal space for pane gaps and outer padding.
const defaultMargin = 4

// focus identifies the active app pane.
type focus int

const (
	focusTree focus = iota
	focusSearch
	focusDetails
	focusFilter
)

// components groups the child models owned by the app model.
type components struct {
	search  search.Search
	filter  filter.Modal
	status  status.Status
	tree    tree.Tree
	details details.Details
}

// size stores the latest terminal dimensions.
type size struct {
	width  int
	height int
}

// treeControls stores the query and filters applied to the tree.
type treeControls struct {
	query   string
	filters map[ui.Action]bool
}

// Model is the top-level Bubble Tea model for the Terraview app.
type Model struct {
	theme      ui.Theme
	size       size
	focus      focus
	controls   treeControls
	components components
	help       help.Model
	logger     *slog.Logger
}

// New creates an initialized app model from a planview root node.
func New(root *planview.Node, logger *slog.Logger) Model {
	t := ui.DefaultTheme()

	c := components{
		search:  search.New(t),
		filter:  filter.New(t),
		status:  status.New(t),
		tree:    tree.New(t),
		details: details.New(t),
	}

	controls := treeControls{
		filters: make(map[ui.Action]bool),
	}

	m := Model{
		theme:      t,
		size:       size{},
		focus:      focusTree,
		controls:   controls,
		components: c,
		help:       help.New(),
		logger:     logger,
	}

	children := make([]*planview.Node, 0)
	if root != nil {
		children = root.Children
	}

	m.applyLayout(0, 0)

	stats := buildStats(root)
	m.components.status.SetStats(stats)
	logger.Debug("plan stats built", "create", stats.Create, "update", stats.Update, "delete", stats.Delete, "replace", stats.Replace, "no_op", stats.NoOp, "errors", stats.Errors)

	filterOpts := buildFilterOptions(children)
	m.components.filter.SetOptions(filterOpts)
	logger.Debug("filter options built", "count", len(filterOpts))

	treeRoot := buildTreeNode(root)
	m.components.tree.SetRoot(treeRoot)
	logger.Debug("tree root built", "has_root", treeRoot != nil)

	m.refreshTreeFromControls()

	logger.Debug("app model initialized")
	return m
}

// applyLayout stores the terminal dimensions and sizes child components.
func (m *Model) applyLayout(width, height int) {
	m.size.width = width
	m.size.height = height

	contentWidth := max(0, width-defaultMargin)
	m.help.SetWidth(contentWidth)
	m.components.search.SetWidth(contentWidth)
	m.components.status.SetWidth(contentWidth)

	treeWidth, treeHeight := treePaneSize(width, height)
	m.components.tree.SetSize(treeWidth, treeHeight)

	detailsPaneWidth := detailsWidth(m.size.width, treeWidth)
	m.components.details.SetSize(detailsPaneWidth, treeHeight)
}

// refreshTreeFromControls reapplies the current query and filters to the tree.
func (m *Model) refreshTreeFromControls() {
	m.logger.Debug("refreshing tree", "query", m.controls.query, "filter_count", len(m.controls.filters))
	m.components.tree.SetCriteria(m.controls.query, m.controls.filters)
	m.syncTreeOutputs()
}

// syncTreeOutputs updates selection-dependent UI state from the tree.
func (m *Model) syncTreeOutputs() {
	selected := m.components.tree.Selected()
	visible := m.components.tree.VisibleResourceCount()
	m.logger.Debug("tree outputs synced", "visible_resources", visible, "has_selection", selected != nil)
	m.components.details.SetContent(buildDetailsContent(selected))
	m.components.search.SetMatches(visible)
}

// generalFooterBindings returns app-level key bindings shown in the footer.
func (m *Model) generalFooterBindings() []key.Binding {
	bindings := []key.Binding{
		keys.Quit,
		keys.Escape,
		keys.Search,
		keys.Filter,
	}

	if m.focus == focusTree || m.focus == focusDetails {
		bindings = append(bindings, keys.LeftPane, keys.RightPane)
	}

	return bindings
}

// specificFooterBindings returns focus-specific key bindings shown in the footer.
func (m *Model) specificFooterBindings() []key.Binding {
	switch m.focus {
	case focusTree:
		return tree.KeyMap().ShortHelp()
	case focusDetails:
		return details.KeyMap().ShortHelp()
	case focusFilter:
		return filter.KeyMap().ShortHelp()
	default:
		return nil
	}
}
