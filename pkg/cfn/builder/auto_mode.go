package builder

import (
	_ "embed"
	"errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"

	"goformation/v4"
	gfn "goformation/v4/cloudformation"
	gfnt "goformation/v4/cloudformation/types"

	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

const (
	iamPolicyAmazonEKSComputePolicy       = "AmazonEKSComputePolicy"
	iamPolicyAmazonEKSBlockStoragePolicy  = "AmazonEKSBlockStoragePolicy"
	iamPolicyAmazonEKSLoadBalancingPolicy = "AmazonEKSLoadBalancingPolicy"
	iamPolicyAmazonEKSNetworkingPolicy    = "AmazonEKSNetworkingPolicy"
)

// AutoModeIAMPolicies is a list of policies required by EKS Auto Mode.
var AutoModeIAMPolicies = []string{iamPolicyAmazonEKSComputePolicy, iamPolicyAmazonEKSBlockStoragePolicy,
	iamPolicyAmazonEKSLoadBalancingPolicy, iamPolicyAmazonEKSNetworkingPolicy}

//go:embed roles/auto-mode-node-role.yaml
var autoModeNodeRoleTemplate []byte

type AutoModeRefs struct {
	NodeRole *gfnt.Value
}

func AddAutoModeResources(clusterTemplate *gfn.Template) (AutoModeRefs, error) {
	template, err := goformation.ParseYAML(autoModeNodeRoleTemplate)
	if err != nil {
		return AutoModeRefs{}, err
	}
	for resourceName, resource := range template.Resources {
		clusterTemplate.Resources[resourceName] = resource
	}
	for key, output := range template.Outputs {
		clusterTemplate.Outputs[key] = output
	}
	return AutoModeRefs{
		NodeRole: gfnt.MakeFnGetAttString("AutoModeNodeRole", "Arn"),
	}, nil
}

func CreateAutoModeResourceSet() (*AutoModeResourceSet, error) {
	template, err := goformation.ParseYAML(autoModeNodeRoleTemplate)
	if err != nil {
		return nil, err
	}
	template.Mappings = map[string]interface{}{
		servicePrincipalPartitionMapName: api.Partitions.ServicePrincipalPartitionMappings(),
	}
	return &AutoModeResourceSet{
		template: template,
	}, nil
}

type AutoModeResourceSet struct {
	template    *gfn.Template
	nodeRoleARN string
}

func (e *AutoModeResourceSet) RenderJSON() ([]byte, error) {
	return e.template.JSON()
}

func (e *AutoModeResourceSet) WithIAM() bool {
	return true
}

func (e *AutoModeResourceSet) WithNamedIAM() bool {
	return false
}

func (e *AutoModeResourceSet) GetAllOutputs(stack cfntypes.Stack) error {
	nodeRoleARN, found := GetAutoModeOutputs(stack)
	if !found {
		return errors.New("node role ARN output not found in Auto Mode stack")
	}
	e.nodeRoleARN = nodeRoleARN
	return nil
}

func (e *AutoModeResourceSet) GetAutoModeRoleARN() string {
	return e.nodeRoleARN
}

func GetAutoModeOutputs(stack cfntypes.Stack) (string, bool) {
	const nodeRoleOutputName = "AutoModeNodeRoleARN"
	for _, output := range stack.Outputs {
		if *output.OutputKey == nodeRoleOutputName {
			return *output.OutputValue, true
		}
	}
	return "", false
}
