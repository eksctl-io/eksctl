# EKS Pod Identity Associations

## Introduction

AWS EKS has introduced a new enhanced mechanism called Pod Identity Association for cluster administrators to configure Kubernetes applications to receive IAM permissions required to connect with AWS services outside of the cluster. Pod Identity Association leverages IRSA, however, it makes it configurable directly through EKS API, eliminating the need for using IAM API altogether.

As a result, IAM roles no longer need to reference an [OIDC provider](/usage/iamserviceaccounts/#how-it-works) and hence won't be tied to a single cluster anymore. This means, IAM roles can now be used across multiple EKS clusters without the need to update the role trust policy each time a new cluster is created. This in turn, eliminates the need for role duplication and simplifies the process of automating IRSA altogether.

## Prerequisites

Behind the scenes, the implementation of pod identity associations is running an agent as a daemonset on the worker nodes. To run the pre-requisite agent on the cluster, EKS provides a new add-on called EKS Pod Identity Agent. Therefore, creating pod identity associations (in general, and with `eksctl`) requires the `eks-pod-identity-agent` addon pre-installed on the cluster. This addon can be [created using `eksctl`](/usage/addons/#creating-addons) in the same fashion any other supported addon is, e.g.

```
eksctl create addon --cluster my-cluster --name eks-pod-identity-agent
```

Additionally, if using a pre-existing IAM role when creating a pod identity association, you must configure the role to trust the newly introduced EKS service principal (`pods.eks.amazonaws.com`). An example IAM trust policy can be found below:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Service": "pods.eks.amazonaws.com"
            },
            "Action": [
                "sts:AssumeRole",
                "sts:TagSession"
            ]
        }
    ]
}
```

If instead you do not provide the ARN of an existing role to the create command, `eksctl` will create one behind the scenes and configure the above trust policy.

## Creating Pod Identity Associations

For manipulating pod identity associations, `eksctl` has added a new field under `iam.podIdentityAssociations`, e.g.

```yaml
iam:
  podIdentityAssociations:
  - namespace: <string> #required
    serviceAccountName: <string> #required
    createServiceAccount: true #optional, default is false
    roleARN: <string> #required if none of permissionPolicyARNs, permissionPolicy and wellKnownPolicies is specified. Also, cannot be used together with any of the three other referenced fields.
    roleName: <string> #optional, generated automatically if not provided, ignored if roleARN is provided
    permissionPolicy: {} #optional
    permissionPolicyARNs: [] #optional
    wellKnownPolicies: {} #optional
    permissionsBoundaryARN: <string> #optional
    tags: {} #optional
```

For a complete example, refer to [pod-identity-associations.yaml](https://github.com/eksctl-io/eksctl/blob/main/examples/39-pod-identity-association.yaml).

???+ note
    Apart from `permissionPolicy` which is used as an inline policy document, all other fields have a CLI flag counterpart.

Creating pod identity associations can be achieved in the following ways. During cluster creation, by specifying the desired pod identity associations as part of the config file and running:

```
eksctl create cluster -f config.yaml
```

Post cluster creation, using either a config file e.g.

```
eksctl create podidentityassociation -f config.yaml
```

OR using CLI flags e.g.

```bash
eksctl create podidentityassociation \
    --cluster my-cluster \
    --namespace default \
    --service-account-name s3-reader \
    --permission-policy-arns="arn:aws:iam::111122223333:policy/permission-policy-1, arn:aws:iam::111122223333:policy/permission-policy-2" \
    --well-known-policies="autoScaler,externalDNS" \
    --permissions-boundary-arn arn:aws:iam::111122223333:policy/permissions-boundary
```

???+ note
    Only a single IAM role can be associated with a service account at a time. Therefore, trying to create a second pod identity association for the same service account will result in an error.

## Fetching Pod Identity Associations

To retrieve all pod identity associations for a certain cluster, run one of the following commands:

```
eksctl get podidentityassociation -f config.yaml
```

OR

```
eksctl get podidentityassociation --cluster my-cluster
```

Additionally, to retrieve only the pod identity associations within a given namespace, use the `--namespace` flag, e.g.

```
eksctl get podidentityassociation --cluster my-cluster --namespace default
```

Finally, to retrieve a single association, corresponding to a certain K8s service account, also include the `--service-account-name` to the command above, i.e.

```
eksctl get podidentityassociation --cluster my-cluster --namespace default --service-account-name s3-reader
```

## Updating Pod Identity Associations

To update the IAM role of one or more pod identity associations, either pass the new `roleARN(s)` to the config file e.g.

```yaml
iam:
  podIdentityAssociations:
    - namespace: default
      serviceAccountName: s3-reader
      roleARN: new-role-arn-1
    - namespace: dev
      serviceAccountName: app-cache-access
      roleARN: new-role-arn-2
