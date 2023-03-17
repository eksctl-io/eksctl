package credentials_test

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gofrs/flock"

	"github.com/spf13/afero"

	"github.com/aws/aws-sdk-go/aws/credentials"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/pkg/credentials"
	"github.com/weaveworks/eksctl/pkg/credentials/fakes"
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
	newFileCacheProvider := func(profile string, c *credentials.Credentials, clock Clock, cacheDir string) (FileCacheProvider, error) {
		return NewFileCacheProvider(profile, c, clock, afero.NewOsFs(), func(path string) Flock {
			return flock.New(path)
		}, filepath.Join(cacheDir, "credentials.yaml"))
	}
	Context("credential cache has being used", func() {
		var (
			tmp string
			err error
		)
		BeforeEach(func() {
			tmp, err = os.MkdirTemp("", "filecache")
			Expect(err).NotTo(HaveOccurred())
			_ = os.Setenv(EksctlCacheFilenameEnvName, filepath.Join(tmp, "credentials.yaml"))
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
			p, err := newFileCacheProvider("profile", c, fakeClock, tmp)
			Expect(err).NotTo(HaveOccurred())
			value, err := p.Retrieve()
			Expect(err).NotTo(HaveOccurred())
			Expect(value.AccessKeyID).To(Equal("id"))
			Expect(value.SecretAccessKey).To(Equal("secret"))
			Expect(value.SessionToken).To(Equal("token"))
			Expect(p.IsExpired()).NotTo(BeTrue())
			content, err := os.ReadFile(filepath.Join(tmp, "credentials.yaml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(Equal(`profiles:
  profile:
    credential:
      accesskeyid: id
      secretaccesskey: secret
      sessiontoken: token
      providername: stubProvider
    expiration: 0001-01-01T00:00:00Z
`))
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
				p, err := newFileCacheProvider("profile", c, fakeClock, tmp)
				Expect(err).NotTo(HaveOccurred())
				Expect(p.IsExpired()).To(BeTrue())
			})
		})
		When("the cache file already exists", func() {
			It("will retrieve its content", func() {
				content := []byte(`profiles:
  profile:
    credential:
      accesskeyid: storedID
      secretaccesskey: storedSecret
      sessiontoken: storedToken
      providername: stubProvider
    expiration: 0001-01-01T00:00:00Z
`)
				err := os.WriteFile(filepath.Join(tmp, "credentials.yaml"), content, 0700)
				Expect(err).NotTo(HaveOccurred())
				c := credentials.NewCredentials(&stubProviderExpirer{})
				fakeClock := &fakes.FakeClock{}
				p, err := newFileCacheProvider("profile", c, fakeClock, tmp)
				Expect(err).NotTo(HaveOccurred())
				creds, err := p.Retrieve()
				Expect(err).NotTo(HaveOccurred())
				Expect(creds.AccessKeyID).To(Equal("storedID"))
				Expect(creds.SecretAccessKey).To(Equal("storedSecret"))
				Expect(creds.SessionToken).To(Equal("storedToken"))
			})
		})
		When("no underlying credentials have been supplied", func() {
			It("returns an appropriate error", func() {
				fakeClock := &fakes.FakeClock{}
				_, err := newFileCacheProvider("profile", nil, fakeClock, tmp)
				Expect(err).To(MatchError("no underlying Credentials object provided"))
			})
		})
		When("the underlying credentials provider doesn't support caching", func() {
			It("won't create a cache file", func() {
				fakeClock := &fakes.FakeClock{}
				fakeClock.NowReturns(time.Date(9999, 1, 1, 1, 1, 1, 1, time.UTC))
				p, err := newFileCacheProvider("profile", credentials.NewStaticCredentials("id", "secret", "token"), fakeClock, tmp)
				Expect(err).NotTo(HaveOccurred())
				_, err = p.Retrieve()
				Expect(err).NotTo(HaveOccurred())
				_, err = os.Stat(filepath.Join(tmp, "credentials.yaml"))
				Expect(os.IsNotExist(err)).To(BeTrue())
			})

		})
		When("the cache file's permission is too broad", func() {
			It("will refuse to use that file", func() {
				content := []byte(`test:`)
				err := os.WriteFile(filepath.Join(tmp, "credentials.yaml"), content, 0777)
				Expect(err).NotTo(HaveOccurred())
				c := credentials.NewCredentials(&stubProviderExpirer{})
				fakeClock := &fakes.FakeClock{}
				_, err = newFileCacheProvider("profile", c, fakeClock, tmp)
				Expect(err).To(MatchError(ContainSubstring(fmt.Sprintf("cache file %s is not private", filepath.Join(tmp, "credentials.yaml")))))
			})
		})
		When("the cache data has been corrupted", func() {
			It("will return an appropriate error", func() {
				content := []byte(`not valid yaml`)
				err := os.WriteFile(filepath.Join(tmp, "credentials.yaml"), content, 0600)
				Expect(err).NotTo(HaveOccurred())
				c := credentials.NewCredentials(&stubProviderExpirer{})
				fakeClock := &fakes.FakeClock{}
				_, err = newFileCacheProvider("profile", c, fakeClock, tmp)
				Expect(err).To(MatchError(ContainSubstring("unable to parse file")))
			})
		})
	})
})
