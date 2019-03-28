package create

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ami"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func configureNodeGroups(ngFilter *cmdutils.NodeGroupFilter, nodeGroups []*api.NodeGroup, cmd *cobra.Command) error {
	return ngFilter.CheckEachNodeGroup(nodeGroups, func(i int, ng *api.NodeGroup) error {
		if ng.AllowSSH && ng.SSHPublicKeyPath == "" {
			return fmt.Errorf("--ssh-public-key must be non-empty string")
		}

		if cmd.Flag("ssh-public-key").Changed {
			ng.AllowSSH = true
		}

		// generate nodegroup name or use flag
		ng.Name = NodeGroupName(ng.Name, "")

		return nil
	})
}

func setNodeGroupDefaults(ngFilter *cmdutils.NodeGroupFilter, nodeGroups []*api.NodeGroup) error {
	return ngFilter.CheckEachNodeGroup(nodeGroups, func(i int, ng *api.NodeGroup) error {

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
	})
}
