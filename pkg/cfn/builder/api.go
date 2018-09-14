package builder

import (
	"fmt"
	"reflect"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	gfn "github.com/awslabs/goformation/cloudformation"
)

const (
	awsStackName = "AWS::StackName"
	awsAccountID = "AWS::AccountId"
	awsRegion    = "AWS::Region"

	ParamClusterName      = "ClusterName"
	ParamClusterStackName = cfnOutputClusterStackName
	ParamNodeGroupID      = "NodeGroupID"

	clusterTemplateDescription                = "EKS cluster"
	clusterTemplateDescriptionDefaultFeatures = " (with dedicated VPC & IAM role) "

	nodeGroupTemplateDescription                = "EKS nodes"
	nodeGroupTemplateDescriptionDefaultFeatures = " (Amazon Linux 2 with SSH) "

	templateDescriptionSuffix = " [created and managed by eksctl]"
)

// ResourceSet is an interface which cluster and nodegroup builders
// must implement
type ResourceSet interface {
	AddAllResources() error
	WithIAM() bool
	RenderJSON() ([]byte, error)
	GetAllOutputs(cfn.Stack) error
}

type resourceSet struct {
	template *gfn.Template
	outputs  []string
	withIAM  bool
}

var refStackName = gfn.Ref(awsStackName)

func newResourceSet() *resourceSet {
	return &resourceSet{
		template: gfn.NewTemplate(),
	}
}

// makeName is syntactic sugar for {"Fn::Sub": "${AWS::Stack}-<name>"}
func makeName(suffix string) string {
	return gfn.Sub(fmt.Sprintf("${%s}-%s", awsStackName, suffix))
}

func (r *resourceSet) newParameter(name, valueType, defaultValue string) string {
	p := map[string]string{"Type": valueType}
	if defaultValue != "" {
		p["Default"] = defaultValue
	}
	r.template.Parameters[name] = p
	return gfn.Ref(name)
}

func (r *resourceSet) newStringParameter(name, defaultValue string) string {
	return r.newParameter(name, "String", defaultValue)
}

func (r *resourceSet) newNumberParameter(name, defaultValue string) string {
	return r.newParameter(name, "Number", defaultValue)
}

// makeAutoNameTag create a new Name tag in the following format:
// {Key: "Name", Value: !Sub "${AWS::StackName}/<logicalResourceName>"}
func makeAutoNameTag(suffix string) gfn.Tag {
	return gfn.Tag{
		Key:   "Name",
		Value: gfn.Sub(fmt.Sprintf("${%s}/%s", awsStackName, suffix)),
	}
}

// maybeSetNameTag adds a Name tag to any resource that supports tags
// it calls makeAutoNameTag to format the tag value
func maybeSetNameTag(name string, resource interface{}) {
	e := reflect.ValueOf(resource).Elem()
	if e.Kind() == reflect.Struct {
		f := e.FieldByName("Tags")
		if f.IsValid() && f.CanSet() {
			tag := reflect.ValueOf(makeAutoNameTag(name))
			if f.Type() == reflect.ValueOf([]gfn.Tag{}).Type() {
				f.Set(reflect.Append(f, tag))
			}
		}
	}
}

// newResource adds a resource, and adds Name tag if possible, it returns a reference
func (r *resourceSet) newResource(name string, resource interface{}) string {
	maybeSetNameTag(name, resource)
	r.template.Resources[name] = resource
	return gfn.Ref(name)
}

// renderJSON renders template as JSON
func (r *resourceSet) renderJSON() ([]byte, error) {
	return r.template.JSON()
}
