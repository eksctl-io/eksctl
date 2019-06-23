package v1alpha5

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
		if IsEnabled(ng.IAM.WithAddonPolicies.AutoScaler) {
			return fmt.Errorf("%s.withAddonPolicies.autoScaler cannot be set at the same time", p)
		}
		if IsEnabled(ng.IAM.WithAddonPolicies.ExternalDNS) {
			return fmt.Errorf("%s.withAddonPolicies.externalDNS cannot be set at the same time", p)
		}
		if IsEnabled(ng.IAM.WithAddonPolicies.CertManager) {
			return fmt.Errorf("%s.withAddonPolicies.certManager cannot be set at the same time", p)
		}
		if IsEnabled(ng.IAM.WithAddonPolicies.ImageBuilder) {
			return fmt.Errorf("%s.imageBuilder cannot be set at the same time", p)
		}
		if IsEnabled(ng.IAM.WithAddonPolicies.AppMesh) {
			return fmt.Errorf("%s.AppMesh cannot be set at the same time", p)
		}
		if IsEnabled(ng.IAM.WithAddonPolicies.EBS) {
			return fmt.Errorf("%s.ebs cannot be set at the same time", p)
		}
		if IsEnabled(ng.IAM.WithAddonPolicies.FSX) {
			return fmt.Errorf("%s.fsx cannot be set at the same time", p)
		}
		if IsEnabled(ng.IAM.WithAddonPolicies.EFS) {
			return fmt.Errorf("%s.efs cannot be set at the same time", p)
		}
		if IsEnabled(ng.IAM.WithAddonPolicies.ALBIngress) {
			return fmt.Errorf("%s.albIngress cannot be set at the same time", p)
		}
		if IsEnabled(ng.IAM.WithAddonPolicies.XRay) {
			return fmt.Errorf("%s.xRay cannot be set at the same time", p)
		}
		if IsEnabled(ng.IAM.WithAddonPolicies.CloudWatch) {
			return fmt.Errorf("%s.cloudWatch cannot be set at the same time", p)
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

	if ng.VolumeSize == nil {
		if IsSetAndNonEmptyString(ng.VolumeType) {
			return fmt.Errorf("%s.volumeType can not be set without %s.volumeSize", path, path)
		}
		if IsSetAndNonEmptyString(ng.VolumeName) {
			return fmt.Errorf("%s.volumeName can not be set without %s.volumeSize", path, path)
		}
	}

	if ng.IAM != nil {
		if err := validateNodeGroupIAM(i, ng, ng.IAM.InstanceProfileARN, "instanceProfileARN", path); err != nil {
			return err
		}
		if err := validateNodeGroupIAM(i, ng, ng.IAM.InstanceRoleARN, "instanceRoleARN", path); err != nil {
			return err
		}

		if err := ValidateNodeGroupLabels(ng); err != nil {
			return err
		}

		if err := validateNodeGroupSSH(ng.SSH); err != nil {
			return fmt.Errorf("only one ssh public key can be specified per node-group")
		}
	}

	if err := validateNodeGroupKubeletExtraConfig(ng.KubeletExtraConfig); err != nil {
		return err
	}

	if err := validateInstancesDistribution(ng); err != nil {
		return err
	}

	return nil
}

func validateInstancesDistribution(ng *NodeGroup) error {
	if ng.InstancesDistribution == nil {
		return nil
	}

	if ng.InstanceType != "" && ng.InstanceType != "mixed" {
		return fmt.Errorf("instanceType should be \"mixed\" or unset when using the mixed instances feature")
	}

	distribution := ng.InstancesDistribution
	if distribution.InstanceTypes == nil || len(distribution.InstanceTypes) == 0 {
		return fmt.Errorf("at least two instance types have to be specified for mixed nodegroups")
	}

	allInstanceTypes := make(map[string]bool)
	for _, instanceType := range distribution.InstanceTypes {
		allInstanceTypes[instanceType] = true
	}

	if len(allInstanceTypes) < 2 || len(allInstanceTypes) > 20 {
		return fmt.Errorf("mixed nodegroups should have between 2 and 20 different instance types")
	}

	if distribution.OnDemandBaseCapacity != nil && *distribution.OnDemandBaseCapacity < 0 {
		return fmt.Errorf("onDemandBaseCapacity should be 0 or more")
	}

	if distribution.OnDemandPercentageAboveBaseCapacity != nil && (*distribution.OnDemandPercentageAboveBaseCapacity < 0 || *distribution.OnDemandPercentageAboveBaseCapacity > 100) {
		return fmt.Errorf("percentageAboveBase should be between 0 and 100")
	}

	if distribution.SpotInstancePools != nil && (*distribution.SpotInstancePools < 1 || *distribution.SpotInstancePools > 20) {
		return fmt.Errorf("spotInstancePools should be between 1 and 20")
	}

	return nil
}

func validateNodeGroupSSH(SSH *NodeGroupSSH) error {
	if SSH == nil {
		return nil
	}
	numSSHFlagsEnabled := countEnabledFields(
		SSH.PublicKeyPath,
		SSH.PublicKey,
		SSH.PublicKeyName)

	if numSSHFlagsEnabled > 1 {
		return fmt.Errorf("only one ssh public key can be specified per node-group")
	}
	return nil
}

func countEnabledFields(fields ...*string) int {
	count := 0
	for _, flag := range fields {
		if flag != nil && *flag != "" {
			count++
		}
	}
	return count
}

func validateNodeGroupKubeletExtraConfig(kubeletConfig *NodeGroupKubeletConfig) error {
	if kubeletConfig == nil {
		return nil
	}

	var kubeletForbiddenFields = map[string]struct{}{
		"kind":               {},
		"apiVersion":         {},
		"address":            {},
		"clusterDomain":      {},
		"authentication":     {},
		"authorization":      {},
		"serverTLSBootstrap": {},
		"featureGates":       {},
	}

	for k := range *kubeletConfig {
		if _, exists := kubeletForbiddenFields[k]; exists {
			return fmt.Errorf("cannot override %q in kubelet config, as it's critical to eksctl functionality", k)
		}
	}
	return nil
}
