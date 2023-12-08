# EKS Pod Identity Associations

## Introduction

AWS EKS has introduced a new enhanced mechanism called Pod Identity Association for cluster administrators to configure Kubernetes applications to receive IAM permissions required to connect with AWS services outside of the cluster. Pod Identity Association leverages IRSA however, it makes it configurable directly through EKS API, eliminating the need for using IAM API altogether. 

As a result, IAM roles no longer need to reference an [OIDC provider](/usage/iamserviceaccounts/#how-it-works) and hence won't be tied to a single cluster anymore. This means, IAM roles can now be used across multiple EKS clusters without the need to update the role trust policy each time a new cluster is created. This in turn, eliminates the need for role duplication and simplifies the process of automating IRSA altogether.

## Prerequisites

Behind the scenes, the implementation of pod identity associations is running an agent as a daemonset on the worker nodes. To run the pre-requisite agent on the cluster, EKS provides a new add-on called EKS Pod Identity Agent. Therefore, creating pod identity associations (with `eksctl`) requires the `eks-pod-identity-agent` addon pre-installed on the cluster. This addon can be [created using `eksctl`](/usage/addons/#creating-addons) in the same fashion any other supported addon is, e.g.

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

## Migrating existing iamserviceaccounts to pod identity associations

`eksctl` has introduced a new utils command for migrating existing IAM Roles for service accounts to pod identity associations, i.e.

```
eksctl utils migrate-to-pod-identity --cluster my-cluster --approve
```

Behind the scenes, the command will apply the following steps:

- install the `eks-pod-identity-agent` addon if not already active on the cluster
- identify all IAM Roles that are associated with K8s service accounts
- update the IAM trust policy of all roles, with an additional trusted entity, pointing to the new EKS Service principal (and, optionally, remove exising OIDC provider trust relationship)
- create pod identity associations between all identified roles and the respective service accounts 

Running the command without the `--approve` flag will only output a plan consisting of a set of tasks reflecting the steps above, e.g. 

```bash
[ℹ]  (plan) would migrate 2 iamserviceaccount(s) to pod identity association(s) by executing the following tasks
[ℹ]  (plan) 
3 sequential tasks: { install eks-pod-identity-agent addon, 
    2 parallel sub-tasks: { 
        update trust policy for owned role "eksctl-my-cluster-addon-iamserv-Role1-beYhlhzpwQte",
        update trust policy for unowned role "Unowned-Role1",
    }, 
    2 parallel sub-tasks: { 
        create pod identity association for service account "default/sa1",
        create pod identity association for service account "default/sa2",
    } 
}
[ℹ]  all tasks were skipped
[!]  no changes were applied, run again with '--approve' to apply the changes
```

Additionally, to delete the existing OIDC provider trust relationship from all IAM Roles, run the command with `--remove-oidc-provider-trust-relationship` flag, e.g.

```
eksctl utils migrate-to-pod-identity --cluster my-cluster --approve --remove-oidc-provider-trust-relationship
```


## Further references

[Official AWS Blog Post](https://aws.amazon.com/blogs/aws/amazon-eks-pod-identity-simplifies-iam-permissions-for-applications-on-amazon-eks-clusters/)

[Official AWS userdocs](https://docs.aws.amazon.com/eks/latest/userguide/pod-identities.html)