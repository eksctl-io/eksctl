package create

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ami"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/utils"
)

const (
	randNodeGroupNameLength     = 8
	randNodeGroupNameComponents = "abcdef0123456789"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// NodeGroupName generates a name string when a and b are empty strings.
// If either a or b are non-empty, it returns whichever is non-empty.
// If neither a nor b are empty, it returns empty name, to indicate
// ambiguous usage.
// It uses a different naming scheme from ClusterName, so that users can
// easily distinguish a cluster name from nodegroup name.
func NodeGroupName(a, b string) string {
	return utils.UseNameOrGenerate(a, b, func() string {
		name := make([]byte, randNodeGroupNameLength)
		for i := 0; i < randNodeGroupNameLength; i++ {
			name[i] = randNodeGroupNameComponents[r.Intn(len(randNodeGroupNameComponents))]
		}
		return fmt.Sprintf("ng-%s", string(name))
	})
}

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
