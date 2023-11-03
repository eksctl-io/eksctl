package accessentry

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_stack_creator.go . StackCreator
type StackCreator interface {
	CreateStack(ctx context.Context, stackName string, resourceSet builder.ResourceSetReader, tags, parameters map[string]string, errs chan error) error
	TroubleshootStackFailureCause(ctx context.Context, stack *cfntypes.Stack, desiredStatus cfntypes.StackStatus)
}

type accessEntryTask struct {
	info         string
	clusterName  string
	accessEntry  api.AccessEntry
	stackCreator StackCreator
	ctx          context.Context
}

func (t *accessEntryTask) Describe() string { return t.info }

func (t *accessEntryTask) Do(errorCh chan error) error {
	defer close(errorCh)
	rs := builder.NewAccessEntryResourceSet(t.clusterName, t.accessEntry)
	if err := rs.AddAllResources(); err != nil {
		return err
	}
	principalARN := t.accessEntry.PrincipalARN.String()
	logger.Info("creating access entry for principal ARN %q", principalARN)
	stackErrCh := make(chan error)
	stackName := MakeStackName(t.clusterName, t.accessEntry)
	if err := t.stackCreator.CreateStack(t.ctx, stackName, rs, nil, nil, stackErrCh); err != nil {
		return err
	}
	select {
	case err := <-stackErrCh:
		if err != nil {
			return t.troubleshootFailure(stackName, err)
		}
		logger.Info("created access entry for principal ARN %q", principalARN)
		return nil
	case <-t.ctx.Done():
		return fmt.Errorf("timed out waiting for access entry %q: %w", principalARN, t.ctx.Err())
	}
}

func (t *accessEntryTask) troubleshootFailure(stackName string, err error) error {
	stack := &cfntypes.Stack{
		StackName: aws.String(stackName),
	}
	t.stackCreator.TroubleshootStackFailureCause(t.ctx, stack, cfntypes.StackStatusCreateComplete)
	if strings.Contains(err.Error(), "waiter state transitioned to Failure") {
		return fmt.Errorf("failed to create access entry for principal ARN %q", t.accessEntry.PrincipalARN.String())
	}
	return err
}
