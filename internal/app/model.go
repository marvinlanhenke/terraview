package app

import (
	"github.com/marvinlanhenke/terraview/internal/plan"
	"github.com/marvinlanhenke/terraview/internal/ui/details"
	"github.com/marvinlanhenke/terraview/internal/ui/filter"
	"github.com/marvinlanhenke/terraview/internal/ui/search"
	"github.com/marvinlanhenke/terraview/internal/ui/summary"
	"github.com/marvinlanhenke/terraview/internal/ui/theme"
	"github.com/marvinlanhenke/terraview/internal/ui/tree"
)

const defaultMargin = 4

type Size struct {
	width  int
	height int
}

type Focus int

const (
	FocusTree Focus = iota
	FocusSearch
	FocusDetails
	FocusFilter
)

type Components struct {
	search  search.Search
	filter  filter.FilterModal
	summary summary.Summary
	tree    tree.Tree
	details details.Details
}

type Model struct {
	theme      theme.Theme
	size       Size
	focus      Focus
	components Components
}

func New(root *plan.Node) Model {
	t := theme.Default()

	c := Components{
		search:  search.New(t),
		filter:  filter.New(t),
		summary: summary.New(t),
		tree:    tree.New(t),
		details: details.New(t),
	}

	m := Model{
		theme:      t,
		size:       Size{},
		focus:      FocusTree,
		components: c,
	}

	m.components.search.SetWidth(m.size.width - defaultMargin)

	m.components.summary.SetWidth(m.size.width - defaultMargin)

	m.components.filter.SetOptions(root.Children)

	treeWidth, treeHeight := treePaneSize(0, 0)
	m.components.tree.SetSize(treeWidth, treeHeight)
	m.components.tree.SetRoot(root)
	m.components.tree.ApplyFilters(c.filter.GetFilters())

	detailsWidth := detailsWidth(m.size.width, treeWidth)
	detailsHeight := treeHeight
	m.components.details.SetSize(detailsWidth, detailsHeight)

	return m
}
