package autonomousmode

import (
	"context"
	"errors"
	"fmt"

	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

// ClusterStackDescriber describes a cluster's CloudFormation stack.
type ClusterStackDescriber interface {
	DescribeClusterStack(ctx context.Context) (*cfntypes.Stack, error)
	ClusterHasDedicatedVPC(ctx context.Context) (bool, error)
}

// A VPCImporter imports the VPC for an existing cluster.
type VPCImporter interface {
	LoadClusterVPC(ctx context.Context, spec *api.ClusterConfig, stack *manager.Stack, ignoreDrift bool) error
}

// SubnetsLoader loads subnets from an existing VPC.
type SubnetsLoader struct {
	ClusterStackDescriber ClusterStackDescriber
	VPCImporter           VPCImporter
	IgnoreMissingSubnets  bool
}

func (s *SubnetsLoader) LoadSubnets(ctx context.Context, clusterConfig *api.ClusterConfig) ([]string, bool, error) {
	hasDedicatedVPC, err := s.ClusterStackDescriber.ClusterHasDedicatedVPC(ctx)
	if err != nil {
		return nil, false, fmt.Errorf("checking if cluster has a dedicated VPC: %w", err)
	}
	if !hasDedicatedVPC {
		return nil, false, nil
	}
	found, err := s.importVPC(ctx, clusterConfig)
	if err != nil {
		return nil, false, err
	}
	if !found {
		return nil, true, nil
	}
	subnetIDs := clusterConfig.VPC.Subnets.Private.WithIDs()
	switch {
	case len(subnetIDs) == 0:
		return nil, false, errors.New("expected to find private subnets in cluster stack")
	case len(subnetIDs) < 2:
		return nil, false, fmt.Errorf("Autonomous Mode requires at least two private subnets; got %v", subnetIDs)
	default:
		return subnetIDs, true, nil
	}
}

func (s *SubnetsLoader) importVPC(ctx context.Context, clusterConfig *api.ClusterConfig) (bool, error) {
	clusterStack, err := s.ClusterStackDescriber.DescribeClusterStack(ctx)
	if err != nil {
		var stackNotFoundErr *manager.StackNotFoundErr
		if errors.As(err, &stackNotFoundErr) {
			return false, nil
		}
		return false, fmt.Errorf("describing cluster stack: %w", err)
	}

	if err := s.VPCImporter.LoadClusterVPC(ctx, clusterConfig, clusterStack, false); err != nil {
		var stackDriftErr *vpc.StackDriftError
		const msg = "loading cluster VPC"
		if errors.As(err, &stackDriftErr) {
			if s.IgnoreMissingSubnets {
				return false, nil
			}
			return false, fmt.Errorf("%s: %w; to skip patching NodeClass to use private subnets and ignore this error, "+
				"please retry the command with --ignore-missing-subnets and patch the NodeClass "+
				"resource manually if you do not want to use cluster subnets for Autonomous Mode", msg, err)
		}
		return false, fmt.Errorf("%s: %w", msg, err)
	}
	return true, nil
}
