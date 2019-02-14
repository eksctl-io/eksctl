package iam

import (
	"fmt"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	awsiam "github.com/aws/aws-sdk-go/service/iam"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
)

// ImportInstanceRoleFromProfileARN fetches first role ARN from instance profile
func ImportInstanceRoleFromProfileARN(provider api.ClusterProvider, ng *api.NodeGroup, profileARN string) error {
	partsOfProfileARN := strings.Split(profileARN, "/")

	if len(partsOfProfileARN) != 2 {
		return fmt.Errorf("unexpected format of instance profile ARN: %q", profileARN)
	}
	profileName := partsOfProfileARN[1]
	input := &awsiam.GetInstanceProfileInput{
		InstanceProfileName: &profileName,
	}
	output, err := provider.IAM().GetInstanceProfile(input)
	if err != nil {
		return errors.Wrap(err, "importing instance role ARN")
	}

	roles := output.InstanceProfile.Roles
	if len(roles) == 0 {
		return fmt.Errorf("instance profile %q has no roles", profileName)
	}
	if len(roles) > 1 {
		logger.Debug("instance profile %q has %d roles, only first role will be used (%#v)", profileName, roles)
	}

	ng.IAM.InstanceRoleARN = *output.InstanceProfile.Roles[0].Arn
	return nil
}

// UseFromNodeGroup retrieves the IAM configuration from an existing node group
// based on stack outputs
func UseFromNodeGroup(provider api.ClusterProvider, stack *cfn.Stack, ng *api.NodeGroup) error {
	if ng.IAM == nil {
		ng.IAM = &api.NodeGroupIAM{}
	}

	requiredCollectors := map[string]outputs.Collector{
		outputs.NodeGroupInstanceRoleARN: func(v string) error {
			ng.IAM.InstanceRoleARN = v
			return nil
		},
	}
	return outputs.Collect(*stack, requiredCollectors, nil)
}
