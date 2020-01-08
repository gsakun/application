package client

const (
	AccessRule_ConstraintType        = "accessRule_Constraint"
	AccessRule_ConstraintFieldKey    = "key"
	AccessRule_ConstraintFieldValues = "values"
)

type AccessRule_Constraint struct {
	Key    string   `json:"key,omitempty" yaml:"key,omitempty"`
	Values []string `json:"values,omitempty" yaml:"values,omitempty"`
}
