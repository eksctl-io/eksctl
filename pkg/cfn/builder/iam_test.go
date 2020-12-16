package builder_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	cft "github.com/weaveworks/eksctl/pkg/cfn/template"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"

	. "github.com/weaveworks/eksctl/pkg/cfn/template/matchers"

	. "github.com/weaveworks/eksctl/pkg/cfn/builder"
)

var _ = Describe("template builder for IAM", func() {
	Describe("IAMServiceAccount", func() {
		var (
			oidc *iamoidc.OpenIDConnectManager
			cfg  *api.ClusterConfig
			err  error
		)

		BeforeEach(func() {
			oidc, err = iamoidc.NewOpenIDConnectManager(nil, "456123987123", "https://oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E", "aws")
			Expect(err).ToNot(HaveOccurred())

			oidc.ProviderARN = "arn:aws:iam::456123987123:oidc-provider/oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E"

			cfg = api.NewClusterConfig()

			cfg.IAM.WithOIDC = api.Enabled()
			cfg.IAM.ServiceAccounts = []*api.ClusterIAMServiceAccount{}
		})

		It("can construct an iamserviceaccount addon template with one managed policy", func() {
			serviceAccount := &api.ClusterIAMServiceAccount{}

			serviceAccount.Name = "sa-1"

			serviceAccount.AttachPolicyARNs = []string{"arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess"}

			appendServiceAccountToClusterConfig(cfg, serviceAccount)

			rs := NewIAMServiceAccountResourceSet(serviceAccount, oidc)

			templateBody := []byte{}

			Expect(rs).To(RenderWithoutErrors(&templateBody))

			t := cft.NewTemplate()

			Expect(t).To(LoadBytesWithoutErrors(templateBody))

			Expect(t.Description).To(Equal("IAM role for serviceaccount \"default/sa-1\" [created and managed by eksctl]"))

			Expect(t.Resources).To(HaveLen(1))
			Expect(t.Outputs).To(HaveLen(1))

			Expect(t).To(HaveResource("Role1", "AWS::IAM::Role"))

			Expect(t).To(HaveResourceWithPropertyValue("Role1", "AssumeRolePolicyDocument", expectedServiceAccountAssumeRolePolicyDocument))
			Expect(t).To(HaveResourceWithPropertyValue("Role1", "ManagedPolicyArns", `[
			"arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess"
		]`))

			Expect(t).To(HaveOutputWithValue("Role1", `{ "Fn::GetAtt": "Role1.Arn" }`))
		})

		It("can constuct an iamserviceaccount addon template with one inline policy", func() {
			serviceAccount := &api.ClusterIAMServiceAccount{}

			serviceAccount.Name = "sa-1"

			serviceAccount.AttachPolicy = cft.MakePolicyDocument(
				cft.MapOfInterfaces{
					"Effect": "Allow",
					"Action": []string{
						"s3:Get*",
					},
					"Resource": "*",
				},
			)

			appendServiceAccountToClusterConfig(cfg, serviceAccount)

			rs := NewIAMServiceAccountResourceSet(serviceAccount, oidc)

			templateBody := []byte{}

			Expect(rs).To(RenderWithoutErrors(&templateBody))

			t := cft.NewTemplate()

			Expect(t).To(LoadBytesWithoutErrors(templateBody))

			Expect(t.Description).To(Equal("IAM role for serviceaccount \"default/sa-1\" [created and managed by eksctl]"))

			Expect(t.Resources).To(HaveLen(2))
			Expect(t.Outputs).To(HaveLen(1))

			Expect(t).To(HaveResource("Role1", "AWS::IAM::Role"))
			Expect(t).To(HaveResource("Policy1", "AWS::IAM::Policy"))

			Expect(t).ToNot(HaveResourceWithProperties("Role1", "ManagedPolicyArns"))

			Expect(t).To(HaveResourceWithPropertyValue("Role1", "AssumeRolePolicyDocument", expectedServiceAccountAssumeRolePolicyDocument))
			Expect(t).To(HaveResourceWithPropertyValue("Policy1", "PolicyName", `{ "Fn::Sub": "${AWS::StackName}-Policy1" }`))
			Expect(t).To(HaveResourceWithPropertyValue("Policy1", "PolicyDocument", `{
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Action": [
                        "s3:Get*"
                    ],
                    "Resource": "*"
                }
            ]
        }`))

			Expect(t).To(HaveOutputWithValue("Role1", `{ "Fn::GetAtt": "Role1.Arn" }`))
		})

		It("can constuct an iamserviceaccount addon template with two managed policies and one inline policy", func() {
			serviceAccount := &api.ClusterIAMServiceAccount{}

			serviceAccount.Name = "sa-1"

			serviceAccount.AttachPolicyARNs = []string{
				"arn:aws:iam::aws:policy/AmazonDynamoDBReadOnlyAccess",
				"arn:aws:iam::aws:policy/AmazonElastiCacheFullAccess",
			}

			serviceAccount.AttachPolicy = cft.MakePolicyDocument(
				cft.MapOfInterfaces{
					"Effect": "Allow",
					"Action": []string{
						"s3:Get*",
					},
					"Resource": "*",
				},
			)

			appendServiceAccountToClusterConfig(cfg, serviceAccount)

			rs := NewIAMServiceAccountResourceSet(serviceAccount, oidc)

			templateBody := []byte{}

			Expect(rs).To(RenderWithoutErrors(&templateBody))

			t := cft.NewTemplate()

			Expect(t).To(LoadBytesWithoutErrors(templateBody))

			Expect(t.Description).To(Equal("IAM role for serviceaccount \"default/sa-1\" [created and managed by eksctl]"))

			Expect(t.Resources).To(HaveLen(2))
			Expect(t.Outputs).To(HaveLen(1))

			Expect(t).To(HaveResource("Role1", "AWS::IAM::Role"))
			Expect(t).To(HaveResource("Policy1", "AWS::IAM::Policy"))

			Expect(t).To(HaveResourceWithPropertyValue("Role1", "ManagedPolicyArns", `[
			"arn:aws:iam::aws:policy/AmazonDynamoDBReadOnlyAccess",
			"arn:aws:iam::aws:policy/AmazonElastiCacheFullAccess"
		]`))

			Expect(t).To(HaveResourceWithPropertyValue("Role1", "AssumeRolePolicyDocument", expectedServiceAccountAssumeRolePolicyDocument))
			Expect(t).To(HaveResourceWithPropertyValue("Policy1", "PolicyName", `{ "Fn::Sub": "${AWS::StackName}-Policy1" }`))
			Expect(t).To(HaveResourceWithPropertyValue("Policy1", "PolicyDocument", `{
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Action": [
                        "s3:Get*"
                    ],
                    "Resource": "*"
                }
            ]
        }`))

			Expect(t).To(HaveOutputWithValue("Role1", `{ "Fn::GetAtt": "Role1.Arn" }`))
		})

		It("can construct an iamserviceaccount addon template with one managed policy and a permissions boundary", func() {
			serviceAccount := &api.ClusterIAMServiceAccount{}

			serviceAccount.Name = "sa-1"

			serviceAccount.AttachPolicyARNs = []string{"arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess"}

			serviceAccount.PermissionsBoundary = "arn:aws:iam::aws:policy/policy/boundary"

			appendServiceAccountToClusterConfig(cfg, serviceAccount)

			rs := NewIAMServiceAccountResourceSet(serviceAccount, oidc)

			templateBody := []byte{}

			Expect(rs).To(RenderWithoutErrors(&templateBody))

			t := cft.NewTemplate()

			Expect(t).To(LoadBytesWithoutErrors(templateBody))

			Expect(t.Description).To(Equal("IAM role for serviceaccount \"default/sa-1\" [created and managed by eksctl]"))

			Expect(t.Resources).To(HaveLen(1))
			Expect(t.Outputs).To(HaveLen(1))

			Expect(t).To(HaveResource("Role1", "AWS::IAM::Role"))

			Expect(t).To(HaveResourceWithPropertyValue("Role1", "AssumeRolePolicyDocument", expectedServiceAccountAssumeRolePolicyDocument))
			Expect(t).To(HaveResourceWithPropertyValue("Role1", "ManagedPolicyArns", `[
			"arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess"
		]`))
			Expect(t).To(HaveResourceWithPropertyValue("Role1", "PermissionsBoundary", `"arn:aws:iam::aws:policy/policy/boundary"`))

			Expect(t).To(HaveOutputWithValue("Role1", `{ "Fn::GetAtt": "Role1.Arn" }`))
		})

		It("can parse an iamserviceaccount addon template", func() {
			t := cft.NewTemplate()

			Expect(t).To(LoadFileWithoutErrors("../template/testdata/addon-example-1.json"))

			Expect(t.Description).To(Equal("IAM role for serviceaccount \"default/sa-1\" [created and managed by eksctl]"))

			Expect(t.Resources).To(HaveLen(1))
			Expect(t.Outputs).To(HaveLen(1))

			Expect(t).To(HaveResource("Role1", "AWS::IAM::Role"))

			Expect(t).To(HaveResourceWithPropertyValue("Role1", "ManagedPolicyArns", `[ "arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess" ]`))

			Expect(t).To(HaveResourceWithPropertyValue("Role1", "AssumeRolePolicyDocument", expectedServiceAccountAssumeRolePolicyDocument))

			Expect(t).To(HaveOutputWithValue("Role1", `{ "Fn::GetAtt": "Role1.Arn" }`))
		})
	})

	Describe("IAMRole", func() {
		var (
			oidc *iamoidc.OpenIDConnectManager
			cfg  *api.ClusterConfig
			err  error
		)

		BeforeEach(func() {
			oidc, err = iamoidc.NewOpenIDConnectManager(nil, "456123987123", "https://oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E", "aws")
			Expect(err).ToNot(HaveOccurred())

			oidc.ProviderARN = "arn:aws:iam::456123987123:oidc-provider/oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E"

			cfg = api.NewClusterConfig()

			cfg.IAM.WithOIDC = api.Enabled()
		})

		It("can construct an iamrole template with attachPolicyARNs", func() {
			arns := []string{"arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess"}

			rs := NewIAMRoleResourceSetWithAttachPolicyARNs("VPC-addon", "", "", arns, oidc)

			templateBody := []byte{}

			Expect(rs).To(RenderWithoutErrors(&templateBody))

			t := cft.NewTemplate()

			Expect(t).To(LoadBytesWithoutErrors(templateBody))

			Expect(t.Description).To(Equal("IAM role for \"VPC-addon\" [created and managed by eksctl]"))

			Expect(t.Resources).To(HaveLen(1))
			Expect(t.Outputs).To(HaveLen(1))

			Expect(t).To(HaveResource("Role1", "AWS::IAM::Role"))

			Expect(t).To(HaveResourceWithPropertyValue("Role1", "AssumeRolePolicyDocument", expectedAssumeRolePolicyDocument))
			Expect(t).To(HaveResourceWithPropertyValue("Role1", "ManagedPolicyArns", `[
			"arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess"
		]`))

			Expect(t).To(HaveOutputWithValue("Role1", `{ "Fn::GetAtt": "Role1.Arn" }`))
		})

		It("can construct an iamrole template with attachPolicy", func() {
			attachPolicy := cft.MakePolicyDocument(
				cft.MapOfInterfaces{
					"Effect": "Allow",
					"Action": []string{
						"s3:Get*",
					},
					"Resource": "*",
				},
			)

			rs := NewIAMRoleResourceSetWithAttachPolicy("VPC-addon", "", "", attachPolicy, oidc)

			templateBody := []byte{}

			Expect(rs).To(RenderWithoutErrors(&templateBody))

			t := cft.NewTemplate()

			Expect(t).To(LoadBytesWithoutErrors(templateBody))

			Expect(t.Description).To(Equal("IAM role for \"VPC-addon\" [created and managed by eksctl]"))

			Expect(t.Resources).To(HaveLen(2))
			Expect(t.Outputs).To(HaveLen(1))

			Expect(t).To(HaveResource("Role1", "AWS::IAM::Role"))
			Expect(t).To(HaveResource("Policy1", "AWS::IAM::Policy"))

			Expect(t).ToNot(HaveResourceWithProperties("Role1", "ManagedPolicyArns"))

			Expect(t).To(HaveResourceWithPropertyValue("Role1", "AssumeRolePolicyDocument", expectedAssumeRolePolicyDocument))
			Expect(t).To(HaveResourceWithPropertyValue("Policy1", "PolicyName", `{ "Fn::Sub": "${AWS::StackName}-Policy1" }`))
			Expect(t).To(HaveResourceWithPropertyValue("Policy1", "PolicyDocument", `{
		   "Version": "2012-10-17",
		   "Statement": [
		       {
		           "Effect": "Allow",
		           "Action": [
		               "s3:Get*"
		           ],
		           "Resource": "*"
		       }
		   ]
		}`))

			Expect(t).To(HaveOutputWithValue("Role1", `{ "Fn::GetAtt": "Role1.Arn" }`))
		})
	})
})

