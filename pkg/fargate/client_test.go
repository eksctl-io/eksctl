package fargate_test

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks/mocks"
	"github.com/weaveworks/eksctl/pkg/fargate"
	"github.com/weaveworks/eksctl/pkg/utils/retry"
	"github.com/weaveworks/eksctl/pkg/utils/strings"
)

const clusterName = "non-existing-test-cluster"

var _ = Describe("fargate", func() {
	Describe("Client", func() {
		Describe("CreateProfile", func() {
			It("fails fast if the provided profile is nil", func() {
				client := fargate.NewClient(clusterName, &mocks.EKSAPI{})
				waitForCreation := false
				err := client.CreateProfile(nil, waitForCreation)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid Fargate profile: nil"))
			})

			It("creates the provided profile without tag", func() {
				client := fargate.NewClient(clusterName, mockForCreateFargateProfileWithoutTag())
				waitForCreation := false
				err := client.CreateProfile(testFargateProfileWithoutTag(), waitForCreation)
				Expect(err).To(Not(HaveOccurred()))
			})

			It("creates the provided profile", func() {
				client := fargate.NewClient(clusterName, mockForCreateFargateProfile())
				waitForCreation := false
				err := client.CreateProfile(testFargateProfile(), waitForCreation)
				Expect(err).To(Not(HaveOccurred()))
			})

			It("fails by wrapping the root error with some additional context for clarity", func() {
				client := fargate.NewClient(clusterName, mockForFailureOnCreateFargateProfile())
				waitForCreation := false
				err := client.CreateProfile(testFargateProfile(), waitForCreation)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to create Fargate profile \"default\": the Internet broke down!"))
			})

			It("waits for the full creation of the profile when configured to do so", func() {
				retryPolicy := &retry.ConstantBackoff{
					// Retry up to 5 times, not waiting at all, in order to speed tests up.
					Time: 0, TimeUnit: time.Second, MaxRetries: 5,
				}
				numRetriesAfterCreation := 3 // < MaxRetries
				client := fargate.NewClientWithRetryPolicy(clusterName, mockForCreateFargateProfileWithWait(numRetriesAfterCreation), retryPolicy)
				waitForCreation := true
				err := client.CreateProfile(testFargateProfile(), waitForCreation)
				Expect(err).To(Not(HaveOccurred()))
			})

			It("returns an error when waiting for the creation of the profile times out", func() {
				retryPolicy := &retry.ConstantBackoff{
					// Retry up to 5 times, not waiting at all, in order to speed tests up.
					Time: 0, TimeUnit: time.Second, MaxRetries: 5,
				}
				numRetriesAfterCreation := 5 // == MaxRetries, i.e. we will time out.
				client := fargate.NewClientWithRetryPolicy(clusterName, mockForCreateFargateProfileWithWait(numRetriesAfterCreation), retryPolicy)
				waitForCreation := true
				err := client.CreateProfile(testFargateProfile(), waitForCreation)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("timed out while waiting for Fargate profile \"default\"'s creation"))
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
				Expect(err.Error()).To(Equal("failed to get Fargate profile(s) for cluster \"non-existing-test-cluster\": the Internet broke down!"))
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
				Expect(err.Error()).To(Equal("failed to get Fargate profile \"test-red\": ResourceNotFoundException: No Fargate Profile found with name: test-red."))
				Expect(out).To(BeNil())
			})
		})

		Describe("DeleteProfile", func() {
			It("fails fast if the provided profile name is empty", func() {
				client := fargate.NewClient(clusterName, &mocks.EKSAPI{})
				waitForDeletion := false
				err := client.DeleteProfile("", waitForDeletion)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid Fargate profile name: empty"))
			})

			It("deletes the profile corresponding to the provided name", func() {
				profileName := "test-green"
				client := fargate.NewClient(clusterName, mockForDeleteFargateProfile(profileName))
				waitForDeletion := false
				err := client.DeleteProfile(profileName, waitForDeletion)
				Expect(err).To(Not(HaveOccurred()))
			})

			It("fails by wrapping the root error with some additional context for clarity", func() {
				profileName := "test-green"
				client := fargate.NewClient(clusterName, mockForFailureOnDeleteFargateProfile(profileName))
				waitForDeletion := false
				err := client.DeleteProfile(profileName, waitForDeletion)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to delete Fargate profile \"test-green\": the Internet broke down!"))
			})

			It("waits for the full deletion of the profile when configured to do so", func() {
				profileName := "test-green"
				retryPolicy := &retry.ConstantBackoff{
					// Retry up to 5 times, not waiting at all, in order to speed tests up.
					Time: 0, TimeUnit: time.Second, MaxRetries: 5,
				}
				numRetriesBeforeDeletion := 3 // < MaxRetries
				client := fargate.NewClientWithRetryPolicy(clusterName, mockForDeleteFargateProfileWithWait(profileName, numRetriesBeforeDeletion), retryPolicy)
				waitForDeletion := true
				err := client.DeleteProfile(profileName, waitForDeletion)
				Expect(err).To(Not(HaveOccurred()))
			})

			It("returns an error when waiting for the full deletion of the profile times out", func() {
				profileName := "test-green"
				retryPolicy := &retry.ConstantBackoff{
					// Retry up to 5 times, not waiting at all, in order to speed tests up.
					Time: 0, TimeUnit: time.Second, MaxRetries: 5,
				}
				numRetriesBeforeDeletion := 5 // == MaxRetries, i.e. we will time out.
				client := fargate.NewClientWithRetryPolicy(clusterName, mockForDeleteFargateProfileWithWait(profileName, numRetriesBeforeDeletion), retryPolicy)
				waitForDeletion := true
				err := client.DeleteProfile(profileName, waitForDeletion)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("timed out while waiting for Fargate profile \"test-green\"'s deletion"))
			})
		})
	})
})

