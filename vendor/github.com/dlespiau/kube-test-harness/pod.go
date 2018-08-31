package harness

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	appsv1 "k8s.io/api/apps/v1beta2"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
)

func (test *Test) listPods(namespace string, options metav1.ListOptions) (*v1.PodList, error) {
	return test.harness.kubeClient.Core().Pods(namespace).List(options)
}

// ListPods returns the list of pods in namespace matching options.
func (test *Test) ListPods(namespace string, options metav1.ListOptions) *v1.PodList {
	pl, err := test.listPods(namespace, options)
	test.err(err)
	return pl
}

func (test *Test) listPodsFromDeployment(d *appsv1.Deployment) (*v1.PodList, error) {
	// XXX: there must a better way to do this?
	selector, err := selectorToString(d.Spec.Selector)
	if err != nil {
		return nil, err
	}
	return test.listPods(d.Namespace, metav1.ListOptions{
		LabelSelector: selector,
	})
}

// ListPodsFromDeployment returns the list of pods created by a deployment.
func (test *Test) ListPodsFromDeployment(d *appsv1.Deployment) *v1.PodList {
	pl, err := test.listPodsFromDeployment(d)
	test.err(err)
	return pl
}

// PodReady returns whether a pod is running and each container has is in the
// ready state.
func (test *Test) PodReady(pod v1.Pod) (bool, error) {
	switch pod.Status.Phase {
	case v1.PodFailed, v1.PodSucceeded:
		return false, fmt.Errorf("pod completed")
	case v1.PodRunning:
		for _, cond := range pod.Status.Conditions {
			if cond.Type != v1.PodReady {
				continue
			}
			return cond.Status == v1.ConditionTrue, nil
		}
		return false, fmt.Errorf("pod ready condition not found")
	}
	return false, nil
}

// WaitForPodsReady waits for a selection of Pods to be running and each
// container to pass its readiness check.
func (test *Test) WaitForPodsReady(namespace string, opts metav1.ListOptions, expectedReplicas int, timeout time.Duration) error {
	return wait.Poll(time.Second, timeout, func() (bool, error) {
		pl, err := test.harness.kubeClient.Core().Pods(namespace).List(opts)
		if err != nil {
			return false, err
		}

		runningAndReady := 0
		for _, p := range pl.Items {
			isRunningAndReady, err := test.PodReady(p)
			if err != nil {
				return false, err
			}

			if isRunningAndReady {
				runningAndReady++
			}
		}

		if runningAndReady == expectedReplicas {
			return true, nil
		}
		return false, nil
	})
}

// PodLogs writes the container logs on w. If the pod has a single container,
// containerName is optional and can be set to "".
func (test *Test) PodLogs(w io.Writer, pod *v1.Pod, containerName string) error {
	if containerName == "" {
		if len(pod.Spec.Containers) != 1 {
			return fmt.Errorf("logs: no container name specified and found %d containers", len(pod.Spec.Containers))
		}
		containerName = pod.Spec.Containers[0].Name
	}

	logs, err := test.harness.kubeClient.Core().RESTClient().Get().
		Resource("pods").
		Namespace(pod.Namespace).
		Name(pod.Name).SubResource("log").
		Param("container", containerName).
		Stream()
	if err != nil {
		return err
	}

	_, err = io.Copy(w, logs)
	return err
}

// PodProxyGet returns a Request that can used to perform an HTTP GET to a pod
// through the API server proxy. Port can be a port name or the port number.
//
// If port is "", the first port found in the containers spec will be used.
func (test *Test) PodProxyGet(pod *v1.Pod, port, path string) *rest.Request {
	name := pod.Name
	if port == "" {
		name += ":" + port
	}

	return test.harness.kubeClient.
		CoreV1().
		RESTClient().
		Get().
		Namespace(pod.Namespace).
		Resource("pods").
		Name(name).
		Suffix("proxy" + path)
}

func (test *Test) podProxyGetJSON(pod *v1.Pod, port, path string, v interface{}) error {
	data, err := test.PodProxyGet(pod, port, path).DoRaw()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// PodProxyGetJSON is a convenience function around PodProxyGet that also
// unmarshals the response body into v.
func (test *Test) PodProxyGetJSON(pod *v1.Pod, port, path string, v interface{}) {
	test.err(test.podProxyGetJSON(pod, port, path, v))
}
