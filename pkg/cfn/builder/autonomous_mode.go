package builder

import (
	_ "embed"
	"errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"

	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/weaveworks/goformation/v4"
	gfn "github.com/weaveworks/goformation/v4/cloudformation"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"
)

//go:embed roles/autonomous-mode-node-role.yaml
var autonomousModeNodeRoleTemplate []byte

type AutonomousModeRefs struct {
	NodeRole *gfnt.Value
}

func AddAutonomousModeResources(clusterTemplate *gfn.Template) (AutonomousModeRefs, error) {
	template, err := goformation.ParseYAML(autonomousModeNodeRoleTemplate)
	if err != nil {
		return AutonomousModeRefs{}, err
	}
	for resourceName, resource := range template.Resources {
		clusterTemplate.Resources[resourceName] = resource
	}
	for key, output := range template.Outputs {
		clusterTemplate.Outputs[key] = output
	}
	return AutonomousModeRefs{
		NodeRole: gfnt.MakeFnGetAttString("AutonomousModeNodeRole", "Arn"),
	}, nil
}

func CreateAutonomousModeResourceSet() (*AutonomousModeResourceSet, error) {
	template, err := goformation.ParseYAML(autonomousModeNodeRoleTemplate)
	if err != nil {
		return nil, err
	}
	template.Mappings = map[string]interface{}{
		servicePrincipalPartitionMapName: api.Partitions.ServicePrincipalPartitionMappings(),
	}
	return &AutonomousModeResourceSet{
		template: template,
	}, nil
}

type AutonomousModeResourceSet struct {
	template    *gfn.Template
	nodeRoleARN string
}

func (e *AutonomousModeResourceSet) RenderJSON() ([]byte, error) {
	return e.template.JSON()
}

func (e *AutonomousModeResourceSet) WithIAM() bool {
	return true
}

func (e *AutonomousModeResourceSet) WithNamedIAM() bool {
	return false
}

func (e *AutonomousModeResourceSet) GetAllOutputs(stack cfntypes.Stack) error {
	nodeRoleARN, found := GetAutonomousModeOutputs(stack)
	if !found {
		return errors.New("node role ARN output not found in Autonomous Mode stack")
	}
	e.nodeRoleARN = nodeRoleARN
	return nil
}

func (e *AutonomousModeResourceSet) GetAutonomousModeRoleARN() string {
	return e.nodeRoleARN
}

func GetAutonomousModeOutputs(stack cfntypes.Stack) (string, bool) {
	const nodeRoleOutputName = "AutonomousModeNodeRoleARN"
	for _, output := range stack.Outputs {
		if *output.OutputKey == nodeRoleOutputName {
			return *output.OutputValue, true
		}
	}
	return "", false
}
