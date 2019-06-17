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
func UpdateAWSNode(rawClient kubernetes.RawClientInterface, region, controlPlaneVersion string, plan bool) (bool, error) {
	_, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(AWSNode, metav1.GetOptions{})
	if err != nil {
		if apierrs.IsNotFound(err) {
			logger.Warning("%q was not found", AWSNode)
			return false, nil
		}
		return false, errors.Wrapf(err, "getting %q", AWSNode)
	}

	// if DaemonSets is present, go through our list of assets
	assetName := AWSNode
	if strings.HasPrefix(controlPlaneVersion, "1.10.") {
		assetName += "-1.10"
	}
	list, err := LoadAsset(assetName, "yaml")
	if err != nil {
		return false, err
	}

	for _, rawObj := range list.Items {
		resource, err := rawClient.NewRawResource(rawObj)
		if err != nil {
			return false, err
		}
		if resource.GVK.Kind == "DaemonSet" {
			image := &resource.Info.Object.(*appsv1.DaemonSet).Spec.Template.Spec.Containers[0].Image
			imageParts := strings.Split(*image, ":")

			if len(imageParts) != 2 {
				return false, fmt.Errorf("unexpected image format %q for %q", *image, KubeProxy)
			}

			if strings.HasPrefix(imageParts[0], awsNodeImagePrefix) &&
				strings.HasSuffix(imageParts[0], awsNodeImageSuffix) {
				*image = awsNodeImagePrefix + region + awsNodeImageSuffix + ":" + imageParts[1]
			}
		}

		if resource.GVK.Kind == "CustomResourceDefinition" && plan {
			// eniconfigs.crd.k8s.amazonaws.com CRD is only partially defined in the
			// manifest, and causes a range of issue in plan mode, we can skip it
			logger.Info(resource.LogAction(plan, "replaced"))
			continue
		}

		status, err := resource.CreateOrReplace(plan)
		if err != nil {
			return false, err
		}
		logger.Info(status)
	}

	if plan {
		logger.Critical("(plan) %q is not up-to-date", AWSNode)
		return true, nil
	}

	logger.Info("%q is now up-to-date", AWSNode)
	return false, nil
}