```

and run:

```
eksctl update podidentityassociation -f config.yaml
```

OR (to update a single association) pass the new `--role-arn` via CLI flags:

```
eksctl update podidentityassociation --cluster my-cluster --namespace default --service-account-name s3-reader --role-arn new-role-arn
```

## Deleting Pod Identity Associations

To delete one or more pod identity associations, either pass `namespace(s)` and `serviceAccountName(s)` to the config file e.g.

```yaml
iam:
  podIdentityAssociations:
    - namespace: default
      serviceAccountName: s3-reader
    - namespace: dev
      serviceAccountName: app-cache-access
```

and run:

```
eksctl delete podidentityassociation -f config.yaml
```

OR (to delete a single association) pass the `--namespace` and `--service-account-name` via CLI flags:

```
eksctl delete podidentityassociation --cluster my-cluster --namespace default --service-account-name s3-reader
```

## EKS Add-ons support for pod identity associations

EKS Add-ons also support receiving IAM permissions via EKS Pod Identity Associations. The config file exposes three fields that allow configuring these: `addon.podIdentityAssociations`, `addonsConfig.autoApplyPodIdentityAssociations` and `addon.useDefaultPodIdentityAssociations`. You can either explicitly configure the desired pod identity associations, using `addon.podIdentityAssociations`, or have `eksctl` automatically resolve (and apply) the recommended pod identity configuration, using either `addonsConfig.autoApplyPodIdentityAssociations` or `addon.useDefaultPodIdentityAssociations`.

???+ note
    Not all EKS Add-ons will support pod identity associations at launch. For this case, required IAM permissions shall continue to be provided using [IRSA settings](/usage/addons/#creating-addons-and-providing-iam-permissions-via-irsa).

### Creating addons with IAM permissions

When creating an addon that requires IAM permissions, `eksctl` will first check if either pod identity associations or IRSA settings are being explicitly configured as part of the config file, and if so, use one of those to configure the permissions for the addon. e.g.

```yaml
addons:
- name: vpc-cni
  podIdentityAssociations:
  - serviceAccountName: aws-node
    permissionPolicyARNs: ["arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"]
```

and run

```bash
eksctl create addon -f config.yaml
2024-05-13 15:38:58 [ℹ] pod identity associations are set for "vpc-cni" addon; will use these to configure required IAM permissions
```

???+ note
    Setting both pod identities and IRSA at the same time is not allowed, and will result in a validation error.

For EKS Add-ons that support pod identities, `eksctl` offers the option to automatically configure any recommended IAM permissions, on addon creation. This can be achieved by simply setting `addonsConfig.autoApplyPodIdentityAssociations: true` in the config file. e.g.

```yaml
addonsConfig:
  autoApplyPodIdentityAssociations: true
# bear in mind that if either pod identity or IRSA configuration is explicitly set in the config file,
# or if the addon does not support pod identities,
# addonsConfig.autoApplyPodIdentityAssociations won't have any effect.
addons:
- name: vpc-cni
```

and run

```bash
eksctl create addon -f config.yaml
2024-05-13 15:38:58 [ℹ] "addonsConfig.autoApplyPodIdentityAssociations" is set to true; will lookup recommended pod identity configuration for "vpc-cni" addon
```

Equivalently, the same can be done via CLI flags e.g.

```bash
eksctl create addon --cluster my-cluster --name vpc-cni --auto-apply-pod-identity-associations
```

To migrate an existing addon to use pod identity with the recommended IAM policies, use

```yaml
addons:
- name: vpc-cni
  useDefaultPodIdentityAssociations: true
