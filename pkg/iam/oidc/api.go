package iamoidc

import (
	"context"
	"crypto/sha1"
	"crypto/tls"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/smithy-go"
	cft "github.com/weaveworks/eksctl/pkg/cfn/template"
	cf "github.com/weaveworks/goformation/v4/cloudformation/cloudformation"
	gfniam "github.com/weaveworks/goformation/v4/cloudformation/iam"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/weaveworks/eksctl/pkg/awsapi"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/pkg/errors"
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
	oidcThumbprint     string

	ProviderARN string

	iam awsapi.IAM
	cf  awsapi.CloudFormation
}

// UnsupportedOIDCError represents an unsupported OIDC error
type UnsupportedOIDCError struct {
	Message string
}

func (u *UnsupportedOIDCError) Error() string {
	return u.Message
}

// NewOpenIDConnectManager constructs a new IAM OIDC manager instance.
// It returns an error if the issuer URL is invalid
func NewOpenIDConnectManager(iamapi awsapi.IAM, cf awsapi.CloudFormation, accountID, issuer, partition string, tags map[string]string, oidcThumbprint string) (*OpenIDConnectManager, error) {
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
		iam:            iamapi,
		accountID:      accountID,
		partition:      partition,
		tags:           tags,
		audience:       defaultAudience,
		issuerURL:      issuerURL,
		cf:             cf,
		oidcThumbprint: oidcThumbprint,
	}

	if oidcThumbprint != "" {
		m.issuerCAThumbprint = oidcThumbprint
	}

	return m, nil
}

// CheckProviderExists will return true when the provider exists, it may return errors
// if it was unable to call IAM API
func (m *OpenIDConnectManager) CheckProviderExists(ctx context.Context) (bool, error) {
	if m.oidcThumbprint != "" {
		arn, err := m.getIDPArn()
		if err != nil {
			return false, err
		}

		if arn == "" {
			return false, nil
		}

		return true, nil
	}

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

func (m *OpenIDConnectManager) getStackName() *string {
	for k, v := range m.tags {
		if k == "alpha.eksctl.io/cluster-name" {
			return aws.String(fmt.Sprintf("eksctl-%s-oidc-provider", v))
		}
	}

	return aws.String(fmt.Sprintf("oidc-provider-unknown"))
}

func (m *OpenIDConnectManager) getIDPArn() (string, error) {
	_, err := m.cf.DescribeStacks(context.Background(), &cloudformation.DescribeStacksInput{
		StackName: m.getStackName(),
	})

	if err != nil {
		//var stackNotFound *types.StackNotFoundException
		//if errors.As(err, &stackNotFound) {
		if strings.Contains(err.Error(), "does not exist") {
			return "", nil
		}
		return "", err
	}

	// The stack exists, get the OIDC Provider ARN
	response, err := m.cf.DescribeStackResources(context.Background(), &cloudformation.DescribeStackResourcesInput{
		StackName: m.getStackName(),
	})

	if err != nil {
		return "", err
	}

	for _, resource := range response.StackResources {
		provider := gfniam.OIDCProvider{}
		if resource.ResourceType != nil && *resource.ResourceType == provider.AWSCloudFormationType() {
			fmt.Printf("OIDC Provider ARN: %s\n", *resource.PhysicalResourceId)
			m.ProviderARN = *resource.PhysicalResourceId
			break
		}
	}

	return m.ProviderARN, nil
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

	if m.oidcThumbprint != "" {
		var t []cf.Tag
		for k, v := range m.tags {
			t = append(t, cf.Tag{
				Key:   gfnt.NewString(k),
				Value: gfnt.NewString(v),
			})
		}

		hostWithoutPort := strings.Split(m.issuerURL.Host, ":")[0]
		newURL := *m.issuerURL
		newURL.Host = hostWithoutPort

		fmt.Println(m.issuerURL.String())
		fmt.Println(hostWithoutPort)
		fmt.Println(newURL.String())

		oidcProvider := &gfniam.OIDCProvider{
			ClientIdList:   gfnt.NewStringSlice(m.audience),
			Tags:           t,
			ThumbprintList: gfnt.NewStringSlice(m.oidcThumbprint),
			Url:            gfnt.NewString(newURL.String()),
		}

		body, err := oidcProvider.MarshalJSON()
		if err != nil {
			return err
		}

		template := fmt.Sprintf(`{
	"Resources": {
		"MyOIDCProvider": %s
	}
}`, string(string(body)))

		stack := &cloudformation.CreateStackInput{
			StackName:    m.getStackName(),
			TemplateBody: aws.String(string(template)),
		}

		output, err := m.cf.CreateStack(ctx, stack)
		if err != nil {
			return err
		}

		for {
			stack, err := m.cf.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
				StackName: m.getStackName(),
			})

			if err != nil {
				return err
			}

			if len(stack.Stacks) == 0 {
				continue
			}

			fmt.Println(stack.Stacks[0].StackStatus)

			if stack.Stacks[0].StackStatus == types.StackStatusCreateComplete {
				fmt.Println("Stack creation complete")
				break
			} else if stack.Stacks[0].StackStatus == types.StackStatusRollbackComplete || stack.Stacks[0].StackStatus == types.StackStatusRollbackFailed || stack.Stacks[0].StackStatus == types.StackStatusCreateFailed {
				fmt.Println("Stack creation failed")
				break
			}

			// Wait for 5 seconds before checking the status again
			time.Sleep(5 * time.Second)
			fmt.Println("waiting for stack to complete")
		}

		fmt.Println("created oidc stack")
		fmt.Println(output)

		_, err = m.getIDPArn()
		if err != nil {
			return err
		}

		return nil
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
	if m.oidcThumbprint != "" {
		_, err := m.cf.DeleteStack(ctx, &cloudformation.DeleteStackInput{
			StackName: m.getStackName(),
		}, nil)

		return err
	}
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
	if m.issuerCAThumbprint != "" {
		return nil
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: m.insecureSkipVerify,
				MinVersion:         tls.VersionTLS12,
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

// IsAccessDeniedError returns true if err is an AccessDenied error.
func IsAccessDeniedError(err error) bool {
	var apiErr smithy.APIError
	return errors.As(err, &apiErr) && apiErr.ErrorCode() == "AccessDenied"
}
