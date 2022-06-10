package connector

import (
	"bytes"
	"context"
	"encoding/base64"
	"strings"
	"time"

	"github.com/aws/smithy-go"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	awsarn "github.com/aws/aws-sdk-go/aws/arn"

	"github.com/cenk/backoff"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"sigs.k8s.io/aws-iam-authenticator/pkg/arn"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

const (
	connectorPolicyName = "eks-connector-agent"
)

type ExternalCluster struct {
	Name             string
	Provider         string
	ConnectorRoleARN string
}

type EKSConnector struct {
	Provider         provider
	ManifestTemplate ManifestTemplate
}

type provider interface {
	EKS() awsapi.EKS
	STS() awsapi.STS
	STSPresigner() api.STSPresigner
	IAM() awsapi.IAM
	Region() string
}

type ManifestList struct {
	ConnectorResources     ManifestFile
	ClusterRoleResources   ManifestFile
	ConsoleAccessResources ManifestFile
	Expiry                 time.Time
	IAMIdentityARN         string
}

// RegisterCluster registers the specified external cluster with EKS and returns a list of Kubernetes resources
// for EKS Connector.
func (c *EKSConnector) RegisterCluster(ctx context.Context, cluster ExternalCluster) (*ManifestList, error) {
	cluster.Provider = strings.ToUpper(cluster.Provider)
	if err := validateProvider(cluster.Provider); err != nil {
		return nil, err
	}

	_, err := c.Provider.EKS().DescribeCluster(ctx, &eks.DescribeClusterInput{
		Name: aws.String(cluster.Name),
	})

	if err != nil {
		var notFoundError *ekstypes.ResourceNotFoundException
		if !errors.As(err, &notFoundError) {
			return nil, errors.New("unexpected error calling DescribeCluster")
		}
	} else {
		return nil, errors.Errorf("cluster already exists; deregister the cluster first using `eksctl deregister cluster --name %s --region %s` and try again", cluster.Name, c.Provider.Region())
	}

	connectorRoleARN := cluster.ConnectorRoleARN
	if connectorRoleARN == "" {
		var err error
		connectorRoleARN, err = c.createConnectorRole(ctx, cluster)
		if err != nil {
			return nil, errors.Wrap(err, "error creating IAM role for EKS Connector")
		}
	}

	registerOutput, err := c.registerCluster(ctx, cluster, connectorRoleARN)
	if err != nil {
		if cluster.ConnectorRoleARN == "" {
			if deleteErr := c.deleteRoleByARN(ctx, connectorRoleARN); deleteErr != nil {
				err = errors.Wrap(err, deleteErr.Error())
			}
		}
		return nil, errors.Wrap(err, "error calling RegisterCluster")
	}
	return c.createManifests(ctx, registerOutput.Cluster)
}

func (c *EKSConnector) createManifests(ctx context.Context, cluster *ekstypes.Cluster) (*ManifestList, error) {
	stsOutput, err := c.Provider.STS().GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}

	connectorResources := c.parseConnectorTemplate(cluster)
	roleARN, err := arn.Canonicalize(*stsOutput.Arn)
	if err != nil {
		return nil, errors.Wrap(err, "error canonicalizing IAM role ARN")
	}

	clusterRoleResources := c.applyRoleARN(c.ManifestTemplate.ClusterRole, roleARN)
	consoleAccessResources := c.applyRoleARN(c.ManifestTemplate.ConsoleAccess, roleARN)

	for _, m := range []ManifestFile{connectorResources, clusterRoleResources, consoleAccessResources} {
		if _, err := kubernetes.NewList(m.Data); err != nil {
			return nil, errors.Wrapf(err, "unexpected error parsing manifests for EKS Connector: %s", m.Filename)
		}
	}

	return &ManifestList{
		ConnectorResources:     connectorResources,
		ClusterRoleResources:   clusterRoleResources,
		ConsoleAccessResources: consoleAccessResources,
		Expiry:                 *cluster.ConnectorConfig.ActivationExpiry,
		IAMIdentityARN:         roleARN,
	}, nil
}

// ValidProviders returns a list of supported providers.
func ValidProviders() []ekstypes.ConnectorConfigProvider {
	var providerConfig ekstypes.ConnectorConfigProvider
	return providerConfig.Values()
}

