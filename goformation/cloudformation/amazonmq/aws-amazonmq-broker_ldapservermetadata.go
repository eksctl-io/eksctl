package amazonmq

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Broker_LdapServerMetadata AWS CloudFormation Resource (AWS::AmazonMQ::Broker.LdapServerMetadata)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-ldapservermetadata.html
type Broker_LdapServerMetadata struct {

	// Hosts AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-ldapservermetadata.html#cfn-amazonmq-broker-ldapservermetadata-hosts
	Hosts *types.Value `json:"Hosts,omitempty"`

	// RoleBase AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-ldapservermetadata.html#cfn-amazonmq-broker-ldapservermetadata-rolebase
	RoleBase *types.Value `json:"RoleBase,omitempty"`

	// RoleName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-ldapservermetadata.html#cfn-amazonmq-broker-ldapservermetadata-rolename
	RoleName *types.Value `json:"RoleName,omitempty"`

	// RoleSearchMatching AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-ldapservermetadata.html#cfn-amazonmq-broker-ldapservermetadata-rolesearchmatching
	RoleSearchMatching *types.Value `json:"RoleSearchMatching,omitempty"`

	// RoleSearchSubtree AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-ldapservermetadata.html#cfn-amazonmq-broker-ldapservermetadata-rolesearchsubtree
	RoleSearchSubtree *types.Value `json:"RoleSearchSubtree,omitempty"`

	// ServiceAccountPassword AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-ldapservermetadata.html#cfn-amazonmq-broker-ldapservermetadata-serviceaccountpassword
	ServiceAccountPassword *types.Value `json:"ServiceAccountPassword,omitempty"`

	// ServiceAccountUsername AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-ldapservermetadata.html#cfn-amazonmq-broker-ldapservermetadata-serviceaccountusername
	ServiceAccountUsername *types.Value `json:"ServiceAccountUsername,omitempty"`

	// UserBase AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-ldapservermetadata.html#cfn-amazonmq-broker-ldapservermetadata-userbase
	UserBase *types.Value `json:"UserBase,omitempty"`

	// UserRoleName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-ldapservermetadata.html#cfn-amazonmq-broker-ldapservermetadata-userrolename
	UserRoleName *types.Value `json:"UserRoleName,omitempty"`

	// UserSearchMatching AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-ldapservermetadata.html#cfn-amazonmq-broker-ldapservermetadata-usersearchmatching
	UserSearchMatching *types.Value `json:"UserSearchMatching,omitempty"`

	// UserSearchSubtree AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-ldapservermetadata.html#cfn-amazonmq-broker-ldapservermetadata-usersearchsubtree
	UserSearchSubtree *types.Value `json:"UserSearchSubtree,omitempty"`

	// AWSCloudFormationDeletionPolicy represents a CloudFormation DeletionPolicy
	AWSCloudFormationDeletionPolicy policies.DeletionPolicy `json:"-"`

	// AWSCloudFormationUpdateReplacePolicy represents a CloudFormation UpdateReplacePolicy
	AWSCloudFormationUpdateReplacePolicy policies.UpdateReplacePolicy `json:"-"`

	// AWSCloudFormationDependsOn stores the logical ID of the resources to be created before this resource
	AWSCloudFormationDependsOn []string `json:"-"`

	// AWSCloudFormationMetadata stores structured data associated with this resource
	AWSCloudFormationMetadata map[string]interface{} `json:"-"`

	// AWSCloudFormationCondition stores the logical ID of the condition that must be satisfied for this resource to be created
	AWSCloudFormationCondition string `json:"-"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *Broker_LdapServerMetadata) AWSCloudFormationType() string {
	return "AWS::AmazonMQ::Broker.LdapServerMetadata"
}
