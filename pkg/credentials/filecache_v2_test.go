package credentials_test

import (
	"context"
	"os"
	"path"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscredentials "github.com/aws/aws-sdk-go/aws/credentials"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"

	"github.com/weaveworks/eksctl/pkg/credentials"
	"github.com/weaveworks/eksctl/pkg/credentials/fakes"
)

//counterfeiter:generate -o fakes/fake_aws_credentials_provider.go . provider
type provider interface {
	aws.CredentialsProvider
}

var cacheFilePath string

var _ = BeforeSuite(func() {
	homeDir, err := os.UserHomeDir()
	Expect(err).NotTo(HaveOccurred())
	cacheFilePath = path.Join(homeDir, ".eksctl", "cache", "credentials.yaml")
})

var _ = Describe("FileCacheV2", func() {

	type fileCacheEntry struct {
		createProvider func() provider
		setupCache     func(afero.Fs) error

		expectedRetrieveErr      string
		expectedErr              string
		expectedCacheData        string
		expectedCacheFileMissing bool
	}

	type cachedCredential struct {
		Credential awscredentials.Value
		Expiration time.Time
	}

	makeFlock := func(_ string) credentials.Flock {
		fl := &fakes.FakeFlock{}
		fl.TryRLockContextReturns(true, nil)
		fl.TryLockContextReturns(true, nil)
		fl.UnlockReturns(nil)
		return fl
	}

	DescribeTable("FileCacheV2 credentials caching", func(e fileCacheEntry) {
		fs := afero.NewMemMapFs()
		if e.setupCache != nil {
			Expect(e.setupCache(fs)).To(Succeed())
		}

		c := &fakes.FakeClock{}
		c.NowReturns(time.Date(42, 1, 1, 0, 0, 0, 0, time.UTC))

		f, err := credentials.NewFileCacheV2(e.createProvider(), "test", fs, makeFlock, c, cacheFilePath)
		if e.expectedErr != "" {
			Expect(err).To(MatchError(ContainSubstring(e.expectedErr)))
			return
		}
		Expect(err).NotTo(HaveOccurred())

		_, err = f.Retrieve(context.TODO())
		if e.expectedRetrieveErr != "" {
			Expect(err).To(MatchError(ContainSubstring(e.expectedRetrieveErr)))
			return
		}
		Expect(err).NotTo(HaveOccurred())

		data, err := afero.ReadFile(fs, cacheFilePath)
		if e.expectedCacheFileMissing {
			Expect(err).To(HaveOccurred())
			Expect(os.IsNotExist(err)).To(BeTrue())
			return
		}
		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).To(Equal(e.expectedCacheData))
	},
		Entry("credentials that can expire are cached", fileCacheEntry{
			createProvider: func() provider {
				p := &fakes.FakeProvider{}
				p.RetrieveReturns(aws.Credentials{
					AccessKeyID:     "k123",
					SecretAccessKey: "s123",
					SessionToken:    "t123",
					Source:          "eksctl-test",
					CanExpire:       true,
					Expires:         time.Time{},
				}, nil)
				return p
			},
			expectedCacheData: `profiles:
  test:
    credential:
      accesskeyid: k123
      secretaccesskey: s123
      sessiontoken: t123
      providername: eksctl-test
    expiration: 0001-01-01T00:00:00Z
`,
		}),

		Entry("cached credentials from file are used when available", fileCacheEntry{
			setupCache: func(fs afero.Fs) error {
				if err := fs.MkdirAll(cacheFilePath, 0700); err != nil {
					return err
				}

				data, err := yaml.Marshal(map[string]map[string]cachedCredential{
					"profiles": {
						"test": {
							Credential: awscredentials.Value{
								AccessKeyID:     "k123",
								SecretAccessKey: "s123",
								SessionToken:    "t123",
								ProviderName:    "eksctl-test",
							},
							Expiration: time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC),
						},
					},
				})
				if err != nil {
					return err
				}
				return afero.WriteFile(fs, cacheFilePath, data, 0700)
			},

			createProvider: func() provider {
				p := &fakes.FakeProvider{}
				p.RetrieveReturns(aws.Credentials{}, errors.New("unexpected call to Retrieve"))
				return p
			},
			expectedCacheData: `profiles:
  test:
    credential:
      accesskeyid: k123
      secretaccesskey: s123
      sessiontoken: t123
      providername: eksctl-test
    expiration: 9999-01-01T00:00:00Z
`,
		}),

		Entry("cached credentials from file are not used when expired", fileCacheEntry{
			setupCache: func(fs afero.Fs) error {
				if err := fs.MkdirAll(cacheFilePath, 0700); err != nil {
					return err
				}

				data, err := yaml.Marshal(map[string]map[string]cachedCredential{
					"profiles": {
						"test": {
							Credential: awscredentials.Value{
								AccessKeyID:     "k123",
								SecretAccessKey: "s123",
								SessionToken:    "t123",
								ProviderName:    "eksctl-test",
							},
							Expiration: time.Time{},
						},
					},
				})
				if err != nil {
					return err
				}
				return afero.WriteFile(fs, cacheFilePath, data, 0700)
			},

			createProvider: func() provider {
				p := &fakes.FakeProvider{}
				p.RetrieveReturns(aws.Credentials{
					AccessKeyID:     "a567",
					SecretAccessKey: "s567",
					SessionToken:    "t567",
					Source:          "eksctl-test",
					CanExpire:       true,
					Expires:         time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC),
				}, nil)
				return p
			},
			expectedCacheData: `profiles:
  test:
    credential:
      accesskeyid: a567
      secretaccesskey: s567
      sessiontoken: t567
      providername: eksctl-test
    expiration: 9999-01-01T00:00:00Z
`,
		}),

		Entry("cached credentials for a different profile are not used and non-expiring credentials are not cached", fileCacheEntry{
			setupCache: func(fs afero.Fs) error {
				if err := fs.MkdirAll(cacheFilePath, 0700); err != nil {
					return err
				}

				data, err := yaml.Marshal(map[string]map[string]cachedCredential{
					"profiles": {
						"eksctl": {
							Credential: awscredentials.Value{
								AccessKeyID:     "a123",
								SecretAccessKey: "s123",
								SessionToken:    "t123",
								ProviderName:    "eksctl-test",
							},
							Expiration: time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC),
						},
					},
				})
				if err != nil {
					return err
				}
				return afero.WriteFile(fs, cacheFilePath, data, 0700)
			},

			createProvider: func() provider {
				p := &fakes.FakeProvider{}
				p.RetrieveReturns(aws.Credentials{
					AccessKeyID:     "a999",
					SecretAccessKey: "s999",
					SessionToken:    "t999",
					Source:          "eksctl-test",
					CanExpire:       false,
				}, nil)
				return p
			},
			expectedCacheData: `profiles:
  eksctl:
    credential:
      accesskeyid: a123
      secretaccesskey: s123
      sessiontoken: t123
      providername: eksctl-test
    expiration: 9999-01-01T00:00:00Z
`,
		}),

		Entry("propagate error from provider", fileCacheEntry{
			createProvider: func() provider {
				return aws.AnonymousCredentials{}
			},
			expectedRetrieveErr: "not a valid credential provider",
		}),

		Entry("credentials that do not expire are not cached", fileCacheEntry{
			createProvider: func() provider {
				f := &fakes.FakeProvider{}
				f.RetrieveReturns(aws.Credentials{
					AccessKeyID:     "a123",
					SecretAccessKey: "s123",
					SessionToken:    "t123",
					Source:          "eksctl-test",
					CanExpire:       false,
				}, nil)
				return f
			},
			expectedCacheFileMissing: true,
		}),
	)
})
