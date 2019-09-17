// +build integration

package integration_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	awseks "github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

const (
	createCluster    = `CreateCluster`
	updateCluster    = `UpdateCluster`
	endpointPubTmpl  = `EndpointPublicAccess: %v`
	endpointPrivTmpl = `EndpointPrivateAccess: %v`
)

var False = false

func setEndpointConfig(cfg *api.ClusterConfig, privateAccess, publicAccess bool) {
	cfg.VPC.ClusterEndpoints.PrivateAccess = &privateAccess
	cfg.VPC.ClusterEndpoints.PublicAccess = &publicAccess
}

func generateName(prefix string) string {
	if clusterName == "" {
		clusterName = cmdutils.ClusterName("", "")
	}
	return fmt.Sprintf("%v-%v", prefix, clusterName)
}

func setMetadata(cfg *api.ClusterConfig, name, region string) {
	cfg.Metadata.Name = name
	cfg.Metadata.Region = region
}

func printSeparator(op, clname string) {
	fmt.Println("\n==================================================================")
	fmt.Printf("%v cluster %v\n", op, clname)
	fmt.Println("==================================================================")
}

var _ = Describe("(Integration) Create and Update Cluster with Endpoint Configs", func() {
	type EndpointAccessCases struct {
		Name    string
		Private bool
		Public  bool
		Type    string // createCluster or updateCluster
		Output  string
		Error   error
		Delete  bool
	}

	FDescribeTable("Can create/update Cluster Endpoint Access",
		func(e EndpointAccessCases) {
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
				printSeparator("Creating", clName)
				cmd := eksctlCreateCmd.WithArgs(
					"cluster",
					"--verbose", "2",
					"--config-file", tmpfile.Name(),
					"--without-nodegroup",
				)
				if e.Error != nil {
					Expect(cmd).ShouldNot(RunSuccessfully())
					return
				}
				Expect(cmd).Should(RunSuccessfully())
				awsSession := NewSession(region)
				Eventually(awsSession, timeOut, pollInterval).Should(
					HaveExistingCluster(clName, awseks.ClusterStatusActive, version))
			} else if e.Type == updateCluster {
				printSeparator("Updating", clName)
				utilsCmd := eksctlUtilsCmd.WithArgs(
					"update-cluster-api-access",
					"--name", clName,
					fmt.Sprintf("--private-access=%v", e.Private),
					fmt.Sprintf("--public-access=%v", e.Public),
					fmt.Sprintf("--approve"),
				)
				if e.Error != nil {
					Expect(utilsCmd).ShouldNot(RunSuccessfully())
					return
				}
				Expect(utilsCmd).Should(RunSuccessfully())
			}
			printSeparator("Getting", clName)
			getCmd := eksctlGetCmd.WithArgs(
				"cluster",
				"--name", clName,
				"-o", "yaml",
			)
			Expect(getCmd).To(RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring(endpointPubTmpl, e.Public)),
				ContainElement(ContainSubstring(endpointPrivTmpl, e.Private)),
			))
			if e.Delete {
				printSeparator("Deleting", clName)
				// nned to update public access to allow access to delete when it isn't allowed
				if e.Public == false {
					utilsCmd := eksctlUtilsCmd.WithArgs(
						"update-cluster-api-access",
						"--name", clName,
						fmt.Sprintf("--public-access=%v", true),
						fmt.Sprintf("--approve"),
					)
					Expect(utilsCmd).Should(RunSuccessfully())
				}
				deleteCmd := eksctlDeleteCmd.WithArgs(
					"cluster",
					"--name", clName,
					"--wait",
				)
				Expect(deleteCmd).Should(RunSuccessfully())
				Expect(getCmd).ShouldNot(RunSuccessfully())
			}
		},
		Entry("Create cluster1 with Private=false, Public=true", EndpointAccessCases{
			Name:    "cluster1",
			Private: false,
			Public:  true,
			Type:    createCluster,
			Error:   nil,
			Delete:  false,
		}),
		Entry("Create cluster2 with Private=true, Public=false", EndpointAccessCases{
			Name:    "cluster2",
			Private: true,
			Public:  false,
			Type:    createCluster,
			Error:   errors.New(api.PrivateOnlyUseUtilsMsg()),
			Delete:  true, // In case the cluster gets create because of a bug
		}),
		Entry("Create cluster 3 Private=true, Public=true", EndpointAccessCases{
			Name:    "cluster3",
			Private: true,
			Public:  true,
			Type:    createCluster,
			Error:   nil,
			Delete:  false,
		}),
		Entry("Create cluster4 with Private=false, Public=false (should error)", EndpointAccessCases{
			Name:    "cluster4",
			Private: false,
			Public:  false,
			Type:    createCluster,
			Error:   errors.New(api.NoAccessMsg(
				&api.ClusterEndpoints{PrivateAccess: &False, PublicAccess: &False,},
			)),
			Delete:  true, // In case the cluster gets created because of a bug.
		}),
		Entry("Update cluster1 to Private=true, Public=false", EndpointAccessCases{
			Name:    "cluster1",
			Private: true,
			Public:  false,
			Type:    updateCluster,
			Error:   nil,
			Delete:  true,
		}),
		Entry("Update cluster3 to Private=true, Public=false", EndpointAccessCases{
			Name:    "cluster3",
			Private: true,
			Public:  false,
			Type:    updateCluster,
			Error:   nil,
			Delete:  false,
		}),
		Entry("Update cluster3 to Private=false, Public=false (should error)", EndpointAccessCases{
			Name:    "cluster3",
			Private: false,
			Public:  false,
			Type:    updateCluster,
			Error:   errors.New("Unable to make requested changes.  Either public access or private access must be enabled"),
			Delete:  false,
		}),
		Entry("Update cluster3 to Private=false, Public=false (should error)", EndpointAccessCases{
			Name:    "cluster3",
			Private: false,
			Public:  true,
			Type:    updateCluster,
			Error:   nil,
			Delete:  true,
		}),
	)
})
