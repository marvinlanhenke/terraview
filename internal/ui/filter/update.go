package filter

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

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
				intent = Intent{
					Action:    selected.Action,
					HasToggle: true,
				}
			}

		// Reset Filters
		case key.Matches(msg, keys.reset):
			intent = Intent{
				Reset: true,
			}
		}
	}

	f.clampCursor()

	return intent, nil
}
