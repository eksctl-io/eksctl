package accessentry_test

import (
	"context"
	"encoding/base32"
	"fmt"
	"strings"

	"github.com/weaveworks/eksctl/pkg/actions/accessentry"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/actions/accessentry/fakes"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
)

var _ = Describe("Access Entry", func() {

	type accessEntryTest struct {
		accessEntry api.AccessEntry
		clusterName string
	}

	DescribeTable("stack creation", func(ae accessEntryTest) {
		var s fakes.FakeStackCreator
		s.CreateStackStub = func(ctx context.Context, stackName string, r builder.ResourceSetReader, tags map[string]string, parameters map[string]string, errorCh chan error) error {
			defer close(errorCh)
			prefix := fmt.Sprintf("eksctl-%s-accessentry-", ae.clusterName)
			idx := strings.Index(stackName, prefix)
			if idx < 0 {
				return fmt.Errorf("expected stack name to have prefix %q", prefix)
			}
			suffix := stackName[idx+len(prefix):]
			_, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(suffix)
			if err != nil {
				return fmt.Errorf("expected stack name to have a base32-encoded suffix: %w", err)
			}
			return nil
		}

		accessEntryCreator := &accessentry.Creator{
			ClusterName:  ae.clusterName,
			StackCreator: &s,
		}
		Expect(accessEntryCreator.Create(context.Background(), []api.AccessEntry{ae.accessEntry})).To(Succeed())
	},
		Entry("access entry 1", accessEntryTest{
			accessEntry: api.AccessEntry{
				PrincipalARN:       api.MustParseARN("arn:aws:iam::111122223333:role/role-1"),
				KubernetesGroups:   []string{"viewers"},
				KubernetesUsername: "user1",
			},
			clusterName: "access-1",
		}),

		Entry("access entry 2", accessEntryTest{
			accessEntry: api.AccessEntry{
				PrincipalARN: api.MustParseARN("arn:aws:iam::111122223333:role/role-2"),
			},
			clusterName: "access-2",
		}),
	)
})
