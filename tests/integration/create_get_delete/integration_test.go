// +build integration

package create_get_delete

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	. "github.com/weaveworks/eksctl/pkg/testutils/matchers"

	"github.com/aws/aws-sdk-go/aws"
	awssess "github.com/aws/aws-sdk-go/aws/session"
	awseks "github.com/aws/aws-sdk-go/service/eks"
)

const (
	clusterName         = "int-cluster"
	createTimeoutInMins = 20
	eksRegion           = "us-west-2"
)

var (
	pathToEksCtl string
	skipCreation bool
)

func TestCreateIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration - Create Suite")
}

var _ = Describe("Create (Integration)", func() {
	var (
		kubeConfigPath *os.File
	)

	AfterSuite(func() {
		gexec.KillAndWait()
	})

	BeforeEach(func() {
		kubeConfigPath, _ = ioutil.TempFile("", "")
	})

	AfterEach(func() {
		os.Remove(kubeConfigPath.Name())
	})

	Describe("when creating a cluster with 1 node", func() {
		var (
			err     error
			session *gexec.Session
		)

		It("should not return an error", func() {
			if skipCreation {
				fmt.Printf("Creation test skip: %t\n", skipCreation)
				return
			}

			args := []string{"create", "cluster", "-n", clusterName, "-t", "t2.medium", "-N", "1", "-r", eksRegion, "--kubeconfig", kubeConfigPath.Name()}
			//fmt.Printf("Path: %s\n", pathToEksCtl)
			//fmt.Printf("Args: %s\n", args)

			command := exec.Command(pathToEksCtl, args...)
			session, err = gexec.Start(command, os.Stdout, os.Stdout)

			if err != nil {
				Fail("error starting process", 1)
			}

			session.Wait(createTimeoutInMins * time.Minute)
			//Expect(session).Should(gexec.Exit())
			Expect(session.ExitCode()).Should(Equal(0))
		})

		It("should have created an EKS cluster", func() {
			session := newSession()
			Expect(session).To(HaveEksCluster(clusterName, awseks.ClusterStatusActive, "1.10"))
		})

		It("should have the required cloudformation stacks", func() {
			session := newSession()

			Expect(session).To(HaveCfnStack(fmt.Sprintf("EKS-%s-VPC", clusterName)))
			Expect(session).To(HaveCfnStack(fmt.Sprintf("EKS-%s-ControlPlane", clusterName)))
			Expect(session).To(HaveCfnStack(fmt.Sprintf("EKS-%s-ServiceRole", clusterName)))
			Expect(session).To(HaveCfnStack(fmt.Sprintf("EKS-%s-DefaultNodeGroup", clusterName)))
		})

	})
})

func newSession() *awssess.Session {
	config := aws.NewConfig()
	config = config.WithRegion(eksRegion)
	opts := awssess.Options{
		Config: *config,
	}
	return awssess.Must(awssess.NewSessionWithOptions(opts))
}

func init() {
	flag.StringVar(&pathToEksCtl, "eksctl-path", "./eksctl", "Path to eksctl")
	flag.BoolVar(&skipCreation, "skip-creation", false, "Skip the creation tests. Useful for debugging the tests")
}
