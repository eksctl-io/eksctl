package manager

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/eks/api"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"

	"github.com/kubicorn/kubicorn/pkg/logger"
)

const (
	ClusterNameTag = "eksctl.cluster.k8s.io/v1alpha1/cluster-name"
	NodeGroupTagID = "eksctl.cluster.k8s.io/v1alpha1/nodegroup-id"
)

type Stack = cloudformation.Stack

type StackCollection struct {
	cfn  cloudformationiface.CloudFormationAPI
	spec *api.ClusterConfig
	tags []*cloudformation.Tag
}

func newTag(key, value string) *cloudformation.Tag {
	return &cloudformation.Tag{Key: &key, Value: &value}
}

// NewStackCollection create a stack manager for a single cluster
func NewStackCollection(provider api.ClusterProvider, spec *api.ClusterConfig) *StackCollection {
	tags := []*cloudformation.Tag{
		newTag(ClusterNameTag, spec.ClusterName),
	}
	for key, value := range spec.Tags {
		tags = append(tags, newTag(key, value))
	}
	logger.Debug("tags = %#v", tags)
	return &StackCollection{
		cfn:  provider.CloudFormation(),
		spec: spec,
		tags: tags,
	}
}

func (c *StackCollection) doCreateStackRequest(i *Stack, templateBody []byte, parameters map[string]string, withIAM bool) error {
	input := &cloudformation.CreateStackInput{
		StackName: i.StackName,
	}

	input.SetTags(c.tags)
	input.SetTemplateBody(string(templateBody))

	if withIAM {
		input.SetCapabilities(aws.StringSlice([]string{cloudformation.CapabilityCapabilityIam}))
	}

	for k, v := range parameters {
		p := &cloudformation.Parameter{
			ParameterKey:   aws.String(k),
			ParameterValue: aws.String(v),
		}
		input.Parameters = append(input.Parameters, p)
	}

	logger.Debug("input = %#v", input)
	s, err := c.cfn.CreateStack(input)
	if err != nil {
		return errors.Wrapf(err, "creating CloudFormation stack %q", *i.StackName)
	}
	logger.Debug("stack = %#v", s)
	i.StackId = s.StackId
	return nil
}

func (c *StackCollection) doUpdateStackRequest(name string, templateBody []byte, parameters map[string]string, withIAM bool) error {
	input := &cloudformation.UpdateStackInput{}

	input.SetStackName(name)
	input.SetTags(c.tags)

	input.SetTemplateBody(string(templateBody))

	if withIAM {
		input.SetCapabilities(aws.StringSlice([]string{cloudformation.CapabilityCapabilityIam}))
	}

	for k, v := range parameters {
		p := &cloudformation.Parameter{
			ParameterKey:   aws.String(k),
			ParameterValue: aws.String(v),
		}
		input.Parameters = append(input.Parameters, p)
	}

	logger.Debug("input = %#v", input)
	s, err := c.cfn.UpdateStack(input)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("updating CloudFormation stack %q", name))
	}
	logger.Debug("stack = %#v", s)

	return nil
}

// CreateStack with given name, stack builder instance and parameters;
// any errors will be written to errs channel, when nil is written,
// assume completion, do not expect more then one error value on the
// channel, it's closed immediately after it is written two
func (c *StackCollection) CreateStack(name string, stack builder.ResourceSet, parameters map[string]string, errs chan error) error {
	i := &Stack{StackName: &name}
	templateBody, err := stack.RenderJSON()
	if err != nil {
		return errors.Wrapf(err, "rendering template for %q stack", *i.StackName)
	}
	logger.Debug("templateBody = %s", string(templateBody))

	if err := c.doCreateStackRequest(i, templateBody, parameters, stack.WithIAM()); err != nil {
		return err
	}

	go c.waitUntilStackIsCreated(i, stack, errs)

	return nil
}

func (c *StackCollection) UpdateStack(name string, template []byte, parameters map[string]string, errs chan error) error {
	// if err := c.doUpdateStackRequest(name, template, parameters, stack.WithIAM()); err != nil {
	if err := c.doUpdateStackRequest(name, template, parameters, false); err != nil {
		return err
	}

	//TODO: waitUntilStackIsUpdated
	//go c.waitUntilStackIsCreated(name, stack, errs)

	return nil
}