```

```bash
$ eksctl update addon -f config.yaml
```

### Updating addons with IAM permissions

When updating an addon, specifying `addon.PodIdentityAssociations` will represent the single source of truth for the state that the addon shall have, after the update operation is completed. Behind the scenes, different types of operations are performed in order to achieve the desired state i.e.

- create pod identites that are present in the config file, but missing on the cluster
- delete existing pod identites that were removed from the config file, together with any associated IAM resources
- update existing pod identities that are also present in the config file, and for which the set of IAM permissions has changed

???+ note
    The lifecycle of pod identity associations owned by EKS Add-ons is directly handled by the EKS Addons API, thus, using `eksctl update podidentityassociation` (to update IAM permissions) or `eksctl delete podidentityassociations` (to remove the association) is not supported for this type of associations. Instead, `eksctl update addon` or `eksctl delete addon` shall be used.

Let's see an example for the above, starting by analyzing the initial pod identity config for the addon:

```bash
eksctl get podidentityassociation --cluster my-cluster --namespace opentelemetry-operator-system --output json
[
    {
        ...
        "ServiceAccountName": "adot-col-prom-metrics",
        "RoleARN": "arn:aws:iam::111122223333:role/eksctl-my-cluster-addon-adot-podident-Role1-JwrGA4mn1Ny8",
        # OwnerARN is populated when the pod identity lifecycle is handled by the EKS Addons API
        "OwnerARN": "arn:aws:eks:us-west-2:111122223333:addon/my-cluster/adot/b2c7bb45-4090-bf34-ec78-a2298b8643f6"
    },
    {
        ...
        "ServiceAccountName": "adot-col-otlp-ingest",
        "RoleARN": "arn:aws:iam::111122223333:role/eksctl-my-cluster-addon-adot-podident-Role1-Xc7qVg5fgCqr",
        "OwnerARN": "arn:aws:eks:us-west-2:111122223333:addon/my-cluster/adot/b2c7bb45-4090-bf34-ec78-a2298b8643f6"
    }
]
```

Now use the below configuration:

```yaml
addons:
- name: adot
  podIdentityAssociations:

  # For the first association, the permissions policy of the role will be updated
  - serviceAccountName: adot-col-prom-metrics
    permissionPolicyARNs:
    #- arn:aws:iam::aws:policy/AmazonPrometheusRemoteWriteAccess
    - arn:aws:iam::aws:policy/CloudWatchAgentServerPolicy

  # The second association will be deleted, as it's been removed from the config file
  #- serviceAccountName: adot-col-otlp-ingest
  #  permissionPolicyARNs:
  #  - arn:aws:iam::aws:policy/AWSXrayWriteOnlyAccess

  # The third association will be created, as it's been added to the config file
  - serviceAccountName: adot-col-container-logs
    permissionPolicyARNs:
    - arn:aws:iam::aws:policy/CloudWatchAgentServerPolicy
```

and run

```bash
eksctl update addon -f config.yaml
...
# updating the permission policy for the first association
2024-05-14 13:27:43 [ℹ]  updating IAM resources stack "eksctl-my-cluster-addon-adot-podidentityrole-adot-col-prom-metrics" for pod identity association "a-reaxk2uz1iknwazwj"
2024-05-14 13:27:44 [ℹ]  waiting for CloudFormation changeset "eksctl-opentelemetry-operator-system-adot-col-prom-metrics-update-1715682463" for stack "eksctl-my-cluster-addon-adot-podidentityrole-adot-col-prom-metrics"
2024-05-14 13:28:47 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-addon-adot-podidentityrole-adot-col-prom-metrics"
2024-05-14 13:28:47 [ℹ]  updated IAM resources stack "eksctl-my-cluster-addon-adot-podidentityrole-adot-col-prom-metrics" for "a-reaxk2uz1iknwazwj"
# creating the IAM role for the second association
2024-05-14 13:28:48 [ℹ]  deploying stack "eksctl-my-cluster-addon-adot-podidentityrole-adot-col-container-logs"
2024-05-14 13:28:48 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-addon-adot-podidentityrole-adot-col-container-logs"
2024-05-14 13:29:19 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-addon-adot-podidentityrole-adot-col-container-logs"
# updating the addon, which handles the pod identity config changes behind the scenes
2024-05-14 13:29:19 [ℹ]  updating addon
# deleting the IAM role for the third association
2024-05-14 13:29:19 [ℹ]  deleting IAM resources for pod identity service account adot-col-otlp-ingest
2024-05-14 13:29:20 [ℹ]  will delete stack "eksctl-my-cluster-addon-adot-podidentityrole-adot-col-otlp-ingest"
2024-05-14 13:29:20 [ℹ]  waiting for stack "eksctl-my-cluster-addon-adot-podidentityrole-adot-col-otlp-ingest" to get deleted
2024-05-14 13:29:51 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-addon-adot-podidentityrole-adot-col-otlp-ingest"
2024-05-14 13:29:51 [ℹ]  deleted IAM resources for addon adot
```

now check that pod identity config was updated correctly

```bash
eksctl get podidentityassociation --cluster my-cluster --output json
[
    {
        ...
        "ServiceAccountName": "adot-col-prom-metrics",
        "RoleARN": "arn:aws:iam::111122223333:role/eksctl-my-cluster-addon-adot-podident-Role1-nQAlp0KktS2A",
        "OwnerARN": "arn:aws:eks:us-west-2:111122223333:addon/my-cluster/adot/1ec7bb63-8c4e-ca0a-f947-310c4b55052e"
    },
    {
        ...
        "ServiceAccountName": "adot-col-otlp-ingest",
        "RoleARN": "arn:aws:iam::111122223333:role/eksctl-my-cluster-addon-adot-podident-Role1-1k1XhAdziGzX",
        "OwnerARN": "arn:aws:eks:us-west-2:111122223333:addon/my-cluster/adot/1ec7bb63-8c4e-ca0a-f947-310c4b55052e"
    }
]
```


To remove all pod identity associations from an addon, `addon.PodIdentityAssociations` must be explicitly set to `[]`, e.g.

```yaml
addons:
- name: vpc-cni
  # omitting the `podIdentityAssociations` field from the config file,
  # instead of explicitly setting it to [], will result in a validation error
  podIdentityAssociations: []
