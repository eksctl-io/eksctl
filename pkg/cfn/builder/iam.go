package builder

import (
	"fmt"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	gfn "github.com/awslabs/goformation/cloudformation"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	cft "github.com/weaveworks/eksctl/pkg/cfn/template"
	"github.com/weaveworks/eksctl/pkg/iam"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
)

const (
	iamPolicyAmazonEKSServicePolicyARN = "arn:aws:iam::aws:policy/AmazonEKSServicePolicy"
	iamPolicyAmazonEKSClusterPolicyARN = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"

	iamPolicyAmazonEKSWorkerNodePolicyARN           = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
	iamPolicyAmazonEKSCNIPolicyARN                  = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
	iamPolicyAmazonEC2ContainerRegistryPowerUserARN = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPowerUser"
	iamPolicyAmazonEC2ContainerRegistryReadOnlyARN  = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
	iamPolicyCloudWatchAgentServerPolicyARN         = "arn:aws:iam::aws:policy/CloudWatchAgentServerPolicy"

	iamPolicyAmazonEKSFargatePodExecutionRolePolicyARN = "arn:aws:iam::aws:policy/AmazonEKSFargatePodExecutionRolePolicy"
)

const (
	cfnIAMInstanceRoleName    = "NodeInstanceRole"
	cfnIAMInstanceProfileName = "NodeInstanceProfile"
)

var (
	iamDefaultNodePolicyARNs = []string{
		iamPolicyAmazonEKSWorkerNodePolicyARN,
		iamPolicyAmazonEKSCNIPolicyARN,
	}
)

