# Connection tracking timeouts

Security groups track each connection to and from a node so that return traffic is allowed
automatically. Tracked connections that sit idle for longer than EC2's timeout are removed
from the connection tracking table, and traffic on them is dropped.

The default idle timeout for established TCP connections depends on the instance's Nitro
generation. Nitro v6 instance types (the `*8i` families, excluding P6e-GB200) use 350
seconds, while earlier generations use 432,000 seconds. Workloads that hold long-lived but
mostly idle TCP connections - WebSocket streams, gRPC channels, database connection pools -
can therefore start losing connections after moving to a newer instance type.

`connectionTracking` sets the timeouts on the network interfaces of a nodegroup's nodes, so
you can keep the longer timeout, or pick a shorter one:

```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: connection-tracking
  region: us-west-2

managedNodeGroups:
  - name: ng-1
    instanceType: c8in.8xlarge
    connectionTracking:
      tcpEstablishedTimeout: 432000

nodeGroups:
  - name: ng-2
    instanceType: m8i.4xlarge
    connectionTracking:
      tcpEstablishedTimeout: 3600
      udpStreamTimeout: 180
      udpTimeout: 60
```

At least one timeout must be set. Timeouts you leave out keep EC2's default for the
instance type.

| Field | Description | Range (seconds) |
|-------|-------------|-----------------|
| `tcpEstablishedTimeout` | Idle TCP connections in an established state | 60 - 432000 |
| `udpStreamTimeout` | Idle UDP flows classified as streams, which have seen more than one request-response transaction | 60 - 180 |
| `udpTimeout` | Idle UDP flows that have seen traffic only in a single direction or a single request-response transaction | 30 - 60 |

The timeouts apply to the network interfaces `eksctl` puts in the launch template it
generates, which covers the node's primary interface and, for EFA-enabled nodegroups, the
EFA interfaces. The Amazon VPC CNI plugin copies the primary interface's settings onto the
interfaces it creates for pods, so pod traffic inherits them too. That copying was added in
VPC CNI `v1.21.2` and `v1.22.1`; on older versions the timeouts apply to node traffic only.

Because `connectionTracking` is applied through the launch template that `eksctl` generates,
it cannot be combined with a nodegroup that supplies its own launch template through
`launchTemplate.id`. Set the timeouts in your own launch template in that case.

!!! note
    Changing `connectionTracking` on an existing nodegroup does not update the nodes that
    are already running. Create a new nodegroup, or replace the existing nodes, to pick up
    the new timeouts.

## Further information

- [Connection tracking][ec2-conntrack] in the Amazon EC2 User Guide
- [Nitro instance types][ec2-nitro] and their default timeouts

[ec2-conntrack]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/security-group-connection-tracking.html
[ec2-nitro]: https://docs.aws.amazon.com/ec2/latest/instancetypes/ec2-nitro-instances.html
