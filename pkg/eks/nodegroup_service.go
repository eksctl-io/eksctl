package eks

import (
	"reflect"

	"github.com/aws/amazon-ec2-instance-selector/v2/pkg/bytequantity"
	"github.com/aws/amazon-ec2-instance-selector/v2/pkg/selector"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/ami"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ssh"
)

// A NodeGroupService provides helpers for nodegroup creation
type NodeGroupService struct {
	cluster  *api.ClusterConfig
	provider api.ClusterProvider
}

// NewNodeGroupService creates a new NodeGroupService
func NewNodeGroupService(clusterConfig *api.ClusterConfig, provider api.ClusterProvider) *NodeGroupService {
	return &NodeGroupService{
		cluster:  clusterConfig,
		provider: provider,
	}
}

const defaultCPUArch = "x86_64"

// Normalize normalizes nodegroups
func (m *NodeGroupService) Normalize(nodePools []api.NodePool) error {
	for _, np := range nodePools {
		switch ng := np.(type) {
		case *api.NodeGroup:
			// resolve AMI
			if !api.IsAMI(ng.AMI) {
				if err := ResolveAMI(m.provider, m.cluster.Metadata.Version, ng); err != nil {
					return err
				}
			}
			logger.Info("nodegroup %q will use %q [%s/%s]", ng.Name, ng.AMI, ng.AMIFamily, m.cluster.Metadata.Version)
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
		publicKeyName, err := ssh.LoadKey(ng.SSH, m.cluster.Metadata.Name, ng.Name, m.provider.EC2())
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
func (m *NodeGroupService) ExpandInstanceSelectorOptions(nodePools []api.NodePool) error {
	sess, ok := m.provider.ConfigProvider().(*session.Session)
	if !ok {
		return errors.Errorf("expected ConfigProvider to be of type %T; got %T", &session.Session{}, m.provider.ConfigProvider())
	}
	instanceSelector := selector.New(sess)

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

		instanceTypes, err := m.expandInstanceSelector(instanceSelector, baseNG.InstanceSelector)
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

func (m *NodeGroupService) expandInstanceSelector(instanceSelector *selector.Selector, ins *api.InstanceSelector) ([]string, error) {
	makeRange := func(val int) *selector.IntRangeFilter {
		return &selector.IntRangeFilter{
			LowerBound: val,
			UpperBound: val,
		}
	}

	filters := selector.Filters{
		Service: aws.String("eks"),
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

	instanceTypes, err := instanceSelector.Filter(filters)
	if err != nil {
		return nil, errors.Wrap(err, "error querying instance types for the specified instance selector criteria")
	}
	if len(instanceTypes) == 0 {
		return nil, errors.New("instance selector criteria matched no instances; consider broadening your criteria so that more instance types are returned")
	}

	return instanceTypes, nil
}
