package filter

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/kris-nova/logger"
	"k8s.io/apimachinery/pkg/util/sets"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate . NodegroupFilter
// NodegroupFilter is an interface that holds filter configuration
type NodegroupFilter interface {
	SetOnlyLocal(ctx context.Context, eksAPI awsapi.EKS, lister StackLister, clusterConfig *api.ClusterConfig) error
	Match(ngName string) bool
	LogInfo(cfg *api.ClusterConfig)
}

// StackLister lists nodegroup stacks
type StackLister interface {
	ListNodeGroupStacks(ctx context.Context) ([]manager.NodeGroupStack, error)
}

// NodeGroupFilter holds filter configuration
type NodeGroupFilter struct {
	delegate         *Filter
	onlyLocal        bool
	onlyRemote       bool
	localNodegroups  sets.String
	remoteNodegroups sets.String
}

// NewNodeGroupFilter creates a new NodeGroupFilter struct
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
	return f.AppendExcludeGlobs(excludeGlobExprs...)
}

// AppendIncludeGlobs sets globs for inclusion rules
func (f *NodeGroupFilter) AppendIncludeGlobs(ngNames []string, globExprs ...string) error {
	return f.delegate.doAppendIncludeGlobs(ngNames, "nodegroup", globExprs...)
}

// AppendExcludeGlobs sets globs for inclusion rules
func (f *NodeGroupFilter) AppendExcludeGlobs(globExprs ...string) error {
	return f.delegate.AppendExcludeGlobs(globExprs...)
}

// AppendIncludeNames sets globs for inclusion rules
func (f *NodeGroupFilter) AppendIncludeNames(names ...string) {
	f.delegate.AppendIncludeNames(names...)
}

// SetOnlyLocal uses StackLister to list existing nodegroup stacks and configures
// the filter to only include the nodegroups that don't exist in the cluster already.
// Note: they are present in the config file but not in the cluster. This is used by
// the create nodegroup command
func (f *NodeGroupFilter) SetOnlyLocal(ctx context.Context, eksAPI awsapi.EKS, lister StackLister, clusterConfig *api.ClusterConfig) error {
	f.onlyLocal = true

	err := f.loadLocalAndRemoteNodegroups(ctx, eksAPI, lister, clusterConfig)
	if err != nil {
		return err
	}

	// Remote ones will be excluded
	if f.remoteNodegroups.Len() > 0 {
		logger.Info("%d existing %s(s) (%s) will be excluded", f.remoteNodegroups.Len(), "nodegroup", strings.Join(f.remoteNodegroups.List(), ","))
	}
	return nil
}

// SetOnlyRemote uses StackLister to list existing nodegroup stacks and configures
// the filter to exclude nodegroups already defined in the config file. It will include the
//  nodegroups that exist in the cluster but not in the config
func (f *NodeGroupFilter) SetOnlyRemote(ctx context.Context, eksAPI awsapi.EKS, lister StackLister, clusterConfig *api.ClusterConfig) error {
	f.onlyRemote = true

	err := f.loadLocalAndRemoteNodegroups(ctx, eksAPI, lister, clusterConfig)
	if err != nil {
		return err
	}

	// local ones will be excluded
	if f.localNodegroups.Len() > 0 {
		logger.Info("%d %s(s) present in the config file (%s) will be excluded", f.localNodegroups.Len(), "nodegroup", strings.Join(f.localNodegroups.List(), ","))
	}
	return nil
}

// SetExcludeAll sets the ExcludeAll flag in the filter so that no nodegroups are matched
func (f *NodeGroupFilter) SetExcludeAll(excludeAll bool) {
	f.delegate.ExcludeAll = excludeAll
}

// GetExcludeAll returns whether all nodegroups will be excluded
func (f *NodeGroupFilter) GetExcludeAll() bool {
	return f.delegate.ExcludeAll
}

