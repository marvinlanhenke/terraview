package filter

import (
	"github.com/marvinlanhenke/terraview/internal/planview"
	"github.com/marvinlanhenke/terraview/internal/ui/theme"
)

type option struct {
	action planview.Action
	label  string
	count  string
}

type FilterModal struct {
	filters map[planview.Action]bool
	options []option

	cursor int
	width  int
	height int
	styles styles
}

func New(t theme.Theme) FilterModal {
	s := newStyles(t)
	f := make(map[planview.Action]bool)

	return FilterModal{
		filters: f,
		styles:  s,
	}
}

func (f *FilterModal) SetOptions(nodes []*planview.Node) {
	f.options = f.options[:0]

	seen := make(map[planview.Action]struct{})

	for _, n := range nodes {
		if _, exists := seen[n.Action]; !exists {
			option := option{
				action: n.Action,
				label:  n.Label,
				count:  n.LabelCount,
			}
			f.options = append(f.options, option)
		}
		seen[n.Action] = struct{}{}
	}
}

func (f *FilterModal) Selected() *option {
	if len(f.options) == 0 {
		return nil
	}

	return &f.options[f.cursor]
}

func (f *FilterModal) ToggleFilters(actions []planview.Action) {
	for _, action := range actions {
		f.ToggleSingleFilter(action)
	}
}

func (f *FilterModal) ToggleSingleFilter(action planview.Action) {
	before, exists := f.filters[action]

	if !exists {
		f.filters[action] = true
		return
	}

	f.filters[action] = !before
}

func (f *FilterModal) ResetFilters() {
	f.filters = nil
	f.filters = make(map[planview.Action]bool)
}

func (f FilterModal) GetFilters() map[planview.Action]bool {
	return f.filters
}
