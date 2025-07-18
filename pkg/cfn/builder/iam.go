package builder

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/kris-nova/logger"

	gfniam "github.com/weaveworks/eksctl/pkg/goformation/cloudformation/iam"
	gfnrolesanywhere "github.com/weaveworks/eksctl/pkg/goformation/cloudformation/rolesanywhere"
	gfnt "github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	cft "github.com/weaveworks/eksctl/pkg/cfn/template"
	"github.com/weaveworks/eksctl/pkg/iam"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
)

const (
	iamPolicyAmazonEKSClusterPolicy             = "AmazonEKSClusterPolicy"
	iamPolicyAmazonEKSVPCResourceController     = "AmazonEKSVPCResourceController"
	iamPolicyAmazonEKSLocalOutpostClusterPolicy = "AmazonEKSLocalOutpostClusterPolicy"

	iamPolicyAmazonEKSWorkerNodePolicy           = "AmazonEKSWorkerNodePolicy"
	iamPolicyAmazonEKSCNIPolicy                  = "AmazonEKS_CNI_Policy"
	iamPolicyAmazonEC2ContainerRegistryPowerUser = "AmazonEC2ContainerRegistryPowerUser"
	iamPolicyAmazonEC2ContainerRegistryPullOnly  = "AmazonEC2ContainerRegistryPullOnly"
	iamPolicyCloudWatchAgentServerPolicy         = "CloudWatchAgentServerPolicy"
	iamPolicyAmazonSSMManagedInstanceCore        = "AmazonSSMManagedInstanceCore"
	iamPolicyAmazonEBSCSIDriverPolicy            = "service-role/AmazonEBSCSIDriverPolicy"

	iamPolicyAmazonEKSFargatePodExecutionRolePolicy = "AmazonEKSFargatePodExecutionRolePolicy"
)

const (
	cfnIAMInstanceRoleName    = "NodeInstanceRole"
	cfnIAMInstanceProfileName = "NodeInstanceProfile"
)

const (
	TrustAnchor      = "TrustAnchor"
	AnywhereProfile  = "AnywhereProfile"
	IntermediateRole = "IntermediateRole"
	IRARole          = "HybridNodesIRARole"
	SSMRole          = "HybridNodesSSMRole"
)

var (
	iamDefaultNodePolicies = []string{
		iamPolicyAmazonEKSWorkerNodePolicy,
	}
	eksDescribeClusterPolicy = gfniam.Role_Policy{
		PolicyName: gfnt.NewString("EKSDescribeCluster"),
		PolicyDocument: cft.MakePolicyDocument(cft.MapOfInterfaces{
			"Effect": "Allow",
			"Action": []string{
				"eks:DescribeCluster",
			},
			"Resource": "*",
		}),
	}
	ssmRolePolicy = gfniam.Role_Policy{
		PolicyName: gfnt.NewString("SSMRolePolicy"),
		PolicyDocument: cft.MakePolicyDocument(cft.MapOfInterfaces{
			"Effect": "Allow",
			"Action": []string{
				"ssm:DeregisterManagedInstance",
				"ssm:DescribeInstanceInformation",
			},
			"Resource": "*",
		}),
	}
)

func (c *resourceSet) attachAllowPolicy(name string, refRole *gfnt.Value, statements []cft.MapOfInterfaces) {
	c.newResource(name, &gfniam.Policy{
		PolicyName:     makeName(name),
		Roles:          gfnt.NewSlice(refRole),
		PolicyDocument: cft.MakePolicyDocument(statements...),
	})
}

func (c *resourceSet) attachAllowPolicyDocument(name string, refRole *gfnt.Value, document api.InlineDocument) {
	c.newResource(name, &gfniam.Policy{
		PolicyName:     makeName(name),
		Roles:          gfnt.NewSlice(refRole),
		PolicyDocument: document,
	})
}

