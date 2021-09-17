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

func (a *Manager) Create(addon *api.Addon, wait bool) error {
	version := addon.Version
	if version != "" {
		var err error
		version, err = a.getLatestMatchingVersion(addon)
		if err != nil {
			return fmt.Errorf("failed to fetch version %s for addon %s: %w", version, addon.Name, err)
		}
	}
	createAddonInput := &eks.CreateAddonInput{
		AddonName:    &addon.Name,
		AddonVersion: &version,
		ClusterName:  &a.clusterConfig.Metadata.Name,
		//ResolveConflicts: 		"enum":["OVERWRITE","NONE"]
	}

	if addon.Force {
		createAddonInput.ResolveConflicts = aws.String("overwrite")
		logger.Debug("setting resolve conflicts to overwrite")
	} else {
		addonName := strings.ToLower(addon.Name)
		if addonName == "coredns" || addonName == "kube-proxy" || addonName == "vpc-cni" {
			logger.Info("when creating an addon to replace an existing application, e.g. CoreDNS, kube-proxy & VPC-CNI the --force flag will ensure the currently deployed configuration is replaced")
		}
	}

	logger.Debug("addon: %v", addon)
	namespace, serviceAccount := a.getKnownServiceAccountLocation(addon)

	if len(addon.Tags) > 0 {
		createAddonInput.Tags = aws.StringMap(addon.Tags)
	}
	if a.withOIDC {
		if addon.ServiceAccountRoleARN != "" {
			logger.Info("using provided ServiceAccountRoleARN %q", addon.ServiceAccountRoleARN)
			createAddonInput.ServiceAccountRoleArn = &addon.ServiceAccountRoleARN
		} else if hasPoliciesSet(addon) {
			outputRole, err := a.createRole(addon, namespace, serviceAccount)
			if err != nil {
				return err
			}
			createAddonInput.ServiceAccountRoleArn = &outputRole
		} else {
			policies := a.getRecommendedPolicies(addon)
			if len(policies) != 0 {
				logger.Info("creating role using recommended policies")
				addon.AttachPolicyARNs = policies
				resourceSet := builder.NewIAMRoleResourceSetWithAttachPolicyARNs(addon.Name, namespace, serviceAccount, addon.PermissionsBoundary, addon.AttachPolicyARNs, a.oidcManager)
				if err := resourceSet.AddAllResources(); err != nil {
					return err
				}
				err := a.createStack(resourceSet, addon)
				if err != nil {
					return err
				}
				createAddonInput.ServiceAccountRoleArn = &resourceSet.OutputRole
			} else {
				logger.Info("no recommended policies found, proceeding without any IAM")
			}
		}
	} else {
		//if any sort of policy is set or could be set, log a warning
		if addon.ServiceAccountRoleARN != "" || hasPoliciesSet(addon) || len(a.getRecommendedPolicies(addon)) != 0 {
			logger.Warning("OIDC is disabled but policies are required/specified for this addon. Users are responsible for attaching the policies to all nodegroup roles")
		}
	}

	if addon.CanonicalName() == vpcCNIName {
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
	output, err := a.eksAPI.CreateAddon(createAddonInput)
	if err != nil {
		return errors.Wrapf(err, "failed to create addon %q", addon.Name)
	}

	if output != nil {
		logger.Debug("EKS Create Addon output: %s", output.String())
	}

	if wait {
		return a.waitForAddonToBeActive(addon)
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
	switch addon.CanonicalName() {
	case vpcCNIName:
		return []string{fmt.Sprintf("arn:%s:iam::aws:policy/%s", api.Partition(a.clusterConfig.Metadata.Region), api.IAMPolicyAmazonEKSCNIPolicy)}
	default:
		return []string{}
	}
}

func (a *Manager) getKnownServiceAccountLocation(addon *api.Addon) (string, string) {
	// API isn't case sensitive
	switch addon.CanonicalName() {
	case vpcCNIName:
		logger.Debug("found known service account location %s/%s", api.AWSNodeMeta.Namespace, api.AWSNodeMeta.Name)
		return api.AWSNodeMeta.Namespace, api.AWSNodeMeta.Name
	default:
		return "", ""
	}
}

func hasPoliciesSet(addon *api.Addon) bool {
	return len(addon.AttachPolicyARNs) != 0 || addon.WellKnownPolicies.HasPolicy() || addon.AttachPolicy != nil
}

func (a *Manager) createRole(addon *api.Addon, namespace, serviceAccount string) (string, error) {
	resourceSet, err := a.createRoleResourceSet(addon, namespace, serviceAccount)

	if err != nil {
		return "", err
	}

	err = a.createStack(resourceSet, addon)
	if err != nil {
		return "", err
	}
	return resourceSet.OutputRole, nil
}

func (a *Manager) createRoleResourceSet(addon *api.Addon, namespace, serviceAccount string) (*builder.IAMRoleResourceSet, error) {
	var resourceSet *builder.IAMRoleResourceSet
	if len(addon.AttachPolicyARNs) != 0 {
		logger.Info("creating role using provided policies ARNs")
		resourceSet = builder.NewIAMRoleResourceSetWithAttachPolicyARNs(addon.Name, namespace, serviceAccount, addon.PermissionsBoundary, addon.AttachPolicyARNs, a.oidcManager)
	} else if addon.WellKnownPolicies.HasPolicy() {
		logger.Info("creating role using provided well known policies")
		resourceSet = builder.NewIAMRoleResourceSetWithWellKnownPolicies(addon.Name, namespace, serviceAccount, addon.PermissionsBoundary, addon.WellKnownPolicies, a.oidcManager)
	} else {
		logger.Info("creating role using provided policies")
		resourceSet = builder.NewIAMRoleResourceSetWithAttachPolicy(addon.Name, namespace, serviceAccount, addon.PermissionsBoundary, addon.AttachPolicy, a.oidcManager)
	}
	return resourceSet, resourceSet.AddAllResources()
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
