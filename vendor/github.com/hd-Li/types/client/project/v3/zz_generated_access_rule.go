package client

const (
	AccessRuleType             = "accessRule"
	AccessRuleFieldConstraints = "constraints"
	AccessRuleFieldHosts       = "hosts"
	AccessRuleFieldMethods     = "methods"
	AccessRuleFieldNotHosts    = "not_hosts"
	AccessRuleFieldNotMethods  = "not_methods"
	AccessRuleFieldNotPaths    = "not_paths"
	AccessRuleFieldNotPorts    = "not_ports"
	AccessRuleFieldPaths       = "paths"
	AccessRuleFieldPorts       = "ports"
	AccessRuleFieldServices    = "services"
)

type AccessRule struct {
	Constraints []AccessRule_Constraint `json:"constraints,omitempty" yaml:"constraints,omitempty"`
	Hosts       []string                `json:"hosts,omitempty" yaml:"hosts,omitempty"`
	Methods     []string                `json:"methods,omitempty" yaml:"methods,omitempty"`
	NotHosts    []string                `json:"not_hosts,omitempty" yaml:"not_hosts,omitempty"`
	NotMethods  []string                `json:"not_methods,omitempty" yaml:"not_methods,omitempty"`
	NotPaths    []string                `json:"not_paths,omitempty" yaml:"not_paths,omitempty"`
	NotPorts    []int64                 `json:"not_ports,omitempty" yaml:"not_ports,omitempty"`
	Paths       []string                `json:"paths,omitempty" yaml:"paths,omitempty"`
	Ports       []int64                 `json:"ports,omitempty" yaml:"ports,omitempty"`
	Services    []string                `json:"services,omitempty" yaml:"services,omitempty"`
}
