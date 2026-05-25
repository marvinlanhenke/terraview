package details

import (
	"reflect"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/marvinlanhenke/terraview/internal/ui"
)

func TestFlattenChangesSortsAndExpandsNestedMaps(t *testing.T) {
	rows := flattenChanges(ui.ChangeSet{
		Before: map[string]any{
			"nested": map[string]any{
				"a": true,
				"b": "old",
			},
			"removed": "old",
			"z":       1,
		},
		After: map[string]any{
			"added": "new",
			"nested": map[string]any{
				"b": "new",
				"c": 3,
			},
			"z": 2,
		},
	})

	wantPaths := []string{"added", "nested.a", "nested.b", "nested.c", "nested", "removed", "z"}
	if got := changePaths(rows); !reflect.DeepEqual(got, wantPaths) {
		t.Fatalf("expected paths %v, got %v", wantPaths, got)
	}

	requireChange(t, rows, "added", nil, "new")
	requireChange(t, rows, "nested.a", true, nil)
	requireChange(t, rows, "nested.b", "old", "new")
	requireChange(t, rows, "nested.c", nil, 3)
	requireChange(t, rows, "removed", "old", nil)
	requireChange(t, rows, "z", 1, 2)
}

func TestFlattenChangesRendersEmptyMapOnce(t *testing.T) {
	rows := flattenChanges(ui.ChangeSet{
		After: map[string]any{
			"settings": map[string]any{},
		},
	})

	if got := len(rows); got != 1 {
		t.Fatalf("expected one row, got %d: %#v", got, rows)
	}

	requireChange(t, rows, "settings", nil, map[string]any{})
}

func TestGetJsonStrFormatsNilAndPrefixesJSON(t *testing.T) {
	if got, want := getJsonStr(nil, "  "), "  null"; got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}

	got := getJsonStr(map[string]any{"name": "web"}, "  ")
	want := "  {\n   \"name\": \"web\"\n  }"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestViewRendersEmptyStates(t *testing.T) {
	tests := []struct {
		name    string
		content Content
		want    string
	}{
		{
			name: "none",
			want: "Select a resource to inspect changes.",
		},
		{
			name: "group",
			content: Content{
				Key:  "group",
				Kind: KindGroup,
			},
			want: "Select a resource to inspect changes.",
		},
		{
			name: "resource",
			content: Content{
				Key:  "resource",
				Kind: KindResource,
			},
			want: "No changed attributes for this resource.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			details := newDetails(80, 8)
			if tt.content.Key != "" || tt.content.Kind != KindNone {
				details.SetContent(tt.content)
			}

			view := strippedView(details.View())

			requireContains(t, view, "Details")
			requireContains(t, view, "Diff")
			requireContains(t, view, tt.want)
		})
	}
}

func TestViewRendersDiffContent(t *testing.T) {
	details := newDetails(100, 20)
	details.SetContent(Content{
		Key:  "resource",
		Kind: KindResource,
		Changes: ui.ChangeSet{
			Before: map[string]any{"name": "old"},
			After:  map[string]any{"name": "new"},
		},
	})

	view := strippedView(details.View())

	requireContains(t, view, "Details")
	requireContains(t, view, "Diff")
	requireContains(t, view, "attribute: name")
	requireContains(t, view, "before:")
	requireContains(t, view, "\"old\"")
	requireContains(t, view, "after:")
	requireContains(t, view, "\"new\"")
}

func TestUpdateTogglesBetweenDiffAndPlan(t *testing.T) {
	details := newDetails(100, 20)
	details.SetContent(Content{
		Key:     "resource",
		Kind:    KindResource,
		Payload: map[string]any{"address": "aws_instance.web"},
		Changes: ui.ChangeSet{
			Before: map[string]any{"name": "old"},
			After:  map[string]any{"name": "new"},
		},
	})

	if details.showPlan {
		t.Fatal("expected details to start in diff mode")
	}

	details.Update(keyText("p"))
	if !details.showPlan {
		t.Fatal("expected toggle key to show plan")
	}

	view := strippedView(details.View())
	requireContains(t, view, "Plan")
	requireContains(t, view, "aws_instance.web")

	details.Update(keyText("t"))
	if details.showPlan {
		t.Fatal("expected second toggle key to show diff")
	}

	view = strippedView(details.View())
	requireContains(t, view, "Diff")
	requireContains(t, view, "attribute: name")
}

func TestUpdateKeepsErrorContentOnPlan(t *testing.T) {
	details := newDetails(100, 20)
	details.SetContent(Content{
		Key:     "diagnostic",
		Kind:    KindResource,
		Payload: map[string]any{"severity": "error"},
		IsError: true,
	})

	if !details.showPlan {
		t.Fatal("expected error content to start in plan mode")
	}

	details.Update(keyText("p"))
	if !details.showPlan {
		t.Fatal("expected error content to ignore plan toggle")
	}

	view := strippedView(details.View())
	requireContains(t, view, "Plan")
	requireContains(t, view, "error")
}

func TestFocusAndBlurUpdateState(t *testing.T) {
	details := newDetails(80, 8)

	details.Focus()
	if !details.focus {
		t.Fatal("expected details to be focused")
	}

	details.Blur()
	if details.focus {
		t.Fatal("expected details to be blurred")
	}
}

func TestKeyMapProvidesHelpBindings(t *testing.T) {
	keyMap := KeyMap()

	if got := len(keyMap.ShortHelp()); got != 1 {
		t.Fatalf("expected 1 short help binding, got %d", got)
	}

	if got := len(keyMap.FullHelp()); got != 1 {
		t.Fatalf("expected 1 full help group, got %d", got)
	}
}

func newDetails(width, height int) Details {
	details := New(ui.DefaultTheme())
	details.SetSize(width, height)

	return details
}

func changePaths(rows []change) []string {
	paths := make([]string, len(rows))
	for i, row := range rows {
		paths[i] = row.path
	}

	return paths
}

func requireChange(t *testing.T, rows []change, path string, before, after any) {
	t.Helper()

	for _, row := range rows {
		if row.path != path {
			continue
		}

		if !reflect.DeepEqual(row.before, before) || !reflect.DeepEqual(row.after, after) {
			t.Fatalf("expected change %q before=%#v after=%#v, got before=%#v after=%#v", path, before, after, row.before, row.after)
		}

		return
	}

	t.Fatalf("expected change %q in %#v", path, rows)
}

func keyText(text string) tea.KeyPressMsg {
	runes := []rune(text)
	return tea.KeyPressMsg(tea.Key{Text: text, Code: runes[0]})
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
