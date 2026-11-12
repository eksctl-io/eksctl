package accessentry

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/iam"
)

// {{SessionName}} mappings are shared by Fargate, hybrid and AmazonLinux2027 nodes, so their access entry type is unknown.
var _ = DescribeTable("Migrate {{SessionName}} iamidentitymapping", func(groups []string) {
	identity, err := iam.NewIdentity("arn:aws:iam::111122223333:role/test", authconfigmap.RoleNodeGroupUsernameInstanceID, groups)
	Expect(err).NotTo(HaveOccurred())

	entries, skipAPIMode := doFilterAccessEntries([]iam.Identity{identity}, nil)
	Expect(entries).To(BeEmpty())
	Expect(skipAPIMode).To(BeTrue())
},
	Entry("Fargate pod execution role", []string{"system:bootstrappers", "system:nodes", "system:node-proxier"}),
	Entry("hybrid nodes or AmazonLinux2027 node role", authconfigmap.RoleNodeGroupGroups),
)
