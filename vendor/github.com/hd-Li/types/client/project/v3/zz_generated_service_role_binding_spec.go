package client

const (
	ServiceRoleBindingSpecType          = "serviceRoleBindingSpec"
	ServiceRoleBindingSpecFieldMode     = "mode"
	ServiceRoleBindingSpecFieldRole     = "role"
	ServiceRoleBindingSpecFieldRoleRef  = "roleRef"
	ServiceRoleBindingSpecFieldSubjects = "subjects"
)

type ServiceRoleBindingSpec struct {
	Mode     int64     `json:"mode,omitempty" yaml:"mode,omitempty"`
	Role     string    `json:"role,omitempty" yaml:"role,omitempty"`
	RoleRef  *RoleRef  `json:"roleRef,omitempty" yaml:"roleRef,omitempty"`
	Subjects []Subject `json:"subjects,omitempty" yaml:"subjects,omitempty"`
}
