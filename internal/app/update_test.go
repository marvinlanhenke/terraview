package app

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/marvinlanhenke/terraview/internal/planview"
	"github.com/marvinlanhenke/terraview/internal/ui"
	"github.com/marvinlanhenke/terraview/internal/ui/filter"
)

func TestRouteKeyPressConsumesSearchActivation(t *testing.T) {
	m := New(testPlanRoot())

	cmd, consumed := m.routeKeyPress(keyText("/"))
	if !consumed {
		t.Fatal("expected search activation to consume the slash key")
	}

	if cmd == nil {
		t.Fatal("expected search activation command")
	}

	if m.focus != focusSearch {
		t.Fatalf("expected search focus, got %v", m.focus)
	}

	if !m.components.search.Focused() {
		t.Fatal("expected search input to be focused")
	}

	if got := m.components.search.Value(); got != "" {
		t.Fatalf("expected slash key to be consumed before text input, got query %q", got)
	}
}

func TestRouteKeyPressConsumesSearchEscape(t *testing.T) {
	m := New(testPlanRoot())
	m = updateModel(t, m, keyText("/"))
	m = updateModel(t, m, keyText("w"))

	_, consumed := m.routeKeyPress(keySpecial(tea.KeyEscape))
	if !consumed {
		t.Fatal("expected search escape to be consumed")
	}

	if m.focus != focusTree {
		t.Fatalf("expected tree focus, got %v", m.focus)
	}

	if m.controls.query != "" {
		t.Fatalf("expected search query to be cleared, got %q", m.controls.query)
	}

	if got := m.components.search.Value(); got != "" {
		t.Fatalf("expected search input to be cleared, got %q", got)
	}

	if m.components.search.Focused() {
		t.Fatal("expected search input to be blurred")
	}
}

func TestUpdateSearchInputUpdatesTreeControls(t *testing.T) {
	m := New(testPlanRoot())

	m = updateModel(t, m, keyText("/"))
	m = updateModel(t, m, keyText("w"))
	m = updateModel(t, m, keyText("e"))
	m = updateModel(t, m, keyText("b"))

	if m.focus != focusSearch {
		t.Fatalf("expected search focus, got %v", m.focus)
	}

	if m.controls.query != "web" {
		t.Fatalf("expected query %q, got %q", "web", m.controls.query)
	}
}

func TestUpdateSearchEnterDoesNotExpandTree(t *testing.T) {
	m := New(testPlanRoot())
	if got := m.components.tree.VisibleResourceCount(); got != 0 {
		t.Fatalf("expected collapsed tree to have no visible resources, got %d", got)
	}

	m = updateModel(t, m, keyText("/"))
	m = updateModel(t, m, keySpecial(tea.KeyEnter))

	if m.focus != focusTree {
		t.Fatalf("expected tree focus after applying search, got %v", m.focus)
	}

	if got := m.components.tree.VisibleResourceCount(); got != 0 {
		t.Fatalf("expected search enter to avoid expanding tree, got %d visible resources", got)
	}
}

func TestUpdateRightPaneOnlyRoutesFromTreeFocus(t *testing.T) {
	m := New(testPlanRoot())
	m.focus = focusSearch
	m.components.search.Focus()

	m = updateModel(t, m, keyCtrl('l'))

	if m.focus != focusSearch {
		t.Fatalf("expected right-pane key to stay in search focus, got %v", m.focus)
	}
}

func TestUpdateCanFocusDetailsForSelectedResource(t *testing.T) {
	m := New(testPlanRoot())

	m = updateModel(t, m, keyText("e"))
	m = updateModel(t, m, keySpecial(tea.KeyDown))
	m = updateModel(t, m, keyCtrl('l'))

	if m.focus != focusDetails {
		t.Fatalf("expected details focus, got %v", m.focus)
	}
}

