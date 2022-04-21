package builder

import (
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	gfn "github.com/weaveworks/goformation/v4/cloudformation"
	gfncfn "github.com/weaveworks/goformation/v4/cloudformation/cloudformation"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
)

const (
	clusterTemplateDescription   = "EKS cluster"
	nodeGroupTemplateDescription = "EKS nodes"
	templateDescriptionSuffix    = "[created and managed by eksctl]"
)

type awsCloudFormationResource struct {
	Type         string
	Properties   map[string]interface{}
	UpdatePolicy map[string]map[string]interface{} `json:",omitempty"`
	DependsOn    []string                          `json:",omitempty"`
}

func (r *awsCloudFormationResource) AWSCloudFormationType() string {
	return r.Type
}

// ResourceSet is an interface which cluster and nodegroup builders
// must implement
type ResourceSet interface {
	AddAllResources() error
	ResourceSetReader
}

// ResourceSetReader contains the set of operations required to create a stack and to collect outputs.
type ResourceSetReader interface {
	RenderJSON() ([]byte, error)
	WithIAM() bool
	WithNamedIAM() bool
	GetAllOutputs(types.Stack) error
}

type resourceSet struct {
	template     *gfn.Template
	outputs      *outputs.CollectorSet
	withIAM      bool
	withNamedIAM bool
}

func newResourceSet() *resourceSet {
	return &resourceSet{
		template: gfn.NewTemplate(),
		outputs:  outputs.NewCollectorSet(nil),
	}
}

// makeName is syntactic sugar for {"Fn::Sub": "${AWS::Stack}-<name>"}
func makeName(suffix string) *gfnt.Value {
	return gfnt.MakeFnSubString(fmt.Sprintf("${%s}-%s", gfnt.StackName, suffix))
}

// makeSlice makes a slice from a list of string arguments
func makeStringSlice(s ...string) []*gfnt.Value {
	slice := []*gfnt.Value{}
	for _, i := range s {
		slice = append(slice, gfnt.NewString(i))
	}
	return slice
}

// makeAutoNameTag create a new Name tag in the following format:
// {Key: "Name", Value: !Sub "${AWS::StackName}/<logicalResourceName>"}
func makeAutoNameTag(suffix string) gfncfn.Tag {
	return gfncfn.Tag{
		Key:   gfnt.NewString("Name"),
		Value: gfnt.MakeFnSubString(fmt.Sprintf("${%s}/%s", gfnt.StackName, suffix)),
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
			if f.Type() == reflect.ValueOf([]gfncfn.Tag{}).Type() {
				f.Set(reflect.Append(f, tag))
			}
		}
	}
}

// newResource adds a resource, and adds Name tag if possible, it returns a reference
func (r *resourceSet) newResource(name string, resource gfn.Resource) *gfnt.Value {
	maybeSetNameTag(name, resource)
	r.template.Resources[name] = resource
	return gfnt.MakeRef(name)
}

// renderJSON renders template as JSON
func (r *resourceSet) renderJSON() ([]byte, error) {
	return r.template.JSON()
}
