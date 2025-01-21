package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"text/template"
)

// ResourceGenerator takes AWS CloudFormation Resource Specification
// documents, and generates Go structs and a JSON Schema from them.
type ResourceGenerator struct {
	primaryURL   string
	fragmentUrls map[string]string
	Results      *ResourceGeneratorResults
}

// GeneratedResource represents a CloudFormation resource that has been generated
// and includes the name (e.g. AWS::S3::Bucket), Package Name (e.g. s3) and struct name (e.g. Bucket)
type GeneratedResource struct {
	Name        string
	BaseName    string
	PackageName string
	StructName  string
}

// ResourceGeneratorResults contains a summary of the items generated
type ResourceGeneratorResults struct {
	AllResources     []GeneratedResource
	UpdatedResources []GeneratedResource
	UpdatedSchemas   map[string]string
	ProcessedCount   int
}

var (
	// ResourcesThatSupportUpdatePolicies defines which CloudFormation resources support UpdatePolicies
	// see: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-attribute-updatepolicy.html
	ResourcesThatSupportUpdatePolicies = []string{
		"AWS::AutoScaling::AutoScalingGroup",
		"AWS::Lambda::Alias",
	}

	// ResourcesThatSupportCreationPolicies defines which CloudFormation resources support CreationPolicies
	// see: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-attribute-creationpolicy.html
	ResourcesThatSupportCreationPolicies = []string{
		"AWS::AutoScaling::AutoScalingGroup",
		"AWS::EC2::Instance",
		"AWS::CloudFormation::WaitCondition",
	}
)

// NewResourceGenerator contains a primary AWS CloudFormation Resource Specification
// document and an array of fragment Resource Specification documents (such as transforms),
// and generates Go structs and a JSON Schema from them.
// The input can be a mix of URLs (https://) or files (file://).
func NewResourceGenerator(primaryURL string, fragmentUrls map[string]string) (*ResourceGenerator, error) {

	rg := &ResourceGenerator{
		primaryURL:   primaryURL,
		fragmentUrls: fragmentUrls,
		Results: &ResourceGeneratorResults{
			UpdatedResources: []GeneratedResource{},
			AllResources:     []GeneratedResource{},
			UpdatedSchemas:   map[string]string{},
			ProcessedCount:   0,
		},
	}

	return rg, nil
}

// Generate generates Go structs and a JSON Schema from the AWS CloudFormation
// Resource Specifications provided to NewResourceGenerator
func (rg *ResourceGenerator) Generate() error {

	// Process the primary template first, since the primary template resources
	// are added to the JSON schema for fragment transform specs
	fmt.Printf("Downloading cloudformation specification from %s\n", rg.primaryURL)
	primaryData, err := rg.downloadSpec(rg.primaryURL)
	if err != nil {
		return err
	}
	primarySpec, err := rg.processSpec("cloudformation", primaryData)
	if err != nil {
		return err
	}
	if err := rg.generateJSONSchema("cloudformation", primarySpec); err != nil {
		return err
	}

	for name, url := range rg.fragmentUrls {
		fmt.Printf("Downloading %s specification from %s\n", name, url)
		data, err := rg.downloadSpec(url)
		if err != nil {
			return err
		}
		spec, err := rg.processSpec(name, data)
		if err != nil {
			return err
		}
		// Append main CloudFormation schema resource types and property types to this fragment
		for k, v := range primarySpec.Resources {
			spec.Resources[k] = v
		}
		for k, v := range primarySpec.Properties {
			spec.Properties[k] = v
		}
		if err := rg.generateJSONSchema(name, spec); err != nil {
			return err
		}
	}

	if err := rg.generateAllResourcesMap(rg.Results.AllResources); err != nil {
		return err
	}

	return nil
}

