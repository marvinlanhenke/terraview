package terraform_test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/marvinlanhenke/terraview/internal/terraform"
)

func TestParseAllFixtures(t *testing.T) {
	fixtureNames := fixtureNames(t)

	if len(fixtureNames) == 0 {
		t.Fatal("expected at least one plan fixture")
	}

	for _, name := range fixtureNames {
		t.Run(strings.TrimSuffix(name, filepath.Ext(name)), func(t *testing.T) {
			plan := parseFixture(t, name)

			if plan.FormatVersion == "" {
				t.Fatal("expected format version to be populated")
			}

			if plan.TerraformVersion == "" {
				t.Fatal("expected terraform version to be populated")
			}

			if len(plan.Configuration.ProviderConfig) == 0 {
				t.Fatal("expected provider config to be populated")
			}
		})
	}
}

func TestParseFixtureContracts(t *testing.T) {
	tests := []struct {
		name   string
		assert func(t *testing.T, plan terraform.Plan)
	}{
		{
			name: "create-only.json",
			assert: func(t *testing.T, plan terraform.Plan) {
				requireResourceActionCounts(t, plan.ResourceChanges, map[string]int{"create": 7})
				requireOutputActionCounts(t, plan.OutputChanges, map[string]int{"create": 3})

				if got := len(plan.PlannedValues.RootModule.Resources); got != 4 {
					t.Fatalf("expected 4 root planned resources, got %d", got)
				}

				if got := len(plan.PlannedValues.RootModule.ChildModules); got != 2 {
					t.Fatalf("expected 2 planned child modules, got %d", got)
				}

				network := plannedModuleByAddress(t, plan.PlannedValues.RootModule.ChildModules, "module.network")
				if got := len(network.Resources); got != 2 {
					t.Fatalf("expected 2 planned network resources, got %d", got)
				}

				compute := plannedModuleByAddress(t, plan.PlannedValues.RootModule.ChildModules, "module.compute")
				if got := len(compute.Resources); got != 1 {
					t.Fatalf("expected 1 planned compute resource, got %d", got)
				}

				networkCall := moduleCallByName(t, plan.Configuration.RootModule.ModuleCalls, "network")
				if networkCall.Source != "./modules/network" {
					t.Fatalf("expected network module source ./modules/network, got %q", networkCall.Source)
				}

				if got := len(networkCall.Module.Resources); got != 2 {
					t.Fatalf("expected 2 configured network resources, got %d", got)
				}

				computeCall := moduleCallByName(t, plan.Configuration.RootModule.ModuleCalls, "compute")
				if computeCall.Source != "./modules/compute" {
					t.Fatalf("expected compute module source ./modules/compute, got %q", computeCall.Source)
				}

				if got := len(computeCall.Module.Resources); got != 1 {
					t.Fatalf("expected 1 configured compute resource, got %d", got)
				}

				cluster := resourceChangeByAddress(t, plan, "module.compute.aws_ecs_cluster.main")
				requireActions(t, cluster.Change.Actions, "create")
				if cluster.ModuleAddress != "module.compute" {
					t.Fatalf("expected compute module address, got %q", cluster.ModuleAddress)
				}

				clusterOutput := outputChangeByName(t, plan, "ecs_cluster_arn")
				requireActions(t, clusterOutput.Actions, "create")
				if clusterOutput.After != nil {
					t.Fatalf("expected ecs_cluster_arn after to be nil, got %#v", clusterOutput.After)
				}

				if !requireBool(t, clusterOutput.AfterUnknown) {
					t.Fatal("expected ecs_cluster_arn after_unknown to be true")
				}
			},
		},
		{
			name: "update-only.json",
			assert: func(t *testing.T, plan terraform.Plan) {
				requireResourceActionCounts(t, plan.ResourceChanges, map[string]int{"update": 2})
				requireOutputActionCounts(t, plan.OutputChanges, map[string]int{"update": 1})

				securityGroup := resourceChangeByAddress(t, plan, "aws_security_group.api")
				requireActions(t, securityGroup.Change.Actions, "update")

				tags := requireMap(t, securityGroup.Change.After["tags"])
				if tags["owner"] != "platform" {
					t.Fatalf("expected owner tag to be platform, got %#v", tags["owner"])
				}

				if got := len(securityGroup.Change.AfterUnknown); got != 0 {
					t.Fatalf("expected no unknown values for security group update, got %d", got)
				}

				output := outputChangeByName(t, plan, "log_retention_days")
				requireActions(t, output.Actions, "update")

				if got := requireFloat64(t, output.Before); got != 30 {
					t.Fatalf("expected output before to be 30, got %v", got)
				}

				if got := requireFloat64(t, output.After); got != 90 {
					t.Fatalf("expected output after to be 90, got %v", got)
				}

				if requireBool(t, output.AfterUnknown) {
					t.Fatal("expected output after_unknown to be false")
				}
			},
		},
		{
			name: "delete-only.json",
			assert: func(t *testing.T, plan terraform.Plan) {
				requireResourceActionCounts(t, plan.ResourceChanges, map[string]int{"delete": 2})
				requireOutputActionCounts(t, plan.OutputChanges, map[string]int{"delete": 1})

				if got := len(plan.PlannedValues.RootModule.Resources); got != 0 {
					t.Fatalf("expected no planned root resources, got %d", got)
				}

				for _, change := range plan.ResourceChanges {
					requireActions(t, change.Change.Actions, "delete")

					if change.Change.After != nil {
						t.Fatalf("expected delete after value to be nil for %s, got %#v", change.Address, change.Change.After)
					}
				}

				output := outputChangeByName(t, plan, "legacy_nat_public_ip")
				requireActions(t, output.Actions, "delete")

				if output.After != nil {
					t.Fatalf("expected deleted output after to be nil, got %#v", output.After)
				}
			},
		},
		{
			name: "replace-only.json",
			assert: func(t *testing.T, plan terraform.Plan) {
				requireResourceActionCounts(t, plan.ResourceChanges, map[string]int{"delete,create": 2})
				requireOutputActionCounts(t, plan.OutputChanges, map[string]int{"delete,create": 1})

				database := resourceChangeByAddress(t, plan, "aws_db_instance.primary")
				requireActions(t, database.Change.Actions, "delete", "create")

				if !reflect.DeepEqual(database.Change.ReplacePaths, [][]string{{"engine_version"}}) {
					t.Fatalf("expected engine_version replace path, got %#v", database.Change.ReplacePaths)
				}

				if !requireBool(t, database.Change.AfterUnknown["resource_id"]) {
					t.Fatal("expected database resource_id to remain unknown")
				}

				subnet := resourceChangeByAddress(t, plan, "aws_subnet.private_a")
				requireActions(t, subnet.Change.Actions, "delete", "create")

				if !reflect.DeepEqual(subnet.Change.ReplacePaths, [][]string{{"cidr_block"}}) {
					t.Fatalf("expected cidr_block replace path, got %#v", subnet.Change.ReplacePaths)
				}

				output := outputChangeByName(t, plan, "db_endpoint")
				requireActions(t, output.Actions, "delete", "create")

				if output.After != nil {
					t.Fatalf("expected replacement output after to be nil, got %#v", output.After)
				}

				if !requireBool(t, output.AfterUnknown) {
					t.Fatal("expected db_endpoint after_unknown to be true")
				}
			},
		},
		{
			name: "noop-only.json",
			assert: func(t *testing.T, plan terraform.Plan) {
				requireResourceActionCounts(t, plan.ResourceChanges, map[string]int{"no-op": 2})
				requireOutputActionCounts(t, plan.OutputChanges, map[string]int{"no-op": 1})

				role := resourceChangeByAddress(t, plan, "aws_iam_role.app")
				requireActions(t, role.Change.Actions, "no-op")

				if !reflect.DeepEqual(role.Change.Before, role.Change.After) {
					t.Fatalf("expected no-op before and after to match, got before=%#v after=%#v", role.Change.Before, role.Change.After)
				}
			},
		},
		{
			name: "errors-only.json",
			assert: func(t *testing.T, plan terraform.Plan) {
				if got := len(plan.ResourceChanges); got != 0 {
					t.Fatalf("expected no resource changes, got %d", got)
				}

				if got := len(plan.OutputChanges); got != 0 {
					t.Fatalf("expected no output changes, got %d", got)
				}

				if got := len(plan.Diagnostics); got != 3 {
					t.Fatalf("expected 3 diagnostics, got %d", got)
				}

				for _, diagnostic := range plan.Diagnostics {
					if diagnostic.Severity != "error" {
						t.Fatalf("expected error diagnostic severity, got %q", diagnostic.Severity)
					}

					if len(diagnostic.Range) == 0 {
						t.Fatalf("expected diagnostic range for %q", diagnostic.Summary)
					}

					if len(diagnostic.Snippet) == 0 {
						t.Fatalf("expected diagnostic snippet for %q", diagnostic.Summary)
					}
				}

				provider := plan.Configuration.ProviderConfig["aws"]
				if provider.FullName != "registry.terraform.io/hashicorp/aws" {
					t.Fatalf("expected provider full name to be preserved, got %q", provider.FullName)
				}
			},
		},
		{
			name: "errors-partial.json",
			assert: func(t *testing.T, plan terraform.Plan) {
				requireResourceActionCounts(t, plan.ResourceChanges, map[string]int{"create": 1, "update": 1})
				requireOutputActionCounts(t, plan.OutputChanges, map[string]int{"create": 1})

				if got := len(plan.Diagnostics); got != 1 {
					t.Fatalf("expected 1 diagnostic, got %d", got)
				}

				diagnostic := plan.Diagnostics[0]
				if diagnostic.Summary != "Error reading IAM policy document" {
					t.Fatalf("unexpected diagnostic summary %q", diagnostic.Summary)
				}

				bucket := resourceChangeByAddress(t, plan, "aws_s3_bucket.logs")
				requireActions(t, bucket.Change.Actions, "create")

				if !requireBool(t, bucket.Change.AfterUnknown["bucket_domain_name"]) {
					t.Fatal("expected bucket_domain_name to remain unknown")
				}
			},
		},
		{
			name: "output-changes.json",
			assert: func(t *testing.T, plan terraform.Plan) {
				if got := len(plan.ResourceChanges); got != 0 {
					t.Fatalf("expected no resource changes, got %d", got)
				}

				requireOutputActionCounts(t, plan.OutputChanges, map[string]int{
					"no-op":  1,
					"create": 1,
					"delete": 1,
					"update": 1,
				})

				if got := len(plan.PlannedValues.RootModule.Resources); got != 1 {
					t.Fatalf("expected 1 planned root resource, got %d", got)
				}

				platform := outputChangeByName(t, plan, "platform_metadata")
				requireActions(t, platform.Actions, "update")

				after := requireMap(t, platform.After)
				if after["managed_by"] != "terraform" {
					t.Fatalf("expected managed_by to be terraform, got %#v", after["managed_by"])
				}
			},
		},
		{
			name: "mixed-single.json",
			assert: func(t *testing.T, plan terraform.Plan) {
				requireResourceActionCounts(t, plan.ResourceChanges, map[string]int{
					"no-op":         1,
					"create":        1,
					"update":        1,
					"delete":        1,
					"delete,create": 1,
				})
				requireOutputActionCounts(t, plan.OutputChanges, map[string]int{"create": 1, "update": 1})

				if got := len(plan.PlannedValues.RootModule.Resources); got != 2 {
					t.Fatalf("expected 2 root planned resources, got %d", got)
				}

				compute := plannedModuleByAddress(t, plan.PlannedValues.RootModule.ChildModules, "module.compute")
				if got := len(compute.Resources); got != 1 {
					t.Fatalf("expected 1 planned compute resource, got %d", got)
				}

				computeChange := resourceChangeByAddress(t, plan, "module.compute.aws_instance.web")
				requireActions(t, computeChange.Change.Actions, "update")
				if computeChange.ModuleAddress != "module.compute" {
					t.Fatalf("expected compute module address, got %q", computeChange.ModuleAddress)
				}

				database := resourceChangeByAddress(t, plan, "aws_db_instance.main")
				requireActions(t, database.Change.Actions, "delete", "create")

				if !reflect.DeepEqual(database.Change.ReplacePaths, [][]string{{"engine_version"}}) {
					t.Fatalf("expected engine_version replace path, got %#v", database.Change.ReplacePaths)
				}

				computeCall := moduleCallByName(t, plan.Configuration.RootModule.ModuleCalls, "compute")
				if computeCall.Source != "./modules/compute" {
					t.Fatalf("expected compute module source ./modules/compute, got %q", computeCall.Source)
				}
			},
		},
		{
			name: "mixed-multiple.json",
			assert: func(t *testing.T, plan terraform.Plan) {
				requireResourceActionCounts(t, plan.ResourceChanges, map[string]int{
					"no-op":         1,
					"create":        3,
					"update":        3,
					"delete":        2,
					"delete,create": 2,
				})
				requireOutputActionCounts(t, plan.OutputChanges, map[string]int{
					"create":        2,
					"update":        1,
					"delete,create": 1,
					"delete":        1,
				})

				if got := len(plan.PlannedValues.RootModule.Resources); got != 4 {
					t.Fatalf("expected 4 root planned resources, got %d", got)
				}

				network := plannedModuleByAddress(t, plan.PlannedValues.RootModule.ChildModules, "module.network")
				if got := len(network.Resources); got != 2 {
					t.Fatalf("expected 2 planned network resources, got %d", got)
				}

				databaseModule := plannedModuleByAddress(t, plan.PlannedValues.RootModule.ChildModules, "module.database")
				if got := len(databaseModule.Resources); got != 1 {
					t.Fatalf("expected 1 planned database resource, got %d", got)
				}

				privateSubnet := resourceChangeByAddress(t, plan, "module.network.aws_subnet.private_a")
				requireActions(t, privateSubnet.Change.Actions, "delete", "create")
				if privateSubnet.ModuleAddress != "module.network" {
					t.Fatalf("expected network module address, got %q", privateSubnet.ModuleAddress)
				}

				if !reflect.DeepEqual(privateSubnet.Change.ReplacePaths, [][]string{{"cidr_block"}}) {
					t.Fatalf("expected cidr_block replace path, got %#v", privateSubnet.Change.ReplacePaths)
				}

				kmsKey := resourceChangeByAddress(t, plan, "aws_kms_key.app")
				requireActions(t, kmsKey.Change.Actions, "delete", "create")
				if !reflect.DeepEqual(kmsKey.Change.ReplacePaths, [][]string{{"deletion_window_in_days"}}) {
					t.Fatalf("expected deletion_window_in_days replace path, got %#v", kmsKey.Change.ReplacePaths)
				}

				database := resourceChangeByAddress(t, plan, "module.database.aws_db_instance.primary")
				requireActions(t, database.Change.Actions, "update")
				if database.ModuleAddress != "module.database" {
					t.Fatalf("expected database module address, got %q", database.ModuleAddress)
				}

				databaseCall := moduleCallByName(t, plan.Configuration.RootModule.ModuleCalls, "database")
				if databaseCall.Source != "./modules/database" {
					t.Fatalf("expected database module source ./modules/database, got %q", databaseCall.Source)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(strings.TrimSuffix(tc.name, filepath.Ext(tc.name)), func(t *testing.T) {
			plan := parseFixture(t, tc.name)
			tc.assert(t, plan)
		})
	}
}

func TestParseInvalidJSON(t *testing.T) {
	truncatedFixture := readFixture(t, "mixed-single.json")

	tests := []struct {
		name  string
		input []byte
	}{
		{
			name:  "malformed",
			input: []byte("{not-json}"),
		},
		{
			name:  "truncated fixture",
			input: truncatedFixture[:len(truncatedFixture)/2],
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := terraform.Parse(tc.input); err == nil {
				t.Fatal("expected parse error")
			}
		})
	}
}

func fixtureNames(t *testing.T) []string {
	t.Helper()

	dir := fixturesDir(t)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read fixture dir %q: %v", dir, err)
	}

	var names []string
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		names = append(names, entry.Name())
	}

	sort.Strings(names)

	return names
}

