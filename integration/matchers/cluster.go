package matchers

import (
	"strings"

	"github.com/weaveworks/eksctl/integration/runner"

	. "github.com/onsi/gomega"
)

type GetClusterOutput struct {
	ClusterName   string
	Region        string
	EksctlCreated string
}

// AssertContainsCluster asserts that the output of the specified command contains the specified cluster
func AssertContainsCluster(cmd runner.Cmd, out GetClusterOutput) {
	Expect(cmd).To(runner.RunSuccessfullyWithOutputStringLines(
		ContainElement(WithTransform(func(line string) GetClusterOutput {
			fields := strings.Fields(line)
			const expectedColumns = 3
			Expect(fields).To(HaveLen(expectedColumns), "Expected `get clusters` to return %d columns", expectedColumns)
			return GetClusterOutput{
				ClusterName:   fields[0],
				Region:        fields[1],
				EksctlCreated: fields[2],
			}
		}, Equal(out))),
	))

}
