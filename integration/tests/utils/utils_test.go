//go:build integration
// +build integration

package utils

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/integration/tests"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var params *tests.Params

func init() {
	testing.Init()
	params = tests.NewParams("schema")
}

func TestUtils(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = BeforeSuite(func() {
})

var _ = Describe("Utils", func() {
	Context("schema", func() {
		It("displays the schema", func() {
			cmd := params.EksctlUtilsCmd.WithArgs("schema").WithoutArg("--region", params.Region)
			session := cmd.Run()
			Expect(session.ExitCode()).To(BeZero())
			Expect(string(session.Out.Contents())).To(Equal(api.SchemaJSON))
		})
	})
})
