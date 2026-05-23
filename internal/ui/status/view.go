package status

import (
	"fmt"

	"charm.land/lipgloss/v2"
)

func (s *Status) View() string {
	row := s.styles.base
	bar := s.styles.borderBar

	plan := fmt.Sprintf("⚑ Plan: %s", s.stats)
	filterIcon := "⊙"

	if s.hasActiveFilter() {
		filterIcon = "⚲"
	}

	filter := fmt.Sprintf("%s Filter: %d", filterIcon, s.activeFilterCount)

	filterCol := row.Align(lipgloss.Right).Render(filter)
	innerWidth := max(0, s.width-bar.GetHorizontalFrameSize())
	planWidth := max(0, innerWidth-lipgloss.Width(filterCol))
	planCol := row.Width(planWidth).Render(plan)

	line := lipgloss.JoinHorizontal(lipgloss.Top, planCol, filterCol)

	return bar.
		Width(s.width).
		Render(line)
}
