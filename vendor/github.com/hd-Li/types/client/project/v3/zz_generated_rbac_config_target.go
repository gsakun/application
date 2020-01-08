package client

const (
	RbacConfigTargetType            = "rbacConfigTarget"
	RbacConfigTargetFieldNamespaces = "namespaces"
	RbacConfigTargetFieldServices   = "services"
)

type RbacConfigTarget struct {
	Namespaces []string `json:"namespaces,omitempty" yaml:"namespaces,omitempty"`
	Services   []string `json:"services,omitempty" yaml:"services,omitempty"`
}
