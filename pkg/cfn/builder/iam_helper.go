package builder

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/arn"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	cft "github.com/weaveworks/eksctl/pkg/cfn/template"
	gfn "github.com/weaveworks/goformation/v4/cloudformation"
	gfniam "github.com/weaveworks/goformation/v4/cloudformation/iam"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"
	"k8s.io/apimachinery/pkg/util/sets"
)

type cfnTemplate interface {
	attachAllowPolicy(name string, refRole *gfnt.Value, statements []cft.MapOfInterfaces)
	attachAllowPolicyDocument(name string, refRole *gfnt.Value, document api.InlineDocument)
	newResource(name string, resource gfn.Resource) *gfnt.Value
}

type managedPolicyForRole struct {
	name string
}

type customPolicyForRole struct {
	Name       string
	Statements []cft.MapOfInterfaces
}

func createWellKnownPolicies(wellKnownPolicies api.WellKnownPolicies) ([]managedPolicyForRole, []customPolicyForRole) {
	var managedPolicies []managedPolicyForRole
	var customPolicies []customPolicyForRole
	if wellKnownPolicies.ImageBuilder {
		managedPolicies = append(managedPolicies,
			managedPolicyForRole{name: iamPolicyAmazonEC2ContainerRegistryPowerUser},
		)
	}
	if wellKnownPolicies.AutoScaler {
		customPolicies = append(customPolicies,
			customPolicyForRole{Name: "PolicyAutoScaling", Statements: autoScalerStatements()},
		)
	}
	if wellKnownPolicies.AWSLoadBalancerController {
		customPolicies = append(customPolicies,
			customPolicyForRole{Name: "PolicyAWSLoadBalancerController", Statements: loadBalancerControllerStatements()},
		)
	}
	if wellKnownPolicies.ExternalDNS {
		customPolicies = append(customPolicies,
			[]customPolicyForRole{
				{Name: "PolicyExternalDNSChangeSet", Statements: changeSetStatements()},
				{Name: "PolicyExternalDNSHostedZones", Statements: externalDNSHostedZonesStatements()},
			}...,
		)
	}
	if wellKnownPolicies.CertManager {
		customPolicies = append(customPolicies,
			[]customPolicyForRole{
				{Name: "PolicyCertManagerChangeSet", Statements: changeSetStatements()},
				{Name: "PolicyCertManagerGetChange", Statements: certManagerGetChangeStatements()},
				{Name: "PolicyCertManagerHostedZones", Statements: certManagerHostedZonesStatements()},
			}...,
		)
	}
	if wellKnownPolicies.EBSCSIController {
		customPolicies = append(customPolicies,
			customPolicyForRole{Name: "PolicyEBSCSIController", Statements: ebsStatements()},
		)
	}
	if wellKnownPolicies.EFSCSIController {
		customPolicies = append(customPolicies,
			customPolicyForRole{Name: "PolicyEFSCSIController", Statements: efsCSIControllerStatements()},
		)
	}
	return managedPolicies, customPolicies
}

// createRole creates an IAM role with policies required for the worker nodes and addons
func createRole(cfnTemplate cfnTemplate, clusterIAMConfig *api.ClusterIAM, iamConfig *api.NodeGroupIAM, managed, forceAddCNIPolicy bool) error {
	managedPolicyARNs, err := makeManagedPolicies(clusterIAMConfig, iamConfig, managed, forceAddCNIPolicy)
	if err != nil {
		return err
	}
	role := gfniam.Role{
		Path:                     gfnt.NewString("/"),
		AssumeRolePolicyDocument: cft.MakeAssumeRolePolicyDocumentForServices(MakeServiceRef("EC2")),
		ManagedPolicyArns:        managedPolicyARNs,
	}

	if iamConfig.InstanceRoleName != "" {
		role.RoleName = gfnt.NewString(iamConfig.InstanceRoleName)
	}

	if iamConfig.InstanceRolePermissionsBoundary != "" {
		role.PermissionsBoundary = gfnt.NewString(iamConfig.InstanceRolePermissionsBoundary)
	}

	refIR := cfnTemplate.newResource(cfnIAMInstanceRoleName, &role)

	if iamConfig.AttachPolicy != nil {
		cfnTemplate.attachAllowPolicyDocument("Policy1", refIR, iamConfig.AttachPolicy)
	}

	if api.IsEnabled(iamConfig.WithAddonPolicies.AutoScaler) {
		cfnTemplate.attachAllowPolicy("PolicyAutoScaling", refIR, autoScalerStatements())
	}

	if api.IsEnabled(iamConfig.WithAddonPolicies.CertManager) {
		cfnTemplate.attachAllowPolicy("PolicyCertManagerChangeSet", refIR, changeSetStatements())
		cfnTemplate.attachAllowPolicy("PolicyCertManagerHostedZones", refIR, certManagerHostedZonesStatements())
		cfnTemplate.attachAllowPolicy("PolicyCertManagerGetChange", refIR, certManagerGetChangeStatements())
	}
	if api.IsEnabled(iamConfig.WithAddonPolicies.ExternalDNS) {
		cfnTemplate.attachAllowPolicy("PolicyExternalDNSChangeSet", refIR, changeSetStatements())
		cfnTemplate.attachAllowPolicy("PolicyExternalDNSHostedZones", refIR, externalDNSHostedZonesStatements())
	}

	if api.IsEnabled(iamConfig.WithAddonPolicies.AppMesh) {
		cfnTemplate.attachAllowPolicy("PolicyAppMesh", refIR, appMeshStatements("appmesh:*"))
	}

	if api.IsEnabled(iamConfig.WithAddonPolicies.AppMeshPreview) {
		cfnTemplate.attachAllowPolicy("PolicyAppMeshPreview", refIR, appMeshStatements("appmesh-preview:*"))
	}

	if api.IsEnabled(iamConfig.WithAddonPolicies.EBS) {
		cfnTemplate.attachAllowPolicy("PolicyEBS", refIR, ebsStatements())
	}

	if api.IsEnabled(iamConfig.WithAddonPolicies.FSX) {
		cfnTemplate.attachAllowPolicy("PolicyFSX", refIR, fsxStatements())
		cfnTemplate.attachAllowPolicy("PolicyServiceLinkRole", refIR, serviceLinkRoleStatements())
	}

	if api.IsEnabled(iamConfig.WithAddonPolicies.EFS) {
		cfnTemplate.attachAllowPolicy("PolicyEFS", refIR, efsStatements())
		cfnTemplate.attachAllowPolicy("PolicyEFSEC2", refIR, efsEc2Statements())
	}

	if api.IsEnabled(iamConfig.WithAddonPolicies.AWSLoadBalancerController) {
		cfnTemplate.attachAllowPolicy("PolicyAWSLoadBalancerController", refIR, loadBalancerControllerStatements())
	}

	if api.IsEnabled(iamConfig.WithAddonPolicies.XRay) {
		cfnTemplate.attachAllowPolicy("PolicyXRay", refIR, xRayStatements())
	}

	return nil
}