func validateProvider(provider string) error {
	validProviders := ValidProviders()
	for _, p := range validProviders {
		if string(p) == provider {
			return nil
		}
	}
	return errors.Errorf("invalid provider %q; must be one of %s", provider, validProviders)
}

func (c *EKSConnector) registerCluster(ctx context.Context, cluster ExternalCluster, connectorRoleARN string) (*eks.RegisterClusterOutput, error) {
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 3 * time.Minute

	var registerOutput *eks.RegisterClusterOutput
	// IAM role takes some time to propagate,
	// RegisterCluster returns `InvalidRequestException: Not existing role` for such cases.
	err := backoff.RetryNotify(func() error {
		var err error

		registerOutput, err = c.Provider.EKS().RegisterCluster(ctx, &eks.RegisterClusterInput{
			Name: aws.String(cluster.Name),
			ConnectorConfig: &ekstypes.ConnectorConfigRequest{
				Provider: ekstypes.ConnectorConfigProvider(cluster.Provider),
				RoleArn:  aws.String(connectorRoleARN),
				// TODO add tags when they're supported by the API.
			},
		})

		if err != nil {
			var oe *smithy.OperationError
			if errors.As(err, &oe) && strings.Contains(oe.Error(), "Not existing role") {
				logger.Debug("IAM role could not be found; retrying RegisterCluster")
				return err
			}
			return backoff.Permanent(err)
		}

		return nil
	}, bo, func(err error, duration time.Duration) {
		logger.Debug("error calling RegisterCluster; retrying in %v: %v", duration, err)
	})

	if err != nil {
		var oe *smithy.OperationError
		if errors.As(err, &oe) && strings.Contains(oe.Error(), "AWSServiceRoleForAmazonEKSConnector is not available") {
			return nil, errors.Wrap(err, "SLR for EKS Connector does not exist; please run `aws iam create-service-linked-role --aws-service-name eks-connector.amazonaws.com` first")
		}
		return nil, err
	}

	return registerOutput, nil
}

func (c *EKSConnector) parseConnectorTemplate(cluster *ekstypes.Cluster) ManifestFile {
	activationCode := base64.StdEncoding.EncodeToString([]byte(*cluster.ConnectorConfig.ActivationCode))
	manifestFile := c.ManifestTemplate.Connector
	connectorResources := applyVariables(manifestFile.Data, "%EKS_ACTIVATION_ID%", *cluster.ConnectorConfig.ActivationId)
	connectorResources = applyVariables(connectorResources, "%EKS_ACTIVATION_CODE%", activationCode)
	connectorResources = applyVariables(connectorResources, "%AWS_REGION%", c.Provider.Region())
	return ManifestFile{
		Data:     connectorResources,
		Filename: manifestFile.Filename,
	}
}

func (c *EKSConnector) applyRoleARN(manifestFile ManifestFile, iamRoleARN string) ManifestFile {
	resources := applyVariables(manifestFile.Data, `%IAM_ARN%`, iamRoleARN)
	return ManifestFile{
		Data:     resources,
		Filename: manifestFile.Filename,
	}
}

func applyVariables(template []byte, field, value string) []byte {
	return bytes.ReplaceAll(template, []byte(field), []byte(value))
}

// DeregisterCluster deregisters the cluster and removes associated IAM resources.
func (c *EKSConnector) DeregisterCluster(ctx context.Context, clusterName string) error {
	clusterOutput, err := c.Provider.EKS().DeregisterCluster(ctx, &eks.DeregisterClusterInput{
		Name: aws.String(clusterName),
	})

	if err != nil {
		var notFoundErr *ekstypes.ResourceNotFoundException
		if errors.As(err, &notFoundErr) {
			return errors.Errorf("cluster %q does not exist", clusterName)
		}
		return errors.Wrap(err, "unexpected error deregistering cluster")
	}

	roleName, err := roleNameFromARN(*clusterOutput.Cluster.ConnectorConfig.RoleArn)
	if err != nil {
		return errors.Wrapf(err, "error parsing role ARN %q", *clusterOutput.Cluster.ConnectorConfig.RoleArn)
	}

	ownsIAMRole, err := c.ownsIAMRole(ctx, clusterName, roleName)
	if err != nil {
		return err
	}
	if !ownsIAMRole {
		return nil
	}

	return c.deleteRole(ctx, roleName)
}

