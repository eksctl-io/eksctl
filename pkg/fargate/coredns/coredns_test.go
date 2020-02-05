package coredns_test

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/fargate/coredns"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/utils/names"
	"github.com/weaveworks/eksctl/pkg/utils/retry"
	v1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
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
		{
			Name: "min-selecting-coredns",
			Selectors: []api.FargateProfileSelector{
				{Namespace: "kube-system"},
			},
		},
	}

	multipleProfilesWithOneSelectingCoreDNS = []*api.FargateProfile{
		{
			Name: "foo",
			Selectors: []api.FargateProfileSelector{
				{Namespace: "foo"},
			},
		},
		{
			Name: "selecting-coredns",
			Selectors: []api.FargateProfileSelector{
				{Namespace: "fooo"},
				{Namespace: "kube-system"},
				{Namespace: "baar"},
			},
		},
		{
			Name: "bar",
			Selectors: []api.FargateProfileSelector{
				{Namespace: "bar"},
			},
		},
	}

	profileNotSelectingCoreDNSBecauseOfLabels = []*api.FargateProfile{
		{
			Name: "not-selecting-coredns-because-of-labels",
			Selectors: []api.FargateProfileSelector{
				{
					Namespace: "kube-system",
					Labels: map[string]string{
						"foo": "bar",
					},
				},
			},
		},
	}

	profileNotSelectingCoreDNSBecauseOfNamespace = []*api.FargateProfile{
		{
			Name: "not-selecting-coredns-because-of-namespace",
			Selectors: []api.FargateProfileSelector{
				{
					Namespace: "default",
				},
			},
		},
	}

	retryPolicy = &retry.ConstantBackoff{
		// Retry without waiting at all, in order to speed tests up.
		Time: 0, TimeUnit: time.Second, MaxRetries: 1,
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
			// Given:
			mockClientset := mockClientsetWith(deployment("ec2", 0, 2))
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

	Describe("WaitForScheduleOnFargate", func() {
		It("should wait for coredns to be scheduled on Fargate and return w/o any error", func() {
			// Given:
			mockClientset := mockClientsetWith(
				deployment("fargate", 2, 2), pod("fargate", v1.PodRunning), pod("fargate", v1.PodRunning),
			)
			// When:
			err := coredns.WaitForScheduleOnFargate(mockClientset, retryPolicy)
			// Then:
			Expect(err).To(Not(HaveOccurred()))
		})

		It("should time out if coredns cannot be scheduled within the allotted time", func() {
			failureCases := [][]runtime.Object{
				{deployment("ec2", 2, 2), pod("ec2", v1.PodRunning), pod("ec2", v1.PodRunning)},
				{deployment("ec2", 0, 2), pod("ec2", v1.PodPending), pod("ec2", v1.PodPending)},
				{deployment("fargate", 0, 2), pod("fargate", v1.PodPending), pod("fargate", v1.PodPending)},
				{deployment("fargate", 0, 2), pod("fargate", v1.PodFailed), pod("fargate", v1.PodFailed)},
				{deployment("fargate", 0, 2), pod("fargate", v1.PodPending), pod("fargate", v1.PodFailed)},
				{deployment("fargate", 1, 2), pod("fargate", v1.PodRunning), pod("fargate", v1.PodPending)},
				{deployment("fargate", 1, 2), pod("fargate", v1.PodRunning), pod("fargate", v1.PodFailed)},
			}
			for _, failureCase := range failureCases {
				// Given:
				mockClientset := mockClientsetWith(failureCase...)
				// When:
				err := coredns.WaitForScheduleOnFargate(mockClientset, retryPolicy)
				// Then:
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("timed out while waiting for \"coredns\" to be scheduled on Fargate"))
			}
		})
	})
})

func mockClientsetWith(objects ...runtime.Object) kubeclient.Interface {
	return fake.NewSimpleClientset(objects...)
}

func deployment(computeType string, numReady, numReplicas int32) *v1beta1.Deployment {
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
			Replicas: &numReplicas,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						coredns.ComputeTypeAnnotationKey: computeType,
					},
				},
			},
		},
		Status: v1beta1.DeploymentStatus{
			ReadyReplicas: numReady,
		},
	}
}

const chars = "abcdef0123456789"

func pod(computeType string, phase v1.PodPhase) *v1.Pod {
	pod := &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: coredns.Namespace,
			Name:      fmt.Sprintf("%s-%s-%s", coredns.Name, names.RandomName(10, chars), names.RandomName(5, chars)),
			Labels: map[string]string{
				"eks.amazonaws.com/component": coredns.Name,
			},
			Annotations: map[string]string{
				coredns.ComputeTypeAnnotationKey: computeType,
			},
		},
		Status: v1.PodStatus{
			Phase: phase,
		},
	}
	if pod.Status.Phase == v1.PodRunning {
		if computeType == "fargate" {
			pod.Spec.NodeName = "fargate-ip-192-168-xxx-yyy.ap-northeast-1.compute.internal"
		} else {
			pod.Spec.NodeName = "ip-192-168-23-122.ap-northeast-1.compute.internal"
		}
	}
	return pod
}
