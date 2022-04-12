package addons

import (
	"context"
	// For go:embed
	_ "embed"
	"fmt"
	"time"

	"github.com/cloudflare/cfssl/csr"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/assetutil"
	"github.com/weaveworks/eksctl/pkg/kubernetes"

	admv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	certsv1beta1 "k8s.io/api/certificates/v1beta1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/typed/certificates/v1beta1"
)

//go:embed assets/vpc-admission-webhook-config.yaml
var vpcAdmissionWebhookConfigYaml []byte

//go:embed assets/vpc-admission-webhook-csr.yaml
var vpcAdmissionWebhookCsrYaml []byte

//go:embed assets/vpc-admission-webhook-dep.yaml
var vpcAdmissionWebhookDepYaml []byte

//go:embed assets/vpc-admission-webhook.yaml
var vpcAdmissionWebhookYaml []byte

//go:embed assets/vpc-resource-controller-dep.yaml
var vpcResourceControllerDepYaml []byte

//go:embed assets/vpc-resource-controller.yaml
var vpcResourceControllerYaml []byte

const (
	vpcControllerNamespace = metav1.NamespaceSystem
	vpcControllerName      = "vpc-resource-controller"
	webhookServiceName     = "vpc-admission-webhook"

	certWaitTimeout = 45 * time.Second
)

// NewVPCController creates a new VPCController
func NewVPCController(rawClient kubernetes.RawClientInterface, irsa IRSAHelper, clusterStatus *api.ClusterStatus, region string, planMode bool) *VPCController {
	return &VPCController{
		rawClient:     rawClient,
		irsa:          irsa,
		clusterStatus: clusterStatus,
		region:        region,
		planMode:      planMode,
	}
}

// A VPCController deploys Windows VPC controller to a cluster
type VPCController struct {
	rawClient     kubernetes.RawClientInterface
	irsa          IRSAHelper
	clusterStatus *api.ClusterStatus
	region        string
	planMode      bool
}

// Deploy deploys VPC controller to the specified cluster
func (v *VPCController) Deploy(ctx context.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if ae, ok := r.(*assetutil.Error); ok {
				err = ae
			} else {
				panic(r)
			}
		}
	}()

	if err := v.deployVPCResourceController(ctx); err != nil {
		return err
	}

	if err := v.generateCert(); err != nil {
		return err
	}

	return v.deployVPCWebhook()
}

type typeAssertionError struct {
	expected interface{}
	got      interface{}
}

func (t *typeAssertionError) Error() string {
	return fmt.Sprintf("expected type to be %T; got %T", t.expected, t.got)
}

func (v *VPCController) generateCert() error {
	var (
		csrName      = fmt.Sprintf("%s.%s", webhookServiceName, vpcControllerNamespace)
		csrClientSet = v.rawClient.ClientSet().CertificatesV1beta1().CertificateSigningRequests()
	)

	hasApprovedCert, err := v.hasApprovedCert()
	if err != nil {
		return err
	}
	if hasApprovedCert {
		// Delete existing CSR if the secret is missing
		_, err := v.rawClient.ClientSet().CoreV1().Secrets(vpcControllerNamespace).Get(context.TODO(), "vpc-admission-webhook-certs", metav1.GetOptions{})
		if err != nil {
			if !apierrors.IsNotFound(err) {
				return err
			}
			if err := csrClientSet.Delete(context.TODO(), csrName, metav1.DeleteOptions{}); err != nil {
				return err
			}
		}
		return nil
	}

	csrPEM, privateKey, err := generateCertReq(webhookServiceName, vpcControllerNamespace)
	if err != nil {
		return errors.Wrap(err, "generating CSR")
	}

	manifest := vpcAdmissionWebhookCsrYaml
	rawExtension, err := kubernetes.NewRawExtension(manifest)
	if err != nil {
		return err
	}

	certificateSigningRequest, ok := rawExtension.Object.(*certsv1beta1.CertificateSigningRequest)
	if !ok {
		return &typeAssertionError{&certsv1beta1.CertificateSigningRequest{}, rawExtension.Object}
	}

	certificateSigningRequest.Spec.Request = csrPEM
	certificateSigningRequest.Name = csrName

	if err := v.applyRawResource(certificateSigningRequest); err != nil {
		return errors.Wrap(err, "creating CertificateSigningRequest")
	}

	certificateSigningRequest.Status.Conditions = []certsv1beta1.CertificateSigningRequestCondition{
		{
			Type:           certsv1beta1.CertificateApproved,
			LastUpdateTime: metav1.NewTime(time.Now()),
			Message:        "This CSR was approved by eksctl",
			Reason:         "eksctl-approve",
		},
	}

	if _, err := csrClientSet.UpdateApproval(context.TODO(), certificateSigningRequest, metav1.UpdateOptions{}); err != nil {
		return errors.Wrap(err, "updating approval")
	}

	logger.Info("waiting for certificate to be available")

	cert, err := watchCSRApproval(csrClientSet, csrName, certWaitTimeout)
	if err != nil {
		return err
	}

	return v.createCertSecrets(privateKey, cert)
}

