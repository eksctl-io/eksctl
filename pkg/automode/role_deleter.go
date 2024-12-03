package automode

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

// A StackDeleter deletes CloudFormation stacks.
type StackDeleter interface {
	DeleteStackSync(context.Context, *cfntypes.Stack) error
	DescribeStack(ctx context.Context, stack *cfntypes.Stack) (*cfntypes.Stack, error)
}

// A RoleDeleter deletes the IAM role created for Auto Mode.
type RoleDeleter struct {
	StackDeleter StackDeleter
	Cluster      *ekstypes.Cluster
}

// DeleteIfRequired deletes the node role used by Auto Mode if it exists.
func (d *RoleDeleter) DeleteIfRequired(ctx context.Context) error {
	if cc := d.Cluster.ComputeConfig; cc == nil || !*cc.Enabled {
		return nil
	}
	stack, err := d.StackDeleter.DescribeStack(ctx, &cfntypes.Stack{StackName: aws.String(makeStackName(*d.Cluster.Name))})
	if err != nil {
		if manager.IsStackDoesNotExistError(err) {
			return nil
		}
		return fmt.Errorf("describing Auto Mode stack: %w", err)
	}
	if err := d.StackDeleter.DeleteStackSync(ctx, stack); err != nil {
		return fmt.Errorf("deleting Auto Mode resources: %w", err)
	}
	return nil
}