// WithIAM states, if IAM roles will be created or not
func (c *ClusterResourceSet) WithIAM() bool {
	return c.rs.withIAM
}

// WithNamedIAM states, if specifically named IAM roles will be created or not
func (c *ClusterResourceSet) WithNamedIAM() bool {
	return c.rs.withNamedIAM
}

func (c *ClusterResourceSet) addResourcesForServiceRole() {
	c.rs.withIAM = true

	var role *gfniam.Role
	if c.spec.IsControlPlaneOnOutposts() {
		role = &gfniam.Role{
			AssumeRolePolicyDocument: cft.MakeAssumeRolePolicyDocumentForServices(
				MakeServiceRef("EC2"),
			),
			ManagedPolicyArns: gfnt.NewSlice(makePolicyARNs(iamPolicyAmazonEKSLocalOutpostClusterPolicy)...),
		}
	} else {
		managedPolicyARNs := []string{iamPolicyAmazonEKSClusterPolicy}
		if !api.IsDisabled(c.spec.IAM.VPCResourceControllerPolicy) {
			managedPolicyARNs = append(managedPolicyARNs, iamPolicyAmazonEKSVPCResourceController)
		}
		if c.spec.IsAutoModeEnabled() {
			managedPolicyARNs = append(managedPolicyARNs, AutoModeIAMPolicies...)
		}
		role = &gfniam.Role{
			AssumeRolePolicyDocument: cft.MakeAssumeRolePolicyDocumentWithAction(
				"sts:TagSession",
				MakeServiceRef("EKS"),
			),
			ManagedPolicyArns: gfnt.NewSlice(makePolicyARNs(managedPolicyARNs...)...),
		}
		if c.spec.IsCustomEksEndpoint() {
			role.AssumeRolePolicyDocument = createBetaAssumeRolePolicy()
		}
	}

	if api.IsSetAndNonEmptyString(c.spec.IAM.ServiceRolePermissionsBoundary) {
		role.PermissionsBoundary = gfnt.NewString(*c.spec.IAM.ServiceRolePermissionsBoundary)
	}
	c.newResource("ServiceRole", role)

	c.rs.defineOutputFromAtt(outputs.ClusterServiceRoleARN, "ServiceRole", "Arn", true, func(v string) error {
		c.spec.IAM.ServiceRoleARN = &v
		return nil
	})
}

func (c *ClusterResourceSet) addIAMRolesAnywhere() {
	trustAnchor := &gfnrolesanywhere.TrustAnchor{
		Enabled: gfnt.NewBoolean(true),
		Name:    makeName("CA"),
		Source: &gfnrolesanywhere.TrustAnchor_Source{
			SourceType: gfnt.NewString("CERTIFICATE_BUNDLE"),
			SourceData: &gfnrolesanywhere.TrustAnchor_SourceData{},
		},
	}
	if c.spec.RemoteNetworkConfig.IAM.CABundleCert != nil {
		trustAnchor.Source.SourceData.X509CertificateData = gfnt.NewString(*c.spec.RemoteNetworkConfig.IAM.CABundleCert)
	}
	anywhereProfile := &gfnrolesanywhere.Profile{
		Enabled: gfnt.NewBoolean(true),
		Name:    makeName("remote-nodes"),
		RoleArns: gfnt.NewSlice(
			gfnt.MakeFnGetAttString(IRARole, "Arn"),
		),
		AcceptRoleSessionName:      gfnt.NewBoolean(true),
		AWSCloudFormationDependsOn: []string{IRARole},
	}
	iraRole := &gfniam.Role{
		AssumeRolePolicyDocument: cft.MakeAssumeRolePolicyDocumentForServicesWithConditionsAndActions(
			cft.MapOfInterfaces{
				"ArnEquals": cft.MapOfInterfaces{
					"aws:SourceArn": gfnt.MakeFnGetAttString(TrustAnchor, "TrustAnchorArn"),
				},
			},
			[]string{
				"sts:TagSession",
				"sts:SetSourceIdentity",
			},
			MakeServiceRef("IRA"),
		),
		Policies: []gfniam.Role_Policy{
			eksDescribeClusterPolicy,
		},
		ManagedPolicyArns: gfnt.NewSlice(makePolicyARNs([]string{
			iamPolicyAmazonEC2ContainerRegistryPullOnly,
			iamPolicyAmazonSSMManagedInstanceCore,
		}...)...),
		AWSCloudFormationDependsOn: []string{TrustAnchor},
	}

	c.newResource(TrustAnchor, trustAnchor)
	c.newResource(AnywhereProfile, anywhereProfile)
	c.newResource(IRARole, iraRole)

	c.rs.defineOutputFromAtt(outputs.RemoteNodesTrustAnchorARN, TrustAnchor, "TrustAnchorArn", true, func(v string) error {
		return nil
	})
	c.rs.defineOutputFromAtt(outputs.RemoteNodesAnywhereProfileARN, AnywhereProfile, "ProfileArn", true, func(v string) error {
		return nil
	})
	c.rs.defineOutputFromAtt(outputs.RemoteNodesRoleARN, IRARole, "Arn", true, func(v string) error {
		c.spec.RemoteNetworkConfig.IAM.RoleARN = &v
		return nil
	})
}

