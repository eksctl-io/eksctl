package eks

import (
	"encoding/json"
	"fmt"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

type clusterConfigTask struct {
	info string
	spec *api.ClusterConfig
	call func(*api.ClusterConfig) error
}


func endpointsEqual(a, b api.ClusterEndpoints) bool {
	ajson, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		return false
	}
	bjson, err := json.MarshalIndent(b, "", "\t")
	if err != nil {
		return false
	}
	return string(ajson) == string(bjson)
}

//DefaultEndpointsMsg returns a message that the EndpointAccess is the same as the default
func DefaultEndpointsMsg(cfg *api.ClusterConfig) string {
	return fmt.Sprintf(
		"Cluster API server endpoint access will use default of {Pulic: true, Private false} for cluster %q in %q", cfg.Metadata.Name, cfg.Metadata.Region)
}

//CustomEndpointsMsg returns a message indicating the EndpointAccess given by the user
func CustomEndpointsMsg(cfg *api.ClusterConfig) string {
	return fmt.Sprintf(
		"Cluster API server endpoint access will use provided values {Public: %v, Private %v} for cluster %q in %q", *cfg.VPC.ClusterEndpoints.PublicAccess, *cfg.VPC.ClusterEndpoints.PrivateAccess, cfg.Metadata.Name, cfg.Metadata.Region)
}

//UpdateEndpointsMsg gives message indicating that they need to use eksctl utils to make this config
func UpdateEndpointsMsg(cfg *api.ClusterConfig) string {
	return fmt.Sprintf(
		"you can update endpoint access with `eksctl utils update-cluster-api-access --region=%s --name=%s --private-access bool --public-access bool", cfg.Metadata.Region, cfg.Metadata.Name)
}

func (t *clusterConfigTask) Describe() string { return t.info }

func (t *clusterConfigTask) Do(errs chan error) error {
	err := t.call(t.spec)
	close(errs)
	return err
}

// AppendExtraClusterConfigTasks returns all tasks for updating cluster configuration or nil if there are no tasks
func (c *ClusterProvider) AppendExtraClusterConfigTasks(cfg *api.ClusterConfig, tasks *manager.TaskTree) {
	newTasks := &manager.TaskTree{
		Parallel:  false,
		IsSubTask: true,
	}
	if !cfg.HasClusterCloudWatchLogging() {
		logger.Info("CloudWatch logging will not be enabled for cluster %q in %q", cfg.Metadata.Name, cfg.Metadata.Region)
		logger.Info("you can enable it with 'eksctl utils update-cluster-logging --region=%s --name=%s'", cfg.Metadata.Region, cfg.Metadata.Name)

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
	c.maybeAppendTasksForEndpointAccessUpdates(cfg, newTasks, tasks)
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

func (c *ClusterProvider) maybeAppendTasksForEndpointAccessUpdates(cfg *api.ClusterConfig, newTasks, tasks *manager.TaskTree) {
	// if a cluster config doesn't have the default api endpoint access, append a new task
	// so that we update the cluster with the new access configuration.  This is a
	// non-CloudFormation context, so we create a task to send it through the EKS API.
	// A caveat is that sending the default endpoint parameters for a cluster as an update will
	// return an error from the EKS API, so we must check for this before sending the request.
	if cfg.HasClusterEndpointAccess() && endpointsEqual(*cfg.VPC.ClusterEndpoints, *api.ClusterEndpointAccessDefaults()) {
		logger.Info(DefaultEndpointsMsg(cfg))
		logger.Info(UpdateEndpointsMsg(cfg))
	} else {
		logger.Info(CustomEndpointsMsg(cfg))
		newTasks.Append(&clusterConfigTask{
			info: "update cluster VPC endpoint access configuration",
			spec: cfg,
			call: c.UpdateClusterConfigForEndpoints,
		})
		tasks.Append(newTasks)
	}
}
