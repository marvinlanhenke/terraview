package tree

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/marvinlanhenke/terraview/internal/ui"
)

func TestNewMatcherMatchesNodeFieldsAndInvalidRegexFallsBackToLiteral(t *testing.T) {
	node := &Node{
		Id:      "aws_instance.web",
		Label:   "Web Instance",
		Kind:    NodeResource,
		Action:  ui.ActionCreate,
		Payload: map[string]any{"owner": "platform"},
	}
	prepareSearchPayloads(node)

	tests := []struct {
		name  string
		query string
		want  bool
	}{
		{name: "empty query", query: "", want: true},
		{name: "literal label match", query: "web instance", want: true},
		{name: "literal action match", query: "CREATE", want: true},
		{name: "payload match", query: "platform", want: true},
		{name: "regex id match", query: "/^aws_.*\\.web$/", want: true},
		{name: "no match", query: "database", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newMatcher(tt.query).matchNode(node); got != tt.want {
				t.Fatalf("expected query %q to match=%t, got %t", tt.query, tt.want, got)
			}
		})
	}

	m := newMatcher("/[/")
	if m.re != nil {
		t.Fatal("expected invalid regex query to fall back to literal matching")
	}

	if !m.matchString("/[/") {
		t.Fatal("expected invalid regex query to match as a literal string")
	}
}

func TestSearchablePayloadHandlesNilStringsJSONAndMarshalErrors(t *testing.T) {
	if got := searchablePayload(nil); got != "null" {
		t.Fatalf("expected nil payload to render as null, got %q", got)
	}

	if got := searchablePayload("plain-text"); got != "plain-text" {
		t.Fatalf("expected string payload to stay unchanged, got %q", got)
	}

	jsonPayload := searchablePayload(map[string]any{"name": "web"})
	if !strings.HasPrefix(jsonPayload, "\n{") {
		t.Fatalf("expected JSON payload to start on a new line, got %q", jsonPayload)
	}

	if !strings.Contains(jsonPayload, `"name":"web"`) {
		t.Fatalf("expected JSON payload to contain marshalled data, got %q", jsonPayload)
	}

	if got := searchablePayload(badPayload{}); got != "bad-payload" {
		t.Fatalf("expected marshal failure to fall back to fmt string, got %q", got)
	}
}

func TestBuildRowsAppliesExpansionSearchAndFilters(t *testing.T) {
	root := testRoot()
	prepareSearchPayloads(root)

	tests := []struct {
		name     string
		expanded map[string]bool
		criteria criteria
		want     []string
	}{
		{
			name: "collapsed tree shows only top-level groups with children",
			want: []string{"create-group", "delete-group"},
		},
		{
			name:     "expanded group shows immediate children",
			expanded: map[string]bool{"create-group": true},
			want:     []string{"create-group", "aws_instance.web", "module.app", "delete-group"},
		},
		{
			name:     "search keeps matching ancestors and descendants visible",
			criteria: criteria{matcher: newMatcher("critical")},
			want:     []string{"create-group", "module.app", "aws_lb.app"},
		},
		{
			name:     "filters only include allowed root groups",
			criteria: criteria{filters: map[ui.Action]bool{ui.ActionDelete: true}},
			want:     []string{"delete-group"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rowIDs(buildRows(root, tt.expanded, tt.criteria)); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("expected row ids %v, got %v", tt.want, got)
			}
		})
	}
}

func TestSetRootRebasesExpandedStatePreparesPayloadsAndClampsCursor(t *testing.T) {
	tree := newTree(72, 8)
	tree.expanded["create-group"] = true
	tree.expanded["module.app"] = true
	tree.expanded["stale"] = true
	tree.cursor = 99

	root := testRoot()
	tree.SetRoot(root)

	wantExpanded := map[string]bool{
		"create-group": true,
		"module.app":   true,
	}
	if !reflect.DeepEqual(tree.expanded, wantExpanded) {
		t.Fatalf("expected expanded ids %v, got %v", wantExpanded, tree.expanded)
	}

	wantRows := []string{"create-group", "aws_instance.web", "module.app", "aws_lb.app", "delete-group"}
	if got := rowIDs(tree.rows); !reflect.DeepEqual(got, wantRows) {
		t.Fatalf("expected rows %v, got %v", wantRows, got)
	}

	if tree.cursor != len(tree.rows)-1 {
		t.Fatalf("expected cursor to clamp to final row %d, got %d", len(tree.rows)-1, tree.cursor)
	}

	if !strings.Contains(root.Children[0].Children[0].searchPayload, "team-web") {
		t.Fatalf("expected search payload to be prepared, got %q", root.Children[0].Children[0].searchPayload)
	}

	if got := tree.Selected(); got == nil || got.Id != "delete-group" {
		t.Fatalf("expected selected node delete-group, got %#v", got)
	}
}

