package manager

// DeprecatedDeleteStackVPC deletes the VPC stack
func (c *StackCollection) DeprecatedDeleteStackVPC(wait bool) error {
	var err error
	stackName := "EKS-" + c.spec.Metadata.Name + "-VPC"

	if wait {
		err = c.WaitDeleteStack(stackName)
	} else {
		_, err = c.DeleteStack(stackName)
	}

	return err
}

// DeprecatedDeleteStackServiceRole deletes the service role stack
func (c *StackCollection) DeprecatedDeleteStackServiceRole(wait bool) error {
	var err error
	stackName := "EKS-" + c.spec.Metadata.Name + "-ServiceRole"

	if wait {
		err = c.WaitDeleteStack(stackName)
	} else {
		_, err = c.DeleteStack(stackName)
	}

	return err
}

// DeprecatedDeleteStackDefaultNodeGroup deletes the default node group stack
func (c *StackCollection) DeprecatedDeleteStackDefaultNodeGroup(wait bool) error {
	var err error
	stackName := "EKS-" + c.spec.Metadata.Name + "-DefaultNodeGroup"

	if wait {
		err = c.WaitDeleteStack(stackName)
	} else {
		_, err = c.DeleteStack(stackName)
	}

	return err
}

// DeprecatedDeleteStackControlPlane deletes the control plane stack
func (c *StackCollection) DeprecatedDeleteStackControlPlane(wait bool) error {
	var err error
	stackName := "EKS-" + c.spec.Metadata.Name + "-ControlPlane"

	if wait {
		err = c.WaitDeleteStack(stackName)
	} else {
		_, err = c.DeleteStack(stackName)
	}

	return err
}