func (f *NodeGroupFilter) loadLocalAndRemoteNodegroups(ctx context.Context, eksAPI awsapi.EKS, lister StackLister, clusterConfig *api.ClusterConfig) error {
	nodeGroupsWithStacks, nodeGroupsWithoutStacks, err := f.findAllNodeGroups(ctx, eksAPI, lister, clusterConfig)
	if err != nil {
		return err
	}
	for _, s := range nodeGroupsWithStacks {
		f.remoteNodegroups.Insert(s.NodeGroupName)
	}

	for _, s := range nodeGroupsWithoutStacks {
		f.remoteNodegroups.Insert(s)
	}

	f.localNodegroups.Insert(clusterConfig.GetAllNodeGroupNames()...)

	// Get local nodegroups and discover if any specified don't have stacks
	for _, localNodeGroup := range clusterConfig.GetAllNodeGroupNames() {
		if !stackExists(nodeGroupsWithStacks, localNodeGroup) && !nodeExists(nodeGroupsWithoutStacks, localNodeGroup) {
			logger.Debug("nodegroup %q present in the given config, but missing in the cluster", localNodeGroup)
		}
	}

	for _, nodeGroupWithoutStack := range nodeGroupsWithoutStacks {
		if !f.localNodegroups.Has(nodeGroupWithoutStack) {
			ngBase := &api.NodeGroupBase{Name: nodeGroupWithoutStack}
			logger.Debug("nodegroup %q present in the cluster, but missing from the given config", nodeGroupWithoutStack)
			clusterConfig.ManagedNodeGroups = append(clusterConfig.ManagedNodeGroups, &api.ManagedNodeGroup{NodeGroupBase: ngBase})
		}
	}

	// Log remote-only nodegroups  AND add them to the cluster config
	for _, s := range nodeGroupsWithStacks {
		remoteNodeGroupName := s.NodeGroupName
		if !f.localNodegroups.Has(remoteNodeGroupName) {
			ngBase := &api.NodeGroupBase{Name: s.NodeGroupName}
			logger.Debug("nodegroup %q present in the cluster, but missing from the given config", s.NodeGroupName)
			if s.Type == api.NodeGroupTypeManaged {
				clusterConfig.ManagedNodeGroups = append(clusterConfig.ManagedNodeGroups, &api.ManagedNodeGroup{NodeGroupBase: ngBase})
			} else {
				clusterConfig.NodeGroups = append(clusterConfig.NodeGroups, &api.NodeGroup{NodeGroupBase: ngBase})
			}
		}
	}
	return nil
}

func (f *NodeGroupFilter) findAllNodeGroups(ctx context.Context, eksAPI awsapi.EKS, lister StackLister, clusterConfig *api.ClusterConfig) ([]manager.NodeGroupStack, []string, error) {
	// Get remote nodegroups stacks
	nodeGroupsWithStacks, err := lister.ListNodeGroupStacks(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Get all nodegroups
	allNodeGroups, err := eksAPI.ListNodegroups(ctx, &eks.ListNodegroupsInput{
		ClusterName: &clusterConfig.Metadata.Name,
	})
	if err != nil {
		return nil, nil, err
	}

	nodeGroupsWithStacksSet := sets.NewString()
	for _, s := range nodeGroupsWithStacks {
		nodeGroupsWithStacksSet.Insert(s.NodeGroupName)
	}

	var nodeGroupsWithoutStacks []string
	for _, nodeGroupName := range allNodeGroups.Nodegroups {
		if !nodeGroupsWithStacksSet.Has(nodeGroupName) {
			nodeGroupsWithoutStacks = append(nodeGroupsWithoutStacks, nodeGroupName)
		}
	}

	return nodeGroupsWithStacks, nodeGroupsWithoutStacks, nil
}

func nodeExists(stacks []string, stackName string) bool {
	for _, s := range stacks {
		if s == stackName {
			return true
		}
	}
	return false
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
func (f *NodeGroupFilter) LogInfo(cfg *api.ClusterConfig) {
	allNames := f.collectNames(cfg.NodeGroups)

	for _, mng := range cfg.ManagedNodeGroups {
		allNames.Insert(mng.NameString())
	}

	included, excluded := f.matchAll(allNames)
	f.delegate.doLogInfo("nodegroup", included, excluded)
}

// matchAll all names against the filter and return two sets of names - included and excluded
func (f *NodeGroupFilter) matchAll(allNames sets.String) (sets.String, sets.String) {

	matching, notMatching := f.delegate.doMatchAll(allNames.List())

	if f.onlyLocal {
		// From the ones that match, pick only the local ones
		included := matching.Intersection(f.onlyLocalNodegroups())
		excluded := allNames.Difference(included)
		return included, excluded
	}

	if f.onlyRemote {
		// From the ones that match, pick only the remote ones
		included := matching.Intersection(f.onlyRemoteNodegroups())
		excluded := allNames.Difference(included)
		return included, excluded
	}

	return matching, notMatching
}

// Match decides whether the given nodegroup is considered included by this filter. It takes into account not only the
// inclusion and exclusion rules (globs) but also the modifiers onlyRemote and onlyLocal.
func (f *NodeGroupFilter) Match(ngName string) bool {
	if f.onlyRemote {
		if !f.onlyRemoteNodegroups().Has(ngName) {
			return false
		}
		return f.delegate.Match(ngName)
	}

	if f.onlyLocal {
		if !f.onlyLocalNodegroups().Has(ngName) {
			return false
		}
		return f.delegate.Match(ngName)
	}

	return f.delegate.Match(ngName)
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

func (*NodeGroupFilter) collectNames(nodeGroups []*api.NodeGroup) sets.String {
	names := sets.NewString()
	for _, ng := range nodeGroups {
		names.Insert(ng.NameString())
	}
	return names
}
