package addon

import (
	"context"
	"fmt"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/kris-nova/logger"

	"k8s.io/apimachinery/pkg/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
)

var (
	updateAddonRecommended = func(supportsPodIDs bool) string {
		path := "`addon.AttachPolicyARNs`, `addon.AttachPolicy` or `addon.WellKnownPolicies`"
		if supportsPodIDs {
			path = "`addon.PodIdentityAssociations`"
		}
		return fmt.Sprintf("add all recommended policies to the config file, under %s, and run `eksctl update addon`", path)
	}
	iamPermissionsRecommended = func(addonName string, supportsPodIDs, shouldUpdateAddon bool) string {
		method := "IRSA"
		if supportsPodIDs {
			method = "pod identity associations"
		}
		commandSuggestion := "run `eksctl utils migrate-to-pod-identity`"
		if shouldUpdateAddon {
			commandSuggestion = updateAddonRecommended(supportsPodIDs)
		}
		return fmt.Sprintf("the recommended way to provide IAM permissions for %q addon is via %s; after addon creation is completed, %s", addonName, method, commandSuggestion)
	}
	IRSADeprecatedWarning = func(addonName string) string {
		return fmt.Sprintf("IRSA has been deprecated; %s", iamPermissionsRecommended(addonName, true, false))
	}
	OIDCDisabledWarning = func(addonName string, supportsPodIDs, isIRSASetExplicitly bool) string {
		irsaUsedMessage := fmt.Sprintf("recommended policies were found for %q addon", addonName)
		if isIRSASetExplicitly {
			irsaUsedMessage = fmt.Sprintf("IRSA config is set for %q addon", addonName)
		}
		suggestion := "users are responsible for attaching the policies to all nodegroup roles"
		if supportsPodIDs {
			suggestion = iamPermissionsRecommended(addonName, true, true)
		}
		return fmt.Sprintf("%s, but since OIDC is disabled on the cluster, eksctl cannot configure the requested permissions; %s", irsaUsedMessage, suggestion)
	}
	IAMPermissionsRequiredWarning = func(addonName string, supportsPodIDs bool) string {
		suggestion := iamPermissionsRecommended(addonName, false, true)
		if supportsPodIDs {
			suggestion = iamPermissionsRecommended(addonName, true, true)
		}
		return fmt.Sprintf("IAM permissions are required for %q addon; %s", addonName, suggestion)
	}
	IAMPermissionsNotRequiredWarning = func(addonName string) string {
		return fmt.Sprintf("IAM permissions are not required for %q addon; any IRSA configuration or pod identity associations will be ignored", addonName)
	}
)

const (
	kubeSystemNamespace   = "kube-system"
	awsNodeServiceAccount = "aws-node"
)

