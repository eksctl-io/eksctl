---
title: FAQ
weight: 220
---

## FAQ

### How can I change the instance type of my nodegroup?

From the point of view of `eksctl`, nodegroups are immutable. This means that once created the only thing `eksctl` 
can do is scale the nodegroup up or down.

To change the instance type, create a new nodegroup with the desired instance type, then drain it so that the 
workloads move to the new one. After that step is complete you can delete the old nodegroup.

