package client

const (
	SubjectType               = "subject"
	SubjectFieldGroup         = "group"
	SubjectFieldGroups        = "groups"
	SubjectFieldIps           = "ips"
	SubjectFieldNames         = "names"
	SubjectFieldNamespaces    = "namespaces"
	SubjectFieldNotGroups     = "not_groups"
	SubjectFieldNotIps        = "not_ips"
	SubjectFieldNotNames      = "not_names"
	SubjectFieldNotNamespaces = "not_namespaces"
	SubjectFieldProperties    = "properties"
	SubjectFieldUser          = "user"
)

type Subject struct {
	Group         string            `json:"group,omitempty" yaml:"group,omitempty"`
	Groups        []string          `json:"groups,omitempty" yaml:"groups,omitempty"`
	Ips           []string          `json:"ips,omitempty" yaml:"ips,omitempty"`
	Names         []string          `json:"names,omitempty" yaml:"names,omitempty"`
	Namespaces    []string          `json:"namespaces,omitempty" yaml:"namespaces,omitempty"`
	NotGroups     []string          `json:"not_groups,omitempty" yaml:"not_groups,omitempty"`
	NotIps        []string          `json:"not_ips,omitempty" yaml:"not_ips,omitempty"`
	NotNames      []string          `json:"not_names,omitempty" yaml:"not_names,omitempty"`
	NotNamespaces []string          `json:"not_namespaces,omitempty" yaml:"not_namespaces,omitempty"`
	Properties    map[string]string `json:"properties,omitempty" yaml:"properties,omitempty"`
	User          string            `json:"user,omitempty" yaml:"user,omitempty"`
}
