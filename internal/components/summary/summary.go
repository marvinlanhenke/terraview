package summary

import (
	"fmt"

	"github.com/marvinlanhenke/terraview/internal/theme"
)

type Stats struct {
	Add     int
	Change  int
	Destroy int
	Replace int
	NoOp    int
	Errors  int
}

func (s Stats) String() string {
	return fmt.Sprintf(
		"+%d ~%d -%d -/+%d =%d !%d",
		s.Add,
		s.Change,
		s.Destroy,
		s.Replace,
		s.NoOp,
		s.Errors,
	)
}

type Filters struct {
	Add     bool
	Change  bool
	Destroy bool
	Replace bool
	NoOp    bool
	Errors  bool
}

func (f Filters) String() string {
	check := func(v bool) string {
		if v {
			return "x"
		} else {
			return " "
		}
	}

	return fmt.Sprintf(
		"[%s] add [%s] change [%s] destroy [%s] replace [%s] noop [%s] errors",
		check(f.Add),
		check(f.Change),
		check(f.Destroy),
		check(f.Replace),
		check(f.NoOp),
		check(f.Errors),
	)
}

type Summary struct {
	stats   Stats
	filters Filters

	width  int
	styles styles
}

func New(t theme.Theme) Summary {
	return Summary{
		styles: newStyles(t),
	}
}

func (s *Summary) SetStats(stats Stats) {
	s.stats = stats
}

func (s *Summary) SetFilters(filters Filters) {
	s.filters = filters
}

func (s *Summary) SetWidth(width int) {
	s.width = max(0, width)
}

func (s *Summary) View() string {
	filters := fmt.Sprintf("Filters: %s", s.filters)
	plan := fmt.Sprintf("Plan: %s", s.stats)

	return s.styles.
		borderBar.
		Width(s.width).
		Render(filters + "\n\n" + plan)
}
