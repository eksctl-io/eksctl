package eks_test

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/fakes"
)

var _ = Describe("Fargate", func() {
	var (
		config     *api.ClusterConfig
		fakeClient *fakes.FakeFargateClient

		clusterName         string
		fargateProfileName  string
		podExecutionRoleArn string
	)

	BeforeEach(func() {
		clusterName = "imaginative-cluster-name"
		fargateProfileName = "a-cool-name-here"
		podExecutionRoleArn = "shrug-emoji"
		fakeClient = new(fakes.FakeFargateClient)

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
			fakeClient.CreateProfileReturns(nil)
			Expect(eks.DoCreateFargateProfiles(context.Background(), config, fakeClient)).To(Succeed())
		})

		When("a profile already exists", func() {
			BeforeEach(func() {
				fakeClient.CreateProfileReturnsOnCall(0, nil)
				Expect(eks.DoCreateFargateProfiles(context.Background(), config, fakeClient)).To(Succeed())
			})

			It("should not error", func() {
				fakeClient.CreateProfileReturnsOnCall(1, &ekstypes.ResourceInUseException{})
				Expect(eks.DoCreateFargateProfiles(context.Background(), config, fakeClient)).To(Succeed())
			})
		})

		When("profile creation fails", func() {
			It("should return the error", func() {
				fakeClient.CreateProfileReturns(errors.New("omigod"))
				Expect(eks.DoCreateFargateProfiles(context.Background(), config, fakeClient)).To(MatchError(ContainSubstring("omigod")))
			})
		})
	})
})