func TestSetCriteriaClonesFiltersAndCountsMatchingResources(t *testing.T) {
	tree := newTree(72, 8)
	tree.SetRoot(testRoot())

	filters := map[ui.Action]bool{ui.ActionDelete: true}
	tree.SetCriteria("", filters)
	filters[ui.ActionCreate] = true

	if tree.filters[ui.ActionCreate] {
		t.Fatalf("expected filters to be cloned, got %#v", tree.filters)
	}

	if got := rowIDs(tree.rows); !reflect.DeepEqual(got, []string{"delete-group"}) {
		t.Fatalf("expected only delete group to remain, got %v", got)
	}

	if got := tree.VisibleResourceCount(); got != 0 {
		t.Fatalf("expected collapsed filtered tree to show 0 resources, got %d", got)
	}

	tree.SetCriteria("critical", nil)

	if len(tree.filters) != 0 {
		t.Fatalf("expected nil filters to become an empty map, got %#v", tree.filters)
	}

	want := []string{"create-group", "module.app", "aws_lb.app"}
	if got := rowIDs(tree.rows); !reflect.DeepEqual(got, want) {
		t.Fatalf("expected search results %v, got %v", want, got)
	}

	if got := tree.VisibleResourceCount(); got != 1 {
		t.Fatalf("expected one visible matching resource, got %d", got)
	}
}

func TestUpdateMovesCursorExpandsAndNavigatesToParent(t *testing.T) {
	tree := newTree(72, 4)
	tree.SetRoot(testRoot())

	tree.Update(keySpecial(tea.KeyDown))
	if tree.cursor != 1 {
		t.Fatalf("expected cursor 1 after down, got %d", tree.cursor)
	}

	tree.Update(keySpecial(tea.KeyDown))
	if tree.cursor != 1 {
		t.Fatalf("expected cursor to clamp at final row, got %d", tree.cursor)
	}

	tree.Update(keySpecial(tea.KeyUp))
	if tree.cursor != 0 {
		t.Fatalf("expected cursor 0 after up, got %d", tree.cursor)
	}

	tree.Update(keySpecial(tea.KeyRight))
	if got := rowIDs(tree.rows); !reflect.DeepEqual(got, []string{"create-group", "aws_instance.web", "module.app", "delete-group"}) {
		t.Fatalf("expected expanded create group, got %v", got)
	}

	tree.Update(keySpecial(tea.KeyDown))
	if got := tree.Selected(); got == nil || got.Id != "aws_instance.web" {
		t.Fatalf("expected resource selection, got %#v", got)
	}

	tree.Update(keySpecial(tea.KeyLeft))
	if tree.cursor != 0 {
		t.Fatalf("expected left on child to move to parent row, got %d", tree.cursor)
	}

	tree.Update(keySpecial(tea.KeyLeft))
	if got := rowIDs(tree.rows); !reflect.DeepEqual(got, []string{"create-group", "delete-group"}) {
		t.Fatalf("expected left on expanded group to collapse it, got %v", got)
	}
}

