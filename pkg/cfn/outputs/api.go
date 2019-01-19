package outputs

import (
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha3"
)

// Stack output names
const (
	// outputs from cluster stack
	ClusterVPC            = "VPC"
	ClusterSecurityGroup  = "SecurityGroup"
	ClusterSubnets        = "Subnets"
	ClusterSubnetsPrivate = string(ClusterSubnets + api.SubnetTopologyPrivate)
	ClusterSubnetsPublic  = string(ClusterSubnets + api.SubnetTopologyPublic)

	ClusterCertificateAuthorityData = "CertificateAuthorityData"
	ClusterEndpoint                 = "Endpoint"
	ClusterARN                      = "ARN"
	ClusterStackName                = "ClusterStackName"
	ClusterSharedNodeSecurityGroup  = "SharedNodeSecurityGroup"

	// outputs from nodegroup stack
	NodeGroupInstanceRoleARN = "InstanceRoleARN"
	// outputs to indicate configuration attributes that may have critical effect
	// on critical effect on forward-compatibility with respect to overal functionality
	// and integrity, e.g. networking
	NodeGroupFeaturePrivateNetworking   = "FeaturePrivateNetworking"
	NodeGroupFeatureSharedSecurityGroup = "FeatureSharedSecurityGroup"
)
