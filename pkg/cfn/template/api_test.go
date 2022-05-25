package template_test

import (
	"fmt"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/pkg/cfn/template"

	. "github.com/weaveworks/eksctl/pkg/cfn/template/matchers"
)

var _ = Describe("CloudFormation template", func() {
	It("can create and render a minimal template", func() {
		t := NewTemplate()
		t.Description = "a template"

		roleRef := t.NewResource("aRole", &IAMRole{
			RoleName:          "foo",
			ManagedPolicyArns: []interface{}{"abc"},
		})

		t.Outputs["aRole"] = Output{
			Value: MakeFnGetAttString("aRole.Arn"),
			Export: &OutputExport{
				Name: MakeFnSubString(fmt.Sprintf("${%s}::aRole", StackName)),
			},
		}

		t.Outputs["foo"] = Output{
			Value: NewString("bar"),
		}

		jsRoleRef, err := roleRef.MarshalJSON()
		Expect(err).NotTo(HaveOccurred())
		Expect(jsRoleRef).To(MatchJSON(`{"Ref": "aRole"}`))

		policyRef := t.NewResource("aPolicy", &IAMPolicy{
			PolicyName: NewString("foo"),
			Roles:      MakeSlice(roleRef),
		})

		jsPolicyRef, err := policyRef.MarshalJSON()
		Expect(err).NotTo(HaveOccurred())
		Expect(jsPolicyRef).To(MatchJSON(`{"Ref": "aPolicy"}`))

		js, err := t.RenderJSON()
		Expect(err).NotTo(HaveOccurred())
		Expect(js).To(MatchJSON([]byte(templateExample1)))
	})

	It("can parse a template", func() {
		t := NewTemplate()

		Expect(t).To(LoadStringWithoutErrors(templateExample1))

		Expect(t.Description).To(Equal("a template"))

		Expect(t).To(HaveResource("aPolicy", "AWS::IAM::Policy"))
		Expect(t).To(HaveResourceWithPropertyValue("aPolicy", "PolicyName", `"foo"`))

		Expect(t).To(HaveResource("aRole", "AWS::IAM::Role"))

		Expect(t).To(HaveResourceWithPropertyValue("aRole", "RoleName", `"foo"`))

		Expect(t).NotTo(HaveResourceWithPropertyValue("aRole", "RoleName", `"bar"`))
		Expect(t).NotTo(HaveResource("aRole", "AWS::Foo::Bar"))
		Expect(t).NotTo(HaveResource("foo", "*"))

		Expect(t).To(HaveOutputs("aRole"))
		Expect(t).NotTo(HaveOutputs("foo", "bar"))

		Expect(t).To(HaveOutputWithValue("aRole", `{ "Fn::GetAtt": "aRole.Arn" }`))
		Expect(t).To(HaveOutputExportedAs("aRole", `{ "Fn::Sub": "${AWS::StackName}::aRole" }`))

		Expect(t).NotTo(HaveOutputExportedAs("aRole", `{}`))
		Expect(t).NotTo(HaveOutputExportedAs("aRole", `{ "Fn::GetAtt": "aRole.Arn" }`))
		Expect(t).NotTo(HaveOutputExportedAs("foo", `{ "Fn::Sub": "${AWS::StackName}::aRole" }`))
	})

	It("can load multiple real templates", func() {
		examples, err := filepath.Glob("testdata/*.json")
		Expect(err).NotTo(HaveOccurred())
		for _, example := range examples {
			Expect(NewTemplate()).To(LoadFileWithoutErrors(example))
		}
	})
})

const templateExample1 = `{
	"AWSTemplateFormatVersion": "2010-09-09",
	"Description": "a template",
	"Resources": {
		"aPolicy": {
		  "Type": "AWS::IAM::Policy",
		  "Properties": {
			"PolicyName": "foo",
			"Roles": [
			  { "Ref": "aRole" }
			]
		  }
		},
		"aRole": {
		  "Type": "AWS::IAM::Role",
		  "Properties": {
			"RoleName": "foo",
			"ManagedPolicyArns": [ "abc" ]
		  }
		}
	},
	"Outputs": {
		"foo": { "Value": "bar" },
		"aRole": {
			"Value": { "Fn::GetAtt": "aRole.Arn" },
			"Export": { "Name": { "Fn::Sub": "${AWS::StackName}::aRole" } }
		}
	}
}`
