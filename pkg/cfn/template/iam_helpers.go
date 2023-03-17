package template

import gfn "github.com/weaveworks/goformation/v4/cloudformation/types"

// AttachPolicy attaches the specified policy document
func (t *Template) AttachPolicy(name string, refRole *Value, policyDoc MapOfInterfaces) {
	t.NewResource(name, &IAMPolicy{
		PolicyName:     MakeName(name),
		Roles:          MakeSlice(refRole),
		PolicyDocument: policyDoc,
	})
}

// MakePolicyDocument constructs a policy with given statements
func MakePolicyDocument(statements ...MapOfInterfaces) MapOfInterfaces {
	return MapOfInterfaces{
		"Version":   "2012-10-17",
		"Statement": statements,
	}
}

// MakeAssumeRolePolicyDocumentForServices constructs a trust policy for given services
func MakeAssumeRolePolicyDocumentForServices(services ...*gfn.Value) MapOfInterfaces {
	return MakePolicyDocument(MapOfInterfaces{
		"Effect": "Allow",
		"Action": []string{"sts:AssumeRole"},
		"Principal": map[string][]*gfn.Value{
			"Service": services,
		},
	})
}

// MakeAssumeRolePolicyDocumentForServices constructs a trust policy for given services with given conditions
func MakeAssumeRolePolicyDocumentForServicesWithConditions(condition MapOfInterfaces, services ...*gfn.Value) MapOfInterfaces {
	return MakePolicyDocument(MapOfInterfaces{
		"Effect":    "Allow",
		"Action":    []string{"sts:AssumeRole"},
		"Condition": condition,
		"Principal": map[string][]*gfn.Value{
			"Service": services,
		},
	})
}

// MakeAssumeRoleWithWebIdentityPolicyDocument constructs a trust policy for given a web identity priovider with given conditions
func MakeAssumeRoleWithWebIdentityPolicyDocument(providerARN string, condition MapOfInterfaces) MapOfInterfaces {
	return MakePolicyDocument(MapOfInterfaces{
		"Effect": "Allow",
		"Action": []string{"sts:AssumeRoleWithWebIdentity"},
		"Principal": map[string]string{
			"Federated": providerARN,
		},
		"Condition": condition,
	})
}
