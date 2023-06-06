package eks

import (
	"context"
	"time"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"

	"github.com/weaveworks/eksctl/pkg/fargate/coredns"
	"github.com/weaveworks/eksctl/pkg/utils/apierrors"
	"github.com/weaveworks/eksctl/pkg/utils/retry"
	"github.com/weaveworks/eksctl/pkg/utils/strings"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fargate_client.go . FargateClient
type FargateClient interface {
	CreateProfile(ctx context.Context, profile *api.FargateProfile, waitForCreation bool) error
	ListProfiles(ctx context.Context) ([]string, error)
}

type fargateProfilesTask struct {
	info            string
	clusterProvider *ClusterProvider
	spec            *api.ClusterConfig
	manager         FargateClient
	ctx             context.Context
}

func (t *fargateProfilesTask) Describe() string { return t.info }

func (t *fargateProfilesTask) Do(errCh chan error) error {
	defer close(errCh)
	if err := DoCreateFargateProfiles(t.ctx, t.spec, t.manager); err != nil {
		return err
	}

	// Add delay after cluster creation to handle a race condition.
	timer := time.NewTimer(30 * time.Second)
	select {
	case <-timer.C:

	case <-t.ctx.Done():
		timer.Stop()
		return t.ctx.Err()
	}

	clientSet, err := t.clusterProvider.NewStdClientSet(t.spec)
	if err != nil {
		return errors.Wrap(err, "failed to get ClientSet")
	}
	if err := ScheduleCoreDNSOnFargateIfRelevant(t.spec, t.clusterProvider, clientSet); err != nil {
		return errors.Wrap(err, "failed to schedule core-dns on fargate")
	}
	return nil
}

// Check if the target Fargate Profile already exists
func targetFargateProfileExists(target string, profiles []string, clusterName string) bool {
	for _, profileName := range profiles {
		if target == profileName {
			logger.Info("Fargate profile %q already exists on EKS cluster %q", target, clusterName)
			return true
		}
	}
	return false
}

// DoCreateFargateProfiles creates fargate profiles as specified in the config
func DoCreateFargateProfiles(ctx context.Context, config *api.ClusterConfig, fargateClient FargateClient) error {
	clusterName := config.Metadata.Name

	// Get existing Farget profiles list
	existingProfiles, err := fargateClient.ListProfiles(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get Fargate Profile list")
	}

	for _, profile := range config.FargateProfiles {
		// Check if target Fargate Profile exists
		if targetFargateProfileExists(profile.Name, existingProfiles, clusterName) {
			continue
		}

		logger.Info("creating Fargate profile %q on EKS cluster %q", profile.Name, clusterName)

		// Default the pod execution role ARN to be the same as the cluster
		// role defined in CloudFormation:
		if profile.PodExecutionRoleARN == "" {
			profile.PodExecutionRoleARN = strings.EmptyIfNil(config.IAM.FargatePodExecutionRoleARN)
		}
		// Linearise the initial creation of Fargate profiles by passing
		// wait = true, as the API otherwise errors out with a ResourceInUseException
		//
		// In the case that a ResourceInUseException is thrown on a profile which was
		// created on an earlier call, we do not error but continue to the next one
		var inUseErr *ekstypes.ResourceInUseException
		err := fargateClient.CreateProfile(ctx, profile, true)
		switch {
		case err == nil:
			logger.Info("created Fargate profile %q on EKS cluster %q", profile.Name, clusterName)
		case errors.As(err, &inUseErr):
			logger.Info("either Fargate profile %q already exists on EKS cluster %q or another profile is being created/deleted, no action taken", profile.Name, clusterName)
		case apierrors.IsAccessDeniedError(err):
			return errors.Wrapf(err, "either account is not authorized to use Fargate or region %s is not supported", config.Metadata.Region)
		default:
			return errors.Wrapf(err, "failed to create Fargate profile %q on EKS cluster %q", profile.Name, clusterName)
		}
	}
	return nil
}

func ScheduleCoreDNSOnFargateIfRelevant(config *api.ClusterConfig, ctl *ClusterProvider, clientSet kubernetes.Interface) error {
	if coredns.IsSchedulableOnFargate(config.FargateProfiles) {
		scheduled, err := coredns.IsScheduledOnFargate(clientSet)
		if err != nil {
			return err
		}
		if !scheduled {
			if err := coredns.ScheduleOnFargate(clientSet); err != nil {
				return err
			}
			retryPolicy := &retry.TimingOutExponentialBackoff{
				Timeout:  ctl.AWSProvider.WaitTimeout(),
				TimeUnit: time.Second,
			}
			if err := coredns.WaitForScheduleOnFargate(clientSet, retryPolicy); err != nil {
				return err
			}
		}
	}
	return nil
}
