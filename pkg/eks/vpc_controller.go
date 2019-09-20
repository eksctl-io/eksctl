package eks

import (
	"fmt"
	"time"

	"github.com/cloudflare/cfssl/csr"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/kubernetes"

	admv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	certsv1beta1 "k8s.io/api/certificates/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TODO assets

const vpcControllerResources = `
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: vpc-resource-controller
rules:
- apiGroups:
  - ""
  resources:
  - nodes
  - nodes/status
  - pods
  - configmaps
  verbs:
  - update
  - get
  - list
  - watch
  - patch
  - create
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: vpc-resource-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: vpc-resource-controller
subjects:
- kind: ServiceAccount
  name: vpc-resource-controller
  namespace: kube-system
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: vpc-resource-controller
  namespace: kube-system
---
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: vpc-resource-controller
  namespace: kube-system
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: vpc-resource-controller
        tier: backend
        track: stable
    spec:
      serviceAccount: vpc-resource-controller
      containers:
      - command:
        - /vpc-resource-controller
        args:
        - -stderrthreshold=info
        image: 940911992744.dkr.ecr.us-west-2.amazonaws.com/eks/vpc-resource-controller:0.2.0
        imagePullPolicy: Always
        livenessProbe:
          failureThreshold: 5
          httpGet:
            host: 127.0.0.1
            path: /healthz
            port: 61779
            scheme: HTTP
          initialDelaySeconds: 30
          periodSeconds: 30
          timeoutSeconds: 5
        name: vpc-resource-controller
        securityContext:
          privileged: true
      hostNetwork: true
      nodeSelector:
        beta.kubernetes.io/os: linux
        beta.kubernetes.io/arch: amd64
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: eks:kube-proxy-windows
  labels:
    k8s-app: kube-proxy
    eks.amazonaws.com/component: kube-proxy
subjects:
  - kind: Group
    name: "eks:kube-proxy-windows"
roleRef:
  kind: ClusterRole
  name: system:node-proxier
  apiGroup: rbac.authorization.k8s.io
`

const vpcWebhookResources = `
apiVersion: v1
kind: Service
metadata:
  name: vpc-admission-webhook-svc
  namespace: kube-system
  labels:
    app: vpc-admission-webhook
spec:
  ports:
  - port: 443
    targetPort: 443
  selector:
    app: vpc-admission-webhook
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: vpc-admission-webhook-deployment
  namespace: kube-system
  labels:
    app: vpc-admission-webhook
spec:
  replicas: 1
  selector:
    matchLabels:
      app: vpc-admission-webhook
  template:
    metadata:
      labels:
        app: vpc-admission-webhook
    spec:
      containers:
        - name: vpc-admission-webhook
          args:
            - -tlsCertFile=/etc/webhook/certs/cert.pem
            - -tlsKeyFile=/etc/webhook/certs/key.pem
            - -OSLabelSelectorOverride=windows
            - -alsologtostderr
            - -v=4
            - 2>&1
          image: 940911992744.dkr.ecr.us-west-2.amazonaws.com/eks/vpc-admission-webhook:0.2.0
          imagePullPolicy: Always
          volumeMounts:
            - name: webhook-certs
              mountPath: /etc/webhook/certs
              readOnly: true
      hostNetwork: true
      nodeSelector:
        beta.kubernetes.io/os: linux
        beta.kubernetes.io/arch: amd64
      volumes:
        - name: webhook-certs
          secret:
            secretName: vpc-admission-webhook-certs
`

const vpcMutatingWebhook = `
apiVersion: admissionregistration.k8s.io/v1beta1
kind: MutatingWebhookConfiguration
metadata:
  name: vpc-admission-webhook-cfg
  namespace: kube-system 
  labels:
    app: vpc-admission-webhook
webhooks:
  - name: vpc-admission-webhook.amazonaws.com
    clientConfig:
      service:
        name: vpc-admission-webhook-svc
        namespace: kube-system
        path: "/mutate"
    rules:
      - operations: [ "CREATE" ]
        apiGroups: [""]
        apiVersions: ["v1"]
        resources: ["pods"]
    failurePolicy: Ignore
`

func NewVPCController(rawClient kubernetes.RawClientInterface, clusterStatus *v1alpha5.ClusterStatus, namespace string) *VPCController {
	return &VPCController{
		rawClient:     rawClient,
		clusterStatus: clusterStatus,
		Namespace:     namespace,
	}
}

type VPCController struct {
	rawClient     kubernetes.RawClientInterface
	clusterStatus *v1alpha5.ClusterStatus
	Namespace     string
}

func (v *VPCController) Deploy() error {
	if err := v.deployResources(vpcControllerResources); err != nil {
		return err
	}

	const webhookName = "vpc-admission-webhook"

	csrPEM, privateKey, err := generateCertReq(webhookName, v.Namespace)
	if err != nil {
		return errors.Wrap(err, "generating CSR")
	}

	certificateSigningRequest := &certsv1beta1.CertificateSigningRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s.%s", webhookName, v.Namespace),
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
	if err := v.createCertSecrets(privateKey, approvedCSR.Status.Certificate); err != nil {
		return err
	}

	if err := v.deployVPCWebhook(); err != nil {
		return errors.Wrap(err, "deploying VPC webhook")
	}

	return nil
}

func (v *VPCController) createCertSecrets(key, cert []byte) error {
	_, err := v.rawClient.ClientSet().CoreV1().Secrets(v.Namespace).Create(&corev1.Secret{
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

func (v *VPCController) deployVPCWebhook() error {
	if err := v.deployResources(vpcWebhookResources); err != nil {
		return err
	}
	rawExtension, err := kubernetes.NewRawExtension([]byte(vpcMutatingWebhook))
	mutatingWebhook, ok := rawExtension.Object.(*admv1beta1.MutatingWebhookConfiguration)
	if !ok {
		return fmt.Errorf("expected type to be %T; got %T", &admv1beta1.MutatingWebhookConfiguration{}, rawExtension.Object)
	}
	mutatingWebhook.Webhooks[0].ClientConfig.CABundle = v.clusterStatus.CertificateAuthorityData
	_, err = v.rawClient.ClientSet().AdmissionregistrationV1beta1().MutatingWebhookConfigurations().Create(mutatingWebhook)
	return err
}

func (v *VPCController) deployResources(manifests string) error {
	list, err := kubernetes.NewList([]byte(manifests))
	if err != nil {
		return err
	}

	for _, item := range list.Items {
		rawResource, err := v.rawClient.NewRawResource(item)
		if err != nil {
			return err
		}

		msg, err := rawResource.CreateOrReplace(false)
		if err != nil {
			return err
		}
		logger.Info(msg)
	}
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
