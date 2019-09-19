---
title: "gitops"
weight: 190
url: usage/experimental/gitops-flux
---

## gitops

[gitops][gitops] is a way to do Kubernetes application delivery. It works by using Git as a single source of truth for
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
$ EKSCTL_EXPERIMENTAL=true ./eksctl install flux --cluster=cluster-1 --region=eu-west-2  --git-url=git@github.com:weaveworks/cluster-1-gitops.git  --git-email=johndoe+flux@weave.works --namespace=flux
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

To deploy a new workload on the cluster using gitops just add a kubernetes manifest to the repository. After a few
minutes you should see the resources appearing in the cluster.

#### Further reading

To learn more about gitops and Flux, check the [Flux documentation][flux]


### Installing components from a Quick Start profile

`eksctl` provides an application development Quick Star profile which can install the following components in your
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
EKSCTL_EXPERIMENTAL=true eksctl generate profile --config-file=<cluster_config_file> --git-url git@github.com:weaveworks/eks-quickstart-app-dev.git --profile-path <output_directory>
```

For example:

```
$ EKSCTL_EXPERIMENTAL=true eksctl generate profile  --config-file 01-simple-cluster.yaml --git-url git@github.com:weaveworks/eks-quickstart-app-dev.git --profile-path my-gitops-repo/base/
[ℹ]  cloning repository "git@github.com:weaveworks/eks-quickstart-app-dev.git":master
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
$ tree base/
base/
├── amazon-cloudwatch
│   ├── cloudwatch-agent-configmap.yaml
│   ├── cloudwatch-agent-daemonset.yaml
│   ├── cloudwatch-agent-rbac.yaml
│   ├── fluentd-configmap-cluster-info.yaml
│   ├── fluentd-configmap-fluentd-config.yaml
│   ├── fluentd-daemonset.yaml
│   └── fluentd-rbac.yaml
├── demo
│   └── helm-release.yaml
├── kubernetes-dashboard
│   ├── dashboard-metrics-scraper-deployment.yaml
│   ├── dashboard-metrics-scraper-service.yaml
│   ├── kubernetes-dashboard-configmap.yaml
│   ├── kubernetes-dashboard-deployment.yaml
│   ├── kubernetes-dashboard-rbac.yaml
│   ├── kubernetes-dashboard-secrets.yaml
│   └── kubernetes-dashboard-service.yaml
├── kube-system
│   ├── alb-ingress-controller-deployment.yaml
│   ├── alb-ingress-controller-rbac.yaml
│   ├── cluster-autoscaler-deployment.yaml
│   └── cluster-autoscaler-rbac.yaml
├── LICENSE
├── monitoring
│   ├── metrics-server.yaml
│   └── prometheus-operator.yaml
├── namespaces
│   ├── amazon-cloudwatch.yaml
│   ├── demo.yaml
│   ├── kubernetes-dashboard.yaml
│   └── monitoring.yaml
└── README.md

