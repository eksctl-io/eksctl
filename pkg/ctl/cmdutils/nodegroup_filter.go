package cmdutils

import (
	"github.com/kris-nova/logger"
	"k8s.io/apimachinery/pkg/util/sets"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

// NodeGroupFilter holds filter configuration
type NodeGroupFilter struct {
	*Filter
}

// NewNodeGroupFilter create new NodeGroupFilter instance
func NewNodeGroupFilter() *NodeGroupFilter {
	return &NodeGroupFilter{
		Filter: &Filter{
			ExcludeAll:   false,
			includeNames: sets.NewString(),
			excludeNames: sets.NewString(),
		},
	}
}

// AppendGlobs appends globs for inclusion and exclusion rules
func (f *NodeGroupFilter) AppendGlobs(includeGlobExprs, excludeGlobExprs, ngNames []string) error {
	if err := f.AppendIncludeGlobs(ngNames, includeGlobExprs...); err != nil {
		return err
	}
	return f.AppendExcludeGlobs(excludeGlobExprs...)
}

// AppendIncludeGlobs sets globs for inclusion rules
func (f *NodeGroupFilter) AppendIncludeGlobs(ngNames []string, globExprs ...string) error {
	return f.doAppendIncludeGlobs(ngNames, "nodegroup", globExprs...)
}

// A stackLister lists nodegroup stacks
type stackLister interface {
	ListNodeGroupStacks() ([]manager.NodeGroupStack, error)
}

// SetExcludeExistingFilter uses stackManager to list existing nodegroup stacks and configures
// the filter accordingly
func (f *NodeGroupFilter) SetExcludeExistingFilter(lister stackLister) error {
	if f.ExcludeAll {
		return nil
	}

	existingStacks, err := lister.ListNodeGroupStacks()
	if err != nil {
		return err
	}

	var ngNames []string
	for _, s := range existingStacks {
		ngNames = append(ngNames, s.NodeGroupName)
	}

	return f.doSetExcludeExistingFilter(ngNames, "nodegroup")
}

// SetIncludeOrExcludeMissingFilter uses stackLister to list existing nodegroup stacks and configures
// the filter to either explicitly exclude or include nodegroups that are missing from given nodeGroups
func (f *NodeGroupFilter) SetIncludeOrExcludeMissingFilter(lister stackLister, includeOnlyMissing bool, clusterConfig *api.ClusterConfig) error {
	stacks, err := lister.ListNodeGroupStacks()
	if err != nil {
		return err
	}
	return f.SetIncludeOrExcludeMissingStackFilter(stacks, includeOnlyMissing, clusterConfig)
}

// SetIncludeOrExcludeMissingStackFilter uses a list of existing nodegroup stacks and configures
// the filter to either explicitly exclude or include nodegroups that are missing from given nodeGroups
func (f *NodeGroupFilter) SetIncludeOrExcludeMissingStackFilter(stacks []manager.NodeGroupStack, includeOnlyMissing bool, clusterConfig *api.ClusterConfig) error {
	local := sets.NewString()

	for _, localNodeGroup := range getAllNodeGroupNames(clusterConfig) {
		local.Insert(localNodeGroup)
		if !stackExists(stacks, localNodeGroup) {
			logger.Info("nodegroup %q present in the given config, but missing in the cluster", localNodeGroup)
			f.AppendExcludeNames(localNodeGroup)
		} else if includeOnlyMissing {
			f.AppendExcludeNames(localNodeGroup)
		}
	}

	for _, s := range stacks {
		remoteNodeGroupName := s.NodeGroupName
		if !local.Has(remoteNodeGroupName) {
			logger.Info("nodegroup %q present in the cluster, but missing from the given config", s.NodeGroupName)
			if includeOnlyMissing {
				if s.Type == api.NodeGroupTypeManaged {
					clusterConfig.ManagedNodeGroups = append(clusterConfig.ManagedNodeGroups, &api.ManagedNodeGroup{Name: s.NodeGroupName})
				} else {
					clusterConfig.NodeGroups = append(clusterConfig.NodeGroups, &api.NodeGroup{Name: s.NodeGroupName})
				}
				// make sure it passes it through the filter, so that one can use `--only-missing` along with `--exclude`
				if f.Match(remoteNodeGroupName) {
					f.AppendIncludeNames(remoteNodeGroupName)
				}
			}
		}
	}

	return nil
}

func stackExists(stacks []manager.NodeGroupStack, stackName string) bool {
	for _, s := range stacks {
		if s.NodeGroupName == stackName {
			return true
		}
	}
	return false
}

// LogInfo prints out a user-friendly message about how filter was applied
func (f *NodeGroupFilter) LogInfo(nodeGroups []*api.NodeGroup) {
	f.doLogInfo("nodegroup", f.collectNames(nodeGroups))
}

// MatchAll all names against the filter and return two sets of names - included and excluded
func (f *NodeGroupFilter) MatchAll(nodeGroups []*api.NodeGroup) (sets.String, sets.String) {
	return f.doMatchAll(f.collectNames(nodeGroups))
}

// FilterMatching matches names against the filter and returns all included node groups
func (f *NodeGroupFilter) FilterMatching(nodeGroups []*api.NodeGroup) []*api.NodeGroup {
	var match []*api.NodeGroup
	for _, ng := range nodeGroups {
		if f.Match(ng.NameString()) {
			match = append(match, ng)
		}
	}
	return match
}

// ForEach iterates over each nodegroup that is included by the filter and calls iterFn
func (f *NodeGroupFilter) ForEach(nodeGroups []*api.NodeGroup, iterFn func(i int, ng *api.NodeGroup) error) error {
	for i, ng := range nodeGroups {
		if f.Match(ng.NameString()) {
			if err := iterFn(i, ng); err != nil {
				return err
			}
		}
	}
	return nil
}

func (*NodeGroupFilter) collectNames(nodeGroups []*api.NodeGroup) []string {
	names := []string{}
	for _, ng := range nodeGroups {
		names = append(names, ng.NameString())
	}
	return names
}
