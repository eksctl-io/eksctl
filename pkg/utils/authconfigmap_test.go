package utils_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/weaveworks/eksctl/pkg/utils"
)

var _ = Describe("Auth ConfigMap", func() {
	Describe("create new ConfigMap", func() {
		It("should create correct configuration for one nodegroup, and add another", func() {

			cm, err := NewAuthConfigMap("arn:aws:iam::122333:role/eksctl-cluster-5a-nodegroup-ng1-p-NodeInstanceRole-NNH3ISP12CX")

			Expect(err).To(Not(HaveOccurred()))
			Expect(cm).To(Not(BeNil()))

			expected := []string{
				`- username: 'system:node:{{EC2PrivateDNSName}}'`,
				`  groups: [ 'system:bootstrappers', 'system:nodes' ]`,
				`  rolearn: 'arn:aws:iam::122333:role/eksctl-cluster-5a-nodegroup-ng1-p-NodeInstanceRole-NNH3ISP12CX'`,
			}

			Expect(cm.ObjectMeta).To(Equal(metav1.ObjectMeta{
				Name:      AuthConfigMapName,
				Namespace: AuthConfigMapNamespace,
			}))

			Expect(cm.Data).To(HaveKey("mapRoles"))
			Expect(cm.Data["mapRoles"]).To(MatchYAML(strings.Join(expected, "\n")))

			err = UpdateAuthConfigMap(cm, "arn:aws:iam::122333:role/eksctl-cluster-5a-nodegroup-ng2-p-NodeInstanceRole-1L35GCVYSTW4E")

			Expect(err).To(Not(HaveOccurred()))
			Expect(cm).To(Not(BeNil()))

			expected = append(expected,
				`- username: 'system:node:{{EC2PrivateDNSName}}'`,
				`  groups: [ 'system:bootstrappers', 'system:nodes' ]`,
				`  rolearn: 'arn:aws:iam::122333:role/eksctl-cluster-5a-nodegroup-ng2-p-NodeInstanceRole-1L35GCVYSTW4E'`,
			)

			Expect(cm.ObjectMeta).To(Equal(metav1.ObjectMeta{
				Name:      AuthConfigMapName,
				Namespace: AuthConfigMapNamespace,
			}))

			Expect(cm.Data).To(HaveKey("mapRoles"))
			Expect(cm.Data["mapRoles"]).To(MatchYAML(strings.Join(expected, "\n")))
		})
	})
})
