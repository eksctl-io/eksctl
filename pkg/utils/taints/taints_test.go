package taints_test

import (
	"testing"

	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/utils/taints"

	corev1 "k8s.io/api/core/v1"
)

func TestUtilsTaints(t *testing.T) {
	testutils.RegisterAndRun(t)
}

type taintsEntry struct {
	key         string
	valueEffect string

	expectedTaint corev1.Taint
	expectedErr   string
}

var _ = DescribeTable("Parse", func(t taintsEntry) {
	parsedTaints, err := taints.Parse(map[string]string{
		t.key: t.valueEffect,
	})
	if t.expectedErr != "" {
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ContainSubstring(t.expectedErr)))
		return
	}

	Expect(err).NotTo(HaveOccurred())
	Expect(parsedTaints[0]).To(Equal(t.expectedTaint))
},
	Entry("valid", taintsEntry{
		key:         "key1",
		valueEffect: "value1:NoSchedule",
		expectedTaint: corev1.Taint{
			Key:    "key1",
			Value:  "value1",
			Effect: corev1.TaintEffectNoSchedule,
		},
	}),

	Entry("missing value", taintsEntry{
		key:         "key2",
		valueEffect: ":NoExecute",
		expectedTaint: corev1.Taint{
			Key:    "key2",
			Effect: corev1.TaintEffectNoExecute,
		},
	}),

	Entry("missing value and effect", taintsEntry{
		key:         "key3",
		expectedErr: "invalid taint effect",
	}),

	Entry("invalid key", taintsEntry{
		key:         "key1=",
		valueEffect: "value1:NoSchedule",
		expectedErr: "invalid taint key",
	}),

	Entry("invalid value", taintsEntry{
		key:         "key1",
		valueEffect: "v&lue:NoSchedule",
		expectedErr: "invalid taint value",
	}),

	Entry("unsupported effect", taintsEntry{
		key:         "key1",
		valueEffect: "value1:NoEffect",
		expectedErr: "invalid taint effect",
	}),
)