func (rg *ResourceGenerator) downloadSpec(location string) ([]byte, error) {
	uri, err := url.Parse(location)
	if err != nil {
		return nil, err
	}

	switch uri.Scheme {
	case "https", "http":
		uri := uri.String()
		response, err := http.Get(uri)
		if err != nil {
			return nil, err
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		return data, nil
	case "file":
		data, err := ioutil.ReadFile(strings.Replace(location, "file://", "", -1))
		if err != nil {
			return nil, err
		}

		return data, nil
	}

	return nil, fmt.Errorf("invalid URL scheme %s", uri.Scheme)
}

func (rg *ResourceGenerator) processSpec(specname string, data []byte) (*CloudFormationResourceSpecification, error) {

	// Unmarshall the JSON specification
	spec := &CloudFormationResourceSpecification{}
	if err := json.Unmarshal(data, spec); err != nil {
		return nil, err
	}

	// Add the resources processed to the ResourceGenerator output
	for name := range spec.Resources {

		sname, err := structName(name)
		if err != nil {
			return nil, err
		}

		pname, err := packageName(name, true)
		if err != nil {
			return nil, err
		}

		basename, err := packageName(name, false)
		if err != nil {
			return nil, err
		}

		rg.Results.AllResources = append(rg.Results.AllResources, GeneratedResource{
			Name:        name,
			BaseName:    basename,
			PackageName: pname,
			StructName:  sname,
		})
	}

	// Write all of the resources in the spec file
	for name, resource := range spec.Resources {
		if err := rg.generateResources(name, resource, false, spec); err != nil {
			return nil, err
		}
	}

	// Write all of the custom types in the spec file
	for name, resource := range spec.Properties {
		if err := rg.generateResources(name, resource, true, spec); err != nil {
			return nil, err
		}
	}

	return spec, nil

}

func (rg *ResourceGenerator) generateAllResourcesMap(resources []GeneratedResource) error {

	// Open the all resources template
	tmpl, err := template.ParseFiles("generate/templates/all.template")
	if err != nil {
		return fmt.Errorf("failed to load resource template: %s", err)
	}

	// Sort the resources by cloudformation type, so that they don't rearrange
	// every time the all.go file is generated
	sort.Slice(resources, func(i, j int) bool {
		return resources[i].Name < resources[j].Name
	})

	templateData := struct {
		Resources []GeneratedResource
	}{
		Resources: resources,
	}

	// Execute the template, writing it to a buffer
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, templateData)
	if err != nil {
		return fmt.Errorf("failed to generate iterable map of all resources: %s", err)
	}

	// Format the generated Go code with gofmt
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to format Go file for iterable map of all resources: %s", err)
	}

	// Write the file contents out
	if err := ioutil.WriteFile("cloudformation/all.go", formatted, 0644); err != nil {
		return fmt.Errorf("failed to write Go file for iterable map of all resources: %s", err)
	}

	return nil

}

func (rg *ResourceGenerator) generateResources(name string, resource Resource, isCustomProperty bool, spec *CloudFormationResourceSpecification) error {

	// Open the resource template
	tmpl, err := template.ParseFiles("generate/templates/resource.template")
	if err != nil {
		return fmt.Errorf("failed to load resource template: %s", err)
	}

	// Check if this resource type allows specifying a CloudFormation UpdatePolicy
	hasUpdatePolicy := false
	for _, res := range ResourcesThatSupportUpdatePolicies {
		if name == res {
			hasUpdatePolicy = true
			break
		}
	}

	// Check if this resource type allows specifying a CloudFormation CreationPolicy
	hasCreationPolicy := false
	for _, res := range ResourcesThatSupportCreationPolicies {
		if name == res {
			hasCreationPolicy = true
			break
		}
	}

	// Pass in the following information into the template
	sname, err := structName(name)
	if err != nil {
		return err
	}

	pname, err := packageName(name, true)
	if err != nil {
		return err
	}

	// Check if this resource has any tag properties (and if it's not in the
	// same package as the `Tag` type
	// note: the property might not always be called 'Tags'
	// see: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dlm-lifecyclepolicy-schedule.html#cfn-dlm-lifecyclepolicy-schedule-tagstoadd
	hasTags := false
	for _, property := range resource.Properties {
		if property.ItemType == "Tag" {
			hasTags = true && pname != "cloudformation"
			break
		}
	}

	nameParts := strings.Split(name, "::")
	nameParts = strings.Split(nameParts[len(nameParts)-1], ".")
	basename := nameParts[0]

	templateData := struct {
		Name              string
		PackageName       string
		StructName        string
		Basename          string
		Resource          Resource
		IsCustomProperty  bool
		Version           string
		HasUpdatePolicy   bool
		HasCreationPolicy bool
		HasTags           bool
	}{
		Name:              name,
		PackageName:       pname,
		StructName:        sname,
		Basename:          basename,
		Resource:          resource,
		IsCustomProperty:  isCustomProperty,
		Version:           spec.ResourceSpecificationVersion,
		HasUpdatePolicy:   hasUpdatePolicy,
		HasCreationPolicy: hasCreationPolicy,
		HasTags:           hasTags,
	}

	// Execute the template, writing it to a buffer
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, templateData)
	if err != nil {
		return fmt.Errorf("failed to generate resource %s: %s", name, err)
	}

	// Format the generated Go code with gofmt
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println(string(buf.Bytes()))
		return fmt.Errorf("failed to format Go file for resource %s: %s", name, err)
	}

	// Check if the file has changed since the last time generate ran
	dir := "cloudformation/" + pname
	fn := dir + "/" + filename(name)
	current, err := ioutil.ReadFile(fn)

	if err != nil || bytes.Compare(formatted, current) != 0 {

		// Create the directory if it doesn't exist
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.Mkdir(dir, 0755)
		}

		// Write the file contents out
		if err := ioutil.WriteFile(fn, formatted, 0644); err != nil {
			return fmt.Errorf("failed to write resource file %s: %s", fn, err)
		}
		// Log the updated resource name to the results
		rg.Results.UpdatedResources = append(rg.Results.UpdatedResources, GeneratedResource{
			Name:        name,
			BaseName:    basename,
			PackageName: pname,
			StructName:  sname,
		})

	}

	rg.Results.ProcessedCount++

	return nil

}