func (a *Manager) Create(ctx context.Context, addon *api.Addon, iamRoleCreator IAMRoleCreator, waitTimeout time.Duration) error {
	// check if the addon is already present as an EKS managed addon
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
		logger.Info("addon %s is already present on the cluster, as an EKS managed addon, skipping creation", addon.Name)
		return nil
	}

	version, requiresIAMPermissions, err := a.getLatestMatchingVersion(ctx, addon)
	addon.Version = version
	if err != nil {
		return fmt.Errorf("failed to fetch version %s for addon %s: %w", version, addon.Name, err)
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

	if len(addon.Tags) > 0 {
		createAddonInput.Tags = addon.Tags
	}

	podIDConfig, supportsPodIDs, err := a.getRecommendedPoliciesForPodID(ctx, addon)
	if err != nil {
		return err
	}

	if requiresIAMPermissions {
		switch {
		case addon.HasPodIDsSet():
			if !supportsPodIDs {
				return fmt.Errorf("%q addon does not support pod identity associations; use IRSA instead", addon.Name)
			}
			logger.Info("pod identity associations were specified for %q addon; will use those to provide required IAM permissions, any IRSA settings will be ignored", addon.Name)
			for _, pia := range *addon.PodIdentityAssociations {
				roleARN, err := iamRoleCreator.Create(ctx, &pia, addon.Name)
				if err != nil {
					return err
				}
				createAddonInput.PodIdentityAssociations = append(createAddonInput.PodIdentityAssociations, ekstypes.AddonPodIdentityAssociations{
					RoleArn:        &roleARN,
					ServiceAccount: &pia.ServiceAccountName,
				})
			}

		case a.clusterConfig.IAM.AutoCreatePodIdentityAssociations && supportsPodIDs:
			logger.Info("\"iam.AutoCreatePodIdentityAssociations\" is set to true; will use recommended policies for %q addon, any IRSA settings will be ignored", addon.Name)

			if addon.CanonicalName() == api.VPCCNIAddon && a.clusterConfig.IPv6Enabled() {
				roleARN, err := iamRoleCreator.Create(ctx, &api.PodIdentityAssociation{
					ServiceAccountName: awsNodeServiceAccount,
					PermissionPolicy:   makeIPv6VPCCNIPolicyDocument(api.Partitions.ForRegion(a.clusterConfig.Metadata.Region)),
				}, addon.Name)
				if err != nil {
					return err
				}
				createAddonInput.PodIdentityAssociations = append(createAddonInput.PodIdentityAssociations, ekstypes.AddonPodIdentityAssociations{
					RoleArn:        &roleARN,
					ServiceAccount: aws.String(awsNodeServiceAccount),
				})
				break
			}

			for _, config := range podIDConfig {
				roleARN, err := iamRoleCreator.Create(ctx, &api.PodIdentityAssociation{
					ServiceAccountName:   *config.ServiceAccount,
					PermissionPolicyARNs: config.RecommendedManagedPolicies,
				}, addon.Name)
				if err != nil {
					return err
				}
				createAddonInput.PodIdentityAssociations = append(createAddonInput.PodIdentityAssociations, ekstypes.AddonPodIdentityAssociations{
					RoleArn:        &roleARN,
					ServiceAccount: config.ServiceAccount,
				})
			}

		case addon.HasIRSASet():
			if !a.withOIDC {
				logger.Warning(OIDCDisabledWarning(addon.Name, supportsPodIDs,
					/* isIRSASetExplicitly */ addon.ServiceAccountRoleARN != "" || addon.HasIRSAPoliciesSet()))
				break
			}

			if supportsPodIDs {
				logger.Warning(IRSADeprecatedWarning(addon.Name))
			}

			if addon.ServiceAccountRoleARN != "" {
				logger.Info("using provided ServiceAccountRoleARN %q", addon.ServiceAccountRoleARN)
				createAddonInput.ServiceAccountRoleArn = &addon.ServiceAccountRoleARN
				break
			}

			if !addon.HasIRSAPoliciesSet() {
				a.setRecommendedPoliciesForIRSA(addon)
				logger.Info("creating role using recommended policies for %q addon", addon.Name)
			} else {
				logger.Info("creating role using provided policies for %q addon", addon.Name)
			}

			namespace, serviceAccount := a.getKnownServiceAccountLocation(addon)
			roleARN, err := a.createRoleForIRSA(ctx, addon, namespace, serviceAccount)
			if err != nil {
				return err
			}
			createAddonInput.ServiceAccountRoleArn = &roleARN

		default:
			logger.Warning(IAMPermissionsRequiredWarning(addon.Name, supportsPodIDs))
		}

	} else if addon.HasPodIDsSet() || addon.HasIRSASet() {
		logger.Warning(IAMPermissionsNotRequiredWarning(addon.Name))
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
		var resourceInUse *ekstypes.ResourceInUseException
		if errors.As(err, &resourceInUse) {
			defer func() {
				deleteAddonIAMTasks, err := NewRemover(a.stackManager).DeleteAddonIAMTasksFiltered(ctx, addon.Name, false)
				if err != nil {
					logger.Warning("failed to cleanup IAM role stacks: %w; please remove any remaining stacks manually", err)
					return
				}
				if err := runAllTasks(deleteAddonIAMTasks); err != nil {
					logger.Warning("failed to cleanup IAM role stacks: %w; please remove any remaining stacks manually", err)
				}
			}()
			var addonServiceAccounts []string
			for _, config := range podIDConfig {
				addonServiceAccounts = append(addonServiceAccounts, fmt.Sprintf("%q", *config.ServiceAccount))
			}
			return fmt.Errorf("creating addon: one or more service accounts corresponding to %q addon is already associated with a different IAM role; please delete all pre-existing pod identity associations corresponding to %s service account(s) in the addon's namespace, then re-try creating the addon", addon.Name, strings.Join(addonServiceAccounts, ","))
		}
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
	clientSet, err := a.createClientSet()
	if err != nil {
		return err
	}
	serviceAccounts := clientSet.CoreV1().ServiceAccounts("kube-system")
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
	clientSet, err := a.createClientSet()
	if err != nil {
		return err
	}
	daemonSets := clientSet.AppsV1().DaemonSets(kubeSystemNamespace)
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

func (a *Manager) getKnownServiceAccountLocation(addon *api.Addon) (string, string) {
	// API isn't case-sensitive.
	switch addon.CanonicalName() {
	case api.VPCCNIAddon:
		logger.Debug("found known service account location %s/%s", api.AWSNodeMeta.Namespace, api.AWSNodeMeta.Name)
		return api.AWSNodeMeta.Namespace, api.AWSNodeMeta.Name
	default:
		return "", ""
	}
}

func (a *Manager) getRecommendedPoliciesForPodID(ctx context.Context, addon *api.Addon) ([]ekstypes.AddonPodIdentityConfiguration, bool, error) {
	output, err := a.eksAPI.DescribeAddonConfiguration(ctx, &eks.DescribeAddonConfigurationInput{
		AddonName:    &addon.Name,
		AddonVersion: &addon.Version,
	})
	if err != nil {
		return nil, false, fmt.Errorf("describing configuration for addon %s: %w", addon.Name, err)
	}
	return output.PodIdentityConfiguration, len(output.PodIdentityConfiguration) != 0, nil
}

func (a *Manager) setRecommendedPoliciesForIRSA(addon *api.Addon) {
	switch addon.CanonicalName() {
	case api.VPCCNIAddon:
		if a.clusterConfig.IPv6Enabled() {
			addon.AttachPolicy = makeIPv6VPCCNIPolicyDocument(api.Partitions.ForRegion(a.clusterConfig.Metadata.Region))
		}
		addon.AttachPolicyARNs = append(addon.AttachPolicyARNs, fmt.Sprintf("arn:%s:iam::aws:policy/%s", api.Partitions.ForRegion(a.clusterConfig.Metadata.Region), api.IAMPolicyAmazonEKSCNIPolicy))
	case api.AWSEBSCSIDriverAddon:
		addon.WellKnownPolicies = api.WellKnownPolicies{
			EBSCSIController: true,
		}
	case api.AWSEFSCSIDriverAddon:
		addon.WellKnownPolicies = api.WellKnownPolicies{
			EFSCSIController: true,
		}
	default:
		return
	}
}

func (a *Manager) createRoleForIRSA(ctx context.Context, addon *api.Addon, namespace, serviceAccount string) (string, error) {
	resourceSet, err := a.createRoleResourceSet(addon, namespace, serviceAccount)
	if err != nil {
		return "", err
	}
	if err := a.createStack(ctx, resourceSet, addon.Name,
		a.makeAddonIRSAName(addon.Name)); err != nil {
		return "", err
	}
	return resourceSet.OutputRole, nil
}

func (a *Manager) createRoleResourceSet(addon *api.Addon, namespace, serviceAccount string) (*builder.IAMRoleResourceSet, error) {
	var resourceSet *builder.IAMRoleResourceSet
	if len(addon.AttachPolicyARNs) != 0 {
		resourceSet = builder.NewIAMRoleResourceSetWithAttachPolicyARNs(addon.Name, namespace, serviceAccount, addon.PermissionsBoundary, addon.AttachPolicyARNs, a.oidcManager)
	} else if addon.WellKnownPolicies.HasPolicy() {
		resourceSet = builder.NewIAMRoleResourceSetWithWellKnownPolicies(addon.Name, namespace, serviceAccount, addon.PermissionsBoundary, addon.WellKnownPolicies, a.oidcManager)
	} else {
		resourceSet = builder.NewIAMRoleResourceSetWithAttachPolicy(addon.Name, namespace, serviceAccount, addon.PermissionsBoundary, addon.AttachPolicy, a.oidcManager)
	}
	return resourceSet, resourceSet.AddAllResources()
}

func (a *Manager) createStack(ctx context.Context, resourceSet builder.ResourceSetReader, addonName, stackName string) error {
	errChan := make(chan error)

	tags := map[string]string{
		api.AddonNameTag: addonName,
	}

	err := a.stackManager.CreateStack(ctx, stackName, resourceSet, tags, nil, errChan)
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
