package summary

type Stats struct {
	Add     int
	Change  int
	Destroy int
	Replace int
	NoOp    int
	Errors  int
}

type Filters struct {
	Add     bool
	Change  bool
	Destroy bool
	Replace bool
	NoOp    bool
	Errors  bool
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

func (s Summary) View() string {
	return "todo summary string styling"
}
