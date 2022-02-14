package cluster

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	"github.com/kris-nova/logger"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/eks"
)

const (
	eksctlCreatedTrue    api.EKSCTLCreated = "True"
	eksctlCreatedFalse   api.EKSCTLCreated = "False"
	eksctlCreatedUnknown api.EKSCTLCreated = "Unknown"
)

type Description struct {
	Name   string
	Region string
	Owned  api.EKSCTLCreated
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_aws_provider.go . ProviderConstructor
type ProviderConstructor func(spec *api.ProviderConfig, clusterSpec *api.ClusterConfig) (*eks.ClusterProvider, error)

//counterfeiter:generate -o fakes/fake_stack_provider.go . StackManagerConstructor
type StackManagerConstructor func(provider api.ClusterProvider, spec *api.ClusterConfig) manager.StackManager

var (
	newClusterProvider ProviderConstructor     = eks.New
	newStackCollection StackManagerConstructor = manager.NewStackCollection
)

func GetClusters(provider api.ClusterProvider, listAllRegions bool, chunkSize int) ([]Description, error) {
	if !listAllRegions {
		return listClusters(provider, int64(chunkSize))
	}

	var clusters []Description
	authorizedRegionsList, err := provider.EC2().DescribeRegions(&ec2.DescribeRegionsInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to describe regions: %w", err)
	}

	authorizedRegions := map[string]struct{}{}
	for _, r := range authorizedRegionsList.Regions {
		authorizedRegions[*r.RegionName] = struct{}{}
	}

	for _, region := range api.SupportedRegions() {
		if _, authorized := authorizedRegions[region]; !authorized {
			continue
		}
		// Reset region and recreate the client.
		ctl, err := newClusterProvider(&api.ProviderConfig{
			Region:      region,
			Profile:     provider.Profile(),
			WaitTimeout: provider.WaitTimeout(),
		}, nil)

		if err != nil {
			logger.Critical("error creating provider in %q region: %v", region, err)
			continue
		}

		newClusters, err := listClusters(ctl.Provider, int64(chunkSize))
		if err != nil {
			logger.Critical("error listing clusters in %q region: %v", region, err)
			continue
		}

		clusters = append(clusters, newClusters...)
	}

	return clusters, nil
}

func listClusters(provider api.ClusterProvider, chunkSize int64) ([]Description, error) {
	var allClusters []Description

	spec := &api.ClusterConfig{Metadata: &api.ClusterMeta{Name: ""}}
	stackManager := newStackCollection(provider, spec)
	allStacks, err := stackManager.ListClusterStackNames()
	if err != nil {
		return nil, fmt.Errorf("failed to list cluster stacks in region %q: %w", provider.Region(), err)
	}

	token := ""
	for {
		clusters, nextToken, err := getClustersRequest(provider, chunkSize, token)
		if err != nil {
			return nil, err
		}

		for _, clusterName := range clusters {
			hasClusterStack, err := stackManager.HasClusterStackUsingCachedList(allStacks, *clusterName)
			managed := eksctlCreatedFalse
			if err != nil {
				managed = eksctlCreatedUnknown
				logger.Warning("error fetching stacks for cluster %s: %v", clusterName, err)
			} else if hasClusterStack {
				managed = eksctlCreatedTrue
			}
			allClusters = append(allClusters, Description{
				Name:   *clusterName,
				Region: provider.Region(),
				Owned:  managed,
			})
		}

		if api.IsSetAndNonEmptyString(nextToken) {
			token = *nextToken
		} else {
			break
		}
	}

	return allClusters, nil
}

func getClustersRequest(provider api.ClusterProvider, chunkSize int64, nextToken string) ([]*string, *string, error) {
	input := &awseks.ListClustersInput{
		MaxResults: &chunkSize,
		Include:    aws.StringSlice([]string{"all"}),
	}
	if nextToken != "" {
		input = input.SetNextToken(nextToken)
	}
	output, err := provider.EKS().ListClusters(input)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list clusters in region %q: %w", provider.Region(), err)
	}
	return output.Clusters, output.NextToken, nil
}
