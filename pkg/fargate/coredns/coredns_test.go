package coredns_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/fargate/coredns"
	"github.com/weaveworks/eksctl/pkg/testutils"
	v1 "k8s.io/api/core/v1"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestSuite(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var (
	minimumProfileSelectingCoreDNS = []*api.FargateProfile{
		&api.FargateProfile{
			Name: "min-selecting-coredns",
			Selectors: []api.FargateProfileSelector{
				api.FargateProfileSelector{Namespace: "kube-system"},
			},
		},
	}

	multipleProfilesWithOneSelectingCoreDNS = []*api.FargateProfile{
		&api.FargateProfile{
			Name: "foo",
			Selectors: []api.FargateProfileSelector{
				api.FargateProfileSelector{Namespace: "foo"},
			},
		},
		&api.FargateProfile{
			Name: "selecting-coredns",
			Selectors: []api.FargateProfileSelector{
				api.FargateProfileSelector{Namespace: "fooo"},
				api.FargateProfileSelector{Namespace: "kube-system"},
				api.FargateProfileSelector{Namespace: "baar"},
			},
		},
		&api.FargateProfile{
			Name: "bar",
			Selectors: []api.FargateProfileSelector{
				api.FargateProfileSelector{Namespace: "bar"},
			},
		},
	}

	profileNotSelectingCoreDNSBecauseOfLabels = []*api.FargateProfile{
		&api.FargateProfile{
			Name: "not-selecting-coredns-because-of-labels",
			Selectors: []api.FargateProfileSelector{
				api.FargateProfileSelector{
					Namespace: "kube-system",
					Labels: map[string]string{
						"foo": "bar",
					},
				},
			},
		},
	}

	profileNotSelectingCoreDNSBecauseOfNamespace = []*api.FargateProfile{
		&api.FargateProfile{
			Name: "not-selecting-coredns-because-of-namespace",
			Selectors: []api.FargateProfileSelector{
				api.FargateProfileSelector{
					Namespace: "default",
				},
			},
		},
	}
)

var _ = Describe("coredns", func() {
	Describe("IsSchedulableOnFargate", func() {
		It("should return true when a Fargate profile matches kube-system and doesn't have any label", func() {
			Expect(coredns.IsSchedulableOnFargate(minimumProfileSelectingCoreDNS)).To(BeTrue())
			Expect(coredns.IsSchedulableOnFargate(multipleProfilesWithOneSelectingCoreDNS)).To(BeTrue())
		})

		It("should return true when provided the default Fargate profile", func() {
			cfg := api.NewClusterConfig()
			cfg.SetDefaultFargateProfile()
			Expect(coredns.IsSchedulableOnFargate(cfg.FargateProfiles)).To(BeTrue())
		})

		It("should return false when a Fargate profile matches kube-system but has labels", func() {
			Expect(coredns.IsSchedulableOnFargate(profileNotSelectingCoreDNSBecauseOfLabels)).To(BeFalse())
		})

		It("should return false when a Fargate profile doesn't match kube-system", func() {
			Expect(coredns.IsSchedulableOnFargate(profileNotSelectingCoreDNSBecauseOfNamespace)).To(BeFalse())
		})

		It("should return false when not provided any Fargate profile", func() {
			Expect(coredns.IsSchedulableOnFargate([]*api.FargateProfile{})).To(BeFalse())
		})
	})

	Describe("ScheduleOnFargate", func() {
		It("should set the compute-type annotation should have been set to 'fargate'", func() {
			mockClientset := mockClientsetWith(deploymentAnnotatedWith("ec2"))
			// Given:
			deployment, err := mockClientset.ExtensionsV1beta1().Deployments(coredns.Namespace).Get(coredns.Name, metav1.GetOptions{})
			Expect(err).To(Not(HaveOccurred()))
			Expect(deployment.Spec.Template.Annotations).To(HaveKeyWithValue(coredns.ComputeTypeAnnotationKey, "ec2"))
			// When:
			err = coredns.ScheduleOnFargate(mockClientset)
			Expect(err).To(Not(HaveOccurred()))
			// Then:
			deployment, err = mockClientset.ExtensionsV1beta1().Deployments(coredns.Namespace).Get(coredns.Name, metav1.GetOptions{})
			Expect(err).To(Not(HaveOccurred()))
			Expect(deployment.Spec.Template.Annotations).To(HaveKeyWithValue(coredns.ComputeTypeAnnotationKey, "fargate"))
		})
	})
})

func mockClientsetWith(objects ...runtime.Object) kubeclient.Interface {
	return fake.NewSimpleClientset(objects...)
}

func deploymentAnnotatedWith(annotationValue string) *v1beta1.Deployment {
	return &v1beta1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: coredns.Namespace,
			Name:      coredns.Name,
		},
		Spec: v1beta1.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						coredns.ComputeTypeAnnotationKey: annotationValue,
					},
				},
			},
		},
	}
}
