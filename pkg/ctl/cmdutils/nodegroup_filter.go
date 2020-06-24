package cmdutils

import (
	"github.com/kris-nova/logger"
	"k8s.io/apimachinery/pkg/util/sets"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

// NodeGroupFilter holds filter configuration
type NodeGroupFilter struct {
	delegate         *Filter
	onlyLocal        bool
	onlyRemote       bool
	localNodegroups  sets.String
	remoteNodegroups sets.String
}

// NewNodeGroupFilter create new NodeGroupFilter instance
func NewNodeGroupFilter() *NodeGroupFilter {
	return &NodeGroupFilter{
		delegate: &Filter{
			ExcludeAll:   false,
			includeNames: sets.NewString(),
			excludeNames: sets.NewString(),
		},
		localNodegroups:  sets.NewString(),
		remoteNodegroups: sets.NewString(),
	}
}

// AppendGlobs appends globs for inclusion and exclusion rules
func (f *NodeGroupFilter) AppendGlobs(includeGlobExprs, excludeGlobExprs, ngNames []string) error {
	if err := f.AppendIncludeGlobs(ngNames, includeGlobExprs...); err != nil {
		return err
	}
	return f.delegate.AppendExcludeGlobs(excludeGlobExprs...)
}

// AppendIncludeGlobs sets globs for inclusion rules
func (f *NodeGroupFilter) AppendIncludeGlobs(ngNames []string, globExprs ...string) error {
	return f.delegate.doAppendIncludeGlobs(ngNames, "nodegroup", globExprs...)
}

// A stackLister lists nodegroup stacks
type stackLister interface {
	ListNodeGroupStacks() ([]manager.NodeGroupStack, error)
}

// SetOnlyLocal uses stackLister to list existing nodegroup stacks and configures
// the filter to only include the local nodegroups. These are the ones that are present in the clusterconfig file
// but not in the cluster
func (f *NodeGroupFilter) SetOnlyLocal(lister stackLister, clusterConfig *api.ClusterConfig) error {
	f.onlyLocal = true

	return f.loadLocalAndRemoteNodegroups(lister, clusterConfig)
}

// SetOnlyRemote uses stackLister to list existing nodegroup stacks and configures
// the filter to either explicitly exclude or include nodegroups that are missing from given nodeGroups
func (f *NodeGroupFilter) SetOnlyRemote(lister stackLister, clusterConfig *api.ClusterConfig) error {
	f.onlyRemote = true

	return f.loadLocalAndRemoteNodegroups(lister, clusterConfig)

	//for _, localNodeGroup := range getAllNodeGroupNames(clusterConfig) {
	//	local.Insert(localNodeGroup)
	//	if !stackExists(stacks, localNodeGroup) {
	//		logger.Info("nodegroup %q present in the given config, but missing in the cluster", localNodeGroup)
	//		f.AppendExcludeNames(localNodeGroup)
	//	} else if includeOnlyMissing {
	//		f.AppendExcludeNames(localNodeGroup)
	//	}
	//}
	//
	//for _, s := range stacks {
	//	remoteNodeGroupName := s.NodeGroupName
	//	if !local.Has(remoteNodeGroupName) {
	//		logger.Info("nodegroup %q present in the cluster, but missing from the given config", s.NodeGroupName)
	//		if includeOnlyMissing {
	//			if s.Type == api.NodeGroupTypeManaged {
	//				clusterConfig.ManagedNodeGroups = append(clusterConfig.ManagedNodeGroups, &api.ManagedNodeGroup{Name: s.NodeGroupName})
	//			} else {
	//				clusterConfig.NodeGroups = append(clusterConfig.NodeGroups, &api.NodeGroup{Name: s.NodeGroupName})
	//			}
	//			// make sure it passes it through the filter, so that one can use `--only-missing` along with `--exclude`
	//			if f.Match(remoteNodeGroupName) {
	//				f.AppendIncludeNames(remoteNodeGroupName)
	//			}
	//		}
	//	}
	//}
}

func (f *NodeGroupFilter) loadLocalAndRemoteNodegroups(lister stackLister, clusterConfig *api.ClusterConfig) error {

	// Get remote nodegroups
	existingStacks, err := lister.ListNodeGroupStacks()
	if err != nil {
		return err
	}
	for _, s := range existingStacks {
		f.remoteNodegroups.Insert(s.NodeGroupName)
	}

	// Get local nodegroups
	for _, localNodeGroup := range getAllNodeGroupNames(clusterConfig) {
		f.localNodegroups.Insert(localNodeGroup)
		if !stackExists(existingStacks, localNodeGroup) {
			logger.Info("nodegroup %q present in the given config, but missing in the cluster", localNodeGroup)
		}
	}

	// Log remote-only nodegroups  AND add them to the cluster config
	for _, s := range existingStacks {
		remoteNodeGroupName := s.NodeGroupName
		if !f.localNodegroups.Has(remoteNodeGroupName) {
			logger.Info("nodegroup %q present in the cluster, but missing from the given config", s.NodeGroupName)
		} else {
			if s.Type == api.NodeGroupTypeManaged {
				clusterConfig.ManagedNodeGroups = append(clusterConfig.ManagedNodeGroups, &api.ManagedNodeGroup{Name: s.NodeGroupName})
			} else {
				clusterConfig.NodeGroups = append(clusterConfig.NodeGroups, &api.NodeGroup{Name: s.NodeGroupName})
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
	f.delegate.doLogInfo("nodegroup", f.collectNames(nodeGroups))
}

// MatchAll all names against the filter and return two sets of names - included and excluded
func (f *NodeGroupFilter) MatchAll(nodeGroups []*api.NodeGroup) (sets.String, sets.String) {
	names := sets.NewString(f.collectNames(nodeGroups)...)

	if f.onlyLocal {
		names = names.Intersection(f.onlyLocalNodegroups())
		return f.delegate.doMatchAll(names.List())
	}

	if f.onlyRemote {
		names = names.Intersection(f.onlyRemoteNodegroups())
		return f.delegate.doMatchAll(names.List())
	}

	return f.delegate.doMatchAll(names.List())
}

// Match decides whether the given nodegroup is considered included by this filter. It takes into account not only the
// inclusion and exclusion rules (globs) but also the modifiers onlyRemote and onlyLocal.
func (f *NodeGroupFilter) Match(nodeGroup *api.NodeGroup) bool {
	if f.onlyRemote {
		if !f.onlyRemoteNodegroups().Has(nodeGroup.NameString()) {
			return false
		}
		return f.Match(nodeGroup)
	}

	if f.onlyLocal {
		if !f.onlyLocalNodegroups().Has(nodeGroup.NameString()) {
			return false
		}
		return f.Match(nodeGroup)
	}

	return f.Match(nodeGroup)
}

func (f *NodeGroupFilter) onlyLocalNodegroups() sets.String {
	return f.localNodegroups.Difference(f.remoteNodegroups)
}

func (f *NodeGroupFilter) onlyRemoteNodegroups() sets.String {
	return f.remoteNodegroups.Difference(f.localNodegroups)
}

// FilterMatching matches names against the filter and returns all included node groups
func (f *NodeGroupFilter) FilterMatching(nodeGroups []*api.NodeGroup) []*api.NodeGroup {
	var match []*api.NodeGroup
	for _, ng := range nodeGroups {
		if f.Match(ng) {
			match = append(match, ng)
		}
	}
	return match
}

// ForEach iterates over each nodegroup that is included by the filter and calls iterFn
func (f *NodeGroupFilter) ForEach(nodeGroups []*api.NodeGroup, iterFn func(i int, ng *api.NodeGroup) error) error {
	for i, ng := range nodeGroups {
		if f.Match(ng) {
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
