package cmdutils

import (
	"fmt"
	"strings"

	"github.com/gobwas/glob"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	"k8s.io/apimachinery/pkg/util/sets"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

// NodeGroupFilter holds filter configuration
type NodeGroupFilter struct {
	ExcludeAll bool // highest priority

	// include filters take precedence
	includeNames    sets.String
	includeGlobs    []glob.Glob
	rawIncludeGlobs []string

	excludeNames    sets.String
	excludeGlobs    []glob.Glob
	rawExcludeGlobs []string
}

// NewNodeGroupFilter create new NodeGroupFilter instance
func NewNodeGroupFilter() *NodeGroupFilter {
	return &NodeGroupFilter{
		ExcludeAll:   false,
		includeNames: sets.NewString(),
		excludeNames: sets.NewString(),
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
	for _, expr := range globExprs {
		compiledExpr, err := glob.Compile(expr)
		if err != nil {
			return errors.Wrapf(err, "parsing glob filter %q", expr)
		}
		f.includeGlobs = append(f.includeGlobs, compiledExpr)
		f.rawIncludeGlobs = append(f.rawIncludeGlobs, expr)
	}
	return f.includeGlobsMatchAnything(nodeGroups)
}

func (f *NodeGroupFilter) includeGlobsMatchAnything(nodeGroups []*api.NodeGroup) error {
	if len(f.includeGlobs) == 0 {
		return nil
	}
	for _, ng := range nodeGroups {
		if f.matchGlobs(ng.Name, f.includeGlobs) {
			return nil
		}
	}
	return fmt.Errorf("no nodegroups match include glob filter specification: %q", strings.Join(f.rawIncludeGlobs, ","))
}

// AppendIncludeNames appends explicit names to the include filter
func (f *NodeGroupFilter) AppendIncludeNames(names ...string) { f.includeNames.Insert(names...) }

// AppendExcludeGlobs sets globs for exclusion rules
func (f *NodeGroupFilter) AppendExcludeGlobs(globExprs ...string) error {
	for _, expr := range globExprs {
		compiledExpr, err := glob.Compile(expr)
		if err != nil {
			return errors.Wrapf(err, "parsing glob filter %q", expr)
		}
		f.excludeGlobs = append(f.excludeGlobs, compiledExpr)
		f.rawExcludeGlobs = append(f.rawExcludeGlobs, expr)

	}
	return nil // exclude filter doesn't have to match anything, so we don't validate it
}

// AppendExcludeNames appends explicit names to the exclude filter
func (f *NodeGroupFilter) AppendExcludeNames(names ...string) { f.excludeNames.Insert(names...) }

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

	f.excludeNames.Insert(existing...)

	for _, name := range existing {
		isAlsoIncluded := f.includeNames.Has(name)
		if f.matchGlobs(name, f.includeGlobs) {
			isAlsoIncluded = true
		}
		if isAlsoIncluded {
			return fmt.Errorf("existing nodegroup %q should be excluded, but matches include filter: %s", name, f.describeIncludeRules())
		}
	}
	return nil
}

// SetIncludeOrExcludeMissingFilter uses stackManager to list existing nodegroup stacks and configures
// the filter to either explictily exluce or include nodegroups that are missing from given nodeGroups
func (f *NodeGroupFilter) SetIncludeOrExcludeMissingFilter(stackManager *manager.StackCollection, includeOnlyMissing bool, nodeGroups *[]*api.NodeGroup) error {
	stacks, err := stackManager.DescribeNodeGroupStacks()
	if err != nil {
		return err
	}

	remote := sets.NewString()
	local := sets.NewString()

	for _, s := range stacks {
		if name := stackManager.GetNodeGroupName(s); name != "" {
			remote.Insert(name)
		}
	}

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

func (*NodeGroupFilter) matchGlobs(name string, exprs []glob.Glob) bool {
	for _, compiledExpr := range exprs {
		if compiledExpr.Match(name) {
			return true
		}
	}
	return false
}

func (f *NodeGroupFilter) hasIncludeRules() bool {
	return f.includeNames.Len()+len(f.includeGlobs) != 0
}

func (f *NodeGroupFilter) describeIncludeRules() string {
	rules := append(f.includeNames.List(), f.rawIncludeGlobs...)
	return fmt.Sprintf("%s", strings.Join(rules, ","))
}

func (f *NodeGroupFilter) hasExcludeRules() bool {
	return f.excludeNames.Len()+len(f.excludeGlobs) != 0
}

func (f *NodeGroupFilter) describeExcludeRules() string {
	rules := append(f.excludeNames.List(), f.rawExcludeGlobs...)
	return fmt.Sprintf("%s", strings.Join(rules, ","))
}

// Match given nodegroup against the filter and returns
// true or false if it has to be included or excluded
func (f *NodeGroupFilter) Match(name string) bool {
	if f.ExcludeAll {
		return false // force exclude
	}

	hasIncludeRules := f.hasIncludeRules()
	hasExcludeRules := f.hasExcludeRules()

	if !hasIncludeRules && !hasExcludeRules {
		return true // empty rules - include
	}

	mustInclude := false // use this override when rules overlap

	if hasIncludeRules {
		mustInclude = f.includeNames.Has(name)
		if f.matchGlobs(name, f.includeGlobs) {
			mustInclude = true
		}
		if !hasExcludeRules {
			// empty exclusion rules - explicit inclusion mode
			return mustInclude
		}
	}

	if hasExcludeRules {
		exclude := f.excludeNames.Has(name)
		if f.matchGlobs(name, f.excludeGlobs) {
			exclude = true
		}
		if exclude && !mustInclude {
			// exclude, unless overridden by an inclusion rule
			return false
		}
	}

	return true // biased to include
}

// MatchAll nodegroups against the filter and return two sets of names - included and excluded
func (f *NodeGroupFilter) MatchAll(nodeGroups []*api.NodeGroup) (sets.String, sets.String) {
	included, excluded := sets.NewString(), sets.NewString()
	if f.ExcludeAll {
		for _, ng := range nodeGroups {
			excluded.Insert(ng.Name)
		}
		return included, excluded
	}
	for _, ng := range nodeGroups {
		if f.Match(ng.Name) {
			included.Insert(ng.Name)
		} else {
			excluded.Insert(ng.Name)
		}
	}
	return included, excluded
}

// LogInfo prints out a user-friendly message about how filter was applied
func (f *NodeGroupFilter) LogInfo(nodeGroups []*api.NodeGroup) {
	logMsg := func(ngSubset sets.String, status string) {
		count := ngSubset.Len()
		list := strings.Join(ngSubset.List(), ", ")
		subject := "nodegroups (%s) were"
		if count == 1 {
			subject = "nodegroup (%s) was"
		}
		logger.Info("%d "+subject+" %s", count, list, status)
	}

	included, excluded := f.MatchAll(nodeGroups)
	if f.hasIncludeRules() {
		logger.Info("include rules: %s", f.describeIncludeRules())
		if included.Len() == 0 {
			logger.Info("no nogroups were included by the filter")
		}
	}
	if included.Len() > 0 {
		logMsg(included, "included")
	}
	if f.hasExcludeRules() {
		logger.Info("exclude rules: %s", f.describeExcludeRules())
		if excluded.Len() == 0 {
			logger.Info("no nogroups were excluded by the filter")
		}
	}
	if excluded.Len() > 0 {
		logMsg(excluded, "excluded")
	}
}

// ForEach iterates over each nodegroup that is included by the filter and calls iterFn
func (f *NodeGroupFilter) ForEach(nodeGroups []*api.NodeGroup, iterFn func(i int, ng *api.NodeGroup) error) error {
	for i, ng := range nodeGroups {
		if f.Match(ng.Name) {
			if err := iterFn(i, ng); err != nil {
				return err
			}
		}
	}
	return nil
}

// ValidateNodeGroupsAndSetDefaults is calls api.ValidateNodeGroup & api.SetNodeGroupDefaults for
// all nodegroups that match the filter
func (f *NodeGroupFilter) ValidateNodeGroupsAndSetDefaults(nodeGroups []*api.NodeGroup) error {
	return f.ForEach(nodeGroups, func(i int, ng *api.NodeGroup) error {
		if err := api.ValidateNodeGroup(i, ng); err != nil {
			return err
		}
		if err := api.SetNodeGroupDefaults(i, ng); err != nil {
			return err
		}
		return nil
	})
}
