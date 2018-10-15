package manager

// DeprecatedDeleteStackVPC deletes the VPC stack
func (c *StackCollection) DeprecatedDeleteStackVPC() error {
	_, err := c.DeleteStack("EKS-" + c.spec.ClusterName + "-VPC")
	return err
}

// DeprecatedDeleteStackServiceRole deletes the service role stack
func (c *StackCollection) DeprecatedDeleteStackServiceRole() error {
	_, err := c.DeleteStack("EKS-" + c.spec.ClusterName + "-ServiceRole")
	return err
}

// DeprecatedDeleteStackDefaultNodeGroup deletes the default node group stack
func (c *StackCollection) DeprecatedDeleteStackDefaultNodeGroup() error {
	_, err := c.DeleteStack("EKS-" + c.spec.ClusterName + "-DefaultNodeGroup")
	return err
}

// DeprecatedDeleteStackControlPlane deletes the control plane stack
func (c *StackCollection) DeprecatedDeleteStackControlPlane() error {
	_, err := c.DeleteStack("EKS-" + c.spec.ClusterName + "-ControlPlane")
	return err
}