func parseFixture(t *testing.T, name string) terraform.Plan {
	t.Helper()

	data := readFixture(t, name)
	plan, err := terraform.Parse(data)
	if err != nil {
		t.Fatalf("failed to parse %s: %v", name, err)
	}

	return plan
}

func readFixture(t *testing.T, name string) []byte {
	t.Helper()

	path := filepath.Join(fixturesDir(t), name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read fixture %q: %v", path, err)
	}

	return data
}

func fixturesDir(t *testing.T) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve test file path")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "testdata", "plans"))
}

func plannedModuleByAddress(t *testing.T, modules []terraform.PlannedModule, address string) terraform.PlannedModule {
	t.Helper()

	for _, module := range modules {
		if module.Address == address {
			return module
		}
	}

	t.Fatalf("planned module %q not found", address)

	return terraform.PlannedModule{}
}

func moduleCallByName(t *testing.T, calls map[string]terraform.ConfigModuleCall, name string) terraform.ConfigModuleCall {
	t.Helper()

	call, ok := calls[name]
	if !ok {
		t.Fatalf("module call %q not found", name)
	}

	return call
}

func resourceChangeByAddress(t *testing.T, plan terraform.Plan, address string) terraform.ResourceChange {
	t.Helper()

	for _, change := range plan.ResourceChanges {
		if change.Address == address {
			return change
		}
	}

	t.Fatalf("resource change %q not found", address)

	return terraform.ResourceChange{}
}

