package client

const (
	ServiceRoleSpecType       = "serviceRoleSpec"
	ServiceRoleSpecFieldRules = "rules"
)

type ServiceRoleSpec struct {
	Rules []AccessRule `json:"rules,omitempty" yaml:"rules,omitempty"`
}
