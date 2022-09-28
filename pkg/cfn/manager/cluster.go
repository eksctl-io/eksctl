package manager

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
)

// MakeChangeSetName builds a consistent name for a changeset.
func (c *StackCollection) MakeChangeSetName(action string) string {
	return fmt.Sprintf("eksctl-%s-%d", action, time.Now().Unix())
}

func (c *StackCollection) MakeClusterStackName() string {
	return c.MakeClusterStackNameFromName(c.spec.Metadata.Name)
}

func (c *StackCollection) MakeClusterStackNameFromName(name string) string {
	return "eksctl-" + name + "-cluster"
}

// createClusterTask creates the cluster
func (c *StackCollection) createClusterTask(ctx context.Context, errs chan error, supportsManagedNodes bool) error {
	name := c.MakeClusterStackName()
	logger.Info("building cluster stack %q", name)
	stack := builder.NewClusterResourceSet(c.ec2API, c.region, c.spec, nil, false)
	if err := stack.AddAllResources(ctx); err != nil {
		return err
	}
	return c.createClusterStack(ctx, name, stack, errs)
}

// DescribeClusterStackIfExists calls ListStacks and filters out cluster stack.
// If the stack does not exist, it returns nil.
func (c *StackCollection) DescribeClusterStackIfExists(ctx context.Context) (*Stack, error) {
	stacks, err := c.ListStacks(ctx)
	if err != nil {
		return nil, err
	}

	if len(stacks) == 0 {
		return nil, nil
	}

	for _, s := range stacks {
		if s.StackStatus == types.StackStatusDeleteComplete {
			continue
		}
		if getClusterName(s) != "" {
			return s, nil
		}
	}
	return nil, nil
}

// DescribeClusterStack returns the cluster stack. If the stack does not exist, it returns an error.
func (c *StackCollection) DescribeClusterStack(ctx context.Context) (*Stack, error) {
	stack, err := c.DescribeClusterStackIfExists(ctx)
	if err != nil {
		return nil, err
	}
	if stack == nil {
		return nil, &StackNotFoundErr{
			ClusterName: c.spec.Metadata.Name,
		}
	}
	return stack, nil
}

// RefreshFargatePodExecutionRoleARN reads the CloudFormation stacks and
// their output values, and sets the Fargate pod execution role ARN to
// the ClusterConfig. If there is no cluster stack found but a fargate stack
// exists, use the output from that stack.
func (c *StackCollection) RefreshFargatePodExecutionRoleARN(ctx context.Context) error {
	fargateOutputs := map[string]outputs.Collector{
		outputs.FargatePodExecutionRoleARN: func(v string) error {
			c.spec.IAM.FargatePodExecutionRoleARN = &v
			return nil
		},
	}
	stack, err := c.DescribeClusterStackIfExists(ctx)
	if err != nil {
		return err
	}
	//check if fargate role is on the cluster stack
	if stack != nil {
		if err := outputs.Collect(*stack, nil, fargateOutputs); err != nil {
			return err
		}

		if c.spec.IAM.FargatePodExecutionRoleARN != nil {
			return nil
		}
	}

	//check if fargate role is in separate stack
	stack, err = c.GetFargateStack(ctx)
	if err != nil {
		return err
	}

	return outputs.Collect(*stack, fargateOutputs, nil)
}

