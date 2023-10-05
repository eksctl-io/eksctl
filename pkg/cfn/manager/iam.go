package manager

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
)

// DescribeIAMServiceAccountStacks calls ListStacks and filters out iamserviceaccounts
func (c *StackCollection) DescribeIAMServiceAccountStacks(ctx context.Context) ([]*Stack, error) {
	stacks, err := c.ListStacks(ctx)
	if err != nil {
		return nil, err
	}

	iamServiceAccountStacks := []*Stack{}
	for _, s := range stacks {
		if s.StackStatus == types.StackStatusDeleteComplete {
			continue
		}
		if GetIAMServiceAccountName(s) != "" {
			iamServiceAccountStacks = append(iamServiceAccountStacks, s)
		}
	}
	logger.Debug("iamserviceaccounts = %v", iamServiceAccountStacks)
	return iamServiceAccountStacks, nil
}

// ListIAMServiceAccountStacks calls DescribeIAMServiceAccountStacks and returns only iamserviceaccount names
func (c *StackCollection) ListIAMServiceAccountStacks(ctx context.Context) ([]string, error) {
	stacks, err := c.DescribeIAMServiceAccountStacks(ctx)
	if err != nil {
		return nil, err
	}

	names := []string{}
	for _, s := range stacks {
		names = append(names, GetIAMServiceAccountName(s))
	}
	return names, nil
}

// GetIAMServiceAccounts calls DescribeIAMServiceAccountStacks and return native iamserviceaccounts
func (c *StackCollection) GetIAMServiceAccounts(ctx context.Context) ([]*api.ClusterIAMServiceAccount, error) {
	stacks, err := c.DescribeIAMServiceAccountStacks(ctx)
	if err != nil {
		return nil, err
	}

	results := []*api.ClusterIAMServiceAccount{}
	for _, s := range stacks {
		meta, err := api.ClusterIAMServiceAccountNameStringToClusterIAMMeta(GetIAMServiceAccountName(s))
		if err != nil {
			return nil, err
		}
		serviceAccount := &api.ClusterIAMServiceAccount{
			ClusterIAMMeta: *meta,
			Status:         &api.ClusterIAMServiceAccountStatus{},
		}

		// TODO: we need to make it easier to fetch full definition of the object,
		// namely: all label, full role definition; we can do that by caching
		// the ClusterConfig time we make an update and a mechanism of validating
		// whether it is up to date;
		// otherwise we could extend this with tedious calls to each of the API,
		// but it's not very feasible and it's best ot create a general solution
		outputCollectors := outputs.NewCollectorSet(map[string]outputs.Collector{
			outputs.IAMServiceAccountRoleName: func(v string) error {
				serviceAccount.Status.RoleARN = &v
				return nil
			},
		})

		if err := outputCollectors.MustCollect(*s); err != nil {
			return nil, err
		}

		results = append(results, serviceAccount)
	}
	return results, nil
}

// GetIAMServiceAccountName will return iamserviceaccount name based on tags
func GetIAMServiceAccountName(s *Stack) string {
	for _, tag := range s.Tags {
		if *tag.Key == api.IAMServiceAccountNameTag {
			return *tag.Value
		}
	}
	return ""
}

func (c *StackCollection) GetIAMAddonsStacks(ctx context.Context) ([]*Stack, error) {
	stacks, err := c.ListStacks(ctx)
	if err != nil {
		return nil, err
	}

	iamAddonStacks := []*Stack{}
	for _, s := range stacks {
		if s.StackStatus == types.StackStatusDeleteComplete {
			continue
		}
		if c.GetIAMAddonName(s) != "" {
			iamAddonStacks = append(iamAddonStacks, s)
		}
	}
	logger.Debug("iamserviceaccounts = %v", iamAddonStacks)
	return iamAddonStacks, nil
}

func (*StackCollection) GetIAMAddonName(s *Stack) string {
	for _, tag := range s.Tags {
		if *tag.Key == api.AddonNameTag {
			return *tag.Value
		}
	}
	return ""
}
