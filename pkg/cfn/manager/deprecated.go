package manager

func (c *StackCollection) DeprecatedDeleteStackVPC() error {
	_, err := c.DeleteStack("EKS-" + c.spec.ClusterName + "-VPC")
	return err
}

func (c *StackCollection) DeprecatedDeleteStackServiceRole() error {
	_, err := c.DeleteStack("EKS-" + c.spec.ClusterName + "-ServiceRole")
	return err
}

func (c *StackCollection) DeprecatedDeleteStackDefaultNodeGroup() error {
	_, err := c.DeleteStack("EKS-" + c.spec.ClusterName + "-DefaultNodeGroup")
	return err
}

func (c *StackCollection) DeprecatedDeleteStackControlPlane() error {
	_, err := c.DeleteStack("EKS-" + c.spec.ClusterName + "-ControlPlane")
	return err
}
