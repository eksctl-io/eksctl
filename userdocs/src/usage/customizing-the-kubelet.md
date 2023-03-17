# Customizing kubelet configuration

## Customizing kubelet configuration

System resources can be reserved through the configuration of the kubelet. This is recommended, because in the case
of resource starvation the kubelet might not be able to evict pods and eventually make the node become `NotReady`. To
 do this, config files can include the `kubeletExtraConfig` field which accepts a free form yaml that will be embedded
 into the `kubelet.yaml`.


Some fields in the `kubelet.yaml` are set by eksctl and therefore are not overwritable, such as the `address`,
`clusterDomain`, `authentication`, `authorization`, or `serverTLSBootstrap`.

The following example config file creates a nodegroup that reserves `300m` vCPU, `300Mi` of memory and `1Gi` of
ephemeral-storage for the kubelet; `300m` vCPU, `300Mi` of memory and `1Gi`of ephemeral storage for OS system
daemons; and kicks in eviction of pods when there is less than `200Mi` of memory available or less than  10% of the
root filesystem.

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: dev-cluster-1
  region: eu-north-1

nodeGroups:
  - name: ng-1
    instanceType: m5a.xlarge
    desiredCapacity: 1
    kubeletExtraConfig:
        kubeReserved:
            cpu: "300m"
            memory: "300Mi"
            ephemeral-storage: "1Gi"
        kubeReservedCgroup: "/kube-reserved"
        systemReserved:
            cpu: "300m"
            memory: "300Mi"
            ephemeral-storage: "1Gi"
        evictionHard:
            memory.available:  "200Mi"
            nodefs.available: "10%"
        featureGates:
            RotateKubeletServerCertificate: true # has to be enabled, otherwise it will be disabled
```

In this example, given instances of type `m5a.xlarge` which have 4 vCPUs and 16GiB of memory, the `Allocatable` amount
of CPUs would be 3.4 and 15.4 GiB of memory. It is important to know that the values specified in the config file for 
the the fields in `kubeletExtraconfig` will completely overwrite the default values specified by eksctl. 
However, omitting one or more `kubeReserved` parameters will cause the missing parameters to be defaulted to sane 
values based on the aws instance type being used.

### A note on the `kubeReserved` calculation for NodeGroups with mixed instances

While it is generally recommended to configure a mixed instance NodeGroup to use instances with the same CPU and RAM
configuration; that's not a strict requirement. Therefore the `kubeReserved` calculation uses the _smallest instance_ in
the `InstanceDistribution.InstanceTypes` field. This way NodeGroups with disparate instance types will not reserve too
many resources on the smallest instance. However, this could lead to a reservation that is too small for the largest
instance type.

!!! warning
    By default `eksctl` sets `featureGates.RotateKubeletServerCertificate=true`, but when custom `featureGates` are
    provided, it will be unset. You should always include `featureGates.RotateKubeletServerCertificate=true`, unless
    you have to disable it.

