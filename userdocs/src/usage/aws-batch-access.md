# Enabling Access for AWS Batch

In order to allow [AWS Batch on Amazon EKS](https://docs.aws.amazon.com/batch/latest/userguide/eks.html)
to perform operations on the Kubernetes API, its SLR needs to be granted the required RBAC permissions.
eksctl provides a command that creates the required RBAC resources for AWS Batch, and updates the `aws-auth`
ConfigMap to bind the role with the SLR for AWS Batch.

```shell
$ eksctl create iamidentitymapping --cluster my-eks-cluster --service-name aws-batch --namespace my-batch-namespace
```

> NOTE: The Kubernetes namespace used by Batch, in this example `my-batch-namespace`, must exist before running this
> command to enable access.  See AWS Batch on Amazon EKS
> [getting started guide](https://docs.aws.amazon.com/batch/latest/userguide/getting-started-eks.html)
> for more information.
