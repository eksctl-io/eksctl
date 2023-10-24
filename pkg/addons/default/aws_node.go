package defaultaddons

import (
	"context"
	"fmt"
	"strings"

	"github.com/blang/semver"
	"github.com/kris-nova/logger"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/weaveworks/eksctl/pkg/addons"
	"github.com/weaveworks/eksctl/pkg/kubernetes"

	// For go:embed
	_ "embed"
)

const (
	// AWSNode is the name of the aws-node addon
	AWSNode = "aws-node"

	awsNodeImageFormatPrefix     = "%s.dkr.ecr.%s.%s/amazon-k8s-cni"
	awsNodeInitImageFormatPrefix = "%s.dkr.ecr.%s.%s/amazon-k8s-cni-init"
)

//go:embed assets/aws-node.yaml
var latestAWSNodeYaml []byte

// DoesAWSNodeSupportMultiArch makes sure awsnode supports ARM nodes
func DoesAWSNodeSupportMultiArch(ctx context.Context, clientSet kubernetes.Interface) (bool, error) {
	clusterDaemonSet, err := getAWSNode(ctx, clientSet)
	if err != nil {
		return false, err
	}
	if clusterDaemonSet == nil {
		return true, nil
	}

	minVersion := semver.Version{
		Major: 1,
		Minor: 6,
		Patch: 3,
	}

	clusterTag, err := addons.ImageTag(clusterDaemonSet.Spec.Template.Spec.Containers[0].Image)
	if err != nil {
		return false, err
	}
	clusterVersion, err := semver.ParseTolerant(clusterTag)
	if err != nil {
		return false, err
	}
	clusterSemverVersion := semver.Version{
		Major: clusterVersion.Major,
		Minor: clusterVersion.Minor,
		Patch: clusterVersion.Patch,
	}

	if clusterSemverVersion.GT(minVersion) ||
		(clusterSemverVersion.EQ(minVersion) && clusterVersion.String() == "1.6.3-eksbuild.1") {
		return true, nil
	}

	return false, nil
}

// UpdateAWSNode will update the `aws-node` add-on and returns true
// if an update is available.
func UpdateAWSNode(ctx context.Context, input AddonInput, plan bool) (bool, error) {
	clusterDaemonSet, err := getAWSNode(ctx, input.RawClient.ClientSet())
	if err != nil {
		return false, err
	}
	if clusterDaemonSet == nil {
		return false, nil
	}

	resourceList, err := newList(latestAWSNodeYaml)
	if err != nil {
		return false, err
	}

	tagMismatch := true
	for _, rawObj := range resourceList.Items {
		resource, err := input.RawClient.NewRawResource(rawObj.Object)
		if err != nil {
			return false, err
		}

		switch resource.GVK.Kind {
		case "DaemonSet":
			daemonSet, ok := resource.Info.Object.(*appsv1.DaemonSet)
			if !ok {
				return false, fmt.Errorf("expected type %T; got %T", &appsv1.Deployment{}, resource.Info.Object)
			}
			container := &daemonSet.Spec.Template.Spec.Containers[0]
			initContainer := &daemonSet.Spec.Template.Spec.InitContainers[0]
			imageParts := strings.Split(container.Image, ":")
			if len(imageParts) != 2 {
				return false, fmt.Errorf("invalid container image: %s", container.Image)
			}

			container.Image = awsNodeImageFormatPrefix + ":" + imageParts[1]
			initContainer.Image = awsNodeInitImageFormatPrefix + ":" + imageParts[1]
			if err := addons.UseRegionalImage(&daemonSet.Spec.Template, input.Region); err != nil {
				return false, err
			}

			containerTagMismatch, err := addons.ImageTagsDiffer(
				container.Image,
				clusterDaemonSet.Spec.Template.Spec.Containers[0].Image,
			)
			if err != nil {
				return false, err
			}

			initContainerTagMismatch := true // Will be true by default if the init containers don't exist
			if len(clusterDaemonSet.Spec.Template.Spec.InitContainers) > 0 {
				initContainerTagMismatch, err = addons.ImageTagsDiffer(
					initContainer.Image,
					clusterDaemonSet.Spec.Template.Spec.InitContainers[0].Image,
				)
				if err != nil {
					return false, err
				}
			}

			tagMismatch = containerTagMismatch || initContainerTagMismatch

		case "CustomResourceDefinition":
			if plan {
				// eniconfigs.crd.k8s.amazonaws.com CRD is only partially defined in the
				// manifest, and causes a range of issue in plan mode, we can skip it
				logger.Info(resource.LogAction(plan, "replaced"))
				continue
			}
		case "ServiceAccount":
			// Leave service account if it exists
			// to avoid overwriting annotations
			_, exists, err := resource.Get()
			if err != nil {
				return false, err
			}
			if exists {
				logger.Info(resource.LogAction(plan, "skipped existing"))
				continue
			}
		}

		status, err := resource.CreateOrReplace(plan)
		if err != nil {
			return false, err
		}
		logger.Info(status)
	}

	if plan {
		if tagMismatch {
			logger.Critical("(plan) %q is not up-to-date", AWSNode)
			return true, nil
		}
		logger.Info("(plan) %q is already up-to-date", AWSNode)
		return false, nil
	}

	logger.Info("%q is now up-to-date", AWSNode)
	return false, nil
}

func getAWSNode(ctx context.Context, clientSet kubernetes.Interface) (*appsv1.DaemonSet, error) {
	d, err := clientSet.AppsV1().DaemonSets(metav1.NamespaceSystem).Get(ctx, AWSNode, metav1.GetOptions{})
	return makeGetError(d, err, AWSNode)
}
