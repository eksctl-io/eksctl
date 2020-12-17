/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package evictor

import (
	"context"
	"fmt"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	daemonSetFatal      = "DaemonSet-managed Pods (use --ignore-daemonsets to ignore)"
	daemonSetWarning    = "ignoring DaemonSet-managed Pods"
	localStorageFatal   = "Pods with local storage (use --delete-local-data to override)"
	localStorageWarning = "deleting Pods with local storage"
	unmanagedFatal      = "Pods not managed by ReplicationController, ReplicaSet, Job, DaemonSet or StatefulSet (use --force to override)"
	unmanagedWarning    = "deleting Pods not managed by ReplicationController, ReplicaSet, Job, DaemonSet or StatefulSet"

	drainPodAnnotation       = "pod.alpha.kubernetes.io/drain"
	drainPodAnnotationForce  = "force"
	drainPodAnnotationIgnore = "ignore"
	drainPodAnnotationNever  = "never"
)

type PodDelete struct {
	Pod    corev1.Pod
	Status PodDeleteStatus
}

type PodDeleteList struct {
	Items []PodDelete
}

func (l *PodDeleteList) Pods() []corev1.Pod {
	pods := []corev1.Pod{}
	for _, i := range l.Items {
		if i.Status.Delete {
			pods = append(pods, i.Pod)
		}
	}
	return pods
}

func (l *PodDeleteList) Warnings() string {
	ps := make(map[string][]string)
	for _, i := range l.Items {
		if i.Status.Reason == podDeleteStatusTypeWarning {
			ps[i.Status.Message] = append(ps[i.Status.Message], fmt.Sprintf("%s/%s", i.Pod.Namespace, i.Pod.Name))
		}
	}

	msgs := []string{}
	for key, pods := range ps {
		msgs = append(msgs, fmt.Sprintf("%s: %s", key, strings.Join(pods, ", ")))
	}
	return strings.Join(msgs, "; ")
}

func (l *PodDeleteList) errors() []error {
	failedPods := make(map[string][]string)
	for _, i := range l.Items {
		if i.Status.Reason == podDeleteStatusTypeError {
			msg := i.Status.Message
			if msg == "" {
				msg = "unexpected error"
			}
			failedPods[msg] = append(failedPods[msg], fmt.Sprintf("%s/%s", i.Pod.Namespace, i.Pod.Name))
		}
	}
	errs := make([]error, 0)
	for msg, pods := range failedPods {
		errs = append(errs, fmt.Errorf("cannot Delete %s: %s", msg, strings.Join(pods, ", ")))
	}
	return errs
}

type PodDeleteStatus struct {
	Delete  bool
	Reason  string
	Message string
}

// Takes a Pod and returns a PodDeleteStatus
type podFilter func(corev1.Pod) PodDeleteStatus

const (
	podDeleteStatusTypeOkay    = "Okay"
	podDeleteStatusTypeSkip    = "Skip"
	podDeleteStatusTypeWarning = "Warning"
	podDeleteStatusTypeError   = "Error"
)

func makePodDeleteStatusOkay() PodDeleteStatus {
	return PodDeleteStatus{
		Delete: true,
		Reason: podDeleteStatusTypeOkay,
	}
}

func makePodDeleteStatusSkip() PodDeleteStatus {
	return PodDeleteStatus{
		Delete: false,
		Reason: podDeleteStatusTypeSkip,
	}
}

func makePodDeleteStatusWithWarning(delete bool, message string) PodDeleteStatus {
	return PodDeleteStatus{
		Delete:  delete,
		Reason:  podDeleteStatusTypeWarning,
		Message: message,
	}
}

func makePodDeleteStatusWithError(message string) PodDeleteStatus {
	return PodDeleteStatus{
		Delete:  false,
		Reason:  podDeleteStatusTypeError,
		Message: message,
	}
}

func (d *Evictor) makeFilters() []podFilter {
	return []podFilter{
		d.annotationFilter,
		d.daemonSetFilter,
		d.mirrorPodFilter,
		d.localStorageFilter,
		d.unreplicatedFilter,
	}
}

