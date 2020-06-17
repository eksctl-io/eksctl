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
		ExpectError  bool
	}

	DescribeTable("When resolving an AMI using the static resolvers",
		func(c ResolveCase) {
			resolver := ami.NewStaticResolver()
			actualAmi, err := resolver.Resolve(c.Region, c.Version, c.InstanceType, c.ImageFamily)
			if c.ExpectError {
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(MatchError(ami.NewErrFailedResolution(c.Region, c.Version, c.InstanceType, c.ImageFamily)))
				errorType := reflect.TypeOf(err).Elem().Name()
				Expect(errorType).To(Equal("ErrFailedResolution"))
			} else {
				Expect(actualAmi).To(HavePrefix("ami"))
				Expect(err).ShouldNot(HaveOccurred())
			}
		},
		Entry("with non-gpu instance and us-west-2", ResolveCase{
			Region:       "us-west-2",
			Version:      "1.14",
			InstanceType: "t2.medium",
			ImageFamily:  "AmazonLinux2",
			ExpectError:  false,
		}),
		Entry("with non-gpu instance and us-east-1", ResolveCase{
			Region:       "us-east-1",
			Version:      "1.14",
			InstanceType: "t2.medium",
			ImageFamily:  "AmazonLinux2",
			ExpectError:  false,
		}),
		Entry("with non-gpu instance and eu-west-1", ResolveCase{
			Region:       "eu-west-1",
			Version:      "1.14",
			InstanceType: "t2.medium",
			ImageFamily:  "AmazonLinux2",
			ExpectError:  false,
		}),
		Entry("with non-gpu instance and non-eks enabled region", ResolveCase{
			Region:       "ap-northeast-3",
			Version:      "1.14",
			InstanceType: "t2.medium",
			ImageFamily:  "AmazonLinux2",
			ExpectError:  true,
		}),
		Entry("with gpu (p2) instance and us-west-2", ResolveCase{
			Region:       "us-west-2",
			Version:      "1.14",
			InstanceType: "p2.xlarge",
			ImageFamily:  "AmazonLinux2",
			ExpectError:  false,
		}),
		Entry("with gpu (p3) instance and us-east-1", ResolveCase{
			Region:       "us-east-1",
			Version:      "1.14",
			InstanceType: "p3.2xlarge",
			ImageFamily:  "AmazonLinux2",
			ExpectError:  false,
		}),
		Entry("with gpu (p2) instance and eu-west-1", ResolveCase{
			Region:       "eu-west-1",
			Version:      "1.14",
			InstanceType: "p2.xlarge",
			ImageFamily:  "AmazonLinux2",
			ExpectError:  false,
		}),
		Entry("with gpu (g3) instance and eu-west-1", ResolveCase{
			Region:       "eu-west-1",
			Version:      "1.14",
			InstanceType: "g3.4xlarge",
			ImageFamily:  "AmazonLinux2",
			ExpectError:  false,
		}),
		Entry("with gpu (g4dn) instance and us-east-1", ResolveCase{
			Region:       "us-east-1",
			Version:      "1.14",
			InstanceType: "g4dn.xlarge",
			ImageFamily:  "AmazonLinux2",
			ExpectError:  false,
		}),
		Entry("with gpu (p3) instance and non-existent region", ResolveCase{
			Region:       "eu-east-1",
			Version:      "1.14",
			InstanceType: "p3.2xlarge",
			ImageFamily:  "AmazonLinux2",
			ExpectError:  true,
		}),
		Entry("with non-gpu instance, us-west-2 and Ubuntu image", ResolveCase{
			Region:       "us-west-2",
			Version:      "1.14",
			InstanceType: "t2.medium",
			ImageFamily:  "Ubuntu1804",
			ExpectError:  false,
		}),
		Entry("with non-gpu instance, us-east-1 and Ubuntu image", ResolveCase{
			Region:       "us-east-1",
			Version:      "1.14",
			InstanceType: "t2.medium",
			ImageFamily:  "Ubuntu1804",
			ExpectError:  false,
		}),
		Entry("with non-gpu instance, eu-west-1 and Ubuntu image", ResolveCase{
			Region:       "eu-west-1",
			Version:      "1.14",
			InstanceType: "t2.medium",
			ImageFamily:  "Ubuntu1804",
			ExpectError:  false,
		}),
		Entry("with non-gpu instance, non-eks enabled region and Ubuntu image", ResolveCase{
			Region:       "ap-northeast-3",
			Version:      "1.14",
			InstanceType: "t2.medium",
			ImageFamily:  "Ubuntu1804",
			ExpectError:  true,
		}),
		Entry("with gpu instance, any region and Ubuntu image", ResolveCase{
			Region:       "us-east-1",
			Version:      "1.14",
			InstanceType: "p2.xlarge",
			ImageFamily:  "Ubuntu1804",
			ExpectError:  true,
		}),
	)
})
