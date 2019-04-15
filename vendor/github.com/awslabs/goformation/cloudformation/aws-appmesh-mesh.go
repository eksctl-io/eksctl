package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSAppMeshMesh AWS CloudFormation Resource (AWS::AppMesh::Mesh)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appmesh-mesh.html
type AWSAppMeshMesh struct {

	// MeshName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appmesh-mesh.html#cfn-appmesh-mesh-meshname
	MeshName *Value `json:"MeshName,omitempty"`

	// Spec AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appmesh-mesh.html#cfn-appmesh-mesh-spec
	Spec *AWSAppMeshMesh_MeshSpec `json:"Spec,omitempty"`

	// Tags AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appmesh-mesh.html#cfn-appmesh-mesh-tags
	Tags []AWSAppMeshMesh_TagRef `json:"Tags,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshMesh) AWSCloudFormationType() string {
	return "AWS::AppMesh::Mesh"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSAppMeshMesh) MarshalJSON() ([]byte, error) {
	type Properties AWSAppMeshMesh
	return json.Marshal(&struct {
		Type       string
		Properties Properties
	}{
		Type:       r.AWSCloudFormationType(),
		Properties: (Properties)(*r),
	})
}

// UnmarshalJSON is a custom JSON unmarshalling hook that strips the outer
// AWS CloudFormation resource object, and just keeps the 'Properties' field.
func (r *AWSAppMeshMesh) UnmarshalJSON(b []byte) error {
	type Properties AWSAppMeshMesh
	res := &struct {
		Type       string
		Properties *Properties
	}{}
	if err := json.Unmarshal(b, &res); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return err
	}

	// If the resource has no Properties set, it could be nil
	if res.Properties != nil {
		*r = AWSAppMeshMesh(*res.Properties)
	}

	return nil
}

// GetAllAWSAppMeshMeshResources retrieves all AWSAppMeshMesh items from an AWS CloudFormation template
func (t *Template) GetAllAWSAppMeshMeshResources() map[string]AWSAppMeshMesh {
	results := map[string]AWSAppMeshMesh{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSAppMeshMesh:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::AppMesh::Mesh" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSAppMeshMesh{}
						if err := result.UnmarshalJSON(b); err == nil {
							results[name] = *result
						}
					}
				}
			}
		}
	}
	return results
}

// GetAWSAppMeshMeshWithName retrieves all AWSAppMeshMesh items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSAppMeshMeshWithName(name string) (AWSAppMeshMesh, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSAppMeshMesh:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::AppMesh::Mesh" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSAppMeshMesh{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSAppMeshMesh{}, errors.New("resource not found")
}
