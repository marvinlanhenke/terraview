package app

import "github.com/marvinlanhenke/terraview/internal/components/search"

type Focus int

const (
	FocusTree Focus = iota
	FocusSearch
	FocusDetails
)

type Model struct {
	width  int
	height int

	focus Focus

	search  search.Search
	summary string
}

func New() Model {
	return Model{
		focus:   FocusTree,
		search:  search.New(),
		summary: "todo summary",
	}
}
