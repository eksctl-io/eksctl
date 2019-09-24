package addons

import (
	"fmt"
	"time"

	"github.com/cloudflare/cfssl/csr"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/kubernetes"

	admv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	certsv1beta1 "k8s.io/api/certificates/v1beta1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	vpcControllerNamespace = "kube-system"
	webhookServiceName     = "vpc-admission-webhook"
)

// NewVPCController creates a new VPCController
func NewVPCController(rawClient kubernetes.RawClientInterface, clusterStatus *api.ClusterStatus, region string, planMode bool) *VPCController {
	return &VPCController{
		rawClient:     rawClient,
		clusterStatus: clusterStatus,
		region:        region,
		planMode:      planMode,
	}
}

// A VPCController deploys Windows VPC controller to a cluster
type VPCController struct {
	rawClient     kubernetes.RawClientInterface
	clusterStatus *api.ClusterStatus
	region        string
	planMode      bool
}

// Deploy deploys Windows VPC controller to the specified cluster
func (v *VPCController) Deploy() error {
	if err := v.deployVPCResourceController(); err != nil {
		return err
	}

	if err := v.generateCert(); err != nil {
		return err
	}

	if err := v.deployVPCWebhook(); err != nil {
		return err
	}

	return nil
}

func (v *VPCController) generateCert() error {
	skipCSRGeneration, err := v.hasApprovedCert()
	if err != nil {
		return err
	}
	if skipCSRGeneration {
		return nil
	}

	csrPEM, privateKey, err := generateCertReq(webhookServiceName, vpcControllerNamespace)
	if err != nil {
		return errors.Wrap(err, "generating CSR")
	}

	manifest, err := vpcAdmissionWebhookCsrYamlBytes()
	if err != nil {
		return err
	}
	rawExtension, err := kubernetes.NewRawExtension(manifest)
	if err != nil {
		return err
	}

	certificateSigningRequest, ok := rawExtension.Object.(*certsv1beta1.CertificateSigningRequest)
	if !ok {
		return fmt.Errorf("expected type to be %T; got %T", &certsv1beta1.CertificateSigningRequest{}, rawExtension.Object)
	}

	certificateSigningRequest.Spec.Request = csrPEM
	certificateSigningRequest.ObjectMeta.Name = fmt.Sprintf("%s.%s", webhookServiceName, vpcControllerNamespace)

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

	csrClientSet := v.rawClient.ClientSet().CertificatesV1beta1().CertificateSigningRequests()

	if _, err := csrClientSet.UpdateApproval(certificateSigningRequest); err != nil {
		return errors.Wrap(err, "updating approval")
	}

	approvedCSR, err := csrClientSet.Get(certificateSigningRequest.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if approvedCSR.Status.Certificate == nil {
		return errors.New("failed to find certificate after approval")
	}
	return v.createCertSecrets(privateKey, approvedCSR.Status.Certificate)
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

func (v *VPCController) deployVPCResourceController() error {
	if err := v.applyResources(vpcResourceControllerYamlBytes); err != nil {
		return err
	}
	return v.applyDeployment(vpcResourceControllerDepYamlBytes)
}

func (v *VPCController) deployVPCWebhook() error {
	if err := v.applyResources(vpcAdmissionWebhookYamlBytes); err != nil {
		return err
	}
	if err := v.applyDeployment(vpcAdmissionWebhookDepYamlBytes); err != nil {
		return err
	}

	manifest, err := vpcAdmissionWebhookConfigYamlBytes()
	if err != nil {
		return err
	}
	rawExtension, err := kubernetes.NewRawExtension(manifest)
	if err != nil {
		return err
	}

	mutatingWebhook, ok := rawExtension.Object.(*admv1beta1.MutatingWebhookConfiguration)
	if !ok {
		return fmt.Errorf("expected type to be %T; got %T", &admv1beta1.MutatingWebhookConfiguration{}, rawExtension.Object)
	}

	mutatingWebhook.Webhooks[0].ClientConfig.CABundle = v.clusterStatus.CertificateAuthorityData
	return v.applyRawResource(rawExtension.Object)
}

func (v *VPCController) hasApprovedCert() (bool, error) {
	csrClientSet := v.rawClient.ClientSet().CertificatesV1beta1().CertificateSigningRequests()
	request, err := csrClientSet.Get(fmt.Sprintf("%s.%s", webhookServiceName, vpcControllerNamespace), metav1.GetOptions{})
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
			_, err := v.rawClient.ClientSet().CoreV1().Secrets(vpcControllerNamespace).Get("vpc-admission-webhook-certs", metav1.GetOptions{})
			if err != nil {
				if !apierrors.IsNotFound(err) {
					return false, err
				}
				return false, nil
			}
			return true, nil
		}
		return false, fmt.Errorf("expected certificate to be approved; got %q", conditions[0].Type)

	case 0:
		return false, nil
	default:
		return false, fmt.Errorf("unexpected number of request conditions: %d", len(conditions))
	}
}

type assetFunc func() ([]byte, error)

func (v *VPCController) applyResources(assetFn assetFunc) error {
	manifests, err := assetFn()
	if err != nil {
		return errors.Wrap(err, "unexpected error reading assets")
	}
	list, err := kubernetes.NewList([]byte(manifests))
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

func (v *VPCController) applyDeployment(assetFn assetFunc) error {
	manifests, err := assetFn()
	if err != nil {
		return errors.Wrap(err, "unexpected error reading assets")
	}
	rawExtension, err := kubernetes.NewRawExtension(manifests)
	if err != nil {
		return err
	}

	deployment, ok := rawExtension.Object.(*appsv1.Deployment)
	if !ok {
		return fmt.Errorf("expected %T; got %T", &appsv1.Deployment{}, rawExtension.Object)
	}
	useRegionalImage(&deployment.Spec.Template, v.region)
	return v.applyRawResource(rawExtension.Object)
}

func (v *VPCController) applyRawResource(object runtime.Object) error {
	if svc, ok := object.(*corev1.Service); ok {
		existingSvc, err := v.rawClient.ClientSet().CoreV1().Services(vpcControllerNamespace).Get(svc.Name, metav1.GetOptions{})
		if err != nil {
			if !apierrors.IsNotFound(err) {
				return errors.Wrapf(err, "retrieving service: %q", svc.Name)
			}
		} else {
			svc.Spec.ClusterIP = existingSvc.Spec.ClusterIP
		}
	}
	rawResource, err := v.rawClient.NewRawResource(object)
	if err != nil {
		return err
	}

	msg, err := rawResource.CreateOrReplace(v.planMode)
	if err != nil {
		return err
	}
	logger.Info(msg)
	return nil
}

// TODO use this for other addons
func useRegionalImage(spec *corev1.PodTemplateSpec, region string) {
	imageFormat := spec.Spec.Containers[0].Image
	// TODO uncomment this call after these container images are available publicly
	// regionalImage := fmt.Sprintf(imageFormat, api.EKSResourceAccountID(region), region)
	regionalImage := fmt.Sprintf(imageFormat, "940911992744", region)
	spec.Spec.Containers[0].Image = regionalImage
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
