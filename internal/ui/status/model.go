// Package status renders the app status bar.
package status

import (
	"fmt"

	"github.com/marvinlanhenke/terraview/internal/ui"
)

// Stats contains Terraform plan action counts shown in the status bar.
type Stats struct {
	Create  int
	Update  int
	Delete  int
	Replace int
	NoOp    int
	Errors  int
}

// String formats Stats using the compact action markers used in the UI.
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

// Status renders plan and filter state for the app status bar.
type Status struct {
	stats             Stats
	activeFilterCount int
	width             int
	styles            styles
}

// New returns an initialized Status using t for styling.
func New(t ui.Theme) Status {
	return Status{
		styles: newStyles(t),
	}
}

// SetStats replaces the plan action counts rendered by Status.
func (s *Status) SetStats(st Stats) {
	s.stats = st
}

// SetActiveFilterCount records how many filters are currently enabled.
func (s *Status) SetActiveFilterCount(count int) {
	s.activeFilterCount = max(0, count)
}

// SetWidth sets the rendered status bar width.
func (s *Status) SetWidth(width int) {
	s.width = max(0, width)
}

// hasActiveFilter reports whether any filter is enabled.
func (s *Status) hasActiveFilter() bool {
	return s.activeFilterCount > 0
}
