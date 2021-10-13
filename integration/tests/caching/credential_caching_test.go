//go:build integration
// +build integration

package caching

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	"github.com/weaveworks/eksctl/pkg/credentials"
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
	Context("accessing cluster related information", func() {
		When("credential caching is disabled", func() {
			var tmp string
			BeforeEach(func() {
				tmp, err := os.MkdirTemp("", "caching_creds")
				Expect(err).NotTo(HaveOccurred())
				_ = os.Setenv(credentials.EksctlCacheFilenameEnvName, filepath.Join(tmp, "credentials.yaml"))
				// Make sure the environment property is not set on the running environment.
				_ = os.Setenv(credentials.EksctlGlobalEnableCachingEnvName, "")
			})
			AfterEach(func() {
				_ = os.RemoveAll(tmp)
			})
			It("should not cache the credentials", func() {
				cmd := params.EksctlGetCmd.WithArgs(
					"cluster",
				).WithoutArg("--region", params.Region)
				Expect(cmd).Should(RunSuccessfully())

				_, err := os.Stat(filepath.Join(tmp, "credentials.yaml"))
				Expect(os.IsNotExist(err)).To(BeTrue())
			})
		})
		// Note: This proves to be a challenge since normal providers like, static and file providers, do not support
		// expiry, therefore, they cannot be cached. This requires an AssumeRole or an EC2 role provider.
		// For now, this part is tested via unit tests.
		XWhen("credential caching is enabled", func() {
			var tmp string
			BeforeEach(func() {
				tmp, err := os.MkdirTemp("", "caching_creds")
				Expect(err).NotTo(HaveOccurred())
				_ = os.Setenv(credentials.EksctlCacheFilenameEnvName, filepath.Join(tmp, "credentials.yaml"))
				_ = os.Setenv(credentials.EksctlGlobalEnableCachingEnvName, "1")
			})
			AfterEach(func() {
				_ = os.RemoveAll(tmp)
			})
			It("should cache the credentials", func() {
				cmd := params.EksctlGetCmd.WithArgs(
					"cluster",
				).WithoutArg("--region", params.Region)
				Expect(cmd).Should(RunSuccessfully())

				content, err := os.ReadFile(filepath.Join(tmp, "credentials.yaml"))
				Expect(err).NotTo(HaveOccurred())
				Expect(content).NotTo(BeEmpty())
			})
		})
	})
})