func (c *EKSConnector) deleteRole(ctx context.Context, roleName string) error {
	logger.Info("deleting IAM role %q", roleName)

	if _, err := c.Provider.IAM().DeleteRolePolicy(ctx, &iam.DeleteRolePolicyInput{
		PolicyName: aws.String(connectorPolicyName),
		RoleName:   aws.String(roleName),
	}); err != nil {
		var notFoundErr *iamtypes.NoSuchEntityException
		if errors.As(err, &notFoundErr) {
			return errors.Errorf("could not find policy %q on IAM role", connectorPolicyName)
		}
		return err
	}

	if _, err := c.Provider.IAM().DeleteRole(ctx, &iam.DeleteRoleInput{
		RoleName: aws.String(roleName),
	}); err != nil {
		return errors.Wrap(err, "error deleting IAM role")
	}

	return nil
}

func (c *EKSConnector) deleteRoleByARN(ctx context.Context, roleARN string) error {
	connectorRoleName, err := roleNameFromARN(roleARN)
	if err != nil {
		return errors.Wrap(err, "error parsing connector role ARN")
	}
	return c.deleteRole(ctx, connectorRoleName)
}

func roleNameFromARN(roleARN string) (string, error) {
	parsed, err := awsarn.Parse(roleARN)
	if err != nil {
		return "", err
	}
	parts := strings.Split(parsed.Resource, "/")
	if len(parts) != 2 {
		return "", errors.New("invalid format for role ARN")
	}
	if parts[0] != "role" {
		return "", errors.Errorf(`expected resource type to be "role"; got %q`, parts[0])
	}
	return parts[1], nil
}

func (c *EKSConnector) ownsIAMRole(ctx context.Context, clusterName, roleName string) (bool, error) {
	roleOutput, err := c.Provider.IAM().GetRole(ctx, &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	})
	if err != nil {
		return false, errors.Wrapf(err, "error getting IAM role %q", roleName)
	}

	for _, tag := range roleOutput.Role.Tags {
		if *tag.Key == api.ClusterNameTag && *tag.Value == clusterName {
			return true, nil
		}
	}
	return false, nil
}

func (c *EKSConnector) createConnectorRole(ctx context.Context, cluster ExternalCluster) (string, error) {
	roleName := makeRoleName()
	logger.Info("creating IAM role %q", *roleName)

	output, err := c.Provider.IAM().CreateRole(ctx, &iam.CreateRoleInput{
		RoleName: roleName,
		AssumeRolePolicyDocument: aws.String(`{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "EKSConnectorAccess",
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "ssm.amazonaws.com"
        ]
      },
      "Action": "sts:AssumeRole"
    }
  ]
}`),
		Tags: []iamtypes.Tag{
			{
				Key:   aws.String(api.ClusterNameTag),
				Value: aws.String(cluster.Name),
			},
		},
	})
	if err != nil {
		return "", errors.Wrap(err, "error creating IAM role")
	}

	waiter := iam.NewRoleExistsWaiter(c.Provider.IAM())
	const maxWaitDuration = 5 * time.Minute
	if err := waiter.Wait(ctx, &iam.GetRoleInput{
		RoleName: roleName,
	}, maxWaitDuration); err != nil {
		return "", err
	}

	_, err = c.Provider.IAM().PutRolePolicy(ctx, &iam.PutRolePolicyInput{
		RoleName:   roleName,
		PolicyName: aws.String(connectorPolicyName),
		PolicyDocument: aws.String(`{
	  "Version": "2012-10-17",
	  "Statement": [
	    {
	      "Sid": "SsmControlChannel",
	      "Effect": "Allow",
	      "Action": [
	        "ssmmessages:CreateControlChannel"
	      ],
	      "Resource": "arn:aws:eks:*:*:cluster/*"
	    },
	    {
	      "Sid": "ssmDataplaneOperations",
	      "Effect": "Allow",
	      "Action": [
	        "ssmmessages:CreateDataChannel",
	        "ssmmessages:OpenDataChannel",
	        "ssmmessages:OpenControlChannel"
	      ],
	      "Resource": "*"
	    }
	  ]
	}`),
	})

	if err != nil {
		return "", err
	}

	return *output.Role.Arn, nil
}

func makeRoleName() *string {
	return aws.String(uniqueName("eksctl-"))
}

func uniqueName(prefix string) string {
	timestamp := strings.Replace(time.Now().Format("20060102150405.000000"), ".", "", 1)
	return prefix + timestamp
}
