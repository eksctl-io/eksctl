package defaultaddons

import (
	"fmt"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	"github.com/weaveworks/eksctl/pkg/printers"

	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	// KubeProxy is the name of the kube-proxy addon
	KubeProxy = "kube-proxy"
)

// UpdateKubeProxyImageTag updates image tag for kube-system:damoneset/kube-proxy based to match controlPlaneVersion
func UpdateKubeProxyImageTag(clientSet kubernetes.Interface, controlPlaneVersion string, dryRun bool) error {
	printer := printers.NewJSONPrinter()

	d, err := clientSet.Apps().DaemonSets(metav1.NamespaceSystem).Get(KubeProxy, metav1.GetOptions{})
	if err != nil {
		if apierrs.IsNotFound(err) {
			logger.Warning("%q was not found", KubeProxy)
			return nil
		}
		return errors.Wrapf(err, "getting %s", KubeProxy)
	}
	if numContainers := len(d.Spec.Template.Spec.Containers); !(numContainers >= 1) {
		return fmt.Errorf("%s has %d containers, expected at least 1", KubeProxy, numContainers)
	}

	if err := printer.LogObj(logger.Debug, KubeProxy+" [current] = \\\n%s\n", d); err != nil {
		return err
	}

	image := &d.Spec.Template.Spec.Containers[0].Image
	imageParts := strings.Split(*image, ":")

	if len(imageParts) != 2 {
		return fmt.Errorf("unexpected image format %q for %q", *image, KubeProxy)
	}

	desiredTag := "v" + controlPlaneVersion

	if imageParts[1] == desiredTag {
		logger.Debug("imageParts = %v, desiredTag = %s", imageParts, desiredTag)
		logger.Info("%q is already up-to-date", KubeProxy)
		return nil
	}

	if dryRun {
		logger.Critical("%q is not up-to-date", KubeProxy)
		return nil
	}

	imageParts[1] = desiredTag
	*image = strings.Join(imageParts, ":")

	if err := printer.LogObj(logger.Debug, KubeProxy+" [updated] = \\\n%s\n", d); err != nil {
		return err
	}
	if _, err := clientSet.Apps().DaemonSets(metav1.NamespaceSystem).Update(d); err != nil {
		return err
	}

	logger.Info("%q is now up-to-date", KubeProxy)
	return nil
}
