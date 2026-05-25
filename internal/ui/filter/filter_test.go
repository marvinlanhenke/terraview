package filter

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/marvinlanhenke/terraview/internal/ui"
)

func TestIntentZeroValueHasNoAction(t *testing.T) {
	var intent Intent

	if intent.Reset() {
		t.Fatal("expected zero-value intent not to reset filters")
	}

	if action, ok := intent.Toggle(); ok {
		t.Fatalf("expected zero-value intent not to toggle filters, got %q", action)
	}
}

func TestToggleIntent(t *testing.T) {
	intent := ToggleIntent(ui.ActionCreate)

	action, ok := intent.Toggle()
	if !ok {
		t.Fatal("expected toggle intent")
	}

	if action != ui.ActionCreate {
		t.Fatalf("expected action %q, got %q", ui.ActionCreate, action)
	}

	if intent.Reset() {
		t.Fatal("expected toggle intent not to reset filters")
	}
}

func TestResetIntent(t *testing.T) {
	intent := ResetIntent()

	if !intent.Reset() {
		t.Fatal("expected reset intent")
	}

	if action, ok := intent.Toggle(); ok {
		t.Fatalf("expected reset intent not to toggle filters, got %q", action)
	}
}

func TestSetOptionsClonesInputAndClampsCursor(t *testing.T) {
	modal := New(ui.DefaultTheme())
	options := testOptions()

	modal.SetOptions(options)
	options[0].Label = "Changed"

	if got := modal.options[0].Label; got != "Create" {
		t.Fatalf("expected options to be cloned, got label %q", got)
	}

	modal.cursor = 10
	modal.SetOptions(options[:1])
	if modal.cursor != 0 {
		t.Fatalf("expected cursor to clamp to 0, got %d", modal.cursor)
	}

	modal.cursor = 3
	modal.SetOptions(nil)
	if modal.cursor != 0 {
		t.Fatalf("expected cursor to reset for empty options, got %d", modal.cursor)
	}
}

func TestUpdateMovesCursorAndClampsAtBounds(t *testing.T) {
	modal := New(ui.DefaultTheme())
	modal.SetOptions(testOptions())

	modal.Update(keySpecial(tea.KeyDown))
	if modal.cursor != 1 {
		t.Fatalf("expected cursor 1 after down, got %d", modal.cursor)
	}

	modal.Update(keySpecial(tea.KeyDown))
	if modal.cursor != 1 {
		t.Fatalf("expected cursor to clamp at final option, got %d", modal.cursor)
	}

	modal.Update(keySpecial(tea.KeyUp))
	if modal.cursor != 0 {
		t.Fatalf("expected cursor 0 after up, got %d", modal.cursor)
	}

	modal.Update(keySpecial(tea.KeyUp))
	if modal.cursor != 0 {
		t.Fatalf("expected cursor to clamp at first option, got %d", modal.cursor)
	}
}

func TestUpdateReturnsToggleIntentForSelectedOption(t *testing.T) {
	modal := New(ui.DefaultTheme())
	modal.SetOptions(testOptions())
	modal.Update(keySpecial(tea.KeyDown))

	intent, cmd := modal.Update(keySpecial(tea.KeyEnter))
	if cmd != nil {
		t.Fatalf("expected nil command, got %T", cmd)
	}

	action, ok := intent.Toggle()
	if !ok {
		t.Fatal("expected toggle intent")
	}

	if action != ui.ActionDelete {
		t.Fatalf("expected action %q, got %q", ui.ActionDelete, action)
	}
}

func TestUpdateReturnsResetIntent(t *testing.T) {
	modal := New(ui.DefaultTheme())
	modal.SetOptions(testOptions())

	intent, cmd := modal.Update(keyText("r"))
	if cmd != nil {
		t.Fatalf("expected nil command, got %T", cmd)
	}

	if !intent.Reset() {
		t.Fatal("expected reset intent")
	}
}

func TestUpdateToggleWithNoOptionsReturnsNoIntent(t *testing.T) {
	modal := New(ui.DefaultTheme())

	intent, cmd := modal.Update(keySpecial(tea.KeyEnter))
	if cmd != nil {
		t.Fatalf("expected nil command, got %T", cmd)
	}

	if intent.Reset() {
		t.Fatal("expected no reset intent")
	}

	if action, ok := intent.Toggle(); ok {
		t.Fatalf("expected no toggle intent, got %q", action)
	}

	if modal.cursor != 0 {
		t.Fatalf("expected cursor to remain 0, got %d", modal.cursor)
	}
}

func TestViewRendersOptionsAndActiveState(t *testing.T) {
	modal := New(ui.DefaultTheme())
	modal.SetOptions(testOptions())

	view := strippedView(modal.View(map[ui.Action]bool{
		ui.ActionDelete: true,
	}))

	requireContains(t, view, "Filter:")
	requireContains(t, view, "Create")
	requireContains(t, view, "Delete")
	requireContains(t, view, "(1/2)")
	requireContains(t, view, "(1/1)")
	requireContains(t, view, "[x]")
	requireContains(t, view, "[ ]")
}

func TestViewTreatsNilActiveMapAsNoActiveFilters(t *testing.T) {
	modal := New(ui.DefaultTheme())
	modal.SetOptions(testOptions())

	view := strippedView(modal.View(nil))

	if strings.Contains(view, "[x]") {
		t.Fatalf("expected no active filters, got %q", view)
	}
}

func TestKeyMapProvidesHelpBindings(t *testing.T) {
	keyMap := KeyMap()

	if got := len(keyMap.ShortHelp()); got != 4 {
		t.Fatalf("expected 4 short help bindings, got %d", got)
	}

	if got := len(keyMap.FullHelp()); got != 2 {
		t.Fatalf("expected 2 full help groups, got %d", got)
	}
}

func testOptions() []Option {
	return []Option{
		{Action: ui.ActionCreate, Label: "Create", Count: "(1/2)"},
		{Action: ui.ActionDelete, Label: "Delete", Count: "(1/1)"},
	}
}

func keyText(text string) tea.KeyPressMsg {
	runes := []rune(text)
	return tea.KeyPressMsg(tea.Key{Text: text, Code: runes[0]})
}

func keySpecial(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: code})
}

func strippedView(view string) string {
	return ansi.Strip(view)
}

func requireContains(t *testing.T, got, want string) {
	t.Helper()

	if !strings.Contains(got, want) {
		t.Fatalf("expected view to contain %q, got %q", want, got)
	}
}
