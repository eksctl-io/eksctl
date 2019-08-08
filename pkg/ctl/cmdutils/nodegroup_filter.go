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
func (f *NodeGroupFilter) AppendGlobs(includeGlobExprs, excludeGlobExprs []string, nodeGroups []*api.NodeGroup) error {
	if err := f.AppendIncludeGlobs(nodeGroups, includeGlobExprs...); err != nil {
		return err
	}
	return f.AppendExcludeGlobs(excludeGlobExprs...)
}

// AppendIncludeGlobs sets globs for inclusion rules
func (f *NodeGroupFilter) AppendIncludeGlobs(nodeGroups []*api.NodeGroup, globExprs ...string) error {
	return f.doAppendIncludeGlobs(f.collectNames(nodeGroups), "nodegroup", globExprs...)
}

// SetExcludeExistingFilter uses stackManager to list existing nodegroup stacks and configures
// the filter accordingly
func (f *NodeGroupFilter) SetExcludeExistingFilter(stackManager *manager.StackCollection) error {
	if f.ExcludeAll {
		return nil
	}

	existing, err := stackManager.ListNodeGroupStacks()
	if err != nil {
		return err
	}

	return f.doSetExcludeExistingFilter(existing, "nodegroup")
}

// SetIncludeOrExcludeMissingFilter uses stackManager to list existing nodegroup stacks and configures
// the filter to either explictily exluce or include nodegroups that are missing from given nodeGroups
func (f *NodeGroupFilter) SetIncludeOrExcludeMissingFilter(stackManager *manager.StackCollection, includeOnlyMissing bool, nodeGroups *[]*api.NodeGroup) error {
	existing, err := stackManager.ListNodeGroupStacks()
	if err != nil {
		return err
	}

	remote := sets.NewString(existing...)
	local := sets.NewString()

	for _, localNodeGroup := range *nodeGroups {
		local.Insert(localNodeGroup.Name)
		if !remote.Has(localNodeGroup.Name) {
			logger.Info("nodegroup %q present in the given config, but missing in the cluster", localNodeGroup.Name)
			f.AppendExcludeNames(localNodeGroup.Name)
		} else if includeOnlyMissing {
			f.AppendExcludeNames(localNodeGroup.Name)
		}
	}

	for remoteNodeGroupName := range remote {
		if !local.Has(remoteNodeGroupName) {
			logger.Info("nodegroup %q present in the cluster, but missing from the given config", remoteNodeGroupName)
			if includeOnlyMissing {
				// append it to the config object, so that `ngFilter.ForEach` knows about it
				*nodeGroups = append(*nodeGroups, &api.NodeGroup{Name: remoteNodeGroupName})
				// make sure it passes it through the filter, so that one can use `--only-missing` along with `--exclude`
				if f.Match(remoteNodeGroupName) {
					f.AppendIncludeNames(remoteNodeGroupName)
				}
			}
		}
	}

	return nil
}

// LogInfo prints out a user-friendly message about how filter was applied
func (f *NodeGroupFilter) LogInfo(nodeGroups []*api.NodeGroup) {
	f.doLogInfo("nodegroup", f.collectNames(nodeGroups))
}

// MatchAll all names against the filter and return two sets of names - included and excluded
func (f *NodeGroupFilter) MatchAll(nodeGroups []*api.NodeGroup) (sets.String, sets.String) {
	return f.doMatchAll(f.collectNames(nodeGroups))
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
