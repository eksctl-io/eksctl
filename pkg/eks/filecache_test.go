package eks_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/fakes"
)

type stubProvider struct {
	creds   credentials.Value
	expired bool
	err     error
}

func (s *stubProvider) Retrieve() (credentials.Value, error) {
	s.expired = false
	s.creds.ProviderName = "stubProvider"
	return s.creds, s.err
}

func (s *stubProvider) IsExpired() bool {
	return s.expired
}

type stubProviderExpirer struct {
	stubProvider
	expiration time.Time
}

func (s *stubProviderExpirer) ExpiresAt() time.Time {
	return s.expiration
}

var _ = Describe("filecache", func() {
	Context("a cached based provider is request", func() {
		var (
			tmp string
			err error
		)
		BeforeEach(func() {
			tmp, err = ioutil.TempDir("", "filecache")
			Expect(err).ToNot(HaveOccurred())
			os.Setenv(EksctlCacheFilenameEnvName, filepath.Join(tmp, "credentials.yaml"))
		})
		AfterEach(func() {
			_ = os.RemoveAll(tmp)
		})
		It("will provide a working file based cache", func() {
			c := credentials.NewCredentials(&stubProviderExpirer{
				stubProvider: stubProvider{
					creds: credentials.Value{
						AccessKeyID:     "id",
						SecretAccessKey: "secret",
						SessionToken:    "token",
						ProviderName:    "stubProvider",
					},
				},
			})
			fakeClock := &fakes.FakeClock{}
			fakeClock.NowReturns(time.Date(1981, 1, 1, 1, 1, 1, 1, time.UTC))
			p, err := NewFileCacheProvider("profile", c, fakeClock)
			Expect(err).ToNot(HaveOccurred())
			value, err := p.Retrieve()
			Expect(err).ToNot(HaveOccurred())
			Expect(value.AccessKeyID).To(Equal("id"))
			Expect(value.SecretAccessKey).To(Equal("secret"))
			Expect(value.SessionToken).To(Equal("token"))
			Expect(p.IsExpired()).NotTo(BeTrue())
			content, err := ioutil.ReadFile(filepath.Join(tmp, "credentials.yaml"))
			Expect(err).ToNot(HaveOccurred())
			Expect(string(content)).To(Equal(`profiles:
  profile:
    credential:
      accesskeyid: id
      secretaccesskey: secret
      sessiontoken: token
      providername: stubProvider
    expiration: 0001-01-01T00:00:00Z
`))
			Expect(p.IsExpired()).NotTo(BeTrue())
		})
		When("the cache expires", func() {
			It("will ask to refresh it", func() {
				c := credentials.NewCredentials(&stubProviderExpirer{
					stubProvider: stubProvider{
						creds: credentials.Value{
							AccessKeyID:     "id",
							SecretAccessKey: "secret",
							SessionToken:    "token",
							ProviderName:    "stubProvider",
						},
					},
				})
				fakeClock := &fakes.FakeClock{}
				fakeClock.NowReturns(time.Date(9999, 1, 1, 1, 1, 1, 1, time.UTC))
				p, err := NewFileCacheProvider("profile", c, fakeClock)
				Expect(err).ToNot(HaveOccurred())
				Expect(p.IsExpired()).To(BeTrue())
			})
		})
	})
})
