package eks

import (
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	storage "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
)

// AddDefaultStorageClass adds the default EBS gp2 storage class to the cluster
func (c *ClusterProvider) AddDefaultStorageClass(clientSet *clientset.Clientset) error {

	rp := corev1.PersistentVolumeReclaimRetain

	scb := &storage.StorageClass{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StorageClass",
			APIVersion: "storage.k8s.io/v1",
		},
		Provisioner: "kubernetes.io/aws-ebs",
		ObjectMeta: metav1.ObjectMeta{
			Name: "gp2",
			Annotations: map[string]string{
				"storageclass.kubernetes.io/is-default-class": "true",
			},
		},
		Parameters: map[string]string{
			"type": "gp2",
		},
		ReclaimPolicy: &rp,
		MountOptions:  []string{"debug"},
	}
	logger.Debug("Creating a StorageClass as default")

	if _, err := clientSet.StorageV1().StorageClasses().Create(scb); err != nil {
		return errors.Wrap(err, "adding default StorageClass of gp2")
	}

	return nil
}
