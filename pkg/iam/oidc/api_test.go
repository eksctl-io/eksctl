package iamoidc

import (
	"crypto/sha1"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awsiam "github.com/aws/aws-sdk-go/service/iam"
	"github.com/stretchr/testify/mock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var thumbprint string

var _ = BeforeSuite(func() {
	session, err := gexec.Start(exec.Command("make", "-C", "testdata", "all"), GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, 3).Should(gexec.Exit())
	rawCert, err := ioutil.ReadFile("testdata/test-server.pem")
	Expect(err).NotTo(HaveOccurred())
	block, rest := pem.Decode(rawCert)
	Expect(rest).To(BeEmpty())
	thumbprint = fmt.Sprintf("%x", sha1.Sum(block.Bytes))
})

var _ = AfterSuite(func() {
	session, err := gexec.Start(exec.Command("make", "-C", "testdata", "clean"), GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, 3).Should(gexec.Exit())
})

var _ = Describe("EKS/IAM API wrapper", func() {
	const (
		exampleIssuer   = "https://exampleIssuer.eksctl.io/id/13EBFE0C5BD60778E91DFE559E02689C"
		fakeProviderARN = "arn:aws:iam::12345:oidc-provider/localhost/"
	)

	Describe("parse OIDC issuer URL and host fingerprint", func() {
		var (
			p *mockprovider.MockProvider

			err error
		)

		BeforeEach(func() {
			p = mockprovider.NewMockProvider()
		})

		It("should get cluster, cache status and get issuer URL", func() {
			oidc, err := NewOpenIDConnectManager(p.IAM(), "12345", exampleIssuer, "aws", nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(oidc.issuerURL.Port()).To(Equal("443"))
			Expect(oidc.issuerURL.Hostname()).To(Equal("exampleIssuer.eksctl.io"))
		})

		It("should handle bad issuer URL", func() {
			_, err = NewOpenIDConnectManager(p.IAM(), "12345", "http://foo\x7f.com/", "aws", nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix("parsing OIDC issuer URL"))
		})

		It("should handle bad issuer URL scheme", func() {
			_, err = NewOpenIDConnectManager(p.IAM(), "12345", "http://foo.com/", "aws", nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix("unsupported URL scheme"))
		})

		It("should get cluster, and fail to connect to fake issue URL", func() {
			oidc, err := NewOpenIDConnectManager(p.IAM(), "12345", "https://localhost:10020/", "aws", nil)
			Expect(err).NotTo(HaveOccurred())

			err = oidc.getIssuerCAThumbprint()
			Expect(err).To(HaveOccurred())
			// Use regex to match URL as go 1.14 have extra double quotes for URL
			// related commit : https://github.com/golang/go/commit/64cfe9fe22113cd6bc05a2c5d0cbe872b1b57860
			Expect(err.Error()).To(MatchRegexp("connecting to issuer OIDC: Get \"?https://localhost:10020/\"?"))
			Expect(err.Error()).To(HaveSuffix("connect: connection refused"))
		})

		It("should get OIDC issuer's CA fingerprint", func() {
			oidc, err := NewOpenIDConnectManager(p.IAM(), "12345", "https://localhost:10028/", "aws", nil)
			Expect(err).NotTo(HaveOccurred())

			srv, err := newServer(oidc.issuerURL.Host)
			Expect(err).NotTo(HaveOccurred())

			go func() {
				_ = srv.serve()
			}()

			oidc.insecureSkipVerify = true

			err = oidc.getIssuerCAThumbprint()
			Expect(srv.close()).To(Succeed())
			Expect(err).NotTo(HaveOccurred())

			Expect(oidc.issuerCAThumbprint).ToNot(BeEmpty())
			Expect(oidc.issuerCAThumbprint).To(Equal(thumbprint))
		})

		It("should get OIDC issuer's CA fingerprint for a URL that returns 403", func() {
			oidc, err := NewOpenIDConnectManager(p.IAM(), "12345", "https://localhost:10029/fake_eks", "aws", nil)
			Expect(err).NotTo(HaveOccurred())

			srv, err := newServer(oidc.issuerURL.Host)
			Expect(err).NotTo(HaveOccurred())

			go func() {
				_ = srv.serve()
			}()

			oidc.insecureSkipVerify = true

			err = oidc.getIssuerCAThumbprint()
			Expect(srv.close()).To(Succeed())
			Expect(err).NotTo(HaveOccurred())

			Expect(oidc.issuerCAThumbprint).ToNot(BeEmpty())
			Expect(oidc.issuerCAThumbprint).To(Equal(thumbprint))
		})
	})

	Describe("create/get/delete tests", func() {
		var (
			p    *mockprovider.MockProvider
			srv  *testServer
			oidc *OpenIDConnectManager

			err error

			fakeProviderCreated = new(bool)
		)

		BeforeEach(func() {
			p = mockprovider.NewMockProvider()

			if fakeProviderCreated == nil {
				*fakeProviderCreated = false
			}

			nonExistentProviderErr := awserr.New(awsiam.ErrCodeNoSuchEntityException, "provider is not there", fmt.Errorf("test"))

			fakeProviderGetOutput := &awsiam.GetOpenIDConnectProviderOutput{
				Url: aws.String("https://localhost:10028/"),
			}

			fakeProviderCreateOutput := &awsiam.CreateOpenIDConnectProviderOutput{
				OpenIDConnectProviderArn: aws.String(fakeProviderARN),
			}

			p.MockIAM().On("GetOpenIDConnectProvider", mock.MatchedBy(func(input *awsiam.GetOpenIDConnectProviderInput) bool {
				return *input.OpenIDConnectProviderArn == exampleIssuer
			})).Return(nil, nonExistentProviderErr)

			p.MockIAM().On("GetOpenIDConnectProvider", mock.MatchedBy(func(input *awsiam.GetOpenIDConnectProviderInput) bool {
				return *input.OpenIDConnectProviderArn == fakeProviderARN && !*fakeProviderCreated
			})).Return(nil, nonExistentProviderErr)

			p.MockIAM().On("GetOpenIDConnectProvider", mock.MatchedBy(func(input *awsiam.GetOpenIDConnectProviderInput) bool {
				return *input.OpenIDConnectProviderArn == fakeProviderARN && *fakeProviderCreated
			})).Return(fakeProviderGetOutput, nil)

			p.MockIAM().On("CreateOpenIDConnectProvider", mock.MatchedBy(func(input *awsiam.CreateOpenIDConnectProviderInput) bool {
				if *input.Url == *fakeProviderGetOutput.Url {
					*fakeProviderCreated = true
					return true
				}
				return false
			})).Return(fakeProviderCreateOutput, nil)

			p.MockIAM().On("DeleteOpenIDConnectProvider", mock.MatchedBy(func(input *awsiam.DeleteOpenIDConnectProviderInput) bool {
				if *input.OpenIDConnectProviderArn == fakeProviderARN {
					*fakeProviderCreated = false
					return true
				}
				return false
			})).Return(nil, nil)

			oidc, err = NewOpenIDConnectManager(p.IAM(), "12345", "https://localhost:10028/", "aws", nil)
			Expect(err).NotTo(HaveOccurred())

			srv, err = newServer(oidc.issuerURL.Host)
			Expect(err).NotTo(HaveOccurred())

			go func() {
				_ = srv.serve()
			}()

			oidc.insecureSkipVerify = true

			err = oidc.CreateProvider()
			Expect(err).NotTo(HaveOccurred())

			exists, err := oidc.CheckProviderExists()
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())
			Expect(oidc.ProviderARN).To(Equal(fakeProviderARN))

		})

		AfterEach(func() {
			Expect(srv.close()).To(Succeed())
			Expect(oidc.DeleteProvider()).NotTo(HaveOccurred())
		})

		It("delete existing OIDC provider and check it no longer exists", func() {
			err = oidc.DeleteProvider()
			Expect(err).NotTo(HaveOccurred())

			exists, err := oidc.CheckProviderExists()
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())
		})

		It("should construct assume role policy document for a service account", func() {
			exists, err := oidc.CheckProviderExists()
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())

			document := oidc.MakeAssumeRolePolicyDocumentWithServiceAccountConditions("test-ns1", "test-sa1")
			Expect(document).ToNot(BeEmpty())

			expected := `{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Effect": "Allow",
						"Principal": {
							"Federated": "` + fakeProviderARN + `"
						},
						"Action": ["sts:AssumeRoleWithWebIdentity"],
						"Condition": {
							"StringEquals": {
								"localhost/:sub": "system:serviceaccount:test-ns1:test-sa1",
								"localhost/:aud": "sts.amazonaws.com"
							}
						}
					}
				]
			}`

			js, err := json.Marshal(document)
			Expect(err).NotTo(HaveOccurred())
			Expect(js).To(MatchJSON(expected))
		})

		It("should construct assume role policy document", func() {
			exists, err := oidc.CheckProviderExists()
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())

			document := oidc.MakeAssumeRolePolicyDocument()
			Expect(document).ToNot(BeEmpty())

			expected := `{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Effect": "Allow",
						"Principal": {
							"Federated": "` + fakeProviderARN + `"
						},
						"Action": ["sts:AssumeRoleWithWebIdentity"],
						"Condition": {
							"StringEquals": {
								"localhost/:aud": "sts.amazonaws.com"
							}
						}
					}
				]
			}`

			js, err := json.Marshal(document)
			Expect(err).NotTo(HaveOccurred())
			Expect(js).To(MatchJSON(expected))
		})

	})

	Describe("Tags support", func() {
		var (
			provider *mockprovider.MockProvider
			srv      *testServer
		)

		BeforeEach(func() {
			provider = mockprovider.NewMockProvider()
			var err error
			srv, err = newServer("localhost:10028")
			Expect(err).NotTo(HaveOccurred())
			go func() {
				_ = srv.serve()
			}()
		})

		JustAfterEach(func() {
			srv.close()
		})

		It("should tag OIDC resources", func() {
			oidc, err := NewOpenIDConnectManager(provider.IAM(), "12345", "https://localhost:10028/", "aws", map[string]string{
				"cluster":  "oidc",
				"resource": "oidc-provider",
			})
			Expect(err).ToNot(HaveOccurred())
			oidc.insecureSkipVerify = true

			var tagsInput []*awsiam.Tag
			provider.MockIAM().On("CreateOpenIDConnectProvider", mock.MatchedBy(func(input *awsiam.CreateOpenIDConnectProviderInput) bool {
				tagsInput = input.Tags
				return true
			})).Return(&awsiam.CreateOpenIDConnectProviderOutput{
				OpenIDConnectProviderArn: aws.String(fakeProviderARN),
			}, nil)

			Expect(oidc.CreateProvider()).To(Succeed())
			Expect(tagsInput).To(ConsistOf([]*awsiam.Tag{
				{
					Key:   aws.String("cluster"),
					Value: aws.String("oidc"),
				},
				{
					Key:   aws.String("resource"),
					Value: aws.String("oidc-provider"),
				},
			}))
		})

	})

	Describe("OIDC AWS partition test", func() {
		var (
			provider *mockprovider.MockProvider
			srv      *testServer
		)

		BeforeEach(func() {
			provider = mockprovider.NewMockProvider()
			var err error
			srv, err = newServer("localhost:10028")
			Expect(err).NotTo(HaveOccurred())
			go func() {
				_ = srv.serve()
			}()
		})

		JustAfterEach(func() {
			srv.close()
		})

		DescribeTable("AssumeRolePolicyDocument should have correct AWS partition and STS domain", func(partition, expectedAudience string) {
			provider.MockIAM().On("CreateOpenIDConnectProvider", mock.MatchedBy(func(input *awsiam.CreateOpenIDConnectProviderInput) bool {
				if len(input.ClientIDList) != 1 {
					return false
				}
				clientID := *input.ClientIDList[0]
				return clientID == defaultAudience
			})).Return(&awsiam.CreateOpenIDConnectProviderOutput{
				OpenIDConnectProviderArn: aws.String(fmt.Sprintf("arn:%s:iam::12345:oidc-provider/localhost/", partition)),
			}, nil)

			oidc, err := NewOpenIDConnectManager(provider.IAM(), "12345", "https://localhost:10028/", partition, nil)
			oidc.insecureSkipVerify = true
			Expect(err).ToNot(HaveOccurred())
			Expect(oidc.CreateProvider()).To(Succeed())

			document := oidc.MakeAssumeRolePolicyDocumentWithServiceAccountConditions("test-ns", "test-sa")
			expected := fmt.Sprintf(`{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Effect": "Allow",
						"Principal": {
							"Federated": %q
						},
						"Action": ["sts:AssumeRoleWithWebIdentity"],
						"Condition": {
							"StringEquals": {
								"localhost/:sub": "system:serviceaccount:test-ns:test-sa",
								"localhost/:aud": %q
							}
						}
					}
				]
			}`, fmt.Sprintf("arn:%s:iam::12345:oidc-provider/localhost/", partition), expectedAudience)

			actual, err := json.Marshal(document)
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(MatchJSON(expected))
		},
			Entry("Default AWS partition", "aws", "sts.amazonaws.com"),
			Entry("AWS China partition", "aws-cn", "sts.amazonaws.com"),
		)
	})
})

type testServer struct {
	listener net.Listener
	server   *http.Server
}

func newServer(host string) (*testServer, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "ok.") })

	mux.HandleFunc("/fake_eks", func(w http.ResponseWriter, r *http.Request) {
		// this is what EKS normally returns to us, as we are an unauthenticated client
		w.WriteHeader(403)
		fmt.Fprintf(w, `{"message":"Missing Authentication Token"}`)
	})

	fmt.Fprintf(GinkgoWriter, "\nserver will listen on %q\n", host)

	// we must construct listener to avoid race condition, as simply calling
	// `go srv.ListenAndServeTLS` doesn't guarantee that sever will be listening
	// right away and can be tested, so we make sure to listen before we return
	listener, err := net.Listen("tcp", host)
	if err != nil {
		return nil, err
	}

	return &testServer{
		listener: listener,
		server: &http.Server{
			Addr:    host,
			Handler: mux,
		},
	}, nil
}

func (s *testServer) serve() error {
	return s.server.ServeTLS(s.listener, "testdata/test-server.pem", "testdata/test-server-key.pem")
}

func (s *testServer) close() error {
	return s.server.Close()
}
