package eks

import (
	"time"

	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"

	"github.com/weaveworks/eksctl/pkg/fargate"
	"github.com/weaveworks/eksctl/pkg/fargate/coredns"
	"github.com/weaveworks/eksctl/pkg/utils/retry"
	"github.com/weaveworks/eksctl/pkg/utils/strings"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fargate_client.go . FargateClient
type FargateClient interface {
	CreateProfile(profile *api.FargateProfile, waitForCreation bool) error
}

type fargateProfilesTask struct {
	info            string
	clusterProvider *ClusterProvider
	spec            *api.ClusterConfig
	manager         FargateClient
}

func (fpt *fargateProfilesTask) Describe() string { return fpt.info }

func (fpt *fargateProfilesTask) Do(errCh chan error) error {
	defer close(errCh)
	if err := DoCreateFargateProfiles(fpt.spec, fpt.manager); err != nil {
		return err
	}

	// Add delay after cluster creation to handle a race condition
	time.Sleep(30 * time.Second)

	clientSet, err := fpt.clusterProvider.NewStdClientSet(fpt.spec)
	if err != nil {
		return errors.Wrap(err, "failed to get ClientSet")
	}
	if err := ScheduleCoreDNSOnFargateIfRelevant(fpt.spec, fpt.clusterProvider, clientSet); err != nil {
		return errors.Wrap(err, "failed to schedule core-dns on fargate")
	}
	return nil
}

// DoCreateFargateProfiles creates fargate profiles as specified in the config
func DoCreateFargateProfiles(config *api.ClusterConfig, fargateClient FargateClient) error {
	clusterName := config.Metadata.Name
	for _, profile := range config.FargateProfiles {
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
		var e *eks.ResourceInUseException
		err := fargateClient.CreateProfile(profile, true)
		switch {
		case err == nil:
			logger.Info("created Fargate profile %q on EKS cluster %q", profile.Name, clusterName)
		case errors.As(err, &e):
			logger.Info("Either Fargate profile %q already exists on EKS cluster %q or another profile is being created/deleted, no action taken", profile.Name, clusterName)
		case fargate.IsUnauthorizedError(err):
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
				Timeout:  ctl.Provider.WaitTimeout(),
				TimeUnit: time.Second,
			}
			if err := coredns.WaitForScheduleOnFargate(clientSet, retryPolicy); err != nil {
				return err
			}
		}
	}
	return nil
}
