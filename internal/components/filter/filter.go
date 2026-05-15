package filter

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/components/tree"
	"github.com/marvinlanhenke/terraview/internal/theme"
)

type FilterModal struct {
	filters map[tree.Action]bool

	cursor int
	width  int
	height int
	styles styles
}

func New(t theme.Theme) FilterModal {
	inner := make(map[tree.Action]bool)

	return FilterModal{
		filters: inner,
	}
}

func (f *FilterModal) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (f FilterModal) View() string {
	return lipgloss.NewStyle().
		Width(28).
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Render("Filters")
}

func (f *FilterModal) ToggleFilters(actions []tree.Action) {
	for _, action := range actions {
		f.ToggleSingleFilter(action)
	}
}

func (f *FilterModal) ToggleSingleFilter(action tree.Action) {
	before, exists := f.filters[action]

	if !exists {
		f.filters[action] = true
		return
	}

	f.filters[action] = !before
}

func (f *FilterModal) ResetFilters() {
	f.filters = nil
	f.filters = make(map[tree.Action]bool)
}

func (f FilterModal) GetFilters() map[tree.Action]bool {
	return f.filters
}
