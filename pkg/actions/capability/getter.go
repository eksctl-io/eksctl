package capability

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
)

type GetterInterface interface {
	Get(ctx context.Context, capabilityName string) ([]CapabilitySummary, error)
}

type Getter struct {
	ClusterName string
	EKSAPI      awsapi.EKS
}

type CapabilitySummary struct {
	api.Capability
	Status  string
	Version string
}

func NewGetter(clusterName string, eksAPI awsapi.EKS) *Getter {
	return &Getter{
		ClusterName: clusterName,
		EKSAPI:      eksAPI,
	}
}

func (g *Getter) Get(ctx context.Context, capabilityName string) ([]CapabilitySummary, error) {
	toBeFetched := []string{}

	if capabilityName != "" {
		toBeFetched = []string{capabilityName}
	} else {
		out, err := g.EKSAPI.ListCapabilities(ctx, &eks.ListCapabilitiesInput{
			ClusterName: &g.ClusterName,
		})
		if err != nil {
			return nil, fmt.Errorf("calling EKS API to list capabilities: %w", err)
		}
		for _, cap := range out.Capabilities {
			if cap.CapabilityName != nil {
				toBeFetched = append(toBeFetched, *cap.CapabilityName)
			}
		}
	}

	var capabilities []CapabilitySummary
	for _, name := range toBeFetched {
		capability, err := g.getIndividualCapability(ctx, name)
		if err != nil {
			return nil, err
		}
		capabilities = append(capabilities, capability)
	}
	return capabilities, nil
}

func (g *Getter) getIndividualCapability(ctx context.Context, capabilityName string) (CapabilitySummary, error) {
	resp, err := g.EKSAPI.DescribeCapability(ctx, &eks.DescribeCapabilityInput{
		ClusterName:    &g.ClusterName,
		CapabilityName: &capabilityName,
	})
	if err != nil {
		return CapabilitySummary{}, fmt.Errorf("calling EKS API to describe capability %s: %w", capabilityName, err)
	}

	in := resp.Capability
	capability := CapabilitySummary{
		Capability: api.Capability{
			Name:    aws.ToString(in.CapabilityName),
			RoleARN: aws.ToString(in.RoleArn),
			Type:    string(in.Type),
		},
		Status:  string(in.Status),
		Version: aws.ToString(in.Version),
	}

	return capability, nil
}
