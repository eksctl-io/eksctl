package v1alpha5

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/aws/aws-sdk-go-v2/aws"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/weaveworks/eksctl/pkg/utils"
)

const (
	IAMPolicyAmazonEKSCNIPolicy = "AmazonEKS_CNI_Policy"
)

const (
	// Data volume, used by kubelet
	bottlerocketDataDisk = "/dev/xvdb"
	// OS volume
	bottlerocketOSDisk = "/dev/xvda"
)

var (
	AWSNodeMeta = ClusterIAMMeta{
		Name:      "aws-node",
		Namespace: "kube-system",
	}
)

// SetClusterConfigDefaults will set defaults for a given cluster
func SetClusterConfigDefaults(cfg *ClusterConfig) {
	if cfg.IAM == nil {
		cfg.IAM = &ClusterIAM{}
	}

	if cfg.IAM.WithOIDC == nil {
		cfg.IAM.WithOIDC = Disabled()
	}

	if cfg.IAM.VPCResourceControllerPolicy == nil {
		cfg.IAM.VPCResourceControllerPolicy = Enabled()
	}

	for _, sa := range cfg.IAM.ServiceAccounts {
		if sa.Namespace == "" {
			sa.Namespace = metav1.NamespaceDefault
		}
	}

	if cfg.HasClusterCloudWatchLogging() && cfg.ContainsWildcardCloudWatchLogging() {
		cfg.CloudWatch.ClusterLogging.EnableTypes = SupportedCloudWatchClusterLogTypes()
	}

	if cfg.AccessConfig == nil {
		cfg.AccessConfig = &AccessConfig{
			AuthenticationMode: ekstypes.AuthenticationModeApiAndConfigMap,
		}
	} else if cfg.AccessConfig.AuthenticationMode == "" {
		cfg.AccessConfig.AuthenticationMode = ekstypes.AuthenticationModeApiAndConfigMap
	}

	if cfg.PrivateCluster == nil {
		cfg.PrivateCluster = &PrivateCluster{}
	}

	if cfg.VPC != nil && cfg.VPC.ManageSharedNodeSecurityGroupRules == nil {
		cfg.VPC.ManageSharedNodeSecurityGroupRules = Enabled()
	}

	if cfg.Karpenter != nil && cfg.Karpenter.CreateServiceAccount == nil {
		cfg.Karpenter.CreateServiceAccount = Disabled()
	}
}

// IAMServiceAccountsWithImplicitServiceAccounts adds implicitly created
// IAM SAs that need to be explicitly deleted.
func IAMServiceAccountsWithImplicitServiceAccounts(cfg *ClusterConfig) []*ClusterIAMServiceAccount {
	serviceAccounts := cfg.IAM.ServiceAccounts
	if IsEnabled(cfg.IAM.WithOIDC) && !vpcCNIAddonSpecified(cfg) {
		var found bool
		for _, sa := range cfg.IAM.ServiceAccounts {
			found = found || (sa.Name == AWSNodeMeta.Name && sa.Namespace == AWSNodeMeta.Namespace)
		}
		if !found {
			awsNode := ClusterIAMServiceAccount{
				ClusterIAMMeta: AWSNodeMeta,
				AttachPolicyARNs: []string{
					fmt.Sprintf("arn:%s:iam::aws:policy/%s", Partitions.ForRegion(cfg.Metadata.Region), IAMPolicyAmazonEKSCNIPolicy),
				},
			}
			serviceAccounts = append(serviceAccounts, &awsNode)
		}
	}
	return serviceAccounts
}

func vpcCNIAddonSpecified(cfg *ClusterConfig) bool {
	for _, a := range cfg.Addons {
		if strings.ToLower(a.Name) == "vpc-cni" {
			return true
		}
	}
	return false
}

// SetNodeGroupDefaults will set defaults for a given nodegroup
func SetNodeGroupDefaults(ng *NodeGroup, meta *ClusterMeta, controlPlaneOnOutposts bool) {
	setNodeGroupBaseDefaults(ng.NodeGroupBase, meta)

	if ng.AMIFamily == "" {
		ng.AMIFamily = DefaultNodeImageFamily
	}

	setVolumeDefaults(ng.NodeGroupBase, controlPlaneOnOutposts, nil)
	setDefaultsForAdditionalVolumes(ng.NodeGroupBase, controlPlaneOnOutposts)

	if ng.SecurityGroups.WithLocal == nil {
		ng.SecurityGroups.WithLocal = Enabled()
	}
	if ng.SecurityGroups.WithShared == nil {
		ng.SecurityGroups.WithShared = Enabled()
	}

	setContainerRuntimeDefault(ng, meta.Version)
}

