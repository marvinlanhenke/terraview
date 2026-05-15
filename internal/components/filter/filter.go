package filter

import (
	"github.com/marvinlanhenke/terraview/internal/components/tree"
	"github.com/marvinlanhenke/terraview/internal/theme"
)

type Filter struct {
	inner map[tree.Action]bool
}

func New(t theme.Theme) Filter {
	inner := make(map[tree.Action]bool)

	return Filter{
		inner: inner,
	}
}

func (f *Filter) ToggleFilters(actions []tree.Action) {
	for _, action := range actions {
		f.ToggleSingleFilter(action)
	}
}

func (f *Filter) ToggleSingleFilter(action tree.Action) {
	before, exists := f.inner[action]

	if !exists {
		f.inner[action] = true
		return
	}

	f.inner[action] = !before
}

func (f *Filter) ResetFilters() {
	f.inner = nil
	f.inner = make(map[tree.Action]bool)
}

func (f Filter) GetFilter() map[tree.Action]bool {
	return f.inner
}
