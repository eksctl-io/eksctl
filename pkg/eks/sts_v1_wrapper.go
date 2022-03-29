package eks

import (
	"context"

	sts2 "github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/pkg/errors"

	"github.com/weaveworks/eksctl/pkg/awsapi"
)

// STSV1Wrapper wraps a sts v1 client using the implementation provided by v2.
type STSV1Wrapper struct {
	client awsapi.STS
}

func (s *STSV1Wrapper) AssumeRole(input *sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error) {
	ret, err := s.client.AssumeRole(context.TODO(), fromV1AssumeRoleInputToV2(input))
	if err != nil {
		return nil, err
	}
	return fromV2AssumeRoleOutputToV1(ret), nil
}

func (s *STSV1Wrapper) AssumeRoleWithContext(context aws.Context, input *sts.AssumeRoleInput, option ...request.Option) (*sts.AssumeRoleOutput, error) {
	ret, err := s.client.AssumeRole(context, fromV1AssumeRoleInputToV2(input))
	if err != nil {
		return nil, err
	}
	return fromV2AssumeRoleOutputToV1(ret), nil
}

func (s *STSV1Wrapper) AssumeRoleRequest(input *sts.AssumeRoleInput) (*request.Request, *sts.AssumeRoleOutput) {
	return &request.Request{}, &sts.AssumeRoleOutput{}
}

func fromV1AssumeRoleInputToV2(input *sts.AssumeRoleInput) *sts2.AssumeRoleInput {
	var policyARNs []types.PolicyDescriptorType
	for _, policy := range input.PolicyArns {
		policyARNs = append(policyARNs, types.PolicyDescriptorType{
			Arn: policy.Arn,
		})
	}
	var tags []types.Tag
	for _, tag := range input.Tags {
		tags = append(tags, types.Tag{
			Key:   tag.Key,
			Value: tag.Value,
		})
	}
	return &sts2.AssumeRoleInput{
		RoleArn:           input.RoleArn,
		RoleSessionName:   input.RoleSessionName,
		DurationSeconds:   aws.Int32(int32(aws.Int64Value(input.DurationSeconds))),
		ExternalId:        input.ExternalId,
		Policy:            input.Policy,
		PolicyArns:        policyARNs,
		SerialNumber:      input.SerialNumber,
		SourceIdentity:    input.SourceIdentity,
		Tags:              tags,
		TokenCode:         input.TokenCode,
		TransitiveTagKeys: aws.StringValueSlice(input.TransitiveTagKeys),
	}
}

func fromV2AssumeRoleOutputToV1(output *sts2.AssumeRoleOutput) *sts.AssumeRoleOutput {
	user := &sts.AssumedRoleUser{}
	if output.AssumedRoleUser != nil {
		user.AssumedRoleId = output.AssumedRoleUser.AssumedRoleId
		user.Arn = output.AssumedRoleUser.Arn
	}
	credentials := &sts.Credentials{}
	if output.Credentials != nil {
		credentials.SecretAccessKey = output.Credentials.SecretAccessKey
		credentials.SessionToken = output.Credentials.SessionToken
		credentials.AccessKeyId = output.Credentials.AccessKeyId
		credentials.Expiration = output.Credentials.Expiration
	}
	return &sts.AssumeRoleOutput{
		AssumedRoleUser:  user,
		Credentials:      credentials,
		PackedPolicySize: aws.Int64(int64(aws.Int32Value(output.PackedPolicySize))),
		SourceIdentity:   output.SourceIdentity,
	}
}

func (s *STSV1Wrapper) AssumeRoleWithSAML(input *sts.AssumeRoleWithSAMLInput) (*sts.AssumeRoleWithSAMLOutput, error) {
	ret, err := s.client.AssumeRoleWithSAML(context.TODO(), fromV1AssumeRoleWithSAMLInputToV2(input))
	if err != nil {
		return nil, err
	}
	return fromV2AssumeRoleWithSAMLOutputToV1(ret), nil
}

func (s *STSV1Wrapper) AssumeRoleWithSAMLWithContext(ctx aws.Context, input *sts.AssumeRoleWithSAMLInput, option ...request.Option) (*sts.AssumeRoleWithSAMLOutput, error) {
	ret, err := s.client.AssumeRoleWithSAML(ctx, fromV1AssumeRoleWithSAMLInputToV2(input))
	if err != nil {
		return nil, err
	}
	return fromV2AssumeRoleWithSAMLOutputToV1(ret), nil
}

func (s *STSV1Wrapper) AssumeRoleWithSAMLRequest(input *sts.AssumeRoleWithSAMLInput) (*request.Request, *sts.AssumeRoleWithSAMLOutput) {
	return &request.Request{}, &sts.AssumeRoleWithSAMLOutput{}
}

