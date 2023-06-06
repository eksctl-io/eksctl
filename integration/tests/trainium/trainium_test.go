//go:build integration
// +build integration

package trainium

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"k8s.io/client-go/kubernetes"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var (
	params                  *tests.Params
	clusterWithNeuronPlugin string
	clusterWithoutPlugin    string
	nodeZones               string
	clusterZones            string
)

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("trn1")
	clusterWithNeuronPlugin = params.ClusterName
	clusterWithoutPlugin = params.NewClusterName("trn1-no-plugin")
}

func TestTrainium(t *testing.T) {
	testutils.RegisterAndRun(t)
}

const initNG = "trn1-ng-0"

var _ = BeforeSuite(func() {
	params.KubeconfigTemp = false
	if params.KubeconfigPath == "" {
		wd, _ := os.Getwd()
		f, _ := os.CreateTemp(wd, "kubeconfig-")
		params.KubeconfigPath = f.Name()
		params.KubeconfigTemp = true
	}

	if !params.SkipCreate {
		cfg := NewConfig(params.Region)
		ctx := context.Background()
		ec2API := ec2.NewFromConfig(cfg)
		nodeZones, clusterZones = getAvailabilityZones(ctx, ec2API)

		cmd := params.EksctlCreateCmd.WithArgs(
			"cluster",
			"--verbose", "4",
			"--name", clusterWithoutPlugin,
			"--zones", clusterZones,
			"--tags", "alpha.eksctl.io/description=eksctl integration test",
			"--install-neuron-plugin=false",
			"--nodegroup-name", initNG,
			"--node-labels", "ng-name="+initNG,
			"--nodes", "1",
			"--node-type", "trn1.2xlarge",
			"--node-zones", nodeZones,
			"--version", params.Version,
			"--kubeconfig", params.KubeconfigPath,
		)
		Expect(cmd).To(RunSuccessfully())

		cmd = params.EksctlCreateCmd.WithArgs(
			"cluster",
			"--verbose", "4",
			"--name", clusterWithNeuronPlugin,
			"--zones", clusterZones,
			"--tags", "alpha.eksctl.io/description=eksctl integration test",
			"--nodegroup-name", initNG,
			"--node-labels", "ng-name="+initNG,
			"--nodes", "1",
			"--node-type", "trn1.2xlarge",
			"--node-zones", nodeZones,
			"--version", params.Version,
			"--kubeconfig", params.KubeconfigPath,
		)
		Expect(cmd).To(RunSuccessfully())
	}
})

