package irsa

import (
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func SetMaybeCreateServiceAccountOrUpdateMetadata(f func(clientSet kubernetes.Interface, meta v1.ObjectMeta) error) {
	maybeCreateServiceAccountOrUpdateMetadata = f
}
