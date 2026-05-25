package search

import (
	tea "charm.land/bubbletea/v2"
)

// Update forwards Bubble Tea messages to the underlying text input.
func (s *Search) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	s.input, cmd = s.input.Update(msg)

	return cmd
}
