package template

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"

	. "github.com/weaveworks/eksctl/pkg/cfn/template"
)

func checkTemplate(actualTemplate interface{}) error {
	if actualTemplate == nil {
		return fmt.Errorf("template is nil")
	}

	if _, ok := actualTemplate.(*Template); !ok {
		return fmt.Errorf("not a template")
	}

	return nil
}

type TemplateLoader struct {
	*commonMatcher

	templateBody []byte
	templatePath string
}

func LoadBytesWithoutErrors(templateBody []byte) types.GomegaMatcher {
	return &TemplateLoader{
		commonMatcher: &commonMatcher{},

		templateBody: templateBody,
	}
}

func LoadStringWithoutErrors(templateBody string) types.GomegaMatcher {
	return &TemplateLoader{
		commonMatcher: &commonMatcher{},

		templateBody: []byte(templateBody),
	}
}

func LoadFileWithoutErrors(templatePath string) types.GomegaMatcher {
	return &TemplateLoader{
		commonMatcher: &commonMatcher{},

		templatePath: templatePath,
	}
}

func (m *TemplateLoader) Match(actualTemplate interface{}) (bool, error) {
	if err := checkTemplate(actualTemplate); err != nil {
		return false, err
	}

	if m.templatePath != "" {
		js, err := ioutil.ReadFile(m.templatePath)
		if err != nil {
			return false, err
		}
		m.templateBody = js
	}

	err := actualTemplate.(*Template).LoadJSON(m.templateBody)
	if err != nil {
		m.err = err
		return false, nil
	}
	return true, nil
}

func (m *TemplateLoader) FailureMessage(_ interface{}) string {
	return m.failureMessageWithError("Expected to load template from JSON without errors")
}

func (m *TemplateLoader) NegatedFailureMessage(_ interface{}) string {
	return fmt.Sprintf("Expected to NOT load template from JSON without errors")
}

type ResourceNameAndTypeMatcher struct {
	*commonMatcher

	resourceName, resourceType string
}

func HaveResource(resourceName, resourceType string) types.GomegaMatcher {
	return &ResourceNameAndTypeMatcher{
		commonMatcher: &commonMatcher{},

		resourceName: resourceName,
		resourceType: resourceType,
	}
}

func (m *ResourceNameAndTypeMatcher) Match(actualTemplate interface{}) (bool, error) {
	if err := checkTemplate(actualTemplate); err != nil {
		return false, err
	}

	actualResource, ok := actualTemplate.(*Template).Resources[m.resourceName]
	if !ok {
		m.err = fmt.Errorf("resource %q not found", m.resourceName)
		return false, nil
	}

	if m.resourceType != "*" && actualResource.Type != m.resourceType {
		m.err = fmt.Errorf("type of resource %q is %q, not %q", m.resourceName, actualResource.Type, m.resourceType)
		return false, nil
	}

	return true, nil
}

func (m *ResourceNameAndTypeMatcher) FailureMessage(_ interface{}) string {
	return m.failureMessageWithError(fmt.Sprintf("Expected the template to have resoruce %q of type %q", m.resourceName, m.resourceType))
}

func (m *ResourceNameAndTypeMatcher) NegatedFailureMessage(_ interface{}) string {
	return fmt.Sprintf("Expected the template to NOT have resoruce %q of type %q", m.resourceName, m.resourceType)
}

type ResourcePropertiesMatcher struct {
	*commonMatcher

	resourceName                string
	propertyName, propertyValue string
}

func HaveResourceWithPropertyValue(resourceName, propertyName, propertyValue string) types.GomegaMatcher {
	return &ResourcePropertiesMatcher{
		commonMatcher: &commonMatcher{},

		resourceName:  resourceName,
		propertyName:  propertyName,
		propertyValue: propertyValue,
	}
}

func (m *ResourcePropertiesMatcher) Match(actualTemplate interface{}) (bool, error) {
	if ok, err := HaveResource(m.resourceName, "*").Match(actualTemplate); !ok {
		return ok, err
	}

	actualProperty := actualTemplate.(*Template).Resources[m.resourceName].Properties

	if actualProperty == nil {
		m.err = fmt.Errorf("resource %q has no properties", m.resourceName)
		return false, nil
	}

	actualValue, ok := actualProperty.(MapOfInterfaces)[m.propertyName]
	if !ok {
		m.err = fmt.Errorf("resource %q does not have property %q", m.resourceName, m.propertyName)
		return false, nil
	}

	if actualValue == nil {
		m.err = fmt.Errorf("property %q of resource %q is nil", m.propertyName, m.resourceName)
		return false, nil
	}

	js, err := json.Marshal(actualValue)
	if err != nil {
		return false, fmt.Errorf("property %q of resource %q cannot be marshalled to JSON: %s", m.propertyName, m.resourceName, err.Error())

	}

	return m.matchJSON(m.propertyValue, js)
}

