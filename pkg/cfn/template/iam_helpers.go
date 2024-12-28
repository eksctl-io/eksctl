package template

import (
	gfn "goformation/v4/cloudformation/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

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

// MakeAssumeRolePolicyDocumentForServices constructs a trust policy for given services.
func MakeAssumeRolePolicyDocumentForServices(services ...*gfn.Value) MapOfInterfaces {
	return MakeAssumeRolePolicyDocumentWithAction("", services...)
}

// MakeAssumeRolePolicyDocumentWithAction constructs a trust policy for given services and action.
func MakeAssumeRolePolicyDocumentWithAction(action string, services ...*gfn.Value) MapOfInterfaces {
	actions := []string{"sts:AssumeRole"}
	if action != "" {
		actions = append(actions, action)
	}
	return MakePolicyDocument(MapOfInterfaces{
		"Effect": "Allow",
		"Action": actions,
		"Principal": map[string][]*gfn.Value{
			"Service": services,
		},
	})
}

// MakeAssumeRolePolicyDocumentForServicesWithConditions constructs a trust policy for given services with given conditions.
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

// MakeAssumeRolePolicyDocumentForPodIdentity constructs a trust policy for roles used in pods identity associations
func MakeAssumeRolePolicyDocumentForPodIdentity() MapOfInterfaces {
	return MakePolicyDocument(MapOfInterfaces{
		"Effect": "Allow",
		"Action": []string{
			"sts:AssumeRole",
			"sts:TagSession",
		},
		"Principal": map[string]string{
			"Service": api.EKSServicePrincipal,
		},
	})
}

// MakeAssumeRolePolicyDocumentForServices constructs a trust policy for given services with given conditions and extra actions
func MakeAssumeRolePolicyDocumentForServicesWithConditionsAndActions(condition MapOfInterfaces, extraActions []string, services ...*gfn.Value) MapOfInterfaces {
	return MakePolicyDocument(MapOfInterfaces{
		"Effect":    "Allow",
		"Action":    append([]string{"sts:AssumeRole"}, extraActions...),
		"Condition": condition,
		"Principal": map[string][]*gfn.Value{
			"Service": services,
		},
	})
}