// SetManagedNodeGroupDefaults sets default values for a ManagedNodeGroup
func SetManagedNodeGroupDefaults(ng *ManagedNodeGroup, meta *ClusterMeta, controlPlaneOnOutposts bool) {
	setNodeGroupBaseDefaults(ng.NodeGroupBase, meta)

	// When using custom AMIs, we want the user to explicitly specify AMI family.
	// Thus, we only setup default AMI family when no custom AMI is being used.
	if ng.AMIFamily == "" && ng.AMI == "" {
		ng.AMIFamily = NodeImageFamilyAmazonLinux2
	}

	if ng.Tags == nil {
		ng.Tags = make(map[string]string)
	}
	ng.Tags[NodeGroupNameTag] = ng.Name
	ng.Tags[NodeGroupTypeTag] = string(NodeGroupTypeManaged)

	setVolumeDefaults(ng.NodeGroupBase, controlPlaneOnOutposts, ng.LaunchTemplate)
	setDefaultsForAdditionalVolumes(ng.NodeGroupBase, controlPlaneOnOutposts)
}

func setNodeGroupBaseDefaults(ng *NodeGroupBase, meta *ClusterMeta) {
	if ng.ScalingConfig == nil {
		ng.ScalingConfig = &ScalingConfig{}
	}
	if ng.SSH == nil {
		ng.SSH = &NodeGroupSSH{
			Allow: Disabled(),
		}
	}
	setSSHDefaults(ng.SSH)

	if ng.SecurityGroups == nil {
		ng.SecurityGroups = &NodeGroupSGs{}
	}

	if ng.IAM == nil {
		ng.IAM = &NodeGroupIAM{}
	}
	setIAMDefaults(ng.IAM)

	if ng.Labels == nil {
		ng.Labels = make(map[string]string)
	}
	setDefaultNodeLabels(ng.Labels, meta.Name, ng.Name)

	if ng.DisableIMDSv1 == nil {
		ng.DisableIMDSv1 = Enabled()
	}
	if ng.DisablePodIMDS == nil {
		ng.DisablePodIMDS = Disabled()
	}
	if ng.InstanceSelector == nil {
		ng.InstanceSelector = &InstanceSelector{}
	}
	normalizeAMIFamily(ng)
	if ng.AMIFamily == NodeImageFamilyBottlerocket {
		setBottlerocketNodeGroupDefaults(ng)
	}
}

func setVolumeDefaults(ng *NodeGroupBase, controlPlaneOnOutposts bool, template *LaunchTemplate) {
	if ng.VolumeType == nil {
		ng.VolumeType = aws.String(getDefaultVolumeType(controlPlaneOnOutposts || ng.OutpostARN != ""))
	}
	if ng.VolumeSize == nil && template == nil {
		ng.VolumeSize = &DefaultNodeVolumeSize
	}

	switch *ng.VolumeType {
	case NodeVolumeTypeGP3:
		if ng.VolumeIOPS == nil {
			ng.VolumeIOPS = aws.Int(DefaultNodeVolumeGP3IOPS)
		}
		if ng.VolumeThroughput == nil {
			ng.VolumeThroughput = aws.Int(DefaultNodeVolumeThroughput)
		}
	case NodeVolumeTypeIO1:
		if ng.VolumeIOPS == nil {
			ng.VolumeIOPS = aws.Int(DefaultNodeVolumeIO1IOPS)
		}
	}

	if ng.AMIFamily == NodeImageFamilyBottlerocket && !IsSetAndNonEmptyString(ng.VolumeName) {
		ng.AdditionalEncryptedVolume = bottlerocketOSDisk
		ng.VolumeName = aws.String(bottlerocketDataDisk)
	}
}

