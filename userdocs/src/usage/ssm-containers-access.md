# Enabling Access for SSM containers

In order to allow SSM to perform operations on the Kubernetes API, its SLR needs to be granted the required RBAC permissions.
eksctl provides a command that creates the required RBAC resources for SSM containers, and updates the `aws-auth` ConfigMap to bind
the role with the SLR for SSM.

```shell
$ eksctl create iamidentitymapping --cluster dev --service-name ssm-containers
```