func makeManagedPolicies(iamCluster *api.ClusterIAM, iamConfig *api.NodeGroupIAM, managed, forceAddCNIPolicy bool) (*gfnt.Value, error) {
	managedPolicyNames := sets.NewString()
	if len(iamConfig.AttachPolicyARNs) == 0 {
		managedPolicyNames.Insert(iamDefaultNodePolicies...)
		if !api.IsEnabled(iamCluster.WithOIDC) || forceAddCNIPolicy {
			managedPolicyNames.Insert(iamPolicyAmazonEKSCNIPolicy)
		}
		if managed {
			// The Managed Nodegroup API requires this managed policy to be present, even though
			// AmazonEC2ContainerRegistryPowerUser (attached if imageBuilder is enabled) contains a superset of the
			// actions allowed by this managed policy
			managedPolicyNames.Insert(iamPolicyAmazonEC2ContainerRegistryReadOnly)
		}
		managedPolicyNames.Insert(iamPolicyAmazonSSMManagedInstanceCore)
	}

	if api.IsEnabled(iamConfig.WithAddonPolicies.ImageBuilder) {
		managedPolicyNames.Insert(iamPolicyAmazonEC2ContainerRegistryPowerUser)
	} else if !managed {
		// attach this policy even if `AttachPolicyARNs` is specified to preserve existing behaviour for unmanaged
		// nodegroups
		managedPolicyNames.Insert(iamPolicyAmazonEC2ContainerRegistryReadOnly)
	}

	if api.IsEnabled(iamConfig.WithAddonPolicies.CloudWatch) {
		managedPolicyNames.Insert(iamPolicyCloudWatchAgentServerPolicy)
	}

	for _, policyARN := range iamConfig.AttachPolicyARNs {
		parsedARN, err := arn.Parse(policyARN)
		if err != nil {
			return nil, err
		}
		start := strings.IndexRune(parsedARN.Resource, '/')
		if start == -1 || start+1 == len(parsedARN.Resource) {
			return nil, fmt.Errorf("failed to find ARN resource name: %s", parsedARN.Resource)
		}
		resourceName := parsedARN.Resource[start+1:]
		managedPolicyNames.Delete(resourceName)
	}

	return gfnt.NewSlice(append(
		makeStringSlice(iamConfig.AttachPolicyARNs...),
		makePolicyARNs(managedPolicyNames.List()...)...,
	)...), nil
}

// NormalizeARN returns the ARN with just the last element in the resource path preserved. If the
// input does not contain at least one forward-slash then the input is returned unmodified.
//
// When providing an existing instanceRoleARN that contains a path other than "/", nodes may
// fail to join the cluster as the AWS IAM Authenticator does not recognize such ARNs declared in
// the aws-auth ConfigMap.
//
// See: https://docs.aws.amazon.com/eks/latest/userguide/troubleshooting.html#troubleshoot-container-runtime-network
func NormalizeARN(arn string) string {
	parts := strings.Split(arn, "/")
	if len(parts) <= 1 {
		return arn
	}
	return fmt.Sprintf("%s/%s", parts[0], parts[len(parts)-1])
}

// AbstractRoleNameFromARN returns the role name from the ARN
func AbstractRoleNameFromARN(arn string) string {
	parts := strings.Split(arn, "/")
	if len(parts) <= 1 {
		return arn
	}
	return parts[len(parts)-1]
}
