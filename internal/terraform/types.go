package terraform

type Plan struct {
	FormatVersion    string                  `json:"format_version"`
	TerraformVersion string                  `json:"terraform_version"`
	PlannedValues    PlannedValues           `json:"planned_values"`
	ResourceChanges  []ResourceChange        `json:"resource_changes"`
	OutputChanges    map[string]OutputChange `json:"output_changes"`
	Configuration    Configuration           `json:"configuration"`
}

type PlannedValues struct {
	RootModule PlannedModule `json:"root_module"`
}

type PlannedModule struct {
	Address      string            `json:"address,omitempty"`
	Resources    []PlannedResource `json:"resources,omitempty"`
	ChildModules []PlannedModule   `json:"child_modules,omitempty"`
}

type PlannedResource struct {
	Address      string         `json:"address"`
	Mode         string         `json:"mode"`
	Type         string         `json:"type"`
	Name         string         `json:"name"`
	ProviderName string         `json:"provider_name"`
	Values       map[string]any `json:"values"`
}

type ResourceChange struct {
	Address       string `json:"address"`
	ModuleAddress string `json:"module_address,omitempty"`
	Mode          string `json:"mode"`
	Type          string `json:"type"`
	Name          string `json:"name"`
	ProviderName  string `json:"provider_name"`
	Change        Change `json:"change"`
}

type Change struct {
	Actions      []string       `json:"actions"`
	Before       map[string]any `json:"before"`
	After        map[string]any `json:"after"`
	AfterUnknown map[string]any `json:"after_unknown"`
	ReplacePaths [][]string     `json:"replace_paths,omitempty"`
}

type OutputChange struct {
	Actions      []string `json:"actions"`
	Before       any      `json:"before"`
	After        any      `json:"after"`
	AfterUnknown any      `json:"after_unknown"`
}

type Configuration struct {
	ProviderConfig map[string]ProviderConfig `json:"provider_config"`
	RootModule     ConfigModule              `json:"root_module"`
}

type ProviderConfig struct {
	Name              string `json:"name"`
	FullName          string `json:"full_name"`
	VersionConstraint string `json:"version_constraint,omitempty"`
}

type ConfigModule struct {
	Resources   []ConfigResource            `json:"resources,omitempty"`
	ModuleCalls map[string]ConfigModuleCall `json:"module_calls,omitempty"`
}

type ConfigResource struct {
	Address           string `json:"address"`
	Mode              string `json:"mode"`
	Type              string `json:"type"`
	Name              string `json:"name"`
	ProviderConfigKey string `json:"provider_config_key"`
}

type ConfigModuleCall struct {
	Source string       `json:"source"`
	Module ConfigModule `json:"module"`
}
