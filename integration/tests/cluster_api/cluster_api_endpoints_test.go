// +build integration

package cluster_api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils"

	awseks "github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("capi-endpoints")
}

func TestSuite(t *testing.T) {
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
	cfg.VPC.ClusterEndpoints.PrivateAccess = &privateAccess
	cfg.VPC.ClusterEndpoints.PublicAccess = &publicAccess
}

func setMetadata(cfg *api.ClusterConfig, name, region string) {
	cfg.Metadata.Name = name
	cfg.Metadata.Region = region
}

var _ = Describe("(Integration) Create and Update Cluster with Endpoint Configs", func() {

	type endpointAccessCase struct {
		Name    string
		Private bool
		Public  bool
		Type    string
		Fails   bool
	}

	DescribeTable("Can create/update Cluster Endpoint Access",
		func(e endpointAccessCase) {
			//create clusterconfig
			cfg := api.NewClusterConfig()
			clName := params.NewClusterName(e.Name)
			setEndpointConfig(cfg, e.Private, e.Public)
			setMetadata(cfg, clName, params.Region)

			// create and populate config file from clusterconfig
			bytes, err := json.Marshal(cfg)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(bytes)).ToNot(BeZero())
			tmpfile, err := ioutil.TempFile("", "clusterendpointtests")
			Expect(err).ToNot(HaveOccurred())

			defer os.Remove(tmpfile.Name())

			_, err = tmpfile.Write(bytes)
			Expect(err).ToNot(HaveOccurred())
			err = tmpfile.Close()
			Expect(err).ToNot(HaveOccurred())

			// create cluster with config file
			if e.Type == createCluster {
				cmd := params.EksctlCreateCmd.WithArgs(
					"cluster",
					"--verbose", "2",
					"--config-file", tmpfile.Name(),
					"--without-nodegroup",
				).WithoutArg("--region", params.Region)
				if e.Fails {
					Expect(cmd).ShouldNot(RunSuccessfully())
					return
				}
				Expect(cmd).Should(RunSuccessfully())
				awsSession := NewSession(params.Region)
				Eventually(awsSession, timeOutSeconds, pollInterval).Should(
					HaveExistingCluster(clName, awseks.ClusterStatusActive, params.Version))
			} else if e.Type == updateCluster {
				utilsCmd := params.EksctlUtilsCmd.
					WithTimeout(timeOutSeconds*time.Second).
					WithArgs(
						"update-cluster-endpoints",
						"--name", clName,
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
				// nned to update public access to allow access to delete when it isn't allowed
				if e.Public == false {
					utilsCmd := params.EksctlUtilsCmd.
						WithTimeout(timeOutSeconds*time.Second).WithArgs(
						"update-cluster-endpoints",
						"--name", clName,
						fmt.Sprintf("--public-access=%v", true),
						fmt.Sprintf("--approve"),
					)
					Expect(utilsCmd).Should(RunSuccessfully())
				}
				deleteCmd := params.EksctlDeleteCmd.WithArgs(
					"cluster",
					"--name", clName,
				)
				Expect(deleteCmd).Should(RunSuccessfully())
				awsSession := NewSession(params.Region)
				Eventually(awsSession, timeOutSeconds, pollInterval).
					ShouldNot(HaveExistingCluster(clName, awseks.ClusterStatusActive, params.Version))
			}
		},
		Entry("Create cluster1, Private=false, Public=true, should succeed", endpointAccessCase{
			Name:    "cluster1",
			Private: false,
			Public:  true,
			Type:    createCluster,
			Fails:   false,
		}),
		Entry("Create cluster2, Private=true, Public=false, should not succeed", endpointAccessCase{
			Name:    "cluster2",
			Private: true,
			Public:  false,
			Type:    createCluster,
			Fails:   true,
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
