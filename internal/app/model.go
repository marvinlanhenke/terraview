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
	return Model{
		focus:   FocusTree,
		search:  search.New(),
		summary: summary.New(),
		tree:    tree.New(),
	}
}
