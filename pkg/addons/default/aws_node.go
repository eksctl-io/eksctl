package defaultaddons

import (
	"fmt"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/addons"
	appsv1 "k8s.io/api/apps/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

const (
	// AWSNode is the name of the aws-node addon
	AWSNode = "aws-node"

	awsNodeImageFormatPrefix = "%s.dkr.ecr.%s.%s/amazon-k8s-cni"
)

// UpdateAWSNode will update the `aws-node` add-on
func UpdateAWSNode(rawClient kubernetes.RawClientInterface, region string, plan bool) (bool, error) {
	_, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(AWSNode, metav1.GetOptions{})
	if err != nil {
		if apierrs.IsNotFound(err) {
			logger.Warning("%q was not found", AWSNode)
			return false, nil
		}
		return false, errors.Wrapf(err, "getting %q", AWSNode)
	}

	// if DaemonSets is present, go through our list of assets
	list, err := LoadAsset(AWSNode, "yaml")
	if err != nil {
		return false, err
	}

	for _, rawObj := range list.Items {
		resource, err := rawClient.NewRawResource(rawObj.Object)
		if err != nil {
			return false, err
		}
		if resource.GVK.Kind == "DaemonSet" {
			daemonSet, ok := resource.Info.Object.(*appsv1.DaemonSet)
			if !ok {
				return false, fmt.Errorf("expected type %T; got %T", &appsv1.Deployment{}, daemonSet)
			}
			container := &daemonSet.Spec.Template.Spec.Containers[0]
			imageParts := strings.Split(container.Image, ":")
			if len(imageParts) != 2 {
				return false, fmt.Errorf("invalid container image: %s", container.Image)
			}

			container.Image = awsNodeImageFormatPrefix + ":" + imageParts[1]
			if err := addons.UseRegionalImage(&daemonSet.Spec.Template, region); err != nil {
				return false, err
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
