package fargate_test

import (
	"fmt"

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
				Expect(out[0]).To(Equal(apiFargateProfile(testBlue)))
				Expect(out[1]).To(Equal(apiFargateProfile(testGreen)))
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
				Expect(err.Error()).To(Equal("failed to get EKS cluster \"non-existing-test-cluster\"'s Fargate profile(s): the Internet broke down!"))
				Expect(out).To(BeNil())
			})
		})

		Describe("ReadProfile", func() {
			It("returns the Fargate profile matching the provided name, if any", func() {
				client := fargate.NewClient(clusterName, mockForReadProfile())
				out, err := client.ReadProfile(testGreen)
				Expect(err).To(Not(HaveOccurred()))
				Expect(out).To(Equal(apiFargateProfile(testGreen)))
			})

			It("returns a 'not found' error if no Fargate profile matched the provided name", func() {
				client := fargate.NewClient(clusterName, mockForEmptyReadProfile())
				out, err := client.ReadProfile(testRed)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to get EKS cluster \"non-existing-test-cluster\"'s Fargate profile \"test-red\": ResourceNotFoundException: No Fargate Profile found with name: test-red."))
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

const (
	testBlue  = "test-blue"
	testGreen = "test-green"
	testRed   = "test-red"
)

func mockForReadProfiles() *mocks.EKSAPI {
	mockClient := mocks.EKSAPI{}
	mockClient.Mock.On("ListFargateProfiles", &eks.ListFargateProfilesInput{
		ClusterName: strings.Pointer(clusterName),
	}).Return(&eks.ListFargateProfilesOutput{
		FargateProfileNames: []*string{
			strings.Pointer(testBlue),
			strings.Pointer(testGreen),
		},
	}, nil)
	mockDescribeFargateProfile(&mockClient, testBlue)
	mockDescribeFargateProfile(&mockClient, testGreen)
	return &mockClient
}

func mockForReadProfile() *mocks.EKSAPI {
	mockClient := &mocks.EKSAPI{}
	mockDescribeFargateProfile(mockClient, testGreen)
	return mockClient
}

func mockDescribeFargateProfile(mockClient *mocks.EKSAPI, name string) {
	mockClient.Mock.On("DescribeFargateProfile", &eks.DescribeFargateProfileInput{
		ClusterName:        strings.Pointer(clusterName),
		FargateProfileName: strings.Pointer(name),
	}).Return(&eks.DescribeFargateProfileOutput{
		FargateProfile: eksFargateProfile(name),
	}, nil)
}

func eksFargateProfile(name string) *eks.FargateProfile {
	return &eks.FargateProfile{
		ClusterName:        strings.Pointer(clusterName),
		FargateProfileName: strings.Pointer(name),
		Selectors: []*eks.FargateProfileSelector{
			&eks.FargateProfileSelector{
				Namespace: strings.Pointer(name),
			},
		},
	}
}

func apiFargateProfile(name string) *api.FargateProfile {
	return &api.FargateProfile{
		Name: name,
		Selectors: []api.FargateProfileSelector{
			api.FargateProfileSelector{
				Namespace: name,
				Labels:    map[string]string{},
			},
		},
		Subnets: []string{},
	}
}

func mockForEmptyReadProfiles() *mocks.EKSAPI {
	mockClient := mocks.EKSAPI{}
	mockClient.Mock.On("ListFargateProfiles", &eks.ListFargateProfilesInput{
		ClusterName: strings.Pointer(clusterName),
	}).Return(&eks.ListFargateProfilesOutput{}, nil)
	return &mockClient
}

func mockForEmptyReadProfile() *mocks.EKSAPI {
	mockClient := mocks.EKSAPI{}
	mockClient.Mock.On("DescribeFargateProfile", &eks.DescribeFargateProfileInput{
		ClusterName:        strings.Pointer(clusterName),
		FargateProfileName: strings.Pointer(testRed),
	}).Return(nil, fmt.Errorf("ResourceNotFoundException: No Fargate Profile found with name: %s.", testRed))
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
