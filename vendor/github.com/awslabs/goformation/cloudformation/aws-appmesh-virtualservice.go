package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSAppMeshVirtualService AWS CloudFormation Resource (AWS::AppMesh::VirtualService)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appmesh-virtualservice.html
type AWSAppMeshVirtualService struct {

	// MeshName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appmesh-virtualservice.html#cfn-appmesh-virtualservice-meshname
	MeshName *Value `json:"MeshName,omitempty"`

	// Spec AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appmesh-virtualservice.html#cfn-appmesh-virtualservice-spec
	Spec *AWSAppMeshVirtualService_VirtualServiceSpec `json:"Spec,omitempty"`

	// Tags AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appmesh-virtualservice.html#cfn-appmesh-virtualservice-tags
	Tags []AWSAppMeshVirtualService_TagRef `json:"Tags,omitempty"`

	// VirtualServiceName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appmesh-virtualservice.html#cfn-appmesh-virtualservice-virtualservicename
	VirtualServiceName *Value `json:"VirtualServiceName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshVirtualService) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualService"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSAppMeshVirtualService) MarshalJSON() ([]byte, error) {
	type Properties AWSAppMeshVirtualService
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
func (r *AWSAppMeshVirtualService) UnmarshalJSON(b []byte) error {
	type Properties AWSAppMeshVirtualService
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
		*r = AWSAppMeshVirtualService(*res.Properties)
	}

	return nil
}

// GetAllAWSAppMeshVirtualServiceResources retrieves all AWSAppMeshVirtualService items from an AWS CloudFormation template
func (t *Template) GetAllAWSAppMeshVirtualServiceResources() map[string]AWSAppMeshVirtualService {
	results := map[string]AWSAppMeshVirtualService{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSAppMeshVirtualService:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::AppMesh::VirtualService" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSAppMeshVirtualService{}
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

// GetAWSAppMeshVirtualServiceWithName retrieves all AWSAppMeshVirtualService items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSAppMeshVirtualServiceWithName(name string) (AWSAppMeshVirtualService, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSAppMeshVirtualService:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::AppMesh::VirtualService" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSAppMeshVirtualService{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSAppMeshVirtualService{}, errors.New("resource not found")
}
