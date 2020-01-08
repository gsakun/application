package client

const (
	ComponentTraitsForOptType               = "componentTraitsForOpt"
	ComponentTraitsForOptFieldEject         = "eject"
	ComponentTraitsForOptFieldIngress       = "ingress"
	ComponentTraitsForOptFieldManualScaler  = "manualScaler"
	ComponentTraitsForOptFieldRateLimit     = "rateLimit"
	ComponentTraitsForOptFieldVolumeMounter = "volumeMounter"
	ComponentTraitsForOptFieldWhiteList     = "whiteList"
)

type ComponentTraitsForOpt struct {
	Eject         []string       `json:"eject,omitempty" yaml:"eject,omitempty"`
	Ingress       *AppIngress    `json:"ingress,omitempty" yaml:"ingress,omitempty"`
	ManualScaler  *ManualScaler  `json:"manualScaler,omitempty" yaml:"manualScaler,omitempty"`
	RateLimit     *RateLimit     `json:"rateLimit,omitempty" yaml:"rateLimit,omitempty"`
	VolumeMounter *VolumeMounter `json:"volumeMounter,omitempty" yaml:"volumeMounter,omitempty"`
	WhiteList     *WhiteList     `json:"whiteList,omitempty" yaml:"whiteList,omitempty"`
}
