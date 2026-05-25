package search

import (
	"fmt"

	"charm.land/lipgloss/v2"
)

type searchLayout struct {
	label              string
	status             string
	banner             string
	inputStyle         lipgloss.Style
	inputColumnWidth   int
	inputViewportWidth int
}

func (s *Search) View() string {
	if s.width <= 0 {
		return ""
	}

	layout := s.layout()

	input := layout.inputStyle.
		Width(layout.inputColumnWidth).
		MaxWidth(layout.inputColumnWidth).
		Render(s.input.View())

	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		layout.label,
		input,
		layout.status,
		" ",
		layout.banner,
	)

	return s.styles.background.
		Width(s.width).
		MaxWidth(s.width).
		Render(row)
}

func (s *Search) layout() searchLayout {
	label := s.styles.nugget.Render("[S]")
	status := s.styles.status.Render(fmt.Sprintf("%d matches", s.matches))
	banner := s.styles.banner.Render("◎─TERRAVIEW─◉")

	inputStyle := s.styles.input
	if s.Focused() {
		inputStyle = s.styles.inputAlt
	}

	inputColumnWidth := max(
		0,
		s.width-lipgloss.Width(label)-lipgloss.Width(status)-lipgloss.Width(" ")-lipgloss.Width(banner),
	)

	inputViewportWidth := max(
		1,
		inputColumnWidth-inputStyle.GetHorizontalFrameSize()-lipgloss.Width(s.input.Prompt),
	)

	return searchLayout{
		label:              label,
		status:             status,
		banner:             banner,
		inputStyle:         inputStyle,
		inputColumnWidth:   inputColumnWidth,
		inputViewportWidth: inputViewportWidth,
	}
}
