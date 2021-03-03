package manager

import "fmt"

type StackNotFoundErr struct {
	ClusterName string
}

func (e *StackNotFoundErr) Error() string {
	return fmt.Sprintf("no eksctl-managed CloudFormation stacks found for %q", e.ClusterName)
}
