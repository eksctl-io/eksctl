package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sts"

	smithyhttp "github.com/aws/smithy-go/transport/http"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/credentials"
)

// TokenGenerator defines a token generator using STS.
type TokenGenerator interface {
	GetWithSTS(ctx context.Context, clusterID string) (Token, error)
}

// Generator provides a token generating functionality using a signed STS CallerIdentity call.
type Generator struct {
	client api.STSPresigner
	clock  credentials.Clock
}

func NewGenerator(client api.STSPresigner, clock credentials.Clock) Generator {
	return Generator{
		client: client,
		clock:  clock,
	}
}

// Token is generated and used by Kubernetes client-go to authenticate with a Kubernetes cluster.
type Token struct {
	Token      string
	Expiration time.Time
}

const (
	clusterIDHeader        = "x-k8s-aws-id"
	presignedURLExpiration = 10 * time.Minute
	v1Prefix               = "k8s-aws-v1."
)

// GetWithSTS returns a token valid for clusterID using the given STS client.
// This implementation follows the steps outlined here:
// https://github.com/kubernetes-sigs/aws-iam-authenticator#api-authorization-from-outside-a-cluster
// We either add this implementation or have to maintain two versions of STS since aws-iam-authenticator is
// not switching over to aws-go-sdk-v2.
func (g Generator) GetWithSTS(ctx context.Context, clusterID string) (Token, error) {
	// generate a sts:GetCallerIdentity request and add our custom cluster ID header
	presignedURLRequest, err := g.client.PresignGetCallerIdentity(ctx, &sts.GetCallerIdentityInput{}, func(presignOptions *sts.PresignOptions) {
		presignOptions.ClientOptions = append(presignOptions.ClientOptions, g.appendPresignHeaderValuesFunc(clusterID))
	})
	if err != nil {
		return Token{}, fmt.Errorf("failed to presign caller identity: %w", err)
	}

	tokenExpiration := g.clock.Now().Local().Add(presignedURLExpiration)
	// Add the token with k8s-aws-v1. prefix.
	return Token{v1Prefix + base64.RawURLEncoding.EncodeToString([]byte(presignedURLRequest.URL)), tokenExpiration}, nil
}

func (g Generator) appendPresignHeaderValuesFunc(clusterID string) func(stsOptions *sts.Options) {
	return func(stsOptions *sts.Options) {
		stsOptions.APIOptions = append(stsOptions.APIOptions,
			// Add clusterId Header.
			smithyhttp.SetHeaderValue(clusterIDHeader, clusterID),
			// Add X-Amz-Expires query param.
			smithyhttp.SetHeaderValue("X-Amz-Expires", "60"))
	}
}
