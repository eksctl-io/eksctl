package coredns

import (
	"fmt"
	"strings"
	"time"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/utils/retry"
	v1 "k8s.io/api/core/v1"
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
	computeTypeFargate       = "fargate"
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

// IsScheduledOnFargate checks if EKS' coredns is scheduled onto Fargate.
func IsScheduledOnFargate(clientSet kubeclient.Interface) (bool, error) {
	isDepOnFargate, err := isDeploymentScheduledOnFargate(clientSet)
	if err != nil {
		return false, err
	}
	arePodsOnFargate, err := arePodsScheduledOnFargate(clientSet)
	if err != nil {
		return false, err
	}
	return isDepOnFargate && arePodsOnFargate, nil
}

func isDeploymentScheduledOnFargate(clientSet kubeclient.Interface) (bool, error) {
	coredns, err := clientSet.AppsV1().Deployments(Namespace).Get(Name, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	if coredns.Spec.Replicas == nil {
		return false, errors.New("nil spec.replicas in coredns deployment")
	}
	computeType, exists := coredns.Spec.Template.Annotations[ComputeTypeAnnotationKey]
	logger.Debug("deployment %q with compute type %q currently has %v/%v replicas running", Name, computeType, coredns.Status.ReadyReplicas, *coredns.Spec.Replicas)
	scheduled := exists &&
		computeType == computeTypeFargate &&
		*coredns.Spec.Replicas == coredns.Status.ReadyReplicas
	if scheduled {
		logger.Info("%q is now scheduled onto Fargate", Name)
	}
	return scheduled, nil
}

func arePodsScheduledOnFargate(clientSet kubeclient.Interface) (bool, error) {
	pods, err := clientSet.CoreV1().Pods(Namespace).List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("eks.amazonaws.com/component = %s", Name),
	})
	if err != nil {
		return false, err
	}
	for _, pod := range pods.Items {
		if !isRunningOnFargate(&pod) {
			return false, nil
		}
	}
	logger.Info("%q pods are now scheduled onto Fargate", Name)
	return true, nil
}

func isRunningOnFargate(pod *v1.Pod) bool {
	computeType, exists := pod.Annotations[ComputeTypeAnnotationKey]
	logger.Debug("pod %q with compute type %q and status %q is scheduled on %q", pod.Name, computeType, pod.Status.Phase, pod.Spec.NodeName)
	return exists &&
		computeType == computeTypeFargate &&
		pod.Status.Phase == v1.PodRunning &&
		strings.HasPrefix(pod.Spec.NodeName, "fargate-ip-")
}

// ScheduleOnFargate modifies EKS' coredns deployment so that it can be scheduled
// on Fargate.
func ScheduleOnFargate(clientSet kubeclient.Interface) error {
	if err := scheduleOnFargate(clientSet); err != nil {
		return errors.Wrapf(err, "failed to make %q deployment schedulable on Fargate", Name)
	}
	logger.Info("%q is now schedulable onto Fargate", Name)
	return nil
}

func scheduleOnFargate(clientSet kubeclient.Interface) error {
	deployments := clientSet.AppsV1().Deployments(Namespace)
	coredns, err := deployments.Get(Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	coredns.Spec.Template.Annotations[ComputeTypeAnnotationKey] = computeTypeFargate
	bytes, err := runtime.Encode(unstructured.UnstructuredJSONScheme, coredns)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal %q deployment", Name)
	}
	patched, err := deployments.Patch(Name, types.MergePatchType, bytes)
	if err != nil {
		return errors.Wrap(err, "failed to patch deployment")
	}
	value, exists := patched.Spec.Template.Annotations[ComputeTypeAnnotationKey]
	if !exists {
		return fmt.Errorf("could not find annotation %q on patched deployment %q: patching must have failed", ComputeTypeAnnotationKey, Name)
	}
	if value != computeTypeFargate {
		return fmt.Errorf("unexpected value %q for annotation %q on %q patched deployment", value, ComputeTypeAnnotationKey, Name)
	}
	return nil
}

// WaitForScheduleOnFargate waits for coredns to be scheduled on Fargate.
// It will wait until it has detected that the scheduling has been successful,
// or until the retry policy times out, whichever happens first.
func WaitForScheduleOnFargate(clientSet kubeclient.Interface, retryPolicy retry.Policy) error {
	// Clone the retry policy to ensure this method is re-entrant/thread-safe:
	retryPolicy = retryPolicy.Clone()
	for !retryPolicy.Done() {
		isScheduled, err := IsScheduledOnFargate(clientSet)
		if err != nil {
			return err
		}
		if isScheduled {
			return nil
		}
		time.Sleep(retryPolicy.Duration())
	}
	return fmt.Errorf("timed out while waiting for %q to be scheduled on Fargate", Name)
}
