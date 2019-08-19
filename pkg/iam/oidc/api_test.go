package iamoidc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awsiam "github.com/aws/aws-sdk-go/service/iam"
	"github.com/stretchr/testify/mock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	//api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = Describe("EKS/IAM API wrapper", func() {
	Describe("can get OIDC issuer URL and host fingerprint", func() {
		var (
			p *mockprovider.MockProvider
			// ctl *eks.ClusterProvider
			// cfg *api.ClusterConfig
			err error

			fakeProviderCreated = new(bool)
		)

		const (
			exampleIssuer   = "https://exampleIssuer.eksctl.io/id/13EBFE0C5BD60778E91DFE559E02689C"
			fakeProviderARN = "arn:aws:iam::12345:oidc-provider/localhost/"
		)

		BeforeEach(func() {
			p = mockprovider.NewMockProvider()

			if fakeProviderCreated == nil {
				*fakeProviderCreated = false
			}

			nonExistentProviderErr := awserr.New(awsiam.ErrCodeNoSuchEntityException, "provider is not there", fmt.Errorf("test"))

			fakeProviderGetOutput := &awsiam.GetOpenIDConnectProviderOutput{
				Url: aws.String("https://localhost:8443/"),
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
			oidc, err := NewOpenIDConnectManager(p.IAM(), "12345", "https://localhost:10020/")
			Expect(err).NotTo(HaveOccurred())

			err = oidc.getIssuerCAThumbprint()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix("connecting to issuer OIDC (https://localhost:10020/): dial tcp"))
		})

		It("should get OIDC issuer's CA fingerprint", func() {
			oidc, err := NewOpenIDConnectManager(p.IAM(), "12345", "https://localhost:6443/")
			Expect(err).NotTo(HaveOccurred())

			srv := newServer(oidc.issuerURL.Host)
			go func() {
				defer GinkgoRecover()
				err = srv.ListenAndServeTLS("testdata/test-server.pem", "testdata/test-server-key.pem")
				Expect(err).NotTo(HaveOccurred())
			}()

			oidc.insecureSkipVerify = true

			err = oidc.getIssuerCAThumbprint()
			Expect(err).NotTo(HaveOccurred())

			Expect(oidc.issuerCAThumbprint).ToNot(BeEmpty())

			Expect(oidc.issuerCAThumbprint).To(Equal("8b453cc675feb77c65163b7a9907d77994386664"))
		})

		It("should create OIDC provider", func() {
			oidc, err := NewOpenIDConnectManager(p.IAM(), "12345", "https://localhost:8443/")
			Expect(err).NotTo(HaveOccurred())

			srv := newServer(oidc.issuerURL.Host)
			go func() {
				defer GinkgoRecover()
				err = srv.ListenAndServeTLS("testdata/test-server.pem", "testdata/test-server-key.pem")
				Expect(err).NotTo(HaveOccurred())
			}()

			oidc.insecureSkipVerify = true

			exists, err := oidc.CheckProviderExists()
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())

			Expect(oidc.ProviderARN).To(BeEmpty())

			err = oidc.CreateProvider()
			Expect(err).NotTo(HaveOccurred())

			Expect(oidc.ProviderARN).To(Equal(fakeProviderARN))
		})

		It("should check OIDC provider exists, delete it and check again", func() {
			oidc, err := NewOpenIDConnectManager(p.IAM(), "12345", "https://localhost:8443/")
			Expect(err).NotTo(HaveOccurred())

			Expect(oidc.ProviderARN).To(BeEmpty())

			exists, err := oidc.CheckProviderExists()
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())

			Expect(oidc.ProviderARN).To(Equal(fakeProviderARN))

			// TODO
			// - m.DeleteProvider

			// exists, err := oidc.CheckProviderExists()
			// Expect(err).NotTo(HaveOccurred())
			// Expect(exists).To(BeFalse())
		})

		It("should construct assume role policy document for a service account", func() {
			oidc, err := NewOpenIDConnectManager(p.IAM(), "12345", "https://localhost:8443/")
			Expect(err).NotTo(HaveOccurred())

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

func newServer(host string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "ok.") })

	return &http.Server{
		Addr:    host,
		Handler: mux,
	}
}
