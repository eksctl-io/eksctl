package ami_test

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/ami"
)

var _ = Describe("AMI Static Resolution", func() {
	type ResolveCase struct {
		Region       string
		InstanceType string
		ExpectedAMI  string
		ExpectError  bool
	}

	DescribeTable("When resolving an AMI using the default resolvers",
		func(c ResolveCase) {
			actualAmi, err := ami.Resolve(c.Region, c.InstanceType)
			Expect(actualAmi).Should(Equal(c.ExpectedAMI))
			if c.ExpectError {
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(MatchError(ami.NewErrFailedResolution(c.Region, c.InstanceType)))
				errorType := reflect.TypeOf(err).Elem().Name()
				Expect(errorType).To(Equal("ErrFailedResolution"))
			} else {
				Expect(err).ShouldNot(HaveOccurred())
			}
		},
		Entry("with non-gpu instance and us-west-2", ResolveCase{
			Region:       "us-west-2",
			InstanceType: "t2.medium",
			ExpectedAMI:  "ami-0a54c984b9f908c81",
			ExpectError:  false,
		}),
		Entry("with non-gpu instance and us-east-1", ResolveCase{
			Region:       "us-east-1",
			InstanceType: "t2.medium",
			ExpectedAMI:  "ami-0440e4f6b9713faf6",
			ExpectError:  false,
		}),
		Entry("with non-gpu instance and eu-west-1", ResolveCase{
			Region:       "eu-west-1",
			InstanceType: "t2.medium",
			ExpectedAMI:  "ami-0c7a4976cb6fafd3a",
			ExpectError:  false,
		}),
		Entry("with non-gpu instance and non-eks enabled region", ResolveCase{
			Region:       "ap-northeast-1",
			InstanceType: "t2.medium",
			ExpectedAMI:  "",
			ExpectError:  true,
		}),
		Entry("with gpu (p2) instance and us-west-2", ResolveCase{
			Region:       "us-west-2",
			InstanceType: "p2.xlarge",
			ExpectedAMI:  "ami-0731694d53ef9604b",
			ExpectError:  false,
		}),
		Entry("with gpu (p3) instance and us-east-1", ResolveCase{
			Region:       "us-east-1",
			InstanceType: "p3.2xlarge",
			ExpectedAMI:  "ami-058bfb8c236caae89",
			ExpectError:  false,
		}),
		Entry("with gpu (p2) instance and eu-west-1", ResolveCase{
			Region:       "eu-west-1",
			InstanceType: "p2.xlarge",
			ExpectedAMI:  "ami-0706dc8a5eed2eed9",
			ExpectError:  false,
		}),
		Entry("with gpu (p3) instance and non-eks enabled region", ResolveCase{
			Region:       "ap-northeast-1",
			InstanceType: "p3.2xlarge",
			ExpectedAMI:  "",
			ExpectError:  true,
		}),
	)
})