var _ = Describe("(Integration) Trainium nodes", func() {
	const (
		newNG = "trn1-ng-1"
	)

	Context("cluster with trn1 nodes", func() {
		Context("by default", func() {
			BeforeEach(func() {
				cmd := params.EksctlUtilsCmd.WithArgs(
					"write-kubeconfig",
					"--verbose", "4",
					"--cluster", clusterWithoutPlugin,
					"--kubeconfig", params.KubeconfigPath,
				)
				Expect(cmd).To(RunSuccessfully())
			})

			It("should have installed the neuron device plugin", func() {
				clientSet := newClientSet(clusterWithNeuronPlugin)
				_, err := clientSet.AppsV1().DaemonSets("kube-system").Get(context.TODO(), "neuron-device-plugin-daemonset", metav1.GetOptions{})
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("should not have installed the nvidia device plugin", func() {
				_, err := newClientSet(clusterWithNeuronPlugin).AppsV1().DaemonSets("kube-system").Get(context.TODO(), "nvidia-device-plugin-daemonset", metav1.GetOptions{})
				Expect(err).Should(BeNotFoundError())
			})
		})

		Context("with --install-neuron-plugin=false", func() {
			BeforeEach(func() {
				cmd := params.EksctlUtilsCmd.WithArgs(
					"write-kubeconfig",
					"--verbose", "4",
					"--cluster", clusterWithoutPlugin,
					"--kubeconfig", params.KubeconfigPath,
				)
				Expect(cmd).To(RunSuccessfully())
			})

			It("should not have installed the neuron device plugin", func() {
				_, err := newClientSet(clusterWithoutPlugin).AppsV1().DaemonSets("kube-system").Get(context.TODO(), "neuron-device-plugin-daemonset", metav1.GetOptions{})
				Expect(err).Should(BeNotFoundError())
			})

			When("adding an unmanaged nodegroup by default", func() {
				params.LogStacksEventsOnFailureForCluster(clusterWithoutPlugin)

				It("should install without error", func() {
					cmd := params.EksctlCreateCmd.WithArgs(
						"nodegroup",
						"--cluster", clusterWithoutPlugin,
						"--managed=false",
						"--nodes", "1",
						"--verbose", "4",
						"--name", newNG,
						"--tags", "alpha.eksctl.io/description=eksctl integration test",
						"--node-labels", "ng-name="+newNG,
						"--nodes", "1",
						"--node-type", "trn1.2xlarge",
						"--node-zones", nodeZones,
						"--version", params.Version,
					)
					Expect(cmd).To(RunSuccessfully())
				})
				It("should install the neuron device plugin", func() {
					_, err := newClientSet(clusterWithoutPlugin).AppsV1().DaemonSets("kube-system").Get(context.TODO(), "neuron-device-plugin-daemonset", metav1.GetOptions{})
					Expect(err).ShouldNot(HaveOccurred())
				})
			})
		})
	})
})

var _ = AfterSuite(func() {
	params.DeleteClusters()
	gexec.KillAndWait(30 * time.Minute)
	if params.KubeconfigTemp {
		os.Remove(params.KubeconfigPath)
	}
	os.RemoveAll(params.TestDirectory)
})

func newClientSet(name string) kubernetes.Interface {
	cfg := &api.ClusterConfig{
		Metadata: &api.ClusterMeta{
			Name:   name,
			Region: params.Region,
		},
	}
	ctx := context.Background()
	ctl, err := eks.New(ctx, &api.ProviderConfig{Region: params.Region}, cfg)
	Expect(err).NotTo(HaveOccurred())

	err = ctl.RefreshClusterStatus(ctx, cfg)
	Expect(err).ShouldNot(HaveOccurred())

	clientSet, err := ctl.NewStdClientSet(cfg)
	Expect(err).ShouldNot(HaveOccurred())
	return clientSet
}

func getAvailabilityZones(ctx context.Context, ec2API awsapi.EC2) (string, string) {
	// Trn1 instance types are only available in one AZ per region at this time
	// TODO: dynamically discover zones as part of the core codebase
	trnInstanceZones := map[string]string{
		"us-west-2": "usw2-az4",
		"us-east-1": "use1-az5",
	}

	input := &ec2.DescribeAvailabilityZonesInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("region-name"),
				Values: []string{params.Region},
			}, {
				Name:   aws.String("state"),
				Values: []string{string(ec2types.AvailabilityZoneStateAvailable)},
			}, {
				Name:   aws.String("zone-type"),
				Values: []string{string(ec2types.LocationTypeAvailabilityZone)},
			},
		},
	}

	// Get all zones for the region
	output, err := ec2API.DescribeAvailabilityZones(ctx, input)
	Expect(err).NotTo(HaveOccurred())
	zones := output.AvailabilityZones

	zoneMap := map[string]struct{}{}
	var nodeZones []string

	// Add the zones to the zoneMap where the instance type is supported
	for _, zone := range zones {
		if *zone.ZoneId == trnInstanceZones[params.Region] {
			nodeZones = append(nodeZones, *zone.ZoneName)
			zoneMap[*zone.ZoneName] = struct{}{}
		}
	}

	// If we have fewer than the minimum required number of availability zones
	// then backfill clusterZones to get to MinRequiredAvailabilityZones
	for i := 0; i < len(zones) && len(zoneMap) < api.MinRequiredAvailabilityZones; i++ {
		zoneMap[*zones[i].ZoneName] = struct{}{}
	}

	var clusterZones []string
	// convert this back to a slice of strings
	for zone := range zoneMap {
		clusterZones = append(clusterZones, zone)
	}

	return strings.Join(nodeZones, ","), strings.Join(clusterZones, ",")
}
