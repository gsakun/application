package client

const (
	RbacConfigSpecType                 = "rbacConfigSpec"
	RbacConfigSpecFieldEnforcementMode = "enforcementMode"
	RbacConfigSpecFieldExclusion       = "exclusion"
	RbacConfigSpecFieldInclusion       = "inclusion"
	RbacConfigSpecFieldMode            = "mode"
)

type RbacConfigSpec struct {
	EnforcementMode int64             `json:"enforcementMode,omitempty" yaml:"enforcementMode,omitempty"`
	Exclusion       *RbacConfigTarget `json:"exclusion,omitempty" yaml:"exclusion,omitempty"`
	Inclusion       *RbacConfigTarget `json:"inclusion,omitempty" yaml:"inclusion,omitempty"`
	Mode            int64             `json:"mode,omitempty" yaml:"mode,omitempty"`
}
