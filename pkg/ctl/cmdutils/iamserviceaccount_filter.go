package cmdutils

import (
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/kris-nova/logger"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

// IAMServiceAccountFilter holds filter configuration
type IAMServiceAccountFilter struct {
	*Filter
}

// NewIAMServiceAccountFilter create new ServiceAccountFilter instance
func NewIAMServiceAccountFilter() *IAMServiceAccountFilter {
	return &IAMServiceAccountFilter{
		Filter: &Filter{
			ExcludeAll:   false,
			includeNames: sets.NewString(),
			excludeNames: sets.NewString(),
		},
	}
}

// AppendGlobs appends globs for inclusion and exclusion rules
func (f *IAMServiceAccountFilter) AppendGlobs(includeGlobExprs, excludeGlobExprs []string, serviceAccounts []*api.ClusterIAMServiceAccount) error {
	if err := f.AppendIncludeGlobs(serviceAccounts, includeGlobExprs...); err != nil {
		return err
	}
	return f.AppendExcludeGlobs(excludeGlobExprs...)
}

// AppendIncludeGlobs sets globs for inclusion rules
func (f *IAMServiceAccountFilter) AppendIncludeGlobs(serviceAccounts []*api.ClusterIAMServiceAccount, globExprs ...string) error {
	return f.doAppendIncludeGlobs(f.collectNames(serviceAccounts), "iamserviceaccount", globExprs...)
}

// SetExcludeExistingFilter uses stackManager to list existing nodegroup stacks and configures
// the filter accordingly
func (f *IAMServiceAccountFilter) SetExcludeExistingFilter(stackManager *manager.StackCollection) error {
	if f.ExcludeAll {
		return nil
	}

	existing, err := stackManager.ListIAMServiceAccountStacks()
	if err != nil {
		return err
	}

	return f.doSetExcludeExistingFilter(existing, "iamserviceaccount")
}

// SetIncludeOrExcludeMissingFilter uses stackManager to list existing iamserviceaccount stacks and configures
// the filter to either explictily exluce or include iamserviceaccounts that are missing from given serviceAccounts
func (f *IAMServiceAccountFilter) SetIncludeOrExcludeMissingFilter(stackManager *manager.StackCollection, includeOnlyMissing bool, serviceAccounts *[]*api.ClusterIAMServiceAccount) error {
	existing, err := stackManager.ListIAMServiceAccountStacks()
	if err != nil {
		return err
	}

	remote := sets.NewString(existing...)
	local := sets.NewString()

	for _, localServiceAccount := range *serviceAccounts {
		localServiceAccountName := localServiceAccount.NameString()
		local.Insert(localServiceAccountName)
		if !remote.Has(localServiceAccountName) {
			logger.Info("iamserviceaccounts %q present in the given config, but missing in the cluster", localServiceAccountName)
			f.AppendExcludeNames(localServiceAccountName)
		} else if includeOnlyMissing {
			f.AppendExcludeNames(localServiceAccountName)
		}
	}

	for remoteServiceAccountName := range remote {
		if !local.Has(remoteServiceAccountName) {
			logger.Info("iamserviceaccounts %q present in the cluster, but missing from the given config", remoteServiceAccountName)
			if includeOnlyMissing {
				// append it to the config object, so that `saFilter.ForEach` knows about it
				meta, err := api.ClusterIAMServiceAccountNameStringToObjectMeta(remoteServiceAccountName)
				if err != nil {
					return err
				}
				remoteServiceAccount := &api.ClusterIAMServiceAccount{
					ObjectMeta: *meta,
				}
				*serviceAccounts = append(*serviceAccounts, remoteServiceAccount)
				// make sure it passes it through the filter, so that one can use `--only-missing` along with `--exclude`
				if f.Match(remoteServiceAccountName) {
					f.AppendIncludeNames(remoteServiceAccountName)
				}
			}
		}
	}

	return nil
}

// LogInfo prints out a user-friendly message about how filter was applied
func (f *IAMServiceAccountFilter) LogInfo(serviceAccounts []*api.ClusterIAMServiceAccount) {
	f.doLogInfo("iamserviceaccount", f.collectNames(serviceAccounts))
}

// MatchAll all names against the filter and return two sets of names - included and excluded
func (f *IAMServiceAccountFilter) MatchAll(serviceAccounts []*api.ClusterIAMServiceAccount) (sets.String, sets.String) {
	return f.doMatchAll(f.collectNames(serviceAccounts))
}

// ForEach iterates over each nodegroup that is included by the filter and calls iterFn
func (f *IAMServiceAccountFilter) ForEach(serviceAccounts []*api.ClusterIAMServiceAccount, iterFn func(i int, sa *api.ClusterIAMServiceAccount) error) error {
	for i, sa := range serviceAccounts {
		if f.Match(sa.NameString()) {
			if err := iterFn(i, sa); err != nil {
				return err
			}
		}
	}
	return nil
}

func (*IAMServiceAccountFilter) collectNames(serviceAccounts []*api.ClusterIAMServiceAccount) []string {
	names := []string{}
	for _, sa := range serviceAccounts {
		names = append(names, sa.NameString())
	}
	return names
}
