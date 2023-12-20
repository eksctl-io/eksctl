package v1alpha5

import (
	"encoding/json"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Commonly-used constants
const (
	AnnotationEKSRoleARN = "eks.amazonaws.com/role-arn"
	EKSServicePrincipal  = "pods.eks.amazonaws.com"
)

var EKSServicePrincipalTrustStatement = IAMStatement{
	Effect: "Allow",
	Action: []string{
		"sts:AssumeRole",
		"sts:TagSession",
	},
	Principal: map[string]CustomStringSlice{
		"Service": []string{EKSServicePrincipal},
	},
}

// ClusterIAM holds all IAM attributes of a cluster
type ClusterIAM struct {
	// +optional
	ServiceRoleARN *string `json:"serviceRoleARN,omitempty"`

	// permissions boundary for all identity-based entities created by eksctl.
	// See [AWS Permission Boundary](https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies_boundaries.html)
	// +optional
	ServiceRolePermissionsBoundary *string `json:"serviceRolePermissionsBoundary,omitempty"`

	// role used by pods to access AWS APIs. This role is added to the Kubernetes RBAC for authorization.
	// See [Pod Execution Role](https://docs.aws.amazon.com/eks/latest/userguide/pod-execution-role.html)
	// +optional
	FargatePodExecutionRoleARN *string `json:"fargatePodExecutionRoleARN,omitempty"`

	// permissions boundary for the fargate pod execution role`. See [EKS Fargate Support](/usage/fargate-support/)
	// +optional
	FargatePodExecutionRolePermissionsBoundary *string `json:"fargatePodExecutionRolePermissionsBoundary,omitempty"`

	// enables the IAM OIDC provider as well as IRSA for the Amazon CNI plugin
	// +optional
	WithOIDC *bool `json:"withOIDC,omitempty"`

	// service accounts to create in the cluster.
	// See [IAM Service Accounts](/usage/iamserviceaccounts/#usage-with-config-files)
	// +optional
	ServiceAccounts []*ClusterIAMServiceAccount `json:"serviceAccounts,omitempty"`

	// pod identity associations to create in the cluster.
	// See [Pod Identity Associations](TBD)
	// +optional
	PodIdentityAssociations []PodIdentityAssociation `json:"podIdentityAssociations,omitempty"`

	// VPCResourceControllerPolicy attaches the IAM policy
	// necessary to run the VPC controller in the control plane
	// Defaults to `true`
	VPCResourceControllerPolicy *bool `json:"vpcResourceControllerPolicy,omitempty"`
}

// ClusterIAMMeta holds information we can use to create ObjectMeta for service
// accounts
type ClusterIAMMeta struct {
	// +optional
	Name string `json:"name,omitempty"`

	// +optional
	Namespace string `json:"namespace,omitempty"`

	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
}

// AsObjectMeta gives us the k8s ObjectMeta needed to create the service account
func (iamMeta *ClusterIAMMeta) AsObjectMeta() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:        iamMeta.Name,
		Namespace:   iamMeta.Namespace,
		Annotations: iamMeta.Annotations,
		Labels:      iamMeta.Labels,
	}
}

// ClusterIAMServiceAccount holds an IAM service account metadata and configuration
type ClusterIAMServiceAccount struct {
	ClusterIAMMeta `json:"metadata,omitempty"`

	// list of ARNs of the IAM policies to attach
	// +optional
	AttachPolicyARNs []string `json:"attachPolicyARNs,omitempty"`

	WellKnownPolicies WellKnownPolicies `json:"wellKnownPolicies,omitempty"`

	// AttachPolicy holds a policy document to attach to this service account
	// +optional
	AttachPolicy InlineDocument `json:"attachPolicy,omitempty"`

	// ARN of the role to attach to the service account
	AttachRoleARN string `json:"attachRoleARN,omitempty"`

	// ARN of the permissions boundary to associate with the service account
	// +optional
	PermissionsBoundary string `json:"permissionsBoundary,omitempty"`

	// +optional
	Status *ClusterIAMServiceAccountStatus `json:"status,omitempty"`

	// Specific role name instead of the Cloudformation-generated role name
	// +optional
	RoleName string `json:"roleName,omitempty"`

	// Specify if only the IAM Service Account role should be created without creating/annotating the service account
	// +optional
	RoleOnly *bool `json:"roleOnly,omitempty"`

	// AWS tags for the service account
	// +optional
	Tags map[string]string `json:"tags,omitempty"`
}

// ClusterIAMServiceAccountStatus holds status of the IAM service account
type ClusterIAMServiceAccountStatus struct {
	// +optional
	RoleARN *string `json:"roleARN,omitempty"`
	// +optional
	StackName *string `json:"stackName,omitempty"`
	// +optional
	Tags map[string]string `json:"tags,omitempty"`
	// +optional
	Capabilities []string `json:"capabilities,omitempty"`
}

