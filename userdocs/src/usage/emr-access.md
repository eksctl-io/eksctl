# Enabling Access for Amazon EMR

In order to allow [EMR](https://aws.amazon.com/emr) to perform operations on the Kubernetes API, its SLR needs to be granted the required RBAC permissions.
eksctl provides a command that creates the required RBAC resources for EMR, and updates the `aws-auth` ConfigMap to bind
the role with the SLR for EMR.

```shell
$ eksctl create iamidentitymapping --cluster dev --service-name emr-containers --namespace default
```
