package kubernetes

import (
	"fmt"
	"time"

	"github.com/blang/semver"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/pkg/errors"
	"github.com/weaveworks/logger"

	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/discovery"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

// Interface is an alias to avoid having to import k8s.io/client-go/kubernetes
// along with this package, so that most of our packages only care to import
// our kubernetes package
type Interface = kubeclient.Interface

// ClientSetGetter is an interface used for anything that requires
// to obtain Kubernetes client whe it's not possible to pass the
// client directly
type ClientSetGetter interface {
	ClientSet() (Interface, error)
}

// CachedClientSet provides a basic implementation of ClientSetGetter
// where the client is a field of a struct
type CachedClientSet struct {
	CachedClientSet Interface
}

// NewCachedClientSet costructs a new CachedClientSets
func NewCachedClientSet(clientSet Interface) *CachedClientSet {
	return &CachedClientSet{CachedClientSet: clientSet}
}

// ClientSet returns g.CachedClientSet or an error it is nil
func (g *CachedClientSet) ClientSet() (Interface, error) {
	if g.CachedClientSet == nil {
		return nil, fmt.Errorf("no client instance provided")
	}
	return g.CachedClientSet, nil
}

// CallbackClientSet provides an implementation of ClientSetGetter
// where the client is provided via a callback
type CallbackClientSet struct {
	Callback func() (Interface, error)
}

// ClientSet returns g.ClientSet or an error it is nil
func (g *CallbackClientSet) ClientSet() (Interface, error) {
	return g.Callback()
}

// RawClient stores information about the client config
type RawClient struct {
	mapper    meta.RESTMapper
	config    *restclient.Config
	clientSet Interface
}

// RawClientInterface defines high level abstraction for RawClient for testing
type RawClientInterface interface {
	ClientSet() Interface
	NewRawResource(runtime.Object) (*RawResource, error)
}

// RawResource holds info about a resource along with a type-specific raw client instance
type RawResource struct {
	Helper *resource.Helper
	Info   *resource.Info
	GVK    *schema.GroupVersionKind
}

// NewRawClient creates a new raw REST client
func NewRawClient(clientSet Interface, config *restclient.Config) (*RawClient, error) {
	c := &RawClient{
		config:    config,
		clientSet: clientSet,
	}

	return c.new()
}

func (c *RawClient) new() (*RawClient, error) {
	apiGroupResources, err := restmapper.GetAPIGroupResources(c.ClientSet().Discovery())
	if err != nil {
		return nil, errors.Wrap(err, "getting list of API resources for raw REST client")
	}

	c.mapper = restmapper.NewDiscoveryRESTMapper(apiGroupResources)

	if c.config.APIPath == "" {
		c.config.APIPath = "/api"
	}
	if c.config.NegotiatedSerializer == nil {
		c.config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	}

	if err := restclient.SetKubernetesDefaults(c.config); err != nil {
		return nil, errors.Wrap(err, "applying defaults for REST client")
	}
	return c, nil
}

func getServerVersion(discoveryClient discovery.DiscoveryInterface) (string, error) {
	v, err := discoveryClient.ServerVersion()
	if err != nil {
		return "", errors.Wrapf(err, "getting Kubernetes API version")
	}

	sv, err := semver.ParseTolerant(v.GitVersion)
	if err != nil {
		return "", errors.Wrapf(err, "parsing Kubernetes API version")
	}

	sv.Pre = nil // clear extra info

	return sv.String(), nil
}

// ServerVersion will use discovery API to fetch version of Kubernetes control plane
func (c *RawClient) ServerVersion() (string, error) {
	return getServerVersion(c.ClientSet().Discovery())
}

// ClientSet returns the underlying ClientSet
func (c *RawClient) ClientSet() Interface { return c.clientSet }

// NewHelperFor construct a raw client helper instance for a give gvk
// (it's based on k8s.io/kubernetes/pkg/kubectl/cmd/util/factory_client_access.go)
func (c *RawClient) NewHelperFor(gvk schema.GroupVersionKind) (*resource.Helper, error) {
	mapping, err := c.mapper.RESTMapping(gvk.GroupKind(), gvk.Version, "")
	if err != nil {
		return nil, errors.Wrapf(err, "constructing REST client mapping for %s", gvk.String())
	}

	switch gvk.Group {
	case corev1.GroupName:
		c.config.APIPath = "/api"
	default:
		c.config.APIPath = "/apis"
	}
	gv := gvk.GroupVersion()
	c.config.GroupVersion = &gv

	client, err := restclient.RESTClientFor(c.config)
	if err != nil {
		return nil, errors.Wrapf(err, "constructing REST client for %s", gvk.String())
	}

	return resource.NewHelper(client, mapping), nil
}

// NewRawResource constructs a type-specific instance or RawClient for object
func (c *RawClient) NewRawResource(object runtime.Object) (*RawResource, error) {
	gvk := object.GetObjectKind().GroupVersionKind()

	metaObj, ok := object.(metav1.Object)
	if !ok {
		return nil, fmt.Errorf("cannot convert object of type %T to metav1.Object", object)
	}

	info := &resource.Info{
		Name:      metaObj.GetName(),
		Namespace: metaObj.GetNamespace(),
		Object:    object,
	}

	helper, err := c.NewHelperFor(gvk)
	if err != nil {
		return nil, err
	}

	r := &RawResource{
		Helper: helper,
		Info:   info,
		GVK:    &gvk,
	}

	return r, nil
}

// CreateOrReplace will check if the resources in the provided manifest exists,
// and create or update them as needed.
func (c *RawClient) CreateOrReplace(manifest []byte, plan bool) error {
	objects, err := NewRawExtensions(manifest)
	if err != nil {
		return err
	}
	for _, object := range objects {
		if err := c.createOrReplaceObject(object, plan); err != nil {
			return err
		}
	}
	return nil
}

func (c *RawClient) createOrReplaceObject(object runtime.RawExtension, plan bool) error {
	resource, err := c.NewRawResource(object.Object)
	if err != nil {
		return err
	}
	status, err := resource.CreateOrReplace(plan)
	if err != nil {
		return err
	}
	logger.Info(status)
	return nil
}

// Delete attempts to delete the Kubernetes resources in the provided manifest,
// or do nothing if they do not exist.
func (c *RawClient) Delete(manifest []byte) error {
	objects, err := NewRawExtensions(manifest)
	if err != nil {
		return err
	}
	for _, object := range objects {
		if err := c.deleteObject(object); err != nil {
			return err
		}
	}
	return nil
}

func (c *RawClient) deleteObject(object runtime.RawExtension) error {
	resource, err := c.NewRawResource(object.Object)
	if err != nil {
		return err
	}
	status, err := resource.DeleteSync()
	if err != nil {
		return err
	}
	if status != "" {
		logger.Info(status)
	}
	return nil
}

// Exists checks if the Kubernetes resources in the provided manifest exist or
// not, and returns a map[<namespace>]map[<name>]bool to indicate each
// resource's existence.
func (c *RawClient) Exists(manifest []byte) (map[string]map[string]bool, error) {
	objects, err := NewRawExtensions(manifest)
	if err != nil {
		return nil, err
	}
	existence := map[string]map[string]bool{}
	for _, object := range objects {
		resource, err := c.NewRawResource(object.Object)
		if err != nil {
			return nil, err
		}
		exists, err := resource.Exists()
		if err != nil {
			return nil, err
		}
		if _, ok := existence[resource.Info.Namespace]; !ok {
			existence[resource.Info.Namespace] = map[string]bool{}
		}
		existence[resource.Info.Namespace][resource.Info.Name] = exists
	}
	return existence, nil
}

// String returns a canonical name of the resource
func (r *RawResource) String() string {
	description := ""
	if r.Info.Namespace != "" {
		description += r.Info.Namespace + ":"
	}
	description += r.GVK.Kind
	if r.GVK.Group != "" {
		description += "." + r.GVK.Group
	}
	description += "/" + r.Info.Name
	return description
}

// LogAction returns an info message that can be used to log a particular actions
func (r *RawResource) LogAction(plan bool, verb string) string {
	if plan {
		return fmt.Sprintf("(plan) would have %s %q", verb, r)
	}
	return fmt.Sprintf("%s %q", verb, r)
}

// CreateOrReplace will check if the given resource exists, and create or update it as needed
func (r *RawResource) CreateOrReplace(plan bool) (string, error) {
	_, exists, err := r.Get()
	if err != nil {
		return "", errors.Wrap(err, "unexpected non-404 error")
	}
	if !exists {
		if !plan {
			_, err := r.Helper.Create(r.Info.Namespace, true, r.Info.Object)
			if err != nil {
				return "", err
			}
		}
		return r.LogAction(plan, "created"), nil
	}

	convertedObj, err := scheme.Scheme.ConvertToVersion(r.Info.Object, r.GVK.GroupVersion())
	if err != nil {
		return "", errors.Wrapf(err, "converting object")
	}
	scheme.Scheme.Default(convertedObj)
	if !plan {
		if _, err := r.Helper.Replace(r.Info.Namespace, r.Info.Name, true, r.Info.Object); err != nil {
			return "", err
		}
	}

	return r.LogAction(plan, "replaced"), nil
}

/*

	This doesn't work yet. We need to find a way to do defaulting properly, what we have now seems to cause following behaviour and nothing seems to make it go away.

	2019-03-08T10:24:27Z [▶]  oldData["kube-system:DaemonSet.extensions/aws-node"] = {"metadata":{"name":"aws-node","namespace":"kube-system","selfLink":"/apis/extensions/v1beta1/namespaces/kube-system/daemonsets/aws-node","uid":"5ec47d60-3f39-11e9-a23b-0616486ecb7e","resourceVersion":"353986","generation":2,"creationTimestamp":"2019-03-05T11:25:24Z","labels":{"k8s-app":"aws-node"}},"spec":{"selector":{"matchLabels":{"k8s-app":"aws-node"}},"template":{"metadata":{"creationTimestamp":null,"labels":{"k8s-app":"aws-node"},"annotations":{"scheduler.alpha.kubernetes.io/critical-pod":""}},"spec":{"volumes":[{"name":"cni-bin-dir","hostPath":{"path":"/opt/cni/bin","type":""}},{"name":"cni-net-dir","hostPath":{"path":"/etc/cni/net.d","type":""}},{"name":"log-dir","hostPath":{"path":"/var/log","type":""}},{"name":"dockersock","hostPath":{"path":"/var/run/docker.sock","type":""}}],"containers":[{"name":"aws-node","image":"602401143452.dkr.ecr.us-west-2.amazonaws.com/amazon-k8s-cni:v1.3.2","ports":[{"name":"metrics","hostPort":61678,"containerPort":61678,"protocol":"TCP"}],"env":[{"name":"AWS_VPC_K8S_CNI_LOGLEVEL","value":"DEBUG"},{"name":"MY_NODE_NAME","valueFrom":{"fieldRef":{"apiVersion":"v1","fieldPath":"spec.nodeName"}}},{"name":"WATCH_NAMESPACE","valueFrom":{"fieldRef":{"apiVersion":"v1","fieldPath":"metadata.namespace"}}}],"resources":{"requests":{"cpu":"10m"}},"volumeMounts":[{"name":"cni-bin-dir","mountPath":"/host/opt/cni/bin"},{"name":"cni-net-dir","mountPath":"/host/etc/cni/net.d"},{"name":"log-dir","mountPath":"/host/var/log"},{"name":"dockersock","mountPath":"/var/run/docker.sock"}],"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"Always","securityContext":{"privileged":true}}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","serviceAccountName":"aws-node","serviceAccount":"aws-node","hostNetwork":true,"securityContext":{},"schedulerName":"default-scheduler","tolerations":[{"operator":"Exists"}]}},"updateStrategy":{"type":"RollingUpdate","rollingUpdate":{"maxUnavailable":1}},"templateGeneration":2,"revisionHistoryLimit":10},"status":{"currentNumberScheduled":2,"numberMisscheduled":0,"desiredNumberScheduled":2,"numberReady":2,"observedGeneration":2,"updatedNumberScheduled":2,"numberAvailable":2}}

	2019-03-08T10:24:27Z [▶]  newData["kube-system:DaemonSet.extensions/aws-node"] = {"kind":"DaemonSet","apiVersion":"extensions/v1beta1","metadata":{"name":"aws-node","namespace":"kube-system","creationTimestamp":null,"labels":{"k8s-app":"aws-node"}},"spec":{"selector":{"matchLabels":{"k8s-app":"aws-node"}},"template":{"metadata":{"creationTimestamp":null,"labels":{"k8s-app":"aws-node"},"annotations":{"scheduler.alpha.kubernetes.io/critical-pod":""}},"spec":{"volumes":[{"name":"cni-bin-dir","hostPath":{"path":"/opt/cni/bin"}},{"name":"cni-net-dir","hostPath":{"path":"/etc/cni/net.d"}},{"name":"log-dir","hostPath":{"path":"/var/log"}},{"name":"dockersock","hostPath":{"path":"/var/run/docker.sock"}}],"containers":[{"name":"aws-node","image":"602401143452.dkr.ecr.us-west-2.amazonaws.com/amazon-k8s-cni:v1.3.2","ports":[{"name":"metrics","containerPort":61678}],"env":[{"name":"AWS_VPC_K8S_CNI_LOGLEVEL","value":"DEBUG"},{"name":"MY_NODE_NAME","valueFrom":{"fieldRef":{"fieldPath":"spec.nodeName"}}},{"name":"WATCH_NAMESPACE","valueFrom":{"fieldRef":{"fieldPath":"metadata.namespace"}}}],"resources":{"requests":{"cpu":"10m"}},"volumeMounts":[{"name":"cni-bin-dir","mountPath":"/host/opt/cni/bin"},{"name":"cni-net-dir","mountPath":"/host/etc/cni/net.d"},{"name":"log-dir","mountPath":"/host/var/log"},{"name":"dockersock","mountPath":"/var/run/docker.sock"}],"imagePullPolicy":"Always","securityContext":{"privileged":true}}],"serviceAccountName":"aws-node","hostNetwork":true,"tolerations":[{"operator":"Exists"}]}},"updateStrategy":{"type":"RollingUpdate"}},"status":{"currentNumberScheduled":0,"numberMisscheduled":0,"desiredNumberScheduled":0,"numberReady":0}}

	2019-03-08T10:24:27Z [▶]  patch["kube-system:DaemonSet.extensions/aws-node"] = {"apiVersion":"extensions/v1beta1","kind":"DaemonSet","metadata":{"creationTimestamp":null,"generation":null,"resourceVersion":null,"selfLink":null,"uid":null},"spec":{"revisionHistoryLimit":null,"template":{"spec":{"$setElementOrder/containers":[{"name":"aws-node"}],"$setElementOrder/volumes":[{"name":"cni-bin-dir"},{"name":"cni-net-dir"},{"name":"log-dir"},{"name":"dockersock"}],"containers":[{"$setElementOrder/env":[{"name":"AWS_VPC_K8S_CNI_LOGLEVEL"},{"name":"MY_NODE_NAME"},{"name":"WATCH_NAMESPACE"}],"$setElementOrder/ports":[{"containerPort":61678}],"env":[{"name":"MY_NODE_NAME","valueFrom":{"fieldRef":{"apiVersion":null}}},{"name":"WATCH_NAMESPACE","valueFrom":{"fieldRef":{"apiVersion":null}}}],"name":"aws-node","ports":[{"containerPort":61678,"hostPort":null,"protocol":null}],"terminationMessagePath":null,"terminationMessagePolicy":null}],"dnsPolicy":null,"restartPolicy":null,"schedulerName":null,"securityContext":null,"serviceAccount":null,"terminationGracePeriodSeconds":null,"volumes":[{"$retainKeys":["hostPath","name"],"hostPath":{"type":null},"name":"cni-bin-dir"},{"$retainKeys":["hostPath","name"],"hostPath":{"type":null},"name":"cni-net-dir"},{"$retainKeys":["hostPath","name"],"hostPath":{"type":null},"name":"log-dir"},{"$retainKeys":["hostPath","name"],"hostPath":{"type":null},"name":"dockersock"}]}},"templateGeneration":null,"updateStrategy":{"rollingUpdate":null}},"status":{"currentNumberScheduled":0,"desiredNumberScheduled":0,"numberAvailable":null,"numberReady":0,"observedGeneration":null,"updatedNumberScheduled":null}}
	2019-03-08T10:24:27Z [✖]  DaemonSet.apps "aws-node" is invalid: [spec.template.spec.volumes[0].hostPath.path: Required value, spec.template.spec.volumes[1].hostPath.path: Required value, spec.template.spec.volumes[2].hostPath.path: Required value, spec.template.spec.volumes[3].hostPath.path: Required value, spec.template.spec.containers[0].image: Required value, spec.template.spec.containers[0].env[0].valueFrom.fieldRef.fieldPath: Required value, spec.template.spec.containers[0].env[1].valueFrom.fieldRef.fieldPath: Required value]
*/

// CreatePatchOrReplace attempts patching the resource before replacing it
// TODO: it needs more testing and the issue with strategic patch has to be
// understood before we decide whether to use it or not
func (r *RawResource) CreatePatchOrReplace() error {
	msg := func(verb string) { logger.Info("%s %q", verb, r) }

	oldObj, exists, err := r.Get()
	if err != nil {
		return err
	}

	if !exists {
		_, err := r.Helper.Create(r.Info.Namespace, true, r.Info.Object)
		if err != nil {
			return err
		}
		msg("created")
		return nil
	}

	oldData, err := runtime.Encode(unstructured.UnstructuredJSONScheme, oldObj)
	if err != nil {
		return err
	}

	convertedObj, err := scheme.Scheme.ConvertToVersion(r.Info.Object.DeepCopyObject(), r.GVK.GroupVersion())
	if err != nil {
		return errors.Wrapf(err, "converting object")
	}
	scheme.Scheme.Default(convertedObj)
	newData, err := runtime.Encode(unstructured.UnstructuredJSONScheme, convertedObj)
	if err != nil {
		return err
	}

	// lookupPatchMeta, err := strategicpatch.NewPatchMetaFromStruct(r.Info.Object)
	// if err != nil {
	// 	return err
	// }
	// patch, err := strategicpatch.CreateThreeWayMergePatch(oldData, newData, newData, lookupPatchMeta, true)

	// versionedObject, err := scheme.Scheme.New(*r.gvk)
	// if err != nil {
	// 	return err
	// }
	// patch, err := strategicpatch.CreateTwoWayMergePatch(oldData, newData, convertedObj)

	patch, err := jsonpatch.CreateMergePatch(oldData, newData)
	if err != nil {
		logger.Warning("could create patch for %q: %v", r, err)
		_, err := r.Helper.Replace(r.Info.Namespace, r.Info.Name, false, r.Info.Object)
		if err != nil {
			return err
		}
		msg("replaced")
		return nil
	}
	logger.Debug("oldData[%q] = %s", r, oldData)
	logger.Debug("newData[%q] = %s", r, newData)
	logger.Debug("patch[%q] = %s", r, patch)
	_, err = r.Helper.Patch(r.Info.Namespace, r.Info.Name, types.MergePatchType, patch, nil)
	if err != nil {
		return err
	}
	msg("patched")
	return nil
}

// DeleteSync attempts to delete this Kubernetes resource, or returns doing
// nothing if it does not exist. It blocks until the resource has been deleted.
func (r *RawResource) DeleteSync() (string, error) {
	_, exists, err := r.Get()
	if err != nil {
		return "", err
	}
	if !exists {
		return "", nil
	}
	propagationPolicy := metav1.DeletePropagationForeground
	if _, err := r.Helper.DeleteWithOptions(r.Info.Namespace, r.Info.Name, &metav1.DeleteOptions{
		PropagationPolicy: &propagationPolicy,
	}); err != nil {
		return "", err
	}
	if err := r.waitForDeletion(); err != nil {
		return "", err
	}
	return r.LogAction(false, "deleted"), nil
}

const maxWaitingTime = 2 * 60 * time.Second

func (r *RawResource) waitForDeletion() error {
	// Wait for the resource's deletion, typically to avoid "races" as much as
	// possible on eksctl's side, as objects may be still "TERMINATING" while
	// eksctl then tries to create them again.
	waitingTime := maxWaitingTime
	checkInterval := 1 * time.Second
	for waitingTime > 0 {
		_, exists, err := r.Get()
		if err != nil {
			return err
		}
		if !exists {
			return nil
		}
		time.Sleep(checkInterval)
		waitingTime = waitingTime - checkInterval
		checkInterval = 2 * checkInterval
	}
	return fmt.Errorf("waited for %v's deletion, but could not confirm it within %v", r, maxWaitingTime)
}

// Exists checks if this Kubernetes resource exists or not, and returns true if
// so, or false otherwise.
func (r *RawResource) Exists() (bool, error) {
	_, exists, err := r.Get()
	return exists, err
}

// Get returns the Kubernetes resource from the server
func (r *RawResource) Get() (runtime.Object, bool, error) {
	obj, err := r.Helper.Get(r.Info.Namespace, r.Info.Name)
	if err != nil {
		if apierrs.IsNotFound(err) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return obj, true, nil
}
