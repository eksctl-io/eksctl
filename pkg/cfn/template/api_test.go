package template

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CloudFormation template", func() {
	It("can create and render a minimal template", func() {
		t := NewTemplate()
		t.Description = "a template"

		roleRef := t.NewResource("aRole", &IAMRole{
			RoleName:          "foo",
			ManagedPolicyArns: []string{"abc"},
		})

		jsRoleRef, err := roleRef.MarshalJSON()
		Expect(err).ToNot(HaveOccurred())
		Expect(jsRoleRef).To(MatchJSON(`{"Ref": "aRole"}`))

		policyRef := t.NewResource("aPolicy", &IAMPolicy{
			PolicyName: NewString("foo"),
			Roles:      MakeSlice(roleRef),
		})

		jsPolicyRef, err := policyRef.MarshalJSON()
		Expect(err).ToNot(HaveOccurred())
		Expect(jsPolicyRef).To(MatchJSON(`{"Ref": "aPolicy"}`))

		js, err := t.RenderJSON()
		Expect(err).ToNot(HaveOccurred())
		Expect(js).To(MatchJSON([]byte(templateExample1)))
	})

	It("can parse a template", func() {
		t := NewTemplate()
		err := t.LoadJSON([]byte(templateExample1))
		Expect(err).ToNot(HaveOccurred())

		Expect(t.Description).To(Equal("a template"))

		Expect(t.Resources["aPolicy"].Type).To(Equal("AWS::IAM::Policy"))
		Expect(t.Resources["aPolicy"].Properties).To(HaveKeyWithValue("PolicyName", "foo"))

		Expect(t.Resources["aRole"].Type).To(Equal("AWS::IAM::Role"))
		Expect(t.Resources["aRole"].Properties).To(HaveKeyWithValue("RoleName", "foo"))
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
	}
}`
