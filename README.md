# EKS Quickstart App Dev

This repo contains an initial set of cluster components to be installed and
configured by [eksctl](https://eksctl.io) through GitOps.

## Components

- ALB ingress controller -- to easily expose services to the World.
- [Cluster autoscaler](https://github.com/kubernetes/autoscaler/tree/master/cluster-autoscaler) -- to [automatically add/remove nodes](https://aws.amazon.com/premiumsupport/knowledge-center/eks-cluster-autoscaler-setup/) to/from your cluster based on its usage.
- [Prometheus](https://prometheus.io/) (its [Alertmanager](https://prometheus.io/docs/alerting/alertmanager/), its [operator](https://github.com/coreos/prometheus-operator), its [`node-exporter`](https://github.com/prometheus/node_exporter), [`kube-state-metrics`](https://github.com/kubernetes/kube-state-metrics), and [`metrics-server`](https://github.com/kubernetes-incubator/metrics-server)) -- for powerful metrics & alerts.
- [Grafana](https://grafana.com) -- for a rich way to visualize metrics via dashboards you can create, explore, and share.
- [Kubernetes dashboard](https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/) -- Kubernetes' standard dashboard.
- [Fluentd](https://www.fluentd.org/) & Amazon's [CloudWatch agent](https://aws.amazon.com/cloudwatch/) -- for cluster & containers' [log collection, aggregation & analytics in CloudWatch](https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/Container-Insights-setup-logs.html).
- [podinfo](https://github.com/stefanprodan/podinfo) --  a toy demo application.

## Pre-requisites

A running EKS cluster with [IAM policies](https://eksctl.io/usage/iam-policies/) for:

- ALB ingress
- auto-scaler
- CloudWatch

[Here](https://github.com/weaveworks/eksctl/blob/master/examples/eks-quickstart-app-dev.yaml) is a sample `ClusterConfig` manifest that shows how to enable these policies.

**N.B.**: policies are configured at node group level.
Therefore, depending on your use-case, you may want to:

- add these policies to all node groups,
- add [node selectors](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/) to the ALB ingress, auto-scaler and CloudWatch pods, so that they are deployed on the nodes configured with these policies.

## How to access workloads

For security reasons, this quickstart profile does not expose any workload publicly. However, should you want to access one of the workloads, various solutions are possible.

### Port-forwarding

You could port-forward into a pod, so that you (and _only_ you) could access it locally.

For example, for `demo/podinfo`:

- run:
    ```console
    kubectl --namespace demo port-forward service/podinfo 9898:9898
    ```
- go to http://localhost:9898

### Ingress

You could expose a service publicly, _at your own risks_, via ALB ingress.

**N.B.**: the ALB ingress controller requires services:

- to be of `NodePort` type,
- to have the following annotations:
    ```yaml
    annotations:
      kubernetes.io/ingress.class: alb
      alb.ingress.kubernetes.io/scheme: internet-facing
    ```

#### `NodePort` services

For any `NodePort` service:

```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: ${name}
  namespace: ${namespace}
  annotations:
    kubernetes.io/ingress.class: alb
    alb.ingress.kubernetes.io/scheme: internet-facing
  labels:
    app: ${service-app-selector}
spec:
  rules:
    - http:
        paths:
          - path: /*
            backend:
              serviceName: ${service-name}
              servicePort: 80
```

A few minutes after deploying the above `Ingress` object, you should be able to see the public URL for the service:
```console
$ kubectl get ingress --namespace demo podinfo
NAME      HOSTS   ADDRESS                                                                     PORTS   AGE
podinfo   *       xxxxxxxx-${namespace}-${name}-xxxx-xxxxxxxxxx.${region}.elb.amazonaws.com   80      1s
```

#### `HelmRelease` objects

For `HelmRelease` objects, you would have to configure `spec.values.service` and `spec.values.ingress`, e.g. for `demo/podinfo`:

```yaml
apiVersion: helm.fluxcd.io/v1
kind: HelmRelease
metadata:
  name: podinfo
  namespace: demo
spec:
  releaseName: podinfo
  chart:
    git: https://github.com/stefanprodan/podinfo
    ref: 3.0.0
    path: charts/podinfo
  values:
    service:
      enabled: true
      type: NodePort
    ingress:
      enabled: true
      annotations:
        kubernetes.io/ingress.class: alb
        alb.ingress.kubernetes.io/scheme: internet-facing
      path: /*
```

**N.B.**: the above `HelmRelease`

- changes the type of `podinfo`'s service from its default value, `ClusterIP`, to `NodePort`,
- adds the annotations required for the ALB ingress controller to expose the service, and
- exposes all of `podinfo`'s URLs, so that all assets can be served over HTTP.

A few minutes after deploying the above `HelmRelease` object, you should be able to see the following `Ingress` object, and the public URL for `podinfo`:

```console
$ kubectl get ingress --namespace demo podinfo
NAME      HOSTS   ADDRESS                                                             PORTS   AGE
podinfo   *       xxxxxxxx-demo-podinfo-xxxx-xxxxxxxxxx.${region}.elb.amazonaws.com   80      1s
```

## Securing your endpoints
For a production-grade deployment, it's recommended to secure your endpoints with SSL. See [Ingress annotations for SSL](https://kubernetes-sigs.github.io/aws-alb-ingress-controller/guide/ingress/annotation/#ssl).

Any sensitive service that needs to be exposed must have some form of authentication. To add authentication to Grafana for e.g., see [Grafana configuration](https://github.com/helm/charts/tree/master/stable/prometheus-operator#grafana).
To add authentication to other components, please consult their documentation.

## Get in touch

[Create an issue](https://github.com/weaveworks/eks-quickstart-app-dev/issues/new), or
login to [Weave Community Slack (#eksctl)][slackchan] ([signup][slackjoin]).

[slackjoin]: https://slack.weave.works/
[slackchan]: https://weave-community.slack.com/messages/eksctl/
