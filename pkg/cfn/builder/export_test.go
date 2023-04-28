package builder

import (
	gfn "github.com/weaveworks/goformation/v4/cloudformation"
)

func NewRS() *resourceSet {
	return newResourceSet()
}

func GetTemplate(rs *resourceSet) *gfn.Template {
	return rs.template
}

func GetIAMRoleName() string {
	return cfnIAMInstanceRoleName
}
