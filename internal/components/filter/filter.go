package filter

import (
	"github.com/marvinlanhenke/terraview/internal/components/tree"
	"github.com/marvinlanhenke/terraview/internal/theme"
)

type FilterModal struct {
	inner map[tree.Action]bool
}

func New(t theme.Theme) FilterModal {
	inner := make(map[tree.Action]bool)

	return FilterModal{
		inner: inner,
	}
}

func (f *FilterModal) ToggleFilters(actions []tree.Action) {
	for _, action := range actions {
		f.ToggleSingleFilter(action)
	}
}

func (f *FilterModal) ToggleSingleFilter(action tree.Action) {
	before, exists := f.inner[action]

	if !exists {
		f.inner[action] = true
		return
	}

	f.inner[action] = !before
}

func (f *FilterModal) ResetFilters() {
	f.inner = nil
	f.inner = make(map[tree.Action]bool)
}

func (f FilterModal) GetFilter() map[tree.Action]bool {
	return f.inner
}
