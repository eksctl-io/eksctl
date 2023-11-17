package cmdutils

import (
	"errors"
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

// PodIdentityAssociationOptions holds the options for deleting a pod identity association.
type PodIdentityAssociationOptions struct {
	// Namespace is the namespace the service account belongs to.
	Namespace string
	// ServiceAccountName is the name of the Kubernetes ServiceAccount.
	ServiceAccountName string
}

func validatePodIdentityAssociation(l *commonClusterConfigLoader, options PodIdentityAssociationOptions) error {
	if l.ClusterConfig.Metadata.Name == "" {
		return ErrMustBeSet(ClusterNameFlag(l.Cmd))
	}
	if options.Namespace == "" {
		return errors.New("--namespace is required")
	}
	if options.ServiceAccountName == "" {
		return errors.New("--service-account-name is required")
	}
	return nil
}

func validatePodIdentityAssociationForConfig(clusterConfig *api.ClusterConfig) error {
	if clusterConfig.IAM == nil || len(clusterConfig.IAM.PodIdentityAssociations) == 0 {
		return errors.New("no iam.podIdentityAssociations specified in the config file")
	}
	return nil
}

// NewDeletePodIdentityAssociationLoader will load config or use flags for `eksctl delete podidentityassociation`.
func NewDeletePodIdentityAssociationLoader(cmd *Cmd, options PodIdentityAssociationOptions) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)
	l.flagsIncompatibleWithConfigFile.Insert("namespace", "service-account-name")

	l.validateWithoutConfigFile = func() error {
		return validatePodIdentityAssociation(l, options)
	}

	l.validateWithConfigFile = func() error {
		return validatePodIdentityAssociationForConfig(l.ClusterConfig)
	}
	return l
}

// UpdatePodIdentityAssociationOptions holds the options for updating a pod identity association.
type UpdatePodIdentityAssociationOptions struct {
	PodIdentityAssociationOptions
	// RoleARN is the IAM role ARN to be associated with the pod.
	RoleARN string
}

// NewUpdatePodIdentityAssociationLoader will load config or use flags for `eksctl update podidentityassociation`.
func NewUpdatePodIdentityAssociationLoader(cmd *Cmd, options UpdatePodIdentityAssociationOptions) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)
	l.flagsIncompatibleWithConfigFile.Insert("namespace", "service-account-name", "role-arn")

	l.validateWithoutConfigFile = func() error {
		if err := validatePodIdentityAssociation(l, options.PodIdentityAssociationOptions); err != nil {
			return err
		}
		if options.RoleARN == "" {
			return errors.New("--role-arn is required")
		}
		return nil
	}

	l.validateWithConfigFile = func() error {
		return validatePodIdentityAssociationForConfig(l.ClusterConfig)
	}
	return l
}
