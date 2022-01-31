# Additional Volume Mappings

As an additional configuration option, when dealing with volume mappings, it's possible to configure extra mappings
when the nodegroup is created.

To do this, set the following field:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: dev-cluster
  region: eu-north-1

managedNodeGroups:
  - name: ng-1-workers
    labels: { role: workers }
    instanceType: m5.xlarge
    desiredCapacity: 10
    volumeSize: 80
    additionalVolumes:
      volumeSize: 80
      volumeName: '/tmp/mount-1'
      volumeType: 'gp3'
      volumeEncrypted: true
      volumeKmsKeyID: 'id'
      volumeIOPS: 3000
      volumeThroughput: 125
```
