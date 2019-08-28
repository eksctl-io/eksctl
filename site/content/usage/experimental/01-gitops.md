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
and an existing EKS cluster. Then run the following command:

```console
EKSCTL_EXPERIMENTAL=true eksctl install flux --cluster=<cluster_name> --region=<region> --git-url=<git_repo> --git-email=<git_user_email>
```

Or use a config file:
```console
EKSCTL_EXPERIMENTAL=true eksctl install flux -f examples/01-simple-cluster.yaml --git-url=git@github.com:weaveworks/cluster-1-gitops.git --git-email=johndoe+flux@weave.works
```

Note that, by default, `eksctl install flux` installs [Helm](https://helm.sh/) server components to the cluster (it
installs [Tiller](https://helm.sh/docs/glossary/#tiller) and the [Flux Helm Operator](https://github.com/fluxcd/helm-operator)). To
disable the installation of the Helm server components, pass the flag `--with-helm=false`.

Full example:

```console
$ EKSCTL_EXPERIMENTAL=true ./eksctl install flux --cluster=cluster-1 --region=eu-west-2  --git-url=git@github.com:weaveworks/cluster-1-gitops.git  --git-email=johndoe+flux@weave.works--namespace=flux
[ℹ]  Generating public key infrastructure for the Helm Operator and Tiller
[ℹ]    this may take up to a minute, please be patient
[!]  Public key infrastructure files were written into directory "/var/folders/zt/sh1tk7ts24sc6dybr5z9qtfh0000gn/T/eksctl-helm-pki330304977"
[!]  please move the files into a safe place or delete them
[ℹ]  Generating manifests
[ℹ]  Cloning git@github.com:weaveworks/cluster-1-gitops.git
Cloning into '/var/folders/zt/sh1tk7ts24sc6dybr5z9qtfh0000gn/T/eksctl-install-flux-clone-142184188'...
remote: Enumerating objects: 74, done.
remote: Counting objects: 100% (74/74), done.
remote: Compressing objects: 100% (55/55), done.
remote: Total 74 (delta 19), reused 69 (delta 17), pack-reused 0
Receiving objects: 100% (74/74), 30.57 KiB | 381.00 KiB/s, done.
Resolving deltas: 100% (19/19), done.
[ℹ]  Writing Flux manifests
[ℹ]  Applying manifests
[ℹ]  created "Namespace/flux"
[ℹ]  created "flux:Secret/flux-git-deploy"
[ℹ]  created "flux:Deployment.apps/memcached"
[ℹ]  created "flux:ServiceAccount/flux"
[ℹ]  created "ClusterRole.rbac.authorization.k8s.io/flux"
[ℹ]  created "ClusterRoleBinding.rbac.authorization.k8s.io/flux"
[ℹ]  created "flux:ConfigMap/flux-helm-tls-ca-config"
[ℹ]  created "flux:Deployment.extensions/tiller-deploy"
[ℹ]  created "flux:Service/tiller-deploy"
[ℹ]  created "CustomResourceDefinition.apiextensions.k8s.io/helmreleases.helm.fluxcd.io"
[ℹ]  created "flux:ServiceAccount/tiller"
[ℹ]  created "ClusterRoleBinding.rbac.authorization.k8s.io/tiller"
[ℹ]  created "flux:ServiceAccount/helm"
[ℹ]  created "flux:Role.rbac.authorization.k8s.io/tiller-user"
[ℹ]  created "kube-system:RoleBinding.rbac.authorization.k8s.io/tiller-user-binding"
[ℹ]  created "flux:Deployment.apps/flux"
[ℹ]  created "flux:Service/memcached"
[ℹ]  created "flux:Deployment.apps/flux-helm-operator"
[ℹ]  created "flux:ServiceAccount/flux-helm-operator"
[ℹ]  created "ClusterRole.rbac.authorization.k8s.io/flux-helm-operator"
[ℹ]  created "ClusterRoleBinding.rbac.authorization.k8s.io/flux-helm-operator"
[ℹ]  Applying Helm TLS Secret(s)
[ℹ]  created "flux:Secret/flux-helm-tls-cert"
[ℹ]  created "flux:Secret/tiller-secret"
[!]  Note: certificate secrets aren't added to the Git repository for security reasons
[ℹ]  Waiting for Helm Operator to start
ERROR: logging before flag.Parse: E0820 16:05:12.218007   98823 portforward.go:331] an error occurred forwarding 60356 -> 3030: error forwarding port 3030 to pod b1a872e7e6a7f86567488d66c1a880fcfa26179143115b102041e0ee77fe6f9e, uid : exit status 1: 2019/08/20 14:05:12 socat[2873] E connect(5, AF=2 127.0.0.1:3030, 16): Connection refused
[!]  Helm Operator is not ready yet (Get http://127.0.0.1:60356/healthz: EOF), retrying ...
[ℹ]  Helm Operator started successfully
[ℹ]  Waiting for Flux to start
[ℹ]  Flux started successfully
[ℹ]  Committing and pushing manifests to git@github.com:weaveworks/cluster-1-gitops.git
[master ec43024] Add Initial Flux configuration
 Author: Flux <johndoe+flux@weave.works>
14 files changed, 694 insertions(+)
Enumerating objects: 11, done.
Counting objects: 100% (11/11), done.
Delta compression using up to 4 threads
Compressing objects: 100% (6/6), done.
Writing objects: 100% (6/6), 2.09 KiB | 2.09 MiB/s, done.
Total 6 (delta 3), reused 0 (delta 0)
remote: Resolving deltas: 100% (3/3), completed with 3 local objects.
To github.com:weaveworks/cluster-1-gitops.git
   5fe1eb8..ec43024  master -> master
[ℹ]  Flux will only operate properly once it has write-access to the Git repository
[ℹ]  please configure git@github.com:weaveworks/cluster-1-gitops.git  so that the following Flux SSH public key has write access to it
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDYYsPuHzo1L29u3zhr4uAOF29HNyMcS8zJmOTDNZC4EiIwa5BXgg/IBDKudxQ+NBJ7mknPlNv17cqo4ncEq1xiQidfaUawwx3xxtDkZWam5nCBMXEJwkr4VXx/6QQ9Z1QGXpaFwdoVRcY/kM4NaxM54pEh5m43yeqkcpRMKraE0EgbdqFNNARN8rIEHY/giDorCrXp7e6AbzBgZSvc/in7Ul9FQhJ6K4+7QuMFpJt3O/N8KDumoTG0e5ssJGp5L1ugIqhzqvbHdmHVfnXsEvq6cR1SJtYKi2GLCscypoF3XahfjK+xGV/92a1E7X+6fHXSq+bdOKfBc4Z3f9NBwz0v

```

At this point Flux and the Helm server components should be installed in the specified cluster. The only thing left to
do is to give Flux write access to the repository. Configure your repository to allow write access to that ssh key,
for example, through the Deploy keys if it lives in GitHub.

```console
$ kubectl get pods --namespace flux
NAME                       READY   STATUS    RESTARTS   AGE
flux-699cc7f4cb-9qc45      1/1     Running   0          29m
memcached-958f745c-qdfgz   1/1     Running   0          29m
```


#### Adding a workload

To deploy a new workload on the cluster using GitOps just add a kubernetes manifest to the repository. After a few
minutes you should see the resources appearing in the cluster.

#### Further reading

To learn more about GitOps and Flux, check the [Flux documentation][flux]


### Installing components from a quickstart

`eksctl` provides an application development quickstart profile which can install the following components in your
cluster:
  - Metrics Server
  - Prometheus
  - Grafana
  - Kubernetes Dashboard
  - FluentD with connection to CloudWatch logs
  - CNI, present by default in EKS clusters
  - Cluster Autoscaler
  - ALB ingress controller
  - Podinfo as a demo application

To install those components the command `generate profile` can be used:

```console
EKSCTL_EXPERIMENTAL=true eksctl generate profile --config-file=<cluster_config_file> --git-url git@github.com:weaveworks/eks-gitops-example.git --profile-path <output_directory>
```

For example:

```
$ EKSCTL_EXPERIMENTAL=true eksctl generate profile  --config-file 01-simple-cluster.yaml --git-url git@github.com:weaveworks/eks-gitops-example.git --profile-path my-gitops-repo/base/
[ℹ]  cloning repository "git@github.com:weaveworks/eks-gitops-example.git":master
Cloning into '/tmp/quickstart-224631067'...
warning: templates not found /home/.../.git_template
remote: Enumerating objects: 75, done.
remote: Counting objects: 100% (75/75), done.
remote: Compressing objects: 100% (63/63), done.
remote: Total 75 (delta 25), reused 49 (delta 11), pack-reused 0
Receiving objects: 100% (75/75), 19.02 KiB | 1.19 MiB/s, done.
Resolving deltas: 100% (25/25), done.
[ℹ]  processing template files in repository
[ℹ]  writing new manifests to "base/"

$ tree my-gitops-repo/base
base
├── amazon-cloudwatch
│   ├── 0-namespace.yaml
│   ├── config-map.yaml
│   ├── daemonset.yaml
│   ├── fluentd-config-map.yaml
│   ├── fluentd.yml
│   └── service-account.yaml
├── demo
│   ├── 00-namespace.yaml
│   └── helm-release.yaml
├── kube-system
│   ├── aws-alb-ingress-controller
│   │   ├── deployment.yaml
│   │   └── rbac-role.yaml
│   ├── cluster-autoscaler-autodiscover.yaml
│   └── kubernetes-dashboard.yaml
├── LICENSE
├── monitoring
│   ├── 00-namespace.yaml
│   ├── metrics-server.yaml
│   └── prometheus.yaml
└── README.md
```

After running the command, add, commit and push the files:

```bash
cd my-gitops-repo/
git add .
git commit -m "Add application development quickstart components"
git push origin master
```

After a few minutes, Flux and Helm should have installed all the components in your cluster.


[flux]: https://docs.fluxcd.io/en/latest/