func (c *StackCollection) describeStack(i *Stack) (*Stack, error) {
	input := &cloudformation.DescribeStacksInput{
		StackName: i.StackName,
	}
	if i.StackId != nil {
		input.StackName = i.StackId
	}
	resp, err := c.cfn.DescribeStacks(input)
	if err != nil {
		return nil, errors.Wrapf(err, "describing CloudFormation stack %q", *i.StackName)
	}
	return resp.Stacks[0], nil
}

// ListStacks gets all of CloudFormation stacks
func (c *StackCollection) ListStacks(nameRegex string, statusFilters ...string) ([]*Stack, error) {
	var (
		subErr error
		stack  *Stack
	)

	re, err := regexp.Compile(nameRegex)
	if err != nil {
		return nil, errors.Wrap(err, "cannot list stacks")
	}
	input := &cloudformation.ListStacksInput{}
	if len(statusFilters) > 0 {
		input.StackStatusFilter = aws.StringSlice(statusFilters)
	}
	stacks := []*Stack{}

	pager := func(p *cloudformation.ListStacksOutput, last bool) (shouldContinue bool) {
		for _, s := range p.StackSummaries {
			if re.MatchString(*s.StackName) {
				stack, subErr = c.describeStack(&Stack{StackName: s.StackName, StackId: s.StackId})
				if subErr != nil {
					return false
				}
				stacks = append(stacks, stack)
			}
		}
		return true
	}
	if err := c.cfn.ListStacksPages(input, pager); err != nil {
		return nil, err
	}
	if subErr != nil {
		return nil, subErr
	}
	return stacks, nil
}

// ListReadyStacks gets all of CloudFormation stacks with READY status
func (c *StackCollection) ListReadyStacks(nameRegex string) ([]*Stack, error) {
	return c.ListStacks(nameRegex, cloudformation.StackStatusCreateComplete)
}

// DeleteStack kills a stack by name without waiting for DELETED status
func (c *StackCollection) DeleteStack(name string) (*Stack, error) {
	i := &Stack{StackName: &name}
	s, err := c.describeStack(i)
	if err != nil {
		return nil, errors.Wrapf(err, "not able to get stack %q for deletion", name)
	}
	i.StackId = s.StackId
	for _, tag := range s.Tags {
		if *tag.Key == ClusterNameTag && *tag.Value == c.spec.ClusterName {
			input := &cloudformation.DeleteStackInput{
				StackName: i.StackId,
			}

			if _, err := c.cfn.DeleteStack(input); err != nil {
				return nil, errors.Wrapf(err, "not able to delete stack %q", name)
			}
			logger.Info("will delete stack %q", name)
			return i, nil
		}
	}

	return nil, fmt.Errorf("cannot delete stack %q as it doesn't bare our %q tag", *s.StackName,
		fmt.Sprintf("%s:%s", ClusterNameTag, c.spec.ClusterName))
}

// WaitDeleteStack kills a stack by name and waits for DELETED status
func (c *StackCollection) WaitDeleteStack(name string) error {
	i, err := c.DeleteStack(name)
	if err != nil {
		return err
	}

	logger.Info("waiting for stack %q to get deleted", *i.StackName)

	return c.doWaitUntilStackIsDeleted(i)
}

func (c *StackCollection) DescribeStacks(name string) ([]*Stack, error) {
	stacks, err := c.ListStacks(fmt.Sprintf("^(eksclt|EKS)-%s-((cluster|nodegroup)-\\d+|(VPC|ServiceRole|DefaultNodeGroup))$", name))
	if err != nil {
		return nil, errors.Wrapf(err, "describing CloudFormation stacks for %q", name)
	}
	return stacks, nil
}

func (c *StackCollection) DescribeStackEvents(i *Stack) ([]*cloudformation.StackEvent, error) {
	input := &cloudformation.DescribeStackEventsInput{
		StackName: i.StackId,
	}

	resp, err := c.cfn.DescribeStackEvents(input)
	if err != nil {
		return nil, errors.Wrapf(err, "describing CloudFormation stack %q events", *i.StackName)
	}
	return resp.StackEvents, nil
}
