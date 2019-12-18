# GPU Support

!!!warning
    If you'd like to use GPU instance types, that is, [p2](https://aws.amazon.com/ec2/instance-types/p2/) or [p3](https://aws.amazon.com/ec2/instance-types/p3/) types, then the first thing you need to do is subscribe to the [EKS-optimized AMI with GPU Support](https://aws.amazon.com/marketplace/pp/B07GRHFXGM). If you don't do this then node creation will fail.

After subscribing to the AMI you can create a cluster specifying the GPU instance type you'd like to use for the nodes. For example:

```
eksctl create cluster --node-type=p2.xlarge
```

The AMI resolvers (`static`, `auto` and `auto-ssm`) will see that you want to use a GPU instance type (p2 or p3 only) and they will select the correct AMI.

Once the cluster is created you will need to install the [NVIDIA Kubernetes device plugin](https://github.com/NVIDIA/k8s-device-plugin). Check the repo for the most up to date instructions but you should be able to run this:

```
kubectl create -f https://raw.githubusercontent.com/NVIDIA/k8s-device-plugin/v1.11/nvidia-device-plugin.yml
```

!!!note
    Once `addon` support has been added as part of 0.2.0 it is envisioned that there will be a addon to install the NVIDIA Kubernetes Device Plugin. This addon could potentially be installed automatically as we know an GPU instance type is being used.
