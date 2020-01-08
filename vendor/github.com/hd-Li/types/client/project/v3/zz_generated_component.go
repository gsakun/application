package client

const (
	ComponentType                  = "component"
	ComponentFieldArch             = "arch"
	ComponentFieldContainers       = "containers"
	ComponentFieldDevTraits        = "devTraits"
	ComponentFieldName             = "name"
	ComponentFieldOptTraits        = "optTraits"
	ComponentFieldOsType           = "osType"
	ComponentFieldParameters       = "parameters"
	ComponentFieldVersion          = "version"
	ComponentFieldWorkloadSettings = "workloadSetings"
	ComponentFieldWorkloadType     = "workloadType"
)

type Component struct {
	Arch             string                 `json:"arch,omitempty" yaml:"arch,omitempty"`
	Containers       []ComponentContainer   `json:"containers,omitempty" yaml:"containers,omitempty"`
	DevTraits        *ComponentTraitsForDev `json:"devTraits,omitempty" yaml:"devTraits,omitempty"`
	Name             string                 `json:"name,omitempty" yaml:"name,omitempty"`
	OptTraits        *ComponentTraitsForOpt `json:"optTraits,omitempty" yaml:"optTraits,omitempty"`
	OsType           string                 `json:"osType,omitempty" yaml:"osType,omitempty"`
	Parameters       []Parameter            `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Version          string                 `json:"version,omitempty" yaml:"version,omitempty"`
	WorkloadSettings []WorkloadSetting      `json:"workloadSetings,omitempty" yaml:"workloadSetings,omitempty"`
	WorkloadType     string                 `json:"workloadType,omitempty" yaml:"workloadType,omitempty"`
}
