package create

import (
	"fmt"
	"strings"

	"github.com/gobwas/glob"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/weaveworks/eksctl/pkg/ami"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/eks"
)

// NodeGroupFilter holds filter configuration
type NodeGroupFilter struct {
	IgnoreAllExisting bool

	existing sets.String
	only     []glob.Glob
	onlySpec string
}

// NewNodeGroupFilter create new NodeGroupFilter instance
func NewNodeGroupFilter() *NodeGroupFilter {
	return &NodeGroupFilter{
		IgnoreAllExisting: true,

		existing: sets.NewString(),
	}
}

// ApplyOnlyFilter parses given globs for exclusive filtering
func (f *NodeGroupFilter) ApplyOnlyFilter(globExprs []string, cfg *api.ClusterConfig) error {
	for _, expr := range globExprs {
		compiledExpr, err := glob.Compile(expr)
		if err != nil {
			return errors.Wrapf(err, "parsing glob filter %q", expr)
		}
		f.only = append(f.only, compiledExpr)
	}
	f.onlySpec = strings.Join(globExprs, ",")
	return f.onlyFilterMatchesAnything(cfg)
}

func (f *NodeGroupFilter) onlyFilterMatchesAnything(cfg *api.ClusterConfig) error {
	if len(f.only) == 0 {
		return nil
	}
	for i := range cfg.NodeGroups {
		for _, g := range f.only {
			if g.Match(cfg.NodeGroups[i].Name) {
				return nil
			}
		}
	}
	return fmt.Errorf("no nodegroups match filter specification: %q", f.onlySpec)
}

// ApplyExistingFilter uses stackManager to list existing nodegroup stacks and configures
// the filter accordingly
func (f *NodeGroupFilter) ApplyExistingFilter(stackManager *manager.StackCollection) error {
	if !f.IgnoreAllExisting {
		return nil
	}

	existing, err := stackManager.ListNodeGroupStacks()
	if err != nil {
		return err
	}

	f.existing.Insert(existing...)

	return nil
}

// Match checks given nodegroup against the filter
func (f *NodeGroupFilter) Match(ng *api.NodeGroup) bool {
	if f.IgnoreAllExisting && f.existing.Has(ng.Name) {
		return false
	}

	for _, g := range f.only {
		if g.Match(ng.Name) {
			return true // return first match
		}
	}

	// if no globs were given, match everything,
	// if false - we haven't matched anything so far
	return len(f.only) == 0
}

// MatchAll checks all nodegroups against the filter and returns all of
// matching names as set
func (f *NodeGroupFilter) MatchAll(cfg *api.ClusterConfig) sets.String {
	names := sets.NewString()
	for _, ng := range cfg.NodeGroups {
		if f.Match(ng) {
			names.Insert(ng.Name)
		}
	}
	return names
}

// LogInfo prints out a user-friendly message about how filter was applied
func (f *NodeGroupFilter) LogInfo(cfg *api.ClusterConfig) {
	count := f.MatchAll(cfg).Len()
	filteredOutCount := len(cfg.NodeGroups) - count
	if filteredOutCount > 0 {
		reasons := []string{}
		if f.onlySpec != "" {
			reasons = append(reasons, fmt.Sprintf("--only=%q was given", f.onlySpec))
		}
		if existingCount := f.existing.Len(); existingCount > 0 {
			reasons = append(reasons, fmt.Sprintf("%d nodegroup(s) (%s) already exist", existingCount, strings.Join(f.existing.List(), ", ")))
		}
		logger.Info("%d nodegroup(s) were filtered out: %s", filteredOutCount, strings.Join(reasons, ", "))
	}
}

// CheckEachNodeGroup iterates over each nodegroup and calls check function
// (this is needed to avoid common goroutine-for-loop pitfall)
func CheckEachNodeGroup(f *NodeGroupFilter, cfg *api.ClusterConfig, check func(i int, ng *api.NodeGroup) error) error {
	for i, ng := range cfg.NodeGroups {
		if f.Match(ng) {
			if err := check(i, ng); err != nil {
				return err
			}
		}
	}
	return nil
}

