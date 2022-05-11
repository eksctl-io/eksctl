package iamoidc

import (
	"context"
	"crypto/sha1"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"

	"github.com/weaveworks/eksctl/pkg/awsapi"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/pkg/errors"

	cft "github.com/weaveworks/eksctl/pkg/cfn/template"
)

const defaultAudience = "sts.amazonaws.com"

// OpenIDConnectManager hold information about IAM OIDC integration
type OpenIDConnectManager struct {
	accountID string
	partition string
	audience  string
	tags      map[string]string

	issuerURL          *url.URL
	insecureSkipVerify bool
	issuerCAThumbprint string

	ProviderARN string

	iam awsapi.IAM
}

// NewOpenIDConnectManager constructs a new IAM OIDC manager instance.
// It returns an error if the issuer URL is invalid
func NewOpenIDConnectManager(iamapi awsapi.IAM, accountID, issuer, partition string, tags map[string]string) (*OpenIDConnectManager, error) {
	issuerURL, err := url.Parse(issuer)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing OIDC issuer URL")
	}

	if issuerURL.Scheme != "https" {
		return nil, fmt.Errorf("unsupported URL scheme %q", issuerURL.Scheme)
	}

	if issuerURL.Port() == "" {
		issuerURL.Host += ":443"
	}

	m := &OpenIDConnectManager{
		iam:       iamapi,
		accountID: accountID,
		partition: partition,
		tags:      tags,
		audience:  defaultAudience,
		issuerURL: issuerURL,
	}
	return m, nil
}

// CheckProviderExists will return true when the provider exists, it may return errors
// if it was unable to call IAM API
func (m *OpenIDConnectManager) CheckProviderExists(ctx context.Context) (bool, error) {
	input := &iam.GetOpenIDConnectProviderInput{
		OpenIDConnectProviderArn: aws.String(
			fmt.Sprintf("arn:%s:iam::%s:oidc-provider/%s", m.partition, m.accountID, m.hostnameAndPath()),
		),
	}
	_, err := m.iam.GetOpenIDConnectProvider(ctx, input)
	if err != nil {
		var oe *iamtypes.NoSuchEntityException
		if errors.As(err, &oe) {
			return false, nil
		}
		return false, err
	}
	m.ProviderARN = *input.OpenIDConnectProviderArn
	return true, nil
}

// CreateProvider will retrieve CA root certificate and compute its thumbprint for the
// by connecting to it and create the provider using IAM API
func (m *OpenIDConnectManager) CreateProvider(ctx context.Context) error {
	if err := m.getIssuerCAThumbprint(); err != nil {
		return err
	}

	var tags []iamtypes.Tag
	for k, v := range m.tags {
		tags = append(tags, iamtypes.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}

	input := &iam.CreateOpenIDConnectProviderInput{
		ClientIDList:   []string{m.audience},
		ThumbprintList: []string{m.issuerCAThumbprint},
		// It has no name or tags, it's keyed to the URL
		Url:  aws.String(m.issuerURL.String()),
		Tags: tags,
	}
	output, err := m.iam.CreateOpenIDConnectProvider(ctx, input)
	if err != nil {
		return errors.Wrap(err, "creating OIDC provider")
	}
	m.ProviderARN = *output.OpenIDConnectProviderArn
	return nil
}

// DeleteProvider will delete the provider using IAM API, it may return an error
// the API call fails
func (m *OpenIDConnectManager) DeleteProvider(ctx context.Context) error {
	// TODO: the ARN is deterministic, but we need to consider tracking
	// it somehow; it's possible to get a dangling resource if cluster
	// deletion was done by a version of eksctl that is not OIDC-aware,
	// as we don't use CloudFormation;
	// finding dangling resource will require looking at all clusters...
	input := &iam.DeleteOpenIDConnectProviderInput{
		OpenIDConnectProviderArn: &m.ProviderARN,
	}
	if _, err := m.iam.DeleteOpenIDConnectProvider(ctx, input); err != nil {
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
	defer response.Body.Close()
	if response.TLS != nil {
		if numCerts := len(response.TLS.PeerCertificates); numCerts >= 1 {
			root := response.TLS.PeerCertificates[numCerts-1]
			m.issuerCAThumbprint = fmt.Sprintf("%x", sha1.Sum(root.Raw))
			return nil
		}
	}
	return fmt.Errorf("unable to get OIDC issuer's certificate")
}

// MakeAssumeRolePolicyDocumentWithServiceAccountConditions constructs a trust policy document for the given
// provider
func (m *OpenIDConnectManager) MakeAssumeRolePolicyDocumentWithServiceAccountConditions(serviceAccountNamespace, serviceAccountName string) cft.MapOfInterfaces {
	subject := fmt.Sprintf("system:serviceaccount:%s:%s", serviceAccountNamespace, serviceAccountName)
	return cft.MakeAssumeRoleWithWebIdentityPolicyDocument(m.ProviderARN, cft.MapOfInterfaces{
		"StringEquals": map[string]string{
			m.hostnameAndPath() + ":sub": subject,
			m.hostnameAndPath() + ":aud": m.audience,
		},
	})
}

func (m *OpenIDConnectManager) MakeAssumeRolePolicyDocument() cft.MapOfInterfaces {
	return cft.MakeAssumeRoleWithWebIdentityPolicyDocument(m.ProviderARN, cft.MapOfInterfaces{
		"StringEquals": map[string]string{
			m.hostnameAndPath() + ":aud": m.audience,
		},
	})
}

func (m *OpenIDConnectManager) hostnameAndPath() string {
	return m.issuerURL.Hostname() + m.issuerURL.Path
}
