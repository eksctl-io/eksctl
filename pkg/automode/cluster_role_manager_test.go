package automode_test

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/weaveworks/eksctl/pkg/automode"
	"github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"
	"github.com/weaveworks/eksctl/pkg/eks/mocksv2"
)

var _ = Describe("Cluster Role Manager", func() {
	var (
		iamClient    *mocksv2.IAM
		stackManager *fakes.FakeStackManager
		roleManager  *automode.ClusterRoleManager
		attached     []string
		detached     []string
	)

	attachedPolicies := func(names ...string) {
		var policies []iamtypes.AttachedPolicy
		for _, n := range names {
			policies = append(policies, iamtypes.AttachedPolicy{PolicyName: aws.String(n)})
		}
		iamClient.On("ListAttachedRolePolicies", mock.Anything, mock.Anything, mock.Anything).
			Return(&iam.ListAttachedRolePoliciesOutput{AttachedPolicies: policies}, nil)
	}

	BeforeEach(func() {
		attached, detached = nil, nil
		iamClient = &mocksv2.IAM{}
		stackManager = &fakes.FakeStackManager{}
		stackManager.GetStackTemplateReturns(`{"Resources": {"ServiceRole": {"Type": "AWS::IAM::Role"}}}`, nil)
		roleManager = &automode.ClusterRoleManager{
			IAMRoleManager:  iamClient,
			StackManager:    stackManager,
			Region:          "us-west-2",
			ClusterRoleName: "cluster-role",
		}
		iamClient.On("AttachRolePolicy", mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				attached = append(attached, *args[1].(*iam.AttachRolePolicyInput).PolicyArn)
			}).Return(&iam.AttachRolePolicyOutput{}, nil)
		iamClient.On("DetachRolePolicy", mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				detached = append(detached, *args[1].(*iam.DetachRolePolicyInput).PolicyArn)
			}).Return(&iam.DetachRolePolicyOutput{}, nil)
		iamClient.On("UpdateAssumeRolePolicy", mock.Anything, mock.Anything).
			Return(&iam.UpdateAssumeRolePolicyOutput{}, nil)
	})

	It("attaches AmazonEKSBlockStoragePolicyV2 when enabling Auto Mode", func() {
		attachedPolicies("AmazonEKSClusterPolicy")
		Expect(roleManager.UpdateRoleForAutoMode(context.Background())).To(Succeed())
		Expect(attached).To(ConsistOf(
			"arn:aws:iam::aws:policy/AmazonEKSComputePolicy",
			"arn:aws:iam::aws:policy/AmazonEKSBlockStoragePolicyV2",
			"arn:aws:iam::aws:policy/AmazonEKSLoadBalancingPolicy",
			"arn:aws:iam::aws:policy/AmazonEKSNetworkingPolicy",
		))
	})

	It("detaches AmazonEKSBlockStoragePolicy attached by older eksctl versions when disabling Auto Mode", func() {
		attachedPolicies("AmazonEKSClusterPolicy", "AmazonEKSComputePolicy", "AmazonEKSBlockStoragePolicy",
			"AmazonEKSLoadBalancingPolicy", "AmazonEKSNetworkingPolicy")
		Expect(roleManager.DeleteAutoModePolicies(context.Background())).To(Succeed())
		Expect(detached).To(ConsistOf(
			"arn:aws:iam::aws:policy/AmazonEKSComputePolicy",
			"arn:aws:iam::aws:policy/AmazonEKSBlockStoragePolicy",
			"arn:aws:iam::aws:policy/AmazonEKSLoadBalancingPolicy",
			"arn:aws:iam::aws:policy/AmazonEKSNetworkingPolicy",
		))
	})

	It("detaches AmazonEKSBlockStoragePolicyV2 when disabling Auto Mode", func() {
		attachedPolicies("AmazonEKSClusterPolicy", "AmazonEKSComputePolicy", "AmazonEKSBlockStoragePolicyV2",
			"AmazonEKSLoadBalancingPolicy", "AmazonEKSNetworkingPolicy")
		Expect(roleManager.DeleteAutoModePolicies(context.Background())).To(Succeed())
		Expect(detached).To(ConsistOf(
			"arn:aws:iam::aws:policy/AmazonEKSComputePolicy",
			"arn:aws:iam::aws:policy/AmazonEKSBlockStoragePolicyV2",
			"arn:aws:iam::aws:policy/AmazonEKSLoadBalancingPolicy",
			"arn:aws:iam::aws:policy/AmazonEKSNetworkingPolicy",
		))
	})
})