```

and run

```bash
eksctl update addon -f config.yaml
```

### Deleting addons with IAM permissions

Deleting an addon will also remove all pod identities associated with the addon. Deleting the cluster will achieve the same effect, for all addons. Any IAM roles for pod identities, created by `eksctl`, will be deleted as-well.

## Migrating existing iamserviceaccounts and addons to pod identity associations

There is an `eksctl` utils command for migrating existing IAM Roles for service accounts to pod identity associations, i.e.

```
eksctl utils migrate-to-pod-identity --cluster my-cluster --approve
```

Behind the scenes, the command will apply the following steps:

- install the `eks-pod-identity-agent` addon if not already active on the cluster
- identify all IAM Roles that are associated with iamserviceaccounts
- identify all IAM Roles that are associated with EKS addons that support pod identity associations
- update the IAM trust policy of all identified roles, with an additional trusted entity, pointing to the new EKS Service principal (and, optionally, remove exising OIDC provider trust relationship)
- create pod identity associations for filtered roles associated with iamserviceaccounts
- update EKS addons with pod identities (EKS API will create the pod identities behind the scenes)

Running the command without the `--approve` flag will only output a plan consisting of a set of tasks reflecting the steps above, e.g.

```bash
[ℹ]  (plan) would migrate 2 iamserviceaccount(s) and 2 addon(s) to pod identity association(s) by executing the following tasks
[ℹ]  (plan)

3 sequential tasks: { install eks-pod-identity-agent addon,
    ## tasks for migrating the addons
    2 parallel sub-tasks: {
        2 sequential sub-tasks: {
            update trust policy for owned role "eksctl-my-cluster--Role1-DDuMLoeZ8weD",
            migrate addon aws-ebs-csi-driver to pod identity,
        },
        2 sequential sub-tasks: {
            update trust policy for owned role "eksctl-my-cluster--Role1-xYiPFOVp1aeI",
            migrate addon vpc-cni to pod identity,
        },
    },
    ## tasks for migrating the iamserviceaccounts
    2 parallel sub-tasks: {
        2 sequential sub-tasks: {
            update trust policy for owned role "eksctl-my-cluster--Role1-QLXqHcq9O1AR",
            create pod identity association for service account "default/sa1",
        },
        2 sequential sub-tasks: {
            update trust policy for unowned role "Unowned-Role1",
            create pod identity association for service account "default/sa2",
        },
    }
}
[ℹ]  all tasks were skipped
[!]  no changes were applied, run again with '--approve' to apply the changes
```

The existing OIDC provider trust relationship is always being removed from IAM Roles associated with EKS Add-ons. Additionally, to remove the existing OIDC provider trust relationship from IAM Roles associated with iamserviceaccounts, run the command with `--remove-oidc-provider-trust-relationship` flag, e.g.

```
eksctl utils migrate-to-pod-identity --cluster my-cluster --approve --remove-oidc-provider-trust-relationship
```

## Further references

[Official AWS Userdocs for EKS Add-ons support for pod identities](https://docs.aws.amazon.com/eks/latest/userguide/add-ons-iam.html)

[Official AWS Blog Post on Pod Identity Associations](https://aws.amazon.com/blogs/aws/amazon-eks-pod-identity-simplifies-iam-permissions-for-applications-on-amazon-eks-clusters/)

[Official AWS userdocs for Pod Identity Associations](https://docs.aws.amazon.com/eks/latest/userguide/pod-identities.html)
