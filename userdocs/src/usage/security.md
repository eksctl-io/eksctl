# Security

`eksctl` provides some options that can improve the security of your EKS cluster.

## `withOIDC`

Enable [`withOIDC`](/usage/schema/#iam-withOIDC) to automatically create an [IRSA](/usage/iamserviceaccounts/) for the amazon CNI plugin and
limit permissions granted to nodes in your cluster, instead granting the necessary permissions
only to the CNI service account. The background is described in [this AWS
documentation](https://docs.aws.amazon.com/eks/latest/userguide/cni-iam-role.html).

## `disablePodIMDS`

For managed and unmanaged nodegroups, [`disablePodIMDS`](/usage/schema/#nodeGroups-disablePodIMDS) option is available prevents all
non host networking pods running in this nodegroup from making IMDS requests.

???+ note
    This can not be used together with [`withAddonPolicies`](/usage/iam-policies/).

