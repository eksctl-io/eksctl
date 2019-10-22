package eks

import (
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/addons"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

type clusterConfigTask struct {
	info string
	spec *api.ClusterConfig
	call func(*api.ClusterConfig) error
}

func (t *clusterConfigTask) Describe() string { return t.info }

func (t *clusterConfigTask) Do(errs chan error) error {
	err := t.call(t.spec)
	close(errs)
	return err
}

type vpcControllerTask struct {
	info            string
	clusterProvider *ClusterProvider
	spec            *api.ClusterConfig
}

func (v *vpcControllerTask) Describe() string { return v.info }

func (v *vpcControllerTask) Do(errCh chan error) error {
	defer close(errCh)
	rawClient, err := v.clusterProvider.NewRawClient(v.spec)
	if err != nil {
		return err
	}
	vpcController := addons.NewVPCController(rawClient, v.spec.Status, v.clusterProvider.Provider.Region(), false)
	if err := vpcController.Deploy(); err != nil {
		return errors.Wrap(err, "error installing VPC controller")
	}
	return nil
}

// AppendExtraClusterConfigTasks returns all tasks for updating cluster configuration or nil if there are no tasks
func (c *ClusterProvider) AppendExtraClusterConfigTasks(cfg *api.ClusterConfig, installVPCController bool, tasks *manager.TaskTree) {
	newTasks := &manager.TaskTree{
		Parallel:  false,
		IsSubTask: true,
	}
	if !cfg.HasClusterCloudWatchLogging() {
		logger.Info("CloudWatch logging will not be enabled for cluster %q in %q", cfg.Metadata.Name, cfg.Metadata.Region)
		logger.Info("you can enable it with 'eksctl utils update-cluster-logging --region=%s --cluster=%s'", cfg.Metadata.Region, cfg.Metadata.Name)

	} else {
		newTasks.Append(&clusterConfigTask{
			info: "update CloudWatch logging configuration",
			spec: cfg,
			call: c.UpdateClusterConfigForLogging,
		})
	}
	if api.IsEnabled(cfg.IAM.WithOIDC) {
		c.appendCreateTasksForIAMServiceAccounts(cfg, newTasks)
	}
	c.maybeAppendTasksForEndpointAccessUpdates(cfg, newTasks)
	if installVPCController {
		newTasks.Append(&vpcControllerTask{
			info:            "install Windows VPC controller",
			spec:            cfg,
			clusterProvider: c,
		})
	}
	if newTasks.Len() > 0 {
		tasks.Append(newTasks)
	}
}

func (c *ClusterProvider) appendCreateTasksForIAMServiceAccounts(cfg *api.ClusterConfig, tasks *manager.TaskTree) {
	// we don't have all the information to construct full iamoidc.OpenIDConnectManager now,
	// instead we just create a reference that gets updated when first task runs, and gets
	// used by this would be more elegant if it was all done via CloudFormation and we didn't
	// have to put wires across all the things like this; this whole function is needed because
	// we cannot manage certain EKS features with CloudFormation
	eatlyOIDC := &iamoidc.OpenIDConnectManager{}
	tasks.Append(&clusterConfigTask{
		info: "associate IAM OIDC provider",
		spec: cfg,
		call: func(cfg *api.ClusterConfig) error {
			if err := c.RefreshClusterStatus(cfg); err != nil {
				return err
			}
			oidc, err := c.NewOpenIDConnectManager(cfg)
			if err != nil {
				return err
			}
			if err := oidc.CreateProvider(); err != nil {
				return err
			}
			*eatlyOIDC = *oidc
			return nil
		},
	})

	clientSet := &kubernetes.CallbackClientSet{
		Callback: func() (kubernetes.Interface, error) {
			return c.NewStdClientSet(cfg)
		},
	}

	// as this is non-CloudFormation context, we need to construct a new stackManager,
	// given a clientSet getter and OpenIDConnectManager reference we can build out
	// the list of tasks for each of the service accounts that need to be created
	newTasks := c.NewStackManager(cfg).NewTasksToCreateIAMServiceAccounts(cfg.IAM.ServiceAccounts, eatlyOIDC, clientSet)
	newTasks.IsSubTask = true
	tasks.Append(newTasks)
}

func (c *ClusterProvider) maybeAppendTasksForEndpointAccessUpdates(cfg *api.ClusterConfig, tasks *manager.TaskTree) {
	// if a cluster config doesn't have the default api endpoint access, append a new task
	// so that we update the cluster with the new access configuration.  This is a
	// non-CloudFormation context, so we create a task to send it through the EKS API.
	// A caveat is that sending the default endpoint parameters for a cluster as an update will
	// return an error from the EKS API, so we must check for this before sending the request.
	if cfg.HasClusterEndpointAccess() && api.EndpointsEqual(*cfg.VPC.ClusterEndpoints, *api.ClusterEndpointAccessDefaults()) {
		// No tasks to append here as there's no updates to make.
		logger.Info(cfg.DefaultEndpointsMsg())
	} else {
		logger.Info(cfg.CustomEndpointsMsg())

		tasks.Append(&clusterConfigTask{
			info: "update cluster VPC endpoint access configuration",
			spec: cfg,
			call: c.UpdateClusterConfigForEndpoints,
		})
	}
}
