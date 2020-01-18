package iamoidc

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awsiam "github.com/aws/aws-sdk-go/service/iam"
	"github.com/stretchr/testify/mock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

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
			oidc, err := NewOpenIDConnectManager(p.IAM(), "12345", exampleIssuer)
			Expect(err).NotTo(HaveOccurred())
			Expect(oidc.issuerURL.Port()).To(Equal("443"))
			Expect(oidc.issuerURL.Hostname()).To(Equal("exampleIssuer.eksctl.io"))
		})

		It("should handle bad issuer URL", func() {
			_, err = NewOpenIDConnectManager(p.IAM(), "12345", "http://foo\x7f.com/")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix("parsing OIDC issuer URL"))
		})

		It("should handle bad issuer URL scheme", func() {
			_, err = NewOpenIDConnectManager(p.IAM(), "12345", "http://foo.com/")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix("unsupported URL scheme"))
		})

		It("should get cluster, and fail to connect to fake issue URL", func() {
			oidc, err := NewOpenIDConnectManager(p.IAM(), "12345", "https://127.0.01:10020/")
			Expect(err).NotTo(HaveOccurred())

			err = oidc.getIssuerCAThumbprint()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix("connecting to issuer OIDC: Get https://127.0.01:10020/: dial tcp 127.0.0.1:10020: connect: connection refused"))
		})

		It("should get OIDC issuer's CA fingerprint", func() {
			oidc, err := NewOpenIDConnectManager(p.IAM(), "12345", "https://localhost:10028/")
			Expect(err).NotTo(HaveOccurred())

			srv, err := newServer(oidc.issuerURL.Host)
			Expect(err).NotTo(HaveOccurred())

			go srv.serve()

			oidc.insecureSkipVerify = true

			err = oidc.getIssuerCAThumbprint()
			Expect(err).NotTo(HaveOccurred())

			Expect(oidc.issuerCAThumbprint).ToNot(BeEmpty())

			Expect(oidc.issuerCAThumbprint).To(Equal("8b453cc675feb77c65163b7a9907d77994386664"))

			Expect(srv.close()).To(Succeed())
		})

		It("should get OIDC issuer's CA fingerprint for a URL that returns 403", func() {
			oidc, err := NewOpenIDConnectManager(p.IAM(), "12345", "https://localhost:10029/fake_eks")
			Expect(err).NotTo(HaveOccurred())

			srv, err := newServer(oidc.issuerURL.Host)
			Expect(err).NotTo(HaveOccurred())

			go srv.serve()

			oidc.insecureSkipVerify = true

			err = oidc.getIssuerCAThumbprint()
			Expect(err).NotTo(HaveOccurred())

			Expect(oidc.issuerCAThumbprint).ToNot(BeEmpty())

			Expect(oidc.issuerCAThumbprint).To(Equal("8b453cc675feb77c65163b7a9907d77994386664"))

			Expect(srv.close()).To(Succeed())
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
		})

		JustBeforeEach(func() {
			oidc, err = NewOpenIDConnectManager(p.IAM(), "12345", "https://localhost:10028/")
			Expect(err).NotTo(HaveOccurred())

			srv, err = newServer(oidc.issuerURL.Host)
			Expect(err).NotTo(HaveOccurred())

			go srv.serve()

			{
				exists, err := oidc.CheckProviderExists()
				Expect(err).NotTo(HaveOccurred())
				Expect(exists).To(BeFalse())
				Expect(oidc.ProviderARN).To(BeEmpty())
			}

			oidc.insecureSkipVerify = true

			err = oidc.CreateProvider()
			Expect(err).NotTo(HaveOccurred())

			{
				exists, err := oidc.CheckProviderExists()
				Expect(err).NotTo(HaveOccurred())
				Expect(exists).To(BeTrue())
				Expect(oidc.ProviderARN).To(Equal(fakeProviderARN))
			}

		})

		JustAfterEach(func() {
			Expect(srv.close()).To(Succeed())
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

			document := oidc.MakeAssumeRolePolicyDocument("test-ns1", "test-sa1")
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
