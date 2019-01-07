package manager

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha3"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("StackCollection Tasks", func() {
	var (
		cc *api.ClusterConfig
		sc *StackCollection

		p *mockprovider.MockProvider

		call func(chan error, interface{}) error
	)

	testAZs := []string{"us-west-2b", "us-west-2a", "us-west-2c"}

	newClusterConfig := func(clusterName string) *api.ClusterConfig {
		cfg := api.NewClusterConfig()
		ng := cfg.NewNodeGroup()

		cfg.Metadata.Region = "us-west-2"
		cfg.Metadata.Name = clusterName
		cfg.AvailabilityZones = testAZs
		ng.InstanceType = "t2.medium"
		ng.AMIFamily = "AmazonLinux2"

		*cfg.VPC.CIDR = api.DefaultCIDR()

		return cfg
	}

	Describe("RunTask", func() {
		Context("With a cluster name", func() {
			var (
				clusterName string
				errs        []error

				sucecssfulData []string
			)

			BeforeEach(func() {
				clusterName = "test-cluster"

				p = mockprovider.NewMockProvider()

				cc = newClusterConfig(clusterName)

				sc = NewStackCollection(p, cc)

				sucecssfulData = []string{}

				call = func(errs chan error, data interface{}) error {
					s := data.(string)
					if s == "fail" {
						return errors.New("call failed")
					}

					go func() {
						defer close(errs)

						sucecssfulData = append(sucecssfulData, s)

						errs <- nil
					}()

					return nil
				}
			})

			Context("With an unsuccessful call", func() {
				JustBeforeEach(func() {
					errs = sc.RunSingleTask(Task{call, "fail"})
				})

				It("should error", func() {
					Expect(errs).To(HaveLen(1))
					Expect(errs[0]).To(HaveOccurred())
				})
			})

			Context("With a successful call", func() {
				JustBeforeEach(func() {
					errs = sc.RunSingleTask(Task{call, "ok"})
				})

				It("should not error", func() {
					Expect(errs).To(HaveLen(0))
				})

				It("should have made a side-effect with the successful data", func() {
					Expect(sucecssfulData).To(Equal([]string{"ok"}))
				})
			})
		})
	})
})
