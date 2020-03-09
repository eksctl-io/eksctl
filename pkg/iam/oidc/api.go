package iamoidc

import (
	"crypto/sha1"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awsiam "github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/pkg/errors"

	cft "github.com/weaveworks/eksctl/pkg/cfn/template"
)

var defaultAudience = "sts.amazonaws.com"

// OpenIDConnectManager hold information about IAM OIDC integration
type OpenIDConnectManager struct {
	accountID string
	partition string
	audience  string

	issuerURL          *url.URL
	insecureSkipVerify bool
	issuerCAThumbprint string

	ProviderARN string

	iam iamiface.IAMAPI
}

// NewOpenIDConnectManager construct a new IAM OIDC management instance, it can return and error
// when the given issue URL was invalid
func NewOpenIDConnectManager(iamapi iamiface.IAMAPI, accountID, issuer, partition string) (*OpenIDConnectManager, error) {
	issuerURL, err := url.Parse(issuer)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing OIDC issuer URL")
	}

	if issuerURL.Scheme != "https" {
		return nil, fmt.Errorf("unsupported URL scheme %q", issuerURL.Scheme)
	}

	audience := defaultAudience

	if issuerURL.Port() == "" {
		issuerURL.Host += ":443"
	}

	m := &OpenIDConnectManager{
		iam:       iamapi,
		accountID: accountID,
		partition: partition,
		audience:  audience,
		issuerURL: issuerURL,
	}
	return m, nil
}

// CheckProviderExists will return true when the provider exists, it may return errors
// if it was unable to call IAM API
func (m *OpenIDConnectManager) CheckProviderExists() (bool, error) {
	input := &awsiam.GetOpenIDConnectProviderInput{
		OpenIDConnectProviderArn: aws.String(
			fmt.Sprintf("arn:%s:iam::%s:oidc-provider/%s", m.partition, m.accountID, m.hostnameAndPath()),
		),
	}
	_, err := m.iam.GetOpenIDConnectProvider(input)
	if err != nil {
		awsError := err.(awserr.Error)
		if awsError.Code() == awsiam.ErrCodeNoSuchEntityException {
			return false, nil
		}
		return false, err
	}
	m.ProviderARN = *input.OpenIDConnectProviderArn
	return true, nil
}

// CreateProvider will retrieve CA root certificate and compute its thumbprint for the
// by connecting to it and create the provider using IAM API
func (m *OpenIDConnectManager) CreateProvider() error {
	if err := m.getIssuerCAThumbprint(); err != nil {
		return err
	}
	input := &awsiam.CreateOpenIDConnectProviderInput{
		ClientIDList:   aws.StringSlice([]string{m.audience}),
		ThumbprintList: []*string{&m.issuerCAThumbprint},
		// It has no name or tags, it's keyed to the URL
		Url: aws.String(m.issuerURL.String()),
	}
	output, err := m.iam.CreateOpenIDConnectProvider(input)
	if err != nil {
		return errors.Wrap(err, "creating OIDC provider")
	}
	m.ProviderARN = *output.OpenIDConnectProviderArn
	return nil
}

// DeleteProvider will delete the provider using IAM API, it may return an error
// the API call fails
func (m *OpenIDConnectManager) DeleteProvider() error {
	// TODO: the ARN is deterministic, but we need to consider tracking
	// it somehow; it's possible to get a dangling resource if cluster
	// deletion was done by a version of eksctl that is not OIDC-aware,
	// as we don't use CloudFormation;
	// finding dangling resource will require looking at all clusters...
	input := &awsiam.DeleteOpenIDConnectProviderInput{
		OpenIDConnectProviderArn: &m.ProviderARN,
	}
	if _, err := m.iam.DeleteOpenIDConnectProvider(input); err != nil {
		return errors.Wrap(err, "deleting OIDC provider")
	}
	return nil
}

// getIssuerCAThumbprint obtains thumbprint of root CA by connecting to the
// OIDC issuer and parsing certificates
func (m *OpenIDConnectManager) getIssuerCAThumbprint() error {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: m.insecureSkipVerify,
			},
			Proxy: http.ProxyFromEnvironment,
		},
	}

	response, err := client.Get(m.issuerURL.String())
	if err != nil {
		return errors.Wrap(err, "connecting to issuer OIDC")
	}

	if response.TLS != nil {
		if numCerts := len(response.TLS.PeerCertificates); numCerts >= 1 {
			root := response.TLS.PeerCertificates[numCerts-1]
			m.issuerCAThumbprint = fmt.Sprintf("%x", sha1.Sum(root.Raw))
			return nil
		}
	}
	return fmt.Errorf("unable to get OIDC issuer's certificate")
}

// MakeAssumeRolePolicyDocument constructs a trust policy document for the given
// provider
func (m *OpenIDConnectManager) MakeAssumeRolePolicyDocument(serviceAccountNamespace, serviceAccountName string) cft.MapOfInterfaces {
	subject := fmt.Sprintf("system:serviceaccount:%s:%s", serviceAccountNamespace, serviceAccountName)
	return cft.MakeAssumeRoleWithWebIdentityPolicyDocument(m.ProviderARN, cft.MapOfInterfaces{
		"StringEquals": map[string]string{
			m.hostnameAndPath() + ":sub": subject,
			m.hostnameAndPath() + ":aud": m.audience,
		},
	})
}

func (m *OpenIDConnectManager) hostnameAndPath() string {
	return m.issuerURL.Hostname() + m.issuerURL.Path
}
