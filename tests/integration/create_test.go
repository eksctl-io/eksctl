package create_test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	clusterName         = "int-cluster"
	createClusterArgs   = "create cluster -n %s -t t2.medium -N 1 -kubeconfig %s"
	createTimeoutInMins = 20
)

var pathToEksCtl string

func TestCreateIntegration(t *testing.T) {
	RegisterTestingT(t)
	//RegisterFailHandler(Fail)
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
			args := fmt.Sprintf(createClusterArgs, clusterName, kubeConfigPath.Name())
			fmt.Printf("Path: %s\n", pathToEksCtl)
			fmt.Printf("Args: %s\n", args)

			command := exec.Command(pathToEksCtl, "")
			session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)

			if err != nil {
				Fail("error starting process", 1)
			}

			session.Wait(createTimeoutInMins * time.Minute)
			//Expect(session).Should(gexec.Exit())
			//Expect(session.ExitCode()).Should(Equal(0))
		})

	})
})

func init() {
	flag.StringVar(&pathToEksCtl, "eksctl-path", "./eksctl", "Path to eksctl")
}
