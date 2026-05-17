package summary

import (
	"fmt"

	"github.com/marvinlanhenke/terraview/internal/ui/theme"
)

type stats struct {
	add     int
	change  int
	destroy int
	replace int
	noOp    int
	errors  int
}

func (s stats) String() string {
	return fmt.Sprintf(
		"+%d ~%d -%d -/+%d =%d !%d",
		s.add,
		s.change,
		s.destroy,
		s.replace,
		s.noOp,
		s.errors,
	)
}

type Summary struct {
	stats stats

	width  int
	styles styles
}

func New(t theme.Theme) Summary {
	return Summary{
		styles: newStyles(t),
	}
}

func (s *Summary) SetStats(st stats) {
	s.stats = st
}

func (s *Summary) SetWidth(width int) {
	s.width = max(0, width)
}

func (s *Summary) View() string {
	plan := fmt.Sprintf("Plan: %s", s.stats)

	return s.styles.
		borderBar.
		Width(s.width).
		Render(plan)
}
