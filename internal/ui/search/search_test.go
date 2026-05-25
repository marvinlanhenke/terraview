package search

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/marvinlanhenke/terraview/internal/ui"
)

func TestViewReturnsEmptyWithoutPositiveWidth(t *testing.T) {
	searchBar := New(ui.DefaultTheme())

	if got := searchBar.View(); got != "" {
		t.Fatalf("expected empty view before width is set, got %q", got)
	}

	searchBar.SetWidth(-1)
	if got := searchBar.View(); got != "" {
		t.Fatalf("expected empty view for negative width, got %q", got)
	}
}

func TestViewRendersSearchBarContent(t *testing.T) {
	searchBar := newSearchBar(80)
	searchBar.SetMatches(3)

	got := searchBar.View()

	requireContains(t, got, "[S]")
	requireContains(t, got, "search resources...")
	requireContains(t, got, "3 matches")
	requireContains(t, got, "◎─TERRAVIEW─◉")
	requireViewWidth(t, got, 80)
}

func TestUpdateChangesValueWhenFocused(t *testing.T) {
	searchBar := newSearchBar(80)
	cmd := searchBar.Focus()
	if cmd == nil {
		t.Fatal("expected focus command")
	}

	searchBar.Update(keyText("w"))
	searchBar.Update(keyText("e"))
	searchBar.Update(keyText("b"))

	if got := searchBar.Value(); got != "web" {
		t.Fatalf("expected query %q, got %q", "web", got)
	}

	if !searchBar.Focused() {
		t.Fatal("expected search input to stay focused")
	}
}

func TestUpdateIgnoresTextWhenBlurred(t *testing.T) {
	searchBar := newSearchBar(80)

	searchBar.Update(keyText("w"))

	if got := searchBar.Value(); got != "" {
		t.Fatalf("expected blurred search input to ignore text, got %q", got)
	}
}

func TestFocusAndBlurUpdatePlaceholder(t *testing.T) {
	searchBar := newSearchBar(80)
	if got := searchBar.input.Placeholder; got != placeholder {
		t.Fatalf("expected initial placeholder %q, got %q", placeholder, got)
	}

	searchBar.Focus()
	if got := searchBar.input.Placeholder; got != "" {
		t.Fatalf("expected focused placeholder to be cleared, got %q", got)
	}

	searchBar.Blur()
	if got := searchBar.input.Placeholder; got != placeholder {
		t.Fatalf("expected blurred placeholder %q, got %q", placeholder, got)
	}
}

func TestClearResetsSearchStateAndInputWidth(t *testing.T) {
	searchBar := newSearchBar(80)
	searchBar.SetMatches(1_000)
	searchBar.Focus()
	searchBar.Update(keyText("w"))
	searchBar.Update(keyText("e"))
	searchBar.Update(keyText("b"))

	searchBar.Clear()

	if got := searchBar.Value(); got != "" {
		t.Fatalf("expected cleared query, got %q", got)
	}

	if got := searchBar.input.Placeholder; got != placeholder {
		t.Fatalf("expected placeholder %q, got %q", placeholder, got)
	}

	requireContains(t, searchBar.View(), "0 matches")
	requireSyncedInputWidth(t, searchBar)
}

func TestInputWidthStaysSyncedWithLayout(t *testing.T) {
	searchBar := newSearchBar(72)
	requireSyncedInputWidth(t, searchBar)

	if got := searchBar.input.Width(); got >= searchBar.width {
		t.Fatalf("expected input viewport width to be smaller than bar width, got %d >= %d", got, searchBar.width)
	}

	searchBar.SetMatches(12_345)
	requireSyncedInputWidth(t, searchBar)

	searchBar.Focus()
	requireSyncedInputWidth(t, searchBar)

	searchBar.Blur()
	requireSyncedInputWidth(t, searchBar)
}

func newSearchBar(width int) Search {
	searchBar := New(ui.DefaultTheme())
	searchBar.SetWidth(width)

	return searchBar
}

func keyText(text string) tea.KeyPressMsg {
	runes := []rune(text)
	return tea.KeyPressMsg(tea.Key{Text: text, Code: runes[0]})
}

func requireContains(t *testing.T, got, want string) {
	t.Helper()

	if !strings.Contains(ansi.Strip(got), want) {
		t.Fatalf("expected view to contain %q, got %q", want, got)
	}
}

func requireViewWidth(t *testing.T, got string, want int) {
	t.Helper()

	if width := lipgloss.Width(got); width != want {
		t.Fatalf("expected view width %d, got %d in %q", want, width, got)
	}
}

func requireSyncedInputWidth(t *testing.T, searchBar Search) {
	t.Helper()

	if got, want := searchBar.input.Width(), searchBar.layout().inputViewportWidth; got != want {
		t.Fatalf("expected input viewport width %d, got %d", want, got)
	}
}
