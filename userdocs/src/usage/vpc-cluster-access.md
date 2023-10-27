# Cluster Access

## Managing Access to the Kubernetes API Server Endpoints

The default creation of an EKS cluster exposes the Kubernetes API server publicly but not directly from within the
VPC subnets (public=true, private=false). Traffic destined for the API server from within the VPC must first exit the
VPC networks (but not Amazon's network) and then re-enter to reach the API server.

The Kubernetes API server endpoint access for a cluster can be configured for public and private access when creating
the cluster using the cluster config file. Example below:

```yaml
vpc:
  clusterEndpoints:
    publicAccess:  <true|false>
    privateAccess: <true|false>
```

There are some additional caveats when configuring Kubernetes API endpoint access:

1. EKS doesn't allow one to create or update a cluster without at least one private or public access being
   enabled.
1. EKS does allow creating a configuration that allows only private access to be enabled, but eksctl doesn't
   support it during cluster creation as it prevents eksctl from being able to join the worker nodes to the cluster.
1. Updating a cluster to have private only Kubernetes API endpoint access means that Kubernetes commands, by default,
   (e.g. `kubectl`) as well as `eksctl delete cluster`, `eksctl utils write-kubeconfig`, and possibly the command
   `eksctl utils update-kube-proxy` must be run within the cluster VPC.  This requires some changes to various AWS
   resources.  See:
   [EKS user guide](https://docs.aws.amazon.com/en_pv/eks/latest/userguide/cluster-endpoint)
   A user can provide `vpc.extraCIDRs` which will append additional CIDR ranges to the ControlPlaneSecurityGroup,
   allowing subnets outside the VPC to reach the kubernetes API endpoint. Similarly you can provide `vpc.extraIPv6CIDRs`
   to append IPv6 CIDR ranges as well.

The following is an example of how one could configure the Kubernetes API endpoint access using the `utils` sub-command:

```console
eksctl utils update-cluster-vpc-config --cluster=<clustername> --private-access=true --public-access=false
```

!!! warning
    `eksctl utils update-cluster-endpoints` has been deprecated in favour of `eksctl utils update-cluster-vpc-config`
    and will be removed soon.

To update the setting using a `ClusterConfig` file, use:

```console
eksctl utils update-cluster-vpc-config -f config.yaml --approve
```

Note that if you don't pass a flag, it will keep the current value. Once you are satisfied with the proposed changes,
add the `approve` flag to make the change to the running cluster.

## Restricting Access to the EKS Kubernetes Public API endpoint

The default creation of an EKS cluster exposes the Kubernetes API server publicly. To restrict access to the public API
endpoint to a set of CIDRs when creating a cluster, set the `publicAccessCIDRs` field:

```yaml
vpc:
  publicAccessCIDRs: ["1.1.1.1/32", "2.2.2.0/24"]
```

To update the restrictions on an existing cluster, use:

```console
eksctl utils update-cluster-vpc-config --cluster=<cluster> 1.1.1.1/32,2.2.2.0/24
```

!!! warning
    `eksctl utils set-public-access-cidrs` has been deprecated in favour of `eksctl utils update-cluster-vpc-config`
    and will be removed soon.

To update the restrictions using a `ClusterConfig` file, set the new CIDRs in `vpc.publicAccessCIDRs` and run:

```console
eksctl utils update-cluster-vpc-config -f config.yaml
```

!!! warning
    If setting `publicAccessCIDRs` and creating node-groups either `privateAccess` should be set to `true` or
    the nodes' IPs should be added to the `publicAccessCIDRs` list. Otherwise creation will fail with
    `context deadline exceeded` due to the nodes being unable to access the public endpoint and hence failing
    to join the cluster.

???+ note
    This feature only applies to the public endpoint. The
    [API server endpoint access configuration options](https://docs.aws.amazon.com/eks/latest/userguide/cluster-endpoint.html)
    won't change, and you will still have the option to disable the public endpoint so your cluster is not accessible from
    the internet. (Source: https://github.com/aws/containers-roadmap/issues/108#issuecomment-552766489)

    Implementation notes: https://github.com/aws/containers-roadmap/issues/108#issuecomment-552698875


To update both API server endpoint access and public access CIDRs for a cluster in a single command, run:

```console
eksctl utils update-cluster-vpc-config --cluster=<cluster> --public-access=true --private-access=true --public-access-cidrs=1.1.1.1/32,2.2.2.0/24
```

To update the setting using a config file:

```yaml
vpc:
  clusterEndpoints:
    publicAccess:  <true|false>
    privateAccess: <true|false>
  publicAccessCIDRs: ["1.1.1.1/32"]
```

```console
eksctl utils update-cluster-vpc-config --cluster=<cluster> -f config.yaml
```
