//go:build integration
// +build integration

//revive:disable Not changing package name
package cloudwatch_logging

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	clusterutils "github.com/weaveworks/eksctl/integration/utilities/cluster"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	if err := api.Register(); err != nil {
		panic(errors.Wrap(err, "unexpected error registering API scheme"))
	}
	params = tests.NewParams("cloudwatch")
}

func TestCloudWatchLogging(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("(Integration) [CloudWatch Logging test]", func() {
	Describe("CloudWatch logging", func() {
		It("should create a cluster with CloudWatch logging enabled and log retention set", func() {
			cmd := params.EksctlCreateCmd.
				WithArgs(
					"cluster",
					"--config-file=-",
					"--verbose=4",
				).
				WithoutArg("--region", params.Region).
				WithStdin(clusterutils.ReaderFromFile(params.ClusterName, params.Region, "testdata/cloudwatch-cluster.yaml"))

			Expect(cmd).To(RunSuccessfullyWithOutputString(ContainSubstring("set log retention to 545 days for CloudWatch logging")))

			cloudWatchLogs := cloudwatchlogs.NewFromConfig(NewConfig(params.Region))
			logGroups, err := cloudWatchLogs.DescribeLogGroups(context.Background(), &cloudwatchlogs.DescribeLogGroupsInput{
				LogGroupNamePrefix: aws.String(fmt.Sprintf("/aws/eks/%s/cluster", params.ClusterName)),
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(logGroups.LogGroups).To(HaveLen(1))
			Expect(*logGroups.LogGroups[0].RetentionInDays).To(Equal(int32(545)))
		})
	})
})

var _ = AfterSuite(func() {
	params.DeleteClusters()
})
