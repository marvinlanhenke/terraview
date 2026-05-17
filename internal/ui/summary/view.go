package summary

import "fmt"

func (s *Summary) View() string {
	plan := fmt.Sprintf("Plan: %s", s.stats)

	return s.styles.
		borderBar.
		Width(s.width).
		Render(plan)
}
