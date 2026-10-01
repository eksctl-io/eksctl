package testutils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"

	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/eks"
)

func Reader(clusterConfig *api.ClusterConfig) io.Reader {
	data, err := json.Marshal(clusterConfig)
	Expect(err).NotTo(HaveOccurred())
	return bytes.NewReader(data)
}

func ReaderFromFile(clusterName, region, filename string) io.Reader {
	clusterConfig := ParseClusterConfig(clusterName, region, filename)
	return Reader(clusterConfig)
}

func ParseClusterConfig(clusterName, region, filename string) *api.ClusterConfig {
	data, err := os.ReadFile(filename)
	Expect(err).NotTo(HaveOccurred())
	clusterConfig, err := eks.ParseConfig(data)
	Expect(err).NotTo(HaveOccurred())
	clusterConfig.Metadata.Name = clusterName
	clusterConfig.Metadata.Region = region
	return clusterConfig
}

func GetCurrentAndNextVersionsForUpgrade(cvm eks.ClusterVersionsManagerInterface, testVersion string) (currentVersion, nextVersion string) {
	supportedVersions := cvm.SupportedVersions()
	if len(supportedVersions) < 2 {
		Fail("Upgrade test requires at least two supported EKS versions")
	}

	// if latest version is used, fetch previous version to upgrade from
	if testVersion == cvm.LatestVersion() {
		previousVersionIndex := slices.Index(supportedVersions, testVersion) - 1
		currentVersion = supportedVersions[previousVersionIndex]
		nextVersion = testVersion
		return
	}

	// otherwise fetch next version to upgrade to
	nextVersionIndex := slices.Index(supportedVersions, testVersion) + 1
	currentVersion = testVersion
	nextVersion = supportedVersions[nextVersionIndex]
	return
}

const (
	clusterSettleTimeout      = "15m"
	clusterSettlePollInterval = "15s"
)

// WaitForClusterToSettle blocks until clusterName has no EKS update in flight.
//
// EKS reports a cluster as ACTIVE while asynchronous post-create work -- default addon
// installs, access-config changes -- is still running. Any mutating call issued during that
// window is rejected with HTTP 409:
//
//	ResourceInUseException: Cannot VersionUpdate because cluster <name> currently has an
//	update in progress
//
// So a DescribeCluster status of ACTIVE is not a sufficient readiness signal, and neither are
// eksctl's own CanUpdate/CanOperate guards, which only consult that status. ListUpdates is the
// signal that actually reflects the condition named in the error.
//
// Call this after creating a cluster and before the first mutating operation a suite performs
// on it. It is read-only and a no-op on a settled cluster, so it is safe to call speculatively.
func WaitForClusterToSettle(ctx context.Context, eksAPI awsapi.EKS, clusterName string) {
	By(fmt.Sprintf("waiting for cluster %q to have no in-flight updates", clusterName))
	Eventually(func() ([]string, error) {
		return InProgressUpdates(ctx, eksAPI, clusterName)
	}, clusterSettleTimeout, clusterSettlePollInterval).Should(BeEmpty(),
		"expected cluster %q to have no in-flight EKS updates before mutating it", clusterName)
}

// InProgressUpdates returns the IDs of any EKS updates on clusterName that are still in
// progress. An empty result means the cluster is safe to mutate.
func InProgressUpdates(ctx context.Context, eksAPI awsapi.EKS, clusterName string) ([]string, error) {
	var inProgress []string
	var nextToken *string
	for {
		updates, err := eksAPI.ListUpdates(ctx, &awseks.ListUpdatesInput{
			Name:      &clusterName,
			NextToken: nextToken,
		})
		if err != nil {
			return nil, fmt.Errorf("listing updates for cluster %q: %w", clusterName, err)
		}
		for _, id := range updates.UpdateIds {
			// ListUpdates returns completed updates too, so each one has to be described to
			// find out whether it is still running.
			update, err := eksAPI.DescribeUpdate(ctx, &awseks.DescribeUpdateInput{
				Name:     &clusterName,
				UpdateId: &id,
			})
			if err != nil {
				return nil, fmt.Errorf("describing update %q on cluster %q: %w", id, clusterName, err)
			}
			if update.Update != nil && update.Update.Status == ekstypes.UpdateStatusInProgress {
				inProgress = append(inProgress, id)
			}
		}
		if updates.NextToken == nil {
			return inProgress, nil
		}
		nextToken = updates.NextToken
	}
}