func appendServiceAccountToClusterConfig(cfg *api.ClusterConfig, serviceAccount *api.ClusterIAMServiceAccount) {
	cfg.IAM.ServiceAccounts = append(cfg.IAM.ServiceAccounts, serviceAccount)

	api.SetClusterConfigDefaults(cfg)
	err := api.ValidateClusterConfig(cfg)
	Expect(err).ToNot(HaveOccurred())
}

const expectedServiceAccountAssumeRolePolicyDocument = `{
	"Statement": [
	  {
		"Action": [
		  "sts:AssumeRoleWithWebIdentity"
		],
		"Condition": {
		  "StringEquals": {
			"oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E:aud": "sts.amazonaws.com",
			"oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E:sub": "system:serviceaccount:default:sa-1"
		  }
		},
		"Effect": "Allow",
		"Principal": {
		  "Federated": "arn:aws:iam::456123987123:oidc-provider/oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E"
		}
	  }
	],
	"Version": "2012-10-17"
}`

const expectedAssumeRolePolicyDocument = `{
	"Statement": [
	  {
		"Action": [
		  "sts:AssumeRoleWithWebIdentity"
		],
		"Condition": {
		  "StringEquals": {
			"oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E:aud": "sts.amazonaws.com"
		  }
		},
		"Effect": "Allow",
		"Principal": {
		  "Federated": "arn:aws:iam::456123987123:oidc-provider/oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E"
		}
	  }
	],
	"Version": "2012-10-17"
}`
