package defaultaddons

import (
	"fmt"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	appsv1 "k8s.io/api/apps/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

const (
	// AWSNode is the name of the aws-node addon
	AWSNode = "aws-node"

	awsNodeImagePrefix = "602401143452.dkr.ecr."
	awsNodeImageSuffix = ".amazonaws.com/amazon-k8s-cni"
)

// UpdateAWSNode will update the `aws-node` add-on
func UpdateAWSNode(rawClient kubernetes.RawClientInterface, region string) error {
	_, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(AWSNode, metav1.GetOptions{})
	if err != nil {
		if apierrs.IsNotFound(err) {
			logger.Warning("%q was not found", AWSNode)
			return nil
		}
		return errors.Wrapf(err, "getting %q", AWSNode)
	}

	// if DaemonSets is present, go through our list of assets
	list, err := LoadAsset(AWSNode, "yaml")
	if err != nil {
		return err
	}

	for _, rawObj := range list.Items {
		resource, err := rawClient.NewRawResource(rawObj)
		if err != nil {
			return err
		}
		if resource.GVK.Kind == "DaemonSet" {
			image := &resource.Info.Object.(*appsv1.DaemonSet).Spec.Template.Spec.Containers[0].Image
			imageParts := strings.Split(*image, ":")

			if len(imageParts) != 2 {
				return fmt.Errorf("unexpected image format %q for %q", *image, KubeProxy)
			}

			if strings.HasPrefix(imageParts[0], awsNodeImagePrefix) &&
				strings.HasSuffix(imageParts[0], awsNodeImageSuffix) {
				*image = awsNodeImagePrefix + region + awsNodeImageSuffix + ":" + imageParts[1]
			}
		}

		status, err := resource.CreateOrReplace()
		if err != nil {
			return err
		}
		logger.Info(status)
	}

	logger.Info("%q is now up-to-date", AWSNode)
	return nil
}
