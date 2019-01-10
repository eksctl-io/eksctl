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
		Version      string
		InstanceType string
		ImageFamily  string
		ExpectedAMI  string
		ExpectError  bool
	}

	DescribeTable("When resolving an AMI using the default resolvers",
		func(c ResolveCase) {
			actualAmi, err := ami.Resolve(c.Region, c.Version, c.InstanceType, c.ImageFamily)
			Expect(actualAmi).Should(Equal(c.ExpectedAMI))
			if c.ExpectError {
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(MatchError(ami.NewErrFailedResolution(c.Region, c.Version, c.InstanceType, c.ImageFamily)))
				errorType := reflect.TypeOf(err).Elem().Name()
				Expect(errorType).To(Equal("ErrFailedResolution"))
			} else {
				Expect(err).ShouldNot(HaveOccurred())
			}
		},
		Entry("with non-gpu instance and us-west-2", ResolveCase{
			Region:       "us-west-2",
			Version:      "1.10",
			InstanceType: "t2.medium",
			ImageFamily:  "AmazonLinux2",
			ExpectedAMI:  "ami-09e1df3bad220af0b",
			ExpectError:  false,
		}),
		Entry("with non-gpu instance and us-east-1", ResolveCase{
			Region:       "us-east-1",
			Version:      "1.10",
			InstanceType: "t2.medium",
			ImageFamily:  "AmazonLinux2",
			ExpectedAMI:  "ami-04358410d28eaab63",
			ExpectError:  false,
		}),
		Entry("with non-gpu instance and eu-west-1", ResolveCase{
			Region:       "eu-west-1",
			Version:      "1.10",
			InstanceType: "t2.medium",
			ImageFamily:  "AmazonLinux2",
			ExpectedAMI:  "ami-0de10c614955da932",
			ExpectError:  false,
		}),
		Entry("with non-gpu instance and non-eks enabled region", ResolveCase{
			Region:       "ap-northeast-2",
			Version:      "1.10",
			InstanceType: "t2.medium",
			ImageFamily:  "AmazonLinux2",
			ExpectedAMI:  "",
			ExpectError:  true,
		}),
		Entry("with gpu (p2) instance and us-west-2", ResolveCase{
			Region:       "us-west-2",
			Version:      "1.10",
			InstanceType: "p2.xlarge",
			ImageFamily:  "AmazonLinux2",
			ExpectedAMI:  "ami-0ebf0561e61a2be02",
			ExpectError:  false,
		}),
		Entry("with gpu (p3) instance and us-east-1", ResolveCase{
			Region:       "us-east-1",
			Version:      "1.10",
			InstanceType: "p3.2xlarge",
			ImageFamily:  "AmazonLinux2",
			ExpectedAMI:  "ami-0131c0ca222183def",
			ExpectError:  false,
		}),
		Entry("with gpu (p2) instance and eu-west-1", ResolveCase{
			Region:       "eu-west-1",
			Version:      "1.10",
			InstanceType: "p2.xlarge",
			ImageFamily:  "AmazonLinux2",
			ExpectedAMI:  "ami-0dafd3a1dc43781f7",
			ExpectError:  false,
		}),
		Entry("with gpu (p3) instance and non-eks enabled region", ResolveCase{
			Region:       "ap-northeast-2",
			Version:      "1.10",
			InstanceType: "p3.2xlarge",
			ImageFamily:  "AmazonLinux2",
			ExpectedAMI:  "",
			ExpectError:  true,
		}),
		Entry("with non-gpu instance, us-west-2 and Ubuntu image", ResolveCase{
			Region:       "us-west-2",
			Version:      "1.10",
			InstanceType: "t2.medium",
			ImageFamily:  "Ubuntu1804",
			ExpectedAMI:  "ami-6322011b",
			ExpectError:  false,
		}),
		Entry("with non-gpu instance, us-east-1 and Ubuntu image", ResolveCase{
			Region:       "us-east-1",
			Version:      "1.10",
			InstanceType: "t2.medium",
			ImageFamily:  "Ubuntu1804",
			ExpectedAMI:  "ami-06fd8200ac0eb656d",
			ExpectError:  false,
		}),
		Entry("with non-gpu instance, eu-west-1 and Ubuntu image", ResolveCase{
			Region:       "eu-west-1",
			Version:      "1.10",
			InstanceType: "t2.medium",
			ImageFamily:  "Ubuntu1804",
			ExpectedAMI:  "ami-07036622490f7e97b",
			ExpectError:  false,
		}),
		Entry("with non-gpu instance, non-eks enabled region and Ubuntu image", ResolveCase{
			Region:       "ap-northeast-2",
			Version:      "1.10",
			InstanceType: "t2.medium",
			ImageFamily:  "Ubuntu1804",
			ExpectedAMI:  "",
			ExpectError:  true,
		}),
		Entry("with gpu instance, any region and Ubuntu image", ResolveCase{
			Region:       "us-east-1",
			Version:      "1.10",
			InstanceType: "p2.xlarge",
			ImageFamily:  "Ubuntu1804",
			ExpectedAMI:  "",
			ExpectError:  true,
		}),
	)
})
