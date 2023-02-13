package fargate_test

import (
	"context"
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/weaveworks/eksctl/pkg/eks/mocksv2"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"
	"github.com/weaveworks/eksctl/pkg/fargate"
	"github.com/weaveworks/eksctl/pkg/utils/retry"
)

const clusterName = "non-existing-test-cluster"

var retryPolicy = retry.NewTimingOutExponentialBackoff(5 * time.Minute)

var _ = Describe("fargate", func() {
	Describe("Client", func() {
		var neverCalledStackManager *manager.StackCollection
		Describe("CreateProfile", func() {
			It("fails fast if the provided profile is nil", func() {
				client := fargate.NewWithRetryPolicy(clusterName, &mocksv2.EKS{}, &retryPolicy, neverCalledStackManager)
				waitForCreation := false
				err := client.CreateProfile(context.Background(), nil, waitForCreation)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid Fargate profile: nil"))
			})

			It("creates the provided profile without tag", func() {
				client := fargate.NewWithRetryPolicy(clusterName, mockForCreateFargateProfileWithoutTag(), &retryPolicy, neverCalledStackManager)
				waitForCreation := false
				err := client.CreateProfile(context.Background(), createProfileWithoutTag(), waitForCreation)
				Expect(err).To(Not(HaveOccurred()))
			})

			It("creates the provided profile", func() {
				client := fargate.NewWithRetryPolicy(clusterName, mockForCreateFargateProfile(), &retryPolicy, neverCalledStackManager)
				waitForCreation := false
				err := client.CreateProfile(context.Background(), testFargateProfile(), waitForCreation)
				Expect(err).To(Not(HaveOccurred()))
			})

			It("fails by wrapping the root error with some additional context for clarity", func() {
				client := fargate.NewWithRetryPolicy(clusterName, mockForFailureOnCreateFargateProfile(), &retryPolicy, neverCalledStackManager)
				waitForCreation := false
				err := client.CreateProfile(context.Background(), testFargateProfile(), waitForCreation)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to create Fargate profile \"default\": the Internet broke down"))
			})

			It("waits for the full creation of the profile when configured to do so", func() {
				retryPolicy := &retry.ConstantBackoff{
					// Retry up to 5 times, not waiting at all, in order to speed tests up.
					Time: 0, TimeUnit: time.Second, MaxRetries: 5,
				}
				numRetriesAfterCreation := 3 // < MaxRetries
				client := fargate.NewWithRetryPolicy(clusterName, mockForCreateFargateProfileWithWait(numRetriesAfterCreation), retryPolicy, nil)
				waitForCreation := true
				err := client.CreateProfile(context.Background(), testFargateProfile(), waitForCreation)
				Expect(err).To(Not(HaveOccurred()))
			})

			It("returns an error when waiting for the creation of the profile times out", func() {
				retryPolicy := &retry.ConstantBackoff{
					// Retry up to 5 times, not waiting at all, in order to speed tests up.
					Time: 0, TimeUnit: time.Second, MaxRetries: 5,
				}
				numRetriesAfterCreation := 5 // == MaxRetries, i.e. we will time out.
				client := fargate.NewWithRetryPolicy(clusterName, mockForCreateFargateProfileWithWait(numRetriesAfterCreation), retryPolicy, nil)
				waitForCreation := true
				err := client.CreateProfile(context.Background(), testFargateProfile(), waitForCreation)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("timed out while waiting for Fargate profile \"default\"'s creation"))
			})

			It("suggests using a different AZ when profile creation fails with an unsupported AZ error", func() {
				retryPolicy := &retry.ConstantBackoff{
					Time: 0, TimeUnit: time.Second, MaxRetries: 1,
				}
				client := fargate.NewWithRetryPolicy(clusterName, mockCreateFargateProfileAZError(), retryPolicy, nil)
				err := client.CreateProfile(context.Background(), testFargateProfile(), false)
				Expect(err).To(MatchError(ContainSubstring("Fargate Profile creation for the Availability Zone ca-central-1d for Subnet subnet-1234 is not supported; please rerun the command by supplying subnets in the Fargate Profile that do not exist in the unsupported AZ, or recreate the cluster after specifying supported AZs in `availabilityZones`")))
			})
		})

		Describe("ReadProfiles", func() {
			It("returns all Fargate profiles", func() {
				client := fargate.NewWithRetryPolicy(clusterName, mockForReadProfiles(), &retryPolicy, neverCalledStackManager)
				out, err := client.ReadProfiles(context.Background())
				Expect(err).To(Not(HaveOccurred()))
				Expect(out).To(Not(BeNil()))
				Expect(out).To(HaveLen(2))
				Expect(out[0]).To(Equal(apiFargateProfile(testBlue)))
				Expect(out[1]).To(Equal(apiFargateProfile(testGreen)))
			})

			It("returns an empty array if no Fargate profile exists", func() {
				client := fargate.NewWithRetryPolicy(clusterName, mockForEmptyReadProfiles(), &retryPolicy, neverCalledStackManager)
				out, err := client.ReadProfiles(context.Background())
				Expect(err).To(Not(HaveOccurred()))
				Expect(out).To(Not(BeNil()))
				Expect(out).To(HaveLen(0))
			})

			It("fails by wrapping the root error with some additional context for clarity", func() {
				client := fargate.NewWithRetryPolicy(clusterName, mockForFailureOnReadProfiles(), &retryPolicy, neverCalledStackManager)
				out, err := client.ReadProfiles(context.Background())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to get Fargate profile(s) for cluster \"non-existing-test-cluster\": the Internet broke down"))
				Expect(out).To(BeNil())
			})
		})

		Describe("ListProfile", func() {
			It("list empty profile without error", func() {
				client := fargate.NewWithRetryPolicy(clusterName, mockForListEmptyProfile(), &retryPolicy, neverCalledStackManager)
				out, err := client.ListProfiles(context.Background())
				Expect(err).To(Not(HaveOccurred()))
				Expect(out).To(HaveLen(0))
			})

			It("list multiple profiles without error", func() {
				client := fargate.NewWithRetryPolicy(clusterName, mockForListMultipleProfiles(), &retryPolicy, neverCalledStackManager)
				out, err := client.ListProfiles(context.Background())
				Expect(err).To(Not(HaveOccurred()))
				Expect(out).To(HaveLen(4))
			})

			It("list profiles with error", func() {
				client := fargate.NewWithRetryPolicy(clusterName, mockForListProfilesWithError(), &retryPolicy, neverCalledStackManager)
				out, err := client.ListProfiles(context.Background())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to get Fargate profile(s) for cluster \"non-existing-test-cluster\": failed to get Fargate Profile list"))
				Expect(out).To(BeNil())
			})

			It("list all profiles with multiple requests without error", func() {
				client := fargate.NewWithRetryPolicy(clusterName, mockListFargateProfilesMultipleRequest(), &retryPolicy, neverCalledStackManager)
				out, err := client.ListProfiles(context.Background())
				Expect(err).To(Not(HaveOccurred()))
				Expect(out).To(HaveLen(4))
			})
		})

		Describe("ReadProfile", func() {
			It("returns the Fargate profile matching the provided name, if any", func() {
				client := fargate.NewWithRetryPolicy(clusterName, mockForReadProfile(), &retryPolicy, neverCalledStackManager)
				out, err := client.ReadProfile(context.Background(), testGreen)
				Expect(err).To(Not(HaveOccurred()))
				Expect(out).To(Equal(apiFargateProfile(testGreen)))
			})

			It("returns a 'not found' error if no Fargate profile matched the provided name", func() {
				client := fargate.NewWithRetryPolicy(clusterName, mockForEmptyReadProfile(), &retryPolicy, neverCalledStackManager)
				out, err := client.ReadProfile(context.Background(), testRed)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to get Fargate profile \"test-red\": ResourceNotFoundException: No Fargate Profile found with name: test-red"))
				Expect(out).To(BeNil())
			})
		})

		Describe("DeleteProfile", func() {
			It("fails fast if the provided profile name is empty", func() {
				client := fargate.NewWithRetryPolicy(clusterName, &mocksv2.EKS{}, &retryPolicy, neverCalledStackManager)
				waitForDeletion := false
				err := client.DeleteProfile(context.Background(), "", waitForDeletion)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid Fargate profile name: empty"))
			})

			It("deletes the profile corresponding to the provided name", func() {
				profileName := "test-green"
				client := fargate.NewWithRetryPolicy(clusterName, mockForDeleteFargateProfile(profileName), &retryPolicy, neverCalledStackManager)
				waitForDeletion := false
				err := client.DeleteProfile(context.Background(), profileName, waitForDeletion)
				Expect(err).To(Not(HaveOccurred()))
			})

			It("deletes the stack when no fargate profiles remain", func() {
				fakeStackManager := new(fakes.FakeStackManager)
				fakeStackManager.GetFargateStackReturns(&types.Stack{
					StackName: aws.String("my-fargate-profile"),
				}, nil)
				profileName := "test-green"
				client := fargate.NewWithRetryPolicy(clusterName, mockForDeleteFargateProfileWithoutAnyRemaining(profileName), &retryPolicy, fakeStackManager)
				waitForDeletion := false
				err := client.DeleteProfile(context.Background(), profileName, waitForDeletion)
				Expect(err).To(Not(HaveOccurred()))
				Expect(fakeStackManager.GetFargateStackCallCount()).To(Equal(1))
				Expect(fakeStackManager.DeleteStackBySpecCallCount()).To(Equal(1))
				_, stack := fakeStackManager.DeleteStackBySpecArgsForCall(0)
				Expect(*stack.StackName).To(Equal("my-fargate-profile"))
			})

			It("fails by wrapping the root error with some additional context for clarity", func() {
				profileName := "test-green"
				client := fargate.NewWithRetryPolicy(clusterName, mockForFailureOnDeleteFargateProfile(profileName), &retryPolicy, neverCalledStackManager)
				waitForDeletion := false
				err := client.DeleteProfile(context.Background(), profileName, waitForDeletion)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to delete Fargate profile \"test-green\": the Internet broke down"))
			})

			It("waits for the full deletion of the profile when configured to do so", func() {
				profileName := "test-green"
				retryPolicy := &retry.ConstantBackoff{
					// Retry up to 5 times, not waiting at all, in order to speed tests up.
					Time: 0, TimeUnit: time.Second, MaxRetries: 5,
				}
				numRetriesBeforeDeletion := 3 // < MaxRetries
				client := fargate.NewWithRetryPolicy(clusterName, mockForDeleteFargateProfileWithWait(profileName, numRetriesBeforeDeletion), retryPolicy, nil)
				waitForDeletion := true
				err := client.DeleteProfile(context.Background(), profileName, waitForDeletion)
				Expect(err).To(Not(HaveOccurred()))
			})

			It("returns an error when waiting for the full deletion of the profile times out", func() {
				profileName := "test-green"
				retryPolicy := &retry.ConstantBackoff{
					// Retry up to 5 times, not waiting at all, in order to speed tests up.
					Time: 0, TimeUnit: time.Second, MaxRetries: 5,
				}
				numRetriesBeforeDeletion := 5 // == MaxRetries, i.e. we will time out.
				client := fargate.NewWithRetryPolicy(clusterName, mockForDeleteFargateProfileWithWait(profileName, numRetriesBeforeDeletion), retryPolicy, nil)
				waitForDeletion := true
				err := client.DeleteProfile(context.Background(), profileName, waitForDeletion)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("timed out while waiting for Fargate profile \"test-green\"'s deletion"))
			})
		})
	})
})

