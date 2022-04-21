package manager

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	"github.com/weaveworks/eksctl/pkg/cfn/builder"
)

func (c *StackCollection) troubleshootStackFailureCause(ctx context.Context, i *Stack, desiredStatus string) {
	logger.Info("fetching stack events in attempt to troubleshoot the root cause of the failure")
	events, err := c.DescribeStackEvents(ctx, i)
	if err != nil {
		logger.Critical("cannot fetch stack events: %v", err)
		return
	}
	for _, e := range events {
		msg := fmt.Sprintf("%s/%s: %s", *e.ResourceType, *e.LogicalResourceId, e.ResourceStatus)
		if e.ResourceStatusReason != nil {
			msg = fmt.Sprintf("%s – %#v", msg, *e.ResourceStatusReason)
		}
		switch desiredStatus {
		case string(types.StackStatusCreateComplete):
			switch string(e.ResourceStatus) {
			case string(types.ResourceStatusCreateFailed):
				logger.Critical(msg)
			case string(types.ResourceStatusDeleteInProgress):
				logger.Warning(msg)
			default:
				logger.Debug(msg) // only output this when verbose logging is enabled
			}
		case string(types.StackStatusDeleteComplete):
			switch string(e.ResourceStatus) {
			case string(types.ResourceStatusDeleteFailed):
				logger.Critical(msg)
			case string(types.ResourceStatusDeleteSkipped):
				logger.Warning(msg)
			default:
				logger.Debug(msg) // only output this when verbose logging is enabled
			}
		default:
			logger.Info(msg)
		}
	}
}

type noChangeError struct {
	msg string
}

func (e *noChangeError) Error() string {
	return e.msg
}

// DoWaitUntilStackIsCreated blocks until the given stack's
// creation has completed.
func (c *StackCollection) DoWaitUntilStackIsCreated(ctx context.Context, i *Stack) error {
	setCustomRetryer := func(o *cloudformation.StackCreateCompleteWaiterOptions) {
		defaultRetryer := o.Retryable
		o.Retryable = func(ctx context.Context, in *cloudformation.DescribeStacksInput, out *cloudformation.DescribeStacksOutput, err error) (bool, error) {
			logger.Info("waiting for CloudFormation stack %q", *i.StackName)
			return defaultRetryer(ctx, in, out, err)
		}
	}

	waiter := cloudformation.NewStackCreateCompleteWaiter(c.cloudformationAPI)
	return waiter.Wait(ctx, &cloudformation.DescribeStacksInput{
		StackName: i.StackName,
	}, c.waitTimeout, setCustomRetryer)
}

func (c *StackCollection) waitUntilStackIsCreated(ctx context.Context, i *Stack, stack builder.ResourceSetReader, errs chan error) {
	defer close(errs)

	if err := c.DoWaitUntilStackIsCreated(ctx, i); err != nil {
		errs <- err
		return
	}
	s, err := c.DescribeStack(ctx, i)
	if err != nil {
		errs <- err
		return
	}
	if err := stack.GetAllOutputs(*s); err != nil {
		errs <- errors.Wrapf(err, "getting stack %q outputs", *i.StackName)
		return
	}
	errs <- nil
}

func (c *StackCollection) doWaitUntilStackIsDeleted(ctx context.Context, i *Stack) error {
	setCustomRetryer := func(o *cloudformation.StackDeleteCompleteWaiterOptions) {
		defaultRetryer := o.Retryable
		o.Retryable = func(ctx context.Context, in *cloudformation.DescribeStacksInput, out *cloudformation.DescribeStacksOutput, err error) (bool, error) {
			logger.Info("waiting for CloudFormation stack %q", *i.StackName)
			return defaultRetryer(ctx, in, out, err)
		}
	}

	waiter := cloudformation.NewStackDeleteCompleteWaiter(c.cloudformationAPI)
	return waiter.Wait(ctx, &cloudformation.DescribeStacksInput{
		StackName: i.StackName,
	}, c.waitTimeout, setCustomRetryer)
}

func (c *StackCollection) waitUntilStackIsDeleted(ctx context.Context, i *Stack, errs chan error) {
	defer close(errs)

	if err := c.doWaitUntilStackIsDeleted(ctx, i); err != nil {
		errs <- err
		return
	}
	errs <- nil
}

func (c *StackCollection) doWaitUntilStackIsUpdated(ctx context.Context, i *Stack) error {
	setCustomRetryer := func(o *cloudformation.StackUpdateCompleteWaiterOptions) {
		defaultRetryer := o.Retryable
		o.Retryable = func(ctx context.Context, in *cloudformation.DescribeStacksInput, out *cloudformation.DescribeStacksOutput, err error) (bool, error) {
			logger.Info("waiting for CloudFormation stack %q", *i.StackName)
			return defaultRetryer(ctx, in, out, err)
		}
	}

	waiter := cloudformation.NewStackUpdateCompleteWaiter(c.cloudformationAPI)
	return waiter.Wait(ctx, &cloudformation.DescribeStacksInput{
		StackName: i.StackName,
	}, c.waitTimeout, setCustomRetryer)
}

func (c *StackCollection) doWaitUntilChangeSetIsCreated(ctx context.Context, i *Stack, changesetName string) error {
	setCustomRetryer := func(o *cloudformation.ChangeSetCreateCompleteWaiterOptions) {
		defaultRetryer := o.Retryable
		o.Retryable = func(ctx context.Context, in *cloudformation.DescribeChangeSetInput, out *cloudformation.DescribeChangeSetOutput, err error) (bool, error) {
			logger.Info("waiting for CloudFormation changeset %q for stack %q", changesetName, *i.StackName)
			if out.StatusReason != nil && strings.Contains(*out.StatusReason, "The submitted information didn't contain changes") {
				logger.Info("nothing to update")
				return false, &noChangeError{*out.StatusReason}
			}
			return defaultRetryer(ctx, in, out, err)
		}
	}

	waiter := cloudformation.NewChangeSetCreateCompleteWaiter(c.cloudformationAPI, setCustomRetryer)
	return waiter.Wait(ctx, &cloudformation.DescribeChangeSetInput{
		StackName:     i.StackName,
		ChangeSetName: &changesetName,
	}, c.waitTimeout)
}
