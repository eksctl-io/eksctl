package utils

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"

	defaultaddons "github.com/weaveworks/eksctl/pkg/addons/default"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/utils/apierrors"
)

type handleAddonUpdate func(*kubernetes.RawClient, defaultaddons.AddonVersionDescriber) (updateRequired bool, err error)

func updateAddon(ctx context.Context, cmd *cmdutils.Cmd, addonName string, handleUpdate handleAddonUpdate) error {
	if err := cmdutils.NewMetadataLoader(cmd).Load(); err != nil {
		return err
	}
	ctl, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}
	if ok, err := ctl.CanUpdate(cmd.ClusterConfig); !ok {
		return err
	}

	eksAPI := ctl.AWSProvider.EKS()
	switch _, err := eksAPI.DescribeAddon(ctx, &eks.DescribeAddonInput{
		AddonName:   aws.String(addonName),
		ClusterName: aws.String(cmd.ClusterConfig.Metadata.Name),
	}); {
	case err == nil:
		return fmt.Errorf("addon %s is installed as a managed EKS addon; to update it, use `eksctl update addon` instead", addonName)
	case apierrors.IsNotFoundError(err):

	default:
		return fmt.Errorf("error describing addon %s: %w", addonName, err)
	}

	rawClient, err := ctl.NewRawClient(cmd.ClusterConfig)
	if err != nil {
		return err
	}
	updateRequired, err := handleUpdate(rawClient, eksAPI)
	if err != nil {
		return err
	}
	cmdutils.LogPlanModeWarning(cmd.Plan && updateRequired)
	return nil
}
