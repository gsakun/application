package client

const (
	RoleRefType      = "roleRef"
	RoleRefFieldKind = "kind"
	RoleRefFieldName = "name"
)

type RoleRef struct {
	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}
