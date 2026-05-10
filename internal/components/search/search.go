package search

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/theme"
)

const (
	placeholder = "search resources..."
)

type Search struct {
	input   textinput.Model
	focused bool
}

func New() Search {
	input := textinput.New()

	styles := input.Styles()

	styles.Focused.Placeholder = theme.SearchBar
	styles.Focused.Text = theme.SearchInputFocused
	styles.Focused.Prompt = theme.SearchInputFocused

	styles.Blurred.Placeholder = theme.SearchBar.Faint(true)
	styles.Blurred.Text = theme.SearchInput
	styles.Blurred.Prompt = theme.SearchInput

	input.SetStyles(styles)

	input.Placeholder = placeholder
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

func (s *Search) View(width int, matches int) string {
	if width <= 0 {
		return ""
	}

	label := theme.SearchNugget.Render("[S]")
	status := theme.SearchStatus.Render(fmt.Sprintf("%d matches", matches))

	inputStyle := theme.SearchInput

	if s.Focused() {
		inputStyle = theme.SearchInputFocused
	}

	availableWidth := max(0, width-lipgloss.Width(label)-lipgloss.Width(status))

	input := inputStyle.Width(availableWidth).Render(s.input.View())

	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		label,
		input,
		status,
	)

	return theme.SearchBar.Width(width).Render(row)
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
	s.input.Placeholder = placeholder
}

func (s *Search) SetWidth(width int) {
	s.input.SetWidth(max(0, width))
}
