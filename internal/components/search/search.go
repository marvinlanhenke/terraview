package search

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	placeholder = "search resources..."
)

type Search struct {
	input textinput.Model

	width   int
	matches int
}

func New() Search {
	input := textinput.New()

	styles := input.Styles()

	styles.Focused.Placeholder = searchBar
	styles.Focused.Text = searchInputFocused
	styles.Focused.Prompt = searchInputFocused

	styles.Blurred.Placeholder = searchBar.Faint(true)
	styles.Blurred.Text = searchInput
	styles.Blurred.Prompt = searchInput

	input.SetStyles(styles)

	input.Placeholder = placeholder
	input.CharLimit = 256
	input.Blur()

	return Search{input: input}
}

func (s *Search) SetWidth(width int) {
	s.width = max(0, width)
	s.input.SetWidth(s.width)
}

func (s *Search) SetMatches(matches int) {
	if s.Value() != "" {
		s.matches = max(0, matches)
	}
}

func (s *Search) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	s.input, cmd = s.input.Update(msg)
	return cmd
}

func (s *Search) View() string {
	if s.width <= 0 {
		return ""
	}

	label := searchNugget.Render("[S]")
	status := searchStatus.Render(fmt.Sprintf("%d matches", s.matches))

	inputStyle := searchInput

	if s.Focused() {
		inputStyle = searchInputFocused
	}

	availableWidth := max(0, s.width-lipgloss.Width(label)-lipgloss.Width(status))

	input := inputStyle.Width(availableWidth).Render(s.input.View())

	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		label,
		input,
		status,
	)

	return searchBar.
		Width(s.width).
		MaxWidth(s.width).
		Render(row)
}

func (s *Search) Focus() tea.Cmd {
	s.input.Placeholder = ""
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
	s.matches = 0
	s.input.Placeholder = placeholder
}
