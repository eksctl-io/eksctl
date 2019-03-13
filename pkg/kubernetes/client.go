package kubernetes

import (
	"fmt"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/genericclioptions/resource"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

// RawClient stores information about the client config
type RawClient struct {
	mapper    meta.RESTMapper
	config    *restclient.Config
	clientSet kubeclient.Interface
}

// RawClientInterface defines high level abstraction for RawClient for testing
type RawClientInterface interface {
	ClientSet() kubeclient.Interface
	NewRawResource(runtime.RawExtension) (*RawResource, error)
}

// RawResource holds info about a resource along with a type-specific raw client instance
type RawResource struct {
	Helper *resource.Helper
	Info   *resource.Info
	GVK    *schema.GroupVersionKind
}

// NewRawClient creates a new raw REST client
func NewRawClient(clientSet kubeclient.Interface, config *restclient.Config) (*RawClient, error) {
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
		c.config.NegotiatedSerializer = &serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	}

	if err := restclient.SetKubernetesDefaults(c.config); err != nil {
		return nil, errors.Wrap(err, "applying defaults for REST client")
	}
	return c, nil
}

// ClientSet returns the underlying ClientSet
func (c *RawClient) ClientSet() kubeclient.Interface { return c.clientSet }

// NewHelperFor construct a raw client helper instance for a give gvk
// (it's based on k8s.io/kubernetes/pkg/kubectl/cmd/util/factory_client_access.go)
func (c *RawClient) NewHelperFor(gvk schema.GroupVersionKind) (*resource.Helper, error) {
	mapping, err := c.mapper.RESTMapping(gvk.GroupKind(), gvk.GroupVersion().Version, "")
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

// NewRawResource constructs a type-specific instance or RawClient for rawObj
func (c *RawClient) NewRawResource(rawObj runtime.RawExtension) (*RawResource, error) {
	gvk := rawObj.Object.GetObjectKind().GroupVersionKind()

	obj, ok := rawObj.Object.(metav1.Object)
	if !ok {
		return nil, fmt.Errorf("cannot conver object of type %T to metav1.Object", rawObj.Object)
	}

	info := &resource.Info{
		Name:      obj.GetName(),
		Namespace: obj.GetNamespace(),
		Object:    rawObj.Object,
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
func (r *RawResource) LogAction(verb string) string {
	return fmt.Sprintf("%s %q", verb, r)
}

// CreateOrReplace will check if the given resource exists, and create or update it as needed
func (r *RawResource) CreateOrReplace() (string, error) {
	create := false
	if _, err := r.Helper.Get(r.Info.Namespace, r.Info.Name, false); err != nil {
		if !apierrs.IsNotFound(err) {
			return "", errors.Wrap(err, "unexpected non-404 error")
		}
		create = true
	}

	if create {
		_, err := r.Helper.Create(r.Info.Namespace, true, r.Info.Object, &metav1.CreateOptions{})
		if err != nil {
			return "", err
		}

		return r.LogAction("created"), nil
	}

	convertedObj, err := scheme.Scheme.ConvertToVersion(r.Info.Object, r.GVK.GroupVersion())
	if err != nil {
		return "", errors.Wrapf(err, "converting object")
	}
	scheme.Scheme.Default(convertedObj)

	if _, err := r.Helper.Replace(r.Info.Namespace, r.Info.Name, true, r.Info.Object); err != nil {
		return "", err
	}

	return r.LogAction("replaced"), nil
}

/*

	This doesn't work yet. We need to find a way to do defaulting properly, what we have now seems to cause following behaviour and nothing seems to make it go away.

	2019-03-08T10:24:27Z [▶]  oldData["kube-system:DaemonSet.extensions/aws-node"] = {"metadata":{"name":"aws-node","namespace":"kube-system","selfLink":"/apis/extensions/v1beta1/namespaces/kube-system/daemonsets/aws-node","uid":"5ec47d60-3f39-11e9-a23b-0616486ecb7e","resourceVersion":"353986","generation":2,"creationTimestamp":"2019-03-05T11:25:24Z","labels":{"k8s-app":"aws-node"}},"spec":{"selector":{"matchLabels":{"k8s-app":"aws-node"}},"template":{"metadata":{"creationTimestamp":null,"labels":{"k8s-app":"aws-node"},"annotations":{"scheduler.alpha.kubernetes.io/critical-pod":""}},"spec":{"volumes":[{"name":"cni-bin-dir","hostPath":{"path":"/opt/cni/bin","type":""}},{"name":"cni-net-dir","hostPath":{"path":"/etc/cni/net.d","type":""}},{"name":"log-dir","hostPath":{"path":"/var/log","type":""}},{"name":"dockersock","hostPath":{"path":"/var/run/docker.sock","type":""}}],"containers":[{"name":"aws-node","image":"602401143452.dkr.ecr.us-west-2.amazonaws.com/amazon-k8s-cni:v1.3.2","ports":[{"name":"metrics","hostPort":61678,"containerPort":61678,"protocol":"TCP"}],"env":[{"name":"AWS_VPC_K8S_CNI_LOGLEVEL","value":"DEBUG"},{"name":"MY_NODE_NAME","valueFrom":{"fieldRef":{"apiVersion":"v1","fieldPath":"spec.nodeName"}}},{"name":"WATCH_NAMESPACE","valueFrom":{"fieldRef":{"apiVersion":"v1","fieldPath":"metadata.namespace"}}}],"resources":{"requests":{"cpu":"10m"}},"volumeMounts":[{"name":"cni-bin-dir","mountPath":"/host/opt/cni/bin"},{"name":"cni-net-dir","mountPath":"/host/etc/cni/net.d"},{"name":"log-dir","mountPath":"/host/var/log"},{"name":"dockersock","mountPath":"/var/run/docker.sock"}],"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"Always","securityContext":{"privileged":true}}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","serviceAccountName":"aws-node","serviceAccount":"aws-node","hostNetwork":true,"securityContext":{},"schedulerName":"default-scheduler","tolerations":[{"operator":"Exists"}]}},"updateStrategy":{"type":"RollingUpdate","rollingUpdate":{"maxUnavailable":1}},"templateGeneration":2,"revisionHistoryLimit":10},"status":{"currentNumberScheduled":2,"numberMisscheduled":0,"desiredNumberScheduled":2,"numberReady":2,"observedGeneration":2,"updatedNumberScheduled":2,"numberAvailable":2}}

	2019-03-08T10:24:27Z [▶]  newData["kube-system:DaemonSet.extensions/aws-node"] = {"kind":"DaemonSet","apiVersion":"extensions/v1beta1","metadata":{"name":"aws-node","namespace":"kube-system","creationTimestamp":null,"labels":{"k8s-app":"aws-node"}},"spec":{"selector":{"matchLabels":{"k8s-app":"aws-node"}},"template":{"metadata":{"creationTimestamp":null,"labels":{"k8s-app":"aws-node"},"annotations":{"scheduler.alpha.kubernetes.io/critical-pod":""}},"spec":{"volumes":[{"name":"cni-bin-dir","hostPath":{"path":"/opt/cni/bin"}},{"name":"cni-net-dir","hostPath":{"path":"/etc/cni/net.d"}},{"name":"log-dir","hostPath":{"path":"/var/log"}},{"name":"dockersock","hostPath":{"path":"/var/run/docker.sock"}}],"containers":[{"name":"aws-node","image":"602401143452.dkr.ecr.us-west-2.amazonaws.com/amazon-k8s-cni:v1.3.2","ports":[{"name":"metrics","containerPort":61678}],"env":[{"name":"AWS_VPC_K8S_CNI_LOGLEVEL","value":"DEBUG"},{"name":"MY_NODE_NAME","valueFrom":{"fieldRef":{"fieldPath":"spec.nodeName"}}},{"name":"WATCH_NAMESPACE","valueFrom":{"fieldRef":{"fieldPath":"metadata.namespace"}}}],"resources":{"requests":{"cpu":"10m"}},"volumeMounts":[{"name":"cni-bin-dir","mountPath":"/host/opt/cni/bin"},{"name":"cni-net-dir","mountPath":"/host/etc/cni/net.d"},{"name":"log-dir","mountPath":"/host/var/log"},{"name":"dockersock","mountPath":"/var/run/docker.sock"}],"imagePullPolicy":"Always","securityContext":{"privileged":true}}],"serviceAccountName":"aws-node","hostNetwork":true,"tolerations":[{"operator":"Exists"}]}},"updateStrategy":{"type":"RollingUpdate"}},"status":{"currentNumberScheduled":0,"numberMisscheduled":0,"desiredNumberScheduled":0,"numberReady":0}}

	2019-03-08T10:24:27Z [▶]  patch["kube-system:DaemonSet.extensions/aws-node"] = {"apiVersion":"extensions/v1beta1","kind":"DaemonSet","metadata":{"creationTimestamp":null,"generation":null,"resourceVersion":null,"selfLink":null,"uid":null},"spec":{"revisionHistoryLimit":null,"template":{"spec":{"$setElementOrder/containers":[{"name":"aws-node"}],"$setElementOrder/volumes":[{"name":"cni-bin-dir"},{"name":"cni-net-dir"},{"name":"log-dir"},{"name":"dockersock"}],"containers":[{"$setElementOrder/env":[{"name":"AWS_VPC_K8S_CNI_LOGLEVEL"},{"name":"MY_NODE_NAME"},{"name":"WATCH_NAMESPACE"}],"$setElementOrder/ports":[{"containerPort":61678}],"env":[{"name":"MY_NODE_NAME","valueFrom":{"fieldRef":{"apiVersion":null}}},{"name":"WATCH_NAMESPACE","valueFrom":{"fieldRef":{"apiVersion":null}}}],"name":"aws-node","ports":[{"containerPort":61678,"hostPort":null,"protocol":null}],"terminationMessagePath":null,"terminationMessagePolicy":null}],"dnsPolicy":null,"restartPolicy":null,"schedulerName":null,"securityContext":null,"serviceAccount":null,"terminationGracePeriodSeconds":null,"volumes":[{"$retainKeys":["hostPath","name"],"hostPath":{"type":null},"name":"cni-bin-dir"},{"$retainKeys":["hostPath","name"],"hostPath":{"type":null},"name":"cni-net-dir"},{"$retainKeys":["hostPath","name"],"hostPath":{"type":null},"name":"log-dir"},{"$retainKeys":["hostPath","name"],"hostPath":{"type":null},"name":"dockersock"}]}},"templateGeneration":null,"updateStrategy":{"rollingUpdate":null}},"status":{"currentNumberScheduled":0,"desiredNumberScheduled":0,"numberAvailable":null,"numberReady":0,"observedGeneration":null,"updatedNumberScheduled":null}}
	2019-03-08T10:24:27Z [✖]  DaemonSet.apps "aws-node" is invalid: [spec.template.spec.volumes[0].hostPath.path: Required value, spec.template.spec.volumes[1].hostPath.path: Required value, spec.template.spec.volumes[2].hostPath.path: Required value, spec.template.spec.volumes[3].hostPath.path: Required value, spec.template.spec.containers[0].image: Required value, spec.template.spec.containers[0].env[0].valueFrom.fieldRef.fieldPath: Required value, spec.template.spec.containers[0].env[1].valueFrom.fieldRef.fieldPath: Required value]
*/

// CreatePatchOrReplace attempts patching the resource before replacing it
// TODO: it needs more testing and the issue sith strategicpatch has to be
// undertood before we decide wheather to use it or not
func (r *RawResource) CreatePatchOrReplace() error {
	msg := func(verb string) { logger.Info("%s %q", verb, r) }

	create := false
	oldObj, err := r.Helper.Get(r.Info.Namespace, r.Info.Name, false)
	if err != nil {
		if !apierrs.IsNotFound(err) {
			create = true
		}
		return err
	}

	if create {
		_, err := r.Helper.Create(r.Info.Namespace, true, r.Info.Object, &metav1.CreateOptions{})
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