func watchCSRApproval(csrClientSet v1beta1.CertificateSigningRequestInterface, csrName string, timeout time.Duration) ([]byte, error) {
	watcher, err := csrClientSet.Watch(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", csrName),
	})

	if err != nil {
		return nil, err
	}

	defer watcher.Stop()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case event, ok := <-watcher.ResultChan():
			if !ok {
				return nil, errors.New("failed waiting for certificate: unexpected close of ResultChan")
			}
			switch event.Type {
			case watch.Added, watch.Modified:
				req := event.Object.(*certsv1beta1.CertificateSigningRequest)
				if cert := req.Status.Certificate; cert != nil {
					return cert, nil
				}
				logger.Warning("certificate not yet available (event: %s)", event.Type)
			}
		case <-timer.C:
			return nil, fmt.Errorf("timed out (after %v) waiting for certificate", timeout)
		}

	}
}

func (v *VPCController) createCertSecrets(key, cert []byte) error {
	secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "vpc-admission-webhook-certs",
			Namespace: vpcControllerNamespace,
		},
		Data: map[string][]byte{
			"key.pem":  key,
			"cert.pem": cert,
		},
	}

	err := v.applyRawResource(secret)
	if err != nil {
		return errors.Wrap(err, "error creating secret")
	}
	return err
}

func makePolicyDocument() map[string]interface{} {
	return map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Effect": "Allow",
				"Action": []string{
					"ec2:AssignPrivateIpAddresses",
					"ec2:DescribeInstances",
					"ec2:DescribeNetworkInterfaces",
					"ec2:UnassignPrivateIpAddresses",
					"ec2:DescribeRouteTables",
					"ec2:DescribeSubnets",
				},
				"Resource": "*",
			},
		},
	}
}

func (v *VPCController) deployVPCResourceController(ctx context.Context) error {
	irsaEnabled, err := v.irsa.IsSupported(ctx)
	if err != nil {
		return err
	}
	if irsaEnabled {
		sa := &api.ClusterIAMServiceAccount{
			ClusterIAMMeta: api.ClusterIAMMeta{
				Name:      vpcControllerName,
				Namespace: vpcControllerNamespace,
			},
			AttachPolicy: makePolicyDocument(),
		}
		if err := v.irsa.CreateOrUpdate(ctx, sa); err != nil {
			return errors.Wrap(err, "error enabling IRSA")
		}
	} else {
		// If an OIDC provider isn't associated with the cluster, the VPC controller relies on the managed policy
		// attached to the node role for the AWS VPC CNI plugin.
		sa := kubernetes.NewServiceAccount(metav1.ObjectMeta{
			Name:      vpcControllerName,
			Namespace: vpcControllerNamespace,
		})
		if err := v.applyRawResource(sa); err != nil {
			return err
		}
	}
	if err := v.applyResources(vpcResourceControllerYaml); err != nil {
		return err
	}

	return v.applyDeployment(vpcResourceControllerDepYaml)
}

