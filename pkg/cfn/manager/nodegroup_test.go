package manager

import (
	"github.com/aws/aws-sdk-go/aws"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = Describe("StackCollection NodeGroup", func() {
	Describe("GetNodeGroupType", func() {
		createTags := func(tags map[string]string) []*cfn.Tag {
			cfnTags := make([]*cfn.Tag, 0)
			for k, v := range tags {
				cfnTags = append(cfnTags, &cfn.Tag{
					Key:   aws.String(k),
					Value: aws.String(v),
				})
			}
			return cfnTags
		}

		DescribeTable("with tag for the nodegroup type", func(inputTags map[string]string, expectedType api.NodeGroupType) {
			ngType, err := GetNodeGroupType(createTags(inputTags))

			if expectedType == "" {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(ngType).To(Equal(expectedType))
			}
		},

			Entry("finds the type of a managed nodegroup",
				map[string]string{
					api.NodeGroupNameTag: "mng-1",
					api.NodeGroupTypeTag: "managed",
				},
				api.NodeGroupTypeManaged),

			Entry("finds the type of an un-managed nodegroup",
				map[string]string{
					api.NodeGroupNameTag: "ng-1",
					api.NodeGroupTypeTag: "unmanaged",
				},
				api.NodeGroupTypeUnmanaged),

			Entry("finds the type of a legacy un-managed nodegroup",
				map[string]string{
					api.OldNodeGroupNameTag: "ng-1",
					api.NodeGroupTypeTag:    "unmanaged",
				},
				api.NodeGroupTypeUnmanaged),

			Entry("finds the type of a legacy un-managed nodegroup",
				map[string]string{
					api.OldNodeGroupIDTag: "ng-1",
					api.NodeGroupTypeTag:  "unmanaged",
				},
				api.NodeGroupTypeUnmanaged),

			Entry("doesn't return the type if the stack tags don't contain any ng name tag",
				map[string]string{
					"some-other-tag": "ng-1",
					"name":           "ng-1",
				},
				api.NodeGroupType("")),
		)
		DescribeTable("for legacy ngs without tag for the type", func(inputTags map[string]string, expectedType api.NodeGroupType) {
			ngType, err := GetNodeGroupType(createTags(inputTags))

			if expectedType == "" {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(ngType).To(Equal(expectedType))
			}
		},

			Entry("legacy ngs with old name tags are un-managed by default",
				map[string]string{
					api.NodeGroupNameTag: "ng-1",
				},
				api.NodeGroupTypeUnmanaged),

			Entry("legacy ngs with old name tags are un-managed by default",
				map[string]string{
					api.OldNodeGroupNameTag: "ng-1",
				},
				api.NodeGroupTypeUnmanaged),

			Entry("legacy ngs with old name tag group-id are un-managed by default",
				map[string]string{
					api.OldNodeGroupIDTag: "ng-1",
				},
				api.NodeGroupTypeUnmanaged),

			Entry("doesn't return the type if the stack tags don't contain any ng name tag",
				map[string]string{
					"some-other-tag": "ng-1",
					"name":           "ng-1",
				},
				api.NodeGroupType("")),
		)
	})
})
