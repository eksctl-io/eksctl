package filter

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/weaveworks/eksctl/pkg/actions/accessentry"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// An AccessEntryLister lists access entries.
//
//counterfeiter:generate . AccessEntryLister
type AccessEntryLister interface {
	// ListAccessEntryStackNames lists stack names for all access entries in the specified cluster.
	ListAccessEntryStackNames(ctx context.Context, clusterName string) ([]string, error)
}

// AccessEntry filters out existing access entry resources.
type AccessEntry struct {
	Lister      AccessEntryLister
	ClusterName string
}

// FilterOutExistingStacks returns a set of api.AccessEntry resources that do not have a corresponding stack.
func (a *AccessEntry) FilterOutExistingStacks(ctx context.Context, accessEntries []api.AccessEntry) ([]api.AccessEntry, error) {
	stackNames, err := a.Lister.ListAccessEntryStackNames(ctx, a.ClusterName)
	if err != nil {
		return nil, fmt.Errorf("error listing access entry stacks: %w", err)
	}

	existingStacks := sets.NewString(stackNames...)
	var filtered []api.AccessEntry
	for _, ae := range accessEntries {
		stackName := accessentry.MakeStackName(a.ClusterName, ae)
		if !existingStacks.Has(stackName) {
			filtered = append(filtered, ae)
		}
	}
	return filtered, nil
}
