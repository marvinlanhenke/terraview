package filter

import (
	"maps"

	"github.com/marvinlanhenke/terraview/internal/planview"
	"github.com/marvinlanhenke/terraview/internal/ui/theme"
)

type option struct {
	action planview.Action
	label  string
	count  string
}

type Modal struct {
	filters map[planview.Action]bool
	options []option

	cursor int
	styles styles
}

func New(t theme.Theme) Modal {
	s := newStyles(t)
	f := make(map[planview.Action]bool)

	return Modal{
		filters: f,
		styles:  s,
	}
}

func (f *Modal) SetOptions(nodes []*planview.Node) {
	f.options = f.options[:0]

	seen := make(map[planview.Action]struct{})

	for _, n := range nodes {
		if n == nil {
			continue
		}

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

	f.clampCursor()
}

func (f *Modal) Filters() map[planview.Action]bool {
	filters := make(map[planview.Action]bool, len(f.filters))
	maps.Copy(filters, f.filters)

	return filters
}

func (f *Modal) resetFilters() {
	f.filters = nil
	f.filters = make(map[planview.Action]bool)
}

func (f *Modal) toggleSingleFilter(action planview.Action) {
	before, exists := f.filters[action]

	if !exists {
		f.filters[action] = true
		return
	}

	f.filters[action] = !before
}

func (f *Modal) selected() *option {
	if len(f.options) == 0 {
		return nil
	}

	return &f.options[f.cursor]
}

func (f *Modal) clampCursor() {
	if len(f.options) == 0 {
		f.cursor = 0
		return
	}

	if f.cursor < 0 {
		f.cursor = 0
	}

	if f.cursor >= len(f.options) {
		f.cursor = len(f.options) - 1
	}
}