func (rg *ResourceGenerator) generateJSONSchema(specname string, spec *CloudFormationResourceSpecification) error {

	// Open the schema template and setup a counter function that will
	// available in the template to be used to detect when trailing commas
	// are required in the JSON when looping through maps
	tmpl, err := template.New("schema.template").Funcs(template.FuncMap{
		"counter": counter,
	}).ParseFiles("generate/templates/schema.template")

	var buf bytes.Buffer

	// Execute the template, writing it to file
	err = tmpl.Execute(&buf, spec)
	if err != nil {
		return fmt.Errorf("failed to generate JSON Schema: %s", err)
	}

	// Parse it to JSON objects and back again to format it
	var j interface{}
	if err := json.Unmarshal(buf.Bytes(), &j); err != nil {
		return fmt.Errorf("failed to unmarshal JSON Schema: %s", err)
	}

	formatted, err := json.MarshalIndent(j, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON Schema: %s", err)

	}

	filename := fmt.Sprintf("schema/%s.schema.json", specname)
	if err := ioutil.WriteFile(filename, formatted, 0644); err != nil {
		return fmt.Errorf("failed to write JSON Schema: %s", err)
	}

	// Also create a Go importable version
	var gocode []byte
	gocode = append(gocode, []byte("package schema\n")...)
	gocode = append(gocode, []byte("\n")...)
	gocode = append(gocode, []byte("// "+strings.Title(specname)+"Schema defined a JSON Schema that can be used to validate CloudFormation/SAM templates\n")...)
	gocode = append(gocode, []byte("var "+strings.Title(specname)+"Schema = `")...)
	gocode = append(gocode, formatted...)
	gocode = append(gocode, []byte("`\n")...)
	gofilename := fmt.Sprintf("schema/%s.go", specname)
	if err := ioutil.WriteFile(gofilename, gocode, 0644); err != nil {
		return fmt.Errorf("failed to write Go version of JSON Schema: %s", err)
	}

	rg.Results.UpdatedSchemas[filename] = specname

	return nil

}

func generatePolymorphicProperty(typename string, name string, property Property) {

	// Open the polymorphic property template
	tmpl, err := template.New("polymorphic-property.template").Funcs(template.FuncMap{
		"convertToGoType": convertTypeToGo,
		"convertToPureGoType": convertTypeToPureGo,
	}).ParseFiles("generate/templates/polymorphic-property.template")

	nameParts := strings.Split(name, "_")

	types := append([]string{}, property.PrimitiveTypes...)
	types = append(types, property.PrimitiveItemTypes...)
	types = append(types, property.ItemTypes...)
	types = append(types, property.Types...)

	packageName, err := packageName(typename, true)
	if err != nil {
		fmt.Printf("Error: Invalid CloudFormation resource %s\n%s\n", typename, err)
		os.Exit(1)
	}

	templateData := struct {
		Name        string
		PackageName string
		Basename    string
		Property    Property
		Types       []string
		TypesJoined string
	}{
		Name:        name,
		PackageName: packageName,
		Basename:    nameParts[0],
		Property:    property,
		Types:       types,
		TypesJoined: conjoin("or", types),
	}

	// Execute the template, writing it to file
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, templateData)
	if err != nil {
		fmt.Printf("Error: Failed to generate polymorphic property %s\n%s\n", name, err)
		os.Exit(1)
	}

	// Format the generated Go file with gofmt
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Printf("Error: Failed to format Go file for resource %s\n%s\n", name, err)
		os.Exit(1)
	}

	// Ensure the package directory exists
	dir := "cloudformation/" + packageName
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}

	// Write the file out
	if err := ioutil.WriteFile(dir+"/"+filename(name), formatted, 0644); err != nil {
		fmt.Printf("Error: Failed to write JSON Schema\n%s\n", err)
		os.Exit(1)
	}

}

// counter is used within the JSON Schema template to determin whether or not
// to put a comma after a JSON resource (i.e. if it's the last element, then no comma)
// see: http://android.wekeepcoding.com/article/10126058/Go+template+remove+the+last+comma+in+range+loop
func counter(length int) func() int {
	i := length
	return func() int {
		i--
		return i
	}
}

func conjoin(conj string, items []string) string {
	if len(items) == 0 {
		return ""
	}
	if len(items) == 1 {
		return items[0]
	}
	if len(items) == 2 { // "a and b" not "a, and b"
		return items[0] + " " + conj + " " + items[1]
	}

	sep := ", "
	pieces := []string{items[0]}
	for _, item := range items[1 : len(items)-1] {
		pieces = append(pieces, sep, item)
	}
	pieces = append(pieces, sep, conj, " ", items[len(items)-1])

	return strings.Join(pieces, "")
}
