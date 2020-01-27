package template

// IAMPolicy represents a CloudFormation AWS::IAM::Policy resource
type IAMPolicy struct {
	PolicyName *Value `json:",omitempty"`

	Roles          []*Value        `json:",omitempty"`
	PolicyDocument MapOfInterfaces `json:",omitempty"`
}

// Type will return the full type name for the resource
func (r *IAMPolicy) Type() string {
	return "AWS::IAM::Policy"
}

// Properties will return the properties of the resource
func (r *IAMPolicy) Properties() interface{} {
	return r
}

// IAMRole represents a CloudFormation AWS::IAM::Role resource
type IAMRole struct {
	RoleName string `json:",omitempty"`

	Path string `json:",omitempty"`

	AssumeRolePolicyDocument MapOfInterfaces `json:",omitempty"`
	ManagedPolicyArns        []string        `json:",omitempty"`
	PermissionsBoundary      string          `json:",omitempty"`
}

// Type will return the full type name for the resource
func (r *IAMRole) Type() string {
	return "AWS::IAM::Role"
}

// Properties will return the properties of the resource
func (r *IAMRole) Properties() interface{} {
	return r
}
