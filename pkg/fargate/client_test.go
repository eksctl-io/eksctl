package fargate_test

import (
	"github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks/mocks"
	"github.com/weaveworks/eksctl/pkg/fargate"
	"github.com/weaveworks/eksctl/pkg/utils/strings"
)

const clusterName = "non-existing-test-cluster"

var _ = Describe("fargate", func() {
	Describe("Client", func() {
		Describe("CreateProfile", func() {
			It("fails fast if the provided profile is nil", func() {
				client := fargate.NewClient(clusterName, &mocks.EKSAPI{})
				err := client.CreateProfile(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid Fargate profile: nil"))
			})

			It("creates the provided profile", func() {
				client := fargate.NewClient(clusterName, mockForCreateFargateProfile())
				err := client.CreateProfile(testFargateProfile())
				Expect(err).To(Not(HaveOccurred()))
			})

			It("fails by wrapping the root error with some additional context for clarity", func() {
				client := fargate.NewClient(clusterName, mockForFailureOnCreateFargateProfile())
				err := client.CreateProfile(testFargateProfile())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to create Fargate profile \"default\" in cluster \"non-existing-test-cluster\": the Internet broke down!"))
			})
		})

		Describe("ReadProfiles", func() {
			It("returns all Fargate profiles", func() {
				client := fargate.NewClient(clusterName, mockForReadProfiles())
				out, err := client.ReadProfiles()
				Expect(err).To(Not(HaveOccurred()))
				Expect(out).To(Not(BeNil()))
				Expect(out).To(HaveLen(2))
			})

			It("returns an empty array if no Fargate profile exists", func() {
				client := fargate.NewClient(clusterName, mockForEmptyReadProfiles())
				out, err := client.ReadProfiles()
				Expect(err).To(Not(HaveOccurred()))
				Expect(out).To(Not(BeNil()))
				Expect(out).To(HaveLen(0))
			})

			It("fails by wrapping the root error with some additional context for clarity", func() {
				client := fargate.NewClient(clusterName, mockForFailureOnReadProfiles())
				out, err := client.ReadProfiles()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to get EKS cluster \"non-existing-test-cluster\"'s Fargate profile(s) (current token: <nil>): the Internet broke down!"))
				Expect(out).To(BeNil())
			})
		})

		Describe("DeleteProfile", func() {
			It("fails fast if the provided profile name is empty", func() {
				client := fargate.NewClient(clusterName, &mocks.EKSAPI{})
				err := client.DeleteProfile("")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid Fargate profile name: empty"))
			})

			It("deletes the profile corresponding to the provided name", func() {
				profileName := "test-green"
				client := fargate.NewClient(clusterName, mockForDeleteFargateProfile(profileName))
				err := client.DeleteProfile(profileName)
				Expect(err).To(Not(HaveOccurred()))
			})

			It("fails by wrapping the root error with some additional context for clarity", func() {
				profileName := "test-green"
				client := fargate.NewClient(clusterName, mockForFailureOnDeleteFargateProfile(profileName))
				err := client.DeleteProfile(profileName)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to delete Fargate profile \"test-green\" from cluster \"non-existing-test-cluster\": the Internet broke down!"))
			})
		})
	})
})

func mockForCreateFargateProfile() *mocks.EKSAPI {
	mockClient := mocks.EKSAPI{}
	mockClient.Mock.On("CreateFargateProfile", testCreateFargateProfileInput()).
		Return(&eks.CreateFargateProfileOutput{}, nil)
	return &mockClient
}

func mockForFailureOnCreateFargateProfile() *mocks.EKSAPI {
	mockClient := mocks.EKSAPI{}
	mockClient.Mock.On("CreateFargateProfile", testCreateFargateProfileInput()).
		Return(nil, errors.New("the Internet broke down!"))
	return &mockClient
}

func testFargateProfile() *api.FargateProfile {
	return &api.FargateProfile{
		Name: "default",
		Selectors: []api.FargateProfileSelector{
			api.FargateProfileSelector{
				Namespace: "kube-system",
				Labels: map[string]string{
					"app": "my-app",
					"env": "test",
				},
			},
		},
	}
}