func (c *resourceSet) attachAllowPolicy(name string, refRole *gfn.Value, resources interface{}, actions []string) {
	c.newResource(name, &gfn.AWSIAMPolicy{
		PolicyName: makeName(name),
		Roles:      makeSlice(refRole),
		PolicyDocument: cft.MakePolicyDocument(map[string]interface{}{
			"Effect":   "Allow",
			"Resource": resources,
			"Action":   actions,
		}),
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

	refSR := c.newResource("ServiceRole", &gfn.AWSIAMRole{
		AssumeRolePolicyDocument: cft.MakeAssumeRolePolicyDocumentForServices(
			"eks.amazonaws.com",
			// Ensure that EKS can schedule pods onto Fargate, should the user
			// define so-called "Fargate profiles" in order to do so:
			"eks-fargate-pods.amazonaws.com",
		),
		ManagedPolicyArns: makeStringSlice(
			iamPolicyAmazonEKSServicePolicyARN,
			iamPolicyAmazonEKSClusterPolicyARN,
		),
	})
	c.rs.attachAllowPolicy("PolicyNLB", refSR, "*", []string{
		"elasticloadbalancing:*",
		"ec2:CreateSecurityGroup",
		"ec2:Describe*",
	})
	c.rs.attachAllowPolicy("PolicyCloudWatchMetrics", refSR, "*", []string{
		"cloudwatch:PutMetricData",
	})

	c.rs.defineOutputFromAtt(outputs.ClusterServiceRoleARN, "ServiceRole.Arn", true, func(v string) error {
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

func (n *NodeGroupResourceSet) addResourcesForIAM() {
	if n.spec.IAM.InstanceProfileARN != "" {
		n.rs.withIAM = false
		n.rs.withNamedIAM = false

		// if instance profile is given, as well as the role, we simply use both and export the role
		// (which is needed in order to authorise the nodegroup)
		n.instanceProfileARN = gfn.NewString(n.spec.IAM.InstanceProfileARN)
		if n.spec.IAM.InstanceRoleARN != "" {
			n.rs.defineOutputWithoutCollector(outputs.NodeGroupInstanceProfileARN, n.spec.IAM.InstanceProfileARN, true)
			n.rs.defineOutputWithoutCollector(outputs.NodeGroupInstanceRoleARN, n.spec.IAM.InstanceRoleARN, true)
			return
		}
		// if instance role is not given, export profile and use the getter to call importer function
		n.rs.defineOutput(outputs.NodeGroupInstanceProfileARN, n.spec.IAM.InstanceProfileARN, true, func(v string) error {
			return iam.ImportInstanceRoleFromProfileARN(n.provider, n.spec, v)
		})

		return
	}

	n.rs.withIAM = true

	if n.spec.IAM.InstanceRoleARN != "" {
		// if role is set, but profile isn't - create profile
		n.newResource(cfnIAMInstanceProfileName, &gfn.AWSIAMInstanceProfile{
			Path:  gfn.NewString("/"),
			Roles: makeStringSlice(n.spec.IAM.InstanceRoleARN),
		})
		n.instanceProfileARN = gfn.MakeFnGetAttString(makeAttrAccessor(cfnIAMInstanceProfileName, "Arn"))
		n.rs.defineOutputFromAtt(outputs.NodeGroupInstanceProfileARN, makeAttrAccessor(cfnIAMInstanceProfileName, "Arn"), true, func(v string) error {
			n.spec.IAM.InstanceProfileARN = v
			return nil
		})
		n.rs.defineOutputWithoutCollector(outputs.NodeGroupInstanceRoleARN, n.spec.IAM.InstanceRoleARN, true)
		return
	}

	// if neither role nor profile is given - create both

	if n.spec.IAM.InstanceRoleName != "" {
		// setting role name requires additional capabilities
		n.rs.withNamedIAM = true
	}

	createRole(n.rs, n.spec.IAM, false)

	n.newResource(cfnIAMInstanceProfileName, &gfn.AWSIAMInstanceProfile{
		Path:  gfn.NewString("/"),
		Roles: makeSlice(gfn.MakeRef(cfnIAMInstanceRoleName)),
	})
	n.instanceProfileARN = gfn.MakeFnGetAttString(makeAttrAccessor(cfnIAMInstanceProfileName, "Arn"))

	n.rs.defineOutputFromAtt(outputs.NodeGroupInstanceProfileARN, makeAttrAccessor(cfnIAMInstanceProfileName, "Arn"), true, func(v string) error {
		n.spec.IAM.InstanceProfileARN = v
		return nil
	})
	n.rs.defineOutputFromAtt(outputs.NodeGroupInstanceRoleARN, makeAttrAccessor(cfnIAMInstanceRoleName, "Arn"), true, func(v string) error {
		n.spec.IAM.InstanceRoleARN = v
		return nil
	})
}

// IAMServiceAccountResourceSet holds iamserviceaccount stack build-time information
type IAMServiceAccountResourceSet struct {
	template *cft.Template
	spec     *api.ClusterIAMServiceAccount
	oidc     *iamoidc.OpenIDConnectManager
	outputs  *outputs.CollectorSet
}

// NewIAMServiceAccountResourceSet builds iamserviceaccount stack from the give spec
func NewIAMServiceAccountResourceSet(spec *api.ClusterIAMServiceAccount, oidc *iamoidc.OpenIDConnectManager) *IAMServiceAccountResourceSet {
	return &IAMServiceAccountResourceSet{
		template: cft.NewTemplate(),
		spec:     spec,
		oidc:     oidc,
	}
}

// WithIAM returns true
func (*IAMServiceAccountResourceSet) WithIAM() bool { return true }

// WithNamedIAM returns false
func (*IAMServiceAccountResourceSet) WithNamedIAM() bool { return false }

// AddAllResources adds all resources for the stack
func (rs *IAMServiceAccountResourceSet) AddAllResources() error {
	rs.template.Description = fmt.Sprintf(
		"IAM role for serviceaccount %q %s",
		rs.spec.NameString(),
		templateDescriptionSuffix,
	)

	// we use a single stack for each service account, but there maybe a few roles in it
	// so will need to give them unique names
	// we will need to consider using a large stack for all the roles, but that needs some
	// testing and potentially a better stack mutation strategy
	role := &cft.IAMRole{
		AssumeRolePolicyDocument: rs.oidc.MakeAssumeRolePolicyDocument(rs.spec.Namespace, rs.spec.Name),
	}
	role.ManagedPolicyArns = append(role.ManagedPolicyArns, rs.spec.AttachPolicyARNs...)

	roleRef := rs.template.NewResource("Role1", role)

	// TODO: declare output collector automatically when all stack builders migrated to our template package
	rs.template.Outputs["Role1"] = cft.Output{
		Value: cft.MakeFnGetAttString("Role1.Arn"),
	}
	rs.outputs = outputs.NewCollectorSet(map[string]outputs.Collector{
		"Role1": func(v string) error {
			rs.spec.Status = &api.ClusterIAMServiceAccountStatus{
				RoleARN: &v,
			}
			return nil
		},
	})

	if len(rs.spec.AttachPolicy) != 0 {
		rs.template.AttachPolicy("Policy1", roleRef, cft.MapOfInterfaces(rs.spec.AttachPolicy))
	}

	return nil
}

// RenderJSON will render iamserviceaccount stack as JSON
func (rs *IAMServiceAccountResourceSet) RenderJSON() ([]byte, error) {
	return rs.template.RenderJSON()
}

// GetAllOutputs will get all outputs from iamserviceaccount stack
func (rs *IAMServiceAccountResourceSet) GetAllOutputs(stack cfn.Stack) error {
	return rs.outputs.MustCollect(stack)
}