func mockForCreateFargateProfile() *mocks.EKSAPI {
	mockClient := mocks.EKSAPI{}
	mockCreateFargateProfile(&mockClient)
	return &mockClient
}

func mockForCreateFargateProfileWithoutTag() *mocks.EKSAPI {
	mockClient := mocks.EKSAPI{}
	mockCreateFargateProfileWithoutTag(&mockClient)
	return &mockClient
}

func mockForCreateFargateProfileWithWait(numRetries int) *mocks.EKSAPI {
	mockClient := mocks.EKSAPI{}
	mockCreateFargateProfile(&mockClient)
	// Simulate a couple calls to AWS' API before the profile actually gets created:
	for i := 0; i < numRetries; i++ {
		mockDescribeFargateProfile(&mockClient, "default", "CREATING")
	}
	mockDescribeFargateProfile(&mockClient, "default", "ACTIVE") // At this point, the profile has been created.
	return &mockClient
}

func mockCreateFargateProfile(mockClient *mocks.EKSAPI) {
	mockClient.Mock.On("CreateFargateProfile", testCreateFargateProfileInput()).
		Return(&eks.CreateFargateProfileOutput{}, nil)
}

func mockCreateFargateProfileWithoutTag(mockClient *mocks.EKSAPI) {
	mockClient.Mock.On("CreateFargateProfile", testCreateFargateProfileInputWithoutTag()).
		Return(&eks.CreateFargateProfileOutput{}, nil)
}

func mockForFailureOnCreateFargateProfile() *mocks.EKSAPI {
	mockClient := mocks.EKSAPI{}
	mockClient.Mock.On("CreateFargateProfile", testCreateFargateProfileInput()).
		Return(nil, errors.New("the Internet broke down!"))
	return &mockClient
}

func testFargateProfileWithoutTag() *api.FargateProfile {
	return &api.FargateProfile{
		Name: "default",
		Selectors: []api.FargateProfileSelector{
			{
				Namespace: "kube-system",
				Labels: map[string]string{
					"app": "my-app",
					"env": "test",
				},
			},
		},
	}
}

func testFargateProfile() *api.FargateProfile {
	return &api.FargateProfile{
		Name: "default",
		Selectors: []api.FargateProfileSelector{
			{
				Namespace: "kube-system",
				Labels: map[string]string{
					"app": "my-app",
					"env": "test",
				},
			},
		},
		Tags: map[string]string{
			"env": "test",
		},
	}
}

func testCreateFargateProfileInput() *eks.CreateFargateProfileInput {
	return &eks.CreateFargateProfileInput{
		ClusterName:        strings.Pointer(clusterName),
		FargateProfileName: strings.Pointer("default"),
		Selectors: []*eks.FargateProfileSelector{
			{
				Namespace: strings.Pointer("kube-system"),
				Labels: map[string]*string{
					"app": strings.Pointer("my-app"),
					"env": strings.Pointer("test"),
				},
			},
		},
		Tags: map[string]*string{
			"env": strings.Pointer("test"),
		},
	}
}

