package eks_test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	pkg_eks "github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/fakes"
)

var _ = Describe("Fargate", func() {
	var (
		config      *api.ClusterConfig
		fakeManager *fakes.FakeFargateManager

		clusterName         string
		fargateProfileName  string
		podExecutionRoleArn string
	)

	BeforeEach(func() {
		clusterName = "imaginative-cluster-name"
		fargateProfileName = "a-cool-name-here"
		podExecutionRoleArn = "shrug-emoji"
		fakeManager = new(fakes.FakeFargateManager)

		config = &api.ClusterConfig{
			Metadata: &api.ClusterMeta{
				Name: clusterName,
			},
			IAM: &api.ClusterIAM{
				FargatePodExecutionRoleARN: aws.String("pod-eran"),
			},
			FargateProfiles: []*api.FargateProfile{
				{
					Name:                fargateProfileName,
					PodExecutionRoleARN: podExecutionRoleArn,
				},
			},
		}
	})

	Describe("DoCreateFargateProfiles", func() {
		It("should create profiles", func() {
			fakeManager.CreateProfileReturns(nil)
			Expect(pkg_eks.DoCreateFargateProfiles(config, fakeManager)).To(Succeed())
		})

		When("a profile already exists", func() {
			BeforeEach(func() {
				fakeManager.CreateProfileReturnsOnCall(0, nil)
				Expect(pkg_eks.DoCreateFargateProfiles(config, fakeManager)).To(Succeed())
			})

			It("should not error", func() {
				fakeManager.CreateProfileReturnsOnCall(1, &eks.ResourceInUseException{})
				Expect(pkg_eks.DoCreateFargateProfiles(config, fakeManager)).To(Succeed())
			})
		})

		When("profile creation fails", func() {
			It("should return the error", func() {
				fakeManager.CreateProfileReturns(errors.New("omigod"))
				Expect(pkg_eks.DoCreateFargateProfiles(config, fakeManager)).To(MatchError(ContainSubstring("omigod")))
			})
		})
	})
})
