package search

import (
	tea "charm.land/bubbletea/v2"
)

func (s *Search) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	s.input, cmd = s.input.Update(msg)

	return cmd
}