func (c *ClusterResourceSet) addSSM() {
	role := &gfniam.Role{
		AssumeRolePolicyDocument: cft.MakeAssumeRolePolicyDocumentForServices(
			MakeServiceRef("SSM"),
		),
		Policies: []gfniam.Role_Policy{
			ssmRolePolicy,
			eksDescribeClusterPolicy,
		},
		ManagedPolicyArns: gfnt.NewSlice(makePolicyARNs([]string{
			iamPolicyAmazonSSMManagedInstanceCore,
			iamPolicyAmazonEC2ContainerRegistryPullOnly,
		}...)...),
	}
	c.newResource(SSMRole, role)
	c.rs.defineOutputFromAtt(outputs.RemoteNodesRoleARN, SSMRole, "Arn", true, func(v string) error {
		c.spec.RemoteNetworkConfig.IAM.RoleARN = &v
		return nil
	})
}

func (c *ClusterResourceSet) addResourcesForRemoteNodesRole() {
	c.rs.withIAM = true
	switch strings.ToLower(*c.spec.RemoteNetworkConfig.IAM.Provider) {
	case api.SSMProvider:
		c.addSSM()
	case api.IRAProvider:
		c.addIAMRolesAnywhere()
	default:
		// Validations should ensure this is never reached
	}
}

func (c *ClusterResourceSet) addResourcesForIAM() {
	c.rs.withIAM = false
	c.rs.withNamedIAM = false

	if !api.IsSetAndNonEmptyString(c.spec.IAM.ServiceRoleARN) {
		c.addResourcesForServiceRole()
	} else {
		c.rs.defineOutputWithoutCollector(outputs.ClusterServiceRoleARN, c.spec.IAM.ServiceRoleARN, true)
	}

	if c.spec.HasRemoteNetworkingConfigured() {
		if !api.IsSetAndNonEmptyString(c.spec.RemoteNetworkConfig.IAM.RoleARN) {
			c.addResourcesForRemoteNodesRole()
		} else {
			c.rs.defineOutputWithoutCollector(outputs.RemoteNodesRoleARN, c.spec.RemoteNetworkConfig.IAM.RoleARN, true)
		}
	}
}

// WithIAM states, if IAM roles will be created or not
func (n *NodeGroupResourceSet) WithIAM() bool {
	return n.rs.withIAM
}

// WithNamedIAM states, if specifically named IAM roles will be created or not
func (n *NodeGroupResourceSet) WithNamedIAM() bool {
	return n.rs.withNamedIAM
}

