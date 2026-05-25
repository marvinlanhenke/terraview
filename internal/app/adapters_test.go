package app

import (
	"reflect"
	"testing"

	"github.com/marvinlanhenke/terraview/internal/planview"
	"github.com/marvinlanhenke/terraview/internal/ui"
	"github.com/marvinlanhenke/terraview/internal/ui/details"
	"github.com/marvinlanhenke/terraview/internal/ui/filter"
	"github.com/marvinlanhenke/terraview/internal/ui/status"
	"github.com/marvinlanhenke/terraview/internal/ui/tree"
)

func TestBuildTreeNodeAdaptsPlanTree(t *testing.T) {
	payload := map[string]any{"name": "web"}
	root := &planview.Node{
		Id:     "root",
		Label:  "Root",
		Kind:   planview.NodeGroup,
		Action: planview.ActionNoOp,
		Children: []*planview.Node{
			{
				Id:         "create-group",
				Label:      "Create",
				LabelCount: "(1/1)",
				Kind:       planview.NodeGroup,
				Action:     planview.ActionCreate,
				Children: []*planview.Node{
					{
						Id:      "aws_instance.web",
						Label:   "aws_instance.web",
						Kind:    planview.NodeResource,
						Action:  planview.ActionCreate,
						Payload: payload,
					},
				},
			},
		},
	}

	got := buildTreeNode(root)
	if got == nil {
		t.Fatal("expected tree node")
	}

	if got.Id != root.Id || got.Label != root.Label || got.Kind != tree.NodeGroup || got.Action != ui.ActionNoOp {
		t.Fatalf("unexpected root node: %#v", got)
	}

	if len(got.Children) != 1 {
		t.Fatalf("expected one action group, got %d", len(got.Children))
	}

	group := got.Children[0]
	if group.Id != "create-group" || group.LabelCount != "(1/1)" || group.Kind != tree.NodeGroup || group.Action != ui.ActionCreate {
		t.Fatalf("unexpected action group: %#v", group)
	}

	if len(group.Children) != 1 {
		t.Fatalf("expected one resource, got %d", len(group.Children))
	}

	resource := group.Children[0]
	if resource.Id != "aws_instance.web" || resource.Kind != tree.NodeResource || resource.Action != ui.ActionCreate {
		t.Fatalf("unexpected resource node: %#v", resource)
	}

	if !reflect.DeepEqual(resource.Payload, payload) {
		t.Fatalf("expected payload %#v, got %#v", payload, resource.Payload)
	}
}

func TestBuildTreeNodeNil(t *testing.T) {
	if got := buildTreeNode(nil); got != nil {
		t.Fatalf("expected nil tree node, got %#v", got)
	}
}

func TestBuildStatsCountsResourceActions(t *testing.T) {
	root := &planview.Node{
		Kind: planview.NodeGroup,
		Children: []*planview.Node{
			{Kind: planview.NodeGroup, Action: planview.ActionCreate},
			{Kind: planview.NodeResource, Action: planview.ActionCreate},
			{Kind: planview.NodeResource, Action: planview.ActionCreate},
			{Kind: planview.NodeResource, Action: planview.ActionUpdate},
			{Kind: planview.NodeResource, Action: planview.ActionDelete},
			{Kind: planview.NodeResource, Action: planview.ActionReplace},
			{Kind: planview.NodeResource, Action: planview.ActionNoOp},
			{Kind: planview.NodeResource, Action: planview.ActionError},
		},
	}

	got := buildStats(root)
	want := status.Stats{Create: 2, Update: 1, Delete: 1, Replace: 1, NoOp: 1, Errors: 1}

	if got != want {
		t.Fatalf("expected stats %#v, got %#v", want, got)
	}
}

func TestBuildStatsNil(t *testing.T) {
	got := buildStats(nil)
	if got != (status.Stats{}) {
		t.Fatalf("expected empty stats, got %#v", got)
	}
}

func TestBuildDetailsContent(t *testing.T) {
	changes := ui.ChangeSet{
		Before: map[string]any{"size": "small"},
		After:  map[string]any{"size": "large"},
	}
	payload := map[string]any{"address": "aws_instance.web"}

	tests := []struct {
		name string
		node *tree.Node
		want details.Content
	}{
		{
			name: "nil",
			want: details.Content{Kind: details.KindNone},
		},
		{
			name: "group",
			node: &tree.Node{Id: "create-group", Label: "Create", Kind: tree.NodeGroup},
			want: details.Content{Key: "create-group", Label: "Create", Kind: details.KindGroup},
		},
		{
			name: "resource",
			node: &tree.Node{
				Id:      "aws_instance.web",
				Label:   "aws_instance.web",
				Kind:    tree.NodeResource,
				Action:  ui.ActionUpdate,
				Changes: changes,
				Payload: payload,
			},
			want: details.Content{
				Key:     "aws_instance.web",
				Label:   "aws_instance.web",
				Kind:    details.KindResource,
				Changes: changes,
				Payload: payload,
			},
		},
		{
			name: "error resource",
			node: &tree.Node{Id: "diagnostic", Label: "Diagnostic", Kind: tree.NodeResource, Action: ui.ActionError},
			want: details.Content{Key: "diagnostic", Label: "Diagnostic", Kind: details.KindResource, IsError: true},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := buildDetailsContent(tc.node)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("expected content %#v, got %#v", tc.want, got)
			}
		})
	}
}

func TestBuildFilterOptions(t *testing.T) {
	groups := []*planview.Node{
		nil,
		{Label: "Create", LabelCount: "(2/3)", Action: planview.ActionCreate},
		{Label: "Errors", LabelCount: "(1/3)", Action: planview.ActionError},
	}

	got := buildFilterOptions(groups)
	want := []filter.Option{
		{Action: ui.ActionCreate, Label: "Create", Count: "(2/3)"},
		{Action: ui.ActionError, Label: "Errors", Count: "(1/3)"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected options %#v, got %#v", want, got)
	}
}

func TestConvertPlanAction(t *testing.T) {
	tests := []struct {
		name string
		in   planview.Action
		want ui.Action
	}{
		{name: "create", in: planview.ActionCreate, want: ui.ActionCreate},
		{name: "update", in: planview.ActionUpdate, want: ui.ActionUpdate},
		{name: "delete", in: planview.ActionDelete, want: ui.ActionDelete},
		{name: "replace", in: planview.ActionReplace, want: ui.ActionReplace},
		{name: "no-op", in: planview.ActionNoOp, want: ui.ActionNoOp},
		{name: "error", in: planview.ActionError, want: ui.ActionError},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := convertPlanAction(tc.in); got != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, got)
			}
		})
	}
}

func TestConvertPlanActionPanicsOnUnknown(t *testing.T) {
	requirePanic(t, func() {
		convertPlanAction(planview.Action("unknown"))
	})
}

func TestConvertPlanNodeKind(t *testing.T) {
	tests := []struct {
		name string
		in   planview.NodeKind
		want tree.NodeKind
	}{
		{name: "group", in: planview.NodeGroup, want: tree.NodeGroup},
		{name: "resource", in: planview.NodeResource, want: tree.NodeResource},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := convertPlanNodeKind(tc.in); got != tc.want {
				t.Fatalf("expected %d, got %d", tc.want, got)
			}
		})
	}
}

func TestConvertPlanNodeKindPanicsOnUnknown(t *testing.T) {
	requirePanic(t, func() {
		convertPlanNodeKind(planview.NodeKind(99))
	})
}

func requirePanic(t *testing.T, fn func()) {
	t.Helper()

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()

	fn()
}
