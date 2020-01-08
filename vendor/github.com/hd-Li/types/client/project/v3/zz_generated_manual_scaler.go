package client

const (
	ManualScalerType          = "manualScaler"
	ManualScalerFieldReplicas = "replicas"
)

type ManualScaler struct {
	Replicas int64 `json:"replicas,omitempty" yaml:"replicas,omitempty"`
}
