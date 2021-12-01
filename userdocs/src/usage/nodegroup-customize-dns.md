# Nodegroups with custom DNS

There are two ways of overwriting the DNS server IP address used for all the internal and external DNS lookups. This
is the equivalent of the `--cluster-dns` flag for the `kubelet`.

The first, is through the `clusterDNS` field. [Config files](../schema) accepts a `string` field called
`clusterDNS` with the IP address of the DNS server to use.
This will be passed to the `kubelet` that in turn will pass it to the pods through the `/etc/resolv.conf` file.

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: cluster-1
  region: eu-north-1

nodeGroups:
- name: ng-1
  clusterDNS: 169.254.20.10
```

Note that this configuration only accepts one IP address. To specify more than one address, use the
[`kubeletExtraConfig` parameter](../customizing-the-kubelet):

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: cluster-1
  region: eu-north-1

nodeGroups:
  - name: ng-1
    kubeletExtraConfig:
      clusterDNS: ["169.254.20.10","172.20.0.10"]
```
