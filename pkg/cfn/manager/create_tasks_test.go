package manager_test

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("CreateTasks", func() {
	var subnetIDs = []string{"123", "456"}
	var clusterConfig *api.ClusterConfig
	Context("AssignIpv6AddressOnCreationTask", func() {
		BeforeEach(func() {
			clusterConfig = api.NewClusterConfig()
			clusterConfig.VPC.Subnets = &api.ClusterSubnets{}
			clusterConfig.VPC.Subnets.Public = map[string]api.AZSubnetSpec{
				"0": {ID: subnetIDs[0]},
				"1": {ID: subnetIDs[1]},
			}
		})

		It("sets AssignIpv6AddressOnCreation to true for all public subnets", func() {
			modifySubnetAttributeCallCount := 0
			p := mockprovider.NewMockProvider()
			mockCall1 := p.MockEC2().On("ModifySubnetAttribute", &ec2.ModifySubnetAttributeInput{
				AssignIpv6AddressOnCreation: &ec2.AttributeBooleanValue{
					Value: aws.Bool(true),
				},
				SubnetId: aws.String(subnetIDs[0]),
			}).Return(&ec2.ModifySubnetAttributeOutput{}, nil)

			mockCall1.RunFn = func(_ mock.Arguments) {
				modifySubnetAttributeCallCount++
			}

			mockCall2 := p.MockEC2().On("ModifySubnetAttribute", &ec2.ModifySubnetAttributeInput{
				AssignIpv6AddressOnCreation: &ec2.AttributeBooleanValue{
					Value: aws.Bool(true),
				},
				SubnetId: aws.String(subnetIDs[1]),
			}).Return(&ec2.ModifySubnetAttributeOutput{}, nil)

			mockCall2.RunFn = func(_ mock.Arguments) {
				modifySubnetAttributeCallCount++
			}

			task := manager.AssignIpv6AddressOnCreationTask{
				EC2API:        p.EC2(),
				ClusterConfig: clusterConfig,
			}
			errorCh := make(chan error)
			err := task.Do(errorCh)
			Expect(err).NotTo(HaveOccurred())
			Expect(modifySubnetAttributeCallCount).To(Equal(2))

			By("closing the error channel")
			Eventually(errorCh).Should(BeClosed())
		})

		When("the API call errors", func() {
			It("errors", func() {
				p := mockprovider.NewMockProvider()
				p.MockEC2().On("ModifySubnetAttribute", &ec2.ModifySubnetAttributeInput{
					AssignIpv6AddressOnCreation: &ec2.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String(subnetIDs[0]),
				}).Return(&ec2.ModifySubnetAttributeOutput{}, fmt.Errorf("foo"))

				task := manager.AssignIpv6AddressOnCreationTask{
					EC2API:        p.EC2(),
					ClusterConfig: clusterConfig,
				}
				errorCh := make(chan error)
				err := task.Do(errorCh)
				Expect(err).To(MatchError("failed to update subnet \"123\": foo"))

				By("closing the error channel")
				Eventually(errorCh).Should(BeClosed())
			})
		})
	})
})
