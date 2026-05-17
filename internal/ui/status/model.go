package status

import (
	"fmt"

	"github.com/marvinlanhenke/terraview/internal/ui/theme"
)

type actions struct {
	add     int
	change  int
	destroy int
	replace int
	noOp    int
	errors  int
}

func (s actions) String() string {
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

type Status struct {
	actions actions

	width  int
	styles styles
}

func New(t theme.Theme) Status {
	return Status{
		styles: newStyles(t),
	}
}

func (s *Status) SetStats(st actions) {
	s.actions = st
}

func (s *Status) SetWidth(width int) {
	s.width = max(0, width)
}
