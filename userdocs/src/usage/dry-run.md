# Dry Run

The dry-run feature allows you to inspect and change the instances matched by the instance selector before proceeding
to creating a nodegroup.

When `eksctl create cluster` is called with the instance selector options and `--dry-run`, eksctl will output a
ClusterConfig file containing a nodegroup representing the CLI options and the instance types set to the instances
matched by the instance selector resource criteria.

```shell
$ eksctl create cluster --name development --dry-run


apiVersion: eksctl.io/v1alpha5
cloudWatch:
  clusterLogging: {}
iam:
  vpcResourceControllerPolicy: true
  withOIDC: false
kind: ClusterConfig
managedNodeGroups:
- amiFamily: AmazonLinux2
  desiredCapacity: 2
  disableIMDSv1: true
  disablePodIMDS: false
  iam:
    withAddonPolicies:
      albIngress: false
      appMesh: false
      appMeshPreview: false
      autoScaler: false
      certManager: false
      cloudWatch: false
      ebs: false
      efs: false
      externalDNS: false
      fsx: false
      imageBuilder: false
      xRay: false
  instanceSelector: {}
  instanceType: m5.large
  labels:
    alpha.eksctl.io/cluster-name: development
    alpha.eksctl.io/nodegroup-name: ng-4aba8a47
  maxSize: 2
  minSize: 2
  name: ng-4aba8a47
  privateNetworking: false
  securityGroups:
    withLocal: null
    withShared: null
  ssh:
    allow: false
    enableSsm: false
    publicKeyPath: ""
  tags:
    alpha.eksctl.io/nodegroup-name: ng-4aba8a47
    alpha.eksctl.io/nodegroup-type: managed
  volumeIOPS: 3000
  volumeSize: 80
  volumeThroughput: 125
  volumeType: gp3
metadata:
  name: development
  region: us-west-2
  version: "1.24"
privateCluster:
  enabled: false
vpc:
  autoAllocateIPv6: false
  cidr: 192.168.0.0/16
  clusterEndpoints:
    privateAccess: false
    publicAccess: true
  manageSharedNodeSecurityGroupRules: true
  nat:
    gateway: Single
```

The generated ClusterConfig can then be passed to `eksctl create cluster`:

```console
$ eksctl create cluster -f generated-cluster.yaml
```

When a ClusterConfig file is passed with `--dry-run`, eksctl will output a ClusterConfig file containing the values set in the file.

???+ note
    There are certain one-off options that cannot be represented in the ClusterConfig file, e.g., `--install-vpc-controllers`. It is expected that `eksctl create cluster --<options...> --dry-run` > config.yaml followed by `eksctl create cluster -f config.yaml` would be equivalent to running the first command without `--dry-run`. eksctl therefore disallows passing options that cannot be represented in the config file when `--dry-run` is passed. If you need to pass an AWS profile, set the `AWS_PROFILE` environment variable, instead of passing the `--profile` CLI option.
