---
title: "Cluster Upgrades"
weight: 30
url: usage/cluster-upgrades
---

## Cluster upgrades

An _`eksctl`-managed_ cluster can be upgraded in 3 easy steps:

1. update control plane version with `eksctl update cluster`
2. update default add-ons:
    - `kube-proxy`
    - `aws-node`
    - `coredns`
3. replace each of the nodegroups by creating a new one and deleting the old one

Please make sure to read this section in full before you proceed.

> NOTE: Kubernetes supports version drift of up-to two minor versions during upgrade
process. So nodes can be up to two minor versions ahead or behind the control plane
version. You can only upgrade the control plane one minor version at a time, but
nodes can be upgraded more than one minor version at a time, provided the nodes stay
within two minor versions of the control plane.

### Updating control plane version

Control plane version updates must be done for one minor version at a time.

To update control plane to the next available version run:

```
eksctl update cluster --name=<clusterName>
```

This command will not apply any changes right away, you will need to re-run it with
`--approve` to apply the changes.

### Updating nodegroups

You should update nodegroups only after you ran `eksctl update cluster`.

If you have a simple cluster with just an initial nodegroup (i.e. created with
`eksctl create cluster`), the process is very simple.

Get the name of old nodegroup:

```
eksctl get nodegroups --cluster=<clusterName>
```

> NOTE: you should see only one nodegroup here, if you see more - read the next section

Create new nodegroup:

```
eksctl create nodegroup --cluster=<clusterName>
```

Delete old nodegroup:

```
eksctl delete nodegroup --cluster=<clusterName> --name=<oldNodeGroupName>
```

> NOTE: this will drain all pods from that nodegroup before the instances are deleted.

#### Updating multiple nodegroups

If you have multiple nodegroups, it's your responsibility to track how each one was configured.
You can do this by using config files, but if you haven't used it already, you will need to inspect
your cluster to find out how each nodegroup was configured.

In general terms, you are looking to:

- review nodegroups you have and which ones can be deleted or must be replaced for the new version
- note down configuration of each nodegroup, consider using config file to ease upgrades next time

To create a new nodegroup:

```
eksctl create nodegroup --cluster=<clusterName> --name=<newNodeGroupName>
```

To delete old nodegroup:

```
eksctl delete nodegroup --cluster=<clusterName> --name=<oldNodeGroupName>
```

#### Updating multiple nodegroups with config file

If you are using config file, you will need to do the following.

Edit config file to add new nodegroups, and remove old nodegroups.
If you just want to update nodegroups and keep the same configuration,
you can just change nodegroup names, e.g. append `-v2` to the name.

To create all of new nodegroups defined in the config file, run:

```
eksctl create nodegroup --config-file=<path>
```

Once you have new nodegroups in place, you can delete old ones:

```
eksctl delete nodegroup --config-file=<path> --only-missing
```

> NOTE: first run is in plan mode, if you are happy with the proposed
> changes, re-run with `--approve`.

### Updating default add-ons

There are 3 default add-ons that get included in each EKS cluster, the process for updating each of them is different, hence
there are 3 distinct commands that you will need to run.

> NOTE: all of the following commands accept `--config-file`.

> NOTE: by default each of these commands runs in plan mode,
> if you are happy with the proposed changes, re-run with `--approve`.

To update `kube-proxy`, run:

```
eksctl utils update-kube-proxy
```

To update `aws-node`, run:

```
eksctl utils update-aws-node
```

To update `coredns`, run:

```
eksctl utils update-coredns
```

Once upgraded, be sure to run `kubectl get pods -n kube-system` and check if all addon pods are in ready state, you should see
something like this:

```
NAME                       READY   STATUS    RESTARTS   AGE
aws-node-g5ghn             1/1     Running   0          2m
aws-node-zfc9s             1/1     Running   0          2m
coredns-7bcbfc4774-g6gg8   1/1     Running   0          1m
coredns-7bcbfc4774-hftng   1/1     Running   0          1m
kube-proxy-djkp7           1/1     Running   0          3m
kube-proxy-mpdsp           1/1     Running   0          3m
```
