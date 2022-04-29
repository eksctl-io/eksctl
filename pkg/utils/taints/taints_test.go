package taints_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/utils/taints"

	corev1 "k8s.io/api/core/v1"
)

func TestUtilsTaints(t *testing.T) {
	testutils.RegisterAndRun(t)
}

type validateTaintsEntry struct {
	taint       corev1.Taint
	expectedErr string
}

type parseTaintsEntry struct {
	key         string
	valueEffect string

	expectedTaint corev1.Taint
}

var _ = Describe("Taints", func() {
	DescribeTable("Validate", func(t validateTaintsEntry) {
		err := taints.Validate(t.taint)
		if t.expectedErr != "" {
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring(t.expectedErr)))
			return
		}

		Expect(err).NotTo(HaveOccurred())
	},
		Entry("valid", validateTaintsEntry{
			taint: corev1.Taint{
				Key:    "key1",
				Value:  "value1",
				Effect: corev1.TaintEffectNoSchedule,
			},
		}),

		Entry("missing value", validateTaintsEntry{
			taint: corev1.Taint{
				Key:    "key2",
				Effect: corev1.TaintEffectNoExecute,
			},
		}),

		Entry("missing value and effect", validateTaintsEntry{
			taint: corev1.Taint{
				Key: "key3",
			},
			expectedErr: "invalid taint effect",
		}),

		Entry("invalid key", validateTaintsEntry{
			taint: corev1.Taint{
				Key:    "key1=",
				Value:  "value1",
				Effect: corev1.TaintEffectNoSchedule,
			},
			expectedErr: "invalid taint key",
		}),

		Entry("invalid value", validateTaintsEntry{
			taint: corev1.Taint{
				Key:    "key1",
				Value:  "v&lue",
				Effect: corev1.TaintEffectNoSchedule,
			},
			expectedErr: "invalid taint value",
		}),

		Entry("unsupported effect", validateTaintsEntry{
			taint: corev1.Taint{
				Key:    "key1",
				Value:  "value1",
				Effect: "NoEffect",
			},
			expectedErr: "invalid taint effect",
		}),
	)

	DescribeTable("Parse", func(t parseTaintsEntry) {
		parsedTaints := taints.Parse(map[string]string{
			t.key: t.valueEffect,
		})
		Expect(parsedTaints[0]).To(Equal(t.expectedTaint))
	},
		Entry("valid", parseTaintsEntry{
			key:         "key1",
			valueEffect: "value1:NoSchedule",
			expectedTaint: corev1.Taint{
				Key:    "key1",
				Value:  "value1",
				Effect: corev1.TaintEffectNoSchedule,
			},
		}),

		Entry("missing value", parseTaintsEntry{
			key:         "key2",
			valueEffect: ":NoExecute",
			expectedTaint: corev1.Taint{
				Key:    "key2",
				Effect: corev1.TaintEffectNoExecute,
			},
		}),
	)
})
