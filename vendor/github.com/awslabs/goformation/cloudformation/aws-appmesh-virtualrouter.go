package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSAppMeshVirtualRouter AWS CloudFormation Resource (AWS::AppMesh::VirtualRouter)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appmesh-virtualrouter.html
type AWSAppMeshVirtualRouter struct {

	// MeshName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appmesh-virtualrouter.html#cfn-appmesh-virtualrouter-meshname
	MeshName *Value `json:"MeshName,omitempty"`

	// Spec AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appmesh-virtualrouter.html#cfn-appmesh-virtualrouter-spec
	Spec *AWSAppMeshVirtualRouter_VirtualRouterSpec `json:"Spec,omitempty"`

	// Tags AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appmesh-virtualrouter.html#cfn-appmesh-virtualrouter-tags
	Tags []AWSAppMeshVirtualRouter_TagRef `json:"Tags,omitempty"`

	// VirtualRouterName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appmesh-virtualrouter.html#cfn-appmesh-virtualrouter-virtualroutername
	VirtualRouterName *Value `json:"VirtualRouterName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshVirtualRouter) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualRouter"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSAppMeshVirtualRouter) MarshalJSON() ([]byte, error) {
	type Properties AWSAppMeshVirtualRouter
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
func (r *AWSAppMeshVirtualRouter) UnmarshalJSON(b []byte) error {
	type Properties AWSAppMeshVirtualRouter
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
		*r = AWSAppMeshVirtualRouter(*res.Properties)
	}

	return nil
}

// GetAllAWSAppMeshVirtualRouterResources retrieves all AWSAppMeshVirtualRouter items from an AWS CloudFormation template
func (t *Template) GetAllAWSAppMeshVirtualRouterResources() map[string]AWSAppMeshVirtualRouter {
	results := map[string]AWSAppMeshVirtualRouter{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSAppMeshVirtualRouter:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::AppMesh::VirtualRouter" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSAppMeshVirtualRouter{}
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

// GetAWSAppMeshVirtualRouterWithName retrieves all AWSAppMeshVirtualRouter items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSAppMeshVirtualRouterWithName(name string) (AWSAppMeshVirtualRouter, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSAppMeshVirtualRouter:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::AppMesh::VirtualRouter" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSAppMeshVirtualRouter{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSAppMeshVirtualRouter{}, errors.New("resource not found")
}
