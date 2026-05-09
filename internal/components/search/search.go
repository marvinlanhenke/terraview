package search

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type Search struct {
	input   textinput.Model
	focused bool
}

func New() Search {
	input := textinput.New()
	input.Placeholder = "search resources..."
	input.CharLimit = 256
	input.Blur()

	return Search{input: input}
}

func (s *Search) Init() tea.Cmd {
	return textinput.Blink
}

func (s *Search) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	s.input, cmd = s.input.Update(msg)
	return cmd
}

func (s *Search) View() string {
	return s.input.View()
}

func (s *Search) Focus() tea.Cmd {
	s.input.Focus()
	return textinput.Blink
}

func (s *Search) Focused() bool {
	return s.input.Focused()
}

func (s *Search) Blur() {
	s.input.Blur()
}

func (s *Search) Value() string {
	return strings.TrimSpace(s.input.Value())
}

func (s *Search) Clear() {
	s.input.SetValue("")
}

func (s *Search) SetWidth(width int) {
	s.input.SetWidth(max(width, 0))
}
