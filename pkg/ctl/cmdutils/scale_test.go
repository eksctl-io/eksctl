package cmdutils

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

type scaleNodeGroupCase struct {
	name  string
	error error
}

var _ = Describe("scale node group config file loader", func() {
	newCmd := func() *cobra.Command {
		return &cobra.Command{
			Use: "test",
			Run: func(_ *cobra.Command, _ []string) {},
		}
	}

	DescribeTable("create nodegroup successfully",
		func(params scaleNodeGroupCase) {
			cmd := &Cmd{
				CobraCommand:      newCmd(),
				ClusterConfigFile: "test_data/scale-ng-test.yaml",
				ClusterConfig:     api.NewClusterConfig(),
				ProviderConfig:    &api.ProviderConfig{},
				NameArg:           params.name,
			}

			ng := api.NewNodeGroup()
			err := NewScaleNodeGroupLoader(cmd, ng).Load()
			if params.error != nil {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(params.error.Error()))
			} else {
				Expect(err).ToNot(HaveOccurred())
			}
		},
		Entry("one node group matched", scaleNodeGroupCase{
			name: "ng-all-details",
		}),
		Entry("no node group matched", scaleNodeGroupCase{
			name:  "123123",
			error: fmt.Errorf("node group 123123 not found"),
		}),
		Entry("with no desired capacity", scaleNodeGroupCase{
			name:  "ng-no-desired-capacity",
			error: fmt.Errorf("number of nodes must be 0 or greater"),
		}),
		Entry("with no minSize and no maxSize", scaleNodeGroupCase{
			name: "ng-no-min-max",
		}),
		Entry("ng with minSize", scaleNodeGroupCase{
			name: "ng-with-min",
		}),
		Entry("ng with wrong value for minSize", scaleNodeGroupCase{
			name:  "ng-with-wrong-min",
			error: fmt.Errorf("minimum number of nodes must be less than or equal to number of nodes"),
		}),
		Entry("ng with maxSize", scaleNodeGroupCase{
			name: "ng-with-max",
		}),
		Entry("ng with wrong value for maxSize", scaleNodeGroupCase{
			name:  "ng-with-wrong-max",
			error: fmt.Errorf("maximum number of nodes must be greater than or equal to number of nodes"),
		}),
		Entry("ng with desired nodes outside [minSize, maxSize]", scaleNodeGroupCase{
			name:  "ng-with-wrong-desired",
			error: fmt.Errorf("number of nodes must be within range of min nodes and max nodes"),
		}),
	)
})
