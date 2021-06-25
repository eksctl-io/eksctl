package v1alpha5

// WellKnownPolicies for attaching common IAM policies
type WellKnownPolicies struct {
	// ImageBuilder allows for full ECR (Elastic Container Registry) access.
	ImageBuilder bool `json:"imageBuilder,inline"`
	// AutoScaler adds policies for cluster-autoscaler. See [autoscaler AWS
	// docs](https://docs.aws.amazon.com/eks/latest/userguide/cluster-autoscaler.html).
	AutoScaler bool `json:"autoScaler,inline"`
	// AWSLoadBalancerController adds policies for using the
	// aws-load-balancer-controller. See [Load Balancer
	// docs](https://docs.aws.amazon.com/eks/latest/userguide/aws-load-balancer-controller.html).
	AWSLoadBalancerController bool `json:"awsLoadBalancerController,inline"`
	// ExternalDNS adds external-dns policies for Amazon Route 53.
	// See [external-dns
	// docs](https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/aws.md).
	ExternalDNS bool `json:"externalDNS,inline"`
	// CertManager adds cert-manager policies. See [cert-manager
	// docs](https://cert-manager.io/docs/configuration/acme/dns01/route53).
	CertManager bool `json:"certManager,inline"`
	// EBSCSIController adds policies for using the
	// ebs-csi-controller. See [aws-ebs-csi-driver
	// docs](https://github.com/kubernetes-sigs/aws-ebs-csi-driver#set-up-driver-permission).
	EBSCSIController bool `json:"ebsCSIController,inline"`
	// EFSCSIController adds policies for using the
	// efs-csi-controller. See [aws-efs-csi-driver
	// docs](https://aws.amazon.com/blogs/containers/introducing-efs-csi-dynamic-provisioning).
	EFSCSIController bool `json:"efsCSIController,inline"`
}

func (p *WellKnownPolicies) HasPolicy() bool {
	return p.ImageBuilder || p.AutoScaler || p.AWSLoadBalancerController || p.ExternalDNS || p.CertManager || p.EBSCSIController || p.EFSCSIController
}