func fromV1AssumeRoleWithSAMLInputToV2(input *sts.AssumeRoleWithSAMLInput) *sts2.AssumeRoleWithSAMLInput {
	var policyARNs []types.PolicyDescriptorType
	for _, policy := range input.PolicyArns {
		policyARNs = append(policyARNs, types.PolicyDescriptorType{
			Arn: policy.Arn,
		})
	}
	return &sts2.AssumeRoleWithSAMLInput{
		RoleArn:         input.RoleArn,
		DurationSeconds: aws.Int32(int32(aws.Int64Value(input.DurationSeconds))),
		Policy:          input.Policy,
		PolicyArns:      policyARNs,
		SAMLAssertion:   input.SAMLAssertion,
		PrincipalArn:    input.PrincipalArn,
	}
}

func fromV2AssumeRoleWithSAMLOutputToV1(output *sts2.AssumeRoleWithSAMLOutput) *sts.AssumeRoleWithSAMLOutput {
	user := &sts.AssumedRoleUser{}
	if output.AssumedRoleUser != nil {
		user.AssumedRoleId = output.AssumedRoleUser.AssumedRoleId
		user.Arn = output.AssumedRoleUser.Arn
	}
	credentials := &sts.Credentials{}
	if output.Credentials != nil {
		credentials.SecretAccessKey = output.Credentials.SecretAccessKey
		credentials.SessionToken = output.Credentials.SessionToken
		credentials.AccessKeyId = output.Credentials.AccessKeyId
		credentials.Expiration = output.Credentials.Expiration
	}
	return &sts.AssumeRoleWithSAMLOutput{
		AssumedRoleUser:  user,
		Audience:         nil,
		Credentials:      credentials,
		Issuer:           output.Issuer,
		NameQualifier:    output.NameQualifier,
		PackedPolicySize: aws.Int64(int64(aws.Int32Value(output.PackedPolicySize))),
		SourceIdentity:   output.SourceIdentity,
		Subject:          output.Subject,
		SubjectType:      output.SubjectType,
	}
}

func (s *STSV1Wrapper) AssumeRoleWithWebIdentity(input *sts.AssumeRoleWithWebIdentityInput) (*sts.AssumeRoleWithWebIdentityOutput, error) {
	ret, err := s.client.AssumeRoleWithWebIdentity(context.TODO(), fromV1AssumeRoleWithWebIdentityInputToV2(input))
	if err != nil {
		return nil, err
	}
	return fromV2AssumeRoleWithWebIdentityOutputToV1(ret), nil
}

func (s *STSV1Wrapper) AssumeRoleWithWebIdentityWithContext(ctx aws.Context, input *sts.AssumeRoleWithWebIdentityInput, option ...request.Option) (*sts.AssumeRoleWithWebIdentityOutput, error) {
	ret, err := s.client.AssumeRoleWithWebIdentity(ctx, fromV1AssumeRoleWithWebIdentityInputToV2(input))
	if err != nil {
		return nil, err
	}
	return fromV2AssumeRoleWithWebIdentityOutputToV1(ret), nil
}

func (s *STSV1Wrapper) AssumeRoleWithWebIdentityRequest(input *sts.AssumeRoleWithWebIdentityInput) (*request.Request, *sts.AssumeRoleWithWebIdentityOutput) {
	return &request.Request{}, &sts.AssumeRoleWithWebIdentityOutput{}
}

func fromV1AssumeRoleWithWebIdentityInputToV2(input *sts.AssumeRoleWithWebIdentityInput) *sts2.AssumeRoleWithWebIdentityInput {
	var policyARNs []types.PolicyDescriptorType
	for _, policy := range input.PolicyArns {
		policyARNs = append(policyARNs, types.PolicyDescriptorType{
			Arn: policy.Arn,
		})
	}
	return &sts2.AssumeRoleWithWebIdentityInput{
		RoleArn:          input.RoleArn,
		RoleSessionName:  input.RoleSessionName,
		WebIdentityToken: input.WebIdentityToken,
		DurationSeconds:  aws.Int32(int32(aws.Int64Value(input.DurationSeconds))),
		Policy:           input.Policy,
		PolicyArns:       policyARNs,
		ProviderId:       input.ProviderId,
	}
}

