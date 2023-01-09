package addon

import (
	"context"
	"fmt"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/kris-nova/logger"

	"k8s.io/apimachinery/pkg/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
)

const (
	kubeSystemNamespace = "kube-system"
)

func (a *Manager) Create(ctx context.Context, addon *api.Addon, waitTimeout time.Duration) error {
	// First check if the addon is already present as an EKS managed addon
	// in a state different from CREATE_FAILED, and if so, don't re-create
	var notFoundErr *ekstypes.ResourceNotFoundException
	summary, err := a.eksAPI.DescribeAddon(ctx, &eks.DescribeAddonInput{
		AddonName:   &addon.Name,
		ClusterName: &a.clusterConfig.Metadata.Name,
	})
	if err != nil && !errors.As(err, &notFoundErr) {
		return err
	}

	// if the addon already exists AND it is not in CREATE_FAILED state
	if err == nil && summary.Addon.Status != ekstypes.AddonStatusCreateFailed {
		logger.Info("Addon %s is already present in this cluster, as an EKS managed addon, and won't be re-created", addon.Name)
		return nil
	}

	version := addon.Version
	if version != "" {
		var err error
		version, err = a.getLatestMatchingVersion(ctx, addon)
		if err != nil {
			return fmt.Errorf("failed to fetch version %s for addon %s: %w", version, addon.Name, err)
		}
	}
	var configurationValues *string
	if addon.ConfigurationValues != "" {
		configurationValues = &addon.ConfigurationValues
	}
	createAddonInput := &eks.CreateAddonInput{
		AddonName:           &addon.Name,
		AddonVersion:        &version,
		ClusterName:         &a.clusterConfig.Metadata.Name,
		ResolveConflicts:    addon.ResolveConflicts,
		ConfigurationValues: configurationValues,
	}

	if addon.Force {
		createAddonInput.ResolveConflicts = ekstypes.ResolveConflictsOverwrite
	} else {
		addonName := strings.ToLower(addon.Name)
		if addonName == "coredns" || addonName == "kube-proxy" || addonName == "vpc-cni" {
			logger.Info("when creating an addon to replace an existing application, e.g. CoreDNS, kube-proxy & VPC-CNI the --force flag will ensure the currently deployed configuration is replaced")
		}
	}

	logger.Debug("resolve conflicts set to %s", createAddonInput.ResolveConflicts)
	logger.Debug("addon: %v", addon)
	namespace, serviceAccount := a.getKnownServiceAccountLocation(addon)

	if len(addon.Tags) > 0 {
		createAddonInput.Tags = addon.Tags
	}
	if a.withOIDC {
		if addon.ServiceAccountRoleARN != "" {
			logger.Info("using provided ServiceAccountRoleARN %q", addon.ServiceAccountRoleARN)
			createAddonInput.ServiceAccountRoleArn = &addon.ServiceAccountRoleARN
		} else if hasPoliciesSet(addon) {
			outputRole, err := a.createRole(ctx, addon, namespace, serviceAccount)
			if err != nil {
				return err
			}
			createAddonInput.ServiceAccountRoleArn = &outputRole
		} else {
			policyDocument, policyARNs, wellKnownPolicies := a.getRecommendedPolicies(addon)
			if len(policyARNs) != 0 || policyDocument != nil || wellKnownPolicies != nil {
				logger.Info("creating role using recommended policies")
				addon.AttachPolicyARNs = policyARNs
				addon.AttachPolicy = policyDocument
				resourceSet := builder.NewIAMRoleResourceSetWithAttachPolicy(addon.Name, namespace, serviceAccount, addon.PermissionsBoundary, addon.AttachPolicy, a.oidcManager)
				if len(policyARNs) != 0 {
					resourceSet = builder.NewIAMRoleResourceSetWithAttachPolicyARNs(addon.Name, namespace, serviceAccount, addon.PermissionsBoundary, addon.AttachPolicyARNs, a.oidcManager)
				}
				if wellKnownPolicies != nil {
					addon.WellKnownPolicies = *wellKnownPolicies
					resourceSet = builder.NewIAMRoleResourceSetWithWellKnownPolicies(addon.Name, namespace, serviceAccount, addon.PermissionsBoundary, addon.WellKnownPolicies, a.oidcManager)
				}
				if err := resourceSet.AddAllResources(); err != nil {
					return err
				}
				err := a.createStack(ctx, resourceSet, addon)
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
		policyDocument, policyARNs, wellKnownPolicies := a.getRecommendedPolicies(addon)
		if addon.ServiceAccountRoleARN != "" || hasPoliciesSet(addon) || len(policyARNs) != 0 || policyDocument != nil || wellKnownPolicies != nil {
			logger.Warning("OIDC is disabled but policies are required/specified for this addon. Users are responsible for attaching the policies to all nodegroup roles")
		}
	}

	if addon.CanonicalName() == api.VPCCNIAddon {
		logger.Debug("patching AWS node")
		err := a.patchAWSNodeSA(ctx)
		if err != nil {
			return err
		}

		err = a.patchAWSNodeDaemonSet(ctx)
		if err != nil {
			return err
		}
	}

	logger.Info("creating addon")
	output, err := a.eksAPI.CreateAddon(ctx, createAddonInput)
	if err != nil {
		return errors.Wrapf(err, "failed to create addon %q", addon.Name)
	}

	if output != nil {
		logger.Debug("EKS Create Addon output: %s", *output.Addon)
	}

	if waitTimeout > 0 {
		return a.waitForAddonToBeActive(ctx, addon, waitTimeout)
	}
	logger.Info("successfully created addon")
	return nil
}

func (a *Manager) patchAWSNodeSA(ctx context.Context) error {
	serviceAccounts := a.clientSet.CoreV1().ServiceAccounts("kube-system")
	sa, err := serviceAccounts.Get(ctx, "aws-node", metav1.GetOptions{})
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

	_, err = serviceAccounts.Patch(ctx, "aws-node", types.JSONPatchType, []byte(fmt.Sprintf(`[{"op": "remove", "path": "/metadata/managedFields/%d"}]`, managerIndex)), metav1.PatchOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to patch sa")
	}

	return nil
}

func (a *Manager) patchAWSNodeDaemonSet(ctx context.Context) error {
	daemonSets := a.clientSet.AppsV1().DaemonSets(kubeSystemNamespace)
	sa, err := daemonSets.Get(ctx, "aws-node", metav1.GetOptions{})
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

	if a.clusterConfig.IPv6Enabled() {
		// Add ENABLE_IPV6 = true and ENABLE_PREFIX_DELEGATION = true
		_, err = daemonSets.Patch(ctx, "aws-node", types.StrategicMergePatchType, []byte(`{
	"spec": {
		"template": {
			"spec": {
				"containers": [{
					"env": [{
						"name": "ENABLE_IPV6",
						"value": "true"
					}, {
						"name": "ENABLE_PREFIX_DELEGATION",
						"value": "true"
					}],
					"name": "aws-node"
				}]
			}
		}
	}
}
`), metav1.PatchOptions{})
		if err != nil {
			return errors.Wrap(err, "failed to patch daemon set")
		}
		// update the daemonset so the next patch can work.
		_, err = daemonSets.Get(ctx, "aws-node", metav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				logger.Debug("could not find aws-node daemon set, skipping patching")
				return nil
			}
			return err
		}
	}

	_, err = daemonSets.Patch(ctx, "aws-node", types.JSONPatchType, []byte(fmt.Sprintf(`[{"op": "remove", "path": "/metadata/managedFields/%d"}]`, managerIndex)), metav1.PatchOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to patch daemon set")
	}

	return nil
}

func (a *Manager) getRecommendedPolicies(addon *api.Addon) (api.InlineDocument, []string, *api.WellKnownPolicies) {
	// API isn't case-sensitive
	switch addon.CanonicalName() {
	case api.VPCCNIAddon:
		if a.clusterConfig.IPv6Enabled() {
			return makeIPv6VPCCNIPolicyDocument(api.Partition(a.clusterConfig.Metadata.Region)), nil, nil
		}
		return nil, []string{fmt.Sprintf("arn:%s:iam::aws:policy/%s", api.Partition(a.clusterConfig.Metadata.Region), api.IAMPolicyAmazonEKSCNIPolicy)}, nil
	case api.AWSEBSCSIDriverAddon:
		return nil, nil, &api.WellKnownPolicies{
			EBSCSIController: true,
		}
	default:
		return nil, nil, nil
	}
}

func (a *Manager) getKnownServiceAccountLocation(addon *api.Addon) (string, string) {
	// API isn't case sensitive
	switch addon.CanonicalName() {
	case api.VPCCNIAddon:
		logger.Debug("found known service account location %s/%s", api.AWSNodeMeta.Namespace, api.AWSNodeMeta.Name)
		return api.AWSNodeMeta.Namespace, api.AWSNodeMeta.Name
	default:
		return "", ""
	}
}

func hasPoliciesSet(addon *api.Addon) bool {
	return len(addon.AttachPolicyARNs) != 0 || addon.WellKnownPolicies.HasPolicy() || addon.AttachPolicy != nil
}

func (a *Manager) createRole(ctx context.Context, addon *api.Addon, namespace, serviceAccount string) (string, error) {
	resourceSet, err := a.createRoleResourceSet(addon, namespace, serviceAccount)

	if err != nil {
		return "", err
	}

	err = a.createStack(ctx, resourceSet, addon)
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

func (a *Manager) createStack(ctx context.Context, resourceSet builder.ResourceSetReader, addon *api.Addon) error {
	errChan := make(chan error)

	tags := map[string]string{
		api.AddonNameTag: addon.Name,
	}

	err := a.stackManager.CreateStack(ctx, a.makeAddonName(addon.Name), resourceSet, tags, nil, errChan)
	if err != nil {
		return err
	}

	return <-errChan
}

func makeIPv6VPCCNIPolicyDocument(partition string) map[string]interface{} {
	return map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Effect": "Allow",
				"Action": []string{
					"ec2:AssignIpv6Addresses",
					"ec2:DescribeInstances",
					"ec2:DescribeTags",
					"ec2:DescribeNetworkInterfaces",
					"ec2:DescribeInstanceTypes",
				},
				"Resource": "*",
			},
			{
				"Effect": "Allow",
				"Action": []string{
					"ec2:CreateTags",
				},
				"Resource": fmt.Sprintf("arn:%s:ec2:*:*:network-interface/*", partition),
			},
		},
	}
}
