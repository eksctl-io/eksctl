package manager

import (
	"fmt"
	"strings"
	"time"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
)

// MakeChangeSetName builds a consistent name for a changeset.
func (c *StackCollection) MakeChangeSetName(action string) string {
	return fmt.Sprintf("eksctl-%s-%d", action, time.Now().Unix())
}

func (c *StackCollection) makeClusterStackName() string {
	return "eksctl-" + c.spec.Metadata.Name + "-cluster"
}

// createClusterTask creates the cluster
func (c *StackCollection) createClusterTask(errs chan error, supportsManagedNodes bool) error {
	name := c.makeClusterStackName()
	logger.Info("building cluster stack %q", name)
	stack := builder.NewClusterResourceSet(c.provider, c.spec, supportsManagedNodes, nil)
	if err := stack.AddAllResources(); err != nil {
		return err
	}

	// Unlike with `createNodeGroupTask`, all tags are already set for the cluster stack
	return c.CreateStack(name, stack, nil, nil, errs)
}

// DescribeClusterStack calls DescribeStacks and filters out cluster stack
func (c *StackCollection) DescribeClusterStack() (*Stack, error) {
	stacks, err := c.DescribeStacks()
	if err != nil {
		return nil, err
	}

	for _, s := range stacks {
		if *s.StackStatus == cfn.StackStatusDeleteComplete {
			continue
		}
		if getClusterName(s) != "" {
			return s, nil
		}
	}
	return nil, c.errStackNotFound()
}

// RefreshFargatePodExecutionRoleARN reads the CloudFormation stacks and
// their output values, and sets the Fargate pod execution role ARN to
// the ClusterConfig.
func (c *StackCollection) RefreshFargatePodExecutionRoleARN() error {
	fargateOutputs := map[string]outputs.Collector{
		outputs.FargatePodExecutionRoleARN: func(v string) error {
			c.spec.IAM.FargatePodExecutionRoleARN = &v
			return nil
		},
	}
	stack, err := c.DescribeClusterStack()
	if err != nil {
		return err
	}
	if err := outputs.Collect(*stack, nil, fargateOutputs); err != nil {
		return err
	}

	if c.spec.IAM.FargatePodExecutionRoleARN == nil {
		logger.Info("Fargate pod execution role is missing, fixing cluster stack to add Fargate resources")
		if err := c.FixClusterCompatibility(); err != nil {
			return errors.Wrap(err, "error fixing cluster compatibility")
		}
	}

	stack, err = c.DescribeClusterStack()
	if err != nil {
		return err
	}
	return outputs.Collect(*stack, fargateOutputs, nil)
}

// AppendNewClusterStackResource will update cluster
// stack with new resources in append-only way
func (c *StackCollection) AppendNewClusterStackResource(plan, supportsManagedNodes bool) (bool, error) {
	name := c.makeClusterStackName()

	// NOTE: currently we can only append new resources to the stack,
	// as there are a few limitations:
	// - it must work with VPC that is imported as well as VPC that
	//   is managed as part of the stack;
	// - CloudFormation cannot yet upgrade EKS control plane itself;

	currentTemplate, err := c.GetStackTemplate(name)
	if err != nil {
		return false, errors.Wrapf(err, "error getting stack template %s", name)
	}

	currentResources := gjson.Get(currentTemplate, resourcesRootPath)
	currentOutputs := gjson.Get(currentTemplate, outputsRootPath)
	currentMappings := gjson.Get(currentTemplate, mappingsRootPath)
	if !currentResources.IsObject() || !currentOutputs.IsObject() {
		return false, fmt.Errorf("unexpected template format of the current stack ")
	}

	logger.Info("re-building cluster stack %q", name)
	newStack := builder.NewClusterResourceSet(c.provider, c.spec, supportsManagedNodes, &currentResources)
	if err := newStack.AddAllResources(); err != nil {
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
	return true, c.UpdateStack(name, c.MakeChangeSetName("update-cluster"), describeUpdate, []byte(currentTemplate), nil)
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
