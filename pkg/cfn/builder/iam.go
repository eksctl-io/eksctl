package builder

import (
	"fmt"

	"github.com/kris-nova/logger"

	"github.com/weaveworks/eksctl/pkg/iam"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	gfniam "github.com/weaveworks/goformation/v4/cloudformation/iam"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	cft "github.com/weaveworks/eksctl/pkg/cfn/template"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
)

const (
	iamPolicyAmazonEKSClusterPolicy         = "AmazonEKSClusterPolicy"
	iamPolicyAmazonEKSVPCResourceController = "AmazonEKSVPCResourceController"

	iamPolicyAmazonEKSWorkerNodePolicy           = "AmazonEKSWorkerNodePolicy"
	iamPolicyAmazonEKSCNIPolicy                  = "AmazonEKS_CNI_Policy"
	iamPolicyAmazonEC2ContainerRegistryPowerUser = "AmazonEC2ContainerRegistryPowerUser"
	iamPolicyAmazonEC2ContainerRegistryReadOnly  = "AmazonEC2ContainerRegistryReadOnly"
	iamPolicyCloudWatchAgentServerPolicy         = "CloudWatchAgentServerPolicy"
	iamPolicyAmazonSSMManagedInstanceCore        = "AmazonSSMManagedInstanceCore"

	iamPolicyAmazonEKSFargatePodExecutionRolePolicy = "AmazonEKSFargatePodExecutionRolePolicy"
)

const (
	cfnIAMInstanceRoleName    = "NodeInstanceRole"
	cfnIAMInstanceProfileName = "NodeInstanceProfile"
)