func (v *VPCController) deployVPCWebhook() error {
	if err := v.applyResources(vpcAdmissionWebhookYaml); err != nil {
		return err
	}
	if err := v.applyDeployment(vpcAdmissionWebhookDepYaml); err != nil {
		return err
	}

	manifest := vpcAdmissionWebhookConfigYaml
	rawExtension, err := kubernetes.NewRawExtension(manifest)
	if err != nil {
		return err
	}

	mutatingWebhook, ok := rawExtension.Object.(*admv1beta1.MutatingWebhookConfiguration)
	if !ok {
		return &typeAssertionError{&admv1beta1.MutatingWebhookConfiguration{}, rawExtension.Object}
	}

	mutatingWebhook.Webhooks[0].ClientConfig.CABundle = v.clusterStatus.CertificateAuthorityData
	return v.applyRawResource(rawExtension.Object)
}

func (v *VPCController) hasApprovedCert() (bool, error) {
	csrClientSet := v.rawClient.ClientSet().CertificatesV1beta1().CertificateSigningRequests()
	request, err := csrClientSet.Get(context.TODO(), fmt.Sprintf("%s.%s", webhookServiceName, vpcControllerNamespace), metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return false, err
		}
		return false, nil
	}

	conditions := request.Status.Conditions
	switch len(conditions) {
	case 1:
		if conditions[0].Type == certsv1beta1.CertificateApproved {
			return true, nil
		}
		return false, fmt.Errorf("expected certificate to be approved; got %q", conditions[0].Type)

	case 0:
		return false, nil
	default:
		return false, fmt.Errorf("unexpected number of request conditions: %d", len(conditions))
	}
}

func (v *VPCController) applyResources(manifests []byte) error {
	list, err := kubernetes.NewList(manifests)
	if err != nil {
		return err
	}

	for _, item := range list.Items {
		if err := v.applyRawResource(item.Object); err != nil {
			return err
		}
	}
	return nil
}

func (v *VPCController) applyDeployment(manifests []byte) error {
	rawExtension, err := kubernetes.NewRawExtension(manifests)
	if err != nil {
		return err
	}

	deployment, ok := rawExtension.Object.(*appsv1.Deployment)
	if !ok {
		return &typeAssertionError{&appsv1.Deployment{}, rawExtension.Object}
	}
	if err := UseRegionalImage(&deployment.Spec.Template, v.region); err != nil {
		return err
	}
	return v.applyRawResource(rawExtension.Object)
}

func (v *VPCController) applyRawResource(object runtime.Object) error {
	rawResource, err := v.rawClient.NewRawResource(object)
	if err != nil {
		return err
	}

	switch newObject := object.(type) {
	case *corev1.Service:
		r, found, err := rawResource.Get()
		if err != nil {
			return err
		}
		if found {
			service, ok := r.(*corev1.Service)
			if !ok {
				return &typeAssertionError{&corev1.Service{}, r}
			}
			newObject.Spec.ClusterIP = service.Spec.ClusterIP
			newObject.SetResourceVersion(service.GetResourceVersion())
		}
	case *admv1beta1.MutatingWebhookConfiguration:
		r, found, err := rawResource.Get()
		if err != nil {
			return err
		}
		if found {
			mwc, ok := r.(*admv1beta1.MutatingWebhookConfiguration)
			if !ok {
				return &typeAssertionError{&admv1beta1.MutatingWebhookConfiguration{}, r}
			}
			newObject.SetResourceVersion(mwc.GetResourceVersion())
		}
	}

	msg, err := rawResource.CreateOrReplace(v.planMode)
	if err != nil {
		return err
	}
	logger.Info(msg)
	return nil
}

func generateCertReq(service, namespace string) ([]byte, []byte, error) {
	generator := csr.Generator{
		Validator: func(request *csr.CertificateRequest) error {
			// ignore validation as all required fields are being set
			return nil
		},
	}

	serviceCN := fmt.Sprintf("%s.%s.svc", service, namespace)

	return generator.ProcessRequest(&csr.CertificateRequest{
		KeyRequest: &csr.KeyRequest{
			A: "rsa",
			S: 2048,
		},
		CN: serviceCN,
		Hosts: []string{
			service,
			fmt.Sprintf("%s.%s", service, namespace),
			serviceCN,
		},
	})
}
