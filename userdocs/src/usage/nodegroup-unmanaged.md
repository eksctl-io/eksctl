# Unmanaged nodegroups

In `eksctl`, setting `--managed=false` or using the `nodeGroups` field creates an unmanaged nodegroup. Bear in mind that
unmanaged nodegroups do not appear in the EKS console, which as a general rule only knows about EKS-managed nodegroups.

You should be upgrading nodegroups only after you ran `eksctl upgrade cluster`.
(See [Upgrading clusters](/usage/cluster-upgrade).)

If you have a simple cluster with just an initial nodegroup (i.e. created with
`eksctl create cluster`), the process is very simple:

1. Get the name of old nodegroup:

    ```shell
    eksctl get nodegroups --cluster=<clusterName> --region=<region>
    ```

    ???+ note
        You should see only one nodegroup here, if you see more - read the next section.

2. Create a new nodegroup:

    ```shell
    eksctl create nodegroup --cluster=<clusterName> --region=<region> --name=<newNodeGroupName> --managed=false
    ```

3. Delete the old nodegroup:

    ```shell
    eksctl delete nodegroup --cluster=<clusterName> --region=<region> --name=<oldNodeGroupName>
    ```

    ???+ note
        This will drain all pods from that nodegroup before the instances are deleted. In some scenarios, Pod Disruption Budget (PDB) policies can prevent pods to be evicted. To delete the nodegroup regardless of PDB, one should use the `--disable-eviction` flag, will bypass checking PDB policies.

## Updating multiple nodegroups

If you have multiple nodegroups, it's your responsibility to track how each one was configured.
You can do this by using config files, but if you haven't used it already, you will need to inspect
your cluster to find out how each nodegroup was configured.

In general terms, you are looking to:

- review which nodegroups you have and which ones can be deleted or must be replaced for the new version
- note down configuration of each nodegroup, consider using config file to ease upgrades next time

### Updating with config file

If you are using config file, you will need to do the following.

Edit config file to add new nodegroups, and remove old nodegroups.
If you just want to upgrade nodegroups and keep the same configuration,
you can just change nodegroup names, e.g. append `-v2` to the name.

To create all of new nodegroups defined in the config file, run:

```
eksctl create nodegroup --config-file=<path>
```

Once you have new nodegroups in place, you can delete old ones:

```
eksctl delete nodegroup --config-file=<path> --only-missing
```

???+ note
    First run is in plan mode, if you are happy with the proposed changes, re-run with `--approve`.

## Updating default add-ons

There are 3 default add-ons that get included in each EKS cluster, the process for updating each of them is different, hence
there are 3 distinct commands that you will need to run.

???+ info
    All of the following commands accept `--config-file`.

???+ note
    By default each of these commands runs in plan mode, if you are happy with the proposed changes, re-run with `--approve`.

To update `kube-proxy`, run:

```
eksctl utils update-kube-proxy --cluster=<clusterName>
```

To update `aws-node`, run:

```
eksctl utils update-aws-node --cluster=<clusterName>
```

To update `coredns`, run:

```
eksctl utils update-coredns --cluster=<clusterName>
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
