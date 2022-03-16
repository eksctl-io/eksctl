package kubernetes

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func GetNodegroupKubernetesVersion(nodes v1.NodeInterface, ngName string) (string, error) {
	n, err := nodes.List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", api.NodeGroupNameLabel, ngName),
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to list nodes")
	} else if len(n.Items) == 0 {
		return "", nil
	}

	v := n.Items[0].Status.NodeInfo.KubeletVersion
	if strings.IndexRune(v, '-') > 0 {
		v = v[:strings.IndexRune(v, '-')]
	}
	if v[0] == 'v' {
		v = strings.TrimPrefix(v, "v")
	}

	return v, nil
}