// NameString returns common name string
func (sa *ClusterIAMServiceAccount) NameString() string {
	return sa.Namespace + "/" + sa.Name
}

// ClusterIAMServiceAccountNameStringToClusterIAMMeta constructs metav1.ObjectMeta from <ns>/<name> string
func ClusterIAMServiceAccountNameStringToClusterIAMMeta(name string) (*ClusterIAMMeta, error) {
	nameParts := strings.Split(name, "/")
	if len(nameParts) != 2 {
		return nil, fmt.Errorf("unexpected serviceaccount name format %q", name)
	}
	meta := &ClusterIAMMeta{
		Namespace: nameParts[0],
		Name:      nameParts[1],
	}
	return meta, nil
}

// SetAnnotations sets eks.amazonaws.com/role-arn annotation according to IAM role used
func (sa *ClusterIAMServiceAccount) SetAnnotations() {
	if sa.Annotations == nil {
		sa.Annotations = make(map[string]string)
	}

	if sa.Status != nil && sa.Status.RoleARN != nil {
		sa.Annotations[AnnotationEKSRoleARN] = *sa.Status.RoleARN
	}
}

type PodIdentityAssociation struct {
	Namespace string `json:"namespace"`

	ServiceAccountName string `json:"serviceAccountName"`

	RoleARN string `json:"roleARN"`

	// +optional
	RoleName string `json:"roleName,omitempty"`

	// +optional
	PermissionsBoundaryARN string `json:"permissionsBoundaryARN,omitempty"`

	// +optional
	PermissionPolicyARNs []string `json:"permissionPolicyARNs,omitempty"`

	// +optional
	PermissionPolicy InlineDocument `json:"permissionPolicy,omitempty"`

	// +optional
	WellKnownPolicies WellKnownPolicies `json:"wellKnownPolicies,omitempty"`

	// +optional
	Tags map[string]string `json:"tags,omitempty"`
}

func (p PodIdentityAssociation) NameString() string {
	return p.Namespace + "/" + p.ServiceAccountName
}

// Internal type
// IAMPolicyDocument represents an IAM assume role policy document
type IAMPolicyDocument struct {
	Version    string         `json:"Version"`
	ID         string         `json:"Id,omitempty"`
	Statements []IAMStatement `json:"Statement"`
}

// Internal type
// IAMStatement represents an IAM policy document statement
type IAMStatement struct {
	Sid          string                       `json:"Sid,omitempty"`          // statement ID, service specific
	Effect       string                       `json:"Effect"`                 // Allow or Deny
	Principal    map[string]CustomStringSlice `json:"Principal,omitempty"`    // principal that is allowed or denied
	NotPrincipal map[string]CustomStringSlice `json:"NotPrincipal,omitempty"` // exception to a list of principals
	Action       CustomStringSlice            `json:"Action"`                 // allowed or denied action
	NotAction    CustomStringSlice            `json:"NotAction,omitempty"`    // matches everything except
	Resource     CustomStringSlice            `json:"Resource,omitempty"`     // object or objects that the statement covers
	NotResource  CustomStringSlice            `json:"NotResource,omitempty"`  // matches everything except
	Condition    json.RawMessage              `json:"Condition,omitempty"`    // conditions for when a policy is in effect
}

func (s *IAMStatement) ToMapOfInterfaces() map[string]interface{} {
	mapOfInterfaces := map[string]interface{}{
		"Effect": s.Effect,
		"Action": s.Action,
	}
	if s.Sid != "" {
		mapOfInterfaces["Sid"] = s.Sid
	}
	if s.Principal != nil {
		mapOfInterfaces["Principal"] = s.Principal
	}
	if s.NotPrincipal != nil {
		mapOfInterfaces["NotPrincipal"] = s.NotPrincipal
	}
	if s.NotAction != nil {
		mapOfInterfaces["NotAction"] = s.NotAction
	}
	if s.Resource != nil {
		mapOfInterfaces["Resource"] = s.Resource
	}
	if s.NotResource != nil {
		mapOfInterfaces["NotResource"] = s.NotResource
	}
	if s.Condition != nil {
		mapOfInterfaces["Condition"] = s.Condition
	}
	return mapOfInterfaces
}

// AWS allows string or []string as value, we convert everything to []string to avoid casting
type CustomStringSlice []string

func (c *CustomStringSlice) UnmarshalJSON(b []byte) error {

	var raw interface{}
	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}

	var p []string
	//  value can be string or []string, convert everything to []string
	switch v := raw.(type) {
	case string:
		p = []string{v}
	case []interface{}:
		var items []string
		for _, item := range v {
			items = append(items, fmt.Sprintf("%v", item))
		}
		p = items
	default:
		return fmt.Errorf("invalid %s value element: allowed is only string or []string", c)
	}

	*c = p
	return nil
}
