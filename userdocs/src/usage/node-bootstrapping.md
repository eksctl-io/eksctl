# Node bootstrapping

##  AmazonLinux2023

AL2023 introduced a new node initialization process [nodeadm](https://awslabs.github.io/amazon-eks-ami/nodeadm/) that uses a YAML configuration schema, dropping the use of `/etc/eks/bootstrap.sh` script.

### Default settings

For self-managed nodes and EKS-managed nodes based on custom AMIs, `eksctl` creates a default, minimal, `NodeConfig` and automatically injects it into the nodegroups's launch template userdata. i.e.

```yaml
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary=//

--//
Content-Type: application/node.eks.aws

apiVersion: node.eks.aws/v1alpha1
kind: NodeConfig
spec:
  cluster:
    apiServerEndpoint: https://XXXX.us-west-2.eks.amazonaws.com
    certificateAuthority: XXXX
    cidr: 10.100.0.0/16
    name: my-cluster
  kubelet:
    config:
      clusterDNS:
      - 10.100.0.10
    flags:
    - --node-labels=alpha.eksctl.io/cluster-name=my-cluster,alpha.eksctl.io/nodegroup-name=my-nodegroup
    - --register-with-taints=special=true:NoSchedule

--//--
```

For EKS-managed nodes based on native AMIs, the default `NodeConfig` is being added by EKS MNG under the hood, appended directly to the EC2's userdata. Thus, in this scenario, `eksctl` does not need to include it within the launch template.

### Configuring the bootstrapping process

To set advanced properties of `NodeConfig`, or simply override the default values, eksctl allows you to specify a custom `NodeConfig` via `nodeGroup.overrideBootstrapCommand` or `managedNodeGroup.overrideBootstrapCommand`  e.g.

```yaml
managedNodeGroups:
  - name: mng-1
    amiFamily: AmazonLinux2023
    ami: ami-0253856dd7ab7dbc8
    overrideBootstrapCommand: |
      apiVersion: node.eks.aws/v1alpha1
      kind: NodeConfig
      spec:
        instance:
          localStorage:
            strategy: RAID0
```

This custom config will be prepended to the userdata by eksctl, and merged by `nodeadm` with the default config. Read more about `nodeadm`'s capability of merging multiple configuration objects [here](https://awslabs.github.io/amazon-eks-ami/nodeadm/doc/examples/#merging-multiple-configuration-objects).