func (m *ResourcePropertiesMatcher) FailureMessage(_ interface{}) string {
	return m.failureMessageWithError(fmt.Sprintf("Expected the template to have resoruce %q with property %q: %s", m.resourceName, m.propertyName, m.propertyValue))
}

func (m *ResourcePropertiesMatcher) NegatedFailureMessage(_ interface{}) string {
	return fmt.Sprintf("Expected the template to NOT have resoruce %q with property %q: %s", m.resourceName, m.propertyName, m.propertyValue)
}

type OutputsMatcher struct {
	outputNames []string
}

func HaveOutputs(outputNames ...string) types.GomegaMatcher {
	return &OutputsMatcher{
		outputNames: outputNames,
	}
}

func (m *OutputsMatcher) Match(actualTemplate interface{}) (bool, error) {
	if err := checkTemplate(actualTemplate); err != nil {
		return false, err
	}

	actualOutputs := actualTemplate.(*Template).Outputs

	for _, expectedOutputName := range m.outputNames {
		ok, err := gomega.HaveKey(expectedOutputName).Match(actualOutputs)
		if !ok {
			return ok, err
		}
	}

	return true, nil
}

func (m *OutputsMatcher) FailureMessage(_ interface{}) string {
	return fmt.Sprintf("Expected the template to have outputs %v", m.outputNames)
}

func (m *OutputsMatcher) NegatedFailureMessage(_ interface{}) string {
	return fmt.Sprintf("Expected the template to NOT have outputs %v", m.outputNames)
}

type OutputValueMatcher struct {
	*commonMatcher

	outputName, outputValue string
}

func HaveOutputWithValue(outputName, outputValue string) types.GomegaMatcher {
	return &OutputValueMatcher{
		commonMatcher: &commonMatcher{},

		outputName:  outputName,
		outputValue: outputValue,
	}
}

func (m *OutputValueMatcher) Match(actualTemplate interface{}) (bool, error) {
	ok, err := HaveOutputs(m.outputName).Match(actualTemplate)
	if !ok {
		return ok, err
	}

	actualOutput := actualTemplate.(*Template).Outputs[m.outputName]

	if actualOutput.Value == nil {
		return false, fmt.Errorf("output value is nil")
	}
	actualOutputValue := *actualOutput.Value
	js, err := actualOutputValue.MarshalJSON()
	if err != nil {
		return false, fmt.Errorf("value %v of ouput %q cannot be marshalled to JSON: %s", actualOutputValue, m.outputName, err.Error())
	}

	return m.matchJSON(m.outputValue, js)
}

func (m *OutputValueMatcher) FailureMessage(_ interface{}) string {
	return fmt.Sprintf("Expected the template to have output %q with value '%s'", m.outputName, m.outputValue)
}

func (m *OutputValueMatcher) NegatedFailureMessage(_ interface{}) string {
	return fmt.Sprintf("Expected the template to NOT have output %q with value '%s'", m.outputName, m.outputValue)
}

type OutputExportNameMatcher struct {
	*commonMatcher

	outputName, exportName string
}

func HaveOutputExportedAs(outputName, exportName string) types.GomegaMatcher {
	return &OutputExportNameMatcher{
		commonMatcher: &commonMatcher{},

		outputName: outputName,
		exportName: exportName,
	}
}

func (m *OutputExportNameMatcher) Match(actualTemplate interface{}) (bool, error) {
	ok, err := HaveOutputs(m.outputName).Match(actualTemplate)
	if !ok {
		return ok, err
	}

	actualOutput := actualTemplate.(*Template).Outputs[m.outputName]

	if actualOutput.Value == nil {
		return false, fmt.Errorf("output value is nil")
	}

	if actualOutput.Export == nil || actualOutput.Export.Name == nil {
		return false, nil
	}
	actualExportName := *actualOutput.Export.Name
	js, err := actualExportName.MarshalJSON()
	if err != nil {
		return false, fmt.Errorf("export name %v of ouput %q cannot be marshalled to JSON: %s", actualExportName, m.outputName, err.Error())
	}

	return m.matchJSON(m.exportName, js)
}

func (m *OutputExportNameMatcher) FailureMessage(_ interface{}) string {
	return fmt.Sprintf("Expected the template to have output %q with export name '%s'", m.outputName, m.exportName)
}

func (m *OutputExportNameMatcher) NegatedFailureMessage(_ interface{}) string {
	return fmt.Sprintf("Expected the template to NOT have output %q with export name '%s'", m.outputName, m.exportName)
}

type commonMatcher struct {
	err error
}

func (m *commonMatcher) matchJSON(actual interface{}, js []byte) (bool, error) {
	jsMatcher := gomega.MatchJSON(actual)

	ok, err := jsMatcher.Match(js)
	if err != nil {
		m.err = err
		return false, nil
	}
	if !ok {
		m.err = fmt.Errorf(jsMatcher.FailureMessage(js))
	}
	return ok, nil
}

func (m *commonMatcher) failureMessageWithError(msg string) string {
	if m.err != nil {
		msg += fmt.Sprintf("\n\n%s", m.err.Error())
	}
	return msg
}
