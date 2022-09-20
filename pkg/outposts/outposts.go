package outposts

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/outposts"
	outpoststypes "github.com/aws/aws-sdk-go-v2/service/outposts/types"

	"github.com/weaveworks/eksctl/pkg/awsapi"
	instanceutils "github.com/weaveworks/eksctl/pkg/utils/instance"
)

type Service struct {
	OutpostsAPI awsapi.Outposts
	EC2API      awsapi.EC2
	OutpostID   string

	mu                   sync.Mutex
	instanceTypes        []ec2types.InstanceType
	instanceTypeInfoList []ec2types.InstanceTypeInfo
	smallestInstanceType string
}

// GetSmallestInstanceType retrieves the smallest available instance type on Outposts.
// Instance types that have a smaller vCPU are considered smaller.
func (o *Service) GetSmallestInstanceType(ctx context.Context) (string, error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.smallestInstanceType != "" {
		return o.smallestInstanceType, nil
	}
	instanceTypes, err := o.describeOutpostInstanceTypes(ctx)
	if err != nil {
		return "", err
	}
	o.smallestInstanceType = instanceutils.GetSmallestInstanceType(instanceTypes)
	return o.smallestInstanceType, nil
}

func (o *Service) getOutpostInstanceTypes(ctx context.Context) ([]ec2types.InstanceType, error) {
	if o.instanceTypes != nil {
		return o.instanceTypes, nil
	}

	paginator := outposts.NewGetOutpostInstanceTypesPaginator(o.OutpostsAPI, &outposts.GetOutpostInstanceTypesInput{
		OutpostId: aws.String(o.OutpostID),
	})
	var instanceTypes []ec2types.InstanceType
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("error fetching Outpost instance types: %w", err)
		}
		for _, it := range output.InstanceTypes {
			instanceTypes = append(instanceTypes, ec2types.InstanceType(aws.ToString(it.InstanceType)))
		}
	}
	if len(instanceTypes) == 0 {
		return nil, fmt.Errorf("no instance types found for Outpost %q", o.OutpostID)
	}
	o.instanceTypes = instanceTypes
	return o.instanceTypes, nil
}

func (o *Service) describeOutpostInstanceTypes(ctx context.Context) ([]ec2types.InstanceTypeInfo, error) {
	if o.instanceTypeInfoList != nil {
		return o.instanceTypeInfoList, nil
	}
	instanceTypes, err := o.getOutpostInstanceTypes(ctx)
	if err != nil {
		return nil, err
	}
	var instanceTypeInfoList []ec2types.InstanceTypeInfo
	instanceTypesPaginator := ec2.NewDescribeInstanceTypesPaginator(o.EC2API, &ec2.DescribeInstanceTypesInput{
		InstanceTypes: instanceTypes,
	})
	for instanceTypesPaginator.HasMorePages() {
		output, err := instanceTypesPaginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("error describing Outpost instance types: %w", err)
		}
		instanceTypeInfoList = append(instanceTypeInfoList, output.InstanceTypes...)
	}
	if len(instanceTypeInfoList) == 0 {
		return nil, fmt.Errorf("no instance description found for instance types: %v", instanceTypes)
	}
	o.instanceTypeInfoList = instanceTypeInfoList
	return o.instanceTypeInfoList, nil
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes . OutpostInstance

// OutpostInstance represents an instance running on Outposts.
type OutpostInstance interface {
	// SetInstanceType sets the instance type.
	SetInstanceType(instanceType string)

	// GetInstanceType returns the instance type.
	GetInstanceType() string
}

// SetOrValidateOutpostInstanceType sets the instance type if it is not set, or validates that the specified instance
// type exists in Outposts.
func (o *Service) SetOrValidateOutpostInstanceType(ctx context.Context, oi OutpostInstance) error {
	if instanceType := oi.GetInstanceType(); instanceType != "" {
		return o.ValidateInstanceType(ctx, instanceType)
	}

	smallestInstanceType, err := o.GetSmallestInstanceType(ctx)
	if err != nil {
		return fmt.Errorf("error getting smallest instance type: %w", err)
	}
	oi.SetInstanceType(smallestInstanceType)
	return nil
}

// ValidateInstanceType validates that instanceType is a valid instance type for this Outpost.
func (o *Service) ValidateInstanceType(ctx context.Context, instanceType string) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	instanceTypes, err := o.getOutpostInstanceTypes(ctx)
	if err != nil {
		return fmt.Errorf("error retrieving instance types for Outpost: %w", err)
	}
	for _, it := range instanceTypes {
		if it == ec2types.InstanceType(instanceType) {
			return nil
		}
	}
	return fmt.Errorf("instance type %q does not exist in Outpost %q", instanceType, o.OutpostID)
}

// GetOutpost returns details for this Outpost.
func (o *Service) GetOutpost(ctx context.Context) (*outpoststypes.Outpost, error) {
	outpost, err := o.OutpostsAPI.GetOutpost(ctx, &outposts.GetOutpostInput{
		OutpostId: aws.String(o.OutpostID),
	})
	if err != nil {
		return nil, err
	}
	return outpost.Outpost, nil
}
