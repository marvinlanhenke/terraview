package app

import (
	"github.com/marvinlanhenke/terraview/internal/components/search"
	"github.com/marvinlanhenke/terraview/internal/components/summary"
	"github.com/marvinlanhenke/terraview/internal/components/tree"
)

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

	treeWidth, treeHeight := treePaneSize(0, 0)
	m.tree.SetSize(treeWidth, treeHeight)
	m.tree.SetRoot(getNestedRoot(5, 5))

	return m
}
