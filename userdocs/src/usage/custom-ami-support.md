# Custom AMI support

## Setting the node AMI ID

The `--node-ami` flag enables a number of advanced use cases such as using a custom AMI or querying AWS in realtime to determine which AMI to use.
The flag can be used for both non-GPU and GPU images.

The flag can take the AMI image id for an image to explicitly use. It also can take the following 'special' keywords:

| Keyword   | Description                                                                                                         |
| --------- | ------------------------------------------------------------------------------------------------------------------- |
| auto      | Indicates that the AMI to use for the nodes should be found by querying AWS EC2. This relates to the auto resolver. |
| auto-ssm  | Indicates that the AMI to use for the nodes should be found by querying AWS SSM Parameter Store.                    |


???+ note
    When setting `--node-ami` to an ID string, `eksctl` will assume that a custom AMI has been requested.
    For AmazonLinux2 and Ubuntu nodes, both EKS managed and self-managed, this will mean that `overrideBootstrapCommand` is required.

CLI flag examples:
```sh
eksctl create cluster --node-ami=auto

# with a custom ami id
eksctl create cluster --node-ami=ami-custom1234
```

Config file example:
```yaml
nodeGroups:
  - name: ng1
    instanceType: p2.xlarge
    amiFamily: AmazonLinux2
    ami: auto
  - name: ng2
    instanceType: m5.large
    amiFamily: AmazonLinux2
    ami: ami-custom1234
managedNodeGroups:
  - name: m-ng-2
    amiFamily: AmazonLinux2
    ami: ami-custom1234
    instanceType: m5.large
    overrideBootstrapCommand: |
      #!/bin/bash
      /etc/eks/bootstrap.sh <cluster-name>
```

The `--node-ami` flag can also be used with `eksctl create nodegroup`.

## Setting the node AMI Family

The `--node-ami-family` can take following keywords:

| Keyword                        |                                          Description                                         |
|--------------------------------|:--------------------------------------------------------------------------------------------:|
| AmazonLinux2                   | Indicates that the EKS AMI image based on Amazon Linux 2 should be used (default).           |
| Ubuntu2004                     | Indicates that the EKS AMI image based on Ubuntu 20.04 LTS (Focal) should be used.           |
| Ubuntu1804                     | Indicates that the EKS AMI image based on Ubuntu 18.04 LTS (Bionic) should be used.          |
| Bottlerocket                   | Indicates that the EKS AMI image based on Bottlerocket should be used.                       |
| WindowsServer2019FullContainer | Indicates that the EKS AMI image based on Windows Server 2019 Full Container should be used. |
| WindowsServer2019CoreContainer | Indicates that the EKS AMI image based on Windows Server 2019 Core Container should be used. |
| WindowsServer2022FullContainer | Indicates that the EKS AMI image based on Windows Server 2022 Full Container should be used. |
| WindowsServer2022CoreContainer | Indicates that the EKS AMI image based on Windows Server 2022 Core Container should be used. |

CLI flag example:
```sh
eksctl create cluster --node-ami-family=AmazonLinux2
```

Config file example:
```yaml
nodeGroups:
  - name: ng1
    instanceType: m5.large
    amiFamily: AmazonLinux2
managedNodeGroups:
  - name: m-ng-2
    instanceType: m5.large
    amiFamily: Ubuntu2004
```

The `--node-ami-family` flag can also be used with `eksctl create nodegroup`. `eksctl` requires AMI Family to be explicitly set via config file or via `--node-ami-family` CLI flag, whenever working with a custom AMI.

???+ note
    At the moment, EKS managed nodegroups only support the following AMI Families when working with custom AMIs: `AmazonLinux2`, `Ubuntu1804` and `Ubuntu2004`

## Windows custom AMI support
Only self-managed Windows nodegroups can specify a custom AMI. `amiFamily` should be set to a valid Windows AMI family.

The following PowerShell variables will be available to the bootstrap script:

```
$EKSBootstrapScriptFile
$EKSClusterName
$APIServerEndpoint
$Base64ClusterCA
$ServiceCIDR
$KubeletExtraArgs
$KubeletExtraArgsMap: A hashtable containing arguments for the kubelet, e.g., @{ 'node-labels' = ''; 'register-with-taints' = ''; 'max-pods' = '10'}
$DNSClusterIP
$ContainerRuntime
```

Config file example:
```yaml
nodeGroups:
  - name: custom-windows
    amiFamily: WindowsServer2022FullContainer
    ami: ami-01579b74557facaf7
    overrideBootstrapCommand: |
      & $EKSBootstrapScriptFile -EKSClusterName "$EKSClusterName" -APIServerEndpoint "$APIServerEndpoint" -Base64ClusterCA "$Base64ClusterCA" -ContainerRuntime "containerd" -KubeletExtraArgs "$KubeletExtraArgs" 3>&1 4>&1 5>&1 6>&1
```

## Bottlerocket custom AMI support

For Bottlerocket nodes, the `overrideBootstrapCommand` is not supported. Instead, to designate their own bootstrap container, one should use the `bottlerocket` field as part of the configuration file. E.g.

```yaml
  nodeGroups:
  - name: bottlerocket-ng
    ami: ami-custom1234
    amiFamily: Bottlerocket
    bottlerocket:
      enableAdminContainer: true
      settings:
        bootstrap-containers
          bootstrap
            source: <MY-CONTAINER-URI>
```