var (
	iamDefaultNodePolicies = []string{
		iamPolicyAmazonEKSWorkerNodePolicy,
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

func (c *ClusterResourceSet) addResourcesForIAM() {
	c.rs.withNamedIAM = false

	if api.IsSetAndNonEmptyString(c.spec.IAM.ServiceRoleARN) {
		c.rs.withIAM = false
		c.rs.defineOutputWithoutCollector(outputs.ClusterServiceRoleARN, c.spec.IAM.ServiceRoleARN, true)
		return
	}

	c.rs.withIAM = true

	managedPolicyArns := []string{
		iamPolicyAmazonEKSClusterPolicy,
	}
	if !api.IsDisabled(c.spec.IAM.VPCResourceControllerPolicy) {
		managedPolicyArns = append(managedPolicyArns, iamPolicyAmazonEKSVPCResourceController)
	}

	role := &gfniam.Role{
		AssumeRolePolicyDocument: cft.MakeAssumeRolePolicyDocumentForServices(
			MakeServiceRef("EKS"),
		),
		ManagedPolicyArns: gfnt.NewSlice(makePolicyARNs(managedPolicyArns...)...),
	}
	if api.IsSetAndNonEmptyString(c.spec.IAM.ServiceRolePermissionsBoundary) {
		role.PermissionsBoundary = gfnt.NewString(*c.spec.IAM.ServiceRolePermissionsBoundary)
	}
	refSR := c.newResource("ServiceRole", role)
	c.rs.attachAllowPolicy("PolicyCloudWatchMetrics", refSR, cloudWatchMetricsStatements())
	// These are potentially required for creating load balancers but aren't included in the
	// AmazonEKSClusterPolicy
	// See https://docs.aws.amazon.com/elasticloadbalancing/latest/userguide/elb-api-permissions.html#required-permissions-v2
	// and weaveworks/eksctl#2488
	c.rs.attachAllowPolicy("PolicyELBPermissions", refSR, elbStatements())

	c.rs.defineOutputFromAtt(outputs.ClusterServiceRoleARN, "ServiceRole", "Arn", true, func(v string) error {
		c.spec.IAM.ServiceRoleARN = &v
		return nil
	})
}

// WithIAM states, if IAM roles will be created or not
func (n *NodeGroupResourceSet) WithIAM() bool {
	return n.rs.withIAM
}

// WithNamedIAM states, if specifically named IAM roles will be created or not
func (n *NodeGroupResourceSet) WithNamedIAM() bool {
	return n.rs.withNamedIAM
}

func (n *NodeGroupResourceSet) addResourcesForIAM() error {
	if n.spec.IAM.InstanceProfileARN != "" {
		n.rs.withIAM = false
		n.rs.withNamedIAM = false

		// if instance profile is given, as well as the role, we simply use both and export the role
		// (which is needed in order to authorise the nodegroup)
		n.instanceProfileARN = gfnt.NewString(n.spec.IAM.InstanceProfileARN)
		if n.spec.IAM.InstanceRoleARN != "" {
			n.rs.defineOutputWithoutCollector(outputs.NodeGroupInstanceProfileARN, n.spec.IAM.InstanceProfileARN, true)
			n.rs.defineOutputWithoutCollector(outputs.NodeGroupInstanceRoleARN, n.spec.IAM.InstanceRoleARN, true)
			return nil
		}
		// if instance role is not given, export profile and use the getter to call importer function
		n.rs.defineOutput(outputs.NodeGroupInstanceProfileARN, n.spec.IAM.InstanceProfileARN, true, func(v string) error {
			return iam.ImportInstanceRoleFromProfileARN(n.iamAPI, n.spec, v)
		})

		return nil
	}

	n.rs.withIAM = true

	if n.spec.IAM.InstanceRoleARN != "" {
		roleARN := NormalizeARN(n.spec.IAM.InstanceRoleARN)

		// if role is set, but profile isn't - create profile
		n.newResource(cfnIAMInstanceProfileName, &gfniam.InstanceProfile{
			Path:  gfnt.NewString("/"),
			Roles: gfnt.NewStringSlice(AbstractRoleNameFromARN(roleARN)),
		})
		n.instanceProfileARN = gfnt.MakeFnGetAttString(cfnIAMInstanceProfileName, "Arn")
		n.rs.defineOutputFromAtt(outputs.NodeGroupInstanceProfileARN, cfnIAMInstanceProfileName, "Arn", true, func(v string) error {
			n.spec.IAM.InstanceProfileARN = v
			return nil
		})
		n.rs.defineOutputWithoutCollector(outputs.NodeGroupInstanceRoleARN, roleARN, true)
		return nil
	}

	// if neither role nor profile is given - create both

	if n.spec.IAM.InstanceRoleName != "" {
		// setting role name requires additional capabilities
		n.rs.withNamedIAM = true
	}

	if err := createRole(n.rs, n.clusterSpec.IAM, n.spec.IAM, false, n.forceAddCNIPolicy); err != nil {
		return err
	}

	n.newResource(cfnIAMInstanceProfileName, &gfniam.InstanceProfile{
		Path:  gfnt.NewString("/"),
		Roles: gfnt.NewSlice(gfnt.MakeRef(cfnIAMInstanceRoleName)),
	})
	n.instanceProfileARN = gfnt.MakeFnGetAttString(cfnIAMInstanceProfileName, "Arn")

	n.rs.defineOutputFromAtt(outputs.NodeGroupInstanceProfileARN, cfnIAMInstanceProfileName, "Arn", true, func(v string) error {
		n.spec.IAM.InstanceProfileARN = v
		return nil
	})
	n.rs.defineOutputFromAtt(outputs.NodeGroupInstanceRoleARN, cfnIAMInstanceRoleName, "Arn", true, func(v string) error {
		n.spec.IAM.InstanceRoleARN = v
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

// IAMRoleResourceSet holds IAM Role stack build-time information
type IAMRoleResourceSet struct {
	template            *cft.Template
	oidc                *iamoidc.OpenIDConnectManager
	outputs             *outputs.CollectorSet
	roleName            string
	wellKnownPolicies   api.WellKnownPolicies
	attachPolicyARNs    []string
	attachPolicy        api.InlineDocument
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

	var assumeRolePolicyDocument cft.MapOfInterfaces
	if rs.serviceAccount != "" && rs.namespace != "" {
		logger.Debug("service account location provided: %s/%s, adding sub condition", api.AWSNodeMeta.Namespace, api.AWSNodeMeta.Name)
		assumeRolePolicyDocument = rs.oidc.MakeAssumeRolePolicyDocumentWithServiceAccountConditions(rs.namespace, rs.serviceAccount)
	} else {
		assumeRolePolicyDocument = rs.oidc.MakeAssumeRolePolicyDocument()
	}

	role := &cft.IAMRole{
		AssumeRolePolicyDocument: assumeRolePolicyDocument,
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

// RenderJSON will render iamserviceaccount stack as JSON
func (rs *IAMRoleResourceSet) RenderJSON() ([]byte, error) {
	return rs.template.RenderJSON()
}

// GetAllOutputs will get all outputs from iamserviceaccount stack
func (rs *IAMRoleResourceSet) GetAllOutputs(stack cfn.Stack) error {
	return rs.outputs.MustCollect(stack)
}