func (n *NodeGroupResourceSet) addResourcesForIAM(ctx context.Context) error {
	nodeGroupIAM := n.options.NodeGroup.IAM
	if nodeGroupIAM.InstanceProfileARN != "" {
		n.rs.withIAM = false
		n.rs.withNamedIAM = false

		// if instance profile is given, as well as the role, we simply use both and export the role
		// (which is needed in order to authorise the nodegroup)
		n.instanceProfileARN = gfnt.NewString(nodeGroupIAM.InstanceProfileARN)
		if nodeGroupIAM.InstanceRoleARN != "" {
			n.rs.defineOutputWithoutCollector(outputs.NodeGroupInstanceProfileARN, nodeGroupIAM.InstanceProfileARN, true)
			n.rs.defineOutputWithoutCollector(outputs.NodeGroupInstanceRoleARN, nodeGroupIAM.InstanceRoleARN, true)
			return nil
		}
		// if instance role is not given, export profile and use the getter to call importer function
		n.rs.defineOutput(outputs.NodeGroupInstanceProfileARN, nodeGroupIAM.InstanceProfileARN, true, func(v string) error {
			return iam.ImportInstanceRoleFromProfileARN(ctx, n.iamAPI, n.options.NodeGroup, v)
		})

		return nil
	}

	n.rs.withIAM = true

	if nodeGroupIAM.InstanceRoleARN != "" {
		roleARN := NormalizeARN(nodeGroupIAM.InstanceRoleARN)

		// if role is set, but profile isn't - create profile
		n.newResource(cfnIAMInstanceProfileName, &gfniam.InstanceProfile{
			Path:  gfnt.NewString("/"),
			Roles: gfnt.NewStringSlice(AbstractRoleNameFromARN(roleARN)),
		})
		n.instanceProfileARN = gfnt.MakeFnGetAttString(cfnIAMInstanceProfileName, "Arn")
		n.rs.defineOutputFromAtt(outputs.NodeGroupInstanceProfileARN, cfnIAMInstanceProfileName, "Arn", true, func(v string) error {
			nodeGroupIAM.InstanceProfileARN = v
			return nil
		})
		n.rs.defineOutputWithoutCollector(outputs.NodeGroupInstanceRoleARN, roleARN, true)
		return nil
	}

	// if neither role nor profile is given - create both

	if nodeGroupIAM.InstanceRoleName != "" {
		// setting role name requires additional capabilities
		n.rs.withNamedIAM = true
	}

	if err := createRole(n.rs, n.options.ClusterConfig.IAM, nodeGroupIAM, false, n.options.ForceAddCNIPolicy); err != nil {
		return err
	}

	n.newResource(cfnIAMInstanceProfileName, &gfniam.InstanceProfile{
		Path:  gfnt.NewString("/"),
		Roles: gfnt.NewSlice(gfnt.MakeRef(cfnIAMInstanceRoleName)),
	})
	n.instanceProfileARN = gfnt.MakeFnGetAttString(cfnIAMInstanceProfileName, "Arn")

	n.rs.defineOutputFromAtt(outputs.NodeGroupInstanceProfileARN, cfnIAMInstanceProfileName, "Arn", true, func(v string) error {
		nodeGroupIAM.InstanceProfileARN = v
		return nil
	})
	n.rs.defineOutputFromAtt(outputs.NodeGroupInstanceRoleARN, cfnIAMInstanceRoleName, "Arn", true, func(v string) error {
		nodeGroupIAM.InstanceRoleARN = v
		return nil
	})
	return nil
}

