package status_test

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/marvinlanhenke/terraview/internal/ui/status"
	"github.com/marvinlanhenke/terraview/internal/ui/theme"
)

func TestStatsString(t *testing.T) {
	stats := status.Stats{
		Create:  1,
		Update:  2,
		Delete:  3,
		Replace: 4,
		NoOp:    5,
		Errors:  6,
	}

	const want = "[+1] [~2] [-3] [*4] [=5] [!6]"
	if got := stats.String(); got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestViewReturnsEmptyWithoutPositiveWidth(t *testing.T) {
	statusBar := status.New(theme.Default())

	if got := statusBar.View(); got != "" {
		t.Fatalf("expected empty view before width is set, got %q", got)
	}

	statusBar.SetWidth(-1)
	if got := statusBar.View(); got != "" {
		t.Fatalf("expected empty view for negative width, got %q", got)
	}
}

func TestViewRendersStatsAndInactiveFilter(t *testing.T) {
	statusBar := newStatusBar(80)
	statusBar.SetStats(status.Stats{
		Create:  2,
		Update:  3,
		Delete:  4,
		Replace: 5,
		NoOp:    6,
		Errors:  7,
	})

	got := statusBar.View()

	requireContains(t, got, "⚑ Plan: [+2] [~3] [-4] [*5] [=6] [!7]")
	requireContains(t, got, "⊙ Filter: 0")
	requireNotContains(t, got, "⚲ Filter")
	requireViewWidth(t, got, 80)
}

func TestViewRendersActiveFilterCount(t *testing.T) {
	statusBar := newStatusBar(80)
	statusBar.SetActiveFilterCount(3)

	got := statusBar.View()

	requireContains(t, got, "⚲ Filter: 3")
	requireNotContains(t, got, "⊙ Filter")
	requireViewWidth(t, got, 80)
}

func TestViewClampsNegativeFilterCount(t *testing.T) {
	statusBar := newStatusBar(80)
	statusBar.SetActiveFilterCount(-3)

	got := statusBar.View()

	requireContains(t, got, "⊙ Filter: 0")
	requireNotContains(t, got, "Filter: -")
	requireViewWidth(t, got, 80)
}

func newStatusBar(width int) status.Status {
	statusBar := status.New(theme.Default())
	statusBar.SetWidth(width)

	return statusBar
}

func requireContains(t *testing.T, got, want string) {
	t.Helper()

	if !strings.Contains(got, want) {
		t.Fatalf("expected view to contain %q, got %q", want, got)
	}
}

func requireNotContains(t *testing.T, got, want string) {
	t.Helper()

	if strings.Contains(got, want) {
		t.Fatalf("expected view not to contain %q, got %q", want, got)
	}
}

func requireViewWidth(t *testing.T, got string, want int) {
	t.Helper()

	if width := lipgloss.Width(got); width != want {
		t.Fatalf("expected view width %d, got %d in %q", want, width, got)
	}
}
