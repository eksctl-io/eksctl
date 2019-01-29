package outputs_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"

	//	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	. "github.com/weaveworks/eksctl/pkg/cfn/outputs"
)

func appendOutput(stack *cfn.Stack, k, v string) {
	stack.Outputs = append(stack.Outputs, &cfn.Output{
		OutputKey:   &k,
		OutputValue: &v,
	})
}

var _ = Describe("CloudFormation stack outputs API", func() {

	It("should handle nil args", func() {
		err := Collect(cfn.Stack{Outputs: nil}, nil, nil)
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("should handle required collectors correctly", func() {
		spec := &api.ClusterConfig{
			VPC: &api.ClusterVPC{},
		}

		{
			stack := cfn.Stack{
				Outputs: nil,
			}
			requiredCollectors := map[string]Collector{
				ClusterVPC: func(v string) error {
					spec.VPC.ID = v
					return nil
				},
			}
			err := Collect(stack, requiredCollectors, nil)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("no ouput \"" + ClusterVPC + "\""))
		}

		{
			stack := cfn.Stack{}
			stack.SetStackName("foo")

			appendOutput(&stack, ClusterVPC, "vpc-123")

			{
				requiredCollectors := map[string]Collector{
					ClusterVPC: func(v string) error {
						spec.VPC.ID = v
						return nil
					},
					ClusterSecurityGroup: func(v string) error {
						spec.VPC.SecurityGroup = v
						return nil
					},
				}

				err := Collect(stack, requiredCollectors, nil)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(Equal("no ouput \"" + ClusterSecurityGroup + "\" in stack \"foo\""))
			}

			appendOutput(&stack, ClusterSecurityGroup, "sg-123")

			{
				requiredCollectors := map[string]Collector{
					ClusterVPC: func(v string) error {
						spec.VPC.ID = v
						return nil
					},
					ClusterSecurityGroup: func(v string) error {
						spec.VPC.SecurityGroup = v
						return nil
					},
				}

				err := Collect(stack, requiredCollectors, nil)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(spec.VPC.ID).To(Equal("vpc-123"))
				Expect(spec.VPC.SecurityGroup).To(Equal("sg-123"))
			}
		}

	})

	It("should handle required and optional collectors correctly", func() {
		spec := &api.ClusterConfig{
			VPC: &api.ClusterVPC{},
		}

		{
			stack := cfn.Stack{
				Outputs: nil,
			}
			requiredCollectors := map[string]Collector{
				ClusterVPC: func(v string) error {
					spec.VPC.ID = v
					return nil
				},
			}
			err := Collect(stack, requiredCollectors, nil)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("no ouput \"" + ClusterVPC + "\""))
		}

		{
			stack := cfn.Stack{}
			stack.SetStackName("foo")

			appendOutput(&stack, ClusterVPC, "vpc-123")
			appendOutput(&stack, ClusterSecurityGroup, "sg-123")
			appendOutput(&stack, "test1", "")

			{
				requiredCollectors := map[string]Collector{
					ClusterVPC: func(v string) error {
						spec.VPC.ID = v
						return nil
					},
					ClusterSecurityGroup: func(v string) error {
						spec.VPC.SecurityGroup = v
						return nil
					},
				}

				test1 := false
				test2 := false

				optionalCollectors := map[string]Collector{
					"test1": func(_ string) error {
						test1 = true
						return nil
					},
					"test2": func(_ string) error {
						test2 = true
						return nil
					},
				}

				err := Collect(stack, requiredCollectors, optionalCollectors)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(spec.VPC.ID).To(Equal("vpc-123"))
				Expect(spec.VPC.SecurityGroup).To(Equal("sg-123"))

				Expect(test1).To(BeTrue())
				Expect(test2).To(BeFalse())
			}

			{
				appendOutput(&stack, "test3", "")
				appendOutput(&stack, "test4", "")

				test3 := false
				test4 := false

				optionalCollectors := map[string]Collector{
					"test3": func(_ string) error {
						test3 = true
						return nil
					},
					"test4": func(_ string) error {
						test4 = true
						return nil
					},
				}

				err := Collect(stack, nil, optionalCollectors)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(test3).To(BeTrue())
				Expect(test4).To(BeTrue())
			}
		}
	})
})