func NewIAMRoleResourceSetForServiceAccount(spec *api.ClusterIAMServiceAccount, oidc *iamoidc.OpenIDConnectManager) *IAMRoleResourceSet {
	return &IAMRoleResourceSet{
		template:            cft.NewTemplate(),
		attachPolicy:        spec.AttachPolicy,
		attachPolicyARNs:    spec.AttachPolicyARNs,
		serviceAccount:      spec.Name,
		namespace:           spec.Namespace,
		wellKnownPolicies:   spec.WellKnownPolicies,
		roleName:            spec.RoleName,
		permissionsBoundary: spec.PermissionsBoundary,
		description: fmt.Sprintf(
			"IAM role for serviceaccount %q %s",
			spec.NameString(),
			templateDescriptionSuffix,
		),
		oidc: oidc,
		roleNameCollector: func(v string) error {
			spec.Status = &api.ClusterIAMServiceAccountStatus{
				RoleARN: &v,
			}
			return nil
		},
	}
}

func NewIAMRoleResourceSetForPodIdentity(spec *api.PodIdentityAssociation) *IAMRoleResourceSet {
	return &IAMRoleResourceSet{
		template:            cft.NewTemplate(),
		attachPolicy:        spec.PermissionPolicy,
		attachPolicyARNs:    spec.PermissionPolicyARNs,
		serviceAccount:      spec.ServiceAccountName,
		namespace:           spec.Namespace,
		wellKnownPolicies:   spec.WellKnownPolicies,
		roleName:            spec.RoleName,
		permissionsBoundary: spec.PermissionsBoundaryARN,
		description: fmt.Sprintf(
			"IAM role for pod identity association %s",
			templateDescriptionSuffix,
		),
		roleNameCollector: func(v string) error {
			spec.RoleARN = v
			return nil
		},
	}
}

// IAMRoleResourceSet holds IAM Role stack build-time information
type IAMRoleResourceSet struct {
	template            *cft.Template
	oidc                *iamoidc.OpenIDConnectManager
	outputs             *outputs.CollectorSet
	roleName            string
	wellKnownPolicies   api.WellKnownPolicies
	attachPolicyARNs    []string
	attachPolicy        api.InlineDocument
	trustStatements     []api.IAMStatement
	roleNameCollector   func(string) error
	OutputRole          string
	serviceAccount      string
	namespace           string
	permissionsBoundary string
	description         string
}

// NewIAMRoleResourceSetWithAttachPolicyARNs builds IAM Role stack from the give spec
func NewIAMRoleResourceSetWithAttachPolicyARNs(name, namespace, serviceAccount, permissionsBoundary string, attachPolicyARNs []string, oidc *iamoidc.OpenIDConnectManager) *IAMRoleResourceSet {
	return newIAMRoleResourceSet(name, namespace, serviceAccount, permissionsBoundary, nil, attachPolicyARNs, api.WellKnownPolicies{}, oidc)
}

// NewIAMRoleResourceSetWithAttachPolicy builds IAM Role stack from the give spec
func NewIAMRoleResourceSetWithAttachPolicy(name, namespace, serviceAccount, permissionsBoundary string, attachPolicy api.InlineDocument, oidc *iamoidc.OpenIDConnectManager) *IAMRoleResourceSet {
	return newIAMRoleResourceSet(name, namespace, serviceAccount, permissionsBoundary, attachPolicy, nil, api.WellKnownPolicies{}, oidc)
}

// NewIAMRoleResourceSetWithAttachPolicyARNs builds IAM Role stack from the give spec
func NewIAMRoleResourceSetWithWellKnownPolicies(name, namespace, serviceAccount, permissionsBoundary string, wellKnownPolicies api.WellKnownPolicies, oidc *iamoidc.OpenIDConnectManager) *IAMRoleResourceSet {
	return newIAMRoleResourceSet(name, namespace, serviceAccount, permissionsBoundary, nil, nil, wellKnownPolicies, oidc)
}