func TestUpdateFilterToggleAndReset(t *testing.T) {
	m := New(testPlanRoot())

	m = updateModel(t, m, keyText("f"))
	if m.focus != focusFilter {
		t.Fatalf("expected filter focus, got %v", m.focus)
	}

	m = updateModel(t, m, keySpecial(tea.KeyEnter))
	if !m.controls.filters[ui.ActionCreate] {
		t.Fatalf("expected create filter to be active, got %#v", m.controls.filters)
	}

	m = updateModel(t, m, keyText("r"))
	if len(m.controls.filters) != 0 {
		t.Fatalf("expected filters to reset, got %#v", m.controls.filters)
	}
}

func TestUpdateWindowSizeAppliesLayout(t *testing.T) {
	m := New(testPlanRoot())
	m = updateModel(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})

	if m.size.width != 120 || m.size.height != 40 {
		t.Fatalf("expected size 120x40, got %dx%d", m.size.width, m.size.height)
	}
}

func TestApplySearchQuery(t *testing.T) {
	m := New(testPlanRoot())

	if m.applySearchQuery() {
		t.Fatal("expected unchanged empty query")
	}

	m = updateModel(t, m, keyText("/"))
	m = updateModel(t, m, keyText("a"))

	if m.applySearchQuery() {
		t.Fatal("expected query to already be applied by updateFocused")
	}
}

func TestApplyFilterIntent(t *testing.T) {
	m := New(testPlanRoot())

	if !m.applyFilterIntent(filterIntent(ui.ActionCreate)) {
		t.Fatal("expected create filter toggle to report a change")
	}

	if !m.controls.filters[ui.ActionCreate] {
		t.Fatal("expected create filter to be active")
	}

	if !m.applyFilterIntent(filterIntent(ui.ActionCreate)) {
		t.Fatal("expected create filter removal to report a change")
	}

	if m.controls.filters[ui.ActionCreate] {
		t.Fatal("expected create filter to be inactive")
	}

	if m.applyFilterIntent(filterResetIntent()) {
		t.Fatal("expected empty filter reset to report no change")
	}

	m.controls.filters[ui.ActionDelete] = true
	if !m.applyFilterIntent(filterResetIntent()) {
		t.Fatal("expected non-empty filter reset to report a change")
	}

	if len(m.controls.filters) != 0 {
		t.Fatalf("expected reset filters, got %#v", m.controls.filters)
	}
}

func updateModel(t *testing.T, m Model, msg tea.Msg) Model {
	t.Helper()

	next, _ := m.Update(msg)
	model, ok := next.(Model)
	if !ok {
		t.Fatalf("expected app.Model, got %T", next)
	}

	return model
}

func keyText(text string) tea.KeyPressMsg {
	runes := []rune(text)
	return tea.KeyPressMsg(tea.Key{Text: text, Code: runes[0]})
}

func keySpecial(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: code})
}

func keyCtrl(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: code, Mod: tea.ModCtrl})
}

func filterIntent(action ui.Action) filter.Intent {
	return filter.Intent{Action: action, HasToggle: true}
}

func filterResetIntent() filter.Intent {
	return filter.Intent{Reset: true}
}

func testPlanRoot() *planview.Node {
	return &planview.Node{
		Id:     "root",
		Label:  "Root",
		Kind:   planview.NodeGroup,
		Action: planview.ActionNoOp,
		Children: []*planview.Node{
			{
				Id:         "create-group",
				Label:      "Create",
				LabelCount: "(1/2)",
				Kind:       planview.NodeGroup,
				Action:     planview.ActionCreate,
				Children: []*planview.Node{
					{
						Id:      "aws_instance.web",
						Label:   "aws_instance.web",
						Kind:    planview.NodeResource,
						Action:  planview.ActionCreate,
						Payload: map[string]any{"address": "aws_instance.web"},
					},
				},
			},
			{
				Id:         "delete-group",
				Label:      "Delete",
				LabelCount: "(1/2)",
				Kind:       planview.NodeGroup,
				Action:     planview.ActionDelete,
				Children: []*planview.Node{
					{
						Id:      "aws_s3_bucket.old",
						Label:   "aws_s3_bucket.old",
						Kind:    planview.NodeResource,
						Action:  planview.ActionDelete,
						Payload: map[string]any{"address": "aws_s3_bucket.old"},
					},
				},
			},
		},
	}
}
