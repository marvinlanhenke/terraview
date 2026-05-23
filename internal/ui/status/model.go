package status

import (
	"fmt"

	"github.com/marvinlanhenke/terraview/internal/ui/theme"
)

type Stats struct {
	Create  int
	Update  int
	Delete  int
	Replace int
	NoOp    int
	Errors  int
}

func (s Stats) String() string {
	return fmt.Sprintf(
		"[+%d] [~%d] [-%d] [*%d] [=%d] [!%d]",
		s.Create,
		s.Update,
		s.Delete,
		s.Replace,
		s.NoOp,
		s.Errors,
	)
}

type Status struct {
	stats             *Stats
	activeFilterCount int
	width             int
	styles            styles
}

func New(t theme.Theme) Status {
	return Status{
		styles: newStyles(t),
	}
}

func (s *Status) SetStats(st *Stats) {
	s.stats = st
}

func (s *Status) SetActiveFilterCount(count int) {
	s.activeFilterCount = max(0, count)
}

func (s *Status) SetWidth(width int) {
	s.width = max(0, width)
}

func (s *Status) hasActiveFilter() bool {
	return s.activeFilterCount > 0
}
