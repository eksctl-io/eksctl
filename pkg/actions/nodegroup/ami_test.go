package nodegroup_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
)

var _ = Describe("AMI Release Version", func() {
	type versionCase struct {
		v1     string
		v2     string
		cmp    int
		errMsg string
	}

	compare := func(a, b string) (int, error) {
		v1, err := nodegroup.ParseReleaseVersion(a)
		if err != nil {
			return 0, err
		}
		v2, err := nodegroup.ParseReleaseVersion(b)
		if err != nil {
			return 0, err
		}
		return v1.Compare(v2), nil
	}

	DescribeTable("AMI release version comparison", func(vc versionCase) {
		cmp, err := compare(vc.v1, vc.v2)
		if vc.errMsg != "" {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(vc.errMsg))
			return
		}
		Expect(err).NotTo(HaveOccurred())
		Expect(cmp).To(Equal(vc.cmp))
	},
		Entry("Equal", versionCase{
			v1:  "1.18.8-20201007",
			v2:  "1.18.8-20201007",
			cmp: 0,
		}),
		Entry("Less", versionCase{
			v1:  "1.18.8-20201007",
			v2:  "1.18.9-20201112",
			cmp: -1,
		}),
		Entry("Greater", versionCase{
			v1:  "1.18.25-20201007",
			v2:  "1.18.20-20201007",
			cmp: 1,
		}),
		Entry("Older major version", versionCase{
			v1:  "1.17.9-20200101",
			v2:  "1.18.0-20200101",
			cmp: -1,
		}),
		Entry("Newer minor version", versionCase{
			v1:  "1.18.9-20201112",
			v2:  "1.18.8-20201007",
			cmp: 1,
		}),
		Entry("Malformed version", versionCase{
			v1:     "1.18.9-20200101",
			v2:     "1.18.9",
			errMsg: "unexpected format",
		}),
		Entry("Both versions invalid", versionCase{
			v1:     "a-b",
			v2:     "1.17.d",
			errMsg: "invalid SemVer version",
		}),
	)

})
