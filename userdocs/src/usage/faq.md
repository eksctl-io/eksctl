## Nodegroups

!!! question "Question"
    How can I change the instance type of my nodegroup?

!!! quote "Answer"
    From the point of view of `eksctl`, nodegroups are immutable. This means that once created the only thing `eksctl` can do is scale the nodegroup up or down.

    To change the instance type, create a new nodegroup with the desired instance type, then drain it so that the workloads move to the new one. After that step is complete you can delete the old nodegroup

## Ingress

!!! question "Question"
    How do I set up ingress with `eksctl`?
!!! quote "Answer"
    If the plan is to use AWS ALB Ingress controller, setting `nodegroups[*].iam.withAddonPolicies.albIngress` to `true` will add the required IAM policies to your nodes allowing the controller to provision load balancers. Then you can follow [docs to set up the controller](https://kubernetes-sigs.github.io/aws-alb-ingress-controller/guide/controller/setup/).

    For Nginx Ingress Controller, setup would be the same as [any other Kubernetes cluster](https://kubernetes.github.io/ingress-nginx/deploy/#aws).
