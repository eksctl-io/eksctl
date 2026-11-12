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
	It("uses the AL2023 nodeadm user data", func() {
		cfg, _ := makeDefaultClusterSettings()
		ng := api.NewNodeGroup()
		makeDefaultNPSettings(ng)
		ng.AMIFamily = api.NodeImageFamilyAmazonLinux2027

		bootstrapper, err := nodebootstrap.NewBootstrapper(cfg, ng)
		Expect(err).NotTo(HaveOccurred())
		Expect(bootstrapper).To(BeAssignableToTypeOf(&nodebootstrap.AL2023{}))

		al2023BS := bootstrapper.(*nodebootstrap.AL2023)
		al2023BS.UserDataMimeBoundary = "//"
		userData, err := al2023BS.UserData()
		Expect(err).NotTo(HaveOccurred())
		decoded, err := base64.StdEncoding.DecodeString(userData)
		Expect(err).NotTo(HaveOccurred())
		actual := strings.ReplaceAll(string(decoded), "\r\n", "\n")
		Expect(actual).To(Equal(wrapMIMEParts(nodeConfig)))
	})
})
