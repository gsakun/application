package client

const (
	IngressLBType                = "ingressLB"
	IngressLBFieldConsistentType = "consistentType"
	IngressLBFieldLBType         = "lbType"
)

type IngressLB struct {
	ConsistentType string `json:"consistentType,omitempty" yaml:"consistentType,omitempty"`
	LBType         string `json:"lbType,omitempty" yaml:"lbType,omitempty"`
}
