package drain

func (n *NodeGroupDrainer) SetDrainer(drainer Evictor) {
	n.evictor = drainer
}
