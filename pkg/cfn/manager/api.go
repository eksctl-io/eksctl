package manager

import (
	"fmt"
	"regexp"
	"time"

	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/eks/api"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"

	"github.com/kubicorn/kubicorn/pkg/logger"
)

const (
	ClusterNameTag = "eksctl.cluster.k8s.io/v1alpha1/cluster-name"
)

type Stack = cloudformation.Stack

type StackCollection struct {
	cfn  cloudformationiface.CloudFormationAPI
	spec *api.ClusterConfig
}

func NewStackCollection(provider api.ClusterProvider, spec *api.ClusterConfig) *StackCollection {
	return &StackCollection{
		cfn:  provider.CloudFormation(),
		spec: spec,
	}
}

func (c *StackCollection) CreateStack(name string, templateBody []byte, parameters map[string]string, withIAM bool, stack chan Stack, errs chan error) error {
	input := &cloudformation.CreateStackInput{}
	input.SetStackName(name)
	input.SetTags([]*cloudformation.Tag{
		&cloudformation.Tag{
			Key:   aws.String(ClusterNameTag),
			Value: aws.String(c.spec.ClusterName),
		},
	})
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
		return errors.Wrap(err, fmt.Sprintf("creating CloudFormation stack %q", name))
	}
	logger.Debug("stack = %#v", s)

	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()

		timer := time.NewTimer(c.spec.WaitTimeout)
		defer timer.Stop()

		defer close(errs)
		for {
			select {
			case <-timer.C:
				errs <- fmt.Errorf("timed out creating CloudFormation stack %q after %d", name, c.spec.WaitTimeout)
				return

			case <-ticker.C:
				s, err := c.describeStack(&name)
				if err != nil {
					logger.Warning("continue despite err=%q", err.Error())
					continue
				}
				logger.Debug("stack = %#v", s)
				switch *s.StackStatus {
				case cloudformation.StackStatusCreateInProgress:
					continue
				case cloudformation.StackStatusCreateComplete:
					errs <- nil
					stack <- *s
					return
				case cloudformation.StackStatusCreateFailed:
					fallthrough // TODO: https://github.com/weaveworks/eksctl/issues/24
				default:
					errs <- fmt.Errorf("unexpected status %q while creating CloudFormation stack %q", *s.StackStatus, name)
					// stack <- *s // this usually results in closed channel panic, but we don't need it really
					logger.Debug("stack = %#v", s)
					return
				}
			}
		}
	}()

	return nil
}

func (c *StackCollection) describeStack(name *string) (*Stack, error) {
	input := &cloudformation.DescribeStacksInput{
		StackName: name,
	}
	resp, err := c.cfn.DescribeStacks(input)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("describing CloudFormation stack %q", *name))
	}
	return resp.Stacks[0], nil
}

func (c *StackCollection) ListReadyStacks(nameRegex string) ([]*Stack, error) {
	var (
		subErr error
		stack  *Stack
	)

	re := regexp.MustCompile(nameRegex)
	input := &cloudformation.ListStacksInput{
		StackStatusFilter: aws.StringSlice([]string{cloudformation.StackStatusCreateComplete}),
	}
	stacks := []*Stack{}

	pager := func(p *cloudformation.ListStacksOutput, last bool) (shouldContinue bool) {
		for _, s := range p.StackSummaries {
			if re.MatchString(*s.StackName) {
				stack, subErr = c.describeStack(s.StackName)
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

func (c *StackCollection) DeleteStack(name string) error {
	s, err := c.describeStack(&name)
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

	return fmt.Errorf("cannot delete stack %s as it doesn't bare our %q tag", *s.StackName,
		fmt.Sprintf("%s:%s", ClusterNameTag, c.spec.ClusterName))
}
