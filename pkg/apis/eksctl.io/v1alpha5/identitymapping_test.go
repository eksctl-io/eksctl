package v1alpha5_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = Describe("IAMIdentityMapping", func() {
	It("errors when account is not the only configured option", func() {
		m := v1alpha5.IAMIdentityMapping{
			Account: "0000000000",
			ARN:     "a:fake:arn:hasbeenconfigured",
		}
		err := m.Validate()
		Expect(err).To(HaveOccurred())

	})
	It("errors when namespace is configured without serviceName", func() {
		m := v1alpha5.IAMIdentityMapping{
			Namespace: "emr",
		}
		err := m.Validate()
		Expect(err).To(HaveOccurred())
	})
	It("errors when ARN is configured with serviceName", func() {
		m := v1alpha5.IAMIdentityMapping{
			ServiceName: "emr-containers",
			ARN:         "a:fake:arn:hasbeenconfigured",
		}
		err := m.Validate()
		Expect(err).To(HaveOccurred())
	})
	It("errors when Groups is configured with serviceName", func() {
		m := v1alpha5.IAMIdentityMapping{
			ServiceName: "emr-containers",
			Groups:      []string{"group1"},
		}
		err := m.Validate()
		Expect(err).To(HaveOccurred())
	})
	It("errors when Username is configured with serviceName", func() {
		m := v1alpha5.IAMIdentityMapping{
			ServiceName: "emr-containers",
			Username:    "admin",
		}
		err := m.Validate()
		Expect(err).To(HaveOccurred())
	})
})