func fromV2AssumeRoleWithWebIdentityOutputToV1(output *sts2.AssumeRoleWithWebIdentityOutput) *sts.AssumeRoleWithWebIdentityOutput {
	user := &sts.AssumedRoleUser{}
	if output.AssumedRoleUser != nil {
		user.AssumedRoleId = output.AssumedRoleUser.AssumedRoleId
		user.Arn = output.AssumedRoleUser.Arn
	}
	credentials := &sts.Credentials{}
	if output.Credentials != nil {
		credentials.SecretAccessKey = output.Credentials.SecretAccessKey
		credentials.SessionToken = output.Credentials.SessionToken
		credentials.AccessKeyId = output.Credentials.AccessKeyId
		credentials.Expiration = output.Credentials.Expiration
	}
	return &sts.AssumeRoleWithWebIdentityOutput{
		AssumedRoleUser:             user,
		Audience:                    output.Audience,
		Credentials:                 credentials,
		PackedPolicySize:            aws.Int64(int64(aws.Int32Value(output.PackedPolicySize))),
		Provider:                    output.Provider,
		SourceIdentity:              output.SourceIdentity,
		SubjectFromWebIdentityToken: output.SubjectFromWebIdentityToken,
	}
}

func (s *STSV1Wrapper) DecodeAuthorizationMessage(input *sts.DecodeAuthorizationMessageInput) (*sts.DecodeAuthorizationMessageOutput, error) {
	ret, err := s.client.DecodeAuthorizationMessage(context.TODO(), &sts2.DecodeAuthorizationMessageInput{
		EncodedMessage: input.EncodedMessage,
	})
	if err != nil {
		return nil, err
	}
	return &sts.DecodeAuthorizationMessageOutput{
		DecodedMessage: ret.DecodedMessage,
	}, nil
}

func (s *STSV1Wrapper) DecodeAuthorizationMessageWithContext(ctx aws.Context, input *sts.DecodeAuthorizationMessageInput, option ...request.Option) (*sts.DecodeAuthorizationMessageOutput, error) {
	ret, err := s.client.DecodeAuthorizationMessage(ctx, &sts2.DecodeAuthorizationMessageInput{
		EncodedMessage: input.EncodedMessage,
	})
	if err != nil {
		return nil, err
	}
	return &sts.DecodeAuthorizationMessageOutput{
		DecodedMessage: ret.DecodedMessage,
	}, nil
}

func (s *STSV1Wrapper) DecodeAuthorizationMessageRequest(input *sts.DecodeAuthorizationMessageInput) (*request.Request, *sts.DecodeAuthorizationMessageOutput) {
	return &request.Request{}, &sts.DecodeAuthorizationMessageOutput{}
}

func (s *STSV1Wrapper) GetAccessKeyInfo(input *sts.GetAccessKeyInfoInput) (*sts.GetAccessKeyInfoOutput, error) {
	ret, err := s.client.GetAccessKeyInfo(context.TODO(), &sts2.GetAccessKeyInfoInput{
		AccessKeyId: input.AccessKeyId,
	})
	if err != nil {
		return nil, err
	}
	return &sts.GetAccessKeyInfoOutput{
		Account: ret.Account,
	}, nil
}

func (s *STSV1Wrapper) GetAccessKeyInfoWithContext(ctx aws.Context, input *sts.GetAccessKeyInfoInput, option ...request.Option) (*sts.GetAccessKeyInfoOutput, error) {
	ret, err := s.client.GetAccessKeyInfo(ctx, &sts2.GetAccessKeyInfoInput{
		AccessKeyId: input.AccessKeyId,
	})
	if err != nil {
		return nil, err
	}
	return &sts.GetAccessKeyInfoOutput{
		Account: ret.Account,
	}, nil
}

func (s *STSV1Wrapper) GetAccessKeyInfoRequest(input *sts.GetAccessKeyInfoInput) (*request.Request, *sts.GetAccessKeyInfoOutput) {
	return &request.Request{}, &sts.GetAccessKeyInfoOutput{}
}

func (s *STSV1Wrapper) GetCallerIdentity(input *sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
	ret, err := s.client.GetCallerIdentity(context.TODO(), &sts2.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}
	return &sts.GetCallerIdentityOutput{
		Account: ret.Account,
		Arn:     ret.Arn,
		UserId:  ret.UserId,
	}, nil
}

func (s *STSV1Wrapper) GetCallerIdentityWithContext(ctx aws.Context, input *sts.GetCallerIdentityInput, option ...request.Option) (*sts.GetCallerIdentityOutput, error) {
	ret, err := s.client.GetCallerIdentity(ctx, &sts2.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}
	return &sts.GetCallerIdentityOutput{
		Account: ret.Account,
		Arn:     ret.Arn,
		UserId:  ret.UserId,
	}, nil
}

func (s *STSV1Wrapper) GetCallerIdentityRequest(input *sts.GetCallerIdentityInput) (*request.Request, *sts.GetCallerIdentityOutput) {
	return &request.Request{}, &sts.GetCallerIdentityOutput{}
}

func (s *STSV1Wrapper) GetFederationToken(input *sts.GetFederationTokenInput) (*sts.GetFederationTokenOutput, error) {
	ret, err := s.client.GetFederationToken(context.TODO(), fromV1GetFederationTokenInputToV2(input))
	if err != nil {
		return nil, err
	}
	return fromV2GetFederationTokenOutputToV1(ret), nil
}

