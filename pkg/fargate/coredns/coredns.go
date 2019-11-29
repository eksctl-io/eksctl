package coredns

import (
	"fmt"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	kubeclient "k8s.io/client-go/kubernetes"
)

const (
	// Namespace is the Kubernetes namespace under which CoreDNS lives.
	Namespace = "kube-system"
	// Name is the name of the Kubernetes Deployment object for CoreDNS.
	Name = "coredns"
	// ComputeTypeAnnotationKey is the key of the annotation driving CoreDNS'
	// scheduling.
	ComputeTypeAnnotationKey = "eks.amazonaws.com/compute-type"
	annotationValue          = "fargate"
	errorMsg                 = "failed to make CoreDNS schedulable on Fargate"
)

// IsSchedulableOnFargate analyzes the provided profiles to determine whether
// EKS' coredns deployment should be scheduled onto Fargate.
func IsSchedulableOnFargate(profiles []*api.FargateProfile) bool {
	for _, profile := range profiles {
		for _, selector := range profile.Selectors {
			if selectsCoreDNS(selector) {
				return true
			}
		}
	}
	return false
}

func selectsCoreDNS(selector api.FargateProfileSelector) bool {
	return selector.Namespace == Namespace && len(selector.Labels) == 0
}

// ScheduleOnFargate modifies EKS' coredns deployment so that it can be scheduled
// on Fargate.
func ScheduleOnFargate(clientSet kubeclient.Interface) error {
	deployments := clientSet.ExtensionsV1beta1().Deployments(Namespace)
	coreDNS, err := deployments.Get(Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	value, exists := coreDNS.Spec.Template.Annotations[ComputeTypeAnnotationKey]
	if !exists {
		return fmt.Errorf("%s: could not find annotation %q on CoreDNS", errorMsg, ComputeTypeAnnotationKey)
	}
	logger.Debug("CoreDNS is annotated with %s=%s", ComputeTypeAnnotationKey, value)
	coreDNS.Spec.Template.Annotations[ComputeTypeAnnotationKey] = annotationValue
	bytes, err := runtime.Encode(unstructured.UnstructuredJSONScheme, coreDNS)
	if err != nil {
		return errors.Wrapf(err, "%s: failed to marshal CoreDNS' updated Kubernetes deployment", errorMsg)
	}
	patched, err := deployments.Patch(Name, types.MergePatchType, bytes)
	if err != nil {
		return errors.Wrapf(err, errorMsg)
	}
	if _, exists = patched.Spec.Template.Annotations[ComputeTypeAnnotationKey]; !exists {
		return fmt.Errorf("%s: could not find annotation %q on patched CoreDNS, patching must have failed", errorMsg, ComputeTypeAnnotationKey)
	}
	logger.Info("CoreDNS is now schedulable onto Fargate")
	return nil
}
