package client

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/aws/smithy-go"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/eks/mocksv2"

	"github.com/stretchr/testify/mock"
)

var _ = Describe("ssh public key", func() {
	var (
		clusterName        = "sshtestcluster"
		ngName             = "ng1"
		rsaKey             = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDcSoNjWaJaw+MYBz43lgm12ZGdP+zRs9o0sXAGbiQua6e3JSkAiH4p9YZHmWxCTjckbiEdXN5qcs5OC5KUxYBvnEgor7jEydcKe1ZJXqsm/8CrtnJMTNcO9QVFnXfjvpkNjgNYj+8w9PcFRr0JDgDhRb52JvPWoqywv/Om9s1hpUov0gxDIl6CLLHSk0lmXZEhtVMMJmo0Tu/NlHqdky2DxFgHyNjBcMNpiBd8bs3dA5xf36dY+qgcXBV23i1SCgbqn9xcw1Q0IrHuQ4/QB+PJ5haxUx0bnOTahxSZ+tlEz9EiLwlM8VtKo3ND/giBvGaXuIK2iGDL0kSCRjueM5/3 user@example\n"
		ed25519Key         = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBoB6Gtu8zPAPO1yF4OwysWUD8ZSEQYzMpOT0YvF9qJV user@example\n"
		keyName            = "eksctl-sshtestcluster-nodegroup-ng1-f5:d9:01:88:1e:fb:40:fb:e1:ca:69:fe:2e:31:03:6c"
		ed25519KeyName     = "eksctl-sshtestcluster-nodegroup-ng1-HvE7+gmH78VS53+iPuRDh/gKjVo26OzYU/qOnJWAgyk"
		rsaFingerprint     = "f5:d9:01:88:1e:fb:40:fb:e1:ca:69:fe:2e:31:03:6c"
		ed25519Fingerprint = "HvE7+gmH78VS53+iPuRDh/gKjVo26OzYU/qOnJWAgyk"
		mockEC2            *mocksv2.EC2
	)

	BeforeEach(func() {
		mockEC2 = &mocksv2.EC2{}
	})

	Describe("loading from a file", func() {

		It("should import the key", func() {
			mockDescribeKeyPairs(mockEC2, make(map[string]string))
			mockImportKeyPair(mockEC2, keyName, rsaFingerprint, rsaKey)

			keyName, err := LoadKeyFromFile(context.Background(), "assets/id_rsa_tests1.pub", clusterName, ngName, mockEC2)

			Expect(err).NotTo(HaveOccurred())
			Expect(keyName).To(Equal("eksctl-sshtestcluster-nodegroup-ng1-f5:d9:01:88:1e:fb:40:fb:e1:ca:69:fe:2e:31:03:6c"))
			mockEC2.AssertCalled(GinkgoT(),
				"ImportKeyPair",
				mock.Anything,
				&ec2.ImportKeyPairInput{
					KeyName:           &keyName,
					PublicKeyMaterial: []byte(rsaKey),
				})
		})

		It("should not import key that already exists in EC2", func() {
			mockDescribeKeyPairs(mockEC2, map[string]string{keyName: rsaFingerprint})
			mockImportKeyPairError(mockEC2, errors.New("the key shouldn't be imported in this test"))

			keyName, err := LoadKeyFromFile(context.Background(), "assets/id_rsa_tests1.pub", clusterName, ngName, mockEC2)

			Expect(err).NotTo(HaveOccurred())
			Expect(keyName).To(Equal("eksctl-sshtestcluster-nodegroup-ng1-f5:d9:01:88:1e:fb:40:fb:e1:ca:69:fe:2e:31:03:6c"))
			mockEC2.AssertNotCalled(GinkgoT(), "ImportKeyPair", mock.Anything, mock.Anything)
		})

		It("should return error if a key with same name exists in EC2 with different fingerprint", func() {
			differentFingerprint := "ab:cd"
			mockDescribeKeyPairs(mockEC2, map[string]string{keyName: differentFingerprint})
			mockImportKeyPairError(mockEC2, errors.New("the key shouldn't be imported in this test"))

			_, err := LoadKeyFromFile(context.Background(), "assets/id_rsa_tests1.pub", clusterName, ngName, mockEC2)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("but fingerprints don't match"))
		})

		It("should return error if the file does not exist", func() {
			_, err := LoadKeyFromFile(context.Background(), "assets/file_not_existing.pub", clusterName, ngName, mockEC2)

			Expect(err).To(HaveOccurred())
		})

		When("they key is of type ed25519", func() {
			It("loads the key", func() {
				mockDescribeKeyPairs(mockEC2, make(map[string]string))
				mockImportKeyPair(mockEC2, ed25519KeyName, ed25519Fingerprint, ed25519Key)

				keyName, err := LoadKeyFromFile(context.Background(), "assets/id_ed25519_tests1.pub", clusterName, ngName, mockEC2)

				Expect(err).NotTo(HaveOccurred())
				Expect(keyName).To(Equal("eksctl-sshtestcluster-nodegroup-ng1-HvE7+gmH78VS53+iPuRDh/gKjVo26OzYU/qOnJWAgyk"))
				mockEC2.AssertCalled(GinkgoT(),
					"ImportKeyPair",
					mock.Anything,
					&ec2.ImportKeyPairInput{
						KeyName:           &keyName,
						PublicKeyMaterial: []byte(ed25519Key),
					})
			})
		})

		When("they key is invalid", func() {
			It("errors", func() {
				_, err := LoadKeyFromFile(context.Background(), "assets/invalid.pub", clusterName, ngName, mockEC2)
				Expect(err).To(MatchError(ContainSubstring("parsing key \"assets/invalid.pub\"")))
			})
		})
	})

	Describe("loading by the key content", func() {
		It("should import it", func() {
			mockDescribeKeyPairs(mockEC2, make(map[string]string))
			mockImportKeyPair(mockEC2, keyName, rsaFingerprint, rsaKey)

			keyName, err := LoadKeyByContent(context.Background(), &rsaKey, clusterName, ngName, mockEC2)

			Expect(err).NotTo(HaveOccurred())
			Expect(keyName).To(Equal("eksctl-sshtestcluster-nodegroup-ng1-f5:d9:01:88:1e:fb:40:fb:e1:ca:69:fe:2e:31:03:6c"))
			mockEC2.AssertCalled(GinkgoT(),
				"ImportKeyPair",
				mock.Anything,
				&ec2.ImportKeyPairInput{
					KeyName:           &keyName,
					PublicKeyMaterial: []byte(rsaKey),
				})
		})

		It("should not import key that already exists in EC2", func() {
			mockDescribeKeyPairs(mockEC2, map[string]string{keyName: rsaFingerprint})
			mockImportKeyPairError(mockEC2, errors.New("the key shouldn't be imported in this test"))

			keyName, err := LoadKeyByContent(context.Background(), &rsaKey, clusterName, ngName, mockEC2)

			Expect(err).NotTo(HaveOccurred())
			Expect(keyName).To(Equal("eksctl-sshtestcluster-nodegroup-ng1-f5:d9:01:88:1e:fb:40:fb:e1:ca:69:fe:2e:31:03:6c"))
			mockEC2.AssertNotCalled(GinkgoT(), "ImportKeyPair", mock.Anything, mock.Anything)
		})

		It("should return error if a key with same name exists in EC2 with different fingerprint", func() {
			differentFingerprint := "ab:cd"
			mockDescribeKeyPairs(mockEC2, map[string]string{keyName: differentFingerprint})
			mockImportKeyPairError(mockEC2, errors.New("the key shouldn't be imported in this test"))

			_, err := LoadKeyByContent(context.Background(), &rsaKey, clusterName, ngName, mockEC2)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("but fingerprints don't match"))
		})

	})

	Describe("deletion", func() {
		It("should return error if a key with same name exists in EC2 with different fingerprint", func() {
			keyToDelete1 := "eksctl-sshtestcluster-nodegroup-ng1-ab"
			keyToDelete2 := "eksctl-sshtestcluster-ef"
			existingKeys := map[string]string{
				keyToDelete1:                             "ab", // Must delete
				"eksctl-othercluster-nodegroup-ng1-ab":   "ab", // Different cluster
				"eksctl-sshtestcluster-nodegroup-ng1-bc": "22", // Wrong fingerprint
				keyToDelete2:                             "ef", // Must delete
				"not-a-cluster-key":                      "ab", // Not a key imported by eksctl
			}
			mockDescribeKeyPairs(mockEC2, existingKeys)
			mockDeleteKeyPair(mockEC2)

			DeleteKeys(context.Background(), mockEC2, clusterName)

			mockEC2.AssertNumberOfCalls(GinkgoT(), "DeleteKeyPair", 2)
			mockEC2.AssertCalled(GinkgoT(),
				"DeleteKeyPair",
				mock.Anything,
				&ec2.DeleteKeyPairInput{
					KeyName: &keyToDelete1,
				})
			mockEC2.AssertCalled(GinkgoT(),
				"DeleteKeyPair",
				mock.Anything,
				&ec2.DeleteKeyPairInput{
					KeyName: &keyToDelete2,
				})
		})
	})

	Describe("checking in EC2", func() {
		It("should not fail when key exits", func() {
			mockDescribeKeyPairs(mockEC2, map[string]string{keyName: rsaFingerprint})

			err := CheckKeyExistsInEC2(context.Background(), mockEC2, keyName)

			Expect(err).NotTo(HaveOccurred())
		})

		It("should fail when key does not exist", func() {
			mockDescribeKeyPairs(mockEC2, map[string]string{})

			err := CheckKeyExistsInEC2(context.Background(), mockEC2, keyName)
			Expect(err).To(HaveOccurred())
		})

		It("should fail when EC2 call fails", func() {

			mockDescribeKeyPairsError(mockEC2, &smithy.GenericAPIError{
				Code:    "testError",
				Message: "mock error for test EC2 call",
			})

			err := CheckKeyExistsInEC2(context.Background(), mockEC2, keyName)
			Expect(err).To(HaveOccurred())
		})
	})
})

