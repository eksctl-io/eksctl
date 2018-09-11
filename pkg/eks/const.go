package eks

// Node AMI "magic" values
const (
	// NodeAmiFixed is used to indicate that the fixed (i.e. compiled into eksctl) ami's should be used
	NodeAmiFixed = "fixed"

	// NodeAmiLatest is used to indicate that the latest EKS AMIs should be used for the nodes. This implies
	// that automatic resolution of AMI will occur.
	NodeAmiLatest = "latest"
)
