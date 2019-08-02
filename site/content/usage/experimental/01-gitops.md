---
title: "GitOps"
weight: 190
url: usage/experimental/gitops-flux
---

## GitOps

[GitOps][gitops] is a way to do Kubernetes application delivery. It works by using Git as a single source of truth for
Kubernetes resources. With Git at the center of your delivery pipelines, developers can make pull requests to accelerate
and simplify application deployments and operations tasks to Kubernetes.

`eksctl` provides an easy way to set up gitops in an existing cluster with the `eksctl install flux` command.

[gitops]: https://www.weave.works/technologies/gitops/


### Installing Flux

> This is an experimental feature. To enable it, set the environment variable `EKSCTL_EXPERIMENTAL=true`. 
> Experimental features are not stable and their command name and flags may change.

Installing Flux on the cluster is the first step towards a gitops workflow. To install it, you need a Git repository 
and an existing cluster. Then run the following command:

```bash
EKSCTL_EXPERIMENTAL=true eksctl install flux --name <cluster_name> --region <region> --git-url=<git_repo> --git-email=<git_user_email>
```

Or use a config file:
```bash
EKSCTL_EXPERIMENTAL=true eksctl install flux -f examples/01-simple-cluster.yaml --git-url=git@github.com:weaveworks/cluster-1-gitops.git --git-email=johndoe+flux@weave.works
```

Full example:

```bash
EKSCTL_EXPERIMENTAL=true eksctl install flux --name cluster-1 --region eu-west-2 --git-url=git@github.com:weaveworks/cluster-1-gitops.git --git-email=johndoe+flux@weave.works
$ EKSCTL_EXPERIMENTAL=true eksctl install flux -f examples/01-simple-cluster.yaml --git-url=git@github.com:weaveworks/cluster-1-gitops.git --git-email=johndoe+flux@weave.works
[ℹ]  Cloning git@github.com:weaveworks/cluster-1-gitops.git
Cloning into '/tmp/eksctl-install-flux-clone310610186'...
warning: templates not found /home/johndoe/.git_template
remote: Enumerating objects: 3, done.
remote: Counting objects: 100% (3/3), done.
remote: Total 3 (delta 0), reused 0 (delta 0), pack-reused 0
Receiving objects: 100% (3/3), done.
[ℹ]  Writing Flux manifests
[ℹ]  Installing Flux into the cluster
[ℹ]  created "Namespace/flux"
[ℹ]  created "flux:ServiceAccount/flux"
[ℹ]  replaced "ClusterRole.rbac.authorization.k8s.io/flux"
[ℹ]  replaced "ClusterRoleBinding.rbac.authorization.k8s.io/flux"
[ℹ]  created "flux:Deployment.apps/flux"
[ℹ]  created "flux:Secret/flux-git-deploy"
[ℹ]  created "flux:Deployment.apps/memcached"
[ℹ]  created "flux:Service/memcached"
[ℹ]  Waiting for Flux to start
[!]  Flux is not ready yet (executing HTTP request: executing HTTP request: Post http://127.0.0.1:41483/api/flux/v9/git-repo-config: EOF), retrying ...
[ℹ]  Flux started successfully
[ℹ]  Committing and pushing Flux manifests to git@github.com:weaveworks/cluster-1-gitops.git
[master 4ff0d8f] Add Initial Flux configuration
 Author: Flux <johndoe+flux@weave.works>
 6 files changed, 255 insertions(+)
 create mode 100644 flux/flux-account.yaml
 create mode 100644 flux/flux-deployment.yaml
 create mode 100644 flux/flux-secret.yaml
 create mode 100644 flux/memcache-dep.yaml
 create mode 100644 flux/memcache-svc.yaml
 create mode 100644 flux/namespace.yaml
Counting objects: 9, done.
Delta compression using up to 4 threads.
Compressing objects: 100% (9/9), done.
Writing objects: 100% (9/9), 3.52 KiB | 1.17 MiB/s, done.
Total 9 (delta 0), reused 0 (delta 0)
To github.com:weaveworks/cluster-1-gitops.git
   d41f36d..4ff0d8f  master -> master
[ℹ]  Flux will operate properly only once it has SSH access to: git@github.com:weaveworks/cluster-1-gitops.git
[ℹ]  please add the following Flux public SSH key to your repository:
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCn9l5vVmcRbjZJmBG4GcYTK4w+8NjfMHOUr8W1w7E+PX8ono/cXsr9yohIPRUGKT1JSMXqwOTNNqYQoL6qbS7hGzOdO/IPW3JN1qvbBXLjBB8jo3op4KvudMuImBiE0dPB/mITk43t3WNbzZ33xlS9emtQdQlIno8HTFthohljcW5tUzdpC6Fv43fqt1EdHb8NtJz5oFbYbPuRf7swH0raxrhqKs4HW8VDVVqkROG2i0drg8rSalICbJX1YB3tgMYvP//f9uhWskXh5kuetS541I9gqtJD29pFibYQ1GwjfyAvkPBHTmumXdvb111111JWnfiiT7zCrdjYIEUt/9

```

At this point Flux should already be installed in the specified cluster. The only thing left to do would be to install
the SSH key. 

```bash
$ kubectl get pods --namespace flux
NAME                       READY   STATUS    RESTARTS   AGE
flux-699cc7f4cb-9qc45      1/1     Running   0          29m
memcached-958f745c-qdfgz   1/1     Running   0          29m
```

To finish the installation, Flux needs access to the repository. Configure your repo to allow write access to that ssh
key, for example, through the Deploy keys if it lives in GitHub.


### Adding a workload

To deploy a new workload on the cluster using GitOps just add a kubernetes manifest to the repository. After a few 
minutes you should see the resources appearing in the cluster.

### Further reading

To learn more about GitOps and Flux, check the [Flux documentation][flux]
 

[flux]: https://docs.fluxcd.io/en/latest/
