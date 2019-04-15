package cloudformation

import (
	"encoding/json"
)

// AWSElasticLoadBalancingV2Listener_AuthenticateOidcConfig AWS CloudFormation Resource (AWS::ElasticLoadBalancingV2::Listener.AuthenticateOidcConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-listener-authenticateoidcconfig.html
type AWSElasticLoadBalancingV2Listener_AuthenticateOidcConfig struct {

	// AuthenticationRequestExtraParams AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-listener-authenticateoidcconfig.html#cfn-elasticloadbalancingv2-listener-authenticateoidcconfig-authenticationrequestextraparams
	AuthenticationRequestExtraParams map[string]*Value `json:"AuthenticationRequestExtraParams,omitempty"`

	// AuthorizationEndpoint AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-listener-authenticateoidcconfig.html#cfn-elasticloadbalancingv2-listener-authenticateoidcconfig-authorizationendpoint
	AuthorizationEndpoint *Value `json:"AuthorizationEndpoint,omitempty"`

	// ClientId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-listener-authenticateoidcconfig.html#cfn-elasticloadbalancingv2-listener-authenticateoidcconfig-clientid
	ClientId *Value `json:"ClientId,omitempty"`

	// ClientSecret AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-listener-authenticateoidcconfig.html#cfn-elasticloadbalancingv2-listener-authenticateoidcconfig-clientsecret
	ClientSecret *Value `json:"ClientSecret,omitempty"`

	// Issuer AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-listener-authenticateoidcconfig.html#cfn-elasticloadbalancingv2-listener-authenticateoidcconfig-issuer
	Issuer *Value `json:"Issuer,omitempty"`

	// OnUnauthenticatedRequest AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-listener-authenticateoidcconfig.html#cfn-elasticloadbalancingv2-listener-authenticateoidcconfig-onunauthenticatedrequest
	OnUnauthenticatedRequest *Value `json:"OnUnauthenticatedRequest,omitempty"`

	// Scope AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-listener-authenticateoidcconfig.html#cfn-elasticloadbalancingv2-listener-authenticateoidcconfig-scope
	Scope *Value `json:"Scope,omitempty"`

	// SessionCookieName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-listener-authenticateoidcconfig.html#cfn-elasticloadbalancingv2-listener-authenticateoidcconfig-sessioncookiename
	SessionCookieName *Value `json:"SessionCookieName,omitempty"`

	// SessionTimeout AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-listener-authenticateoidcconfig.html#cfn-elasticloadbalancingv2-listener-authenticateoidcconfig-sessiontimeout
	SessionTimeout *Value `json:"SessionTimeout,omitempty"`

	// TokenEndpoint AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-listener-authenticateoidcconfig.html#cfn-elasticloadbalancingv2-listener-authenticateoidcconfig-tokenendpoint
	TokenEndpoint *Value `json:"TokenEndpoint,omitempty"`

	// UserInfoEndpoint AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-listener-authenticateoidcconfig.html#cfn-elasticloadbalancingv2-listener-authenticateoidcconfig-userinfoendpoint
	UserInfoEndpoint *Value `json:"UserInfoEndpoint,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSElasticLoadBalancingV2Listener_AuthenticateOidcConfig) AWSCloudFormationType() string {
	return "AWS::ElasticLoadBalancingV2::Listener.AuthenticateOidcConfig"
}

func (r *AWSElasticLoadBalancingV2Listener_AuthenticateOidcConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
