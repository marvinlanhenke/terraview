package details

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

// Update applies a Bubble Tea message to the details pane.
func (d *Details) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, keys.toggle):
			if d.content.Payload != nil && !d.content.IsError {
				d.showPlan = !d.showPlan
				d.setHeader()
				d.syncViewport()
			}
		}
	}

	d.viewport, cmd = d.viewport.Update(msg)

	return cmd
}
