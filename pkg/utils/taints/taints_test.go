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
	taint       corev1.Taint
	expectedErr string
}

var _ = DescribeTable("Parse", func(t taintsEntry) {
	err := taints.Validate(t.taint)
	if t.expectedErr != "" {
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ContainSubstring(t.expectedErr)))
		return
	}

	Expect(err).NotTo(HaveOccurred())
},
	Entry("valid", taintsEntry{
		taint: corev1.Taint{
			Key:    "key1",
			Value:  "value1",
			Effect: corev1.TaintEffectNoSchedule,
		},
	}),

	Entry("missing value", taintsEntry{
		taint: corev1.Taint{
			Key:    "key2",
			Effect: corev1.TaintEffectNoExecute,
		},
	}),

	Entry("missing value and effect", taintsEntry{
		taint: corev1.Taint{
			Key: "key3",
		},
		expectedErr: "invalid taint effect",
	}),

	Entry("invalid key", taintsEntry{
		taint: corev1.Taint{
			Key:    "key1=",
			Value:  "value1",
			Effect: corev1.TaintEffectNoSchedule,
		},
		expectedErr: "invalid taint key",
	}),

	Entry("invalid value", taintsEntry{
		taint: corev1.Taint{
			Key:    "key1",
			Value:  "v&lue",
			Effect: corev1.TaintEffectNoSchedule,
		},
		expectedErr: "invalid taint value",
	}),

	Entry("unsupported effect", taintsEntry{
		taint: corev1.Taint{
			Key:    "key1",
			Value:  "value1",
			Effect: "NoEffect",
		},
		expectedErr: "invalid taint effect",
	}),
)
