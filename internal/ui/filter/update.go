package filter

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

func (f *FilterModal) Update(msg tea.Msg) tea.Cmd {
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
			selected := f.Selected()
			if selected != nil {
				f.ToggleSingleFilter(selected.action)
			}

		// Reset Filters
		case key.Matches(msg, keys.reset):
			f.ResetFilters()
		}
	}

	f.clampCursor()

	return nil
}

func (f *FilterModal) clampCursor() {
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
