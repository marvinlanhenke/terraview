package app

import (
	"github.com/marvinlanhenke/terraview/internal/components/search"
	"github.com/marvinlanhenke/terraview/internal/components/summary"
	"github.com/marvinlanhenke/terraview/internal/components/tree"
)

const defaultMargin = 4

type Focus int

const (
	FocusTree Focus = iota
	FocusSearch
	FocusDetails
)

type Model struct {
	width      int
	height     int
	matchCount int

	focus Focus

	search  search.Search
	summary summary.Summary
	tree    tree.Tree
}

func New() Model {
	m := Model{
		focus:   FocusTree,
		search:  search.New(),
		summary: summary.New(),
		tree:    tree.New(),
	}

	m.search.SetWidth(m.width - defaultMargin)
	m.summary.SetWidth(m.width - defaultMargin)

	treeWidth, treeHeight := treePaneSize(0, 0)
	m.tree.SetSize(treeWidth, treeHeight)
	m.tree.SetRoot(getNestedRoot(4, 5))

	return m
}