func hasLocalStorage(pod corev1.Pod) bool {
	for _, volume := range pod.Spec.Volumes {
		if volume.EmptyDir != nil && volume.EmptyDir.Medium != "Memory" {
			return true
		}
	}

	return false
}

func (d *Evictor) annotationFilter(pod corev1.Pod) PodDeleteStatus {
	if v, ok := pod.Annotations[drainPodAnnotation]; ok {
		annotation := fmt.Sprintf("due to annotation %s=%s", drainPodAnnotation, v)
		switch v {
		case drainPodAnnotationForce:
			return makePodDeleteStatusWithWarning(true, "forced "+annotation)
		case drainPodAnnotationIgnore:
			return makePodDeleteStatusWithWarning(false, "ignored "+annotation)
		case drainPodAnnotationNever:
			return makePodDeleteStatusWithError("cannot be drained " + annotation)
		}
	}

	return makePodDeleteStatusOkay()
}

func (d *Evictor) daemonSetFilter(pod corev1.Pod) PodDeleteStatus {
	// Note that we return false in cases where the Pod is DaemonSet managed,
	// regardless of flags.
	//
	// The exception is for pods that are orphaned (the referencing
	// management resource - including DaemonSet - is not found).
	// Such pods will be deleted if --force is used.
	controllerRef := metav1.GetControllerOf(&pod)
	if controllerRef == nil || controllerRef.Kind != appsv1.SchemeGroupVersion.WithKind("DaemonSet").Kind {
		return makePodDeleteStatusOkay()
	}
	// Any finished Pod can be removed.
	if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
		return makePodDeleteStatusOkay()
	}

	if _, err := d.client.AppsV1().DaemonSets(pod.Namespace).Get(context.TODO(), controllerRef.Name, metav1.GetOptions{}); err != nil {
		// remove orphaned pods with a warning if --force is used
		if apierrors.IsNotFound(err) && d.force {
			return makePodDeleteStatusWithWarning(true, err.Error())
		}

		return makePodDeleteStatusWithError(err.Error())
	}

	for _, ignoreDaemonSet := range d.ignoreDaemonSets {
		if controllerRef.Name == ignoreDaemonSet.Name {
			switch ignoreDaemonSet.Namespace {
			case pod.Namespace, metav1.NamespaceAll:
				return makePodDeleteStatusWithWarning(false, daemonSetWarning)
			}
		}
	}

	if !d.ignoreAllDaemonSets {
		return makePodDeleteStatusWithError(daemonSetFatal)
	}

	return makePodDeleteStatusWithWarning(false, daemonSetWarning)
}

func (d *Evictor) mirrorPodFilter(pod corev1.Pod) PodDeleteStatus {
	if _, found := pod.ObjectMeta.Annotations[corev1.MirrorPodAnnotationKey]; found {
		return makePodDeleteStatusSkip()
	}
	return makePodDeleteStatusOkay()
}

func (d *Evictor) localStorageFilter(pod corev1.Pod) PodDeleteStatus {
	if !hasLocalStorage(pod) {
		return makePodDeleteStatusOkay()
	}
	// Any finished Pod can be removed.
	if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
		return makePodDeleteStatusOkay()
	}
	if !d.deleteLocalData {
		return makePodDeleteStatusWithError(localStorageFatal)
	}

	return makePodDeleteStatusWithWarning(true, localStorageWarning)
}

func (d *Evictor) unreplicatedFilter(pod corev1.Pod) PodDeleteStatus {
	// any finished Pod can be removed
	if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
		return makePodDeleteStatusOkay()
	}

	controllerRef := metav1.GetControllerOf(&pod)
	if controllerRef != nil {
		return makePodDeleteStatusOkay()
	}
	if d.force {
		return makePodDeleteStatusWithWarning(true, unmanagedWarning)
	}
	return makePodDeleteStatusWithError(unmanagedFatal)
}
