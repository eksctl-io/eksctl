package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSAppMeshVirtualNode AWS CloudFormation Resource (AWS::AppMesh::VirtualNode)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appmesh-virtualnode.html
type AWSAppMeshVirtualNode struct {

	// MeshName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appmesh-virtualnode.html#cfn-appmesh-virtualnode-meshname
	MeshName *Value `json:"MeshName,omitempty"`

	// Spec AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appmesh-virtualnode.html#cfn-appmesh-virtualnode-spec
	Spec *AWSAppMeshVirtualNode_VirtualNodeSpec `json:"Spec,omitempty"`

	// Tags AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appmesh-virtualnode.html#cfn-appmesh-virtualnode-tags
	Tags []AWSAppMeshVirtualNode_TagRef `json:"Tags,omitempty"`

	// VirtualNodeName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appmesh-virtualnode.html#cfn-appmesh-virtualnode-virtualnodename
	VirtualNodeName *Value `json:"VirtualNodeName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshVirtualNode) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualNode"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSAppMeshVirtualNode) MarshalJSON() ([]byte, error) {
	type Properties AWSAppMeshVirtualNode
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
func (r *AWSAppMeshVirtualNode) UnmarshalJSON(b []byte) error {
	type Properties AWSAppMeshVirtualNode
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
		*r = AWSAppMeshVirtualNode(*res.Properties)
	}

	return nil
}

// GetAllAWSAppMeshVirtualNodeResources retrieves all AWSAppMeshVirtualNode items from an AWS CloudFormation template
func (t *Template) GetAllAWSAppMeshVirtualNodeResources() map[string]AWSAppMeshVirtualNode {
	results := map[string]AWSAppMeshVirtualNode{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSAppMeshVirtualNode:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::AppMesh::VirtualNode" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSAppMeshVirtualNode{}
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

// GetAWSAppMeshVirtualNodeWithName retrieves all AWSAppMeshVirtualNode items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSAppMeshVirtualNodeWithName(name string) (AWSAppMeshVirtualNode, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSAppMeshVirtualNode:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::AppMesh::VirtualNode" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSAppMeshVirtualNode{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSAppMeshVirtualNode{}, errors.New("resource not found")
}
