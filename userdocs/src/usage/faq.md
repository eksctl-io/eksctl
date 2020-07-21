## Nodegroups

!!! question "How can I change the instance type of my nodegroup?"
    From the point of view of `eksctl`, nodegroups are immutable. This means that once created the only thing `eksctl` can do is scale the nodegroup up or down.

    To change the instance type, create a new nodegroup with the desired instance type, then drain it so that the workloads move to the new one. After that step is complete you can delete the old nodegroup

!!! question "How can I see the generated userdata for a nodegroup?"
    First you'll need the name of the Cloudformation stack that manages the
    nodegroup:
    ```console
    $ eksctl utils describe-stacks --region=us-west-2 --cluster NAME
    ```
    You'll see a name similar to `eksctl-CLUSTER_NAME-nodegroup-NODEGROUP_NAME`.

    You can execute the following to get the userdata. Note the final line which
    decodes from base64 and uncompresses the gzipped data.
    ```bash
    NG_STACK=eksctl-scrumptious-monster-1595247364-nodegroup-ng-29b8862f # your stack here
    LAUNCH_TEMPLATE_ID=$(aws cloudformation describe-stack-resources --stack-name $NG_STACK \
    | jq -r '.StackResources | map(select(.LogicalResourceId == "NodeGroupLaunchTemplate")
    | .PhysicalResourceId)[0]')
    aws ec2 describe-launch-template-versions --launch-template-id $LAUNCH_TEMPLATE_ID \
    | jq -r '.LaunchTemplateVersions[0].LaunchTemplateData.UserData' \
    | base64 -d | gunzip
    ```

## Ingress

!!! question "How do I set up ingress with `eksctl`?"
    If the plan is to use AWS ALB Ingress controller, setting `nodegroups[*].iam.withAddonPolicies.albIngress` to `true` will add the required IAM policies to your nodes allowing the controller to provision load balancers. Then you can follow [docs to set up the controller](https://kubernetes-sigs.github.io/aws-alb-ingress-controller/guide/controller/setup/).

    For Nginx Ingress Controller, setup would be the same as [any other Kubernetes cluster](https://kubernetes.github.io/ingress-nginx/deploy/#aws).

## Kubectl

!!! question "I'm using an HTTPS proxy and cluster certificate validation fails, how can I use the system CAs?"
    Set the environment variable `KUBECONFIG_USE_SYSTEM_CA` to make `kubeconfig`
    respect the system certificate authorities.
