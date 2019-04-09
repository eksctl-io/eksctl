package ssh

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/stretchr/testify/mock"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("ssh public key", func() {
	var (
		clusterName  = "sshtestcluster"
		ngName       = "ng1"
		key          = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDcSoNjWaJaw+MYBz43lgm12ZGdP+zRs9o0sXAGbiQua6e3JSkAiH4p9YZHmWxCTjckbiEdXN5qcs5OC5KUxYBvnEgor7jEydcKe1ZJXqsm/8CrtnJMTNcO9QVFnXfjvpkNjgNYj+8w9PcFRr0JDgDhRb52JvPWoqywv/Om9s1hpUov0gxDIl6CLLHSk0lmXZEhtVMMJmo0Tu/NlHqdky2DxFgHyNjBcMNpiBd8bs3dA5xf36dY+qgcXBV23i1SCgbqn9xcw1Q0IrHuQ4/QB+PJ5haxUx0bnOTahxSZ+tlEz9EiLwlM8VtKo3ND/giBvGaXuIK2iGDL0kSCRjueM5/3 user@example\n"
		keyName      = "eksctl-sshtestcluster-nodegroup-ng1-f5:d9:01:88:1e:fb:40:fb:e1:ca:69:fe:2e:31:03:6c"
		fingerprint  = "f5:d9:01:88:1e:fb:40:fb:e1:ca:69:fe:2e:31:03:6c"
		mockProvider *mockprovider.MockProvider
	)

	BeforeEach(func() {
		mockProvider = mockprovider.NewMockProvider()
	})

	Describe("loading from a file", func() {

		It("should import the key", func() {
			mockDescribeKeyPairs(mockProvider, make(map[string]string))
			mockImportKeyPair(mockProvider, keyName, fingerprint, key)

			err := LoadSSHKeyFromFile("assets/id_rsa_tests1.pub", clusterName, ngName, *mockProvider)

			Expect(err).ToNot(HaveOccurred())
			mockProvider.MockEC2().AssertCalled(GinkgoT(),
				"ImportKeyPair",
				&ec2.ImportKeyPairInput{
					KeyName:           &keyName,
					PublicKeyMaterial: []byte(key),
				})
		})

		It("should not import key that already exists in EC2", func() {
			mockDescribeKeyPairs(mockProvider, map[string]string{keyName: fingerprint})
			mockImportKeyPairError(mockProvider, errors.New("the key shouldn't be imported in this test"))

			err := LoadSSHKeyFromFile("assets/id_rsa_tests1.pub", clusterName, ngName, *mockProvider)

			Expect(err).ToNot(HaveOccurred())
			mockProvider.MockEC2().AssertNotCalled(GinkgoT(), "ImportKeyPair", mock.Anything)
		})

		It("should return error if a key with same name exists in EC2 with different fingerprint", func() {
			differentFingerprint := "ab:cd"
			mockDescribeKeyPairs(mockProvider, map[string]string{keyName: differentFingerprint})
			mockImportKeyPairError(mockProvider, errors.New("the key shouldn't be imported in this test"))

			err := LoadSSHKeyFromFile("assets/id_rsa_tests1.pub", clusterName, ngName, *mockProvider)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("but fingerprints don't match"))
		})

		It("should return error if the file does not exist", func() {
			err := LoadSSHKeyFromFile("assets/file_not_existing.pub", clusterName, ngName, *mockProvider)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("loading by the key content", func() {
		It("should import it", func() {
			mockDescribeKeyPairs(mockProvider, make(map[string]string))
			mockImportKeyPair(mockProvider, keyName, fingerprint, key)

			err := LoadSSHKeyByContent(&key, clusterName, ngName, *mockProvider)

			Expect(err).ToNot(HaveOccurred())
			mockProvider.MockEC2().AssertCalled(GinkgoT(),
				"ImportKeyPair",
				&ec2.ImportKeyPairInput{
					KeyName:           &keyName,
					PublicKeyMaterial: []byte(key),
				})
		})

		It("should not import key that already exists in EC2", func() {
			mockDescribeKeyPairs(mockProvider, map[string]string{keyName: fingerprint})
			mockImportKeyPairError(mockProvider, errors.New("the key shouldn't be imported in this test"))

			err := LoadSSHKeyByContent(&key, clusterName, ngName, *mockProvider)

			Expect(err).ToNot(HaveOccurred())
			mockProvider.MockEC2().AssertNotCalled(GinkgoT(), "ImportKeyPair", mock.Anything)
		})

		It("should return error if a key with same name exists in EC2 with different fingerprint", func() {
			differentFingerprint := "ab:cd"
			mockDescribeKeyPairs(mockProvider, map[string]string{keyName: differentFingerprint})
			mockImportKeyPairError(mockProvider, errors.New("the key shouldn't be imported in this test"))

			err := LoadSSHKeyByContent(&key, clusterName, ngName, *mockProvider)

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
			mockDescribeKeyPairs(mockProvider, existingKeys)
			mockDeleteKeyPair(mockProvider)

			DeletePublicSSHKeys(clusterName, mockProvider)

			mockProvider.MockEC2().AssertNumberOfCalls(GinkgoT(), "DeleteKeyPair", 2)
			mockProvider.MockEC2().AssertCalled(GinkgoT(),
				"DeleteKeyPair",
				&ec2.DeleteKeyPairInput{
					KeyName: &keyToDelete1,
				})
			mockProvider.MockEC2().AssertCalled(GinkgoT(),
				"DeleteKeyPair",
				&ec2.DeleteKeyPairInput{
					KeyName: &keyToDelete2,
				})
		})
	})

	Describe("checking in EC2", func() {
		It("should not fail when key exits", func() {
			mockDescribeKeyPairs(mockProvider, map[string]string{keyName: fingerprint})

			err := CheckKeyExistsInEc2(keyName, mockProvider)

			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail when key does not exist", func() {
			mockDescribeKeyPairs(mockProvider, map[string]string{})

			err := CheckKeyExistsInEc2(keyName, mockProvider)

			Expect(err).To(HaveOccurred())
		})

		It("should fail when EC2 call fails", func() {
			mockDescribeKeyPairsError(mockProvider, awserr.New("testError", "mock error for test EC2 call", nil))

			err := CheckKeyExistsInEc2(keyName, mockProvider)

			Expect(err).To(HaveOccurred())
		})
	})
})

func mockDeleteKeyPair(mockProvider *mockprovider.MockProvider) {
	mockProvider.MockEC2().
		On("DeleteKeyPair", mock.Anything).
		Return(&ec2.DeleteKeyPairOutput{}, nil)
}

func mockDescribeKeyPairsError(provider *mockprovider.MockProvider, err error) {
	provider.MockEC2().
		On("DescribeKeyPairs", mock.Anything).
		Return(nil, err)
}

func mockDescribeKeyPairs(provider *mockprovider.MockProvider, keys map[string]string) {
	if len(keys) == 0 {
		provider.MockEC2().
			On("DescribeKeyPairs", mock.Anything).
			Return(nil, awserr.New("InvalidKeyPair.NotFound", "not found", nil))
		return
	}
	provider.MockEC2().
		On("DescribeKeyPairs", mock.Anything).
		Return(&ec2.DescribeKeyPairsOutput{
			KeyPairs: toKeyPairInfo(keys),
		}, nil)
}

func mockImportKeyPairError(mockProvider *mockprovider.MockProvider, err error) {
	if err != nil {
		mockProvider.MockEC2().
			On("ImportKeyPair", mock.Anything).
			Return(nil, err)
		return
	}
}

func mockImportKeyPair(mockProvider *mockprovider.MockProvider, keyName, fingerprint, key string) {
	mockProvider.MockEC2().
		On("ImportKeyPair",
			mock.MatchedBy(func(input *ec2.ImportKeyPairInput) bool {
				return *input.KeyName == keyName &&
					string(input.PublicKeyMaterial) == key
			})).
		Return(&ec2.ImportKeyPairOutput{
			KeyName:        &keyName,
			KeyFingerprint: &fingerprint,
		}, nil)
}

func toKeyPairInfo(keys map[string]string) []*ec2.KeyPairInfo {
	var keyPairs []*ec2.KeyPairInfo
	for k, v := range keys {
		keyName := k
		fingerprint := v
		keyPairs = append(keyPairs, &ec2.KeyPairInfo{
			KeyFingerprint: &fingerprint,
			KeyName:        &keyName,
		})
	}
	return keyPairs
}