// NewNodeGroupChecker validates a new nodegroup and applies defaults
func NewNodeGroupChecker(i int, ng *api.NodeGroup) error {
	if err := api.ValidateNodeGroup(i, ng); err != nil {
		return err
	}

	// apply defaults
	if ng.InstanceType == "" {
		ng.InstanceType = api.DefaultNodeType
	}
	if ng.AMIFamily == "" {
		ng.AMIFamily = ami.ImageFamilyAmazonLinux2
	}
	if ng.AMI == "" {
		ng.AMI = ami.ResolverStatic
	}

	if ng.SecurityGroups == nil {
		ng.SecurityGroups = &api.NodeGroupSGs{
			AttachIDs: []string{},
		}
	}
	if ng.SecurityGroups.WithLocal == nil {
		ng.SecurityGroups.WithLocal = api.NewBoolTrue()
	}
	if ng.SecurityGroups.WithShared == nil {
		ng.SecurityGroups.WithShared = api.NewBoolTrue()
	}

	if ng.AllowSSH {
		if ng.SSHPublicKeyPath == "" {
			ng.SSHPublicKeyPath = defaultSSHPublicKey
		}
	}

	if ng.VolumeSize > 0 {
		if ng.VolumeType == "" {
			ng.VolumeType = api.DefaultNodeVolumeType
		}
	}

	if ng.IAM == nil {
		ng.IAM = &api.NodeGroupIAM{}
	}
	if ng.IAM.WithAddonPolicies.ImageBuilder == nil {
		ng.IAM.WithAddonPolicies.ImageBuilder = api.NewBoolFalse()
	}
	if ng.IAM.WithAddonPolicies.AutoScaler == nil {
		ng.IAM.WithAddonPolicies.AutoScaler = api.NewBoolFalse()
	}
	if ng.IAM.WithAddonPolicies.ExternalDNS == nil {
		ng.IAM.WithAddonPolicies.ExternalDNS = api.NewBoolFalse()
	}

	return nil
}

// When passing the --without-nodegroup option, don't create nodegroups
func skipNodeGroupsIfRequested(cfg *api.ClusterConfig) {
	if withoutNodeGroup {
		cfg.NodeGroups = nil
		logger.Warning("cluster will be created without an initial nodegroup")
	}
}

func checkSubnetsGiven(cfg *api.ClusterConfig) bool {
	return cfg.VPC.Subnets != nil && len(cfg.VPC.Subnets.Private)+len(cfg.VPC.Subnets.Public) != 0
}

func checkSubnetsGivenAsFlags() bool {
	return len(*subnets[api.SubnetTopologyPrivate])+len(*subnets[api.SubnetTopologyPublic]) != 0
}

func checkVersion(ctl *eks.ClusterProvider, meta *api.ClusterMeta) error {
	switch meta.Version {
	case "auto":
		break
	case "":
		meta.Version = "auto"
	case "latest":
		meta.Version = api.LatestVersion
		logger.Info("will use version latest version (%s) for new nodegroup(s)", meta.Version)
	default:
		validVersion := false
		for _, v := range api.SupportedVersions() {
			if meta.Version == v {
				validVersion = true
			}
		}
		if !validVersion {
			return fmt.Errorf("invalid version %s, supported values: auto, latest, %s", meta.Version, strings.Join(api.SupportedVersions(), ", "))
		}
	}

	if v := ctl.ControlPlaneVersion(); v == "" {
		return fmt.Errorf("unable to get control plane version")
	} else if meta.Version == "auto" {
		meta.Version = v
		logger.Info("will use version %s for new nodegroup(s) based on control plane version", meta.Version)
	} else if meta.Version != v {
		hint := "--version=auto"
		if clusterConfigFile != "" {
			hint = "metadata.version: auto"
		}
		logger.Warning("will use version %s for new nodegroup(s), while control plane version is %s; to automatically inherit the version use %q", meta.Version, v, hint)
	}

	return nil
}
