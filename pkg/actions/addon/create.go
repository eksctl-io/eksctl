package addon

import (
	"context"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"k8s.io/apimachinery/pkg/types"
)

const (
	kubeSystemNamespace = "kube-system"
	vpcCNIName          = "vpc-cni"
)

func (a *Manager) Create(addon *api.Addon) error {
	createAddonInput := &eks.CreateAddonInput{
		AddonName:    &addon.Name,
		AddonVersion: &addon.Version,
		ClusterName:  &a.clusterConfig.Metadata.Name,
		//ResolveConflicts: 		"enum":["OVERWRITE","NONE"]
	}

	if addon.Force {
		createAddonInput.ResolveConflicts = aws.String("overwrite")
		logger.Debug("setting resolve conflicts to overwrite")
	}

	logger.Debug("addon: %v", addon)
	namespace, serviceAccount := a.getKnownServiceAccountLocation(addon)

	if a.withOIDC {
		if addon.ServiceAccountRoleARN != "" {
			logger.Info("using provided ServiceAccountRoleARN %q", addon.ServiceAccountRoleARN)
			createAddonInput.ServiceAccountRoleArn = &addon.ServiceAccountRoleARN
		} else if addon.AttachPolicyARNs != nil && len(addon.AttachPolicyARNs) != 0 {
			logger.Info("creating role using provided policies ARNs")
			role, err := a.createRoleUsingAttachPolicyARNs(addon, namespace, serviceAccount)
			if err != nil {
				return err
			}
			createAddonInput.ServiceAccountRoleArn = &role
		} else if addon.AttachPolicy != nil {
			logger.Info("creating role using provided policies")
			role, err := a.createRoleUsingAttachPolicy(addon, namespace, serviceAccount)
			if err != nil {
				return err
			}
			createAddonInput.ServiceAccountRoleArn = &role
		} else {
			policies := a.getRecommendedPolicies(addon)
			if len(policies) != 0 {
				logger.Info("creating role using recommended policies")
				addon.AttachPolicyARNs = policies
				role, err := a.createRoleUsingAttachPolicyARNs(addon, namespace, serviceAccount)
				if err != nil {
					return err
				}
				createAddonInput.ServiceAccountRoleArn = &role
			} else {
				logger.Info("no recommended policies found, proceeding without any IAM")
			}
		}
	} else {
		//if any sort of policy is set or could be set, log a warning
		if addon.ServiceAccountRoleARN != "" ||
			(addon.AttachPolicyARNs != nil && len(addon.AttachPolicyARNs) != 0) ||
			addon.AttachPolicy != nil ||
			len(a.getRecommendedPolicies(addon)) != 0 {
			logger.Warning("OIDC is disabled but policies are required/specified for this addon. Users are responsible for attaching the policies to all nodegroup roles")
		}
	}

	if strings.ToLower(addon.Name) == vpcCNIName {
		logger.Debug("patching AWS node")
		err := a.patchAWSNodeSA()
		if err != nil {
			return err
		}

		err = a.patchAWSNodeDaemonSet()
		if err != nil {
			return err
		}
	}

	logger.Info("creating addon")
	output, err := a.clusterProvider.Provider.EKS().CreateAddon(createAddonInput)
	if err != nil {
		return errors.Wrapf(err, "failed to create addon %q", addon.Name)
	}

	if output != nil {
		logger.Debug("EKS Create Addon output: %s", output.String())
	}

	logger.Info("successfully created addon")
	return nil
}

func (a *Manager) patchAWSNodeSA() error {
	serviceaccounts := a.clientSet.CoreV1().ServiceAccounts("kube-system")
	sa, err := serviceaccounts.Get(context.TODO(), "aws-node", metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			logger.Debug("could not find aws-node SA, skipping patching")
			return nil
		}
		return err
	}
	var managerIndex = -1
	for i, managedFields := range sa.ManagedFields {
		if managedFields.Manager == "eksctl" {
			managerIndex = i
		}
	}
	if managerIndex == -1 {
		logger.Debug("no 'eksctl' managed field found")
		return nil
	}

	_, err = serviceaccounts.Patch(context.TODO(), "aws-node", types.JSONPatchType, []byte(fmt.Sprintf(`[{"op": "remove", "path": "/metadata/managedFields/%d"}]`, managerIndex)), metav1.PatchOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to patch sa")
	}

	return nil
}

func (a *Manager) patchAWSNodeDaemonSet() error {
	daemonsets := a.clientSet.AppsV1().DaemonSets(kubeSystemNamespace)
	sa, err := daemonsets.Get(context.TODO(), "aws-node", metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			logger.Debug("could not find aws-node daemon set, skipping patching")
			return nil
		}
		return err
	}
	var managerIndex = -1
	for i, managedFields := range sa.ManagedFields {
		if managedFields.Manager == "eksctl" {
			managerIndex = i
		}
	}
	if managerIndex == -1 {
		logger.Debug("no 'eksctl' managed field found")
		return nil
	}

	_, err = daemonsets.Patch(context.TODO(), "aws-node", types.JSONPatchType, []byte(fmt.Sprintf(`[{"op": "remove", "path": "/metadata/managedFields/%d"}]`, managerIndex)), metav1.PatchOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to patch daemon set")
	}

	return nil
}
func (a *Manager) getRecommendedPolicies(addon *api.Addon) []string {
	// API isn't case sensitive
	switch strings.ToLower(addon.Name) {
	case vpcCNIName:
		return []string{fmt.Sprintf("arn:%s:iam::aws:policy/%s", api.Partition(a.clusterConfig.Metadata.Region), api.IAMPolicyAmazonEKSCNIPolicy)}
	default:
		return []string{}
	}
}

func (a *Manager) getKnownServiceAccountLocation(addon *api.Addon) (string, string) {
	// API isn't case sensitive
	switch strings.ToLower(addon.Name) {
	case vpcCNIName:
		logger.Debug("found known service account location %s/%s", api.AWSNodeMeta.Namespace, api.AWSNodeMeta.Name)
		return api.AWSNodeMeta.Namespace, api.AWSNodeMeta.Name
	default:
		return "", ""
	}
}

func (a *Manager) createRoleUsingAttachPolicyARNs(addon *api.Addon, namespace, serviceAccount string) (string, error) {
	resourceSet := builder.NewIAMRoleResourceSetWithAttachPolicyARNs(addon.Name, namespace, serviceAccount, addon.AttachPolicyARNs, a.oidcManager)
	err := resourceSet.AddAllResources()
	if err != nil {
		return "", err
	}

	err = a.createStack(resourceSet, addon)
	if err != nil {
		return "", err
	}
	return resourceSet.OutputRole, nil
}

func (a *Manager) createRoleUsingAttachPolicy(addon *api.Addon, namespace, serviceAccount string) (string, error) {
	resourceSet := builder.NewIAMRoleResourceSetWithAttachPolicy(addon.Name, namespace, serviceAccount, addon.AttachPolicy, a.oidcManager)
	err := resourceSet.AddAllResources()
	if err != nil {
		return "", err
	}

	err = a.createStack(resourceSet, addon)
	if err != nil {
		return "", err
	}
	return resourceSet.OutputRole, nil
}

func (a *Manager) createStack(resourceSet builder.ResourceSet, addon *api.Addon) error {
	errChan := make(chan error)

	tags := map[string]string{
		api.AddonNameTag: addon.Name,
	}

	err := a.stackManager.CreateStack(a.makeAddonName(addon.Name), resourceSet, tags, nil, errChan)
	if err != nil {
		return err
	}

	return <-errChan
}