func setDefaultsForAdditionalVolumes(ng *NodeGroupBase, controlPlaneOnOutposts bool) {
	for i, av := range ng.AdditionalVolumes {
		if av.VolumeType == nil {
			ng.AdditionalVolumes[i].VolumeType = aws.String(getDefaultVolumeType(controlPlaneOnOutposts))
		}
		if av.VolumeSize == nil {
			ng.AdditionalVolumes[i].VolumeSize = &DefaultNodeVolumeSize
		}
		if *av.VolumeType == NodeVolumeTypeGP3 {
			if av.VolumeIOPS == nil {
				ng.AdditionalVolumes[i].VolumeIOPS = aws.Int(DefaultNodeVolumeGP3IOPS)
			}
			if av.VolumeThroughput == nil {
				ng.AdditionalVolumes[i].VolumeThroughput = aws.Int(DefaultNodeVolumeThroughput)
			}
		}
		if *av.VolumeType == NodeVolumeTypeIO1 && av.VolumeIOPS == nil {
			ng.AdditionalVolumes[i].VolumeIOPS = aws.Int(DefaultNodeVolumeIO1IOPS)
		}
	}
}

func getDefaultVolumeType(nodeGroupOnOutposts bool) string {
	if nodeGroupOnOutposts {
		return NodeVolumeTypeGP2
	}
	return DefaultNodeVolumeType
}

func setContainerRuntimeDefault(ng *NodeGroup, clusterVersion string) {
	if ng.ContainerRuntime != nil {
		return
	}

	// since clusterVersion is standardised beforehand, we can safely ignore the error
	isDockershimDeprecated, _ := utils.IsMinVersion(DockershimDeprecationVersion, clusterVersion)

	if isDockershimDeprecated {
		ng.ContainerRuntime = aws.String(ContainerRuntimeContainerD)
	} else {
		ng.ContainerRuntime = aws.String(ContainerRuntimeDockerD)
		if IsWindowsImage(ng.AMIFamily) {
			ng.ContainerRuntime = aws.String(ContainerRuntimeDockerForWindows)
		}
	}
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
	if iamConfig.WithAddonPolicies.AWSLoadBalancerController == nil {
		iamConfig.WithAddonPolicies.AWSLoadBalancerController = Disabled()
	}
	if iamConfig.WithAddonPolicies.DeprecatedALBIngress == nil {
		iamConfig.WithAddonPolicies.DeprecatedALBIngress = Disabled()
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

func setBottlerocketNodeGroupDefaults(ng *NodeGroupBase) {
	// Initialize config object if not present.
	if ng.Bottlerocket == nil {
		ng.Bottlerocket = &NodeGroupBottlerocket{}
	}
	if ng.Bottlerocket.Settings == nil {
		ng.Bottlerocket.Settings = &InlineDocument{}
	}

	// Use the SSH settings if the user hasn't explicitly configured the Admin
	// Container. If SSH was enabled, the user will be able to ssh into the
	// Bottlerocket node via the admin container.
	if ng.Bottlerocket.EnableAdminContainer == nil && ng.SSH != nil && IsEnabled(ng.SSH.Allow) {
		ng.Bottlerocket.EnableAdminContainer = Enabled()
	}
}

// DefaultClusterNAT will set the default value for Cluster NAT mode
func DefaultClusterNAT() *ClusterNAT {
	def := ClusterNATDefault
	return &ClusterNAT{
		Gateway: &def,
	}
}

// SetClusterEndpointAccessDefaults sets the default values for cluster endpoint access
func SetClusterEndpointAccessDefaults(vpc *ClusterVPC) {
	endpointAccess := ClusterEndpointAccessDefaults()
	if vpc.ClusterEndpoints == nil {
		vpc.ClusterEndpoints = endpointAccess
		return
	}

	if vpc.ClusterEndpoints.PublicAccess == nil {
		vpc.ClusterEndpoints.PublicAccess = endpointAccess.PublicAccess
	}

	if vpc.ClusterEndpoints.PrivateAccess == nil {
		vpc.ClusterEndpoints.PrivateAccess = endpointAccess.PrivateAccess
	}
}

// ClusterEndpointAccessDefaults returns a ClusterEndpoints pointer with default values set.
func ClusterEndpointAccessDefaults() *ClusterEndpoints {
	return &ClusterEndpoints{
		PrivateAccess: Disabled(),
		PublicAccess:  Enabled(),
	}
}

// SetDefaultFargateProfile configures this ClusterConfig to have a single
// Fargate profile called "default", with two selectors matching respectively
// the "default" and "kube-system" Kubernetes namespaces.
func (c *ClusterConfig) SetDefaultFargateProfile() {
	c.FargateProfiles = []*FargateProfile{
		{
			Name: "fp-default",
			Selectors: []FargateProfileSelector{
				{Namespace: "default"},
				{Namespace: "kube-system"},
			},
		},
	}
}
