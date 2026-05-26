package planview

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/marvinlanhenke/terraview/internal/terraform"
)

// discardLogger returns a no-op logger suitable for tests.
func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestFromTerraformAllFixtures(t *testing.T) {
	for _, name := range fixtureNames(t) {
		plan := parseFixture(t, name)

		t.Run(strings.TrimSuffix(name, filepath.Ext(name)), func(t *testing.T) {
			root, err := FromTerraform(plan, discardLogger())
			if err != nil {
				t.Fatalf("expected fixture to build successfully: %v", err)
			}

			requireRoot(t, root, plan.TerraformVersion)
		})
	}
}

func TestFromTerraformFixtureContracts(t *testing.T) {
	tests := []struct {
		name   string
		assert func(t *testing.T, root *Node)
	}{
		{
			name: "create-only.json",
			assert: func(t *testing.T, root *Node) {
				create := requireGroup(t, root, ActionCreate)
				requireGroupCount(t, create, 7, "(7/7)")

				for _, action := range []Action{ActionUpdate, ActionDelete, ActionReplace, ActionNoOp, ActionError} {
					requireGroupCount(t, requireGroup(t, root, action), 0, "")
				}

				cluster := requireChild(t, create, "module.compute.aws_ecs_cluster.main")
				if !cluster.IsResource() {
					t.Fatal("expected created child to be a resource node")
				}

				payload, ok := cluster.Payload.(terraform.ResourceChange)
				if !ok {
					t.Fatalf("expected resource payload type, got %T", cluster.Payload)
				}

				if payload.ModuleAddress != "module.compute" {
					t.Fatalf("expected compute module address, got %q", payload.ModuleAddress)
				}
			},
		},
		{
			name: "update-only.json",
			assert: func(t *testing.T, root *Node) {
				update := requireGroup(t, root, ActionUpdate)
				requireGroupCount(t, update, 2, "(2/2)")

				securityGroup := requireChild(t, update, "aws_security_group.api")

				beforeTags := requireMap(t, securityGroup.ChangeSetBefore()["tags"])
				afterTags := requireMap(t, securityGroup.ChangeSetAfter()["tags"])

				if beforeTags["owner"] != nil {
					t.Fatalf("expected before owner tag to be nil, got %#v", beforeTags["owner"])
				}

				if afterTags["owner"] != "platform" {
					t.Fatalf("expected after owner tag to be platform, got %#v", afterTags["owner"])
				}

				if got := len(beforeTags); got != 1 {
					t.Fatalf("expected only owner tag in before diff, got %d entries", got)
				}

				if got := len(afterTags); got != 1 {
					t.Fatalf("expected only owner tag in after diff, got %d entries", got)
				}
			},
		},
		{
			name: "delete-only.json",
			assert: func(t *testing.T, root *Node) {
				deleteGroup := requireGroup(t, root, ActionDelete)
				requireGroupCount(t, deleteGroup, 2, "(2/2)")

				bucket := requireChild(t, deleteGroup, "aws_s3_bucket.legacy_artifacts")
				before := bucket.ChangeSetBefore()
				after := bucket.ChangeSetAfter()

				if before["bucket"] != "example-prod-legacy-artifacts" {
					t.Fatalf("unexpected before bucket value %#v", before["bucket"])
				}

				if after["bucket"] != nil {
					t.Fatalf("expected removed bucket after value to be nil, got %#v", after["bucket"])
				}
			},
		},
		{
			name: "replace-only.json",
			assert: func(t *testing.T, root *Node) {
				replace := requireGroup(t, root, ActionReplace)
				requireGroupCount(t, replace, 2, "(2/2)")

				database := requireChild(t, replace, "aws_db_instance.primary")

				if got := database.ChangeSetBefore()["engine_version"]; got != "14.9" {
					t.Fatalf("expected engine_version before diff to be 14.9, got %#v", got)
				}

				if got := database.ChangeSetAfter()["engine_version"]; got != "16.1" {
					t.Fatalf("expected engine_version after diff to be 16.1, got %#v", got)
				}
			},
		},
		{
			name: "noop-only.json",
			assert: func(t *testing.T, root *Node) {
				noOp := requireGroup(t, root, ActionNoOp)
				requireGroupCount(t, noOp, 2, "(2/2)")

				role := requireChild(t, noOp, "aws_iam_role.app")
				if got := len(role.ChangeSetBefore()); got != 0 {
					t.Fatalf("expected no before diff entries, got %d", got)
				}

				if got := len(role.ChangeSetAfter()); got != 0 {
					t.Fatalf("expected no after diff entries, got %d", got)
				}
			},
		},
		{
			name: "errors-only.json",
			assert: func(t *testing.T, root *Node) {
				errorGroup := requireGroup(t, root, ActionError)
				requireGroupCount(t, errorGroup, 3, "(3/3)")

				for _, action := range []Action{ActionCreate, ActionUpdate, ActionDelete, ActionReplace, ActionNoOp} {
					requireGroupCount(t, requireGroup(t, root, action), 0, "")
				}

				diagnostic := requireChildByLabel(t, errorGroup, "Invalid value for variable")
				if _, ok := diagnostic.Payload.(terraform.Diagnostic); !ok {
					t.Fatalf("expected diagnostic payload type, got %T", diagnostic.Payload)
				}

				if diagnostic.ChangeSetBefore() != nil {
					t.Fatalf("expected diagnostic before changes to be nil, got %#v", diagnostic.ChangeSetBefore())
				}
			},
		},
		{
			name: "errors-partial.json",
			assert: func(t *testing.T, root *Node) {
				requireGroupCount(t, requireGroup(t, root, ActionCreate), 1, "(1/3)")
				requireGroupCount(t, requireGroup(t, root, ActionUpdate), 1, "(1/3)")
				requireGroupCount(t, requireGroup(t, root, ActionError), 1, "(1/3)")

				errorNode := requireChildByLabel(t, requireGroup(t, root, ActionError), "Error reading IAM policy document")
				if _, ok := errorNode.Payload.(terraform.Diagnostic); !ok {
					t.Fatalf("expected diagnostic payload type, got %T", errorNode.Payload)
				}
			},
		},
		{
			name: "mixed-single.json",
			assert: func(t *testing.T, root *Node) {
				requireGroupCount(t, requireGroup(t, root, ActionCreate), 1, "(1/5)")
				requireGroupCount(t, requireGroup(t, root, ActionUpdate), 1, "(1/5)")
				requireGroupCount(t, requireGroup(t, root, ActionDelete), 1, "(1/5)")
				requireGroupCount(t, requireGroup(t, root, ActionReplace), 1, "(1/5)")
				requireGroupCount(t, requireGroup(t, root, ActionNoOp), 1, "(1/5)")
				requireGroupCount(t, requireGroup(t, root, ActionError), 0, "")

				update := requireChild(t, requireGroup(t, root, ActionUpdate), "module.compute.aws_instance.web")
				payload, ok := update.Payload.(terraform.ResourceChange)
				if !ok {
					t.Fatalf("expected resource payload type, got %T", update.Payload)
				}

				if payload.ModuleAddress != "module.compute" {
					t.Fatalf("expected compute module address, got %q", payload.ModuleAddress)
				}
			},
		},
		{
			name: "mixed-multiple.json",
			assert: func(t *testing.T, root *Node) {
				requireGroupCount(t, requireGroup(t, root, ActionCreate), 3, "(3/11)")
				requireGroupCount(t, requireGroup(t, root, ActionUpdate), 3, "(3/11)")
				requireGroupCount(t, requireGroup(t, root, ActionDelete), 2, "(2/11)")
				requireGroupCount(t, requireGroup(t, root, ActionReplace), 2, "(2/11)")
				requireGroupCount(t, requireGroup(t, root, ActionNoOp), 1, "(1/11)")
				requireGroupCount(t, requireGroup(t, root, ActionError), 0, "")

				replace := requireChild(t, requireGroup(t, root, ActionReplace), "module.network.aws_subnet.private_a")
				replacePayload, ok := replace.Payload.(terraform.ResourceChange)
				if !ok {
					t.Fatalf("expected resource payload type, got %T", replace.Payload)
				}

				if replacePayload.ModuleAddress != "module.network" {
					t.Fatalf("expected network module address, got %q", replacePayload.ModuleAddress)
				}

				update := requireChild(t, requireGroup(t, root, ActionUpdate), "module.database.aws_db_instance.primary")
				updatePayload, ok := update.Payload.(terraform.ResourceChange)
				if !ok {
					t.Fatalf("expected resource payload type, got %T", update.Payload)
				}

				if updatePayload.ModuleAddress != "module.database" {
					t.Fatalf("expected database module address, got %q", updatePayload.ModuleAddress)
				}
			},
		},
		{
			name: "output-changes.json",
			assert: func(t *testing.T, root *Node) {
				for _, action := range []Action{ActionCreate, ActionUpdate, ActionDelete, ActionReplace, ActionNoOp, ActionError} {
					requireGroupCount(t, requireGroup(t, root, action), 0, "")
				}
			},
		},
	}

	for _, tc := range tests {
		plan := parseFixture(t, tc.name)

		t.Run(strings.TrimSuffix(tc.name, filepath.Ext(tc.name)), func(t *testing.T) {
			root, err := FromTerraform(plan, discardLogger())
			if err != nil {
				t.Fatalf("expected fixture to build successfully: %v", err)
			}

			tc.assert(t, root)
		})
	}
}

