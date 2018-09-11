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
			actualAmi, err := ami.ResolveAMI(c.Region, c.InstanceType)
			Expect(actualAmi).Should(Equal(c.ExpectedAMI))
			if c.ExpectError {
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(MatchError(ami.NewErrFailedAMIResolution(c.Region, c.InstanceType)))
				errorType := reflect.TypeOf(err).Elem().Name()
				Expect(errorType).To(Equal("ErrFailedAMIResolution"))
			} else {
				Expect(err).ShouldNot(HaveOccurred())
			}
		},
		Entry("with non-gpu instance and us-west-2", ResolveCase{
			Region:       "us-west-2",
			InstanceType: "t2.medium",
			ExpectedAMI:  "ami-08cab282f9979fc7a",
			ExpectError:  false,
		}),
		Entry("with non-gpu instance and us-east-1", ResolveCase{
			Region:       "us-east-1",
			InstanceType: "t2.medium",
			ExpectedAMI:  "ami-0b2ae3c6bda8b5c06",
			ExpectError:  false,
		}),
		Entry("with non-gpu instance and eu-west-1", ResolveCase{
			Region:       "eu-west-1",
			InstanceType: "t2.medium",
			ExpectedAMI:  "ami-066110c1a7466949e",
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
			ExpectedAMI:  "ami-0d20f2404b9a1c4d1",
			ExpectError:  false,
		}),
		Entry("with gpu (p3) instance and us-east-1", ResolveCase{
			Region:       "us-east-1",
			InstanceType: "p3.2xlarge",
			ExpectedAMI:  "ami-09fe6fc9106bda972",
			ExpectError:  false,
		}),
		Entry("with gpu (p2) instance and eu-west-1", ResolveCase{
			Region:       "eu-west-1",
			InstanceType: "p2.xlarge",
			ExpectedAMI:  "ami-09e0c6b3d3cf906f1",
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