func mockForCreateFargateProfile() *mocksv2.EKS {
	mockClient := mocksv2.EKS{}
	mockCreateFargateProfile(&mockClient)
	return &mockClient
}

func mockForCreateFargateProfileWithoutTag() *mocksv2.EKS {
	mockClient := mocksv2.EKS{}
	mockCreateFargateProfileWithoutTag(&mockClient)
	return &mockClient
}

func mockForCreateFargateProfileWithWait(numRetries int) *mocksv2.EKS {
	mockClient := mocksv2.EKS{}
	mockCreateFargateProfile(&mockClient)
	// Simulate a couple calls to AWS' API before the profile actually gets created:
	for i := 0; i < numRetries; i++ {
		mockDescribeFargateProfile(&mockClient, "default", "CREATING")
	}
	mockDescribeFargateProfile(&mockClient, "default", "ACTIVE") // At this point, the profile has been created.
	return &mockClient
}

func mockCreateFargateProfile(mockClient *mocksv2.EKS) {
	mockClient.Mock.On("CreateFargateProfile", mock.Anything, testCreateFargateProfileInput()).
		Return(&eks.CreateFargateProfileOutput{}, nil)
}

func mockCreateFargateProfileAZError() *mocksv2.EKS {
	var mockClient mocksv2.EKS
	mockClient.Mock.On("CreateFargateProfile", mock.Anything, testCreateFargateProfileInput()).
		Return(nil, &ekstypes.InvalidParameterException{
			Message:            aws.String("Fargate Profile creation for the Availability Zone ca-central-1d for Subnet subnet-1234 is not supported"),
			FargateProfileName: aws.String("test"),
		})
	return &mockClient
}

