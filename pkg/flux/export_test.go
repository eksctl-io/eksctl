package flux

import "github.com/weaveworks/eksctl/pkg/executor"

func (c *Client) SetExecutor(executor executor.Executor) {
	c.executor = executor
}