func TestFromTerraformIgnoresNonErrorDiagnostics(t *testing.T) {
	plan := terraform.Plan{
		TerraformVersion: "1.8.5",
		Diagnostics: []terraform.Diagnostic{
			{Severity: "warning", Summary: "Skipped warning"},
			{Severity: "error", Summary: "Real error"},
		},
	}

	root, err := FromTerraform(plan, discardLogger())
	if err != nil {
		t.Fatalf("expected diagnostics-only plan to build successfully: %v", err)
	}

	errorGroup := requireGroup(t, root, ActionError)
	requireGroupCount(t, errorGroup, 1, "(1/1)")

	if got := errorGroup.Children[0].Label; got != "Real error" {
		t.Fatalf("expected only error diagnostic to be rendered, got %q", got)
	}
}

func TestFromTerraformInvalidActionReturnsError(t *testing.T) {
	plan := terraform.Plan{
		TerraformVersion: "1.8.5",
		ResourceChanges: []terraform.ResourceChange{
			{
				Address: "aws_s3_bucket.invalid",
				Change:  terraform.Change{Actions: []string{"unexpected"}},
			},
		},
	}

	if _, err := FromTerraform(plan, discardLogger()); err == nil {
		t.Fatal("expected invalid action to return an error")
	}
}

func TestFromTerraformEmptyPlan(t *testing.T) {
	plan := terraform.Plan{TerraformVersion: "1.8.5"}

	root, err := FromTerraform(plan, discardLogger())
	if err != nil {
		t.Fatalf("expected empty plan to build successfully: %v", err)
	}

	requireRoot(t, root, plan.TerraformVersion)

	for _, action := range orderedActions() {
		requireGroupCount(t, requireGroup(t, root, action), 0, "")
	}
}

