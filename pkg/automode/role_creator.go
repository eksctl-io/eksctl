package automode

import (
	"context"
	"fmt"

	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/cfn/builder"
)

// A StackCreator creates CloudFormation stacks.
type StackCreator interface {
	CreateStack(ctx context.Context, stackName string, resourceSet builder.ResourceSetReader, tags, parameters map[string]string, errs chan error) error
	GetClusterStackIfExists(ctx context.Context) (*cfntypes.Stack, error)
}

// A RoleCreator creates an IAM role for nodes launched by Auto Mode.
type RoleCreator struct {
	StackCreator StackCreator
}

// CreateOrImport creates a new role or imports an existing role if it exists in the cluster stack.
func (r *RoleCreator) CreateOrImport(ctx context.Context, clusterName string) (string, error) {
	clusterStack, err := r.StackCreator.GetClusterStackIfExists(ctx)
	if err != nil {
		return "", fmt.Errorf("getting cluster stack: %w", err)
	}
	if clusterStack != nil {
		if nodeRoleARN, found := builder.GetAutoModeOutputs(*clusterStack); found {
			return nodeRoleARN, nil
		}
	}
	resourceSet, err := builder.CreateAutoModeResourceSet()
	if err != nil {
		return "", err
	}
	errCh := make(chan error)
	if err := r.StackCreator.CreateStack(ctx, makeStackName(clusterName), resourceSet, nil, nil, errCh); err != nil {
		return "", err
	}
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case err := <-errCh:
		if err != nil {
			return "", err
		}
	}
	return resourceSet.GetAutoModeRoleARN(), nil
}

func makeStackName(clusterName string) string {
	return fmt.Sprintf("eksctl-%s-auto-mode-role", clusterName)
}
