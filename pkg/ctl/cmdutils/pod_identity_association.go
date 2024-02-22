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

	l.flagsIncompatibleWithConfigFile = sets.New[string](podIdentityAssociationFlagsIncompatibleWithConfigFile...)

	l.validateWithConfigFile = func() error {
		return validatePodIdentityAssociationsForConfig(l.ClusterConfig, true)
	}

	l.validateWithoutConfigFile = func() error {
		if err := validatePodIdentityAssociation(l, PodIdentityAssociationOptions{
			Namespace:          podIdentityAssociation.Namespace,
			ServiceAccountName: podIdentityAssociation.ServiceAccountName,
		}); err != nil {
			return err
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

	l.flagsIncompatibleWithConfigFile = sets.New[string]("cluster")

	l.validateWithoutConfigFile = func() error {
		if cmd.ClusterConfig.Metadata.Name == "" {
			return ErrMustBeSet(ClusterNameFlag(cmd))
		}
		if pia.Namespace == "" && pia.ServiceAccountName != "" {
			return fmt.Errorf("--namespace must be set in order to specify --service-account-name")
		}
		return nil
	}

	l.validateWithConfigFile = func() error {
		if cmd.ClusterConfig.Metadata.Name == "" {
			return ErrMustBeSet(ClusterNameFlag(cmd))
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

func validatePodIdentityAssociationsForConfig(clusterConfig *api.ClusterConfig, isCreate bool) error {
	if clusterConfig.IAM == nil || len(clusterConfig.IAM.PodIdentityAssociations) == 0 {
		return errors.New("no iam.podIdentityAssociations specified in the config file")
	}

	for i, pia := range clusterConfig.IAM.PodIdentityAssociations {
		path := fmt.Sprintf("podIdentityAssociations[%d]", i)
		if pia.Namespace == "" {
			return fmt.Errorf("%s.namespace must be set", path)
		}
		if pia.ServiceAccountName == "" {
			return fmt.Errorf("%s.serviceAccountName must be set", path)
		}

		if !isCreate {
			continue
		}

		if pia.RoleARN == "" &&
			len(pia.PermissionPolicy) == 0 &&
			len(pia.PermissionPolicyARNs) == 0 &&
			!pia.WellKnownPolicies.HasPolicy() {
			return fmt.Errorf("at least one of the following must be specified: %[1]s.roleARN, %[1]s.permissionPolicy, %[1]s.permissionPolicyARNs, %[1]s.wellKnownPolicies", path)
		}
		if pia.RoleARN != "" {
			if len(pia.PermissionPolicy) > 0 {
				return fmt.Errorf("%[1]s.permissionPolicy cannot be specified when %[1]s.roleARN is set", path)
			}
			if len(pia.PermissionPolicyARNs) > 0 {
				return fmt.Errorf("%[1]s.permissionPolicyARNs cannot be specified when %[1]s.roleARN is set", path)
			}
			if pia.WellKnownPolicies.HasPolicy() {
				return fmt.Errorf("%[1]s.wellKnownPolicies cannot be specified when %[1]s.roleARN is set", path)
			}
		}
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
		return validatePodIdentityAssociationsForConfig(l.ClusterConfig, false)
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
		return validatePodIdentityAssociationsForConfig(l.ClusterConfig, false)
	}
	return l
}
