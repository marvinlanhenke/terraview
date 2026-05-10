package summary

import "fmt"

type Stats struct {
	Add     int
	Change  int
	Destroy int
	Replace int
	NoOp    int
	Errors  int
}

// TODO: Prettify output; add color coding
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

// TODO: Add filter query param string
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
}

func New() Summary {
	return Summary{}
}

func (s *Summary) SetStats(stats Stats) {
	s.stats = stats
}

func (s *Summary) SetFilters(filters Filters) {
	s.filters = filters
}

func (s Summary) View(width int) string {
	filters := fmt.Sprintf("Filters: %s", s.filters)
	plan := fmt.Sprintf("Plan: %s", s.stats)
	return summaryBar.Width(width).Render(filters + "\n\n" + plan)
}
