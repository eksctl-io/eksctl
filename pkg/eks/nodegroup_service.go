package eks

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/aws/amazon-ec2-instance-selector/v2/pkg/bytequantity"
	"github.com/aws/amazon-ec2-instance-selector/v2/pkg/selector"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/weaveworks/eksctl/pkg/ami"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/outposts"
	"github.com/weaveworks/eksctl/pkg/ssh"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

// MaxInstanceTypes is the maximum number of instance types you can specify in
// a CloudFormation template.
const maxInstanceTypes = 40

// InstanceSelector selects a set of instance types matching the specified instance selector criteria.
//
//counterfeiter:generate -o fakes/fake_instance_selector.go . InstanceSelector
type InstanceSelector interface {
	// Filter returns a set of instance types matching the specified instance selector filters.
	Filter(context.Context, selector.Filters) ([]string, error)
}

// A NodeGroupService provides helpers for nodegroup creation.
type NodeGroupService struct {
	provider         api.ClusterProvider
	instanceSelector InstanceSelector
	outpostsService  *outposts.Service
}

// NewNodeGroupService creates a new NodeGroupService.
func NewNodeGroupService(provider api.ClusterProvider, instanceSelector InstanceSelector, outpostsService *outposts.Service) *NodeGroupService {
	return &NodeGroupService{
		provider:         provider,
		instanceSelector: instanceSelector,
		outpostsService:  outpostsService,
	}
}

const defaultCPUArch = "x86_64"

// Normalize normalizes nodegroups.
func (n *NodeGroupService) Normalize(ctx context.Context, nodePools []api.NodePool, clusterConfig *api.ClusterConfig) error {
	for _, np := range nodePools {
		switch ng := np.(type) {
		case *api.ManagedNodeGroup:
			if ng.LaunchTemplate == nil && ng.InstanceType == "" && len(ng.InstanceTypes) == 0 && ng.InstanceSelector.IsZero() {
				ng.InstanceType = api.DefaultNodeType
			}
			hasNativeAMIFamilySupport :=
				ng.AMIFamily == api.NodeImageFamilyAmazonLinux2023 ||
					ng.AMIFamily == api.NodeImageFamilyAmazonLinux2 ||
					ng.AMIFamily == api.NodeImageFamilyBottlerocket ||
					api.IsWindowsImage(ng.AMIFamily)

			if !hasNativeAMIFamilySupport && !api.IsAMI(ng.AMI) {
				if err := ResolveAMI(ctx, n.provider, clusterConfig.Metadata.Version, np); err != nil {
					return err
				}
			}

		case *api.NodeGroup:
			if !api.IsAMI(ng.AMI) {
				if err := ResolveAMI(ctx, n.provider, clusterConfig.Metadata.Version, ng); err != nil {
					return err
				}
			}

			if clusterConfig.IsControlPlaneOnOutposts() || ng.OutpostARN != "" {
				if err := n.outpostsService.SetOrValidateOutpostInstanceType(ctx, ng); err != nil {
					return fmt.Errorf("error setting or validating instance type for nodegroup %q: %w", ng.Name, err)
				}
			} else if ng.InstanceType == "" {
				if api.HasMixedInstances(ng) || !ng.InstanceSelector.IsZero() {
					ng.InstanceType = "mixed"
				} else {
					ng.InstanceType = api.DefaultNodeType
				}
			}
		}

		ng := np.BaseNodeGroup()
		// resolve AMI
		logger.Info("nodegroup %q will use %q [%s/%s]", ng.Name, ng.AMI, ng.AMIFamily, clusterConfig.Metadata.Version)

		if ng.AMI != "" {
			if err := ami.Use(ctx, n.provider.EC2(), ng); err != nil {
				return err
			}
		}
		// load or use SSH key - name includes cluster name and the
		// fingerprint, so if unique keys are provided, each will get
		// loaded and used as intended and there is no need to have
		// nodegroup name in the key name
		publicKeyName, err := ssh.LoadKey(ctx, ng.SSH, clusterConfig.Metadata.Name, ng.Name, n.provider.EC2())
		if err != nil {
			return err
		}
		if publicKeyName != "" {
			ng.SSH.PublicKeyName = &publicKeyName
		}
	}
	return nil
}

// ExpandInstanceSelectorOptions sets instance types to instances matched by the instance selector criteria.
func (n *NodeGroupService) ExpandInstanceSelectorOptions(nodePools []api.NodePool, clusterAZs []string) error {
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
		instanceTypes, err := n.expandInstanceSelector(baseNG.InstanceSelector, azs)
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

func (n *NodeGroupService) expandInstanceSelector(ins *api.InstanceSelector, azs []string) ([]string, error) {
	makeRange := func(val int) *selector.Int32RangeFilter {
		return &selector.Int32RangeFilter{
			LowerBound: int32(val),
			UpperBound: int32(val),
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

	filters.CPUArchitecture = (*ec2types.ArchitectureType)(aws.String(cpuArch))

	instanceTypes, err := n.instanceSelector.Filter(context.TODO(), filters)
	if err != nil {
		return nil, errors.Wrap(err, "error querying instance types for the specified instance selector criteria")
	}
	if len(instanceTypes) == 0 {
		return nil, errors.New("instance selector criteria matched no instances; consider broadening your criteria so that more instance types are returned")
	}

	return instanceTypes, nil
}

// DoAllNodegroupStackTasks iterates over nodegroup tasks and returns any errors.
func DoAllNodegroupStackTasks(taskTree *tasks.TaskTree, region, name string) error {
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
// validates configuration, if it find issues it logs messages.
func ValidateExistingNodeGroupsForCompatibility(ctx context.Context, cfg *api.ClusterConfig, stackManager manager.StackManager) error {
	infoByNodeGroup, err := stackManager.DescribeNodeGroupStacksAndResources(ctx)
	if err != nil {
		return errors.Wrap(err, "getting resources for all nodegroup stacks")
	}
	if len(infoByNodeGroup) == 0 {
		return nil
	}

	logger.Info("checking security group configuration for all nodegroups")
	var incompatibleNodeGroups []string
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

	if len(incompatibleNodeGroups) == 0 {
		logger.Info("all nodegroups have up-to-date cloudformation templates")
		return nil
	}

	logger.Critical("found %d nodegroup(s) (%s) without shared security group, cluster networking maybe be broken",
		len(incompatibleNodeGroups), strings.Join(incompatibleNodeGroups, ", "))
	logger.Critical("it's recommended to create new nodegroups, then delete old ones")
	if cfg.VPC.SharedNodeSecurityGroup != "" {
		logger.Critical("as a temporary fix, you can patch the configuration and add each of these nodegroup(s) to %q",
			cfg.VPC.SharedNodeSecurityGroup)
	}

	return nil
}
