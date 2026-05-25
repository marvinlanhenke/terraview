package app

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestViewUsesAppViewSettings(t *testing.T) {
	m := New(testPlanRoot())
	m.applyLayout(100, 30)

	view := m.View()
	if !view.AltScreen {
		t.Fatal("expected alt screen view")
	}

	if view.MouseMode != tea.MouseModeCellMotion {
		t.Fatalf("expected cell motion mouse mode, got %v", view.MouseMode)
	}

	if !strings.Contains(view.Content, "TERRAVIEW") {
		t.Fatalf("expected app banner in rendered content, got %q", view.Content)
	}
}

func TestRenderFilterOverlayIncludesModal(t *testing.T) {
	m := New(testPlanRoot())
	m.applyLayout(100, 30)
	m.focus = focusFilter

	content := m.View().Content
	if !strings.Contains(content, "[ ]") {
		t.Fatalf("expected filter modal content, got %q", content)
	}
}

func TestLayoutHelpers(t *testing.T) {
	treeWidth, treeHeight := treePaneSize(120, 40)
	if treeWidth != 40 || treeHeight != 27 {
		t.Fatalf("expected tree pane 40x27, got %dx%d", treeWidth, treeHeight)
	}

	treeWidth, treeHeight = treePaneSize(0, 0)
	if treeWidth != 20 || treeHeight != 5 {
		t.Fatalf("expected minimum tree pane 20x5, got %dx%d", treeWidth, treeHeight)
	}

	if got := detailsWidth(120, 40); got != 75 {
		t.Fatalf("expected details width 75, got %d", got)
	}

	if got := detailsWidth(0, 40); got != 20 {
		t.Fatalf("expected minimum details width 20, got %d", got)
	}
}