func fixtureNames(t *testing.T) []string {
	t.Helper()

	dir := fixturesDir(t)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read fixture dir %q: %v", dir, err)
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		names = append(names, entry.Name())
	}

	return names
}

func parseFixture(t *testing.T, name string) terraform.Plan {
	t.Helper()

	data, err := os.ReadFile(filepath.Join(fixturesDir(t), name))
	if err != nil {
		t.Fatalf("failed to read fixture %q: %v", name, err)
	}

	plan, err := terraform.Parse(data)
	if err != nil {
		t.Fatalf("failed to parse fixture %q: %v", name, err)
	}

	return plan
}

func fixturesDir(t *testing.T) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve test file path")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "testdata", "plans"))
}

func orderedActions() []Action {
	return []Action{ActionCreate, ActionUpdate, ActionDelete, ActionReplace, ActionNoOp, ActionError}
}

func requireRoot(t *testing.T, root *Node, terraformVersion string) {
	t.Helper()

	if root == nil {
		t.Fatal("expected root node")
	}

	if root.Kind != NodeGroup {
		t.Fatalf("expected root kind %v, got %v", NodeGroup, root.Kind)
	}

	if root.Action != ActionNoOp {
		t.Fatalf("expected root action %q, got %q", ActionNoOp, root.Action)
	}

	if !strings.Contains(root.Label, terraformVersion) {
		t.Fatalf("expected root label to contain version %q, got %q", terraformVersion, root.Label)
	}

	if !root.HasChildren() {
		t.Fatal("expected root to have action groups")
	}

	if root.IsResource() {
		t.Fatal("expected root not to be a resource node")
	}

	if got := len(root.Children); got != len(orderedActions()) {
		t.Fatalf("expected %d action groups, got %d", len(orderedActions()), got)
	}

	for i, want := range orderedActions() {
		group := root.Children[i]

		if group.Action != want {
			t.Fatalf("expected group %d action %q, got %q", i, want, group.Action)
		}

		if group.Kind != NodeGroup {
			t.Fatalf("expected group %d to be NodeGroup, got %v", i, group.Kind)
		}

		if group.Label == "" {
			t.Fatalf("expected group %d label to be populated", i)
		}
	}
}

func requireGroup(t *testing.T, root *Node, action Action) *Node {
	t.Helper()

	for _, child := range root.Children {
		if child.Action == action {
			return child
		}
	}

	t.Fatalf("group for action %q not found", action)
	return nil
}

func requireGroupCount(t *testing.T, group *Node, wantChildren int, wantLabelCount string) {
	t.Helper()

	if got := len(group.Children); got != wantChildren {
		t.Fatalf("expected %q group to have %d children, got %d", group.Action, wantChildren, got)
	}

	if group.LabelCount != wantLabelCount {
		t.Fatalf("expected %q group label count %q, got %q", group.Action, wantLabelCount, group.LabelCount)
	}
}

func requireChild(t *testing.T, group *Node, id string) *Node {
	t.Helper()

	for _, child := range group.Children {
		if child.Id == id {
			return child
		}
	}

	t.Fatalf("child %q not found in %q group", id, group.Action)
	return nil
}

func requireChildByLabel(t *testing.T, group *Node, label string) *Node {
	t.Helper()

	for _, child := range group.Children {
		if child.Label == label {
			return child
		}
	}

	t.Fatalf("child with label %q not found in %q group", label, group.Action)
	return nil
}

func requireMap(t *testing.T, value any) map[string]any {
	t.Helper()

	m, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", value)
	}

	return m
}
