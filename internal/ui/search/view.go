package search

import (
	"charm.land/lipgloss/v2"
	"fmt"
)

func (s *Search) View() string {
	if s.width <= 0 {
		return ""
	}

	label := s.styles.nugget.Render("[S]")
	status := s.styles.status.Render(fmt.Sprintf("%d matches", s.matches))
	banner := s.styles.banner.Render("◎─TERRAVIEW─◉")

	inputStyle := s.styles.input

	if s.Focused() {
		inputStyle = s.styles.inputAlt
	}

	availableWidth := max(0, s.width-lipgloss.Width(label)-lipgloss.Width(status)-lipgloss.Width(banner)-1)

	input := inputStyle.Width(availableWidth).Render(s.input.View())

	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		label,
		input,
		status,
		" ",
		banner,
	)

	return s.styles.
		background.
		Width(s.width).
		MaxWidth(s.width).
		Render(row)
}
