package filter

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

// Update applies a Bubble Tea message and returns any requested filter intent.
func (f *Modal) Update(msg tea.Msg) (Intent, tea.Cmd) {
	intent := Intent{}

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
				intent = ToggleIntent(selected.Action)
			}

		// Reset Filters
		case key.Matches(msg, keys.reset):
			intent = ResetIntent()
		}
	}

	f.clampCursor()

	return intent, nil
}
