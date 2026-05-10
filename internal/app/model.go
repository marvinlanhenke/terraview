package app

import (
	"github.com/marvinlanhenke/terraview/internal/components/search"
	"github.com/marvinlanhenke/terraview/internal/components/summary"
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
}

func New() Model {
	return Model{
		focus:   FocusTree,
		search:  search.New(),
		summary: summary.New(),
	}
}
