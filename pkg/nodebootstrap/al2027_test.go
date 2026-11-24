package nodebootstrap_test

import (
	"encoding/base64"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
)

var _ = Describe("Unmanaged AL2027", func() {
	It("uses the same nodeadm user data as AL2023", func() {
		cfg, _ := makeDefaultClusterSettings()
		ng := api.NewNodeGroup()
		makeDefaultNPSettings(ng)
		ng.AMIFamily = api.NodeImageFamilyAmazonLinux2027

		bootstrapper, err := nodebootstrap.NewBootstrapper(cfg, ng)
		Expect(err).NotTo(HaveOccurred())
		Expect(bootstrapper).To(BeAssignableToTypeOf(&nodebootstrap.Nodeadm{}))

		nodeadmBS := bootstrapper.(*nodebootstrap.Nodeadm)
		nodeadmBS.UserDataMimeBoundary = "//"
		userData, err := nodeadmBS.UserData()
		Expect(err).NotTo(HaveOccurred())
		decoded, err := base64.StdEncoding.DecodeString(userData)
		Expect(err).NotTo(HaveOccurred())
		actual := strings.ReplaceAll(string(decoded), "\r\n", "\n")
		Expect(actual).To(Equal(wrapMIMEParts(nodeConfig)))
	})
})

var _ = Describe("Managed AL2027", func() {
	It("uses the same nodeadm user data as AL2023 with a custom AMI", func() {
		cfg, _ := makeDefaultClusterSettings()
		mng := api.NewManagedNodeGroup()
		makeDefaultNPSettings(mng)
		mng.Taints = append(mng.Taints, api.NodeGroupTaint{
			Key:    "special",
			Value:  "true",
			Effect: "NoSchedule",
		})
		mng.AMIFamily = api.NodeImageFamilyAmazonLinux2027
		mng.AMI = "ami-xxxx"

		bootstrapper, err := nodebootstrap.NewManagedBootstrapper(cfg, mng)
		Expect(err).NotTo(HaveOccurred())
		Expect(bootstrapper).To(BeAssignableToTypeOf(&nodebootstrap.Nodeadm{}))

		nodeadmBS := bootstrapper.(*nodebootstrap.Nodeadm)
		nodeadmBS.UserDataMimeBoundary = "//"
		userData, err := nodeadmBS.UserData()
		Expect(err).NotTo(HaveOccurred())
		decoded, err := base64.StdEncoding.DecodeString(userData)
		Expect(err).NotTo(HaveOccurred())
		actual := strings.ReplaceAll(string(decoded), "\r\n", "\n")
		Expect(actual).To(Equal(wrapMIMEParts(managedNodeConfig)))
	})
})
