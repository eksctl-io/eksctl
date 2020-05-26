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

	DescribeTable("When resolving an AMI using the static resolvers",
		func(c ResolveCase) {
			resolver := ami.NewStaticResolver()
			actualAmi, err := resolver.Resolve(c.Region, c.Version, c.InstanceType, c.ImageFamily)
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
			Version:      "1.14",
			InstanceType: "t2.medium",
			ImageFamily:  "AmazonLinux2",
			ExpectedAMI:  "ami-0c13bb9cbfd007e56",
			ExpectError:  false,
		}),
		Entry("with non-gpu instance and us-east-1", ResolveCase{
			Region:       "us-east-1",
			Version:      "1.14",
			InstanceType: "t2.medium",
			ImageFamily:  "AmazonLinux2",
			ExpectedAMI:  "ami-087a82f6b78a07557",
			ExpectError:  false,
		}),
		Entry("with non-gpu instance and eu-west-1", ResolveCase{
			Region:       "eu-west-1",
			Version:      "1.14",
			InstanceType: "t2.medium",
			ImageFamily:  "AmazonLinux2",
			ExpectedAMI:  "ami-0b9d2c11b47bd8264",
			ExpectError:  false,
		}),
		Entry("with non-gpu instance and non-eks enabled region", ResolveCase{
			Region:       "ap-northeast-3",
			Version:      "1.14",
			InstanceType: "t2.medium",
			ImageFamily:  "AmazonLinux2",
			ExpectedAMI:  "",
			ExpectError:  true,
		}),
		Entry("with gpu (p2) instance and us-west-2", ResolveCase{
			Region:       "us-west-2",
			Version:      "1.14",
			InstanceType: "p2.xlarge",
			ImageFamily:  "AmazonLinux2",
			ExpectedAMI:  "ami-0ad9a8dc09680cfc2",
			ExpectError:  false,
		}),
		Entry("with gpu (p3) instance and us-east-1", ResolveCase{
			Region:       "us-east-1",
			Version:      "1.14",
			InstanceType: "p3.2xlarge",
			ImageFamily:  "AmazonLinux2",
			ExpectedAMI:  "ami-0730212bffaa1732a",
			ExpectError:  false,
		}),
		Entry("with gpu (p2) instance and eu-west-1", ResolveCase{
			Region:       "eu-west-1",
			Version:      "1.14",
			InstanceType: "p2.xlarge",
			ImageFamily:  "AmazonLinux2",
			ExpectedAMI:  "ami-03b8e736123fd2b9b",
			ExpectError:  false,
		}),
		Entry("with gpu (g3) instance and eu-west-1", ResolveCase{
			Region:       "eu-west-1",
			Version:      "1.14",
			InstanceType: "g3.4xlarge",
			ImageFamily:  "AmazonLinux2",
			ExpectedAMI:  "ami-03b8e736123fd2b9b",
			ExpectError:  false,
		}),
		Entry("with gpu (g4dn) instance and us-east-1", ResolveCase{
			Region:       "us-east-1",
			Version:      "1.14",
			InstanceType: "g4dn.xlarge",
			ImageFamily:  "AmazonLinux2",
			ExpectedAMI:  "ami-0730212bffaa1732a",
			ExpectError:  false,
		}),
		Entry("with gpu (p3) instance and non-existent region", ResolveCase{
			Region:       "eu-east-1",
			Version:      "1.14",
			InstanceType: "p3.2xlarge",
			ImageFamily:  "AmazonLinux2",
			ExpectedAMI:  "",
			ExpectError:  true,
		}),
		Entry("with non-gpu instance, us-west-2 and Ubuntu image", ResolveCase{
			Region:       "us-west-2",
			Version:      "1.14",
			InstanceType: "t2.medium",
			ImageFamily:  "Ubuntu1804",
			ExpectedAMI:  "ami-0722493f75c013523",
			ExpectError:  false,
		}),
		Entry("with non-gpu instance, us-east-1 and Ubuntu image", ResolveCase{
			Region:       "us-east-1",
			Version:      "1.14",
			InstanceType: "t2.medium",
			ImageFamily:  "Ubuntu1804",
			ExpectedAMI:  "ami-093abf0b4d2d071df",
			ExpectError:  false,
		}),
		Entry("with non-gpu instance, eu-west-1 and Ubuntu image", ResolveCase{
			Region:       "eu-west-1",
			Version:      "1.14",
			InstanceType: "t2.medium",
			ImageFamily:  "Ubuntu1804",
			ExpectedAMI:  "ami-0dd0a16a2bd0784b8",
			ExpectError:  false,
		}),
		Entry("with non-gpu instance, non-eks enabled region and Ubuntu image", ResolveCase{
			Region:       "ap-northeast-3",
			Version:      "1.14",
			InstanceType: "t2.medium",
			ImageFamily:  "Ubuntu1804",
			ExpectedAMI:  "",
			ExpectError:  true,
		}),
		Entry("with gpu instance, any region and Ubuntu image", ResolveCase{
			Region:       "us-east-1",
			Version:      "1.14",
			InstanceType: "p2.xlarge",
			ImageFamily:  "Ubuntu1804",
			ExpectedAMI:  "",
			ExpectError:  true,
		}),
	)
})
