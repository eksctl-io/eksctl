// +build integration

package integration_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	awseks "github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/utils/names"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

const (
	createCluster    = `Create`
	updateCluster    = `Update`
	deleteCluster    = `Delete`
	endpointPubTmpl  = `EndpointPublicAccess: %v`
	endpointPrivTmpl = `EndpointPrivateAccess: %v`
)

func setEndpointConfig(cfg *api.ClusterConfig, privateAccess, publicAccess bool) {
	cfg.VPC.ClusterEndpoints.PrivateAccess = &privateAccess
	cfg.VPC.ClusterEndpoints.PublicAccess = &publicAccess
}

func generateName(prefix string) string {
	if clusterName == "" {
		clusterName = names.ForCluster("", "")
	}
	return fmt.Sprintf("%v-%v", prefix, clusterName)
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
			clName := generateName(e.Name)
			setEndpointConfig(cfg, e.Private, e.Public)
			setMetadata(cfg, clName, region)

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
				cmd := eksctlCreateCmd.WithArgs(
					"cluster",
					"--verbose", "2",
					"--config-file", tmpfile.Name(),
					"--without-nodegroup",
				).WithoutArg("--region", region)
				if e.Fails {
					Expect(cmd).ShouldNot(RunSuccessfully())
					return
				}
				Expect(cmd).Should(RunSuccessfully())
				awsSession := NewSession(region)
				Eventually(awsSession, timeOutSeconds, pollInterval).Should(
					HaveExistingCluster(clName, awseks.ClusterStatusActive, version))
			} else if e.Type == updateCluster {
				utilsCmd := eksctlUtilsCmd.
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
			getCmd := eksctlGetCmd.WithArgs(
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
					utilsCmd := eksctlUtilsCmd.
						WithTimeout(timeOutSeconds*time.Second).WithArgs(
						"update-cluster-endpoints",
						"--name", clName,
						fmt.Sprintf("--public-access=%v", true),
						fmt.Sprintf("--approve"),
					)
					Expect(utilsCmd).Should(RunSuccessfully())
				}
				deleteCmd := eksctlDeleteCmd.WithArgs(
					"cluster",
					"--name", clName,
				)
				Expect(deleteCmd).Should(RunSuccessfully())
				awsSession := NewSession(region)
				Eventually(awsSession, timeOutSeconds, pollInterval).
					ShouldNot(HaveExistingCluster(clName, awseks.ClusterStatusActive, version))
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
