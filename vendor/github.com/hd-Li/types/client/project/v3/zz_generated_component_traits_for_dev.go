package client

const (
	ComponentTraitsForDevType                 = "componentTraitsForDev"
	ComponentTraitsForDevFieldImagePullConfig = "imagePullConfig"
	ComponentTraitsForDevFieldIngressLB       = "ingressLB"
	ComponentTraitsForDevFieldStaticIP        = "staticIP"
)

type ComponentTraitsForDev struct {
	ImagePullConfig *ImagePullConfig `json:"imagePullConfig,omitempty" yaml:"imagePullConfig,omitempty"`
	IngressLB       *IngressLB       `json:"ingressLB,omitempty" yaml:"ingressLB,omitempty"`
	StaticIP        bool             `json:"staticIP,omitempty" yaml:"staticIP,omitempty"`
}
