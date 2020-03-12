package matchers_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/weaveworks/eksctl/integration/matchers"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

func TestSuite(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("BeNodeGroupsWithNamesWhich", func() {
	It("can marshal the JSON representation of a nodegroup and match on its name", func() {
		Expect(`[
			{
				"StackName": "eksctl-test-cluster-nodegroup-ng-0",
				"Cluster": "test-cluster",
				"Name": "ng-0",
				"MaxSize": 4,
				"MinSize": 4,
				"DesiredCapacity": 4,
				"InstanceType": "m5.xlarge",
				"ImageID": "ami-036f46d54262b5179",
				"CreationTime": "2020-03-10T10:53:05.106Z",
				"NodeInstanceRoleARN": "arn:aws:iam::083751696308:role/eksctl-test-cluster-nodegroup-ng-0-NodeInstanceRole-1IYQ3JS8OKPX1"
			}
		]`).To(BeNodeGroupsWithNamesWhich(
			HaveLen(1),
			ContainElement("ng-0"),
			Not(ContainElement("ng-1")),
		))
	})

	It("can marshal the JSON representation of nodegroups and match on their names", func() {
		Expect(`[
			{
				"Name": "ng-0"
			},
			{
				"Name": "ng-1"
			}
		]`).To(BeNodeGroupsWithNamesWhich(
			HaveLen(2),
			ContainElement("ng-0"),
			ContainElement("ng-1"),
			Not(ContainElement("ng-2")),
		))
	})
})
