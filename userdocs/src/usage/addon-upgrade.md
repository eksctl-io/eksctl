# Default add-on updates

!!! warning "New for 2024"
    eksctl now installs default addons as EKS addons instead of self-managed addons. Read more about its implications in [Cluster creation flexibility for default networking addons](#cluster-creation-flexibility-for-default-networking-addons).

!!! warning "New for 2024"
    For updating addons, `eksctl utils update-*` cannot be used for clusters created with eksctl v0.184.0 and above.
    This guide is only valid for clusters created before this change.

There are 3 default add-ons that get included in each EKS cluster:
- `kube-proxy`
- `aws-node`
- `coredns`

???+ info
    For official EKS addons that are created manually through `eksctl create addons` or upon cluster creation, the way to manage them is
    through `eksctl create/get/update/delete addon`. In such cases, please refer to the docs about [EKS Add-Ons](https://eksctl.io/usage/addons/).

The process for updating each of them is different, hence there are 3 distinct commands that you will need to run.

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

## What "latest" means for each command

It is a common source of confusion to compare the version applied by these commands against
`eksctl utils describe-addon-versions`. That command lists **EKS Managed Add-on versions**
(which follow their own release cadence and versioning, e.g. `v1.18.2-eksbuild.1`), whereas these
commands update **self-managed** default add-ons that run directly in the cluster. The two are not
the same and should not be expected to match.

- `eksctl utils update-kube-proxy` resolves the newest `kube-proxy` version for your cluster's
  control-plane version from the EKS API and applies it.
- `eksctl utils update-aws-node` and `eksctl utils update-coredns` apply the manifest **bundled
  with the eksctl release you are running** (a "known good" version curated by eksctl). They do not
  resolve the newest version from the EKS API, so the applied version can differ from the latest
  one listed by `describe-addon-versions`.

???+ tip
    If you want the newest version shipped by EKS for `vpc-cni`, `coredns` or `kube-proxy`, install
    them as [managed add-ons](https://eksctl.io/usage/addons/) instead, e.g. for VPC CNI run:

    ```
    eksctl create addon --name vpc-cni --cluster <clusterName>
    ```

    and keep it up to date with `eksctl update addon --name vpc-cni --cluster <clusterName> --version latest`.
