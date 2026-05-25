package planview

import "testing"

func TestParseAction(t *testing.T) {
	tests := []struct {
		name    string
		actions []string
		want    Action
		wantErr bool
	}{
		{name: "create", actions: []string{"create"}, want: ActionCreate},
		{name: "update", actions: []string{"update"}, want: ActionUpdate},
		{name: "delete", actions: []string{"delete"}, want: ActionDelete},
		{name: "no-op", actions: []string{"no-op"}, want: ActionNoOp},
		{name: "replace delete create", actions: []string{"delete", "create"}, want: ActionReplace},
		{name: "replace create delete", actions: []string{"create", "delete"}, want: ActionReplace},
		{name: "empty", actions: nil, want: ActionError, wantErr: true},
		{name: "unknown", actions: []string{"unknown"}, want: ActionError, wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseAction(tc.actions)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
			} else if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if got != tc.want {
				t.Fatalf("expected action %q, got %q", tc.want, got)
			}
		})
	}
}
