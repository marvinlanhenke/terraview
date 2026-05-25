// Package terraform defines the subset of Terraform plan JSON used by Terraview.
package terraform

// Plan is the subset of a Terraform plan document that Terraview consumes.
type Plan struct {
	FormatVersion    string                  `json:"format_version"`
	TerraformVersion string                  `json:"terraform_version"`
	PlannedValues    PlannedValues           `json:"planned_values"`
	ResourceChanges  []ResourceChange        `json:"resource_changes"`
	OutputChanges    map[string]OutputChange `json:"output_changes"`
	Diagnostics      []Diagnostic            `json:"diagnostics"`
	Configuration    Configuration           `json:"configuration"`
}

// PlannedValues describes the planned values tree for a Terraform plan.
type PlannedValues struct {
	RootModule PlannedModule `json:"root_module"`
}

// PlannedModule describes a module and its nested child modules.
type PlannedModule struct {
	Address      string            `json:"address,omitempty"`
	Resources    []PlannedResource `json:"resources,omitempty"`
	ChildModules []PlannedModule   `json:"child_modules,omitempty"`
}

// PlannedResource describes a resource instance in planned values.
type PlannedResource struct {
	Address      string         `json:"address"`
	Mode         string         `json:"mode"`
	Type         string         `json:"type"`
	Name         string         `json:"name"`
	ProviderName string         `json:"provider_name"`
	Values       map[string]any `json:"values"`
}

// ResourceChange describes a planned change to a resource instance.
type ResourceChange struct {
	Address       string `json:"address"`
	ModuleAddress string `json:"module_address,omitempty"`
	Mode          string `json:"mode"`
	Type          string `json:"type"`
	Name          string `json:"name"`
	ProviderName  string `json:"provider_name"`
	Change        Change `json:"change"`
}

// Diagnostic describes a diagnostic returned in the plan output.
type Diagnostic struct {
	Severity string         `json:"severity"`
	Summary  string         `json:"summary"`
	Detail   string         `json:"detail"`
	Range    map[string]any `json:"range,omitempty"`
	Snippet  map[string]any `json:"snippet,omitempty"`
}

// Change describes the before and after values for a planned change.
type Change struct {
	Actions      []string       `json:"actions"`
	Before       map[string]any `json:"before"`
	After        map[string]any `json:"after"`
	AfterUnknown map[string]any `json:"after_unknown"`
	ReplacePaths [][]string     `json:"replace_paths,omitempty"`
}

// OutputChange describes a planned change to a root module output.
type OutputChange struct {
	Actions      []string `json:"actions"`
	Before       any      `json:"before"`
	After        any      `json:"after"`
	AfterUnknown any      `json:"after_unknown"`
}

// Configuration describes the parsed configuration included in the plan.
type Configuration struct {
	ProviderConfig map[string]ProviderConfig `json:"provider_config"`
	RootModule     ConfigModule              `json:"root_module"`
}

// ProviderConfig describes a provider configuration referenced by resources.
type ProviderConfig struct {
	Name              string `json:"name"`
	FullName          string `json:"full_name"`
	VersionConstraint string `json:"version_constraint,omitempty"`
}

// ConfigModule describes a module in the configuration tree.
type ConfigModule struct {
	Resources   []ConfigResource            `json:"resources,omitempty"`
	ModuleCalls map[string]ConfigModuleCall `json:"module_calls,omitempty"`
}

// ConfigResource describes a configured resource in a module.
type ConfigResource struct {
	Address           string `json:"address"`
	Mode              string `json:"mode"`
	Type              string `json:"type"`
	Name              string `json:"name"`
	ProviderConfigKey string `json:"provider_config_key"`
}

// ConfigModuleCall describes a module block and its child configuration.
type ConfigModuleCall struct {
	Source string       `json:"source"`
	Module ConfigModule `json:"module"`
}
