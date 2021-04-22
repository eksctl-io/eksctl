package eks

import (
	"reflect"

	"github.com/aws/amazon-ec2-instance-selector/v2/pkg/bytequantity"
	"github.com/aws/amazon-ec2-instance-selector/v2/pkg/selector"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/ami"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ssh"
)

// InstanceSelector selects a set of instance types matching the specified instance selector criteria
//go:generate counterfeiter -o fakes/fake_instance_selector.go . InstanceSelector
type InstanceSelector interface {
	// Filter returns a set of instance types matching the specified instance selector filters
	Filter(selector.Filters) ([]string, error)
}

// A NodeGroupService provides helpers for nodegroup creation
type NodeGroupService struct {
	provider         api.ClusterProvider
	instanceSelector InstanceSelector
}

// NewNodeGroupService creates a new NodeGroupService
func NewNodeGroupService(provider api.ClusterProvider, instanceSelector InstanceSelector) *NodeGroupService {
	return &NodeGroupService{
		provider:         provider,
		instanceSelector: instanceSelector,
	}
}

const defaultCPUArch = "x86_64"

// Normalize normalizes nodegroups
func (m *NodeGroupService) Normalize(nodePools []api.NodePool, clusterMeta *api.ClusterMeta) error {
	for _, np := range nodePools {
		switch ng := np.(type) {
		case *api.NodeGroup:
			// TODO remove
			// This is a temporary hack to go down a legacy bootstrap codepath for Ubuntu
			// and AL2 images
			if ng.AMI != "" {
				logger.Warning("Custom AMI detected for nodegroup %s. Please refer to https://github.com/weaveworks/eksctl/issues/3563 for upcoming breaking changes", ng.Name)
				ng.CustomAMI = true
			}
			// resolve AMI
			if !api.IsAMI(ng.AMI) {
				if err := ResolveAMI(m.provider, clusterMeta.Version, ng); err != nil {
					return err
				}
			}
			logger.Info("nodegroup %q will use %q [%s/%s]", ng.Name, ng.AMI, ng.AMIFamily, clusterMeta.Version)
		}

		ng := np.BaseNodeGroup()
		if ng.AMI != "" {
			if err := ami.Use(m.provider.EC2(), ng); err != nil {
				return err
			}
		}
		// load or use SSH key - name includes cluster name and the
		// fingerprint, so if unique keys are provided, each will get
		// loaded and used as intended and there is no need to have
		// nodegroup name in the key name
		publicKeyName, err := ssh.LoadKey(ng.SSH, clusterMeta.Name, ng.Name, m.provider.EC2())
		if err != nil {
			return err
		}
		if publicKeyName != "" {
			ng.SSH.PublicKeyName = &publicKeyName
		}
	}
	return nil
}

// ExpandInstanceSelectorOptions sets instance types to instances matched by the instance selector criteria
func (m *NodeGroupService) ExpandInstanceSelectorOptions(nodePools []api.NodePool, clusterAZs []string) error {
	instanceTypesMatch := func(a, b []string) bool {
		return reflect.DeepEqual(a, b)
	}

	instanceTypesMismatchErr := func(ng *api.NodeGroupBase, path string) error {
		return errors.Errorf("instance types matched by instance selector criteria do not match %s.instanceTypes for nodegroup %q; either remove instanceSelector or instanceTypes and retry the operation", path, ng.Name)
	}

	for _, np := range nodePools {
		baseNG := np.BaseNodeGroup()
		if baseNG.InstanceSelector == nil || baseNG.InstanceSelector.IsZero() {
			continue
		}

		azs := clusterAZs
		if len(baseNG.AvailabilityZones) != 0 {
			azs = baseNG.AvailabilityZones
		}
		instanceTypes, err := m.expandInstanceSelector(baseNG.InstanceSelector, azs)
		if err != nil {
			return errors.Wrapf(err, "error expanding instance selector options for nodegroup %q", baseNG.Name)
		}

		switch ng := np.(type) {
		case *api.NodeGroup:
			if ng.InstancesDistribution == nil {
				ng.InstancesDistribution = &api.NodeGroupInstancesDistribution{}
			}
			if len(ng.InstancesDistribution.InstanceTypes) > 0 {
				if !instanceTypesMatch(ng.InstancesDistribution.InstanceTypes, instanceTypes) {
					return instanceTypesMismatchErr(baseNG, "nodeGroup.instancesDistribution")
				}
			} else {
				ng.InstancesDistribution.InstanceTypes = instanceTypes
			}

		case *api.ManagedNodeGroup:
			if len(ng.InstanceTypes) > 0 {
				if !instanceTypesMatch(ng.InstanceTypes, instanceTypes) {
					return instanceTypesMismatchErr(baseNG, "managedNodeGroup.instanceTypes")
				}
			} else {
				ng.InstanceTypes = instanceTypes
			}

		default:
			return errors.Errorf("unhandled NodePool type %T", np)
		}
	}
	return nil
}

func (m *NodeGroupService) expandInstanceSelector(ins *api.InstanceSelector, azs []string) ([]string, error) {
	makeRange := func(val int) *selector.IntRangeFilter {
		return &selector.IntRangeFilter{
			LowerBound: val,
			UpperBound: val,
		}
	}

	filters := selector.Filters{
		Service:           aws.String("eks"),
		AvailabilityZones: &azs,
	}
	if ins.VCPUs != 0 {
		filters.VCpusRange = makeRange(ins.VCPUs)
	}
	if ins.Memory != "" {
		memory, err := bytequantity.ParseToByteQuantity(ins.Memory)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid value %q for instanceSelector.memory", ins.Memory)
		}
		filters.MemoryRange = &selector.ByteQuantityRangeFilter{
			LowerBound: memory,
			UpperBound: memory,
		}
	}
	if ins.GPUs != 0 {
		filters.GpusRange = makeRange(ins.GPUs)
	}
	cpuArch := ins.CPUArchitecture
	if cpuArch == "" {
		cpuArch = defaultCPUArch
	}
	filters.CPUArchitecture = aws.String(cpuArch)

	instanceTypes, err := m.instanceSelector.Filter(filters)
	if err != nil {
		return nil, errors.Wrap(err, "error querying instance types for the specified instance selector criteria")
	}
	if len(instanceTypes) == 0 {
		return nil, errors.New("instance selector criteria matched no instances; consider broadening your criteria so that more instance types are returned")
	}

	return instanceTypes, nil
}
