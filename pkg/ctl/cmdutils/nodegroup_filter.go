package cmdutils

import (
	"fmt"
	"strings"

	"github.com/gobwas/glob"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	"k8s.io/apimachinery/pkg/util/sets"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

// NodeGroupFilter holds filter configuration
type NodeGroupFilter struct {
	SkipAll           bool
	IgnoreAllExisting bool

	existing sets.String
	only     []glob.Glob
	onlySpec string
}

// NewNodeGroupFilter create new NodeGroupFilter instance
func NewNodeGroupFilter() *NodeGroupFilter {
	return &NodeGroupFilter{
		IgnoreAllExisting: true,
		SkipAll:           false,
		existing:          sets.NewString(),
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
	for _, ng := range cfg.NodeGroups {
		for _, compiledExpr := range f.only {
			if compiledExpr.Match(ng.Name) {
				return nil
			}
		}
	}
	return fmt.Errorf("no nodegroups match filter specification: %q", f.onlySpec)
}

// ApplyExistingFilter uses stackManager to list existing nodegroup stacks and configures
// the filter accordingly
func (f *NodeGroupFilter) ApplyExistingFilter(stackManager *manager.StackCollection) error {
	if !f.IgnoreAllExisting || f.SkipAll {
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
	if f.SkipAll {
		return false
	}
	if f.IgnoreAllExisting && f.existing.Has(ng.Name) {
		return false
	}

	for _, compiledExpr := range f.only {
		if compiledExpr.Match(ng.Name) {
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
	if f.SkipAll {
		return names
	}
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
		if f.SkipAll {
			reasons = append(reasons, fmt.Sprintf("--without-nodegroup was set"))
		}
		if f.onlySpec != "" {
			reasons = append(reasons, fmt.Sprintf("--only=%q was given", f.onlySpec))
		}
		if existingCount := f.existing.Len(); existingCount > 0 {
			reasons = append(reasons, fmt.Sprintf("%d nodegroup(s) (%s) already exist", existingCount, strings.Join(f.existing.List(), ", ")))
		}
		logger.Info("%d nodegroup(s) were filtered out: %s", filteredOutCount, strings.Join(reasons, ", "))
	}
}

// CheckEachNodeGroup iterates over each nodegroup and calls check function if it matches the filter
func (f *NodeGroupFilter) CheckEachNodeGroup(nodeGroups []*api.NodeGroup, check func(i int, ng *api.NodeGroup) error) error {
	for i, ng := range nodeGroups {
		if f.Match(ng) {
			if err := check(i, ng); err != nil {
				return err
			}
		}
	}
	return nil
}

// ValidateNodeGroupsAndSetDefaults is calls api.ValidateNodeGroup & api.SetNodeGroupDefaults for
// all nodegroups that match the filter
func (f *NodeGroupFilter) ValidateNodeGroupsAndSetDefaults(nodeGroups []*api.NodeGroup) error {
	return f.CheckEachNodeGroup(nodeGroups, func(i int, ng *api.NodeGroup) error {
		if err := api.ValidateNodeGroup(i, ng); err != nil {
			return err
		}
		if err := api.SetNodeGroupDefaults(i, ng); err != nil {
			return err
		}
		return nil
	})
}
