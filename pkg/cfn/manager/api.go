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
	return &StackCollection{
		cfn:  provider.CloudFormation(),
		spec: spec,
		tags: []*cloudformation.Tag{
			newTag(ClusterNameTag, spec.ClusterName),
		},
	}
}

func (c *StackCollection) doCreateStackRequest(name string, templateBody []byte, parameters map[string]string, withIAM bool) error {
	input := &cloudformation.CreateStackInput{}

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
	s, err := c.cfn.CreateStack(input)
	if err != nil {
		return errors.Wrapf(err, "creating CloudFormation stack %q", name)
	}
	logger.Debug("stack = %#v", s)

	return nil
}

// CreateStack with given name, stack builder instance and parameters;
// any errors will be written to errs channel, when nil is written,
// assume completion, do not expect more then one error value on the
// channel, it's closed immediately after it is written two
func (c *StackCollection) CreateStack(name string, stack builder.ResourceSet, parameters map[string]string, errs chan error) error {
	templateBody, err := stack.RenderJSON()
	if err != nil {
		return errors.Wrapf(err, "rendering template for %q stack", name)
	}
	logger.Debug("templateBody = %s", string(templateBody))

	if err := c.doCreateStackRequest(name, templateBody, parameters, stack.WithIAM()); err != nil {
		return err
	}

	go c.waitUntilStackIsCreated(name, stack, errs)

	return nil
}

func (c *StackCollection) describeStack(name string) (*Stack, error) {
	input := &cloudformation.DescribeStacksInput{
		StackName: &name,
	}
	resp, err := c.cfn.DescribeStacks(input)
	if err != nil {
		return nil, errors.Wrapf(err, "describing CloudFormation stack %q", name)
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
				stack, subErr = c.describeStack(*s.StackId)
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
func (c *StackCollection) DeleteStack(name string) error {
	s, err := c.describeStack(name)
	if err != nil {
		return errors.Wrapf(err, "not able to get stack %q for deletion", name)
	}

	for _, tag := range s.Tags {
		if *tag.Key == ClusterNameTag && *tag.Value == c.spec.ClusterName {
			input := &cloudformation.DeleteStackInput{
				StackName: s.StackName,
			}

			if _, err := c.cfn.DeleteStack(input); err != nil {
				return errors.Wrapf(err, "not able to delete stack %q", name)
			}
			logger.Info("will delete stack %q", name)
			return nil
		}
	}

	return fmt.Errorf("cannot delete stack %q as it doesn't bare our %q tag", *s.StackName,
		fmt.Sprintf("%s:%s", ClusterNameTag, c.spec.ClusterName))
}

// WaitDeleteStack kills a stack by name and waits for DELETED status
func (c *StackCollection) WaitDeleteStack(name string) error {
	if err := c.DeleteStack(name); err != nil {
		return err
	}

	logger.Info("waiting for stack %q to get deleted", name)

	return c.doWaitUntilStackIsDeleted(name)
}

func (c *StackCollection) DescribeStacks(name string) ([]*Stack, error) {
	stacks, err := c.ListStacks(fmt.Sprintf("^(eksclt|EKS)-%s-((cluster|nodegroup)-\\d+|(VPC|ServiceRole|DefaultNodeGroup))$", name))
	if err != nil {
		return nil, errors.Wrapf(err, "describing CloudFormation stacks for %q", name)
	}
	return stacks, nil
}

func (c *StackCollection) DescribeStackEvents(s *Stack) ([]*cloudformation.StackEvent, error) {
	input := &cloudformation.DescribeStackEventsInput{
		StackName: s.StackId,
	}

	resp, err := c.cfn.DescribeStackEvents(input)
	if err != nil {
		return nil, errors.Wrapf(err, "describing CloudFormation stack %q events", *s.StackName)
	}
	return resp.StackEvents, nil
}
