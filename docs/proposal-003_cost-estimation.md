> **STATUS**: This proposal is a _working draft_, it will get refined and augment as needed.
> If any non-trivial changes are need to functionality defined here, in particular the user
> experience, those changes should be suggested via a PR to this proposal document.
> Any other changes to the text of the proposal or technical corrections are also very welcome.

The main challenge of the cost estimation is to calculate the price induced by a `cluster create` or a `cluster scale` (V2) command.
For the first version, we can focus on EC2 Instances type cost, EBS Volumes cost, EKS cluster cost per region.
For future versions, we can introduce new features like: EC2 Data Transfer, Elastic Load balancer, etc. 

**API Used**
Amazon provides AWS Price List API : https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/index.json. We can use it to get the current pricing of each service, in our case: EC2, EBS, EKS.

*********************
**Parameters**
`eksctl calculator` uses the same parameters as `eksctl create` and `eksctl scale` (V2) + specifics parameters.

**Existing parameters**

* __name__: ignored. No extra cost.
* __tags__: ignored. No extra cost.
* __zones__: ignored. No extra cost.
* __version__: ignored. No extra cost.
* __region__: Used to calculate the cost. The cost is per region
* __nodes__: Used to calculate the cost unless nodes-min is provided
* __nodes-min__: Used to calculate the cost. Minimum number of instances used
* __nodes-max__: ignored.
* __node-type__: Used to calculate the cost. The cost depends on the instance type
* __node-volume-size__: Used to calculate the cost. The cost depends on the volume size
* __node-volume-type__: Used to calculate the cost. The cost depends on the volume type
* __max-pods-per-node__: ignored. No extra cost.
* __node-ami__: ignored. (V2)
* __node-ami-family__: ignored. (V2)
* __ssh-access__: ignored
* __ssh-public-key__: ignored
* __node-private-networking__: ignored
* __node-security-groups__: ignored. No extra cost.
* __node-labels__: ignored. No extra cost.
* __node-zones__: ignored. No extra cost.
* __temp-node-role-policies__: ignored. No extra cost.
* __temp-node-role-name__: ignored. No extra cost.
* __asg-access__: ignored. No extra cost.
* __external-dns-access__: ignored. No extra cost.
* __full-ecr-access__: ignored. No extra cost.
* __storage-class__: Used to calculate the cost. The cost depends on the storage class
* __vpc-private-subnets__: ignored. (V2)
* __vpc-public-subnets__: ignored. (V2)
* __vpc-cidr__: ignored. (V2)
* __vpc-from-kops-cluster__: ignored. (V2) 

**Specifics parameters**

***Usage***

Used to estimate the usage of the cluster and its nodes. By default, it's a full utilization: `--usage-type=utilization --usage-value=100`.
* __usage-type__: String. Possible values: `demo`, `utilization`, `day`, `week`, `month`.
* __usage-value__: Integer. Min value is `1`. Max value depends on usage type, `24` for `demo` and `day`, `100` for `utilization`, `168` for `week`, `732` for `month`.

Examples:
- Demo ` --usage-type=demo --usage-value=3`
- % Utilized/Month ` --usage-type=utilization --usage-value=54` 
- Hours/Day  ` --usage-type=day --usage-value=24`
- Hours/Week ` --usage-type=week --usage-value=2`
- Hours/Month ` --usage-type=month --usage-value=5`

***Free Tier***
* __free-tier__: yes/no (yes by default) 

***Response example***
`eksctl calculator cluster --name=cluster-1 --region=eu-west-1 --nodes=2`

```
[ℹ]  using region eu-west-1
[ℹ]  using "ami-0ce0ec06e682ee10e" for nodes
[ℹ]  2 nodes m5.large
[ℹ]  free tier "yes"
[ℹ]  usage type "utilization"
[ℹ]  usage value "100%"
[✔]  estimation (per month): 
     - Amazon EC2 Service (US East (N. Virginia)): $ 142.56
     -- Compute: $ 140.56
     -- EBS Volumes: $ 2.00
     - EKS Cluster: $ 146.4
     total: $ 288.96 
```

***Usecase (V2)***

Aims at giving an estimation of a `Common Customer Samples`: free-website, bigdata app, machine learning app, serverless, etc. 
It generates a yml file which contains a recommended configuration.

For example: 

`eksctl calculator --usecase=free-websites --region=eu-west-1 --free-tier=yes`

`eksctl calculator --usecase=addons --region=eu-web-1 --free-tier=yes --addons jenkins-x`
