//go:build integration
// +build integration

//revive:disable Not changing package name
package local_zones

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	clusterutils "github.com/weaveworks/eksctl/integration/utilities/cluster"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("local-zones")
	if err := api.Register(); err != nil {
		panic("unexpected error registering API scheme")
	}
}

func TestLocalZones(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("AWS Local Zones", Ordered, func() {
	var (
		clusterConfig *api.ClusterConfig
		localZones    = []string{"us-west-2-lax-1a", "us-west-2-lax-1b"}
		instanceType  string
	)

	BeforeAll(func() {
		const supportedRegion = "us-west-2"
		if params.Region != supportedRegion {
			Fail(fmt.Sprintf("local zones test is only supported on %s", supportedRegion))
		}
		By("creating a cluster with local zones")
		clusterConfig = clusterutils.ParseClusterConfig(params.ClusterName, params.Region, "testdata/local-zones.yaml")

		By(fmt.Sprintf("selecting an instance type available in %s", strings.Join(localZones, ", ")))
		ec2API := ec2.NewFromConfig(NewConfig(params.Region))

		var err error
		instanceType, err = selectInstanceType(ec2API, localZones)
		Expect(err).NotTo(HaveOccurred())
		fmt.Fprintf(GinkgoWriter, "selected instance type %q\n", instanceType)
		for _, ng := range clusterConfig.NodeGroups {
			ng.InstanceType = instanceType
		}

		cmd := params.EksctlCreateCmd.
			WithArgs(
				"cluster",
				"--config-file=-",
				"--kubeconfig", params.KubeconfigPath,
				"--verbose=4",
			).
			WithoutArg("--region", params.Region).
			WithStdin(clusterutils.Reader(clusterConfig))
		Expect(cmd).To(RunSuccessfully())

		DeferCleanup(params.DeleteClusters)
	})

	params.LogStacksEventsOnFailure()

	It("should create unmanaged nodegroups in local zones", func() {
		desiredCapacity := 2
		clusterConfig.NodeGroups = append(clusterConfig.NodeGroups, &api.NodeGroup{
			LocalZones: localZones,
			NodeGroupBase: &api.NodeGroupBase{
				Name:         "local-zones",
				InstanceType: instanceType,
				ScalingConfig: &api.ScalingConfig{
					DesiredCapacity: aws.Int(desiredCapacity),
				},
			},
		})
		cmd := params.EksctlCreateCmd.
			WithArgs(
				"nodegroup",
				"--config-file=-",
				"--verbose=4",
			).
			WithoutArg("--region", params.Region).
			WithStdin(clusterutils.Reader(clusterConfig))
		Expect(cmd).To(RunSuccessfully())

		config, err := clientcmd.BuildConfigFromFlags("", params.KubeconfigPath)
		Expect(err).NotTo(HaveOccurred())
		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).NotTo(HaveOccurred())

		nodeList := tests.ListNodes(clientset, "local-zones")
		Expect(nodeList.Items).To(HaveLen(desiredCapacity))
		for _, node := range nodeList.Items {
			nodeZone := node.Labels["topology.kubernetes.io/zone"]
			Expect(localZones).To(ContainElement(nodeZone))
		}
	})
})

func selectInstanceType(ec2API awsapi.EC2, zones []string) (string, error) {
	output, err := ec2API.DescribeInstanceTypeOfferings(context.Background(), &ec2.DescribeInstanceTypeOfferingsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("location"),
				Values: zones,
			},
		},
		LocationType: ec2types.LocationTypeAvailabilityZone,
	})
	if err != nil {
		return "", fmt.Errorf("error describing instance type offerings: %w", err)
	}

	preferredInstanceTypes := []ec2types.InstanceType{
		ec2types.InstanceTypeT3Small,
		ec2types.InstanceTypeT3Medium,
		ec2types.InstanceTypeT3Large,
		ec2types.InstanceTypeM5Xlarge,
	}
	for _, it := range preferredInstanceTypes {
		instanceTypeCount := 0
		for _, offering := range output.InstanceTypeOfferings {
			if offering.InstanceType == it {
				instanceTypeCount++
				if instanceTypeCount == len(zones) {
					return string(it), nil
				}
			}
		}
	}
	return "", errors.New("failed to find a preferred instance type available in all zones")
}
