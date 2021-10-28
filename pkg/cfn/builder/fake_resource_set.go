package builder

import (
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
)

type FakeResourceSet struct {
	rs *resourceSet
}

func (f *FakeResourceSet) AddAllResources() error {
	return nil
}

// WithIAM states, if IAM roles will be created or not
func (f *FakeResourceSet) WithIAM() bool {
	return f.rs.withIAM
}

// WithNamedIAM states, if specifically named IAM roles will be created or not
func (f *FakeResourceSet) WithNamedIAM() bool {
	return f.rs.withNamedIAM
}

// RenderJSON returns the rendered JSON
func (f *FakeResourceSet) RenderJSON() ([]byte, error) {
	return f.rs.renderJSON()
}

// GetAllOutputs collects all outputs of the cluster
func (f *FakeResourceSet) GetAllOutputs(stack cfn.Stack) error {
	return f.rs.GetAllOutputs(stack)
}

func NewFakeResourceSet() ResourceSet {
	return &FakeResourceSet{
		rs: newResourceSet(),
	}
}
