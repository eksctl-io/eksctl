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
	// +optional
	ServiceRolePermissionsBoundary *string `json:"serviceRolePermissionsBoundary,omitempty"`
	// +optional
	FargatePodExecutionRoleARN *string `json:"fargatePodExecutionRoleARN,omitempty"`
	// +optional
	FargatePodExecutionRolePermissionsBoundary *string `json:"fargatePodExecutionRolePermissionsBoundary,omitempty"`
	// +optional
	WithOIDC *bool `json:"withOIDC,omitempty"`
	// +optional
	OIDCClientIDList []string `json:"oidcClientIDList,omitempty"`
	// +optional
	ServiceAccounts []*ClusterIAMServiceAccount `json:"serviceAccounts,omitempty"`
}

// ClusterIAMServiceAccount holds an iamserviceaccount metadata and configuration
type ClusterIAMServiceAccount struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// +optional
	AttachPolicyARNs []string `json:"attachPolicyARNs,omitempty"`
	// +optional
	AttachPolicy InlineDocument `json:"attachPolicy,omitempty"`
	// +optional
	PermissionsBoundary string `json:"permissionsBoundary,omitempty"`
	// +optional
	Status *ClusterIAMServiceAccountStatus `json:"status,omitempty"`
}

// ClusterIAMServiceAccountStatus holds status of iamserviceaccount
type ClusterIAMServiceAccountStatus struct {
	// +optional
	RoleARN *string `json:"roleARN,omitempty"`
}

// NameString returns common name string
func (sa *ClusterIAMServiceAccount) NameString() string {
	return sa.Namespace + "/" + sa.Name
}

// ClusterIAMServiceAccountNameStringToObjectMeta constructs metav1.ObjectMeta from <ns>/<name> string
func ClusterIAMServiceAccountNameStringToObjectMeta(name string) (*metav1.ObjectMeta, error) {
	nameParts := strings.Split(name, "/")
	if len(nameParts) != 2 {
		return nil, fmt.Errorf("unexpected serviceaccount name format %q", name)
	}
	meta := &metav1.ObjectMeta{
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
