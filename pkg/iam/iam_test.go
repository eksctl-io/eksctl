package iam_test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/iam"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var (
	ng1 = &api.NodeGroup{
		Name: "ng1",
		IAM: &api.NodeGroupIAM{
			InstanceRoleARN:    "arn:aws:iam::122333:role/eksctl-cluster-ng1-NodeInstanceRole-ASDF",
			InstanceProfileARN: "arn:aws:iam::122333:instance-profile/eksctl-cluster-ng1-instance-profile",
		},
	}

	mng1 = &api.ManagedNodeGroup{
		Name: "mng1",
		IAM: &api.NodeGroupIAM{
			InstanceRoleARN:    "arn:aws:iam::122333:role/eksctl-cluster-mng1-NodeInstanceRole-BLAH",
			InstanceProfileARN: "arn:aws:iam::122333:instance-profile/eksctl-cluster-mng1-instance-profile",
		},
	}
)

var _ = Describe("Get IAM from Node Group", func() {

	p := mockprovider.NewMockProvider()

	cfnDescriptions := []*manager.Stack{
		{
			StackStatus: aws.String(cloudformation.StackStatusCreateComplete),
			StackName:   aws.String(ng1.Name),
			RoleARN:     aws.String(ng1.IAM.InstanceRoleARN),
			Tags: []*cloudformation.Tag{
				{Key: aws.String(api.NodeGroupNameTag), Value: aws.String(ng1.Name)},
				{Key: aws.String(api.NodeGroupTypeTag), Value: aws.String(string(api.NodeGroupTypeUnmanaged))},
			},
			Outputs: []*cloudformation.Output{
				{OutputKey: aws.String("InstanceRoleARN"), OutputValue: aws.String(ng1.IAM.InstanceRoleARN)},
			},
		},
		{
			StackStatus: aws.String(cloudformation.StackStatusCreateComplete),
			StackName:   aws.String(mng1.Name),
			RoleARN:     aws.String(mng1.IAM.InstanceRoleARN),
			Tags: []*cloudformation.Tag{
				{Key: aws.String(api.NodeGroupNameTag), Value: aws.String(mng1.Name)},
				{Key: aws.String(api.NodeGroupTypeTag), Value: aws.String(string(api.NodeGroupTypeManaged))},
			},
			Outputs: []*cloudformation.Output{
				{OutputKey: aws.String("InstanceRoleARN"), OutputValue: aws.String(mng1.IAM.InstanceRoleARN)},
			},
		},
	}

	It("should get IAM configuration from node group", func() {
		ng := &api.NodeGroup{Name: "ng1"}
		err := iam.UseFromNodeGroup(p, cfnDescriptions[0], ng)
		Expect(err).ToNot(HaveOccurred())
		Expect(ng.IAM.InstanceRoleARN).To(Equal(ng1.IAM.InstanceRoleARN))
	})

	It("should get IAM configuration from managed node group", func() {
		mng := &api.ManagedNodeGroup{Name: "mng1"}
		err := iam.UseFromManagedNodeGroup(p, cfnDescriptions[1], mng)
		Expect(err).ToNot(HaveOccurred())
		Expect(mng.IAM.InstanceRoleARN).To(Equal(mng1.IAM.InstanceRoleARN))
	})
})
