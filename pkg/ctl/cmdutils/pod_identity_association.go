package cmdutils

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var (
	podIdentityAssociationFlagsIncompatibleWithConfigFile = []string{
		"cluster",
		"namespace",
		"service-account-name",
		"role-arn",
		"role-name",
		"permission-boundary-arn",
		"permission-policy-arn",
		"well-known-policies",
	}
)

// NewCreatePodIdentityAssociationLoader will load config or use flags for 'eksctl create podidentityassociation'.
func NewCreatePodIdentityAssociationLoader(cmd *Cmd, podIdentityAssociation *api.PodIdentityAssociation) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	l.flagsIncompatibleWithConfigFile = sets.NewString(podIdentityAssociationFlagsIncompatibleWithConfigFile...)

	l.validateWithConfigFile = func() error {
		if len(cmd.ClusterConfig.IAM.PodIdentityAssociations) == 0 {
			return fmt.Errorf("at least one pod identity association is required")
		}
		return nil
	}

	l.validateWithoutConfigFile = func() error {
		if l.ClusterConfig.Metadata.Name == "" {
			return ErrMustBeSet(ClusterNameFlag(cmd))
		}
		if podIdentityAssociation.Namespace == "" {
			return ErrMustBeSet("--namespace")
		}
		if podIdentityAssociation.ServiceAccountName == "" {
			return ErrMustBeSet("--service-account-name")
		}
		if podIdentityAssociation.RoleARN == "" &&
			len(podIdentityAssociation.PermissionPolicyARNs) == 0 &&
			!podIdentityAssociation.WellKnownPolicies.HasPolicy() {
			return fmt.Errorf("at least one of the following flags must be specified: --role-arn, --permission-policy-arns, --well-known-policies")
		}
		if podIdentityAssociation.RoleARN != "" {
			if len(podIdentityAssociation.PermissionPolicyARNs) > 0 {
				return fmt.Errorf("--permission-policy-arns cannot be specified when --role-arn is set")
			}
			if podIdentityAssociation.WellKnownPolicies.HasPolicy() {
				return fmt.Errorf("--well-known-policies cannot be specified when --role-arn is set")
			}
		}
		l.Cmd.ClusterConfig.IAM.PodIdentityAssociations = []api.PodIdentityAssociation{*podIdentityAssociation}
		return nil
	}

	return l
}

func NewGetPodIdentityAssociationLoader(cmd *Cmd, pia *api.PodIdentityAssociation) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	l.validateWithoutConfigFile = func() error {
		if cmd.ClusterConfig.Metadata.Name == "" {
			return ErrMustBeSet(ClusterNameFlag(cmd))
		}
		if pia.Namespace == "" && pia.ServiceAccountName != "" {
			return fmt.Errorf("--namespace must be set in order to specify --service-account-name")
		}
		return nil
	}
	return l
}