func (s *STSV1Wrapper) GetFederationTokenWithContext(ctx aws.Context, input *sts.GetFederationTokenInput, option ...request.Option) (*sts.GetFederationTokenOutput, error) {
	ret, err := s.client.GetFederationToken(ctx, fromV1GetFederationTokenInputToV2(input))
	if err != nil {
		return nil, err
	}
	return fromV2GetFederationTokenOutputToV1(ret), nil
}

func (s *STSV1Wrapper) GetFederationTokenRequest(input *sts.GetFederationTokenInput) (*request.Request, *sts.GetFederationTokenOutput) {
	return &request.Request{}, &sts.GetFederationTokenOutput{}
}

func fromV1GetFederationTokenInputToV2(input *sts.GetFederationTokenInput) *sts2.GetFederationTokenInput {
	var policyARNs []types.PolicyDescriptorType
	for _, policy := range input.PolicyArns {
		policyARNs = append(policyARNs, types.PolicyDescriptorType{
			Arn: policy.Arn,
		})
	}
	var tags []types.Tag
	for _, tag := range input.Tags {
		tags = append(tags, types.Tag{
			Key:   tag.Key,
			Value: tag.Value,
		})
	}
	return &sts2.GetFederationTokenInput{
		Name:            input.Name,
		DurationSeconds: aws.Int32(int32(aws.Int64Value(input.DurationSeconds))),
		Policy:          input.Policy,
		PolicyArns:      policyARNs,
		Tags:            tags,
	}
}

func fromV2GetFederationTokenOutputToV1(output *sts2.GetFederationTokenOutput) *sts.GetFederationTokenOutput {
	credentials := &sts.Credentials{}
	if output.Credentials != nil {
		credentials.SecretAccessKey = output.Credentials.SecretAccessKey
		credentials.SessionToken = output.Credentials.SessionToken
		credentials.AccessKeyId = output.Credentials.AccessKeyId
		credentials.Expiration = output.Credentials.Expiration
	}
	return &sts.GetFederationTokenOutput{
		Credentials:      credentials,
		PackedPolicySize: aws.Int64(int64(aws.Int32Value(output.PackedPolicySize))),
	}
}

func (s *STSV1Wrapper) GetSessionToken(input *sts.GetSessionTokenInput) (*sts.GetSessionTokenOutput, error) {
	ret, err := s.client.GetSessionToken(context.TODO(), &sts2.GetSessionTokenInput{
		DurationSeconds: aws.Int32(int32(aws.Int64Value(input.DurationSeconds))),
		SerialNumber:    input.SerialNumber,
		TokenCode:       input.TokenCode,
	})
	if err != nil {
		return nil, err
	}
	if ret.Credentials == nil {
		return nil, errors.New("no credentials in response")
	}
	return &sts.GetSessionTokenOutput{
		Credentials: &sts.Credentials{
			AccessKeyId:     ret.Credentials.AccessKeyId,
			Expiration:      ret.Credentials.Expiration,
			SecretAccessKey: ret.Credentials.SecretAccessKey,
			SessionToken:    ret.Credentials.SessionToken,
		},
	}, nil
}

func (s *STSV1Wrapper) GetSessionTokenWithContext(ctx aws.Context, input *sts.GetSessionTokenInput, option ...request.Option) (*sts.GetSessionTokenOutput, error) {
	ret, err := s.client.GetSessionToken(ctx, &sts2.GetSessionTokenInput{
		DurationSeconds: aws.Int32(int32(aws.Int64Value(input.DurationSeconds))),
		SerialNumber:    input.SerialNumber,
		TokenCode:       input.TokenCode,
	})
	if err != nil {
		return nil, err
	}
	if ret.Credentials == nil {
		return nil, errors.New("no credentials in response")
	}
	return &sts.GetSessionTokenOutput{
		Credentials: &sts.Credentials{
			AccessKeyId:     ret.Credentials.AccessKeyId,
			Expiration:      ret.Credentials.Expiration,
			SecretAccessKey: ret.Credentials.SecretAccessKey,
			SessionToken:    ret.Credentials.SessionToken,
		},
	}, nil
}

func (s *STSV1Wrapper) GetSessionTokenRequest(input *sts.GetSessionTokenInput) (*request.Request, *sts.GetSessionTokenOutput) {
	return &request.Request{}, &sts.GetSessionTokenOutput{}
}

func NewSTSV1Wrapper(client awsapi.STS) *STSV1Wrapper {
	return &STSV1Wrapper{
		client: client,
	}
}

var _ stsiface.STSAPI = &STSV1Wrapper{}
