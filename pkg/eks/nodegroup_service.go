package eks

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/aws/amazon-ec2-instance-selector/v2/pkg/bytequantity"
	"github.com/aws/amazon-ec2-instance-selector/v2/pkg/selector"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/ami"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/ssh"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

// MaxInstanceTypes is the maximum number of instance types you can specify in
// a CloudFormation template
const maxInstanceTypes = 40

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate -o fakes/fake_instance_selector.go . InstanceSelector
// InstanceSelector selects a set of instance types matching the specified instance selector criteria
type InstanceSelector interface {
	// Filter returns a set of instance types matching the specified instance selector filters
	Filter(selector.Filters) ([]string, error)
}

//counterfeiter:generate -o fakes/fake_nodegroup_initialiser.go . NodeGroupInitialiser
// NodeGroupInitialiser is an interface that provides helpers for nodegroup creation.
type NodeGroupInitialiser interface {
	Normalize(nodePools []api.NodePool, clusterMeta *api.ClusterMeta) error
	ExpandInstanceSelectorOptions(nodePools []api.NodePool, clusterAZs []string) error
	NewAWSSelectorSession(provider api.ClusterProvider)
	ValidateLegacySubnetsForNodeGroups(spec *api.ClusterConfig, provider api.ClusterProvider) error
	DoesAWSNodeUseIRSA(provider api.ClusterProvider, clientSet kubernetes.Interface) (bool, error)
	DoAllNodegroupStackTasks(taskTree *tasks.TaskTree, region, name string) error
	ValidateExistingNodeGroupsForCompatibility(cfg *api.ClusterConfig, stackManager manager.StackManager) error
}

// A NodeGroupService provides helpers for nodegroup creation
type NodeGroupService struct {
	Provider         api.ClusterProvider
	instanceSelector InstanceSelector
}

// NewNodeGroupService creates a new NodeGroupService
func NewNodeGroupService(provider api.ClusterProvider, instanceSelector InstanceSelector) *NodeGroupService {
	return &NodeGroupService{
		Provider:         provider,
		instanceSelector: instanceSelector,
	}
}

const defaultCPUArch = "x86_64"

// NewAWSSelectorSession returns a new instance of Selector provided an aws session
func (m *NodeGroupService) NewAWSSelectorSession(provider api.ClusterProvider) {
	m.instanceSelector = selector.New(provider.Session())
}

// Normalize normalizes nodegroups
func (m *NodeGroupService) Normalize(nodePools []api.NodePool, clusterMeta *api.ClusterMeta) error {
	for _, np := range nodePools {
		switch ng := np.(type) {
		case *api.ManagedNodeGroup:
			hasNativeAMIFamilySupport := ng.AMIFamily == api.NodeImageFamilyAmazonLinux2 || ng.AMIFamily == api.NodeImageFamilyBottlerocket
			if !hasNativeAMIFamilySupport && !api.IsAMI(ng.AMI) {
				if err := ResolveAMI(m.Provider, clusterMeta.Version, np); err != nil {
					return err
				}
			}

		case *api.NodeGroup:
			if !api.IsAMI(ng.AMI) {
				if err := ResolveAMI(m.Provider, clusterMeta.Version, ng); err != nil {
					return err
				}
			} else {
				// TODO remove
				// This is a temporary hack to go down a legacy bootstrap codepath for Ubuntu
				// and AL2 images
				logger.Warning("Custom AMI detected for nodegroup %s. Please refer to https://github.com/weaveworks/eksctl/issues/3563 for upcoming breaking changes", ng.Name)
				ng.CustomAMI = true
			}
		}

		ng := np.BaseNodeGroup()
		// resolve AMI
		logger.Info("nodegroup %q will use %q [%s/%s]", ng.Name, ng.AMI, ng.AMIFamily, clusterMeta.Version)

		if ng.AMI != "" {
			if err := ami.Use(m.Provider.EC2(), ng); err != nil {
				return err
			}
		}
		// load or use SSH key - name includes cluster name and the
		// fingerprint, so if unique keys are provided, each will get
		// loaded and used as intended and there is no need to have
		// nodegroup name in the key name
		publicKeyName, err := ssh.LoadKey(ng.SSH, clusterMeta.Name, ng.Name, m.Provider.EC2())
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

		if len(instanceTypes) > maxInstanceTypes {
			return errors.Errorf("instance selector filters resulted in %d instance types, which is greater than the maximum of %d, please set more selector options", len(instanceTypes), maxInstanceTypes)
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
	if ins.GPUs != nil {
		filters.GpusRange = makeRange(*ins.GPUs)
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

func (m *NodeGroupService) ValidateLegacySubnetsForNodeGroups(spec *api.ClusterConfig, provider api.ClusterProvider) error {
	return vpc.ValidateLegacySubnetsForNodeGroups(spec, provider)
}

// DoAllNodegroupStackTasks iterates over nodegroup tasks and returns any errors.
func (m *NodeGroupService) DoAllNodegroupStackTasks(taskTree *tasks.TaskTree, region, name string) error {
	logger.Info(taskTree.Describe())
	errs := taskTree.DoAllSync()
	if len(errs) > 0 {
		logger.Info("%d error(s) occurred and nodegroups haven't been created properly, you may wish to check CloudFormation console", len(errs))
		logger.Info("to cleanup resources, run 'eksctl delete nodegroup --region=%s --cluster=%s --name=<name>' for each of the failed nodegroup", region, name)
		for _, err := range errs {
			if err != nil {
				logger.Critical("%s\n", err.Error())
			}
		}
		return fmt.Errorf("failed to create nodegroups for cluster %q", name)
	}
	return nil
}

// ValidateExistingNodeGroupsForCompatibility looks at each of the existing nodegroups and
// validates configuration, if it find issues it logs messages
func (m *NodeGroupService) ValidateExistingNodeGroupsForCompatibility(cfg *api.ClusterConfig, stackManager manager.StackManager) error {
	infoByNodeGroup, err := stackManager.DescribeNodeGroupStacksAndResources()
	if err != nil {
		return errors.Wrap(err, "getting resources for all nodegroup stacks")
	}
	if len(infoByNodeGroup) == 0 {
		return nil
	}

	logger.Info("checking security group configuration for all nodegroups")
	incompatibleNodeGroups := []string{}
	for ng, info := range infoByNodeGroup {
		if stackManager.StackStatusIsNotTransitional(info.Stack) {
			isCompatible, err := isNodeGroupCompatible(ng, info)
			if err != nil {
				return err
			}
			if isCompatible {
				logger.Debug("nodegroup %q is compatible", ng)
			} else {
				logger.Debug("nodegroup %q is incompatible", ng)
				incompatibleNodeGroups = append(incompatibleNodeGroups, ng)
			}
		}
	}

	numIncompatibleNodeGroups := len(incompatibleNodeGroups)
	if numIncompatibleNodeGroups == 0 {
		logger.Info("all nodegroups have up-to-date cloudformation templates")
		return nil
	}

	logger.Critical("found %d nodegroup(s) (%s) without shared security group, cluster networking maybe be broken",
		numIncompatibleNodeGroups, strings.Join(incompatibleNodeGroups, ", "))
	logger.Critical("it's recommended to create new nodegroups, then delete old ones")
	if cfg.VPC.SharedNodeSecurityGroup != "" {
		logger.Critical("as a temporary fix, you can patch the configuration and add each of these nodegroup(s) to %q",
			cfg.VPC.SharedNodeSecurityGroup)
	}

	return nil
}