func testCreateFargateProfileInputWithoutTag() *eks.CreateFargateProfileInput {
	return &eks.CreateFargateProfileInput{
		ClusterName:        strings.Pointer(clusterName),
		FargateProfileName: strings.Pointer("default"),
		Selectors: []*eks.FargateProfileSelector{
			{
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
	mockListFargateProfiles(&mockClient, testBlue, testGreen)
	mockDescribeFargateProfile(&mockClient, testBlue, "ACTIVE")
	mockDescribeFargateProfile(&mockClient, testGreen, "ACTIVE")
	return &mockClient
}

func mockListFargateProfiles(mockClient *mocks.EKSAPI, names ...string) {
	profileNames := make([]*string, len(names))
	for i, name := range names {
		profileNames[i] = strings.Pointer(name)
	}
	mockClient.Mock.On("ListFargateProfiles", &eks.ListFargateProfilesInput{
		ClusterName: strings.Pointer(clusterName),
	}).Once().Return(&eks.ListFargateProfilesOutput{
		FargateProfileNames: profileNames,
	}, nil)
}

func mockForReadProfile() *mocks.EKSAPI {
	mockClient := &mocks.EKSAPI{}
	mockDescribeFargateProfile(mockClient, testGreen, "ACTIVE")
	return mockClient
}

func mockDescribeFargateProfile(mockClient *mocks.EKSAPI, name, status string) {
	mockClient.Mock.On("DescribeFargateProfile", &eks.DescribeFargateProfileInput{
		ClusterName:        strings.Pointer(clusterName),
		FargateProfileName: strings.Pointer(name),
	}).Once().Return(&eks.DescribeFargateProfileOutput{
		FargateProfile: eksFargateProfile(name, status),
	}, nil)
}

func eksFargateProfile(name, status string) *eks.FargateProfile {
	return &eks.FargateProfile{
		ClusterName:        strings.Pointer(clusterName),
		FargateProfileName: strings.Pointer(name),
		Selectors: []*eks.FargateProfileSelector{
			{
				Namespace: strings.Pointer(name),
			},
		},
		Status: strings.Pointer(status),
	}
}

func apiFargateProfile(name string) *api.FargateProfile {
	return &api.FargateProfile{
		Name: name,
		Selectors: []api.FargateProfileSelector{
			{
				Namespace: name,
				Labels:    map[string]string{},
			},
		},
		Subnets: []string{},
		Tags:    map[string]string{},
	}
}

func mockForEmptyReadProfiles() *mocks.EKSAPI {
	mockClient := mocks.EKSAPI{}
	mockListFargateProfiles(&mockClient)
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
	mockDeleteFargateProfile(&mockClient, name)
	return &mockClient
}

func mockForDeleteFargateProfileWithWait(name string, numRetries int) *mocks.EKSAPI {
	mockClient := mocks.EKSAPI{}
	mockDeleteFargateProfile(&mockClient, name)
	// Simulate a couple calls to AWS' API before the profile actually gets deleted:
	for i := 0; i < numRetries; i++ {
		mockListFargateProfiles(&mockClient, name)
	}
	mockListFargateProfiles(&mockClient) // At this point, the profile has been deleted.
	return &mockClient
}

func mockDeleteFargateProfile(mockClient *mocks.EKSAPI, name string) {
	mockClient.Mock.On("DeleteFargateProfile", &eks.DeleteFargateProfileInput{
		ClusterName:        strings.Pointer(clusterName),
		FargateProfileName: &name,
	}).Return(&eks.DeleteFargateProfileOutput{
		FargateProfile: &eks.FargateProfile{
			FargateProfileName: &name,
			Status:             strings.Pointer("DELETING"),
		},
	}, nil)
}

func mockForFailureOnDeleteFargateProfile(name string) *mocks.EKSAPI {
	mockClient := mocks.EKSAPI{}
	mockClient.Mock.On("DeleteFargateProfile", &eks.DeleteFargateProfileInput{
		ClusterName:        strings.Pointer(clusterName),
		FargateProfileName: &name,
	}).Return(nil, errors.New("the Internet broke down!"))
	return &mockClient
}
