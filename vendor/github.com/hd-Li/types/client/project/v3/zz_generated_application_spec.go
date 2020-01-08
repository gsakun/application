package client

const (
	ApplicationSpecType            = "applicationSpec"
	ApplicationSpecFieldComponents = "components"
)

type ApplicationSpec struct {
	Components []Component `json:"components,omitempty" yaml:"components,omitempty"`
}
