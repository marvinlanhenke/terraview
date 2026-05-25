package planview

import (
	"reflect"
	"testing"
)

func TestCompareChanges(t *testing.T) {
	tests := []struct {
		name       string
		before     map[string]any
		after      map[string]any
		wantBefore map[string]any
		wantAfter  map[string]any
	}{
		{
			name:       "scalar change",
			before:     map[string]any{"instance_type": "t3.micro"},
			after:      map[string]any{"instance_type": "t3.small"},
			wantBefore: map[string]any{"instance_type": "t3.micro"},
			wantAfter:  map[string]any{"instance_type": "t3.small"},
		},
		{
			name:       "added field",
			before:     map[string]any{"name": "api"},
			after:      map[string]any{"name": "api", "owner": "platform"},
			wantBefore: map[string]any{"owner": nil},
			wantAfter:  map[string]any{"owner": "platform"},
		},
		{
			name:       "removed field",
			before:     map[string]any{"name": "api", "owner": "platform"},
			after:      map[string]any{"name": "api"},
			wantBefore: map[string]any{"owner": "platform"},
			wantAfter:  map[string]any{"owner": nil},
		},
		{
			name: "nested map change",
			before: map[string]any{
				"tags": map[string]any{"env": "prod", "owner": "app"},
			},
			after: map[string]any{
				"tags": map[string]any{"env": "prod", "owner": "platform"},
			},
			wantBefore: map[string]any{
				"tags": map[string]any{"owner": "app"},
			},
			wantAfter: map[string]any{
				"tags": map[string]any{"owner": "platform"},
			},
		},
		{
			name:       "unchanged fields omitted",
			before:     map[string]any{"name": "api", "env": "prod"},
			after:      map[string]any{"name": "api", "env": "prod"},
			wantBefore: map[string]any{},
			wantAfter:  map[string]any{},
		},
		{
			name: "slices compare as whole value",
			before: map[string]any{
				"cidr_blocks": []any{"10.0.0.0/16"},
			},
			after: map[string]any{
				"cidr_blocks": []any{"10.0.0.0/16", "10.1.0.0/16"},
			},
			wantBefore: map[string]any{
				"cidr_blocks": []any{"10.0.0.0/16"},
			},
			wantAfter: map[string]any{
				"cidr_blocks": []any{"10.0.0.0/16", "10.1.0.0/16"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := compareChanges(tc.before, tc.after)

			if !reflect.DeepEqual(got.before, tc.wantBefore) {
				t.Fatalf("unexpected before diff: got %#v want %#v", got.before, tc.wantBefore)
			}

			if !reflect.DeepEqual(got.after, tc.wantAfter) {
				t.Fatalf("unexpected after diff: got %#v want %#v", got.after, tc.wantAfter)
			}
		})
	}
}

func TestNodeChangeSetsReturnClones(t *testing.T) {
	node := &Node{
		changes: changeSet{
			before: map[string]any{
				"tags": map[string]any{"env": "prod"},
				"ports": []any{float64(80), float64(443)},
			},
			after: map[string]any{
				"tags": map[string]any{"env": "stage"},
				"ports": []any{float64(8080), float64(8443)},
			},
		},
	}

	before := node.ChangeSetBefore()
	after := node.ChangeSetAfter()

	requireMap(t, before["tags"])["env"] = "mutated"
	requireMap(t, after["tags"])["env"] = "mutated"
	before["new"] = "value"
	after["new"] = "value"

	beforePorts := before["ports"].([]any)
	afterPorts := after["ports"].([]any)
	beforePorts[0] = float64(1)
	afterPorts[0] = float64(2)

	refetchedBefore := node.ChangeSetBefore()
	refetchedAfter := node.ChangeSetAfter()

	if got := requireMap(t, refetchedBefore["tags"])["env"]; got != "prod" {
		t.Fatalf("expected original before tags to remain unchanged, got %#v", got)
	}

	if got := requireMap(t, refetchedAfter["tags"])["env"]; got != "stage" {
		t.Fatalf("expected original after tags to remain unchanged, got %#v", got)
	}

	if _, ok := refetchedBefore["new"]; ok {
		t.Fatal("expected cloned before map mutation not to leak back into node")
	}

	if _, ok := refetchedAfter["new"]; ok {
		t.Fatal("expected cloned after map mutation not to leak back into node")
	}

	if got := refetchedBefore["ports"].([]any)[0]; got != float64(80) {
		t.Fatalf("expected original before slice to remain unchanged, got %#v", got)
	}

	if got := refetchedAfter["ports"].([]any)[0]; got != float64(8080) {
		t.Fatalf("expected original after slice to remain unchanged, got %#v", got)
	}
}
