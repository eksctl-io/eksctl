package eks

import (
	"time"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"

	"github.com/weaveworks/eksctl/pkg/fargate"
	"github.com/weaveworks/eksctl/pkg/fargate/coredns"
	"github.com/weaveworks/eksctl/pkg/utils/retry"
	"github.com/weaveworks/eksctl/pkg/utils/strings"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

type fargateProfilesTask struct {
	info            string
	clusterProvider *ClusterProvider
	spec            *api.ClusterConfig
}

func (fpt *fargateProfilesTask) Describe() string { return fpt.info }

func (fpt *fargateProfilesTask) Do(errCh chan error) error {
	defer close(errCh)
	if err := DoCreateFargateProfiles(fpt.spec, fpt.clusterProvider); err != nil {
		return err
	}
	clientSet, err := fpt.clusterProvider.NewStdClientSet(fpt.spec)
	if err != nil {
		return err
	}
	if err := ScheduleCoreDNSOnFargateIfRelevant(fpt.spec, fpt.clusterProvider, clientSet); err != nil {
		return err
	}
	return nil
}

// DoCreateFargateProfiles creates fargate profiles as specified in the config
func DoCreateFargateProfiles(config *api.ClusterConfig, ctl *ClusterProvider) error {
	clusterName := config.Metadata.Name
	awsClient := fargate.NewClientWithWaitTimeout(clusterName, ctl.Provider.EKS(), ctl.Provider.WaitTimeout())
	for _, profile := range config.FargateProfiles {
		logger.Info("creating Fargate profile %q on EKS cluster %q", profile.Name, clusterName)

		// Default the pod execution role ARN to be the same as the cluster
		// role defined in CloudFormation:
		if profile.PodExecutionRoleARN == "" {
			profile.PodExecutionRoleARN = strings.EmptyIfNil(config.IAM.FargatePodExecutionRoleARN)
		}
		// Linearise the creation of Fargate profiles by passing
		// wait = true, as the API otherwise errors out with:
		//   ResourceInUseException: Cannot create Fargate Profile
		//   ${name2} because cluster ${clusterName} currently has
		//   Fargate profile ${name1} in status CREATING
		if err := awsClient.CreateProfile(profile, true); err != nil {
			return errors.Wrapf(err, "failed to create Fargate profile %q on EKS cluster %q", profile.Name, clusterName)
		}
		logger.Info("created Fargate profile %q on EKS cluster %q", profile.Name, clusterName)
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