// AppendNewClusterStackResource will update cluster
// stack with new resources in append-only way
func (c *StackCollection) AppendNewClusterStackResource(ctx context.Context, extendForOutposts, plan bool) (bool, error) {
	name := c.MakeClusterStackName()

	// NOTE: currently we can only append new resources to the stack,
	// as there are a few limitations:
	// - it must work with VPC that is imported as well as VPC that
	//   is managed as part of the stack;
	// - CloudFormation cannot yet upgrade EKS control plane itself;

	currentTemplate, err := c.GetStackTemplate(ctx, name)
	if err != nil {
		return false, errors.Wrapf(err, "error getting stack template %s", name)
	}

	currentResources := gjson.Get(currentTemplate, resourcesRootPath)
	currentOutputs := gjson.Get(currentTemplate, outputsRootPath)
	currentMappings := gjson.Get(currentTemplate, mappingsRootPath)
	if !currentResources.IsObject() || !currentOutputs.IsObject() {
		return false, fmt.Errorf("unexpected template format of the current stack ")
	}

	if err := c.importServiceRoleARN(ctx, currentResources); err != nil {
		return false, err
	}

	logger.Info("re-building cluster stack %q", name)
	newStack := builder.NewClusterResourceSet(c.ec2API, c.region, c.spec, &currentResources, extendForOutposts)
	if err := newStack.AddAllResources(ctx); err != nil {
		return false, err
	}

	newTemplate, err := newStack.RenderJSON()
	if err != nil {
		return false, errors.Wrapf(err, "rendering template for %q stack", name)
	}
	logger.Debug("newTemplate = %s", newTemplate)

	newResources := gjson.Get(string(newTemplate), resourcesRootPath)
	newOutputs := gjson.Get(string(newTemplate), outputsRootPath)
	newMappings := gjson.Get(string(newTemplate), mappingsRootPath)
	if !newResources.IsObject() || !newOutputs.IsObject() || !newMappings.IsObject() {
		return false, errors.New("unexpected template format of the new version of the stack")
	}

	logger.Debug("currentTemplate = %s", currentTemplate)

	var iterErr error
	iterFunc := func(list *[]string, root string, currentSet, key, value gjson.Result) bool {
		k := key.String()
		if currentSet.Get(k).Exists() {
			return true
		}
		*list = append(*list, k)
		path := root + "." + k
		currentTemplate, iterErr = sjson.Set(currentTemplate, path, value.Value())
		return iterErr == nil
	}

	var (
		addResources []string
		addOutputs   []string
		addMappings  []string
	)

	newResources.ForEach(func(k, v gjson.Result) bool {
		return iterFunc(&addResources, resourcesRootPath, currentResources, k, v)
	})
	if iterErr != nil {
		return false, errors.Wrap(iterErr, "adding resources to current stack template")
	}
	newOutputs.ForEach(func(k, v gjson.Result) bool {
		return iterFunc(&addOutputs, outputsRootPath, currentOutputs, k, v)
	})
	if iterErr != nil {
		return false, errors.Wrap(iterErr, "adding outputs to current stack template")
	}

	newMappings.ForEach(func(k, v gjson.Result) bool {
		return iterFunc(&addMappings, mappingsRootPath, currentMappings, k, v)
	})
	if iterErr != nil {
		return false, errors.Wrap(iterErr, "adding mappings to current stack template")
	}

	if len(addResources) == 0 && len(addOutputs) == 0 && len(addMappings) == 0 {
		logger.Success("all resources in cluster stack %q are up-to-date", name)
		return false, nil
	}

	logger.Debug("currentTemplate = %s", currentTemplate)

	describeUpdate := fmt.Sprintf("updating stack to add new resources %v and outputs %v", addResources, addOutputs)
	if plan {
		logger.Info("(plan) %s", describeUpdate)
		return true, nil
	}

	err = c.UpdateStack(ctx, UpdateStackOptions{
		StackName:     name,
		ChangeSetName: c.MakeChangeSetName("update-cluster"),
		Description:   describeUpdate,
		TemplateData:  TemplateBody(currentTemplate),
		Wait:          true,
	})
	if err != nil {
		return false, err
	}
	stack, err := c.DescribeStack(ctx, &Stack{
		StackName: aws.String(name),
	})
	if err != nil {
		return false, fmt.Errorf("error describing cluster stack: %w", err)
	}
	if err := newStack.GetAllOutputs(*stack); err != nil {
		return false, fmt.Errorf("error getting outputs for updated cluster stack: %w", err)
	}
	return true, nil
}

// ClusterHasDedicatedVPC returns true if the cluster was created with a dedicated VPC.
func (c *StackCollection) ClusterHasDedicatedVPC(ctx context.Context) (bool, error) {
	stackName := c.MakeClusterStackName()
	currentTemplate, err := c.GetStackTemplate(ctx, stackName)
	if err != nil {
		return false, fmt.Errorf("error getting stack template %q: %w", stackName, err)
	}

	resources := gjson.Get(currentTemplate, resourcesRootPath)
	return resources.IsObject() && resources.Get("VPC").Exists(), nil
}

func (c *StackCollection) importServiceRoleARN(ctx context.Context, resources gjson.Result) error {
	s, err := c.DescribeClusterStackIfExists(ctx)
	if err != nil {
		return err
	}
	usesEksctlCreatedServiceRole := false
	resources.ForEach(func(key, value gjson.Result) bool {
		if key.String() == "ServiceRole" {
			usesEksctlCreatedServiceRole = true
		}
		return true
	})

	if usesEksctlCreatedServiceRole {
		return nil
	}

	for _, o := range s.Outputs {
		if *o.OutputKey == "ServiceRoleARN" {
			c.spec.IAM.ServiceRoleARN = o.OutputValue
		}
	}
	return nil
}

func getClusterName(s *Stack) string {
	if strings.HasSuffix(*s.StackName, "-cluster") {
		if v := getClusterNameTag(s); v != "" {
			return v
		}
	}

	if strings.HasPrefix(*s.StackName, "EKS-") && strings.HasSuffix(*s.StackName, "-ControlPlane") {
		return strings.TrimPrefix("EKS-", strings.TrimSuffix(*s.StackName, "-ControlPlane"))
	}
	return ""
}

func getClusterNameTag(s *Stack) string {
	for _, tag := range s.Tags {
		if *tag.Key == api.ClusterNameTag || *tag.Key == api.OldClusterNameTag {
			return *tag.Value
		}
	}
	return ""
}
