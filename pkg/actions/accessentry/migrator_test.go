package accessentry_test

import (
	"bytes"
	"context"
	"encoding/base32"
	"fmt"
	"os"
	"strings"

	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/kris-nova/logger"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/weaveworks/eksctl/pkg/actions/accessentry"
	"github.com/weaveworks/eksctl/pkg/actions/accessentry/fakes"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	"k8s.io/client-go/kubernetes/fake"
)

type migrateToAccessEntryEntry struct {
	clusterName string
	mockEKS     func(provider *mockprovider.MockProvider)
	// mockIAM                    func(provider *mockprovider.MockProvider)
	mockK8s                    func(clientSet *fake.Clientset)
	validateCustomLoggerOutput func(output string)
	options                    accessentry.MigrationOptions
	expectedErr                string
}

var _ = Describe("Migrate Access Entry", func() {

	var (
		migrator      *accessentry.Migrator
		mockProvider  *mockprovider.MockProvider
		fakeClientset *fake.Clientset
		clusterName   = "test-cluster"
		tgAuthMode    = ekstypes.AuthenticationModeApi
		curAuthMode   = ekstypes.AuthenticationModeConfigMap
	)

	DescribeTable("Migrate", func(ae migrateToAccessEntryEntry) {
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

		mockProvider = mockprovider.NewMockProvider()
		if ae.mockEKS != nil {
			ae.mockEKS(mockProvider)
		}

		fakeClientset = fake.NewSimpleClientset()
		if ae.mockK8s != nil {
			ae.mockK8s(fakeClientset)
		}

		output := &bytes.Buffer{}
		if ae.validateCustomLoggerOutput != nil {
			defer func() {
				logger.Writer = os.Stdout
			}()
			logger.Writer = output
		}

		migrator = accessentry.NewMigrator(
			ae.clusterName,
			mockProvider.MockEKS(),
			mockProvider.MockIAM(),
			fakeClientset,
			*accessEntryCreator,
			curAuthMode,
			tgAuthMode,
		)

		err := migrator.MigrateToAccessEntry(context.Background(), ae.options)

		if ae.expectedErr != "" {
			Expect(err).To(MatchError(ContainSubstring(ae.expectedErr)))
			return
		}

		Expect(err).ToNot(HaveOccurred())

		if ae.validateCustomLoggerOutput != nil {
			ae.validateCustomLoggerOutput(output.String())
		}
	}, Entry("[API Error] Authentication mode update fails", migrateToAccessEntryEntry{
		clusterName: clusterName,
		mockEKS: func(provider *mockprovider.MockProvider) {
			mockProvider.MockEKS().
				On("UpdateClusterConfig", mock.Anything, mock.Anything).
				Run(func(args mock.Arguments) {
					Expect(args).To(HaveLen(2))
					Expect(args[1]).To(BeAssignableToTypeOf(&awseks.UpdateClusterConfigInput{}))
				}).
				Return(nil, fmt.Errorf("failed to update cluster config"))
		},
		expectedErr: "failed to update cluster config",
	}),
	)
})
