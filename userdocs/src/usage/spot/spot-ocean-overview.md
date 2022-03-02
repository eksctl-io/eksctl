# Spot Ocean

[Ocean](https://spot.io/products/ocean/) by [Spot](https://spot.io/) simplifies infrastructure management for Kubernetes.  With robust, container-driven infrastructure auto-scaling and intelligent right-sizing for container resource requirements, operations can literally "set and forget" the underlying cluster.

Ocean seamlessly integrates with your existing nodegroups, as a drop-in replacement for AWS Auto Scaling groups, and allows you to streamline and optimize the entire workflow, from initially creating your cluster to managing and optimizing it on an ongoing basis.

## Features

- **Simplify Cluster Management** —
  Ocean's Virtual Node Groups make it easy to run different infrastructure in a single cluster, which can span multiple AWS VPC availability zones and subnets for high-availability.

- **Container-Driven Autoscaling and Vertical Rightsizing** —
  Auto-detect your container infrastructure requirements so the appropriate instance size or type will always be available. Measure real-time CPU/Memory consumption of your Pods for ongoing resource optimization.

- **Cloud-Native Showback** —
  Gain a granular view of your cluster's cost breakdown (compute and storage) for each one of the cluster's resources such as Namespaces, Deployments, Daemon Sets, Jobs, and Pods.

- **Optimized Pricing and Utilization** —
  Ocean not only intelligently leverages Spot Instances and reserved capacity to reduce costs, but also eliminates underutilized instances with container-driven autoscaling and advanced bin-packing.

## Prerequisites

Make sure you have installed [eksctl](./spot-ocean-eksctl-install.md/#installation) for Spot Ocean and [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/).

## Documentation

If you're new to [Spot](https://spot.io/) and want to get started, please check out our [Getting Started](https://docs.spot.io/getting-started-with-spotinst/) guide, available on the [Spot Help Center](https://docs.spot.io/) website.

## Getting Help

Please use these community resources for getting help:

- Join our [Spot](https://spot.io/) community on [Slack](http://slack.spot.io/).
- Open a GitHub [issue](https://github.com/spotinst/weaveworks-eksctl/issues/new/choose/).
- Ask a question on [Stack Overflow](https://stackoverflow.com/) and tag it with `spot-ocean`.
- Also see [Spot Ocean Setup](./spot-ocean-setup.md/#setup)