func mockCreateFargateProfileWithoutTag(mockClient *mocksv2.EKS) {
	mockClient.Mock.On("CreateFargateProfile", mock.Anything, createEksProfileWithoutTag()).
		Return(&eks.CreateFargateProfileOutput{}, nil)
}

func mockForFailureOnCreateFargateProfile() *mocksv2.EKS {
	mockClient := mocksv2.EKS{}
	mockClient.Mock.On("CreateFargateProfile", mock.Anything, testCreateFargateProfileInput()).
		Return(nil, errors.New("the Internet broke down"))
	return &mockClient
}

func createProfileWithoutTag() *api.FargateProfile {
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
		Status: "ACTIVE",
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
		ClusterName:        aws.String(clusterName),
		FargateProfileName: aws.String("default"),
		Selectors: []ekstypes.FargateProfileSelector{
			{
				Namespace: aws.String("kube-system"),
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

func createEksProfileWithoutTag() *eks.CreateFargateProfileInput {
	return &eks.CreateFargateProfileInput{
		ClusterName:        aws.String(clusterName),
		FargateProfileName: aws.String("default"),
		Selectors: []ekstypes.FargateProfileSelector{
			{
				Namespace: aws.String("kube-system"),
				Labels: map[string]string{
					"app": "my-app",
					"env": "test",
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

func mockForReadProfiles() *mocksv2.EKS {
	mockClient := mocksv2.EKS{}
	mockListFargateProfiles(&mockClient, testBlue, testGreen)
	mockDescribeFargateProfile(&mockClient, testBlue, "ACTIVE")
	mockDescribeFargateProfile(&mockClient, testGreen, "ACTIVE")
	return &mockClient
}

func mockListFargateProfilesMultipleRequest() *mocksv2.EKS {
	mockClient := mocksv2.EKS{}
	mockClient.Mock.On("ListFargateProfiles", mock.Anything, &eks.ListFargateProfilesInput{
		ClusterName: aws.String(clusterName),
	}).Return(&eks.ListFargateProfilesOutput{
		FargateProfileNames: []string{"default"},
		NextToken:           aws.String(testBlue),
	}, nil)

	mockClient.Mock.On("ListFargateProfiles", mock.Anything, &eks.ListFargateProfilesInput{
		ClusterName: aws.String(clusterName),
		NextToken:   aws.String(testBlue),
	}).Return(&eks.ListFargateProfilesOutput{
		FargateProfileNames: []string{testBlue},
		NextToken:           aws.String(testGreen),
	}, nil)

	mockClient.Mock.On("ListFargateProfiles", mock.Anything, &eks.ListFargateProfilesInput{
		ClusterName: aws.String(clusterName),
		NextToken:   aws.String(testGreen),
	}).Return(&eks.ListFargateProfilesOutput{
		FargateProfileNames: []string{testGreen},
		NextToken:           aws.String(testRed),
	}, nil)

	mockClient.Mock.On("ListFargateProfiles", mock.Anything, &eks.ListFargateProfilesInput{
		ClusterName: aws.String(clusterName),
		NextToken:   aws.String(testRed),
	}).Return(&eks.ListFargateProfilesOutput{
		FargateProfileNames: []string{testRed},
	}, nil)
	return &mockClient
}

func mockListFargateProfiles(mockClient *mocksv2.EKS, names ...string) {
	mockClient.Mock.On("ListFargateProfiles", mock.Anything, &eks.ListFargateProfilesInput{
		ClusterName: aws.String(clusterName),
	}).Return(&eks.ListFargateProfilesOutput{
		FargateProfileNames: names,
	}, nil).Once()
}

func mockForReadProfile() *mocksv2.EKS {
	mockClient := &mocksv2.EKS{}
	mockDescribeFargateProfile(mockClient, testGreen, "ACTIVE")
	return mockClient
}

func mockDescribeFargateProfile(mockClient *mocksv2.EKS, name, status string) {
	mockClient.Mock.On("DescribeFargateProfile", mock.Anything, &eks.DescribeFargateProfileInput{
		ClusterName:        aws.String(clusterName),
		FargateProfileName: aws.String(name),
	}).Return(&eks.DescribeFargateProfileOutput{
		FargateProfile: eksFargateProfile(name, status),
	}, nil).Once()
}

func eksFargateProfile(name, status string) *ekstypes.FargateProfile {
	return &ekstypes.FargateProfile{
		ClusterName:        aws.String(clusterName),
		FargateProfileName: aws.String(name),
		Selectors: []ekstypes.FargateProfileSelector{
			{
				Namespace: aws.String(name),
			},
		},
		Status: ekstypes.FargateProfileStatus(status),
	}
}

func apiFargateProfile(name string) *api.FargateProfile {
	return &api.FargateProfile{
		Name: name,
		Selectors: []api.FargateProfileSelector{
			{
				Namespace: name,
			},
		},
		Status: "ACTIVE",
	}
}

func mockForListEmptyProfile() *mocksv2.EKS {
	mockClient := mocksv2.EKS{}
	mockListFargateProfiles(&mockClient)
	return &mockClient
}

func mockForListMultipleProfiles() *mocksv2.EKS {
	mockClient := mocksv2.EKS{}
	mockListFargateProfiles(&mockClient, "default", testBlue, testGreen, testRed)
	return &mockClient
}

func mockForListProfilesWithError() *mocksv2.EKS {
	mockClient := mocksv2.EKS{}
	mockClient.Mock.On("ListFargateProfiles", mock.Anything, &eks.ListFargateProfilesInput{
		ClusterName: aws.String(clusterName),
	}).Return(nil, errors.New("failed to get Fargate Profile list"))
	return &mockClient
}

func mockForEmptyReadProfiles() *mocksv2.EKS {
	mockClient := mocksv2.EKS{}
	mockListFargateProfiles(&mockClient)
	return &mockClient
}

func mockForEmptyReadProfile() *mocksv2.EKS {
	mockClient := mocksv2.EKS{}
	mockClient.Mock.On("DescribeFargateProfile", mock.Anything, &eks.DescribeFargateProfileInput{
		ClusterName:        aws.String(clusterName),
		FargateProfileName: aws.String(testRed),
	}).Return(nil, fmt.Errorf("ResourceNotFoundException: No Fargate Profile found with name: %s", testRed))
	return &mockClient
}

func mockForFailureOnReadProfiles() *mocksv2.EKS {
	mockClient := mocksv2.EKS{}
	mockClient.Mock.On("ListFargateProfiles", mock.Anything, &eks.ListFargateProfilesInput{
		ClusterName: aws.String(clusterName),
	}).Return(nil, errors.New("the Internet broke down"))
	return &mockClient
}

func mockForDeleteFargateProfile(name string) *mocksv2.EKS {
	mockClient := mocksv2.EKS{}
	mockDeleteFargateProfile(&mockClient, name)
	mockListFargateProfiles(&mockClient, "default")
	return &mockClient
}

func mockForDeleteFargateProfileWithoutAnyRemaining(name string) *mocksv2.EKS {
	mockClient := mocksv2.EKS{}
	mockDeleteFargateProfile(&mockClient, name)
	mockListFargateProfiles(&mockClient)
	return &mockClient
}

func mockForDeleteFargateProfileWithWait(name string, numRetries int) *mocksv2.EKS {
	mockClient := mocksv2.EKS{}
	mockDeleteFargateProfile(&mockClient, name)
	// Simulate a couple calls to AWS' API before the profile actually gets deleted:
	for i := 0; i < numRetries; i++ {
		mockListFargateProfiles(&mockClient, name, "default")
	}
	mockListFargateProfiles(&mockClient, "default") // At this point, the profile has been deleted.

	return &mockClient
}

func mockDeleteFargateProfile(mockClient *mocksv2.EKS, name string) {
	mockClient.Mock.On("DeleteFargateProfile", mock.Anything, &eks.DeleteFargateProfileInput{
		ClusterName:        aws.String(clusterName),
		FargateProfileName: &name,
	}).Return(&eks.DeleteFargateProfileOutput{
		FargateProfile: &ekstypes.FargateProfile{
			FargateProfileName: &name,
			Status:             ekstypes.FargateProfileStatusDeleting,
		},
	}, nil)
}

func mockForFailureOnDeleteFargateProfile(name string) *mocksv2.EKS {
	mockClient := mocksv2.EKS{}
	mockClient.Mock.On("DeleteFargateProfile", mock.Anything, &eks.DeleteFargateProfileInput{
		ClusterName:        aws.String(clusterName),
		FargateProfileName: &name,
	}).Return(nil, errors.New("the Internet broke down"))
	return &mockClient
}
