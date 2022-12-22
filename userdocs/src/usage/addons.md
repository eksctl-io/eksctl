# Addons

EKS Add-Ons is a new feature that lets you enable and manage Kubernetes operational
software for your AWS EKS clusters. At launch, EKS add-ons supports controlling the launch and version of the AWS VPC
CNI plugin through the EKS API

## Creating addons

You can specify what addons you want and what policies (if required) to attach to them in your config file:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig
metadata:
  name: example-cluster
  region: us-west-2
  version: "1.20"

iam:
  withOIDC: true

addons:
- name: vpc-cni
  version: 1.7.5 # optional
  attachPolicyARNs: #optional
  - arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy
  serviceAccountRoleARN: arn:aws:iam::aws:policy/AmazonEKSCNIAccess # optional
  tags: # optional
    team: eks
  attachPolicy: # optional
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

!!!note
    In order to attach policies to addons your cluster must have `OIDC` enabled. If it's not enabled we ignore any policies
    attached.


You can then either have these addons created during the cluster creation process:
```console
eksctl create cluster -f config.yaml
```

Or you can create after cluster creation using the config file or CLI flags:

```console
eksctl create addon -f config.yaml
```

```console
eksctl create addon --name vpc-cni --version 1.7.5 --service-account-role-arn=<role-arn>
```

During addon creation, if a self-managed version of the addon already exists on the cluster, you can choose how potential `configMap` conflicts shall be resolved by setting `resolveConflicts` option via the config file. e.g.,

```yaml
addons:
- name: vpc-cni
  attachPolicyARNs:
    - arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy
  resolveConflicts: overwrite
```

For addon create, the `resolveConflicts` field supports two distinct values.

- `overwrite` - EKS overwrites any config changes back to EKS default values
- `none` - EKS doesn't change the value. The create might fail.

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

Setting the version of the addon is optional. If the `version` field is empty in the request sent by `eksctl`, the EKS API will set it to the default version for that specific addon. More information about which version is the default version for specific addons can be found in the AWS documentation about EKS. Note that the default version might not necessarily be the latest version available.

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
These previously discovered options can be configured by passing the `ConfigurationValues` field to the configuration file, and then used with either addon create or update commands. For the moment, all the options the user wants to configure will be passed together in the format of a JSON string, e.g.

```yaml
addons:
- name: coredns
  version: latest
  configurationValues: "{\"replicaCount\":3}"
  resolveConflics: overwrite
```

!!!note
    Bear in mind that when addon configuration values are being modified, configuration conflics will arise.
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
eksctl update addon --name vpc-cni --version 1.8.0 --service-account-role-arn=<new-role>
```

Similarly to addon creation, When updating an addon, you have full control over the config changes that you may have previously applied on that add-on's `configMap`. Specifically, you can preserve, or overwrite them. This optional functionality is available via the same config file field `resolveConflicts`. e.g.,


```yaml
addons:
- name: vpc-cni
  attachPolicyARNs:
    - arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy
  resolveConflicts: preserve
```

For addon update, the `resolveConflicts` field accepts three distinct values.

- `preserve` - EKS preserves the value. If you choose this option, we recommend that you test any field and value changes on a non-production cluster before updating the add-on on your production cluster.
- `overwrite` - EKS overwrites any config changes back to EKS default values
- `none` - EKS doesn't change the value. The update might fail.

## Deleting addons
You can delete an addon by running:
```console
eksctl delete addon --cluster <cluster-name> --name <addon-name
```
This will delete the addon and any IAM roles associated to it.

When you delete your cluster all IAM roles associated to addons are also deleted.
