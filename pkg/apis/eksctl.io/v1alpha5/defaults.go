package v1alpha5

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SetClusterConfigDefaults will set defaults for a given cluster
func SetClusterConfigDefaults(cfg *ClusterConfig) {
	if cfg.IAM == nil {
		cfg.IAM = &ClusterIAM{}
	}

	if cfg.IAM.WithOIDC == nil {
		cfg.IAM.WithOIDC = Disabled()
	}

	for _, sa := range cfg.IAM.ServiceAccounts {
		if sa.Namespace == "" {
			sa.Namespace = metav1.NamespaceDefault
		}
	}

	if cfg.HasClusterCloudWatchLogging() && len(cfg.CloudWatch.ClusterLogging.EnableTypes) == 1 {
		switch cfg.CloudWatch.ClusterLogging.EnableTypes[0] {
		case "all", "*":
			cfg.CloudWatch.ClusterLogging.EnableTypes = SupportedCloudWatchClusterLogTypes()
		}
	}
}

// SetNodeGroupDefaults will set defaults for a given nodegroup
func SetNodeGroupDefaults(ng *NodeGroup, meta *ClusterMeta) {
	if ng.InstanceType == "" {
		if HasMixedInstances(ng) {
			ng.InstanceType = "mixed"
		} else {
			ng.InstanceType = DefaultNodeType
		}
	}
	if ng.AMIFamily == "" {
		ng.AMIFamily = DefaultNodeImageFamily
	}
	if ng.AMI == "" {
		ng.AMI = "static"
	}

	if ng.SecurityGroups == nil {
		ng.SecurityGroups = &NodeGroupSGs{
			AttachIDs: []string{},
		}
	}
	if ng.SecurityGroups.WithLocal == nil {
		ng.SecurityGroups.WithLocal = Enabled()
	}
	if ng.SecurityGroups.WithShared == nil {
		ng.SecurityGroups.WithShared = Enabled()
	}

	if ng.SSH == nil {
		ng.SSH = &NodeGroupSSH{
			Allow: Disabled(),
		}
	}

	setSSHDefaults(ng.SSH)

	if !IsSetAndNonEmptyString(ng.VolumeType) {
		ng.VolumeType = &DefaultNodeVolumeType
	}

	if ng.IAM == nil {
		ng.IAM = &NodeGroupIAM{}
	}

	setIAMDefaults(ng.IAM)

	if ng.Labels == nil {
		ng.Labels = make(map[string]string)
	}
	setDefaultNodeLabels(ng.Labels, meta.Name, ng.Name)
}

// SetManagedNodeGroupDefaults sets default values for a ManagedNodeGroup
func SetManagedNodeGroupDefaults(ng *ManagedNodeGroup, meta *ClusterMeta) {
	if ng.AMIFamily == "" {
		ng.AMIFamily = NodeImageFamilyAmazonLinux2
	}
	if ng.InstanceType == "" {
		ng.InstanceType = DefaultNodeType
	}
	if ng.ScalingConfig == nil {
		ng.ScalingConfig = &ScalingConfig{}
	}
	if ng.SSH == nil {
		ng.SSH = &NodeGroupSSH{
			Allow: Disabled(),
		}
	}
	setSSHDefaults(ng.SSH)

	if ng.IAM == nil {
		ng.IAM = &NodeGroupIAM{}
	}
	setIAMDefaults(ng.IAM)

	if ng.Labels == nil {
		ng.Labels = make(map[string]string)
	}
	setDefaultNodeLabels(ng.Labels, meta.Name, ng.Name)

	if ng.Tags == nil {
		ng.Tags = make(map[string]string)
	}
	ng.Tags[NodeGroupNameTag] = ng.Name
	ng.Tags[NodeGroupTypeTag] = string(NodeGroupTypeManaged)
}

func setIAMDefaults(iamConfig *NodeGroupIAM) {
	if iamConfig.WithAddonPolicies.ImageBuilder == nil {
		iamConfig.WithAddonPolicies.ImageBuilder = Disabled()
	}
	if iamConfig.WithAddonPolicies.AutoScaler == nil {
		iamConfig.WithAddonPolicies.AutoScaler = Disabled()
	}
	if iamConfig.WithAddonPolicies.ExternalDNS == nil {
		iamConfig.WithAddonPolicies.ExternalDNS = Disabled()
	}
	if iamConfig.WithAddonPolicies.CertManager == nil {
		iamConfig.WithAddonPolicies.CertManager = Disabled()
	}
	if iamConfig.WithAddonPolicies.ALBIngress == nil {
		iamConfig.WithAddonPolicies.ALBIngress = Disabled()
	}
	if iamConfig.WithAddonPolicies.XRay == nil {
		iamConfig.WithAddonPolicies.XRay = Disabled()
	}
	if iamConfig.WithAddonPolicies.CloudWatch == nil {
		iamConfig.WithAddonPolicies.CloudWatch = Disabled()
	}
	if iamConfig.WithAddonPolicies.EBS == nil {
		iamConfig.WithAddonPolicies.EBS = Disabled()
	}
	if iamConfig.WithAddonPolicies.FSX == nil {
		iamConfig.WithAddonPolicies.FSX = Disabled()
	}
	if iamConfig.WithAddonPolicies.EFS == nil {
		iamConfig.WithAddonPolicies.EFS = Disabled()
	}
}

func setSSHDefaults(sshConfig *NodeGroupSSH) {
	numSSHFlagsEnabled := countEnabledFields(
		sshConfig.PublicKeyName,
		sshConfig.PublicKeyPath,
		sshConfig.PublicKey)

	if numSSHFlagsEnabled == 0 {
		if IsEnabled(sshConfig.Allow) {
			sshConfig.PublicKeyPath = &DefaultNodeSSHPublicKeyPath
		} else {
			sshConfig.Allow = Disabled()
		}
	} else if !IsDisabled(sshConfig.Allow) {
		// Enable SSH if not explicitly disabled when passing an SSH key
		sshConfig.Allow = Enabled()
	}

}

func setDefaultNodeLabels(labels map[string]string, clusterName, nodeGroupName string) {
	labels[ClusterNameLabel] = clusterName
	labels[NodeGroupNameLabel] = nodeGroupName
}

// DefaultClusterNAT will set the default value for Cluster NAT mode
func DefaultClusterNAT() *ClusterNAT {
	single := ClusterSingleNAT
	return &ClusterNAT{
		Gateway: &single,
	}
}

// ClusterEndpointAccessDefaults returns a ClusterEndpoints pointer with default values set.
func ClusterEndpointAccessDefaults() *ClusterEndpoints {
	return &ClusterEndpoints{
		PrivateAccess: Disabled(),
		PublicAccess:  Enabled(),
	}
}
