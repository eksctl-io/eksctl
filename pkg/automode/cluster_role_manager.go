package automode

import (
	"context"
	_ "embed"
	"fmt"
	"slices"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"

	"github.com/kris-nova/logger"

	"goformation/v4"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

//go:embed assets/eks-auto-mode-policy-document.json
var eksAutoModePolicyDocument string

//go:embed assets/eks-policy-document.json
var eksPolicyDocument string

type IAMRoleManager interface {
	AttachRolePolicy(ctx context.Context, params *iam.AttachRolePolicyInput, optFns ...func(options *iam.Options)) (*iam.AttachRolePolicyOutput, error)
	DetachRolePolicy(ctx context.Context, params *iam.DetachRolePolicyInput, optFns ...func(options *iam.Options)) (*iam.DetachRolePolicyOutput, error)
}

type ClusterRoleManager struct {
	IAMRoleManager  awsapi.IAM
	StackManager    manager.StackManager
	Region          string
	ClusterRoleName string
}

func (c *ClusterRoleManager) UpdateRoleForAutoMode(ctx context.Context) error {
	hasDedicatedClusterRole, err := c.hasDedicatedClusterRole(ctx)
	if err != nil {
		return err
	}
	if !hasDedicatedClusterRole {
		return nil
	}
	existingRolePolicies, err := c.listExistingRolePolicies(ctx)
	if err != nil {
		return err
	}
	for _, p := range builder.AutoModeIAMPolicies {
		if slices.Contains(existingRolePolicies, p) {
			continue
		}
		if _, err := c.IAMRoleManager.AttachRolePolicy(ctx, &iam.AttachRolePolicyInput{
			PolicyArn: aws.String(makePolicyARN(p, c.Region)),
			RoleName:  aws.String(c.ClusterRoleName),
		}); err != nil {
			return fmt.Errorf("attaching role policy %q: %w", p, err)
		}
	}
	return c.updateAssumeRolePolicy(ctx, eksAutoModePolicyDocument)
}

func (c *ClusterRoleManager) DeleteAutoModePolicies(ctx context.Context) error {
	hasDedicatedClusterRole, err := c.hasDedicatedClusterRole(ctx)
	if err != nil {
		return err
	}
	if !hasDedicatedClusterRole {
		return nil
	}
	existingRolePolicies, err := c.listExistingRolePolicies(ctx)
	if err != nil {
		return err
	}
	for _, p := range builder.AutoModeIAMPolicies {
		if !slices.Contains(existingRolePolicies, p) {
			continue
		}
		if _, err := c.IAMRoleManager.DetachRolePolicy(ctx, &iam.DetachRolePolicyInput{
			PolicyArn: aws.String(makePolicyARN(p, c.Region)),
			RoleName:  aws.String(c.ClusterRoleName),
		}); err != nil {
			return fmt.Errorf("detaching role policy %q: %w", p, err)
		}
	}
	return c.updateAssumeRolePolicy(ctx, eksPolicyDocument)
}

func (c *ClusterRoleManager) updateAssumeRolePolicy(ctx context.Context, policyDocument string) error {
	if _, err := c.IAMRoleManager.UpdateAssumeRolePolicy(ctx, &iam.UpdateAssumeRolePolicyInput{
		RoleName:       aws.String(c.ClusterRoleName),
		PolicyDocument: aws.String(policyDocument),
	}); err != nil {
		return fmt.Errorf("updating assume role policy: %w", err)
	}
	return nil
}

func (c *ClusterRoleManager) hasDedicatedClusterRole(ctx context.Context) (bool, error) {
	templateData, err := c.StackManager.GetStackTemplate(ctx, c.StackManager.MakeClusterStackName())
	if err != nil {
		return false, fmt.Errorf("getting cluster stack template: %w", err)
	}
	template, err := goformation.ParseJSON([]byte(templateData))
	if err != nil {
		return false, fmt.Errorf("parsing cluster stack template: %w", err)
	}
	if _, hasDedicatedClusterRole := template.Resources["ServiceRole"]; !hasDedicatedClusterRole {
		logger.Info("cluster role was not created by eksctl; please ensure the IAM policies required by Auto Mode are attached to the cluster role")
		return false, nil
	}
	return true, nil
}

func (c *ClusterRoleManager) listExistingRolePolicies(ctx context.Context) ([]string, error) {
	paginator := iam.NewListAttachedRolePoliciesPaginator(c.IAMRoleManager, &iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(c.ClusterRoleName),
	})
	var existingPolicyNames []string
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("listing role policies: %w", err)
		}
		for _, policy := range output.AttachedPolicies {
			existingPolicyNames = append(existingPolicyNames, *policy.PolicyName)
		}
	}
	return existingPolicyNames, nil
}

func makePolicyARN(policyName, region string) string {
	return fmt.Sprintf("arn:%s:iam::aws:policy/%s", api.Partitions.ForRegion(region), policyName)
}
