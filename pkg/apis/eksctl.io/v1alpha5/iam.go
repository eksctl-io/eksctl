package v1alpha5

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Commonly-used constants
const (
	AnnotationEKSRoleARN = "eks.amazonaws.com/role-arn"
)

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
