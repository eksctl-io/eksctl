package v1alpha4

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation"
)

func validateNodeGroupIAM(i int, ng *NodeGroup, value, fieldName, path string) error {
	if value != "" {
		p := fmt.Sprintf("%s.iam.%s and %s.iam", path, fieldName, path)
		if ng.IAM.InstanceRoleName != "" {
			return fmt.Errorf("%s.instanceRoleName cannot be set at the same time", p)
		}
		if len(ng.IAM.AttachPolicyARNs) != 0 {
			return fmt.Errorf("%s.attachPolicyARNs cannot be set at the same time", p)
		}
		if v := ng.IAM.WithAddonPolicies.AutoScaler; v != nil && *v {
			return fmt.Errorf("%s.withAddonPolicies.autoScaler cannot be set at the same time", p)
		}
		if v := ng.IAM.WithAddonPolicies.ExternalDNS; v != nil && *v {
			return fmt.Errorf("%s.withAddonPolicies.externalDNS cannot be set at the same time", p)
		}
		if v := ng.IAM.WithAddonPolicies.ImageBuilder; v != nil && *v {
			return fmt.Errorf("%s.imageBuilder cannot be set at the same time", p)
		}
		if v := ng.IAM.WithAddonPolicies.AppMesh; v != nil && *v {
			return fmt.Errorf("%s.AppMesh cannot be set at the same time", p)
		}
		if v := ng.IAM.WithAddonPolicies.EBS; v != nil && *v {
			return fmt.Errorf("%s.ebs cannot be set at the same time", p)
		}
		if v := ng.IAM.WithAddonPolicies.FSX; v != nil && *v {
			return fmt.Errorf("%s.fsx cannot be set at the same time", p)
		}
		if v := ng.IAM.WithAddonPolicies.ALBIngress; v != nil && *v {
			return fmt.Errorf("%s.albIngress cannot be set at the same time", p)
		}
	}
	return nil
}

// ValidateNodeGroupLabels uses proper Kubernetes label validation,
// it's designed to make sure users don't pass weird labels to the
// nodes, which would prevent kubelets to startup properly
func ValidateNodeGroupLabels(ng *NodeGroup) error {
	// compact version based on:
	// - https://github.com/kubernetes/kubernetes/blob/v1.13.2/cmd/kubelet/app/options/options.go#L257-L267
	// - https://github.com/kubernetes/kubernetes/blob/v1.13.2/pkg/kubelet/apis/well_known_labels.go
	// we cannot import those package because they break other dependencies

	unknownKubernetesLabels := []string{}

	for l := range ng.Labels {
		labelParts := strings.Split(l, "/")

		if len(labelParts) > 2 {
			return fmt.Errorf("node label key %q is of invalid format, can only use one '/' separator", l)
		}

		if errs := validation.IsQualifiedName(l); len(errs) > 0 {
			return fmt.Errorf("label %q is invalid - %v", l, errs)
		}
		if errs := validation.IsValidLabelValue(ng.Labels[l]); len(errs) > 0 {
			return fmt.Errorf("label %q has invalid value %q - %v", l, ng.Labels[l], errs)
		}

		isKubernetesLabel := false
		allowedKubeletNamespace := false
		if len(labelParts) == 2 {
			ns := labelParts[0]

			for _, domain := range []string{"kubernetes.io", "k8s.io"} {
				if ns == domain || strings.HasSuffix(ns, "."+domain) {
					isKubernetesLabel = true
				}
			}

			for _, domain := range []string{"kubelet.kubernetes.io", "node.kubernetes.io", "node-role.kubernetes.io"} {
				if ns == domain || strings.HasSuffix(ns, "."+domain) {
					allowedKubeletNamespace = true
				}
			}

			if isKubernetesLabel && !allowedKubeletNamespace {
				switch l {
				case
					"kubernetes.io/hostname",
					"kubernetes.io/instance-type",
					"kubernetes.io/os",
					"kubernetes.io/arch",
					"beta.kubernetes.io/instance-type",
					"beta.kubernetes.io/os",
					"beta.kubernetes.io/arch",
					"failure-domain.beta.kubernetes.io/zone",
					"failure-domain.beta.kubernetes.io/region",
					"failure-domain.kubernetes.io/zone",
					"failure-domain.kubernetes.io/region":
				default:
					unknownKubernetesLabels = append(unknownKubernetesLabels, l)
				}
			}
		}
	}

	if len(unknownKubernetesLabels) > 0 {
		return fmt.Errorf("unknown 'kubernetes.io' or 'k8s.io' labels were specified: %v", unknownKubernetesLabels)
	}
	return nil
}

// ValidateNodeGroup checks compatible fileds of a given nodegroup
func ValidateNodeGroup(i int, ng *NodeGroup) error {
	path := fmt.Sprintf("nodegroups[%d]", i)
	if ng.Name == "" {
		return fmt.Errorf("%s.name must be set", path)
	}

	if ng.IAM == nil {
		return nil
	}

	if err := validateNodeGroupIAM(i, ng, ng.IAM.InstanceProfileARN, "instanceProfileARN", path); err != nil {
		return err
	}
	if err := validateNodeGroupIAM(i, ng, ng.IAM.InstanceRoleARN, "instanceRoleARN", path); err != nil {
		return err
	}

	if err := ValidateNodeGroupLabels(ng); err != nil {
		return err
	}

	return nil
}
