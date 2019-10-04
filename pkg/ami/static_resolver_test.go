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
			Version:      "1.12",
			InstanceType: "t2.medium",
			ImageFamily:  "AmazonLinux2",
			ExpectedAMI:  "ami-0b520e822d42998c1",
			ExpectError:  false,
		}),
		Entry("with non-gpu instance and us-east-1", ResolveCase{
			Region:       "us-east-1",
			Version:      "1.12",
			InstanceType: "t2.medium",
			ImageFamily:  "AmazonLinux2",
			ExpectedAMI:  "ami-01e370f796735b244",
			ExpectError:  false,
		}),
		Entry("with non-gpu instance and eu-west-1", ResolveCase{
			Region:       "eu-west-1",
			Version:      "1.12",
			InstanceType: "t2.medium",
			ImageFamily:  "AmazonLinux2",
			ExpectedAMI:  "ami-018b4a3f81f517183",
			ExpectError:  false,
		}),
		Entry("with non-gpu instance and non-eks enabled region", ResolveCase{
			Region:       "ap-northeast-3",
			Version:      "1.12",
			InstanceType: "t2.medium",
			ImageFamily:  "AmazonLinux2",
			ExpectedAMI:  "",
			ExpectError:  true,
		}),
		Entry("with gpu (p2) instance and us-west-2", ResolveCase{
			Region:       "us-west-2",
			Version:      "1.12",
			InstanceType: "p2.xlarge",
			ImageFamily:  "AmazonLinux2",
			ExpectedAMI:  "ami-0c9156d7fcd3c2948",
			ExpectError:  false,
		}),
		Entry("with gpu (p3) instance and us-east-1", ResolveCase{
			Region:       "us-east-1",
			Version:      "1.12",
			InstanceType: "p3.2xlarge",
			ImageFamily:  "AmazonLinux2",
			ExpectedAMI:  "ami-0eb0119f55d589a03",
			ExpectError:  false,
		}),
		Entry("with gpu (p2) instance and eu-west-1", ResolveCase{
			Region:       "eu-west-1",
			Version:      "1.12",
			InstanceType: "p2.xlarge",
			ImageFamily:  "AmazonLinux2",
			ExpectedAMI:  "ami-04668c090ff8c1f50",
			ExpectError:  false,
		}),
		Entry("with gpu (g3) instance and eu-west-1", ResolveCase{
			Region:       "eu-west-1",
			Version:      "1.12",
			InstanceType: "g3.4xlarge",
			ImageFamily:  "AmazonLinux2",
			ExpectedAMI:  "ami-04668c090ff8c1f50",
			ExpectError:  false,
		}),
		Entry("with gpu (g4dn) instance and us-east-1", ResolveCase{
			Region:       "us-east-1",
			Version:      "1.14",
			InstanceType: "g4dn.xlarge",
			ImageFamily:  "AmazonLinux2",
			ExpectedAMI:  "ami-06191d186b526c958",
			ExpectError:  false,
		}),
		Entry("with gpu (p3) instance and non-existent region", ResolveCase{
			Region:       "eu-east-1",
			Version:      "1.12",
			InstanceType: "p3.2xlarge",
			ImageFamily:  "AmazonLinux2",
			ExpectedAMI:  "",
			ExpectError:  true,
		}),
		Entry("with non-gpu instance, us-west-2 and Ubuntu image", ResolveCase{
			Region:       "us-west-2",
			Version:      "1.12",
			InstanceType: "t2.medium",
			ImageFamily:  "Ubuntu1804",
			ExpectedAMI:  "ami-0551b6dc0b4079e1d",
			ExpectError:  false,
		}),
		Entry("with non-gpu instance, us-east-1 and Ubuntu image", ResolveCase{
			Region:       "us-east-1",
			Version:      "1.12",
			InstanceType: "t2.medium",
			ImageFamily:  "Ubuntu1804",
			ExpectedAMI:  "ami-058c6e6f18b1cca2c",
			ExpectError:  false,
		}),
		Entry("with non-gpu instance, eu-west-1 and Ubuntu image", ResolveCase{
			Region:       "eu-west-1",
			Version:      "1.12",
			InstanceType: "t2.medium",
			ImageFamily:  "Ubuntu1804",
			ExpectedAMI:  "ami-0ac97cadeed862a54",
			ExpectError:  false,
		}),
		Entry("with non-gpu instance, non-eks enabled region and Ubuntu image", ResolveCase{
			Region:       "ap-northeast-3",
			Version:      "1.12",
			InstanceType: "t2.medium",
			ImageFamily:  "Ubuntu1804",
			ExpectedAMI:  "",
			ExpectError:  true,
		}),
		Entry("with gpu instance, any region and Ubuntu image", ResolveCase{
			Region:       "us-east-1",
			Version:      "1.12",
			InstanceType: "p2.xlarge",
			ImageFamily:  "Ubuntu1804",
			ExpectedAMI:  "",
			ExpectError:  true,
		}),
	)
})