func testCreateFargateProfileInput() *eks.CreateFargateProfileInput {
	return &eks.CreateFargateProfileInput{
		ClusterName:        strings.Pointer(clusterName),
		FargateProfileName: strings.Pointer("default"),
		Selectors: []*eks.FargateProfileSelector{
			&eks.FargateProfileSelector{
				Namespace: strings.Pointer("kube-system"),
				Labels: map[string]*string{
					"app": strings.Pointer("my-app"),
					"env": strings.Pointer("test"),
				},
			},
		},
	}
}

func mockForReadProfiles() *mocks.EKSAPI {
	mockClient := mocks.EKSAPI{}

	// First "page" of Fargate profiles (of size 1):
	mockClient.Mock.On("ListFargateProfiles", &eks.ListFargateProfilesInput{
		ClusterName: strings.Pointer(clusterName),
	}).Return(&eks.ListFargateProfilesOutput{
		FargateProfiles: []*eks.FargateProfile{
			&eks.FargateProfile{
				ClusterName:        strings.Pointer(clusterName),
				FargateProfileName: strings.Pointer("test-blue"),
				Selectors: []*eks.FargateProfileSelector{
					&eks.FargateProfileSelector{
						Namespace: strings.Pointer("test-blue"),
					},
				},
			},
		},
		NextToken: strings.Pointer("1"), // all items after item #1
	}, nil)

	// Second "page" of Fargate profiles (of size 1):
	mockClient.Mock.On("ListFargateProfiles", &eks.ListFargateProfilesInput{
		ClusterName: strings.Pointer(clusterName),
		NextToken:   strings.Pointer("1"), // all items after item #1
	}).Return(&eks.ListFargateProfilesOutput{
		FargateProfiles: []*eks.FargateProfile{
			&eks.FargateProfile{
				ClusterName:        strings.Pointer(clusterName),
				FargateProfileName: strings.Pointer("test-green"),
				Selectors: []*eks.FargateProfileSelector{
					&eks.FargateProfileSelector{
						Namespace: strings.Pointer("test-green"),
					},
				},
			},
		},
		NextToken: strings.Pointer("2"), // all items after item #2
	}, nil)

	// No more Fargate profile to read:
	mockClient.Mock.On("ListFargateProfiles", &eks.ListFargateProfilesInput{
		ClusterName: strings.Pointer(clusterName),
		NextToken:   strings.Pointer("2"), // all items after item #2
	}).Return(&eks.ListFargateProfilesOutput{
		FargateProfiles: []*eks.FargateProfile{},
	}, nil)

	return &mockClient
}

func mockForEmptyReadProfiles() *mocks.EKSAPI {
	mockClient := mocks.EKSAPI{}
	mockClient.Mock.On("ListFargateProfiles", &eks.ListFargateProfilesInput{
		ClusterName: strings.Pointer(clusterName),
	}).Return(&eks.ListFargateProfilesOutput{}, nil)
	return &mockClient
}

func mockForFailureOnReadProfiles() *mocks.EKSAPI {
	mockClient := mocks.EKSAPI{}
	mockClient.Mock.On("ListFargateProfiles", &eks.ListFargateProfilesInput{
		ClusterName: strings.Pointer(clusterName),
	}).Return(nil, errors.New("the Internet broke down!"))
	return &mockClient
}

func mockForDeleteFargateProfile(name string) *mocks.EKSAPI {
	mockClient := mocks.EKSAPI{}
	mockClient.Mock.On("DeleteFargateProfile", &eks.DeleteFargateProfileInput{
		ClusterName:        strings.Pointer(clusterName),
		FargateProfileName: &name,
	}).Return(&eks.DeleteFargateProfileOutput{}, nil)
	return &mockClient
}

func mockForFailureOnDeleteFargateProfile(name string) *mocks.EKSAPI {
	mockClient := mocks.EKSAPI{}
	mockClient.Mock.On("DeleteFargateProfile", &eks.DeleteFargateProfileInput{
		ClusterName:        strings.Pointer(clusterName),
		FargateProfileName: &name,
	}).Return(nil, errors.New("the Internet broke down!"))
	return &mockClient
}
