package windows_test

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/windows"
	"k8s.io/client-go/kubernetes/fake"
)

type ipamEntry struct {
	existingConfigMapData map[string]string

	expectedConfigMapData map[string]string
}

var _ = DescribeTable("Windows IPAM", func(e ipamEntry) {
	var clientset *fake.Clientset
	if e.existingConfigMapData != nil {
		clientset = fake.NewSimpleClientset(&v1.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ConfigMap",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "amazon-vpc-cni",
				Namespace: "kube-system",
			},
			Data: e.existingConfigMapData,
		})
	} else {
		clientset = fake.NewSimpleClientset()
	}

	ipam := &windows.IPAM{
		Clientset: clientset,
	}
	ctx := context.Background()
	err := ipam.Enable(ctx)
	Expect(err).NotTo(HaveOccurred())

	cm, err := clientset.CoreV1().ConfigMaps("kube-system").Get(ctx, "amazon-vpc-cni", metav1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())
	Expect(cm.Data).To(Equal(e.expectedConfigMapData))

},
	Entry("VPC CNI ConfigMap is missing", ipamEntry{
		expectedConfigMapData: map[string]string{
			"enable-windows-ipam": "true",
		},
	}),

	Entry("VPC CNI ConfigMap has data", ipamEntry{
		existingConfigMapData: map[string]string{
			"VPC_CNI_1": "yes",
			"VPC_CNI_2": "no",
			"other":     "true",
		},
		expectedConfigMapData: map[string]string{
			"VPC_CNI_1":           "yes",
			"VPC_CNI_2":           "no",
			"other":               "true",
			"enable-windows-ipam": "true",
		},
	}),

	Entry("VPC CNI ConfigMap has Windows IPAM already enabled", ipamEntry{
		existingConfigMapData: map[string]string{
			"VPC_CNI_1":           "yes",
			"VPC_CNI_2":           "no",
			"enable-windows-ipam": "true",
		},
		expectedConfigMapData: map[string]string{
			"VPC_CNI_1":           "yes",
			"VPC_CNI_2":           "no",
			"enable-windows-ipam": "true",
		},
	}),

	Entry("VPC CNI ConfigMap has Windows IPAM explicitly disabled", ipamEntry{
		existingConfigMapData: map[string]string{
			"VPC_CNI_1":           "yes",
			"VPC_CNI_2":           "no",
			"enable-windows-ipam": "false",
		},
		expectedConfigMapData: map[string]string{
			"VPC_CNI_1":           "yes",
			"VPC_CNI_2":           "no",
			"enable-windows-ipam": "true",
		},
	}),
)
