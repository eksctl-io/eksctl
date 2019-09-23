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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func NewVPCController(rawClient kubernetes.RawClientInterface, clusterStatus *api.ClusterStatus, region, namespace string) *VPCController {
	return &VPCController{
		rawClient:     rawClient,
		clusterStatus: clusterStatus,
		namespace:     namespace,
		region:        region,
	}
}

type VPCController struct {
	rawClient     kubernetes.RawClientInterface
	clusterStatus *api.ClusterStatus
	region        string
	namespace     string
}

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
	const webhookName = "vpc-admission-webhook"

	csrPEM, privateKey, err := generateCertReq(webhookName, v.namespace)
	if err != nil {
		return errors.Wrap(err, "generating CSR")
	}

	certificateSigningRequest := &certsv1beta1.CertificateSigningRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s.%s", webhookName, v.namespace),
		},
		Spec: certsv1beta1.CertificateSigningRequestSpec{
			Request: csrPEM,
			Usages: []certsv1beta1.KeyUsage{
				certsv1beta1.UsageSigning,
				certsv1beta1.UsageKeyEncipherment,
				certsv1beta1.UsageServerAuth,
			},
			Groups: []string{"system:authenticated"},
		},
	}

	csrClientSet := v.rawClient.ClientSet().CertificatesV1beta1().CertificateSigningRequests()

	// TODO create or replace
	request, err := csrClientSet.Create(certificateSigningRequest)
	if err != nil {
		return errors.Wrap(err, "creating CSR")
	}

	request.Status.Conditions = []certsv1beta1.CertificateSigningRequestCondition{
		{
			Type:           certsv1beta1.CertificateApproved,
			LastUpdateTime: metav1.NewTime(time.Now()),
			Message:        "This CSR was approved by eksctl",
			Reason:         "eksctl-approve",
		},
	}

	if _, err := csrClientSet.UpdateApproval(request); err != nil {
		return errors.Wrap(err, "updating approval")
	}

	approvedCSR, err := csrClientSet.Get(request.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if approvedCSR.Status.Certificate == nil {
		return errors.New("failed to find certificate after approval")
	}
	return v.createCertSecrets(privateKey, approvedCSR.Status.Certificate)
}

func (v *VPCController) createCertSecrets(key, cert []byte) error {
	_, err := v.rawClient.ClientSet().CoreV1().Secrets(v.namespace).Create(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "vpc-admission-webhook-certs",
		},
		Data: map[string][]byte{
			"key.pem":  key,
			"cert.pem": cert,
		},
	})

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
	return v.applyRawResource(rawExtension)
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
		if err := v.applyRawResource(item); err != nil {
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
	return v.applyRawResource(rawExtension)
}

func (v *VPCController) applyRawResource(r runtime.RawExtension) error {
	rawResource, err := v.rawClient.NewRawResource(r)
	if err != nil {
		return err
	}

	msg, err := rawResource.CreateOrReplace(false)
	if err != nil {
		return err
	}
	logger.Info(msg)
	return nil
}

// TODO use the same function for other addons
func useRegionalImage(spec *corev1.PodTemplateSpec, region string) {
	imageFormat := spec.Spec.Containers[0].Image
	regionalImage := fmt.Sprintf(imageFormat, api.EKSResourceAccountID(region), region)
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
