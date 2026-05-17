package filter

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

func (f *Modal) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		// Up
		case key.Matches(msg, keys.up):
			f.cursor--

		// Down
		case key.Matches(msg, keys.down):
			f.cursor++

		// Toggle Filter
		case key.Matches(msg, keys.toggle):
			selected := f.selected()
			if selected != nil {
				f.toggleSingleFilter(selected.action)
			}

		// Reset Filters
		case key.Matches(msg, keys.reset):
			f.resetFilters()
		}
	}

	f.clampCursor()

	return nil
}
