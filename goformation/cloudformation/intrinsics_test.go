package cloudformation_test

import (
	"fmt"
	"strings"

	"goformation/v4"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Goformation", func() {
	Context("with a YAML template that contains an intrinsic", func() {

		tests := []struct {
			Name     string
			Input    string
			Expected string
		}{
			{
				Name:     "Ref",
				Input:    `!Ref tvalue`,
				Expected: `{"Ref":"tvalue"}`,
			},
			{
				Name:     "ImportValue",
				Input:    `!ImportValue tvalue`,
				Expected: `{"Fn::ImportValue":"tvalue"}`,
			},
			{
				Name:     "Base64",
				Input:    `!Base64 tvalue`,
				Expected: `{"Fn::Base64":"tvalue"}`,
			},
			{
				Name:     "GetAZs",
				Input:    `!GetAZs tvalue`,
				Expected: `{"Fn::GetAZs":"tvalue"}`,
			},
			{
				Name:     "Sub",
				Input:    `!Sub tvalue`,
				Expected: `{"Fn::Sub":"tvalue"}`,
			},
			{
				Name:     "GetAtt",
				Input:    `!GetAtt object.property`,
				Expected: `{"Fn::GetAtt":["object","property"]}`,
			},
			{
				Name:     "Split",
				Input:    `!Split [d, sss]`,
				Expected: `{"Fn::Split":["d","sss"]}`,
			},
			{
				Name:     "Equals",
				Input:    `!Equals [a, b]`,
				Expected: `{"Fn::Equals":["a","b"]}`,
			},
			{
				Name:     "CIDR",
				Input:    `!Cidr [a, b, c]`,
				Expected: `{"Fn::Cidr":["a","b","c"]}`,
			},
			{
				Name:     "FindInMap",
				Input:    `!FindInMap [a, b, c]`,
				Expected: `{"Fn::FindInMap":["a","b","c"]}`,
			},
			{
				Name:     "If",
				Input:    `!If [a, b, c]`,
				Expected: `{"Fn::If":["a","b","c"]}`,
			},
			{
				Name:     "Join",
				Input:    `!Join [a, [b, c]]`,
				Expected: `{"Fn::Join":["a",["b","c"]]}`,
			},
			{
				Name:     "Select",
				Input:    `!Select [a, [b, c]]`,
				Expected: `{"Fn::Select":["a",["b","c"]]}`,
			},
			{
				Name:     "And",
				Input:    `!And [a, b, c]`,
				Expected: `{"Fn::And":["a","b","c"]}`,
			},
			{
				Name:     "Not",
				Input:    `!Not [a, b, c]`,
				Expected: `{"Fn::Not":["a","b","c"]}`,
			},
			{
				Name:     "Or",
				Input:    `!Or [a, b, c]`,
				Expected: `{"Fn::Or":["a","b","c"]}`,
			},
		}

		for _, test := range tests {
			test := test

			It("should replace "+test.Name+" with the JSON expanded value", func() {
				templateYamlStr := `
Resources:
  TestBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: %s`
				inputTemplate := fmt.Sprintf(templateYamlStr, test.Input)
				template, err := goformation.ParseYAML([]byte(inputTemplate))

				Expect(err).To(BeNil())

				raw, err := template.JSON()
				output := string(raw)
				Expect(err).To(BeNil())
				output = strings.Replace(output, " ", "", -1)
				output = strings.Replace(output, "\n", "", -1)

				templateJSON := fmt.Sprintf(
					`{"Resources":{"TestBucket":{"Type":"AWS::S3::Bucket","Properties":{"BucketName":%s}}}}`,
					test.Expected,
				)
				Expect(output).To(BeEquivalentTo(templateJSON))
			})
		}
	})
})