func TestExpandAllRecursivelyExpandsAndCollapsesNestedGroups(t *testing.T) {
	tree := newTree(72, 8)
	tree.SetRoot(testRoot())

	tree.expandAll(true)

	wantExpandedRows := []string{
		"create-group",
		"aws_instance.web",
		"module.app",
		"aws_lb.app",
		"delete-group",
		"aws_s3_bucket.old",
	}
	if got := rowIDs(tree.rows); !reflect.DeepEqual(got, wantExpandedRows) {
		t.Fatalf("expected fully expanded rows %v, got %v", wantExpandedRows, got)
	}

	wantExpanded := map[string]bool{
		"create-group": true,
		"module.app":   true,
		"delete-group": true,
	}
	if !reflect.DeepEqual(tree.expanded, wantExpanded) {
		t.Fatalf("expected expanded map %v, got %v", wantExpanded, tree.expanded)
	}

	tree.expandAll(false)

	if got := rowIDs(tree.rows); !reflect.DeepEqual(got, []string{"create-group", "delete-group"}) {
		t.Fatalf("expected collapse all to hide nested rows, got %v", got)
	}

	if len(tree.expanded) != 0 {
		t.Fatalf("expected collapse all to clear expanded state, got %#v", tree.expanded)
	}
}

func TestViewRendersEmptyAndPopulatedStates(t *testing.T) {
	tree := newTree(72, 8)

	emptyView := ansi.Strip(tree.View())
	requireContains(t, emptyView, "Nothing to show...")

	tree.SetRoot(testRoot())
	tree.expandAll(true)
	tree.syncViewport()

	view := ansi.Strip(tree.View())
	requireContains(t, view, "⌘ Resources")
	requireContains(t, view, "Create")
	requireContains(t, view, "aws_lb.app")
}

func TestKeyMapProvidesHelpBindings(t *testing.T) {
	keyMap := KeyMap()

	if got := len(keyMap.ShortHelp()); got != 6 {
		t.Fatalf("expected 6 short help bindings, got %d", got)
	}

	if got := len(keyMap.FullHelp()); got != 3 {
		t.Fatalf("expected 3 full help groups, got %d", got)
	}
}

type badPayload struct{}

func (badPayload) MarshalJSON() ([]byte, error) {
	return nil, errors.New("boom")
}

func (badPayload) String() string {
	return "bad-payload"
}

func newTree(width, height int) Tree {
	tree := New(ui.DefaultTheme())
	tree.SetSize(width, height)

	return tree
}

func testRoot() *Node {
	return &Node{
		Id:     "root",
		Label:  "Root",
		Kind:   NodeGroup,
		Action: ui.ActionNoOp,
		Children: []*Node{
			{
				Id:         "create-group",
				Label:      "Create",
				LabelCount: "(2/3)",
				Kind:       NodeGroup,
				Action:     ui.ActionCreate,
				Children: []*Node{
					{
						Id:      "aws_instance.web",
						Label:   "aws_instance.web",
						Kind:    NodeResource,
						Action:  ui.ActionCreate,
						Payload: map[string]any{"address": "aws_instance.web", "owner": "team-web"},
					},
					{
						Id:     "module.app",
						Label:  "module.app",
						Kind:   NodeGroup,
						Action: ui.ActionCreate,
						Children: []*Node{
							{
								Id:      "aws_lb.app",
								Label:   "aws_lb.app",
								Kind:    NodeResource,
								Action:  ui.ActionUpdate,
								Payload: map[string]any{"note": "critical", "owner": "platform"},
							},
						},
					},
				},
			},
			{
				Id:         "delete-group",
				Label:      "Delete",
				LabelCount: "(1/3)",
				Kind:       NodeGroup,
				Action:     ui.ActionDelete,
				Children: []*Node{
					{
						Id:      "aws_s3_bucket.old",
						Label:   "aws_s3_bucket.old",
						Kind:    NodeResource,
						Action:  ui.ActionDelete,
						Payload: "legacy-bucket",
					},
				},
			},
			{
				Id:     "empty-group",
				Label:  "Empty",
				Kind:   NodeGroup,
				Action: ui.ActionUpdate,
			},
		},
	}
}

func rowIDs(rows []row) []string {
	ids := make([]string, len(rows))
	for i, row := range rows {
		ids[i] = row.node.Id
	}

	return ids
}

func keySpecial(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: code})
}

func requireContains(t *testing.T, got, want string) {
	t.Helper()

	if !strings.Contains(got, want) {
		t.Fatalf("expected %q to contain %q", got, want)
	}
}
