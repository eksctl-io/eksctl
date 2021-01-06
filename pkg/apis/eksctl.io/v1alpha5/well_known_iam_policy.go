package v1alpha5

type WellKnownPolicies struct {
	ImageBuilder    bool `json:"imageBuilder,inline"`
	AutoScaler      bool `json:"autoScaler,inline"`
	AWSLoadBalancer bool `json:"awsLoadBalancer,inline"`
	ExternalDNS     bool `json:"externalDNS,inline"`
	CertManager     bool `json:"certManager,inline"`
}

func (p *WellKnownPolicies) HasPolicy() bool {
	return p.ImageBuilder || p.AutoScaler || p.AWSLoadBalancer || p.ExternalDNS || p.CertManager
}
