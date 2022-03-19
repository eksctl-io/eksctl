package manager

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
)

// makeIAMServiceAccountStackName generates the name of the iamserviceaccount stack identified by its name, isolated by the cluster this StackCollection operates on and 'addon' suffix
func (c *StackCollection) makeIAMServiceAccountStackName(namespace, name string) string {
	return fmt.Sprintf("eksctl-%s-addon-iamserviceaccount-%s-%s", c.spec.Metadata.Name, namespace, name)
}

// stackHasRolledBack alerts of existing stack in rollback status
func (c *StackCollection) stackHasRolledBack(stackName string) (*Stack, error) {
	input := &cfn.DescribeStacksInput{
		StackName: &stackName,
	}
	resp, err := c.cloudformationAPI.DescribeStacks(input)
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if !ok {
			return nil, errors.Wrapf(err, "conversion to an AWS error failed")
		}
		if len(aerr.Code()) == 0 {
			return nil, err
		}
		if aerr.Code() != "ValidationError" {
			return nil, errors.Wrapf(err, "describing CloudFormation stack %q, code %q", stackName, aerr.Code())
		}
	}
	if len(resp.Stacks) == 0 {
		return nil, nil
	}
	for _, s := range resp.Stacks {
		if *(s.StackStatus) == cfn.StackStatusRollbackComplete {
			return s, nil
		}
		if !c.StackStatusIsNotTransitional(s) {
			return nil, errors.Wrapf(err, "stack %q is in a transitional status (%q)", stackName, *(s.StackStatus))
		}
	}
	return nil, nil
}

func (c *StackCollection) deleteRolledbackStack(name string) error {
	rollbackedStack, err := c.stackHasRolledBack(name)
	if err != nil {
		return err
	}
	if rollbackedStack != nil {
		logger.Warning("deleting existing rolled back stack %q", name)
		err = c.DeleteStackSync(rollbackedStack)
		if err != nil {
			return err
		}
	}
	return nil
}

// createIAMServiceAccountTask creates the iamserviceaccount in CloudFormation
func (c *StackCollection) createIAMServiceAccountTask(errs chan error, spec *api.ClusterIAMServiceAccount, oidc *iamoidc.OpenIDConnectManager) error {
	name := c.makeIAMServiceAccountStackName(spec.Namespace, spec.Name)
	logger.Info("building iamserviceaccount stack %q", name)
	stack := builder.NewIAMRoleResourceSetForServiceAccount(spec, oidc)
	if err := stack.AddAllResources(); err != nil {
		return err
	}

	if spec.Tags == nil {
		spec.Tags = make(map[string]string)
	}
	spec.Tags[api.IAMServiceAccountNameTag] = spec.NameString()

	if err := c.deleteRolledbackStack(name); err != nil {
		return err
	}
	if err := c.CreateStack(name, stack, spec.Tags, nil, errs); err != nil {
		logger.Info("an error occurred creating the stack, to cleanup resources, run 'eksctl delete iamserviceaccount --region=%s --name=%s --namespace=%s'", c.spec.Metadata.Region, spec.Name, spec.Namespace)
		return err
	}
	return nil
}

// DescribeIAMServiceAccountStacks calls DescribeStacks and filters out iamserviceaccounts
func (c *StackCollection) DescribeIAMServiceAccountStacks() ([]*Stack, error) {
	stacks, err := c.DescribeStacks()
	if err != nil {
		return nil, err
	}

	iamServiceAccountStacks := []*Stack{}
	for _, s := range stacks {
		if *s.StackStatus == cfn.StackStatusDeleteComplete {
			continue
		}
		if *s.StackStatus == cfn.StackStatusRollbackComplete {
			logger.Warning("found stack %v in ROLLBACK_COMPLETE", *s.StackName)
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
func (c *StackCollection) ListIAMServiceAccountStacks() ([]string, error) {
	stacks, err := c.DescribeIAMServiceAccountStacks()
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
func (c *StackCollection) GetIAMServiceAccounts() ([]*api.ClusterIAMServiceAccount, error) {
	stacks, err := c.DescribeIAMServiceAccountStacks()
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

func (c *StackCollection) GetIAMAddonsStacks() ([]*Stack, error) {
	stacks, err := c.DescribeStacks()
	if err != nil {
		return nil, err
	}

	iamAddonStacks := []*Stack{}
	for _, s := range stacks {
		if *s.StackStatus == cfn.StackStatusDeleteComplete {
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
