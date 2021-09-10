//go:build integration
// +build integration

package caching

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("cache")
}

func TestCredentialsCaching(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("", func() {
	var (
		clusterName string
	)
	BeforeSuite(func() {
		clusterName = params.NewClusterName("cache")
		cmd := params.EksctlCreateCmd.WithArgs(
			"cluster",
			"--name", clusterName,
			"--without-nodegroup",
		).WithoutArg("--region", params.Region)
		Expect(cmd).Should(RunSuccessfully())
	})
	AfterSuite(func() {
		params.DeleteClusters()
		gexec.KillAndWait()
	})
	FContext("creating a new raw client", func() {
		It("should not cache the credentials", func() {
			a, err := eks.New(&api.ProviderConfig{
				CloudFormationRoleARN: "arn",
				Region:                "us-west-2",
				Profile:               "profile",
			}, &api.ClusterConfig{
				Metadata: &api.ClusterMeta{
					Name:   "test-cluster",
					Region: "us-west-2",
				},
			})
			Expect(err).ToNot(HaveOccurred())
			fmt.Println(a.Provider)
		})
		When("credential caching is enabled", func() {
			var tmp string
			BeforeEach(func() {
				tmp, err := ioutil.TempDir("", "caching_creds")
				Expect(err).NotTo(HaveOccurred())
				os.Setenv(eks.EksctlCacheFilenameEnvName, filepath.Join(tmp, "credentials.yaml"))
				os.Setenv(eks.EksctlGlobalEnableCachingEnvName, "1")
			})
			AfterEach(func() {
				_ = os.RemoveAll(tmp)
			})
			It("should cache the credentials for the given profile", func() {

			})
		})
	})
})
