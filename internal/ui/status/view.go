package status

import "fmt"

func (s *Status) View() string {
	plan := fmt.Sprintf("Plan: %s", s.actions)

	return s.styles.
		borderBar.
		Width(s.width).
		Render(plan)
}