```

After running the command, add, commit and push the files:

```bash
cd my-gitops-repo/
git add .
git commit -m "Add application development quick start components"
git push origin master
```

After a few minutes, Flux and Helm should have installed all the components in your cluster.

## Setting up gitops in a repo from a Quick Start

Configuring gitops can be done easily with eksctl. The command `eksctl enable profile` takes an existing EKS cluster and
an empty repository and sets them up with gitops and a specified Quick Start profile. This means that with one command
the cluster will have all the components provided by the Quick Start profile installed in the cluster and you can enjoy
the advantages of gitops moving forward.

The basic command usage looks like this:

```console
EKSCTL_EXPERIMENTAL=true eksctl enable profile --cluster <cluster-name> --region <region> --git-url=<url_to_your_repo> app-dev
```


This command will clone the specified repository in your current working directory and then it will follow these steps:

  1. install Flux, Helm and Tiller in the cluster and add the manifests of those components into the `flux/` folder in your repo
  2. add the component manifests of the Quick Start profile to your repository inside the `base/` folder
  3. commit the Quick Start files and push the changes to the origin remote
  4. once you have given read and write access to your repository to the the SSH key printed by the command, Flux will install the components from the `base/` folder into your cluster

Example:

```
$ EKSCTL_EXPERIMENTAL=true eksctl enable profile --cluster production-cluster --region eu-north-1 --git-url=git@github.com:myorg/production-kubernetes --output-path=/tmp/gitops-repos/  app-dev
[ℹ]  Generating public key infrastructure for the Helm Operator and Tiller
[ℹ]    this may take up to a minute, please be patient
[!]  Public key infrastructure files were written into directory "/tmp/eksctl-helm-pki786744152"
[!]  please move the files into a safe place or delete them
[ℹ]  Generating manifests
[ℹ]  Cloning git@github.com:myorg/production-kubernetes
Cloning into '/tmp/eksctl-install-flux-clone-615092439'...
remote: Enumerating objects: 114, done.
remote: Counting objects: 100% (114/114), done.
remote: Compressing objects: 100% (94/94), done.
remote: Total 114 (delta 36), reused 93 (delta 17), pack-reused 0
Receiving objects: 100% (114/114), 31.43 KiB | 4.49 MiB/s, done.
Resolving deltas: 100% (36/36), done.
[ℹ]  Writing Flux manifests
[ℹ]  Applying manifests
[ℹ]  created "Namespace/flux"
[ℹ]  created "flux:Deployment.apps/flux-helm-operator"
[ℹ]  created "flux:Deployment.apps/flux"
[ℹ]  created "flux:Deployment.apps/memcached"
[ℹ]  created "flux:ConfigMap/flux-helm-tls-ca-config"
[ℹ]  created "flux:Deployment.extensions/tiller-deploy"
[ℹ]  created "flux:Service/tiller-deploy"
[ℹ]  created "CustomResourceDefinition.apiextensions.k8s.io/helmreleases.helm.fluxcd.io"
[ℹ]  created "flux:Service/memcached"
[ℹ]  created "flux:ServiceAccount/flux"
[ℹ]  created "ClusterRole.rbac.authorization.k8s.io/flux"
[ℹ]  created "ClusterRoleBinding.rbac.authorization.k8s.io/flux"
[ℹ]  created "flux:Secret/flux-git-deploy"
[ℹ]  created "flux:ServiceAccount/flux-helm-operator"
[ℹ]  created "ClusterRole.rbac.authorization.k8s.io/flux-helm-operator"
[ℹ]  created "ClusterRoleBinding.rbac.authorization.k8s.io/flux-helm-operator"
[ℹ]  created "flux:ServiceAccount/tiller"
[ℹ]  created "ClusterRoleBinding.rbac.authorization.k8s.io/tiller"
[ℹ]  created "flux:ServiceAccount/helm"
[ℹ]  created "flux:Role.rbac.authorization.k8s.io/tiller-user"
[ℹ]  created "kube-system:RoleBinding.rbac.authorization.k8s.io/tiller-user-binding"
[ℹ]  Applying Helm TLS Secret(s)
[ℹ]  created "flux:Secret/flux-helm-tls-cert"
[ℹ]  created "flux:Secret/tiller-secret"
[!]  Note: certificate secrets aren't added to the Git repository for security reasons
[ℹ]  Waiting for Helm Operator to start
ERROR: logging before flag.Parse: E0822 14:45:28.440236   17028 portforward.go:331] an error occurred forwarding 44915 -> 3030: error forwarding port 3030 to pod 2f6282bf597b345b3ffad8a0447bdd8515d91060335456591d759ad87a976ed2, uid : exit status 1: 2019/08/22 12:45:28 socat[8131] E connect(5, AF=2 127.0.0.1:3030, 16): Connection refused
[!]  Helm Operator is not ready yet (Get http://127.0.0.1:44915/healthz: EOF), retrying ...
[!]  Helm Operator is not ready yet (Get http://127.0.0.1:44915/healthz: EOF), retrying ...
[ℹ]  Helm Operator started successfully
[ℹ]  see https://docs.fluxcd.io/projects/helm-operator for details on how to use the Helm Operator
[ℹ]  Waiting for Flux to start
[ℹ]  Flux started successfully
[ℹ]  see https://docs.fluxcd.io/projects/flux for details on how to use Flux
[ℹ]  Committing and pushing manifests to git@github.com:myorg/production-kubernetes
[master 0985830] Add Initial Flux configuration
 Author: Flux <>
 13 files changed, 727 insertions(+)
 create mode 100644 flux/flux-account.yaml
 create mode 100644 flux/flux-deployment.yaml
 create mode 100644 flux/flux-helm-operator-account.yaml
 create mode 100644 flux/flux-helm-release-crd.yaml
 create mode 100644 flux/flux-namespace.yaml
 create mode 100644 flux/flux-secret.yaml
 create mode 100644 flux/helm-operator-deployment.yaml
 create mode 100644 flux/memcache-dep.yaml
 create mode 100644 flux/memcache-svc.yaml
 create mode 100644 flux/tiller-ca-cert-configmap.yaml
 create mode 100644 flux/tiller-dep.yaml
 create mode 100644 flux/tiller-rbac.yaml
 create mode 100644 flux/tiller-svc.yaml
Counting objects: 16, done.
Delta compression using up to 4 threads.
Compressing objects: 100% (15/15), done.
Writing objects: 100% (16/16), 8.23 KiB | 8.23 MiB/s, done.
Total 16 (delta 1), reused 12 (delta 1)
remote: Resolving deltas: 100% (1/1), done.
To github.com:myorg/production-kubernetes
   3ea1fdc..0985830  master -> master
[ℹ]  Flux will only operate properly once it has write-access to the Git repository
[ℹ]  please configure git@github.com:myorg/production-kubernetes so that the following Flux SSH public key has write access to it
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAdDG4LAEiEOTbT3XVL5sYf0Hy7T30PG2sFReIwrylR7syA+IU9GPf7azgjjbzbQc/5BXTx2E0GotrzDkvCNScuYfw7wXKK87yr5jhPudpNubK9bFsKKwOj7wxO2XsUOceKVRhTKP7VJgpAliCCPK288HvQzIZfWEgbDQjhE0EnFgZVYXKkgye2Cc3MkwiYuZJtuynxipb5rPrY/3Kjywk/vWxLeZ/hvv58mZSdRQwX6zbGGW1h70QA47B+W2076MBQQ1t0H0KKctuS8A1/n+aKjpD4Ne6lXqHDhqi25SBhJxK3zEXhskS9DMW8DYi1xHT2MCjE8HhiVBMRIITyTox
Cloning into '/tmp/gitops-repos/flux-test-3'...
remote: Enumerating objects: 118, done.
remote: Counting objects: 100% (118/118), done.
remote: Compressing objects: 100% (98/98), done.
remote: Total 118 (delta 37), reused 96 (delta 17), pack-reused 0
Receiving objects: 100% (118/118), 33.15 KiB | 1.44 MiB/s, done.
Resolving deltas: 100% (37/37), done.
[ℹ]  cloning repository "git@github.com:weaveworks/eks-quickstart-app-dev.git":master
Cloning into '/tmp/quickstart-365477450'...
remote: Enumerating objects: 127, done.
remote: Counting objects: 100% (127/127), done.                                         
remote: Compressing objects: 100% (95/95), done.
remote: Total 127 (delta 53), reused 92 (delta 30), pack-reused 0
Receiving objects: 100% (127/127), 30.20 KiB | 351.00 KiB/s, done.
Resolving deltas: 100% (53/53), done.
[ℹ]  processing template files in repository
[ℹ]  writing new manifests to "/tmp/gitops-repos/flux-test-3/base"
[master d0810f7] Add app-dev quickstart components
 Author: Flux <>
 27 files changed, 1207 insertions(+)
 create mode 100644 base/LICENSE
 create mode 100644 base/README.md
 create mode 100644 base/amazon-cloudwatch/cloudwatch-agent-configmap.yaml
 create mode 100644 base/amazon-cloudwatch/cloudwatch-agent-daemonset.yaml
 create mode 100644 base/amazon-cloudwatch/cloudwatch-agent-rbac.yaml
 create mode 100644 base/amazon-cloudwatch/fluentd-configmap-cluster-info.yaml
 create mode 100644 base/amazon-cloudwatch/fluentd-configmap-fluentd-config.yaml
 create mode 100644 base/amazon-cloudwatch/fluentd-daemonset.yaml
 create mode 100644 base/amazon-cloudwatch/fluentd-rbac.yaml
 create mode 100644 base/demo/helm-release.yaml
 create mode 100644 base/kube-system/alb-ingress-controller-deployment.yaml
 create mode 100644 base/kube-system/alb-ingress-controller-rbac.yaml
 create mode 100644 base/kube-system/cluster-autoscaler-deployment.yaml
 create mode 100644 base/kube-system/cluster-autoscaler-rbac.yaml
 create mode 100644 base/kubernetes-dashboard/dashboard-metrics-scraper-deployment.yaml
 create mode 100644 base/kubernetes-dashboard/dashboard-metrics-scraper-service.yaml
 create mode 100644 base/kubernetes-dashboard/kubernetes-dashboard-configmap.yaml
 create mode 100644 base/kubernetes-dashboard/kubernetes-dashboard-deployment.yaml
 create mode 100644 base/kubernetes-dashboard/kubernetes-dashboard-rbac.yaml
 create mode 100644 base/kubernetes-dashboard/kubernetes-dashboard-secrets.yaml
 create mode 100644 base/kubernetes-dashboard/kubernetes-dashboard-service.yaml
 create mode 100644 base/monitoring/metrics-server.yaml
 create mode 100644 base/monitoring/prometheus-operator.yaml
 create mode 100644 base/namespaces/amazon-cloudwatch.yaml
 create mode 100644 base/namespaces/demo.yaml
 create mode 100644 base/namespaces/kubernetes-dashboard.yaml
 create mode 100644 base/namespaces/monitoring.yaml
Counting objects: 36, done.
Delta compression using up to 4 threads.
Compressing objects: 100% (27/27), done.
Writing objects: 100% (36/36), 11.17 KiB | 3.72 MiB/s, done.
Total 36 (delta 8), reused 27 (delta 8)
remote: Resolving deltas: 100% (8/8), done.
To github.com:myorg/production-kubernetes
   0985830..d0810f7  master -> master

```

Now the ssh key printed above:

```
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAdDG4LAEiEOTbT3XVL5sYf0Hy7T30PG2sFReIwrylR7syA+IU9GPf7azgjjbzbQc/5BXTx2E0GotrzDkvCNScuYfw7wXKK87yr5jhPudpNubK9bFsKKwOj7wxO2XsUOceKVRhTKP7VJgpAliCCPK288HvQzIZfWEgbDQjhE0EnFgZVYXKkgye2Cc3MkwiYuZJtuynxipb5rPrY/3Kjywk/vWxLeZ/hvv58mZSdRQwX6zbGGW1h70QA47B+W2076MBQQ1t0H0KKctuS8A1/n+aKjpD4Ne6lXqHDhqi25SBhJxK3zEXhskS9DMW8DYi1xHT2MCjE8HhiVBMRIITyTox
```

needs to be added as a deploy key to the chosen Github repository, in this case `github.com:myorg/production-kubernetes`.
Once that is done, Flux will pick up the changes in the repository with the Quick Start components and deploy them to the
 cluster. After a couple of minutes the pods should appear in the cluster:
 

```
$ kube get pods --all-namespaces 
NAMESPACE              NAME                                                      READY   STATUS                       RESTARTS   AGE
amazon-cloudwatch      cloudwatch-agent-qtdmc                                    1/1     Running                      0           4m28s
amazon-cloudwatch      fluentd-cloudwatch-4rwwr                                  1/1     Running                      0           4m28s
demo                   podinfo-75b8547f78-56dll                                  1/1     Running                      0          103s
flux                   flux-56b5664cdd-nfzx2                                     1/1     Running                      0          11m
flux                   flux-helm-operator-6bc7c85bb5-l2nzn                       1/1     Running                      0          11m
flux                   memcached-958f745c-dqllc                                  1/1     Running                      0          11m
flux                   tiller-deploy-7ccc4b4d45-w2mrt                            1/1     Running                      0          11m
kube-system            alb-ingress-controller-6b64bcbbd8-6l7kf                   1/1     Running                      0          4m28s
kube-system            aws-node-l49ct                                            1/1     Running                      0          14m
kube-system            cluster-autoscaler-5b8c96cd98-26z5f                       1/1     Running                      0          4m28s
kube-system            coredns-7d7755744b-4jkp6                                  1/1     Running                      0          21m
kube-system            coredns-7d7755744b-ls5d9                                  1/1     Running                      0          21m
kube-system            kube-proxy-wllff                                          1/1     Running                      0          14m
kubernetes-dashboard   dashboard-metrics-scraper-f7b5dbf7d-rm5z7                 1/1     Running                      0          4m28s
kubernetes-dashboard   kubernetes-dashboard-7447f48f55-94rhg                     1/1     Running                      0          4m28s
monitoring             alertmanager-prometheus-operator-alertmanager-0           2/2     Running                      0          78s
monitoring             metrics-server-7dfc675884-q9qps                           1/1     Running                      0          4m24s
monitoring             prometheus-operator-grafana-9bb769cf-pjk4r                2/2     Running                      0          89s
monitoring             prometheus-operator-kube-state-metrics-79f476bff6-r9m2s   1/1     Running                      0          89s
monitoring             prometheus-operator-operator-58fcb66576-6dwpg             1/1     Running                      0          89s
monitoring             prometheus-operator-prometheus-node-exporter-tllwl        1/1     Running                      0          89s
monitoring             prometheus-prometheus-operator-prometheus-0               3/3     Running                      1          72s
```


All CLI arguments:

| Flag                         | Default value | Type   | Required       | Use                                                           |
|------------------------------|---------------|--------|----------------|---------------------------------------------------------------|
| `--cluster`                  |               | string | required       | name of the EKS cluster to add the nodegroup to               |
| `--name`                     |               | string | required       | name or URL of the Quick Start profile. For example, app-dev  |
| <name positional argument>   |               | string | required       | same as `--name`                                              |
| `--git-url`                  |               | string | required       | URL                                                           |
| `--git-branch`               | master        | string | optional       | Git branch                                                    |
| `--output-path`              | ./            | string | optional       | Path                                                          |
| `--git-user`                 | Flux          | string | optional       | Username                                                      |
| `--git-email`                |               | string | optional       | Email                                                         |
| `--git-private-ssh-key-path` |               | string | optional       | Optional path to the private SSH key to use with Git          |


## Creating your own Quick Start profile

A Quick Start profile is a Git repository that contains Kubernetes manifests that can be installed in a cluster using
gitops (through [Flux][flux]).

These manifests, will probably need some information about the cluster they will be installed in, such as the cluster
name or the AWS region. That's why they are templated using [Go templates][go-templates].

> Please bear in mind that this is an experimental feature and therefore the chosen templating technology can be changed
before the feature becomes stable.

The variables that can be templated are shown below:

| Name                | Template               |
|---------------------|------------------------|
| cluster name        | `{{ .ClusterName }}`   |
| cluster region      | `{{ .Region }}`        |


For example, we could create a config map using these variables:

```yaml
apiVersion: v1
data:
  cluster.name: {{ .ClusterName }}
  logs.region: {{ .Region }}
kind: ConfigMap
metadata:
  name: cluster-info
  namespace: my-namespace
```

Write this into a file with the extension `*.yaml.tmpl` and commit it to your Quick Start repository.
Files with this extension get processed by eksctl before committing them to the user's gitops repository, while the rest get copied unmodified.

Regarding the folder structure inside the Quick Start repository, we recommend using a folder for each `namespace` and
one file per Kubernetes resource.

```
repository-name/
├── kube-system
│   ├── ingress-controller-deployment.yaml.tmpl
│   └── ingress-controller-rbac.yaml.tmpl
└── alerting
    ├── alerting-app-deployment.yaml
    ├── alerting-app-service.yaml.tmpl
    ├── monitoring-sidecar-deployment.yaml
    ├── monitoring-sidecar-service.yaml.tmpl
    ├── cluster-info-configmap.yaml.tmpl
    └── alerting-namespace.yaml

```

Note that some files have the extension `*.yaml` while others have `*.yaml.tmpl`. The last ones are the ones that can
contain template actions while the former are plain `yaml` files.

These files can now be committed and pushed to your Quick Start repository, for example `git@github.com:my-org/production-infra`.

```
cd repository-name/
git add .
git commit -m "Add component templates"
git push origin master
```

Now that the templates are in the remote repository, the Quick Start is ready to be used with `eksctl enable profile`:

```console
EKSCTL_EXPERIMENTAL=true eksctl enable profile --cluster team1 --region eu-west-1 --git-url git@github.com:my-org/team1-cluster --git-email alice@my-org.com git@github.com:my-org/production-infra 
```

In this example we provide `github.com:my-org/production-infra` as the Quick Start profile and 
`github.com:my-org/team1-cluster` as the gitops repository that is connected to the Flux instance in the cluster named
`cluster1`.


For a full example of a Quick Start profile, check out [App Dev][app-dev].


[flux]: https://docs.fluxcd.io/en/latest/
[go-templates]: https://golang.org/pkg/text/template/
[app-dev]: https://github.com/weaveworks/eks-quickstart-app-dev
