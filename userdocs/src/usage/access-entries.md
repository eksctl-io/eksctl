# EKS Access Entries

## Introduction

AWS EKS has introduced a new set of controls, called access entries, for managing access of IAM principals to Kubernetes clusters. `eksctl` has fully integrated with this feature, allowing users to directly associate access policies to certain IAM principals, while doing work behind the scenes for others. More details in the [upcoming section](access-entries.md#how-does-this-affect-different-resources).

EKS predefines several managed access policies that mirror the default Kubernetes user facing roles. Predefined access policies can also include policies with permissions required by other AWS services such as Amazon EMR to run workloads on EKS clusters. See a list of predefined access policies as-well as a detailed description for each of those [here]().

???+ note
    For now, users can only use predefined EKS access policies. For more advanced requirements, one can continue to use `iamIdentityMappings`.
    Bear in mind that the permissions associated with a predefined access policy are subject to change over time. EKS will periodically backfill policies to match upstream permissions.

## How to enable the access entries API?

`eksctl` has added a new `accessConfig.authenticationMode` field, which dictates how cluster access management is achieved, and can be set to one of the following three values: 
  - `CONFIG_MAP` - default in EKS API - only `aws-auth` ConfigMap will be used
  - `API` - only access entries API will be used
  - `API_AND_CONFIG_MAP` - default in `eksctl` - both `aws-auth` ConfigMap and access entries API can be used

e.g.

```yaml
accessConfig:
  authenticationMode: <> 
```

When creating a new cluster with access entries, using `eksctl`, if `authenticationMode` is not provided by the user, it is automatically set to `API_AND_CONFIG_MAP`. Thus, the access entries API will be enabled by default. If instead you want to use access entries on an already existing, non-eksctl created, cluster, where `CONFIG_MAP` option is used, the user will need to first set `authenticationMode` to `API_AND_CONFIG_MAP`. For that, `eksctl` has introduced a new command for updating the cluster authentication mode, which works both with CLI flags e.g.

```
eksctl utils update-authentication-mode --cluster my-cluster --authentication-mode API_AND_CONFIG_MAP
```

and by providing a config file e.g.

```
eksctl utils update-authentication-mode -f config.yaml
```

## How does this affect different resources?

### IAM Entities

Cluster management access for these type of resources falls under user's control. `eksctl` has added a new `accessConfig.accessEntries` field that maps one-to-one to the [Access Entries EKS API](). .e.g.

```yaml
accessConfig:
  authenticationMode: API_AND_CONFIG_MAP
  accessEntries:
    - principalARN: arn:aws:iam::111122223333:user/my-user-name
      kubernetesGroups: # optional Kubernetes groups
        - group1 # groups can used to give permissions via RBAC
        - group2

    - principalARN: arn:aws:iam::111122223333:role/role-name-1
      accessPolicies: # optional access polices
        - policyARN: arn:aws:eks::aws:cluster-access-policy/AmazonEKSViewPolicy
          accessScope:
            type: namespace
            namespaces:
              - default
              - my-namespace
              - dev-*

    - principalARN: arn:aws:iam::111122223333:role/admin-role
      accessPolicies: # optional access polices
        - policyARN: arn:aws:eks::aws:cluster-access-policy/AmazonEKSClusterAdminPolicy
          accessScope:
            type: cluster
```

In addition to associating EKS policies, one can also specify the Kubernetes groups to which an IAM entity belongs, thus granting permissions via RBAC.   

### Managed nodegroups and Fargate

The integration with access entries for these resources will be achieved behind the scenes, by the EKS API. Newly created managed node groups and Fargate pods will create API access entries, rather than using pre-loaded RBAC resources. Existing node groups and Fargate pods will not be changed, and continue to rely on the entries in the aws-auth config map.

### Self-managed nodegroups

For authorizing self-managed nodegroups, `eksctl` will create a unique access entry for each nodegroup with the principal ARN set to the node role ARN and type set to either `EC2_LINUX` or `EC2_WINDOWS` depending on nodegroup amiFamily.


## Managing access entries

### Create access entries

This can be done in two different ways. Either during cluster creation, specifying the desired access entries as part of the config file and running:

```
eksctl create cluster -f config.yaml
```

OR post cluster creation, by running:

```
eksctl create accessentry -f config.yaml
```

An example config file for creating access entries can be found [here](https://github.com/weaveworks/eksctl/blob/main/examples/40-access-entries.yaml).

### Fetch access entries

The user can retieve all access entries associated with a certain cluster by running one of the following:

```
eksctl get accessentry -f config.yaml
```

OR

```
eksctl get accessentry --cluster my-cluster
```

Alternatively, to retrieve only the access entry corresponding to a certain IAM entity one shall use the `--principal-arn` flag. e.g.

```
eksctl get accessentry --cluster my-cluster --principal-arn arn:aws:iam::111122223333:user/admin
```

### Delete access entries

To delete a single access entry at a time use:

```
eksctl delete accessentry --cluster my-cluster --principal-arn arn:aws:iam::111122223333:user/admin
```

To delete multiple access entries, use the `--config-file` flag and specify all the `principalARN's` corresponding with the access entries, under the top-level `accessEntry` field, e.g.

```yaml
...
accessEntry:
  - principalARN: arn:aws:iam::111122223333:user/my-user-name
  - principalARN: arn:aws:iam::111122223333:role/role-name-1
  - principalARN: arn:aws:iam::111122223333:role/admin-role
```

```
eksctl delete accessentry -f config.yaml
```


## Disabling cluster creator admin permissions
`eksctl` has added a new field `accessConfig.bootstrapClusterCreatorAdminPermissions: boolean` that, when set to false, disables granting cluster-admin permissions to the IAM identity creating the cluster. i.e.

add the option to the config file:

```yaml
accessConfig:
  bootstrapClusterCreatorAdminPermissions: false
```

and run:

```
eksctl create cluster -f config.yaml
```