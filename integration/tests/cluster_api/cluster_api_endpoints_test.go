//go:build integration
// +build integration

//revive:disable Not changing package name
package cluster_api

import (
	"fmt"
	"testing"
	"time"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	clusterutils "github.com/weaveworks/eksctl/integration/utilities/cluster"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("capi")
}

func TestClusterAPI(t *testing.T) {
	testutils.RegisterAndRun(t)
}

const (
	createCluster    = `Create`
	updateCluster    = `Update`
	deleteCluster    = `Delete`
	endpointPubTmpl  = `EndpointPublicAccess: %v`
	endpointPrivTmpl = `EndpointPrivateAccess: %v`
	pollInterval     = 15   //seconds
	timeOutSeconds   = 1200 // 20 minutes
)

func setEndpointConfig(cfg *api.ClusterConfig, privateAccess, publicAccess bool) {
	cfg.VPC.ClusterEndpoints = api.ClusterEndpointAccessDefaults()
	cfg.VPC.ClusterEndpoints.PrivateAccess = &privateAccess
	cfg.VPC.ClusterEndpoints.PublicAccess = &publicAccess
}

func setMetadata(cfg *api.ClusterConfig, name, region string) {
	cfg.Metadata.Name = name
	cfg.Metadata.Region = region
}

var _ = Describe("(Integration) Create and Update Cluster with Endpoint Configs", func() {

	clusterNames := map[string]string{}

	type endpointAccessCase struct {
		Name    string
		Private bool
		Public  bool
		Type    string
		Fails   bool
	}

	params.LogStacksEventsOnFailure()

	DescribeTable("Can create/update Cluster Endpoint Access",
		func(e endpointAccessCase) {
			//create clusterconfig
			cfg := api.NewClusterConfig()
			// get or generate unique cluster name:
			clName, ok := clusterNames[e.Name]
			if !ok {
				clName = params.NewClusterName(e.Name)
				clusterNames[e.Name] = clName
			}
			params.ClusterName = clName
			setEndpointConfig(cfg, e.Private, e.Public)
			setMetadata(cfg, clName, params.Region)

			// create cluster with config file
			if e.Type == createCluster {
				cmd := params.EksctlCreateCmd.WithArgs(
					"cluster",
					"--verbose", "2",
					"--config-file", "-",
					"--without-nodegroup",
				).
					WithoutArg("--region", params.Region).
					WithStdin(clusterutils.Reader(cfg))

				if e.Fails {
					Expect(cmd).ShouldNot(RunSuccessfully())
					return
				}
				Expect(cmd).Should(RunSuccessfully())
				awsSession := NewConfig(params.Region)
				Eventually(awsSession, timeOutSeconds, pollInterval).Should(
					HaveExistingCluster(clName, string(ekstypes.ClusterStatusActive), params.Version))
			} else if e.Type == updateCluster {
				utilsCmd := params.EksctlUtilsCmd.
					WithTimeout(timeOutSeconds*time.Second).
					WithArgs(
						"update-cluster-endpoints",
						"--cluster", clName,
						fmt.Sprintf("--private-access=%v", e.Private),
						fmt.Sprintf("--public-access=%v", e.Public),
						"--approve")
				if e.Fails {
					Expect(utilsCmd).ShouldNot(RunSuccessfully())
					return
				}
				Expect(utilsCmd).Should(RunSuccessfully())
			}
			getCmd := params.EksctlGetCmd.WithArgs(
				"cluster",
				"--name", clName,
				"-o", "yaml",
			)
			Expect(getCmd).To(RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring(endpointPubTmpl, e.Public)),
				ContainElement(ContainSubstring(endpointPrivTmpl, e.Private)),
			))
			if e.Type == deleteCluster {
				// need to update public access to allow access to delete when it isn't allowed
				if e.Public == false {
					utilsCmd := params.EksctlUtilsCmd.
						WithTimeout(timeOutSeconds*time.Second).WithArgs(
						"update-cluster-endpoints",
						"--cluster", clName,
						fmt.Sprintf("--public-access=%v", true),
						"--approve",
					)
					Expect(utilsCmd).Should(RunSuccessfully())
				}
				deleteCmd := params.EksctlDeleteCmd.WithArgs(
					"cluster",
					"--name", clName,
				)
				Expect(deleteCmd).Should(RunSuccessfully())
				awsSession := NewConfig(params.Region)
				Eventually(awsSession, timeOutSeconds, pollInterval).
					ShouldNot(HaveExistingCluster(clName, string(ekstypes.ClusterStatusActive), params.Version))
			}
		},
		Entry("Create cluster1, Private=false, Public=true, should succeed", endpointAccessCase{
			Name:    "cluster1",
			Private: false,
			Public:  true,
			Type:    createCluster,
			Fails:   false,
		}),
		Entry("Create cluster3, Private=true, Public=true, should succeed", endpointAccessCase{
			Name:    "cluster3",
			Private: true,
			Public:  true,
			Type:    createCluster,
			Fails:   false,
		}),
		Entry("Create cluster4, Private=false, Public=false, should not succeed", endpointAccessCase{
			Name:    "cluster4",
			Private: false,
			Public:  false,
			Type:    createCluster,
			Fails:   true,
		}),
		Entry("Update cluster1 to Private=true, Public=false, should succeed", endpointAccessCase{
			Name:    "cluster1",
			Private: true,
			Public:  false,
			Type:    updateCluster,
			Fails:   false,
		}),
		Entry("Update cluster3 to Private=true, Public=false, should succeed", endpointAccessCase{
			Name:    "cluster3",
			Private: true,
			Public:  false,
			Type:    updateCluster,
			Fails:   false,
		}),
		Entry("Update cluster3 to Private=false, Public=false, should not succeed", endpointAccessCase{
			Name:    "cluster3",
			Private: false,
			Public:  false,
			Type:    updateCluster,
			Fails:   true,
		}),
		Entry("Update cluster3 to Private=false, Public=true, should succeed", endpointAccessCase{
			Name:    "cluster3",
			Private: false,
			Public:  true,
			Type:    updateCluster,
			Fails:   false,
		}),
		Entry("Delete cluster1, should succeed (test case updates access)", endpointAccessCase{
			Name:    "cluster1",
			Private: true,
			Public:  false,
			Type:    deleteCluster,
			Fails:   false,
		}),
		Entry("Delete cluster3, succeed", endpointAccessCase{
			Name:    "cluster3",
			Private: false,
			Public:  true,
			Type:    deleteCluster,
			Fails:   false,
		}),
	)
})

var _ = AfterSuite(func() {
	params.DeleteClusters()
})
