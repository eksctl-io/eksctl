package manager

func (c *StackCollection) DeprecatedDeleteStackVPC() error {
	return c.DeleteStack("EKS-" + c.spec.ClusterName + "-VPC")
}

func (c *StackCollection) DeprecatedDeleteStackServiceRole() error {
	return c.DeleteStack("EKS-" + c.spec.ClusterName + "-ServiceRole")
}

func (c *StackCollection) DeprecatedDeleteStackDefaultNodeGroup() error {
	return c.DeleteStack("EKS-" + c.spec.ClusterName + "-DefaultNodeGroup")
}

func (c *StackCollection) DeprecatedDeleteStackControlPlane() error {
	return c.DeleteStack("EKS-" + c.spec.ClusterName + "-ControlPlane")
}