func mockDeleteKeyPair(mockEC2 *mocksv2.EC2) {
	mockEC2.
		On("DeleteKeyPair", mock.Anything, mock.Anything).
		Return(&ec2.DeleteKeyPairOutput{}, nil)
}

func mockDescribeKeyPairsError(mockEC2 *mocksv2.EC2, err error) {
	mockEC2.
		On("DescribeKeyPairs", mock.Anything, mock.Anything).
		Return(nil, err)
}

func mockDescribeKeyPairs(mockEC2 *mocksv2.EC2, keys map[string]string) {

	if len(keys) == 0 {
		mockEC2.
			On("DescribeKeyPairs", mock.Anything, mock.Anything).
			Return(nil, &smithy.GenericAPIError{
				Code:    "InvalidKeyPair.NotFound",
				Message: "not found",
			})
		return
	}
	mockEC2.
		On("DescribeKeyPairs", mock.Anything, mock.Anything).
		Return(&ec2.DescribeKeyPairsOutput{
			KeyPairs: toKeyPairInfo(keys),
		}, nil)
}

func mockImportKeyPairError(mockEC2 *mocksv2.EC2, err error) {
	if err != nil {
		mockEC2.
			On("ImportKeyPair", mock.Anything, mock.Anything).
			Return(nil, err)
	}
}

func mockImportKeyPair(mockEC2 *mocksv2.EC2, keyName, fingerprint, key string) {
	mockEC2.
		On("ImportKeyPair",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.ImportKeyPairInput) bool {
				return *input.KeyName == keyName &&
					string(input.PublicKeyMaterial) == key
			})).
		Return(&ec2.ImportKeyPairOutput{
			KeyName:        &keyName,
			KeyFingerprint: &fingerprint,
		}, nil)
}

func toKeyPairInfo(keys map[string]string) []ec2types.KeyPairInfo {
	var keyPairs []ec2types.KeyPairInfo
	for k, v := range keys {
		keyName := k
		fingerprint := v
		keyPairs = append(keyPairs, ec2types.KeyPairInfo{
			KeyFingerprint: &fingerprint,
			KeyName:        &keyName,
		})
	}
	return keyPairs
}
