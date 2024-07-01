# Addons

EKS Add-Ons is a new feature that lets you enable and manage Kubernetes operational
software for your AWS EKS clusters. At launch, EKS add-ons supports controlling the launch and version of the AWS VPC
CNI plugin through the EKS API

## Creating addons (and providing IAM permissions via IRSA)

!!! tip "New for 2024"
    eksctl now supports creating clusters without any default networking addons: [Cluster creation flexibility for default networking addons](#cluster-creation-flexibility-for-default-networking-addons).

!!! warning "New for 2024"
    eksctl now installs default addons as EKS addons instead of self-managed addons. Read more about its implications in [Cluster creation flexibility for default networking addons](#cluster-creation-flexibility-for-default-networking-addons).

!!! tip "New for 2024"
    EKS Add-ons now support receiving IAM permissions, required to connect with AWS services outside of cluster, via [EKS Pod Identity Associations](/usage/pod-identity-associations/#eks-add-ons-support-for-pod-identity-associations)

In your config file, you can specify the addons you want and (if required) the role or policies to attach to them:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig
metadata:
  name: example-cluster
  region: us-west-2

iam:
  withOIDC: true

addons:
- name: vpc-cni
  # all below properties are optional
  version: 1.7.5
  tags:
    team: eks
  # you can specify at most one of:
  attachPolicyARNs:
  - arn:aws:iam::account:policy/AmazonEKS_CNI_Policy
  # or
  serviceAccountRoleARN: arn:aws:iam::account:role/AmazonEKSCNIAccess
  # or
  attachPolicy:
    Statement:
    - Effect: Allow
      Action:
      - ec2:AssignPrivateIpAddresses
      - ec2:AttachNetworkInterface
      - ec2:CreateNetworkInterface
      - ec2:DeleteNetworkInterface
      - ec2:DescribeInstances
      - ec2:DescribeTags
      - ec2:DescribeNetworkInterfaces
      - ec2:DescribeInstanceTypes
      - ec2:DetachNetworkInterface
      - ec2:ModifyNetworkInterfaceAttribute
      - ec2:UnassignPrivateIpAddresses
      Resource: '*'
```

You can specify at most one of `attachPolicy`, `attachPolicyARNs` and `serviceAccountRoleARN`.

If none of these are specified, the addon will be created with a role that has all recommended policies attached.

???+ note
    In order to attach policies to addons your cluster must have `OIDC` enabled. If it's not enabled we ignore any policies
    attached.


You can then either have these addons created during the cluster creation process:
```console
eksctl create cluster -f config.yaml
```

Or create the addons explicitly after cluster creation using the config file or CLI flags:

```console
eksctl create addon -f config.yaml
```

```console
eksctl create addon --name vpc-cni --version 1.7.5 --service-account-role-arn <role-arn>
```

During addon creation, if a self-managed version of the addon already exists on the cluster, you can choose how potential `configMap` conflicts shall be resolved by setting `resolveConflicts` option via the config file, e.g.

```yaml
addons:
- name: vpc-cni
  attachPolicyARNs:
    - arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy
  resolveConflicts: overwrite
```

For addon create, the `resolveConflicts` field supports three distinct values:

- `none` - EKS doesn't change the value. The create might fail.
- `overwrite` - EKS overwrites any config changes back to EKS default values.
- `preserve` - EKS doesn't change the value. The create might fail. (Similarly to `none`, but different from [`preserve` in updating addons](#updating-addons))

## Listing enabled addons

You can see what addons are enabled in your cluster by running:
```console
eksctl get addons --cluster <cluster-name>
```

or

```console
eksctl get addons -f config.yaml
```

## Setting the addon's version

Setting the version of the addon is optional. If the `version` field is left empty `eksctl` will resolve the default version for the addon. More information about which version is the default version for specific addons can be found in the AWS documentation about EKS. Note that the default version might not necessarily be the latest version available.

The addon version can be set to `latest`. Alternatively, the version can be set with the EKS build tag specified, such as `v1.7.5-eksbuild.1` or `v1.7.5-eksbuild.2`. It can also be set to the release version of the addon, such as `v1.7.5` or `1.7.5`, and the `eksbuild` suffix tag will be discovered and set for you.

See the section below on how to discover available addons and their versions.

## Discovering addons
You can discover what addons are available to install on your cluster by running:
```console
eksctl utils describe-addon-versions --cluster <cluster-name>
```

This will discover your cluster's kubernetes version and filter on that. Alternatively if you want to see what
addons are available for a particular kubernetes version you can run:
```console
eksctl utils describe-addon-versions --kubernetes-version <version>
```

You can also discover addons by filtering on their `type`, `owner` and/or `publisher`.
For e.g., to see addons for a particular owner and type you can run:
```console
eksctl utils describe-addon-versions --kubernetes-version 1.22 --types "infra-management, policy-management" --owners "aws-marketplace"
```
The `types`, `owners` and `publishers` flags are optional and can be specified together or individually to filter the results.

## Discovering the configuration schema for addons
After discovering the addon and version, you can view the customization options by fetching its JSON configuration schema.

```console
eksctl utils describe-addon-configuration --name vpc-cni --version v1.12.0-eksbuild.1
```

This returns a JSON schema of the various options available for this addon.

## Working with configuration values
`ConfigurationValues` can be provided in the configuration file during the creation or update of addons. Only JSON and YAML formats are supported.

For eg.,

```yaml
addons:
- name: coredns
  configurationValues: |-
    replicaCount: 2
```

```yaml
addons:
- name: coredns
  version: latest
  configurationValues: "{\"replicaCount\":3}"
  resolveConflicts: overwrite
```

???+ note
    Bear in mind that when addon configuration values are being modified, configuration conflicts will arise.

    Thus, we need to specify how to deal with those by setting the `resolveConflicts` field accordingly.
    As in this scenario we want to modify these values, we'd set `resolveConflicts: overwrite`.

Additionally, the get command will now also retrieve `ConfigurationValues` for the addon. e.g.

```console
eksctl get addon --cluster my-cluster --output yaml
```
```yaml
- ConfigurationValues: '{"replicaCount":3}'
  IAMRole: ""
  Issues: null
  Name: coredns
  NewerVersion: ""
  Status: ACTIVE
  Version: v1.8.7-eksbuild.3
```

## Updating addons
You can update your addons to newer versions and change what policies are attached by running:
```console
eksctl update addon -f config.yaml
```

```console
eksctl update addon --name vpc-cni --version 1.8.0 --service-account-role-arn <new-role>
```

Similarly to addon creation, When updating an addon, you have full control over the config changes that you may have previously applied on that add-on's `configMap`. Specifically, you can preserve, or overwrite them. This optional functionality is available via the same config file field `resolveConflicts`. e.g.,


```yaml
addons:
- name: vpc-cni
  attachPolicyARNs:
    - arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy
  resolveConflicts: preserve
```

For addon update, the `resolveConflicts` field accepts three distinct values:

- `none` - EKS doesn't change the value. The update might fail.
- `overwrite` - EKS overwrites any config changes back to EKS default values.
- `preserve` - EKS preserves the value. If you choose this option, we recommend that you test any field and value changes on a non-production cluster before updating the add-on on your production cluster.

## Deleting addons
You can delete an addon by running:
```console
eksctl delete addon --cluster <cluster-name> --name <addon-name>
```
This will delete the addon and any IAM roles associated to it.

When you delete your cluster all IAM roles associated to addons are also deleted.

## Cluster creation flexibility for default networking addons

When a cluster is created, EKS automatically installs VPC CNI, CoreDNS and kube-proxy as self-managed addons.
To disable this behavior in order to use other CNI plugins like Cilium and Calico, eksctl now supports creating a cluster
without any default networking addons. To create such a cluster, set `addonsConfig.disableDefaultAddons`, as in:

```yaml
addonsConfig:
  disableDefaultAddons: true
```

```shell
$ eksctl create cluster -f cluster.yaml
```

To create a cluster with only CoreDNS and kube-proxy and not VPC CNI, specify the addons explicitly in `addons`
and set `addonsConfig.disableDefaultAddons`, as in:

```yaml
addonsConfig:
  disableDefaultAddons: true
addons:
  - name: kube-proxy
  - name: coredns
```

```shell
$ eksctl create cluster -f cluster.yaml
```

As part of this change, eksctl now installs default addons as EKS addons instead of self-managed addons during cluster creation
if `addonsConfig.disableDefaultAddons` is not explicitly set to true. As such, `eksctl utils update-*` commands can no
longer be used for updating addons for clusters created with eksctl v0.184.0 and above:

- `eksctl utils update-aws-node`
- `eksctl utils update-coredns`
- `eksctl utils update-kube-proxy`

Instead, `eksctl update addon` should be used now.

To learn more, see [EKS documentation][eksdocs].


[eksdocs]: https://aws.amazon.com/about-aws/whats-new/2024/06/amazon-eks-cluster-creation-flexibility-networking-add-ons/
