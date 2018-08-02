package builder

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	gfn "github.com/awslabs/goformation/cloudformation"

	"github.com/kubicorn/kubicorn/pkg/logger"
)

const (
	awsStackName  = "AWS::StackName"
	awsAccountID  = "AWS::AccountId"
	awsRegion     = "AWS::Region"
	fnSub         = "Fn::Sub"
	fnJoin        = "Fn::Join"
	fnSplit       = "Fn::Split"
	fnGetAtt      = "Fn::GetAtt"
	fnImportValue = "Fn::ImportValue"

	ParamClusterName      = "ClusterName"
	ParamClusterStackName = cfnOutputClusterStackName
	ParamNodeGroupID      = "NodeGroupID"

	clusterTemplateDescription                = "EKS cluster"
	clusterTemplateDescriptionDefaultFeatures = " (with dedicated VPC & IAM role) "

	nodeGroupTemplateDescription                = "EKS nodes"
	nodeGroupTemplateDescriptionDefaultFeatures = " (Amazon Linux 2 with SSH) "

	templateDescriptionSuffix = " [created and managed by eksctl]"
)

type resourceSet struct {
	template *gfn.Template
	outputs  []string
}

func newResourceSet() *resourceSet {
	return &resourceSet{
		template: gfn.NewTemplate(),
	}
}

func makeRef(refName string) *gfn.StringIntrinsic {
	return gfn.NewStringRef(refName)
}

func makeSub(expr string) *gfn.StringIntrinsic {
	return gfn.NewStringIntrinsic(fnSub, expr)
}

func makeName(suffix string) *gfn.StringIntrinsic {
	return makeSub(fmt.Sprintf("${%s}-%s", awsStackName, suffix))
}

func makeSlice(i ...*gfn.StringIntrinsic) []*gfn.StringIntrinsic {
	return i
}

func makeStringSlice(s ...string) []*gfn.StringIntrinsic {
	slice := []*gfn.StringIntrinsic{}
	for _, i := range s {
		slice = append(slice, gfn.NewString(i))
	}
	return slice
}

func makeAutoNameTag(suffix string) gfn.Tag {
	return gfn.Tag{
		Key:   gfn.NewString("Name"),
		Value: makeSub(fmt.Sprintf("${%s}/%s", awsStackName, suffix)),
	}
}

var refStackName = makeRef(awsStackName)

func (r *resourceSet) newParameter(name, valueType, defaultValue string) *gfn.StringIntrinsic {
	p := map[string]string{"Type": valueType}
	if defaultValue != "" {
		p["Default"] = defaultValue
	}
	r.template.Parameters[name] = p
	return makeRef(name)
}

func (r *resourceSet) newStringParameter(name, defaultValue string) *gfn.StringIntrinsic {
	return r.newParameter(name, "String", defaultValue)
}

func (r *resourceSet) newResource(name string, resource interface{}) *gfn.StringIntrinsic {
	r.template.Resources[name] = resource
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
	return makeRef(name)
}

func (r *resourceSet) renderJSON() ([]byte, error) {
	return r.template.JSON()
}

func exportName(prefix, output string) string {
	return fmt.Sprintf("${%s}::%s", prefix, output)
}

func (r *resourceSet) newOutput(name string, value interface{}, export bool) {
	o := map[string]interface{}{"Value": value}
	if export {
		o["Export"] = map[string]map[string]string{
			"Name": map[string]string{fnSub: exportName(awsStackName, name)},
		}
	}
	r.template.Outputs[name] = o
	r.outputs = append(r.outputs, name)
}

func (r *resourceSet) newJoinedOutput(name string, value []*gfn.StringIntrinsic, export bool) {
	r.newOutput(name, map[string][]interface{}{fnJoin: []interface{}{",", value}}, export)
}

func (r *resourceSet) newOutputFromAtt(name, att string, export bool) {
	r.newOutput(name, map[string]string{fnGetAtt: att}, export)
}

func makeImportValue(prefix, output string) *gfn.StringIntrinsic {
	return gfn.NewStringIntrinsic(fnImportValue, makeSub(exportName(prefix, output)))
}

func getOutput(stack *cfn.Stack, key string) *string {
	for _, x := range stack.Outputs {
		if *x.OutputKey == key {
			return x.OutputValue
		}
	}
	return nil
}

func setOutput(obj interface{}, key, value string) error {
	e := reflect.ValueOf(obj).Elem()
	if e.Kind() == reflect.Struct {
		f := e.FieldByName(key)
		if f.IsValid() && f.CanSet() {
			switch f.Kind() {
			case reflect.String:
				f.SetString(value)
			case reflect.Slice:
				switch f.Type() {
				case reflect.ValueOf([]string{}).Type():
					f.Set(reflect.ValueOf(strings.Split(value, ",")))
				case reflect.ValueOf([]byte{}).Type():
					data, err := base64.StdEncoding.DecodeString(value)
					if err != nil {
						return errors.Wrapf(err, "decoding value of %q", key)
					}
					f.Set(reflect.ValueOf(data))
				default:
					return fmt.Errorf("unexpected type %q of destination field for %q", f.Type(), key)
				}
			default:
				return fmt.Errorf("unexpected kind %q of destination field for %q", f.Kind(), key)
			}
		} else {
			return fmt.Errorf("cannot set destination field for %q", key)
		}
	} else {
		return fmt.Errorf("cannot use destination interface of type %q", e.Kind())
	}
	return nil
}

func (r *resourceSet) GetAllOutputs(stack cfn.Stack, obj interface{}) error {
	logger.Debug("processing stack outputs")
	for _, key := range r.outputs {
		value := getOutput(&stack, key)
		if value == nil {
			return fmt.Errorf("%s is nil", key)
		}
		if err := setOutput(obj, key, *value); err != nil {
			return errors.Wrap(err, "processing stack outputs")
		}
	}
	return nil
}