func outputChangeByName(t *testing.T, plan terraform.Plan, name string) terraform.OutputChange {
	t.Helper()

	change, ok := plan.OutputChanges[name]
	if !ok {
		t.Fatalf("output change %q not found", name)
	}

	return change
}

func requireResourceActionCounts(t *testing.T, changes []terraform.ResourceChange, want map[string]int) {
	t.Helper()

	got := make(map[string]int)
	for _, change := range changes {
		got[actionKey(change.Change.Actions)]++
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected resource action counts: got %s want %s", formatCounts(got), formatCounts(want))
	}
}

func requireOutputActionCounts(t *testing.T, changes map[string]terraform.OutputChange, want map[string]int) {
	t.Helper()

	got := make(map[string]int)
	for _, change := range changes {
		got[actionKey(change.Actions)]++
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected output action counts: got %s want %s", formatCounts(got), formatCounts(want))
	}
}

func requireActions(t *testing.T, got []string, want ...string) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected actions: got %v want %v", got, want)
	}
}

func requireMap(t *testing.T, value any) map[string]any {
	t.Helper()

	m, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", value)
	}

	return m
}

func requireBool(t *testing.T, value any) bool {
	t.Helper()

	b, ok := value.(bool)
	if !ok {
		t.Fatalf("expected bool, got %T", value)
	}

	return b
}

func requireFloat64(t *testing.T, value any) float64 {
	t.Helper()

	f, ok := value.(float64)
	if !ok {
		t.Fatalf("expected float64, got %T", value)
	}

	return f
}

func actionKey(actions []string) string {
	return strings.Join(actions, ",")
}

func formatCounts(counts map[string]int) string {
	if len(counts) == 0 {
		return "{}"
	}

	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s:%d", key, counts[key]))
	}

	return "{" + strings.Join(parts, ", ") + "}"
}