// NewIAMRoleResourceSetForServiceAccount builds IAM Role stack from the give spec
func newIAMRoleResourceSet(name, namespace, serviceAccount, permissionsBoundary string, attachPolicy api.InlineDocument, attachPolicyARNs []string, wellKnownPolicies api.WellKnownPolicies, oidc *iamoidc.OpenIDConnectManager) *IAMRoleResourceSet {
	rs := &IAMRoleResourceSet{
		template:            cft.NewTemplate(),
		attachPolicyARNs:    attachPolicyARNs,
		attachPolicy:        attachPolicy,
		oidc:                oidc,
		serviceAccount:      serviceAccount,
		namespace:           namespace,
		permissionsBoundary: permissionsBoundary,
		description: fmt.Sprintf(
			"IAM role for %q %s",
			name,
			templateDescriptionSuffix,
		),
		wellKnownPolicies: wellKnownPolicies,
	}

	rs.roleNameCollector = func(v string) error {
		rs.OutputRole = v
		return nil
	}
	return rs
}

// WithIAM returns true
func (*IAMRoleResourceSet) WithIAM() bool { return true }

// WithNamedIAM returns false
func (rs *IAMRoleResourceSet) WithNamedIAM() bool { return rs.roleName != "" }

// AddAllResources adds all resources for the stack
func (rs *IAMRoleResourceSet) AddAllResources() error {
	rs.template.Description = rs.description

	role := &cft.IAMRole{
		AssumeRolePolicyDocument: rs.makeAssumeRolePolicyDocument(),
		PermissionsBoundary:      rs.permissionsBoundary,
		RoleName:                 rs.roleName,
	}

	for _, arn := range rs.attachPolicyARNs {
		role.ManagedPolicyArns = append(role.ManagedPolicyArns, arn)
	}

	managedPolicies, customPolicies := createWellKnownPolicies(rs.wellKnownPolicies)

	for _, p := range managedPolicies {
		role.ManagedPolicyArns = append(role.ManagedPolicyArns, makePolicyARN(p.name))
	}

	roleRef := rs.template.NewResource(outputs.IAMServiceAccountRoleName, role)

	for _, p := range customPolicies {
		doc := cft.MakePolicyDocument(p.Statements...)
		rs.template.AttachPolicy(p.Name, roleRef, doc)
	}

	rs.template.Outputs[outputs.IAMServiceAccountRoleName] = cft.Output{
		Value: cft.MakeFnGetAttString("Role1.Arn"),
	}
	rs.outputs = outputs.NewCollectorSet(map[string]outputs.Collector{
		outputs.IAMServiceAccountRoleName: rs.roleNameCollector,
	})

	if len(rs.attachPolicy) != 0 {
		rs.template.AttachPolicy("Policy1", roleRef, rs.attachPolicy)
	}

	return nil
}

func (rs *IAMRoleResourceSet) makeAssumeRolePolicyDocument() cft.MapOfInterfaces {
	if len(rs.trustStatements) > 0 {
		return cft.MakePolicyDocument(toMapOfInterfaces(rs.trustStatements)...)
	}
	if rs.oidc == nil {
		return cft.MakeAssumeRolePolicyDocumentForPodIdentity()
	}
	if rs.serviceAccount != "" && rs.namespace != "" {
		logger.Debug("service account location provided: %s/%s, adding sub condition", api.AWSNodeMeta.Namespace, api.AWSNodeMeta.Name)
		return rs.oidc.MakeAssumeRolePolicyDocumentWithServiceAccountConditions(rs.namespace, rs.serviceAccount)
	}
	return rs.oidc.MakeAssumeRolePolicyDocument()
}

func toMapOfInterfaces(old []api.IAMStatement) []cft.MapOfInterfaces {
	new := []cft.MapOfInterfaces{}
	for _, s := range old {
		new = append(new, s.ToMapOfInterfaces())
	}
	return new
}

// RenderJSON will render iamserviceaccount stack as JSON
func (rs *IAMRoleResourceSet) RenderJSON() ([]byte, error) {
	return rs.template.RenderJSON()
}

// GetAllOutputs will get all outputs from iamserviceaccount stack
func (rs *IAMRoleResourceSet) GetAllOutputs(stack types.Stack) error {
	return rs.outputs.MustCollect(stack)
}
