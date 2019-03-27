package create

import (
	"github.com/kris-nova/logger"

	"github.com/weaveworks/eksctl/pkg/ami"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
)

// When passing the --without-nodegroup option, don't create nodegroups
func skipNodeGroupsIfRequested(cfg *api.ClusterConfig) {
	if withoutNodeGroup {
		cfg.NodeGroups = nil
		logger.Warning("cluster will be created without an initial nodegroup")
	}
}

// SetNodeGroupDefaults will set defaults for a given nodegroup
func SetNodeGroupDefaults(i int, ng *api.NodeGroup) error {

	if err := api.ValidateNodeGroup(i, ng); err != nil {
		return err
	}

	// apply defaults
	if ng.InstanceType == "" {
		ng.InstanceType = api.DefaultNodeType
	}
	if ng.AMIFamily == "" {
		ng.AMIFamily = ami.ImageFamilyAmazonLinux2
	}
	if ng.AMI == "" {
		ng.AMI = ami.ResolverStatic
	}

	if ng.SecurityGroups == nil {
		ng.SecurityGroups = &api.NodeGroupSGs{
			AttachIDs: []string{},
		}
	}
	if ng.SecurityGroups.WithLocal == nil {
		ng.SecurityGroups.WithLocal = api.NewBoolTrue()
	}
	if ng.SecurityGroups.WithShared == nil {
		ng.SecurityGroups.WithShared = api.NewBoolTrue()
	}

	// Enable SSH when a key is provided
	if ng.SSHPublicKeyPath != "" {
		ng.AllowSSH = true
	}

	if ng.AllowSSH && ng.SSHPublicKeyPath == "" {
		ng.SSHPublicKeyPath = defaultSSHPublicKey
	}

	if ng.VolumeSize > 0 {
		if ng.VolumeType == "" {
			ng.VolumeType = api.DefaultNodeVolumeType
		}
	}

	if ng.IAM == nil {
		ng.IAM = &api.NodeGroupIAM{}
	}
	if ng.IAM.WithAddonPolicies.ImageBuilder == nil {
		ng.IAM.WithAddonPolicies.ImageBuilder = api.NewBoolFalse()
	}
	if ng.IAM.WithAddonPolicies.AutoScaler == nil {
		ng.IAM.WithAddonPolicies.AutoScaler = api.NewBoolFalse()
	}
	if ng.IAM.WithAddonPolicies.ExternalDNS == nil {
		ng.IAM.WithAddonPolicies.ExternalDNS = api.NewBoolFalse()
	}

	return nil
}
